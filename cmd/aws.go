/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/spf13/cobra"
)

const (
	RegionFlag = "region"
	PurgeFlag  = "purge"
)

// awsCmd represents the aws command
var awsCmd = &cobra.Command{
	Use:              "aws",
	TraverseChildren: true,
	Short:            "Interact with AWS cloud objects",
	Long:             ``,
}

func init() {
	RootCmd.AddCommand(awsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	awsCmd.PersistentFlags().String(RegionFlag, "", "The AWS region to work with")
	awsCmd.PersistentFlags().Bool(PurgeFlag, false, "Whether to purge on deletion")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// awsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
