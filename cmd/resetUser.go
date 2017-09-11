// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/menuka94/wso2apim-cli/utils"
	"github.com/spf13/cobra"
)

var resetUserEnvironment string

// ResetUserCmd represents the resetUser command
var ResetUserCmd = &cobra.Command{
	Use:   "reset-user",
	Short: utils.ResetUserCmdShortDesc,
	Long:  utils.ResetUserCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "reset-user called")

		err := utils.RemoveEnvFromKeysFile(resetUserEnvironment, utils.EnvKeysAllFilePath)
		if err != nil {
			utils.HandleErrorAndExit("Error clearing user data for environment " + resetUserEnvironment, err)
		} else {
			fmt.Println("Successfully cleared user data for environment: " + resetUserEnvironment)
		}
	},
}

func init() {
	RootCmd.AddCommand(ResetUserCmd)
	ResetUserCmd.Flags().StringVarP(&resetUserEnvironment, "environment", "e", "", "Clear user details of an environment")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ResetUserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ResetUserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
