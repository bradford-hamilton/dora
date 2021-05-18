package ast

import (
	"fmt"
	"io"
	"strings"
)

// JSONWriter provides the ability to write an AST representation to an io.Writer
type JSONWriter struct {
	writer io.Writer
}

// NewJSONWriter
func NewJSONWriter(writer io.Writer) *JSONWriter {
	return &JSONWriter{
		writer: writer,
	}
}

// WriteJSONString returns the string representation of the JSON in rootNode
func WriteJSONString(rootNode *RootNode) (string, error) {

	var builder strings.Builder
	j := NewJSONWriter(&builder)

	if err := j.appendValue(*rootNode.RootValue); err != nil {
		return "", err
	}

	return builder.String(), nil
}

func (j *JSONWriter) appendValue(item Value) error {
	if err := j.appendStructure(item.PrefixStructure); err != nil {
		return err
	}
	if err := j.appendValueContent(item.Content); err != nil {
		return err
	}
	if err := j.appendStructure(item.SuffixStructure); err != nil {
		return err
	}
	return nil
}
func (j *JSONWriter) appendValueContent(item ValueContent) error {
	switch valueTyped := item.(type) {
	case Object:
		return j.appendObject(valueTyped)
	case Array:
		return j.appendArray(valueTyped)
	case Literal:
		return j.appendLiteral(valueTyped)
	case Value:
		return j.appendValue(valueTyped)
	default:
		return fmt.Errorf("unhandled type in appendValue: %T", valueTyped)
	}
}
func (j *JSONWriter) appendStructure(items []StructuralItem) error {
	for _, item := range items {
		if _, err := fmt.Fprint(j.writer, item.Value); err != nil {
			return err
		}
	}
	return nil
}
func (j *JSONWriter) appendObject(item Object) error {
	if _, err := fmt.Fprint(j.writer, "{"); err != nil {
		return err
	}

	for _, child := range item.Children {
		if err := j.appendProperty(child); err != nil {
			return err
		}
	}

	if err := j.appendStructure(item.SuffixStructure); err != nil {
		return err
	}
	if _, err := fmt.Fprint(j.writer, "}"); err != nil {
		return err
	}

	return nil
}
func (j *JSONWriter) appendProperty(item Property) error {
	if err := j.appendIdentifier(item.Key); err != nil {
		return err
	}
	if _, err := fmt.Fprint(j.writer, ":"); err != nil {
		return err
	}
	if err := j.appendValue(item.Value); err != nil {
		return err
	}
	if item.HasCommaSeparator {
		if _, err := fmt.Fprint(j.writer, ","); err != nil {
			return err
		}
	}
	return nil
}
func (j *JSONWriter) appendIdentifier(item Identifier) error {
	if err := j.appendStructure(item.PrefixStructure); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(j.writer, "%s%s%s", item.Delimiter, item.Value, item.Delimiter); err != nil {
		return err
	}
	if err := j.appendStructure(item.SuffixStructure); err != nil {
		return err
	}
	return nil
}
func (j *JSONWriter) appendLiteral(item Literal) error {
	var valueToWrite string
	if item.OriginalRendering != "" {
		valueToWrite = item.OriginalRendering
	} else {
		switch item.ValueType {
		case StringLiteralValueType:
			valueToWrite = fmt.Sprintf("%s%s%s", item.Delimiter, item.Value.(string), item.Delimiter)
		case BooleanLiteralValueType:
			valueToWrite = fmt.Sprintf("%t", item.Value.(bool))
		case NullLiteralValueType:
			valueToWrite = "null"
		case NumberLiteralValueType:
			valueToWrite = fmt.Sprintf("%v", item.Value)
		default:
			return fmt.Errorf("unhandled Literal Value Type: %v", item.ValueType)
		}
	}

	if _, err := fmt.Fprint(j.writer, valueToWrite); err != nil {
		return err
	}
	return nil
}

func (j *JSONWriter) appendArray(item Array) error {
	if err := j.appendStructure(item.PrefixStructure); err != nil {
		return err
	}
	if _, err := fmt.Fprint(j.writer, "["); err != nil {
		return err
	}
	for _, arrayItem := range item.Children {
		if err := j.appendArrayItem(arrayItem); err != nil {
			return err
		}
	}
	if err := j.appendStructure(item.SuffixStructure); err != nil {
		return err
	}
	if _, err := fmt.Fprint(j.writer, "]"); err != nil {
		return err
	}
	return nil
}
func (j *JSONWriter) appendArrayItem(item ArrayItem) error {
	if err := j.appendStructure(item.PrefixStructure); err != nil {
		return err
	}
	if err := j.appendValueContent(item.Value); err != nil {
		return err
	}
	if err := j.appendStructure(item.PostValueStructure); err != nil {
		return err
	}
	if item.HasCommaSeparator {
		if _, err := fmt.Fprint(j.writer, ","); err != nil {
			return err
		}
	}
	return nil
}
