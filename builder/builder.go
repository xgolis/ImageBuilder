package builder

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
)

type DockerClient struct {
	Client    *client.Client
	Ctx       context.Context
	Path      string
	Username  string
	ImageName string
	ImageTag  string
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

	err = dockerClient.findImage()
	if err != nil {
		return "", fmt.Errorf("error while finding image: %v", err)
	}
	fmt.Printf("\nimage built: %s\n", dockerClient.ImageName)
	// fmt.Println("konec")

	err = dockerClient.pushImage()
	if err != nil {
		return "", fmt.Errorf("error while pushing image: %v", err)
	}
	fmt.Println("malo by to tam byt idk")
	return dockerClient.ImageName, err
}

func (c *DockerClient) buildImage() error {
	// source: https://stackoverflow.com/questions/46878793/golang-docker-api-reports-invalid-argument-while-hitting-imagebuild
	os.MkdirAll("container/"+c.Username, 0755)
	tar := new(archivex.TarFile)
	tar.Create("container/" + c.Username + "/conf.tar")
	tar.AddAll(c.Username, false)
	tar.Close()
	// doriesit vymazavanie

	dockerBuildContext, err := os.Open("container/" + c.Username + "/conf.tar")
	defer dockerBuildContext.Close()

	imageBuildResponse, err := c.Client.ImageBuild(
		c.Ctx,
		dockerBuildContext,
		types.ImageBuildOptions{
			Tags:    []string{"xgolis/" + c.Username},
			Context: dockerBuildContext,
			// ContextDir: "./" + c.Username,
			Dockerfile: "Dockerfile",
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

func (c *DockerClient) findImage() error {
	images, err := c.Client.ImageList(c.Ctx, types.ImageListOptions{
		All: true,
	})
	if err != nil {
		return err
	}

	for i := 0; i < len(images); i++ {
		if images[i].RepoTags[0] == "xgolis/"+c.Username+":latest" {
			c.ImageTag = images[i].RepoTags[0]
			c.ImageName = images[i].ID
			return nil
		}
	}
	return fmt.Errorf("did not find desired image")
}

func getAuthConfig() (*types.AuthConfig, error) {
	username := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")

	if username == "" || password == "" {
		return nil, fmt.Errorf("Missing docker credentials")
	}
	return &types.AuthConfig{
		Username: username,
		Password: password,
	}, nil
}

func (c *DockerClient) pushImage() error {
	authConfig, err := getAuthConfig()
	if err != nil {
		return err
	}

	authConfigEncoded, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(authConfigEncoded)

	body, err := c.Client.ImagePush(c.Ctx, c.ImageTag, types.ImagePushOptions{
		RegistryAuth: authStr,
	})
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}
	defer body.Close()

	return nil
}
