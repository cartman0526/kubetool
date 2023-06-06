package docker

import (
	"github.com/docker/docker/client"
)

func NewDockerClient(host string) *client.Client {
	// ctx := context.Background()
	client, err := client.NewClientWithOpts(client.WithHost("tcp://"+host), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err.Error())
	}
	// defer client.Close()
	return client
}
