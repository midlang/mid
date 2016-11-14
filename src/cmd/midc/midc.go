package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/midlang/mid/src/mid/build"
	"github.com/midlang/mid/src/mid/lexer"
	"github.com/midlang/mid/src/mid/parser"
	"github.com/mkideal/cli"
	"github.com/mkideal/log"
	"github.com/mkideal/log/logger"
)

type argT struct {
	cli.Helper
	Config
	ConfigFile   string            `cli:"c,config" usage:"config filename"`
	LogLevel     logger.Level      `cli:"v,loglevel" usage:"log level for debugging: trace/debug/info/warn/error/fatal" dft:"warn"`
	Inputs       []string          `cli:"I,input" usage:"input directories or files which has suffix SUFFIX"`
	Outdirs      map[string]string `cli:"O,outdir" usage:"output directories for each language, e.g. -Ogo=dir1 -Ocpp=dir2"`
	Extentions   []string          `cli:"X,extension" usage:"extensions, e.g. -Xproto -Xredis -Xmysql:go (only for go generator)"`
	Envvars      map[string]string `cli:"E,env" usage:"custom defined environment variables"`
	ImportPaths  []string          `cli:"P,importpath" usage:"import paths for lookuping imports"`
	TemplatesDir map[string]string `cli:"T,template" usage:"templates directories for each language, e.g. -Tgo=dir1 -Tjava=dir2"`
}

func newArgT() *argT {
	argv := &argT{
		Outdirs:      map[string]string{},
		TemplatesDir: map[string]string{},
		Envvars:      map[string]string{},
		Config:       *newConfig(),
		ConfigFile:   filepath.Join(os.Getenv("HOME"), ".midconfig"),
		ImportPaths:  strings.Split(os.Getenv("MID_IMPORT_PATH"), ":"),
	}
	return argv
}

var root = &cli.Command{
	Name:      "midc",
	Argv:      func() interface{} { return newArgT() },
	Desc:      "midlang compiler",
	NumOption: cli.AtLeast(1),

	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		log.SetLevel(argv.LogLevel)
		if !argv.LogLevel.MoreVerboseThan(log.LvINFO) {
			log.NoHeader()
		}
		log.WithJSON(argv).Debug("argv")

		var (
			blue = ctx.Color().Blue
			cyan = ctx.Color().Cyan
			red  = ctx.Color().Red

			inputs  []string
			plugins []*build.Plugin

			formatPlugin = func(lang, name string) string {
				return "<" + blue(lang) + ":" + cyan(name) + ">"
			}
			sourceFileFilter = func(finfo os.FileInfo) bool {
				return strings.HasSuffix(finfo.Name(), argv.Suffix)
			}
		)

		if err := argv.Config.Load(argv.ConfigFile); err != nil {
			log.Error("load config %s: %v", cyan(argv.ConfigFile), red(err))
		}

		templatesDir := make(map[string]string)
		for lang, dir := range argv.TemplatesDir {
			absDir, err := filepath.Abs(dir)
			if err != nil {
				log.Error("invalid templates dir: %s", red(dir))
				return nil
			}
			templatesDir[lang] = absDir
		}
		argv.TemplatesDir = templatesDir

		// validate source directories and files
		if len(argv.Inputs) == 0 {
			argv.Inputs = []string{"."}
		}
		for _, in := range argv.Inputs {
			finfo, err := os.Lstat(in)
			if err != nil {
				log.Error("input %s: %v", cyan(in), red(err))
				return nil
			}
			if finfo.IsDir() {
				files, err := filesInDir(in, sourceFileFilter)
				if err != nil {
					log.Error("get source files from dir %s: %v", cyan(in), red(err))
					return nil
				}
				inputs = append(inputs, files...)
			} else {
				inputs = append(inputs, in)
			}
		}
		log.Debug("inputs: %v", inputs)

		// lookup plugins
		var hasError bool
		for lang, outdir := range argv.Outdirs {
			if outdir == "" {
				log.Error("language %s output directory is empty", blue(lang))
				hasError = true
			}
			plugin, ok := argv.Plugins.Lookup(lang)
			if !ok {
				log.Error("language plugin %s not found", blue(lang))
				hasError = true
				continue
			}
			if err := plugin.Init(outdir, argv.Extentions, argv.Envvars); err != nil {
				log.Error("init plugin %s: %s", formatPlugin(plugin.Lang, plugin.Name))
				hasError = true
				continue
			}
			plugin.RuntimeConfig.Verbose = argv.LogLevel.String()
			if templatesDir, ok := argv.TemplatesDir[plugin.Lang]; ok {
				// replace default templatesDir
				plugin.TemplatesDir = templatesDir
			}
			if plugin.TemplatesDir == "" {
				if argv.Config.TemplatesRootDir != "" {
					plugin.TemplatesDir = filepath.Join(argv.Config.TemplatesRootDir, plugin.Lang)
				} else {
					log.Error("templates directory of plugin %s missing", formatPlugin(plugin.Lang, plugin.Name))
					hasError = true
					continue
				}
			}
			plugins = append(plugins, plugin)
			for _, x := range argv.Extentions {
				if !plugin.IsSupportExt(x) {
					log.Warn("plugin %s does not support extension %s", formatPlugin(plugin.Lang, plugin.Name), x)
				}
			}
		}
		if hasError {
			return nil
		}

		// build source
		fset := lexer.NewFileSet()
		pkgs, err := parser.ParseFiles(fset, inputs)
		if err != nil {
			log.Error("parse error:\n%v", red(err))
			return nil
		}
		builder, err := build.Build(pkgs, argv.ImportPaths)
		if err != nil {
			log.Error("build error: %v", red(err))
			return nil
		}

		log.Debug("len(pkgs): %d", len(pkgs))
		for name, _ := range pkgs {
			log.Debug("package %s", cyan(name))
		}

		// generate codes
		for _, plugin := range plugins {
			log.Debug("ready execute plugin %s", formatPlugin(plugin.Lang, plugin.Name))
			if err := plugin.Generate(builder, os.Stdout, os.Stderr); err != nil {
				log.Error("plugin %s generate codes error: %v", formatPlugin(plugin.Lang, plugin.Name), red(err))
			}
		}
		return nil
	},
}

func main() {
	defer log.Uninit(log.InitConsole(log.LvTRACE))
	err := root.Run(os.Args[1:])
	log.If(err != nil).Error("%v", err)
}

func filesInDir(dir string, filter func(os.FileInfo) bool) ([]string, error) {
	fd, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	list, err := fd.Readdir(-1)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, d := range list {
		if d.IsDir() || (filter != nil && !filter(d)) {
			continue
		}
		filename := filepath.Join(dir, d.Name())
		files = append(files, filename)
	}
	return files, nil
}
