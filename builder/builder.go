package builder

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	Client   *client.Client
	Ctx      context.Context
	Path     string
	Username string
}

const DefaultDockerHost = "unix:///var/run/docker.sock"

func getDockerClient(client *client.Client, path string, username string) *DockerClient {
	return &DockerClient{
		Client:   client,
		Ctx:      context.Background(),
		Path:     path,
		Username: username,
	}
}

func BuildRepo(dockerfilePath, username string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err, " :unable to init client")
	}

	dockerClient := getDockerClient(cli, dockerfilePath, username)

	err = dockerClient.buildImage()
	if err != nil {
		return "", err
	}
	image, err := dockerClient.findImage()
	if err != nil {
		return "", fmt.Errorf("error while finding image: %v", err)
	}
	fmt.Printf("image built: %s", image)
	// fmt.Println("konec")

	return image, err
}

func (c *DockerClient) buildImage() error {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dockerFile := "Dockerfile"
	dockerFileReader, err := os.Open(c.Path)
	if err != nil {
		log.Fatal(err, " :unable to open Dockerfile")
	}
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		log.Fatal(err, " :unable to read dockerfile")
	}

	tarHeader := &tar.Header{
		Name: dockerFile,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		log.Fatal(err, " :unable to write tar header")
	}
	_, err = tw.Write(readDockerFile)
	if err != nil {
		log.Fatal(err, " :unable to write tar body")
	}
	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	imageBuildResponse, err := c.Client.ImageBuild(
		c.Ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Tags:       []string{c.Username},
			Context:    dockerFileTarReader,
			Dockerfile: dockerFile,
			Remove:     true})
	if err != nil {
		log.Fatal(err, " :unable to build docker image")
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}
	return nil
}

func (c *DockerClient) findImage() (string, error) {
	images, err := c.Client.ImageList(c.Ctx, types.ImageListOptions{
		All: true,
	})
	if err != nil {
		return "", err
	}

	for i := 0; i < len(images); i++ {
		if images[i].RepoTags[0] == c.Username+":latest" {
			return images[i].ID, nil
		}
	}
	return "", fmt.Errorf("did not find desired image")
}
