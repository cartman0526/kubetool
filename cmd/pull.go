/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cartman0526/kubetool/cmd/types"
	"github.com/cartman0526/kubetool/pkg/config"
	"github.com/cartman0526/kubetool/pkg/docker"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var option types.Option

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull image from outside registry or internet, -f to choose images list",
	Run: func(cmd *cobra.Command, args []string) {
		pullImages()
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.Flags().StringVarP(&option.ConfigFile, "filename", "f", "", "specify the location of image list file")

}

func pullImages() {
	ctx := context.Background()
	vip := config.LoadConfig(option.ConfigFile)
	client := docker.NewDockerClient(vip.GetString("host"))
	defer client.Close()

	authConfig := dockertypes.AuthConfig{
		Username: vip.GetString("registry.username"),
		Password: vip.GetString("registry.password"),
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	images := vip.GetStringSlice("images")
	for _, image := range images {
		fmt.Printf("开始拉取镜像: %s\n", image)
		out, err := client.ImagePull(ctx, image, dockertypes.ImagePullOptions{RegistryAuth: authStr})
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		defer out.Close()
		io.Copy(os.Stdout, out)
	}
}
