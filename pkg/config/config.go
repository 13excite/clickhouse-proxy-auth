package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

const DefaultConfigPath = "/etc/ch-proxy-auth.yaml"

// Config contains settings of the logger
type Config struct {
	ServerHost        string              `yaml:"server_host"`
	ServerPort        string              `yaml:"server_port"`
	Level             string              `yaml:"log_level"`
	Encoding          string              `yaml:"log_encoding"`
	Color             bool                `yaml:"log_color"`
	DisableStacktrace bool                `yaml:"log_disable_stacktrace"`
	DevMode           bool                `yaml:"log_dev_mode"`
	DisableCaller     bool                `yaml:"log_disable_caller"`
	IgnorePaths       []string            `yaml:"log_ignore_paths"`
	OutputPaths       []string            `yaml:"log_output_paths"` // will skip it for while
	ErrorOutputPaths  []string            `yaml:"log_err_output_paths"`
	HostToCluster     map[string]string   `yaml:"host_to_cluster"`
	NetAclClusters    map[string][]string `yaml:"net_acl_clusters"`
}

// Defaults initializes default logger settings
func (conf *Config) Defaults() {
	conf.ServerHost = "127.0.0.1"
	conf.ServerPort = "8081"
	conf.Level = "info"
	conf.Encoding = "console"
	conf.Color = false
	conf.DisableStacktrace = true
	conf.DevMode = true
	conf.DisableCaller = false
	conf.IgnorePaths = []string{"/version"}
	conf.OutputPaths = []string{"stderr"}
	conf.ErrorOutputPaths = []string{"stderr"}
}

// ReadConfigFile reading and parsing configuration yaml file
func (conf *Config) ReadConfigFile(configPath string) {
	if configPath == "" {
		configPath = DefaultConfigPath
	}
	yamlConfig, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlConfig, &conf)
	if err != nil {
		log.Fatal(fmt.Errorf("could not unmarshal config %v", conf), err)
	}
}

// GetHostsClustersMap returns map of hosts to clusters
func (conf *Config) GetHostsClustersMap() map[string]string {
	return conf.HostToCluster
}

// GetNetAclClusters returns map of clusters to subnets
func (conf *Config) GetNetAclClusters() map[string][]string {
	return conf.NetAclClusters
}
