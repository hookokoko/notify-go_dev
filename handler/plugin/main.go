package main

import (
	"context"
	"errors"
	"fmt"
	"plugin"
)

type PluginContext struct {
}

func (p *PluginContext) GetConf(name string, value interface{}) error {
	fmt.Printf("main_conf_name: %s, main_conf_value: %v. \n", name, value)
	return nil
}

func (p *PluginContext) GetLogger() interface{} {
	fmt.Printf("main_logger. \n")
	return "get logger conf"
}

type PluginChannel struct {
	Name        string                          `json:"name"`
	FuncExecute func(ctx context.Context) error `json:"-"`
	Err         error                           `json:"-"`
}

func NewPluginChannel(name string) *PluginChannel {
	p := &PluginChannel{Name: name}
	err := p.loadSo(fmt.Sprintf("plugin/plugin/%s.so", name))
	if err != nil {
		panic(fmt.Errorf("[%s]load plugin error: %v.\n", name, err))
	}
	return p
}

func (pc *PluginChannel) loadSo(path string) error {
	pl, openErr := plugin.Open(path)
	if openErr != nil {
		return openErr
	}

	sym, lookErr := pl.Lookup("Execute")
	if lookErr != nil {
		return lookErr
	}

	execute, ok := sym.(func(ctx context.Context) error)
	if !ok {
		fmt.Printf("%+v\n", sym)
		return errors.New("execute is not type Func_Execute")
	}
	pc.FuncExecute = execute
	return nil
}

func main() {
	pc := NewPluginChannel("email")
	err := pc.loadSo(fmt.Sprintf("plugin/plugin/%s.so", pc.Name))
	if err != nil {
		fmt.Printf("load plugin error: %v. \n", err)
		return
	}
	ctx := context.WithValue(context.TODO(), "plugin", &PluginContext{})
	err = pc.FuncExecute(ctx)
	if err != nil {
		fmt.Printf("execute plugin error: %v. \n", err)
		return
	}
}
