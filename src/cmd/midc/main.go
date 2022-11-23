package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gopherd/log"
	"github.com/midlang/mid/src/mid"
	"github.com/midlang/mid/src/mid/build"
	"github.com/midlang/mid/src/mid/lexer"
	"github.com/midlang/mid/src/mid/parser"
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Config
	Version      bool              `cli:"!v,version" usage:"display version information"`
	ConfigFile   string            `cli:"c,config" usage:"config filename"`
	LogLevel     log.Level         `cli:"log" usage:"log level for debugging: trace/debug/info/warn/error/fatal" dft:"warn"`
	Outdirs      map[string]string `cli:"O,outdir" usage:"output directories for each language, e.g. -Ogo=dir1 -Ocpp=dir2"`
	Extensions   []string          `cli:"X,extension" usage:"extensions, e.g. -Xmeta -Xcodec"`
	Envvars      map[string]string `cli:"E,env" usage:"custom defined environment variables"`
	ImportPaths  []string          `cli:"I,importpath" usage:"import paths for lookuping imports"`
	TemplateKind string            `cli:"K,tempkind" usage:"template kind, a directory name" dft:"default"`
	TemplatesDir map[string]string `cli:"T,template" usage:"templates directories for each language, e.g. -Tgo=dir1 -Tjava=dir2"`
	PluginFiles  map[string]string `cli:"P" usage:"plugin generator file"`
	IdAllocator  string            `cli:"id-allocator" usage:"id allocator name and options,supported allocators: file"`
	IdFor        string            `cli:"id-for" usage:"specific bean kinds which should be allocated a id"`

	Inputs []string `cli:"-"`
}

func newArgT() *argT {
	argv := &argT{
		Outdirs:      map[string]string{},
		TemplatesDir: map[string]string{},
		Envvars:      map[string]string{},
		Config:       *newConfig(),
	}
	if s := os.Getenv("MID_IMPORT_PATH"); s != "" {
		argv.ImportPaths = strings.Split(s, string(filepath.ListSeparator))
	}
	return argv
}

