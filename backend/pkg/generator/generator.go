package generator

import "github.com/google/uuid"

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateId() string {
	return uuid.NewString()
}
