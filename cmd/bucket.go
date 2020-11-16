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
	"github.com/redradrat/cloud-objects/aws/s3"

	"fmt"

	"github.com/spf13/cobra"
)

var bucketName string

// bucketCmd represents the bucket command
var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Args:  OnlyCloudObjectAction(),
	Short: "Interact with the S3 bucket cloud object",
	Long: `Interact with the S3 bucket cloud object. For example:

	*) cloud-objects aws s3 bucket create --name testbucket

	*) cloud-objects aws s3 bucket delete --name testbucket`,
	Run: func(cmd *cobra.Command, args []string) {
		session, err := GetSession(cmd)
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		ins, err := s3.NewBucket(bucketName, session)
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		spec := s3.SaneS3Bucket()

		_, err = HandleCloudObject(ins, &spec, CloudObjectAction(args[0]), false)
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		fmt.Println(ins.Status().String())
	},
}

func init() {
	s3Cmd.AddCommand(bucketCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//bucketCmd.PersistentFlags().String("create", "", "Create an RDS bucket")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	bucketCmd.Flags().StringVarP(&bucketName, "name", "n", "", "The name of the bucket")
}
