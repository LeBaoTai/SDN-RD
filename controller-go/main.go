package main

import "github.com/LeBaoTai/myco-controller/internal/controller"

type ConfigPath struct {
	Paths []string `yaml:"paths"`
}

func main() {
	controller := controller.NewController()
	controller.Start()
}
