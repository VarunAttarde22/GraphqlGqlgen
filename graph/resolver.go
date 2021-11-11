package graph

import "github.com/VarunAttarde22/hackernews/graph/model"

//go:generate go run github.com/99designs/gqlgen
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	nodes []*model.Node
	covid []*model.Covid
}
