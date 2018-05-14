package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type Docker struct {
	Client *client.Client
}

func New() (*Docker, error) {
	client, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &Docker{Client: client}, nil
}

func (d *Docker) ImagePull(image string) error {
	out, err := d.Client.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	fmt.Println("Done")

	io.Copy(os.Stdout, out)
	return nil
}

func (d *Docker) ImageIsLoaded(image string) (bool, error) {
	result, err := d.Client.ImageSearch(context.Background(), image, types.ImageSearchOptions{Limit: 1})
	if err != nil {
		return false, err
	}

	if len(result) != 0 {
		return true, nil
	}
	return false, nil
}

func (d *Docker) CreateContainer(image string, env []string, port string, secret string) (string, error) {
	isLoaded, err := d.ImageIsLoaded(image)
	if err != nil {
		return "", err
	}

	if !isLoaded {
		if err := d.ImagePull(image); err != nil {
			return "", err
		}
	}

	//mysql -u root -pDE8002B5 -e 'status'
	portBindings := getPortBinding(port)
	resp, err := d.Client.ContainerCreate(context.Background(), &container.Config{
		Image: image,
		Env:   env,
		Healthcheck: &container.HealthConfig{
			Test:     []string{"CMD-SHELL", "mysql -u root -p" + secret + " -e 'status'"},
			Timeout:  60 * time.Second,
			Interval: 5 * time.Second,
			Retries:  10,
		},
	}, &container.HostConfig{PortBindings: portBindings}, nil, "")

	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (d *Docker) StartContainer(containerID string) error {
	if err := d.Client.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	status := ""
	attempts := 0
	for status != "healthy" && attempts <= 10 {
		r, _ := d.Client.ContainerInspect(context.Background(), containerID)
		status = r.State.Health.Status
		attempts = r.State.Health.FailingStreak
		time.Sleep(5 * time.Second)
	}

	if status == "unhealthy" || attempts == 10 {
		return errors.New("Cannot create container")
	}
	fmt.Println("Mysql now running")
	return nil
}

func (d *Docker) ExecCommand(containerID string, cmd []string) error {
	exceConfig := types.ExecConfig{
		User:   "root",
		Detach: false,
		Tty:    false,
		Cmd:    cmd,
	}
	exceResp, err := d.Client.ContainerExecCreate(context.Background(), containerID, exceConfig)
	if err != nil {
		return nil
	}

	sc := types.ExecStartCheck{
		Tty:    true,
		Detach: false,
	}
	if err := d.Client.ContainerExecStart(context.Background(), exceResp.ID, sc); err != nil {
		return err
	}

	info, err := d.Client.ContainerExecInspect(context.Background(), exceResp.ID)
	if err != nil {
		return err
	}
	if info.ExitCode == 125 || info.ExitCode == 126 || info.ExitCode == 127 {
		return errors.New("Cannot execute command")
	}

	return nil
}

func getPortBinding(port string) map[nat.Port][]nat.PortBinding {
	return map[nat.Port][]nat.PortBinding{
		"3306/tcp": []nat.PortBinding{nat.PortBinding{HostPort: port}}}
}
