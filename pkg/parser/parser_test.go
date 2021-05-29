package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bradford-hamilton/dora/pkg/ast"
	"github.com/bradford-hamilton/dora/pkg/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParsingJSONObjectChildren(t *testing.T) {
	tests := [...]struct {
		input       string
		childrenLen int
	}{
		{input: "{\"key0\": \"value0\"}", childrenLen: 1},
		{input: "{\"key0\": \"value0\" }", childrenLen: 1},
		{input: "{\"key1\": \"value1\", \"key2\": \"value2\"}", childrenLen: 2},
		{input: "{\"key3\": [\"value3\", \"value4\"]}", childrenLen: 1},
		{input: "{\"key4\": [\"value5\", {\"key5\": \"value6\"}]}", childrenLen: 1},
		{input: "{\"key5\":\" value7\", \"key6\": \"value7\"}", childrenLen: 2},
		{input: "{\"key5\":\" value7\", \"key6\": \"value7\", \"key7\": \"value8\"}", childrenLen: 3},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseJSON()
		if err != nil {
			t.Fatalf("Failed to parse program. Error: %v", err)
		}

		rv := *program.RootValue
		val := rv.Content.(ast.Object)

		checkParserErrors(t, p)

		if len(val.Children) != tt.childrenLen {
			t.Fatalf("The length of the children does not contain 1 statement. Got: %d", len(val.Children))
		}
	}
}

func TestParsingJSONArrayChildren(t *testing.T) {
	tests := [...]struct {
		input       string
		childrenLen int
	}{
		{input: "[\"item1\"]", childrenLen: 1},
		{input: "[\"item2\", \"item3\"]", childrenLen: 2},
		{input: "[\"item4\", \"item5\", \"item6\"]", childrenLen: 3},
		{input: "[\"item7\", \"item8\", {\"key1\": \"value1\"}]", childrenLen: 3},
		{input: "[\"item9\", \"item10\", {\"key1\": \"value1\"}, \"item11\"]", childrenLen: 4},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseJSON()
		if err != nil {
			t.Fatalf("Failed to parse program. Error: %v", err)
		}

		rv := *program.RootValue
		val := rv.Content.(ast.Array)

		checkParserErrors(t, p)

		if len(val.Children) != tt.childrenLen {
			t.Fatalf("The length of the children does not contain 1 statement. Got: %d", len(val.Children))
		}
	}
}

func TestParseAndWriteFull(t *testing.T) {
	input := `// Initial comment
{
	"items": {
		"item": [{
			'id': '0001', // <-- using single quotes here
			"type": "donut",
			"name": "Cake",
			// Add
			// comments
			// again
			"cpu": 55,
			"batters"    : {
				"batter": [{
					"id": false,
					"name": null,
					"fun": true
				}]
			},
			"names": ["catstack", "lampcat", "langlang"]
		}]
		/* Add
		 * a
		 * block
		 * comment
		*/
	},
	"version": 0.1,
	"number": 11.4,
	"negativeNum": -5,
	"escapeString": "I'm some \"string\" thats escaped"
}`
	rewritten, err := parseAndOutputString(input)
	assert.NoError(t, err)

	assert.Equal(t, input, rewritten)
}

func TestParseAndWriteMinimal(t *testing.T) {
	input := `{}`
	rewritten, err := parseAndOutputString(input)
	assert.NoError(t, err)

	assert.Equal(t, input, rewritten)
}
func TestParseAndWriteObjectWithSingleProperty(t *testing.T) {
	input := `{
		"prop1"  : "value1"
	}`
	rewritten, err := parseAndOutputString(input)
	if assert.NoError(t, err) {
		assert.Equal(t, input, rewritten)
	}
}
func TestParseAndWriteObjectWithFloatProperty(t *testing.T) {
	input := `{
		"prop1" : 1.23,
		"prop2" : 1.234567890 	
	}`
	rewritten, err := parseAndOutputString(input)
	if assert.NoError(t, err) {
		assert.Equal(t, input, rewritten)
	}
}
func TestParseAndWriteObjectWithIntProperty(t *testing.T) {
	input := `{
		"prop1" : 123,
	}`
	rewritten, err := parseAndOutputString(input)
	if assert.NoError(t, err) {
		assert.Equal(t, input, rewritten)
	}
}
func TestParseAndWriteObjectWithSingleQuotedProperty(t *testing.T) {
	input := `{
		'prop1'  : 'value1'
	}`
	rewritten, err := parseAndOutputString(input)
	if assert.NoError(t, err) {
		assert.Equal(t, input, rewritten)
	}
}
func TestParseAndWriteObjectWithMultipleProperties(t *testing.T) {
	input := `{
		"prop1"  : "value1",
		"prop2"  : "value2"   ,   
		"prop3"  : "value3"
	}`
	rewritten, err := parseAndOutputString(input)
	if assert.NoError(t, err) {
		assert.Equal(t, input, rewritten)
	}
}
func TestParseAndWriteObjectWithTrailingComma(t *testing.T) {
	input := `{
		"prop1"  : "value1",
	}`
	rewritten, err := parseAndOutputString(input)
	if assert.NoError(t, err) {
		assert.Equal(t, input, rewritten)
	}
}
func TestParseAndWriteObjectWithArray(t *testing.T) {
	input := `{
		"prop1"  : [
			"one",
			"two",
			"three"
		]
	}`
	rewritten, err := parseAndOutputString(input)
	if assert.NoError(t, err) {
		assert.Equal(t, input, rewritten)
	}
}