var root = &cli.Command{
	Name:        "midc",
	Argv:        func() interface{} { return newArgT() },
	Desc:        "midlang compiler - compile source files and generate code or documentation",
	CanSubRoute: true,

	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		if argv.Version {
			ctx.String("v%v\n", mid.Meta["version"])
			return nil
		}

		// initialize log
		log.SetLevel(argv.LogLevel)

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
			var paths []string
			if homeDir, err := os.UserHomeDir(); err == nil && homeDir != "" {
				paths = append(paths, homeDir)
			}
			paths = append(paths, "/usr/local/etc", "/etc")
			for _, dir := range paths {
				fullpath := filepath.Join(dir, ".midconfig")
				if tmpInfo, err := os.Lstat(fullpath); err == nil && tmpInfo != nil && !tmpInfo.IsDir() {
					argv.ConfigFile = fullpath
					break
				}
			}
		}
		if argv.ConfigFile != "" {
			if err := argv.Config.Load(argv.ConfigFile); err != nil {
				log.Error().
					String("filename", cyan(argv.ConfigFile)).
					String("error", red(err)).
					Print("load config failed")
				return nil
			}
		}
		// set MidRoot
		if argv.Config.MidRoot == "" {
			argv.Config.MidRoot = filepath.Join(os.Getenv("HOME"), ".mid")
		}

		// load extensions
		extensionsDir := filepath.Join(argv.Config.MidRoot, "extensions")
		extensions, err := loadExtensions(extensionsDir, argv)
		if err != nil {
			return nil
		}

		// migrate templates directories
		templatesDir := make(map[string]string)
		for lang, dir := range argv.TemplatesDir {
			absDir, err := filepath.Abs(dir)
			if err != nil {
				log.Error().
					String("dir", red(dir)).
					Printf("invalid templates directory")
				return nil
			}
			log.Debug().
				String("lang", cyan(lang)).
				String("dir", dir).
				Print("language templates directory")
			templatesDir[lang] = absDir
		}
		argv.TemplatesDir = templatesDir

		// validate source directories and files
		argv.Inputs = ctx.Args()
		if len(argv.Inputs) == 0 {
			argv.Inputs = []string{"."}
		}
		for _, in := range argv.Inputs {
			finfo, err := os.Lstat(in)
			if err != nil {
				log.Error().
					String("input", cyan(in)).
					String("error", red(err)).
					Print("stat error")
				return nil
			}
			if finfo.IsDir() {
				files, err := filesInDir(in, sourceFileFilter)
				if err != nil {
					log.Error().
						String("input", cyan(in)).
						String("error", red(err)).
						Print("get source files from dir error")
					return nil
				}
				inputs = append(inputs, files...)
			} else {
				inputs = append(inputs, in)
			}
		}

		// lookup plugins
		var hasError bool
		for lang, outdir := range argv.Outdirs {
			if outdir == "" {
				log.Error().
					String("lang", blue(lang)).
					Print("language output directory is empty")
				hasError = true
			}
			var (
				plugin *build.Plugin
				ok     bool
			)
			if pluginFilename, found := argv.PluginFiles[lang]; found {
				plugin = &build.Plugin{
					Lang: lang,
					Bin:  pluginFilename,
					Name: "temp",
				}
			} else {
				if argv.Plugins == nil {
					log.Error().
						String("lang", blue(lang)).
						Print("language plugin not found")
					hasError = true
					continue
				}
				plugin, ok = argv.Plugins.Lookup(lang)
				if !ok {
					log.Error().
						String("lang", blue(lang)).
						Print("language plugin not found")
					hasError = true
					continue
				}
			}
			if err := plugin.Init(); err != nil {
				log.Error().
					String("plugin", formatPlugin(plugin.Lang, plugin.Name)).
					Error("error", err).
					Print("init plugin error")
				hasError = true
				continue
			}
			oldOutdir := outdir
			outdir, err = filepath.Abs(outdir)
			if err != nil {
				log.Error().
					String("outdir", oldOutdir).
					Error("error", err).
					Print("get abs of outdir error")
				hasError = true
				continue
			}
			// initialize RuntimeConfig for plugin
			plugin.RuntimeConfig.Outdir = outdir
			plugin.RuntimeConfig.ExtentionsDir = extensionsDir
			plugin.RuntimeConfig.Extensions = extensions
			plugin.RuntimeConfig.Envvars = argv.Envvars
			plugin.RuntimeConfig.Verbose = argv.LogLevel.String()
			if templatesDir, ok := argv.TemplatesDir[plugin.Lang]; ok {
				// replace default templatesDir
				plugin.TemplatesDir = templatesDir
			}
			if plugin.TemplatesDir == "" {
				// if templatesDir is empty
				templatesRootDir := filepath.Join(argv.MidRoot, "templates")
				fullpath := filepath.Join(templatesRootDir, argv.TemplateKind, plugin.Lang)
				plugin.TemplatesDir, err = filepath.Abs(fullpath)
				if err != nil {
					log.Error().
						String("fullpath", fullpath).
						Error("error", err).
						Print("get abs of path error")
					return nil
				}
				if plugin.TemplatesDir == "" {
					log.Error().
						String("plugin", formatPlugin(plugin.Lang, plugin.Name)).
						Print("templates directory of plugin missing")
					hasError = true
					continue
				}
			}
			plugins = append(plugins, plugin)
		}
		if hasError {
			return nil
		}

		// build source
		fset := lexer.NewFileSet()
		pkgs, err := parser.ParseFiles(fset, argv.ImportPaths, inputs)
		if err != nil {
			log.Error().
				String("error", red(err)).
				Print("parse error")
			return nil
		}
		builder, err := build.Build(pkgs)
		if err != nil {
			log.Error().
				String("error", red(err)).
				Print("build error")
			return nil
		}

		// allocate id for beans which kind contained in argv.IdFor
		allocatorInfos := strings.SplitN(argv.IdAllocator, ":", 2)
		if argv.IdAllocator != "" && len(allocatorInfos) != 0 {
			allocatorName := allocatorInfos[0]
			allocatorOpts := ""
			if len(allocatorInfos) == 2 {
				allocatorOpts = allocatorInfos[1]
			}
			allocator, err := build.NewBeanIdAllocator(allocatorName, allocatorOpts)
			if err != nil {
				log.Error().
					Error("error", err).
					Print("new bean id allocator error")
				return err
			}
			idFor := make(map[string]bool)
			for _, f := range strings.Split(argv.IdFor, ",") {
				idFor[strings.TrimSpace(f)] = true
			}
			for _, pkg := range builder.Packages {
				for _, file := range pkg.Files {
					for _, bean := range file.Beans {
						if idFor[bean.Kind] {
							bean.Id = allocator.Allocate(build.JoinBeanKey(pkg.Name, bean.Name))
						}
					}
				}
			}
			if err := allocator.Output(nil); err != nil {
				log.Error().
					Error("error", err).
					Print("id allocator output error")
				return err
			}
		}

		// generate codes
		for _, plugin := range plugins {
			log.Debug().
				String("plugin", formatPlugin(plugin.Lang, plugin.Name)).
				Print("ready execute plugin")
			if err := plugin.Generate(builder, os.Stdout, os.Stderr); err != nil {
				log.Error().
					String("plugin", formatPlugin(plugin.Lang, plugin.Name)).
					String("error", red(err)).
					Print("plugin generate codes error")
			}
		}
		return nil
	},
}

func main() {
	err := root.Run(os.Args[1:])
	log.If(err != nil).Error().
		Error("error", err).
		Print("run error")
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

func loadExtensions(extensionsDir string, argv *argT) (extensions []build.Extension, err error) {
	loaded := map[string]bool{}
	deps := map[string]bool{}
	shouldLoadExts := make([]string, len(argv.Extensions))
	copy(shouldLoadExts, argv.Extensions)
	for len(shouldLoadExts) > 0 {
		var exts []build.Extension
		exts, err = build.LoadExtensions(extensionsDir, shouldLoadExts)
		if err != nil {
			log.Error().
				Printf("load extensions error: %v", err)
			break
		}
		extensions = append(extensions, exts...)
		// load deps
		shouldLoadExts = shouldLoadExts[0:0]
		for _, ext := range exts {
			loaded[ext.Path] = true
		}
		for _, ext := range exts {
			for _, dep := range ext.Deps {
				if !loaded[dep] && !deps[dep] {
					shouldLoadExts = append(shouldLoadExts, dep)
				}
				deps[dep] = true
			}
		}
	}
	return
}
