package merge

import (
	"fmt"

	"github.com/bradford-hamilton/dora/pkg/ast"
)

func MergeJSON(baseDocument ast.RootNode, mergeDocument ast.RootNode) (*ast.RootNode, error) {

	result := baseDocument

	newContent, err := mergeValues(*result.RootValue, *mergeDocument.RootValue, "$")
	if err != nil {
		return nil, err
	}
	result.RootValue = &newContent
	return &result, nil
}

func mergeValues(baseValue ast.Value, mergeValue ast.Value, currentPath string) (ast.Value, error) {

	result := baseValue
	switch resultContent := (baseValue.Content).(type) {
	case ast.Object:
		switch mergeContent := mergeValue.Content.(type) {
		case ast.Object:
			for _, mergeChild := range mergeContent.Children {
				resultChildIndex, resultChild := getChildByKey(resultContent, mergeChild.Key.Value)
				if resultChild == nil {
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
					resultChild, err := mergeValues(resultChild.Value, mergeChild.Value, currentPath+"."+mergeChild.Key.Value)
					if err != nil {
						return ast.Value{}, err
					}
					resultContent.Children[resultChildIndex].Value = resultChild
				}
			}
			result.Content = resultContent
			return result, nil
		default:
			return ast.Value{}, fmt.Errorf("mis-matched types at %q. base type: %T, merge type: %T", currentPath, resultContent, mergeContent)
		}
	case ast.Literal:
		return mergeValue, nil
	default:
		return ast.Value{}, fmt.Errorf("unhandled type at %q. base type: %T", currentPath, resultContent)
	}
}

func getChildByKey(object ast.Object, key string) (int, *ast.Property) {
	for index, child := range object.Children {
		if child.Key.Value == key {
			return index, &child
		}
	}
	return -1, nil
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
