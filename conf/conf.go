package conf

import (
	"io/ioutil"
	"log"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var Conf *conf

type conf struct {
	Proxy     string `yaml:"proxy"`
	ApiKey    string `yaml:"apiKey"`
	SecretKey string `yaml:"secretKey"`
}

func init() {
	var err error
	Conf = new(conf)

	yamlFileBytes, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		logrus.Fatalln(err)
	}

	err = yaml.Unmarshal(yamlFileBytes, Conf)
	if err != nil {
		log.Fatal(err)
	}
}
