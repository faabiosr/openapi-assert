package assert

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type (
	BaseTestSuite struct {
		suite.Suite

		assert   *assert.Assertions
		filePath string

		doc Document
	}
)

func (s *BaseTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.filePath = "./fixtures/docs.json"
	s.doc, _ = LoadFromURI(s.filePath)
}
