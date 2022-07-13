package conf

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var (
	Conf *Config

	configFilePath = "./config.yml"
)

type Config struct {
	Pprof     bool   `yaml:"pprof"`
	Proxy     string `yaml:"proxy"`
	ApiKey    string `yaml:"apiKey"`
	SecretKey string `yaml:"secretKey"`
}

func init() {
	var (
		err error
	)
	Conf = new(Config)

	yamlFileBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(yamlFileBytes, Conf)
	if err != nil {
		log.Fatal(err)
	}

}

func WriteYAML(conf Config) {
	var (
		err error
	)

	out, err := yaml.Marshal(&conf)
	if err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile(configFilePath, out, 0777)
	if err != nil {
		log.Fatalln(err)
	}
}
