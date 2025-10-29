package main

import (
	fxApp "github.com/lkgiovani/growth_technical_challenge/infra/fx"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fxApp.Module,
	).Run()
}
