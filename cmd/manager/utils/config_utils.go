package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/ugcompsoc/apid/internal/config"
	"gopkg.in/yaml.v2"
)

func PrintConfig(c *config.Config, printSecrets bool) (string, error) {
	if c == nil {
		return "", errors.New("A nil config struct was given")
	}
	if !printSecrets {
		if len(c.Database.Username) != 0 {
			c.Database.Username = "********"
		}
		if len(c.Database.Password) != 0 {
			c.Database.Password = "********"
		}
	}
	cYaml, err := yaml.Marshal(&c)
	// im not sure how we would reach this error unless the config is nil
	// but we check for that above so we dont check for the length on nil
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not marshall the config struct: %s\n", err))
	}
	return string(cYaml), nil
}

func VerifyFilename(filename string) error {
	filenameRegex, err := regexp.Compile("^([a-z]+).yml$")
	// unless this regex is changed, this error will never be reached
	if err != nil {
		return err
	}
	if !filenameRegex.MatchString(filename) {
		return errors.New("The filename is not in the form [NAME].yml")
	}
	return nil
}

func ExtractFile(absoluteFilePath string) (*config.Config, error) {
	file, err := ioutil.ReadFile(absoluteFilePath)
	if err != nil {
		return nil, err
	}
	fileStr := string(file)
	if len(fileStr) == 0 {
		return nil, errors.New("The file apid.yml is completely empty. What do you want me to do with this?")
	}
	c := config.Config{}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
