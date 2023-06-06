/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cartman0526/kubetool/pkg/config"
	"github.com/cartman0526/kubetool/pkg/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		pushImages()
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringVarP(&option.ConfigFile, "filename", "f", "", "specify the location of image list file")
	pushCmd.Flags().StringVarP(&option.Target, "tag", "t", "", "specify the image's tag")
}

func pushImages() {
	ctx := context.Background()
	vip := config.LoadConfig(option.ConfigFile)
	client := docker.NewDockerClient(vip.GetString("host"))
	defer client.Close()
	images := getImages(ctx, client)

	authConfig := types.AuthConfig{
		Username: vip.GetString("registry.username"),
		Password: vip.GetString("registry.password"),
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	var out io.ReadCloser
	for _, image := range images {
		imageResp := strings.Split(image, "/")

		switch len(imageResp) {
		case 3:
			//fmt.Println(imageResp)
			imageName := vip.GetString("registry.address") + "/" + option.Target + "/" + imageResp[2]
			client.ImageTag(ctx, image, imageName)
			out, err = client.ImagePush(ctx, imageName, types.ImagePushOptions{RegistryAuth: authStr})
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			defer out.Close()
			io.Copy(os.Stdout, out)
		case 2:
			imageName := vip.GetString("registry.address") + "/" + option.Target + "/" + imageResp[1]
			client.ImageTag(ctx, image, imageName)
			out, err = client.ImagePush(ctx, imageName, types.ImagePushOptions{RegistryAuth: authStr})
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			defer out.Close()
			io.Copy(os.Stdout, out)
		case 1:
			imageName := vip.GetString("registry.address") + "/" + option.Target + "/" + imageResp[0]
			client.ImageTag(ctx, image, imageName)
			out, err = client.ImagePush(ctx, imageName, types.ImagePushOptions{RegistryAuth: authStr})
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			defer out.Close()
			io.Copy(os.Stdout, out)
		}
	}
}

func getImages(ctx context.Context, client *client.Client) []string {
	var imageSlice []string
	imageList, err := client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, image := range imageList {
		imageSlice = append(imageSlice, image.RepoTags[0])
	}
	return imageSlice
}
