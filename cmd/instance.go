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

var instanceName string

// instanceCmd represents the instance command
var instanceCmd = &cobra.Command{
	Use:   "instance",
	Args:  OnlyCloudObjectAction(),
	Short: "Interact with the RDS instance cloud object",
	Long: `Interact with the RDS instance cloud object. For example:

	*) cloud-objects aws rds instance create --name testinstance

	*) cloud-objects aws rds instance delete --name testinstance`,
	Run: func(cmd *cobra.Command, args []string) {

		arg := args[0]
		switch CloudObjectAction(arg) {
		case CreateCloudObjectAction:
		}
	},
}

func init() {
	rdsCmd.AddCommand(instanceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//instanceCmd.PersistentFlags().String("create", "", "Create an RDS instance")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	instanceCmd.Flags().StringVarP(&instanceName, "--name", "n", "", "The name of the instance")
}
