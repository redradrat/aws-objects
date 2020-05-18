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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/redradrat/cloud-objects/aws/rds"
)

var (
	subnetGroupName string
	subnetIDs       []string
)

// subnetgroupCmd represents the subnetgroup command
var subnetgroupCmd = &cobra.Command{
	Use:   "subnetgroup",
	Args:  OnlyCloudObjectAction(),
	Short: "Interact with the RDS subnetgroup cloud object",
	Long: `Interact with the RDS subnetgroup cloud object. For example:

	*) cloud-objects aws rds subnetgroup create --name testsubnetgroup

	*) cloud-objects aws rds subnetgroup delete --name testsubnetgroup`,
	Run: func(cmd *cobra.Command, args []string) {
		session, err := GetSession(cmd)
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		sg, err := rds.NewSubnetGroup(subnetGroupName, session)
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		spec := rds.SubnetGroupSpec{
			Description: "A test RDS DB SubnetGroup",
			SubnetIDs:   subnetIDs,
			Tags: map[string]string{
				"Test": "test",
			},
		}

		_, err = HandleCloudObject(sg, &spec, CloudObjectAction(args[0]), false)
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		fmt.Println(sg.Status())

	},
}

func init() {
	rdsCmd.AddCommand(subnetgroupCmd)

	// SubnetGroupName
	subnetgroupCmd.Flags().StringVarP(
		&subnetGroupName,
		"name",
		"n",
		"",
		"The name for the RDS DB SubnetGroup",
	)

	// SubnetGroup SubnetIDs
	subnetgroupCmd.Flags().StringSliceVar(
		&subnetIDs,
		"subnets",
		[]string{},
		"The list of SubnetIDs this RDS DB SubnetGroup should contain",
	)

}
