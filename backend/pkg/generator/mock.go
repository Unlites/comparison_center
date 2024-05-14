package generator

import "github.com/stretchr/testify/mock"

type MockGenerator struct {
	mock.Mock
}

func NewMockGenerator() *MockGenerator {
	return &MockGenerator{}
}

func (g *MockGenerator) GenerateId() string {
	args := g.Called()

	return args.String(0)
}
