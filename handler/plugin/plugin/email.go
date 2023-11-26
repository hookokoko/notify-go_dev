package main

import (
	"context"
	"fmt"
)

var (
	Name        = "email"
	Version     = "0.0.1"
	Description = "email plugin"
)

// go build -ldflags "-pluginpath=github.com/ecodeclub/notify-go/experience/handler/plugin/plugin" -buildmode=plugin -o email.so

type PluginFunc interface {
	GetConf(name string, value interface{}) error
	GetLogger() interface{}
}

func Execute(f context.Context) error {
	pf := f.Value("plugin").(PluginFunc)

	fmt.Println("get conf, ", pf.GetConf("account", 12345))
	fmt.Println("get logger, ", pf.GetLogger())

	fmt.Println("execute ok")
	return nil
}