func TestParseAndWriteObjectWithArrayWithTrailingComma(t *testing.T) {
	input := `{
		"prop1"  : [
			"one",
			"two",
			"three",
		]
	}`
	rewritten, err := parseAndOutputString(input)
	if assert.NoError(t, err) {
		assert.Equal(t, input, rewritten)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}

func parseAndOutputString(input string) (string, error) {
	l := lexer.New(input)
	p := New(l)
	j, err := p.ParseJSON()
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	if err := appendValue(&builder, *j.RootValue); err != nil {
		return "", err
	}

	return builder.String(), nil
}
func appendValue(builder *strings.Builder, item ast.Value) error {
	if err := appendStructure(builder, item.PrefixStructure); err != nil {
		return err
	}
	if err := appendValueContent(builder, item.Content); err != nil {
		return err
	}
	if err := appendStructure(builder, item.SuffixStructure); err != nil {
		return err
	}
	return nil
}
func appendValueContent(builder *strings.Builder, item ast.ValueContent) error {
	switch valueTyped := item.(type) {
	case ast.Object:
		return appendObject(builder, valueTyped)
	case ast.Array:
		return appendArray(builder, valueTyped)
	case ast.Literal:
		return appendLiteral(builder, valueTyped)
	case ast.Value:
		return appendValue(builder, valueTyped)
	default:
		return fmt.Errorf("Unhandled type in appendValue: %T", valueTyped)
	}
}
func appendStructure(builder *strings.Builder, items []ast.StructuralItem) error {
	for _, item := range items {
		if _, err := builder.WriteString(item.Value); err != nil {
			return err
		}
	}
	return nil
}
func appendObject(builder *strings.Builder, item ast.Object) error {
	if _, err := builder.WriteString("{"); err != nil {
		return err
	}

	for _, child := range item.Children {
		if err := appendProperty(builder, child); err != nil {
			return err
		}
	}

	if err := appendStructure(builder, item.SuffixStructure); err != nil {
		return err
	}
	if _, err := builder.WriteString("}"); err != nil {
		return err
	}

	return nil
}
func appendProperty(builder *strings.Builder, item ast.Property) error {
	if err := appendStructure(builder, item.PrefixStructure); err != nil {
		return err
	}
	if err := appendIdentifier(builder, item.Key); err != nil {
		return err
	}
	if err := appendStructure(builder, item.PostKeyStructure); err != nil {
		return err
	}
	if _, err := builder.WriteString(":"); err != nil {
		return err
	}
	if err := appendStructure(builder, item.PreValueStructure); err != nil {
		return err
	}
	if err := appendValueContent(builder, item.Value); err != nil {
		return err
	}
	if err := appendStructure(builder, item.PostValueStructure); err != nil {
		return err
	}
	if item.HasCommaSeparator {
		if _, err := builder.WriteString(","); err != nil {
			return err
		}
	}
	return nil
}
func appendIdentifier(builder *strings.Builder, item ast.Identifier) error {
	valueToWrite := fmt.Sprintf("%s%s%s", item.Delimiter, item.Value, item.Delimiter)
	if _, err := builder.WriteString(valueToWrite); err != nil {
		return err
	}
	return nil
}
func appendLiteral(builder *strings.Builder, item ast.Literal) error {
	var valueToWrite string
	if item.OriginalRendering != "" {
		valueToWrite = item.OriginalRendering
	} else {
		switch item.ValueType {
		case ast.StringLiteralValueType:
			valueToWrite = fmt.Sprintf("%s%s%s", item.Delimiter, item.Value.(string), item.Delimiter)
		case ast.BooleanLiteralValueType:
			valueToWrite = fmt.Sprintf("%t", item.Value.(bool))
		case ast.NullLiteralValueType:
			valueToWrite = "null"
		case ast.NumberLiteralValueType:
			valueToWrite = fmt.Sprintf("%v", item.Value)
		default:
			return fmt.Errorf("Unhandled Literal Value Type: %v", item.ValueType)
		}
	}

	if _, err := builder.WriteString(valueToWrite); err != nil {
		return err
	}
	return nil
}

func appendArray(builder *strings.Builder, item ast.Array) error {
	if err := appendStructure(builder, item.PrefixStructure); err != nil {
		return err
	}
	if _, err := builder.WriteString("["); err != nil {
		return err
	}
	for _, arrayItem := range item.Children {
		if err := appendArrayItem(builder, arrayItem); err != nil {
			return err
		}
	}
	if err := appendStructure(builder, item.SuffixStructure); err != nil {
		return err
	}
	if _, err := builder.WriteString("]"); err != nil {
		return err
	}
	return nil
}
func appendArrayItem(builder *strings.Builder, item ast.ArrayItem) error {
	if err := appendStructure(builder, item.PrefixStructure); err != nil {
		return err
	}
	if err := appendValueContent(builder, item.Value); err != nil {
		return err
	}
	if err := appendStructure(builder, item.PostValueStructure); err != nil {
		return err
	}
	if item.HasCommaSeparator {
		if _, err := builder.WriteString(","); err != nil {
			return err
		}
	}
	return nil
}
