package params

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
)

// Configuration represents endpoint config
type Configuration struct {
	// RetryTimeOut for endpoint
	RetryTimeOut *int `yaml:"retryTimeOut,omitempty" json:"retryTimeOut,omitempty"`
	// RetryDelay for endpoint
	RetryDelay *int `yaml:"retryDelay,omitempty" json:"retryDelay,omitempty"`
	// Factor used for config
	Factor *int `yaml:"factor,omitempty" json:"factor,omitempty"`
}

// Endpoint details
type Endpoint struct {
	// Type of the endpoints
	EndpointType string `json:"endpoint_type,omitempty"`
	// Url of the endpoint
	Url *string `yaml:"url" json:"url"`
	// Config of endpoint
	Config *Configuration `yaml:"config,omitempty" json:"config,omitempty"`
}

// EndpointData contains details about endpoints
type EndpointData struct {
	// Type of the endpoints
	EndpointType string `json:"endpoint_type"`
	// Production endpoint
	Production *Endpoint `yaml:"production" json:"production_endpoints,omitempty"`
	// Sandbox endpoint
	Sandbox *Endpoint `yaml:"sandbox" json:"sandbox_endpoints,omitempty"`
}

// ApiIdentifier stores API Identifier details
type APIIdentifier struct {
	// Name of the provider of the API
	ProviderName string `json:"providerName"`
	// Name of the API
	APIName string `json:"apiName"`
	// Version of the API
	Version string `json:"version"`
}

type Environment struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"configs"`
}

// ApiParams represents environments defined in configuration file
type ApiParams struct {
	// Environments contains all environments in a configuration
	Environments []Environment `yaml:"environments"`
	Deploy       APIVCSParams  `yaml:"deploy"`
}

type ApiProductParams struct {
	Deploy ApiProductVCSParams `yaml:"deploy"`
}

type ApplicationParams struct {
	Deploy ApplicationVCSParams `yaml:"deploy"`
}

// ------------------- Structs for VCS Import Params ----------------------------------

type ApplicationVCSParams struct {
	Import ApplicationImportParams `yaml:"import"`
}

type APIVCSParams struct {
	Import APIImportParams `yaml:"import"`
}

type ApiProductVCSParams struct {
	Import APIProductImportParams `yaml:"import"`
}

type APIImportParams struct {
	Update           bool `yaml:"update"`
	PreserveProvider bool `yaml:"preserveProvider"`
	RotateRevision   bool `yaml:"rotateRevision"`
}

type APIProductImportParams struct {
	ImportAPIs       bool `yaml:"importApis"`
	UpdateAPIs       bool `yaml:"updateApis"`
	UpdateAPIProduct bool `yaml:"updateApiProduct"`
	PreserveProvider bool `yaml:"preserveProvider"`
	RotateRevision   bool `yaml:"rotateRevision"`
}

type ApplicationImportParams struct {
	Update            bool   `yaml:"update"`
	TargetOwner       string `yaml:"targetOwner"`
	PreserveOwner     bool   `yaml:"preserveOwner"`
	SkipKeys          bool   `yaml:"skipKeys"`
	SkipSubscriptions bool   `yaml:"skipSubscriptions"`
}

type ProjectParams struct {
	Type                       string          `yaml:"type"`
	AbsolutePath               string          `yaml:"absolutePath,omitempty"`
	RelativePath               string          `yaml:"relativePath,omitempty"`
	NickName                   string          `yaml:"nickName,omitempty"`
	FailedDuringPreviousDeploy bool            `yaml:"failedDuringPreviousDeploy,omitempty"`
	Deleted                    bool            `yaml:"deleted,omitempty"`
	MetaData                   *utils.MetaData `yaml:"metaData,omitempty"`
}

type ProjectInfo struct {
	Owner   string `yaml:"owner,omitempty"`
	Name    string `yaml:"name,omitempty"`
	Version string `yaml:"version,omitempty"`
}

// ---------------- End of Structs for Project Details ---------------------------------

// APIEndpointConfig contains details about endpoints in an API
type APIEndpointConfig struct {
	// EPConfig is representing endpoint configuration
	EPConfig string `json:"endpointConfig"`
}

// loads the given file in path and substitutes environment variables that are defined as ${var} or $var in the file.
//	returns the file as string.
func GetEnvSubstitutedFileContent(path string) (string, error) {
	r, err := os.Open(path)
	defer func() {
		_ = r.Close()
	}()
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	str, err := utils.EnvSubstituteForCurlyBraces(string(data))
	if err != nil {
		return "", err
	}
	return str, nil
}

// LoadApiParamsFromDirectory loads an API Project configuration YAML file located in path when the root
// directory is provided instead of yaml file.
//	It returns an error or a valid ApiParams
func LoadApiParamsFromDirectory(path string) (*ApiParams, error) {
	paramsFilePath := filepath.Join(path, utils.ParamFile)
	utils.Logln(utils.LogPrefixInfo + "Loading params from " + paramsFilePath)
	fileContent, err := GetEnvSubstitutedFileContent(paramsFilePath)
	if err != nil {
		return nil, err
	}

	apiParams := &ApiParams{}
	err = yaml.Unmarshal([]byte(fileContent), &apiParams)
	if err != nil {
		return nil, err
	}

	return apiParams, err
}

// LoadApiParamsFromFile loads an API Project configuration YAML file located in path.
//	It returns an error or a valid ApiParams
func LoadApiParamsFromFile(path string) (*ApiParams, error) {
	fileContent, err := GetEnvSubstitutedFileContent(path)
	if err != nil {
		return nil, err
	}

	apiParams := &ApiParams{}
	err = yaml.Unmarshal([]byte(fileContent), &apiParams)
	if err != nil {
		return nil, err
	}

	return apiParams, err
}

// LoadApiProductParamsFromFile loads an API Product project configuration YAML file located in path.
//	It returns an error or a valid ApiProductParams
func LoadApiProductParamsFromFile(path string) (*ApiProductParams, error) {
	fileContent, err := GetEnvSubstitutedFileContent(path)
	if err != nil {
		return nil, err
	}

	apiParams := &ApiProductParams{}
	err = yaml.Unmarshal([]byte(fileContent), &apiParams)
	if err != nil {
		return nil, err
	}

	return apiParams, err
}

// LoadApplicationParamsFromFile loads an Application project configuration YAML file located in path.
//	It returns an error or a valid ApplicationParams
func LoadApplicationParamsFromFile(path string) (*ApplicationParams, error) {
	fileContent, err := GetEnvSubstitutedFileContent(path)
	if err != nil {
		return nil, err
	}

	apiParams := &ApplicationParams{}
	err = yaml.Unmarshal([]byte(fileContent), &apiParams)
	if err != nil {
		return nil, err
	}

	return apiParams, err
}

// ExtractAPIEndpointConfig extracts API endpoint information from a slice of byte b
func ExtractAPIEndpointConfig(b []byte) (string, error) {
	apiConfig := &APIEndpointConfig{}
	err := json.Unmarshal(b, &apiConfig)
	if err != nil {
		return "", err
	}

	return apiConfig.EPConfig, err
}

// GetEnv returns the EndpointData associated for key in the ApiParams, if not found returns nil
func (config ApiParams) GetEnv(key string) *Environment {
	for index, env := range config.Environments {
		if env.Name == key {
			return &config.Environments[index]
		}
	}
	return nil
}
