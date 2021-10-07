/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hysios/log"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	pb "cskyzn.com/pkg/bimgserver/rpc"
	"cskyzn.com/pkg/bimgserver/server"
)

var (
	cfgFile string
	addr    string
	client  bool
	target  string
	image   string
	method  string
	size    int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bimgserver",
	Short: "启动图像处理服务",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if !client {
			handler := pb.NewBimgServerServer(server.NewServer(":8090"))
			// You can use any mux you like - NewHelloWorldServer gives you an http.Handler.
			mux := http.NewServeMux()
			// The generated code includes a method, PathPrefix(), which
			// can be used to mount your service on a mux.
			mux.Handle(handler.PathPrefix(), handler)
			log.Infof("bimgserver start at addr %s", addr)
			http.ListenAndServe(addr, mux)
		} else {
			client := pb.NewBimgServerProtobufClient(target, &http.Client{})
			b, err := ioutil.ReadFile(image)
			if err != nil {
				log.Fatalf("open image %s failed %s", image, err)
			}
			switch method {
			case "thumbnail":
				resp, err := client.Thumbnail(context.Background(), &pb.ThumbnailReq{Content: b, Pixels: int32(size)})
				if err != nil {
					log.Fatalf("thumbnail error %s", err)
				}
				log.Info(ioutil.WriteFile(extractFile(image, ".out"), resp.Content, os.ModePerm))
			default:
				panic("nonimplement")
			}
		}
	},
}

func extractFile(filename string, add string) string {
	dir, file := filepath.Split(filename)
	base := filepath.Base(file)
	ext := filepath.Ext(file)
	return filepath.Join(dir, base+add+ext)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bimgserver.yaml)")
	rootCmd.PersistentFlags().StringVar(&addr, "addr", ":9080", "bimg server addr")
	rootCmd.PersistentFlags().BoolVar(&client, "client", false, "client mode")
	rootCmd.PersistentFlags().StringVar(&target, "target", "http://localhost:9080", "connect server addr")
	rootCmd.PersistentFlags().StringVar(&image, "image", "", "process image")
	rootCmd.PersistentFlags().IntVar(&size, "size", 64, "resize/thumb size")
	rootCmd.PersistentFlags().StringVar(&method, "method", "thumbnail", "image process method [thumbnail, resize]")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bimgserver" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bimgserver")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
