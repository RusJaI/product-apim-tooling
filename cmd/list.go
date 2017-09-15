/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
*/

package cmd

import (
	"fmt"

	"crypto/tls"
	"encoding/json"
	"github.com/go-resty/resty"
	"github.com/menuka94/wso2apim-cli/utils"
	"github.com/spf13/cobra"
	"errors"
)

var listEnvironment string
var listCmdUsername string
var listCmdPassword string

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: utils.ListCmdShortDesc,
	Long:  utils.ListCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln("list called")

		accessToken, apiManagerEndpoint, preCommandErr := utils.ExecutePreCommand(listEnvironment, listCmdUsername, listCmdPassword)

		if preCommandErr == nil {
			count, apis, err := GetAPIList("", accessToken, apiManagerEndpoint)

			if err == nil {
				fmt.Println("No. of APIs:", count)
				for _, api := range apis {
					fmt.Println(api.Name + " v" + api.Version)
				}
			} else {
				utils.Logln(utils.LogPrefixError + "Getting List of APIs", err)
			}
		} else {
			utils.Logln(utils.LogPrefixError + "calling 'list' "+ preCommandErr.Error())
			fmt.Println("Error calling 'list'", preCommandErr.Error())
		}
	},
}

func GetAPIList(query string, accessToken string, apiManagerEndpoint string) (int32, []utils.API, error) {
	url := apiManagerEndpoint

	// append '/' to the end if there isn't one already
	if url != "" && string(url[len(url)-1]) != "/" {
		url += "/"
	}
	url += "apis?query=" + query
	fmt.Println("URL:", url)

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in HTTPS certificates
	resp, err := resty.R().
		SetHeaders(headers).
		Get(url)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to " + url, err)
	}

	utils.Logln(utils.LogPrefixInfo + "GetAPIList(): Response:", resp.Status())

	if resp.StatusCode() == 200 {
		apiListResponse := &utils.APIListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &apiListResponse)

		if unmarshalError != nil {
			fmt.Println("UnmarshalError")
			utils.HandleErrorAndExit(utils.LogPrefixError + "Unmarshall Error", unmarshalError)
		}

		return apiListResponse.Count, apiListResponse.List, nil
	} else {
		return 0, nil, errors.New(resp.Status())
	}

}

func init() {
	RootCmd.AddCommand(ListCmd)
	ListCmd.Flags().StringVarP(&listEnvironment, "environment", "e", "", "Environment to be searched")
	ListCmd.Flags().StringVarP(&listCmdUsername, "usrename", "u", "", "Username")
	ListCmd.Flags().StringVarP(&listCmdPassword, "password", "p", "", "Password")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
