package main

import (
	"github.com/LeBaoTai/SDN-RD/internal/controller/controller"
)

func main() {
	controller := controller.NewController()
	controller.Start()
}
