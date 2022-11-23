package main

import (
	"fmt"

	"github.com/gopherd/log"
	"github.com/mkideal/pkg/errors"

	"github.com/midlang/mid/src/genutil"
	"github.com/midlang/mid/src/mid/build"
)

func main() {
	log.Start(log.WithSync(true), log.WithLevel(log.LevelWarn))

	plugin, config, builder, err := build.ParseFlags()
	log.If(err != nil).Fatal().
		Error("err", err).
		Print("ParseFlags error")
	log.Debug().
		Any("plugin", plugin).
		Any("config", config).
		Any("builder", builder).
		Print("running plugin")

	err = generate(builder, plugin, config)
	log.If(err != nil).Error().
		Error("err", err).
		Print("generate error")
}

func generate(builder *build.Builder, plugin build.Plugin, config build.PluginRuntimeConfig) (err error) {
	defer func() {
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

	// initialize generator
	genutil.Init(buildType, plugin, config)

	pkgs := builder.Packages
	for _, pkg := range pkgs {
		if _, err := genutil.GeneratePackage(pkg); err != nil {
			return err
		}
	}
	return nil
}
