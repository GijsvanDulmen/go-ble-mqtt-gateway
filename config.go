package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type MacList struct {
	Mac string `json:"mac"`
}

type NameList struct {
	Name string `json:"name"`
}

type ConfigFile struct {
	MacWhitelist   []MacList
	MacBlacklist   []MacList
	BlacklistNames []NameList
}

type Config struct {
	configFile ConfigFile
	hasConfig  bool
}

func NewConfig() Config {
	useConfig := os.Getenv("CONFIG")

	fmt.Println(useConfig)

	var configFile ConfigFile
	if useConfig != "" {
		jsonFile, err := os.Open(useConfig)

		// if we os.Open returns an error then handle it
		if err != nil {
			log.Fatalf("Could not open config file")
		}

		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log.Fatalf("Could not read config file")
		}

		err = json.Unmarshal(byteValue, &configFile)
		if err != nil {
			log.Fatalf("Could not parse config file")
		}
		fmt.Println("Done reading config")
	}

	fmt.Println(configFile.MacWhitelist)
	fmt.Println(configFile.MacBlacklist)
	fmt.Println(configFile.BlacklistNames)
	return Config{
		configFile: configFile,
		hasConfig:  useConfig != "",
	}
}

func (c Config) MatchesAgainstConfig(addr string, name string) bool {
	if !c.hasConfig {
		return true
	}

	for _, item := range c.configFile.BlacklistNames {
		if item.Name == name {
			return false
		}
	}

	for _, item := range c.configFile.MacBlacklist {
		if item.Mac == addr {
			return false
		}
	}

	// apply whitelist if there is any
	if len(c.configFile.MacWhitelist) > 0 {
		for _, item := range c.configFile.MacWhitelist {
			if item.Mac == addr {
				return true
			}
		}
		return false
	}
	return true
}
