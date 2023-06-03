/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate hybrid architecture image",
	Long:  `command to generate hybrid architecture image with multiple architecture images`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err error
		)

		if err = generate(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	generateCmd.Flags().BoolVarP(&insecure, "insecure", "i", false, "insecure registry")
	generateCmd.Flags().StringVarP(&registry, "registry", "r", "registry.paas", "registry address")
	generateCmd.Flags().StringVarP(&amd64, "amd64", "", "", "file path of amd64 images")
	generateCmd.Flags().StringVarP(&arm64, "arm64", "", "", "file path of arm64 images")
}

var (
	insecure               bool
	registry, amd64, arm64 string
)

func generate() error {
	var (
		images = &Images{}
		err    error
	)

	images = images.Initialize(registry, insecure)
	// load arm64 images
	if err = images.Load("arm64", arm64); err != nil {
		return err
	}
	// load amd64 images
	if err = images.Load("amd64", amd64); err != nil {
		return err
	}

	for name, images := range images.Manifests() {
		if err = (&Manifest{}).Initialize(name, insecure, images).Generate(); err != nil {
			return err
		}
	}
	return nil
}
