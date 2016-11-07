package build

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type PluginRuntimeConfig struct {
	Outdir     string
	Extentions []string
	Envvars    map[string]string
}

func (config PluginRuntimeConfig) Encode() string {
	data, _ := json.Marshal(config)
	return string(data)
}

func (config *PluginRuntimeConfig) Decode(s string) error {
	return json.Unmarshal([]byte(s), config)
}

type Plugin struct {
	Lang                string   `json:"lang"`
	Name                string   `json:"name"`
	Bin                 string   `json:"bin"`
	SupportedExtentions []string `json:"supported_exts"`

	RuntimeConfig PluginRuntimeConfig `json:"-"`
}

func (plugin *Plugin) Init(outdir string, extensions []string, envvars map[string]string) error {
	bin, err := exec.LookPath(plugin.Bin)
	if err != nil {
		return err
	}
	plugin.Bin = bin
	plugin.RuntimeConfig.Outdir = outdir
	plugin.RuntimeConfig.Extentions = extensions
	plugin.RuntimeConfig.Envvars = envvars
	return nil
}

func (plugin *Plugin) IsSupportExt(ext string) bool {
	for _, x := range plugin.SupportedExtentions {
		if ext == x {
			return true
		}
	}
	return false
}

// /path/to/bin <out_dir> <source_string>
func (plugin Plugin) Generate(builder *Builder, stdout, stderr io.Writer) error {
	source := builder.Encode()
	runtimeConfig := plugin.RuntimeConfig.Encode()
	cmd := exec.Command(plugin.Bin, runtimeConfig, source)
	if stdout != nil {
		cmd.Stdout = stdout
	}
	if stderr != nil {
		cmd.Stderr = stderr
	}
	return cmd.Run()
}

type PluginSet struct {
	plugins []*Plugin
}

func NewPluginSet() *PluginSet {
	return &PluginSet{plugins: []*Plugin{}}
}

func (pset *PluginSet) Register(plugin *Plugin) error {
	if _, existed := pset.Lookup(plugin.Lang); existed {
		return fmt.Errorf("plugin %s existed", plugin.Lang)
	}
	pset.plugins = append(pset.plugins, plugin)
	return nil
}

func (pset *PluginSet) Lookup(lang string) (*Plugin, bool) {
	for _, plugin := range pset.plugins {
		if plugin.Lang == lang {
			return plugin, true
		}
	}
	return nil, false
}

func ParseFlags() (config PluginRuntimeConfig, builder *Builder, err error) {
	args := os.Args[1:]
	if len(args) != 2 {
		err = fmt.Errorf("usage: %s <json_ecnoded_config> <gob_and_base64_encoded_source>", os.Args[0])
		return
	}
	if err = config.Decode(args[0]); err != nil {
		err = fmt.Errorf("decode config error: %v", err)
		return
	}
	source := args[1]
	builder = new(Builder)
	if err = builder.Decode(source); err != nil {
		err = fmt.Errorf("decode source error: %v", err)
	}
	return
}
