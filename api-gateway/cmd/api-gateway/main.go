package main

import (
	"github.com/mrevds/pizza-app/api-gateway/internal/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(app.Register()).Run()
}
