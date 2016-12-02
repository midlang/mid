package build

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
)

type PluginRuntimeConfig struct {
	Verbose    string
	Outdir     string
	Extentions []Extention
	Envvars    map[string]string
}

func (config PluginRuntimeConfig) Encode() string {
	data, _ := json.Marshal(config)
	return string(data)
}

func (config *PluginRuntimeConfig) Decode(s string) error {
	return json.Unmarshal([]byte(s), config)
}

func (config *PluginRuntimeConfig) Getenv(name string) string {
	if config.Envvars == nil {
		return ""
	}
	return config.Envvars[name]
}

func (config *PluginRuntimeConfig) BoolEnv(name string) bool {
	if config.Envvars != nil {
		v, found := config.Envvars[name]
		if found && v == "" {
			return true
		}
	}
	s := config.Getenv(name)
	b, _ := strconv.ParseBool(s)
	return b
}

func (config *PluginRuntimeConfig) IntEnv(name string) int64 {
	s := config.Getenv(name)
	i, _ := strconv.ParseInt(s, 0, 64)
	return i
}

func (config *PluginRuntimeConfig) UintEnv(name string) uint64 {
	s := config.Getenv(name)
	u, _ := strconv.ParseUint(s, 0, 64)
	return u
}

func (config *PluginRuntimeConfig) FloatEnv(name string) float64 {
	s := config.Getenv(name)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

type Plugin struct {
	Lang         string `json:"lang"`
	Name         string `json:"name"`
	Bin          string `json:"bin"`
	TemplatesDir string `json:"templates,omitempty"`

	RuntimeConfig PluginRuntimeConfig `json:"-"`
}

func (plugin *Plugin) Init(outdir string, extensions []Extention, envvars map[string]string) error {
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

func (plugin Plugin) Generate(builder *Builder, stdout, stderr io.Writer) error {
	source := builder.Encode()
	runtimeConfig := plugin.RuntimeConfig.Encode()
	encodedPlugin, err := json.Marshal(plugin)
	if err != nil {
		return err
	}
	cmd := exec.Command(plugin.Bin, string(encodedPlugin), runtimeConfig, source)
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

func (pset *PluginSet) Len() int {
	return len(pset.plugins)
}

func (pset *PluginSet) Register(plugin Plugin) error {
	if _, existed := pset.Lookup(plugin.Lang); existed {
		return fmt.Errorf("plugin %s existed", plugin.Lang)
	}
	pset.plugins = append(pset.plugins, &plugin)
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

func ParseFlags() (plugin Plugin, config PluginRuntimeConfig, builder *Builder, err error) {
	args := os.Args[1:]
	if len(args) != 3 {
		err = fmt.Errorf("usage: %s <json_ecnoded_plugin> <json_ecnoded_config> <gob_and_base64_encoded_source>", os.Args[0])
		return
	}
	if err = json.Unmarshal([]byte(args[0]), &plugin); err != nil {
		err = fmt.Errorf("decode plugin error: %v", err)
		return
	}
	if err = config.Decode(args[1]); err != nil {
		err = fmt.Errorf("decode config error: %v", err)
		return
	}
	builder = new(Builder)
	if err = builder.Decode(args[2]); err != nil {
		err = fmt.Errorf("decode source error: %v", err)
	}
	return
}
