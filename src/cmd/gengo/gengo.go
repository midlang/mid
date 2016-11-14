package main

import (
	"fmt"
	"path/filepath"

	"github.com/midlang/mid/src/mid/build"
	"github.com/midlang/mid/src/tools/generator"
	"github.com/mkideal/log"
	"github.com/mkideal/pkg/errors"
)

func main() {
	defer log.Uninit(log.InitColoredConsole(log.LvWARN))

	plugin, config, builder, err := build.ParseFlags()
	log.If(err != nil).Fatal("ParseFlags: %v", err)

	level := log.SetLevelFromString(config.Verbose)
	if !level.MoreVerboseThan(log.LvWARN) {
		log.NoHeader()
	}
	log.WithJSON(plugin).Debug("plugin")
	log.WithJSON(config).Debug("config")
	log.WithJSON(builder).Trace("builder")

	err = generate(builder, plugin, config)
	log.If(err != nil).Error("generate go: %v", err)
}

func generate(builder *build.Builder, plugin build.Plugin, config build.PluginRuntimeConfig) (err error) {
	defer func() {
		return
		if e := recover(); e != nil {
			switch x := e.(type) {
			case error:
				err = x
			case string:
				err = errors.Error(x)
			default:
				err = fmt.Errorf("%v", x)
			}
		}
	}()

	infos, err := generator.OpenTemplatesDir(plugin.Lang, plugin.TemplatesDir)
	if err != nil {
		return errors.Throw(err.Error())
	}
	if len(infos) == 0 {
		return nil
	}

	pkgs := builder.Packages
	for _, pkg := range pkgs {
		outdir := filepath.Join(config.Outdir, pkg.Name)
		if files, err := generator.GeneratePackage(outdir, pkg, plugin, config, buildType, infos); err != nil {
			return err
		} else {
			for file := range files {
				if err = generator.GoFmt(file); err != nil {
					return errors.Throw("gofmt file `" + file + "` error: " + err.Error())
				}
			}
		}
	}
	return nil
}
