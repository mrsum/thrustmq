package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type exhaustConfigStruct struct {
	Port          int
	TurbineBlades int
	RecurringFlux int
}

type intakeConfigStruct struct {
	Port             int
	CompressorBlades int
}

type ConfigStruct struct {
	Intake   intakeConfigStruct
	Exhaust  exhaustConfigStruct
	Filename string
	Logfile  string
	Debug    bool
}

func loadConfig() ConfigStruct {
	cfg_file := "./config.json"

	raw, err := ioutil.ReadFile(cfg_file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var config ConfigStruct
	json.Unmarshal(raw, &config)

	return config
}

var Config = loadConfig()
