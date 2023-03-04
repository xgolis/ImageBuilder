package builder

import "github.com/docker/docker/client"

func BuildRepo(dockerfilePath string) (string, error) {
	// client, err := docker.NewClient("localhost:4342")
	// if err != nil { https://pkg.go.dev/github.com/docker/docker
	// 	return "", err
	// }

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// t := time.Now()
	// inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	// tr := tar.NewWriter(inputbuf)
	// tr.WriteHeader(&tar.Header{Name: dockerfilePath, Size: 10, ModTime: t, AccessTime: t, ChangeTime: t})
	// // tr.Write([]byte("FROM base\n"))
	// tr.Close()
	// opts := docker.BuildImageOptions{
	// 	Name:         "test",
	// 	InputStream:  inputbuf,
	// 	OutputStream: outputbuf,
	// }
	// if err := client.BuildImage(opts); err != nil {
	// 	return "", err
	// }

	return "testik", nil
}
