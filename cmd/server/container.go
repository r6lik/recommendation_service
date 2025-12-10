package main

import "go.uber.org/fx"

func BuildContainer() *fx.App {
	return fx.New()
}
