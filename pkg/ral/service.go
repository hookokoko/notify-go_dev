package ral

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Service struct {
	Resources []Resource `toml:"Resources"`
}

type Resource struct {
	Name         string      `toml:"Name"`
	ConnTimeOut  int         `toml:"ConnTimeOut"`
	WriteTimeOut int         `toml:"WriteTimeOut"`
	ReadTimeOut  int         `toml:"ReadTimeOut"`
	Retry        int         `toml:"Retry"`
	Protocol     string      `toml:"Protocol"`
	Converter    string      `toml:"Converter"`
	Interface    []Interface `toml:"Interface"`
}

type Interface struct {
	Name   string `toml:"Name"`
	URL    string `toml:"Url"`
	Method string `toml:"Method"`
	Host   string `toml:"Host"`
	Port   string `toml:"Port"`
}

func NewService(file string) Service {
	service := new(Service)
	if _, err := toml.DecodeFile(file, service); err != nil {
		log.Printf("[ral] 初始化失败 %v.", err)
	}
	return *service
}
