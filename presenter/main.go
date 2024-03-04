package main

import "github.com/mindwingx/graph-processor/bootstrap"

func main() {
	service := bootstrap.NewApp()
	service.Init()
	service.Start()
}
