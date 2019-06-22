package mongosync

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var Config = &config{}

type config struct {
	Path           string
	Log            string                 `yaml:"log"`
	CheckpointPath string                 `yaml:"checkpoint_path"`
	Mongo          mongoConfig            `yaml:"mongo"`
	DocManagers    map[string]interface{} `yaml:"doc_managers"`
}

type mongoConfig struct {
	Url string `yaml:"url"`
}

func InitConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	err = yaml.NewDecoder(f).Decode(Config)
	if err != nil {
		log.Println(err)
	}
	err = f.Close()
	if err != nil {
		log.Println(err)
	}
	return nil
}
