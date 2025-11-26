package main

import (
	"context"

	"github.com/moremoneymod/auth/internal/app"
	desc "github.com/moremoneymod/auth/pkg/auth_v1"
)

type server struct {
	desc.UnimplementedAuthV1Server
}

func main() {
	ctx := context.Background()
	app, err := app.NewApp(ctx)
	if err != nil {
		panic(err)
	}
	err = app.Run()
	if err != nil {
		panic(err)
	}
}
