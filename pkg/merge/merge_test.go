package merge

import (
	"testing"

	"github.com/bradford-hamilton/dora/pkg/ast"
	"github.com/bradford-hamilton/dora/pkg/lexer"
	"github.com/bradford-hamilton/dora/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestMergeSimpleObjectNoConflicts(t *testing.T) {

	baseInput := `{
	"prop1" : "Hello"

}`
	newInput := `{
	"prop2" : "World"


}`
	expectedOutput := `{
	"prop1" : "Hello",
	"prop2" : "World"


}`

	testMerge(t, baseInput, newInput, expectedOutput)
}
func TestMergeSimpleObjectTrailingCommaOnBaseNoConflicts(t *testing.T) {

	baseInput := `{
	"prop1" : "Hello",

}`
	newInput := `{
	"prop2" : "World"


}`
	expectedOutput := `{
	"prop1" : "Hello",
	"prop2" : "World"


}`

	testMerge(t, baseInput, newInput, expectedOutput)
}
func TestMergeSimpleObjectTrailingCommaOnMergeNoConflicts(t *testing.T) {

	baseInput := `{
	"prop1" : "Hello"

}`
	newInput := `{
	"prop2" : "World",


}`
	expectedOutput := `{
	"prop1" : "Hello",
	"prop2" : "World",


}`

	testMerge(t, baseInput, newInput, expectedOutput)
}
func TestMergeSimpleObjectTrailingCommaOnBothNoConflicts(t *testing.T) {

	baseInput := `{
	"prop1" : "Hello",
}`
	newInput := `{
	"prop2" : "World",


}`
	expectedOutput := `{
	"prop1" : "Hello",
	"prop2" : "World",


}`

	testMerge(t, baseInput, newInput, expectedOutput)
}

func TestMergeSimpleObjectWithConflict(t *testing.T) {

	baseInput := `{
	"prop1" : "Hello"
}`
	newInput := `{
	"prop1" : "Goodbye"
}`
	expectedOutput := `{
	"prop1" : "Goodbye"
}`

	testMerge(t, baseInput, newInput, expectedOutput)
}

func testMerge(t *testing.T, baseInput string, newInput string, expectedOutput string) {

	l := lexer.New(baseInput)
	p := parser.New(l)
	baseDocument, err := p.ParseJSON()
	if !assert.NoError(t, err) {
		return
	}

	l = lexer.New(newInput)
	p = parser.New(l)
	newDocument, err := p.ParseJSON()
	if !assert.NoError(t, err) {
		return
	}

	mergedDocument, err := MergeJSON(baseDocument, newDocument)
	if !assert.NoError(t, err) {
		return
	}

	mergedJSON, err := ast.WriteJSONString(mergedDocument)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, expectedOutput, mergedJSON)
}
