//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"errors"

	"github.com/google/wire"
)

type Foo struct {
	X int
}

func ProvideFoo() Foo {
	return Foo{X: 42}
}

type Bar struct {
	X int
}

func ProvideBar(foo Foo) Bar {
	return Bar{X: -foo.X}
}

type Baz struct {
	X int
}

// ProvideBaz returns a value if Bar is not Zero
func ProvideBaz(ctx context.Context, bar Bar) (Baz, error) {
	if bar.X == 0 {
		return Baz{}, errors.New("cannot provide baz when bar is zero")
	}
	return Baz{X: bar.X}, nil
}

// Providers can be grouped into provider sets. This is useful if several providers
// will frequently be used together.

var SuperSet = wire.NewSet(ProvideFoo, ProvideBar, ProvideBaz)

func initializeBar(ctx context.Context) (Baz, error) {
	wire.Build(SuperSet)
	return Baz{}, nil
}
