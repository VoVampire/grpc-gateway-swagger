package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Server struct {
	GrpcHost    string
	GrpcPort    int32
	GatewayHost string
	GatewayPort int32
}

type Config struct {
	Server
}

var Settings = &Config{}

func InitConfig(configPath string) {
	configPath, err := filepath.Abs(configPath)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error happened: %v \n", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(data, Settings)
	if err != nil {
		fmt.Printf("Error happened: %v \n", err)
		os.Exit(1)
	}
}
