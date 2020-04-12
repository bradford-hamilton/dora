package dora

import (
	"testing"
)

func TestScanQueryTokens(t *testing.T) {
	tests := []struct {
		input         []rune
		expectedToken []queryToken
	}{
		{
			[]rune("$.item1[2].innerKey"),
			[]queryToken{
				{accessType: "object", keyReq: "item1"},
				{accessType: "array", indexReq: 2},
				{accessType: "object", keyReq: "innerKey"},
			},
		},
		{
			[]rune("$[25].item3"),
			[]queryToken{
				{accessType: "array", indexReq: 25},
				{accessType: "object", keyReq: "item3"},
			},
		},
		{
			[]rune("$[7].item4.innerKey"),
			[]queryToken{
				{accessType: "array", indexReq: 7},
				{accessType: "object", keyReq: "item4"},
				{accessType: "object", keyReq: "innerKey"},
			},
		},
		{
			[]rune("$.item1[2].innerKey.anotherValue"),
			[]queryToken{
				{accessType: "object", keyReq: "item1"},
				{accessType: "array", indexReq: 2},
				{accessType: "object", keyReq: "innerKey"},
				{accessType: "object", keyReq: "anotherValue"},
			},
		},
		{
			[]rune("$[0].item1[2].coolKey.neatValue[16]"),
			[]queryToken{
				{accessType: "array", indexReq: 0},
				{accessType: "object", keyReq: "item1"},
				{accessType: "array", indexReq: 2},
				{accessType: "object", keyReq: "coolKey"},
				{accessType: "object", keyReq: "neatValue"},
				{accessType: "array", indexReq: 16},
			},
		},
	}

	for _, tt := range tests {
		tokens, err := scanQueryTokens(tt.input)
		if err != nil {
			t.Fatalf("Failed to scan tokens. Error: %v", err)
		}

		for i, tok := range tokens {
			if tok.accessType != tt.expectedToken[i].accessType {
				t.Fatalf("Expected access type of %s, got: %s", tt.expectedToken[i].accessType, tok.accessType)
			}
			if tok.keyReq != tt.expectedToken[i].keyReq {
				t.Fatalf("Expected keyReq of %s, got: %s", tt.expectedToken[i].keyReq, tok.keyReq)
			}
			if tok.indexReq != tt.expectedToken[i].indexReq {
				t.Fatalf("Expected indexReq of %d, got: %d", tt.expectedToken[i].indexReq, tok.indexReq)
			}
		}
	}
}
