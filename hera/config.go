package hera

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	Services map[string]*Service `yaml:"services"`
}

func (config *Config) prime() {
	for _, service := range config.Services {
		service.prime()
	}
}

type Service struct {
	Command           string   `yaml:"command"`
	Watch             []string `yaml:"watch"`
	Exclude           []string `yaml:"exclude"`
	prefixesToWatch   []string
	prefixesToExclude []string
}

func (service *Service) prime() {
	wd, _ := os.Getwd()

	clean := func(path string) string {
		if path == "." {
			path = ""
		}
		path = strings.TrimPrefix(path, "/")
		path = fmt.Sprintf("%s/%s", wd, path)
		return path
	}

	for _, path := range service.Watch {
		service.prefixesToWatch = append(service.prefixesToWatch, clean(path))
	}

	for _, path := range service.Exclude {
		service.prefixesToExclude = append(service.prefixesToExclude, clean(path))
	}
}

func (service *Service) shouldTriggerUpdate(fileName string) bool {
	for _, exclude := range service.prefixesToExclude {
		if strings.HasPrefix(fileName, exclude) {
			return false
		}
	}

	for _, path := range service.prefixesToWatch {
		if strings.Contains(fileName, path) {
			file, err := os.OpenFile(".config/tmp/hera.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}

			log.New(file, "", log.LstdFlags).Println(fileName)
			if err := file.Close(); err != nil {
				panic(err)
			}

			return true
		}
	}

	return false
}
