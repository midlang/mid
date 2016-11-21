package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/midlang/mid/src/mid"
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
	Version      bool              `cli:"!v,version" usage:"display version information"`
	ConfigFile   string            `cli:"c,config" usage:"config filename"`
	LogLevel     logger.Level      `cli:"log,loglevel" usage:"log level for debugging: trace/debug/info/warn/error/fatal" dft:"warn"`
	Inputs       []string          `cli:"I,input" usage:"input directories or files which has suffix SUFFIX"`
	Outdirs      map[string]string `cli:"O,outdir" usage:"output directories for each language, e.g. -Ogo=dir1 -Ocpp=dir2"`
	Extentions   []string          `cli:"X,extension" usage:"extensions, e.g. -Xproto -Xredis -Xmysql:go (only for go generator)"`
	Envvars      map[string]string `cli:"E,env" usage:"custom defined environment variables"`
	ImportPaths  []string          `cli:"P,importpath" usage:"import paths for lookuping imports"`
	TemplateKind string            `cli:"K,tk,template-kind" usage:"template kind, a directory name" dft:"default"`
	TemplatesDir map[string]string `cli:"T,template" usage:"templates directories for each language, e.g. -Tgo=dir1 -Tjava=dir2"`
}

func newArgT() *argT {
	argv := &argT{
		Outdirs:      map[string]string{},
		TemplatesDir: map[string]string{},
		Envvars:      map[string]string{},
		Config:       *newConfig(),
		ImportPaths:  strings.Split(os.Getenv("MID_IMPORT_PATH"), ":"),
	}
	return argv
}

var root = &cli.Command{
	Name:      "midc",
	Argv:      func() interface{} { return newArgT() },
	Desc:      "midlang compiler - compile source files and generate other languages code",
	NumOption: cli.AtLeast(1),

	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		if argv.Version {
			ctx.String("v%s\n", mid.Meta["version"])
			return nil
		}
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

		// load config file
		if argv.ConfigFile == "" {
			for _, dir := range []string{os.Getenv("HOME"), "/etc", "/usr/local/etc"} {
				fullpath := filepath.Join(dir, "midconfig")
				if tmpInfo, err := os.Lstat(fullpath); err == nil && tmpInfo != nil && !tmpInfo.IsDir() {
					argv.ConfigFile = fullpath
					break
				}
			}
		}
		if argv.ConfigFile == "" {
			log.Error("missing config file")
			return nil
		}
		if err := argv.Config.Load(argv.ConfigFile); err != nil {
			log.Error("load config %s: %v", cyan(argv.ConfigFile), red(err))
			return nil
		}

		// migrate templates directories
		templatesDir := make(map[string]string)
		for lang, dir := range argv.TemplatesDir {
			absDir, err := filepath.Abs(dir)
			if err != nil {
				log.Error("invalid templates directory: %s", red(dir))
				return nil
			}
			log.Debug("language %s templates directory: %s", cyan(lang), dir)
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
				log.Error("init plugin %s: %v", formatPlugin(plugin.Lang, plugin.Name), err)
				hasError = true
				continue
			}
			plugin.RuntimeConfig.Verbose = argv.LogLevel.String()
			if templatesDir, ok := argv.TemplatesDir[plugin.Lang]; ok {
				// replace default templatesDir
				plugin.TemplatesDir = templatesDir
			}
			if plugin.TemplatesDir == "" {
				// if templatesDir is empty
				var pendingDir []string
				if argv.Config.TemplatesRootDir != "" {
					pendingDir = append(pendingDir, argv.Config.TemplatesRootDir)
				} else {
					pendingDir = []string{
						filepath.Join(os.Getenv("HOME"), "mid_templates"),
						filepath.Join("/etc", "mid_templates"),
						filepath.Join("/usr/local/usr", "mid_templates"),
					}
				}
				for _, dir := range pendingDir {
					fullpath := filepath.Join(dir, argv.TemplateKind, plugin.Lang)
					log.Trace("try lookup templates dir for plugin %s in directory %s, fullpath=%s", plugin.Lang, dir, fullpath)
					tmpInfo, err := os.Lstat(fullpath)
					if err != nil || tmpInfo == nil || !tmpInfo.IsDir() {
						continue
					}
					plugin.TemplatesDir, err = filepath.Abs(fullpath)
					if err != nil {
						log.Error("get abs of path `%s` error: %v", fullpath, err)
						return nil
					}
					break
				}
				if plugin.TemplatesDir == "" {
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
