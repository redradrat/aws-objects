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
	"github.com/redradrat/cloud-objects/aws/kms"

	"fmt"

	"github.com/spf13/cobra"
)

var keyName string
var keyUsage kms.KeyUsage
var keyType kms.KeyType

// keyCmd represents the key command
var keyCmd = &cobra.Command{
	Use:   "key",
	Args:  OnlyCloudObjectAction(),
	Short: "Interact with the KMS key cloud object",
	Long: `Interact with the KMS key cloud object. For example:

	*) cloud-objects aws kms key create --name testkey

	*) cloud-objects aws kms key delete --name testkey`,
	Run: func(cmd *cobra.Command, args []string) {
		session, err := GetSession(cmd)
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		key, err := kms.NewKey(keyName, session)
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}

		spec := kms.KeySpec{
			KeyUsage: keyUsage,
			KeyType:  keyType,
		}

		prg, err := cmd.Flags().GetBool(PurgeFlag)
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		_, err = HandleCloudObject(key, &spec, CloudObjectAction(args[0]), prg)
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		fmt.Println(key.Status())
	},
}

func init() {
	kmsCmd.AddCommand(keyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//keyCmd.PersistentFlags().String("create", "", "Create an KMS key")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	keyCmd.Flags().StringVarP(&keyName, "name", "n", "", "The name of the key")
	keyCmd.Flags().StringVar((*string)(&keyUsage), "keyUsage", string(kms.EncryptDecryptKeyUsage),
		"The purpose of the key (e.g. 'ENCRYPT_DECRYPT', 'SIGN_VERIFY')")
	keyCmd.Flags().StringVar((*string)(&keyType), "keyType", string(kms.SymmetricDefaultKeyType),
		"The key type to use (e.g. 'SYMMETRIC_DEFAULT', 'RSA_2048', ...)")

}
