package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/midlang/mid/src/mid/build"
)

var (
	midconfig = filepath.Join(os.Getenv("HOME"), ".midconfig")
)

type Config struct {
	Suffix string `cli:"suffix" usage:"source file suffix" dft".mid" name:"SUFFIX"`

	Plugins       *build.PluginSet `json:"-" cli:"-"`
	LoadedPlugins []build.Plugin   `json:"plugins" cli:"-"`
}

func newConfig() *Config {
	cfg := new(Config)
	cfg.Plugins = build.NewPluginSet()
	return cfg
}

func (cfg *Config) Load(filename string) error {
	if filename == "" {
		filename = midconfig
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return err
	}
	for _, p := range cfg.LoadedPlugins {
		if err := cfg.Plugins.Register(&p); err != nil {
			return err
		}
	}
	return nil
}
