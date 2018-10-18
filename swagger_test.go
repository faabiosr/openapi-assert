package assert

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type (
	SwaggerTestSuite struct {
		BaseTestSuite
	}
)

func (s *SwaggerTestSuite) TestLoadFromUriWithEmptyParam() {
	doc, err := LoadFromURI("")

	s.assert.Nil(doc)
	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrSwaggerLoad)
}

func (s *SwaggerTestSuite) TestLoadFromUriWithInvalidFile() {
	doc, err := LoadFromURI("./fixtures/invalid-doc.json")

	s.assert.Nil(doc)
	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrSwaggerExpand)
}

func (s *SwaggerTestSuite) TestLoadFromUri() {
	doc, err := LoadFromURI(s.filePath)

	s.assert.IsType(&Swagger{}, doc)
	s.assert.Implements(new(Document), doc)
	s.assert.NoError(err)
}

func (s *SwaggerTestSuite) TestFindPathWithBrokenDocument() {
	doc, _ := LoadFromURI("./fixtures/invalid-path.json")

	path, err := doc.FindPath("/api/food/a")

	s.assert.Empty(path)
	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrResourceURI)
}

func (s *SwaggerTestSuite) TestFindPathWithInvalidPath() {
	path, err := s.doc.FindPath("/some")

	s.assert.Empty(path)
	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrResourceURI)
}

func (s *SwaggerTestSuite) TestFindPath() {
	path, err := s.doc.FindPath("/api/pets")

	s.assert.NotEmpty(path)
	s.assert.NoError(err)
}

func (s *SwaggerTestSuite) TestFindNodeWithoutSegment() {
	node, err := s.doc.FindNode("")

	s.assert.Empty(node)
	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrNodeNotFound)
}

func (s *SwaggerTestSuite) TestFindNode() {
	node, err := s.doc.FindNode("paths")

	s.assert.NotEmpty(node)
	s.assert.NoError(err)
}

func TestSwaggerTestSuite(t *testing.T) {
	suite.Run(t, new(SwaggerTestSuite))
}
