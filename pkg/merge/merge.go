package merge

import (
	"fmt"

	"github.com/bradford-hamilton/dora/pkg/ast"
)

func MergeJSON(baseDocument ast.RootNode, mergeDocument ast.RootNode) (*ast.RootNode, error) {

	result := baseDocument

	newContent, err := mergeValueContent(*result.RootValue, *mergeDocument.RootValue, "$")
	if err != nil {
		return nil, err
	}
	result.RootValue = &newContent
	return &result, nil
}

func mergeValueContent(baseValue ast.Value, mergeValue ast.Value, currentPath string) (ast.Value, error) {

	result := baseValue
	switch resultContent := (baseValue.Content).(type) {
	case ast.Object:
		switch mergeContent := mergeValue.Content.(type) {
		case ast.Object:
			for _, mergeChild := range mergeContent.Children {
				baseChild := getChildByKey(resultContent, mergeChild.Key.Value)
				if baseChild == nil {
					lastChildIndex := len(resultContent.Children) - 1
					if resultContent.Children[lastChildIndex].HasCommaSeparator {
						resultContent.SuffixStructure = append(stripWhiteSpace(resultContent.SuffixStructure), mergeContent.SuffixStructure...)
					} else {
						// Add in comma
						resultContent.Children[lastChildIndex].HasCommaSeparator = true
						resultContent.Children[lastChildIndex].Value.SuffixStructure = stripWhiteSpace(resultContent.Children[lastChildIndex].Value.SuffixStructure)
						if mergeChild.HasCommaSeparator {
							resultContent.SuffixStructure = append(stripWhiteSpace(resultContent.SuffixStructure), mergeContent.SuffixStructure...)
						}
					}
					resultContent.Children = append(resultContent.Children, mergeChild)
				} else {
					// TODO - handle merging object properties
				}
			}
			result.Content = resultContent
			return result, nil
		default:
			return ast.Value{}, fmt.Errorf("mis-matched types at %q. base type: %T, merge type: %T", currentPath, resultContent, mergeContent)
		}
	default:
		return ast.Value{}, fmt.Errorf("unhandled type at %q. base type: %T", currentPath, resultContent)
	}
}

func getChildByKey(object ast.Object, key string) *ast.Property {
	for _, child := range object.Children {
		if child.Key.Value == key {
			return &child
		}
	}
	return nil
}

func stripWhiteSpace(structuralItems []ast.StructuralItem) []ast.StructuralItem {
	var lastNonWhitespaceIndex int
	for lastNonWhitespaceIndex := len(structuralItems) - 1; lastNonWhitespaceIndex >= 0; lastNonWhitespaceIndex-- {
		if structuralItems[lastNonWhitespaceIndex].ItemType != ast.WhitespaceStructuralItemType {
			break
		}
	}
	return structuralItems[0:lastNonWhitespaceIndex]
}
