package builder

import (
	"archive/tar"
	"bytes"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

func BuildRepo(dockerfilePath string) (string, error) {
	client, err := docker.NewClient("localhost:4342")
	if err != nil {
		return "", nil
	}

	t := time.Now()
	inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputbuf)
	tr.WriteHeader(&tar.Header{Name: dockerfilePath, Size: 10, ModTime: t, AccessTime: t, ChangeTime: t})
	// tr.Write([]byte("FROM base\n"))
	tr.Close()
	opts := docker.BuildImageOptions{
		Name:         "test",
		InputStream:  inputbuf,
		OutputStream: outputbuf,
	}
	if err := client.BuildImage(opts); err != nil {
		return "", err
	}

	return "testik", nil
}
