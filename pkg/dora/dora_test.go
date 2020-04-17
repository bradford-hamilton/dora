package dora

import (
	"fmt"
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
				{accessType: ObjectAccess, keyReq: "item1"},
				{accessType: ArrayAccess, indexReq: 2},
				{accessType: ObjectAccess, keyReq: "innerKey"},
			},
		},
		{
			[]rune("$[25].item3"),
			[]queryToken{
				{accessType: ArrayAccess, indexReq: 25},
				{accessType: ObjectAccess, keyReq: "item3"},
			},
		},
		{
			[]rune("$[7].item4.innerKey"),
			[]queryToken{
				{accessType: ArrayAccess, indexReq: 7},
				{accessType: ObjectAccess, keyReq: "item4"},
				{accessType: ObjectAccess, keyReq: "innerKey"},
			},
		},
		{
			[]rune("$.item1[2].innerKey.anotherValue"),
			[]queryToken{
				{accessType: ObjectAccess, keyReq: "item1"},
				{accessType: ArrayAccess, indexReq: 2},
				{accessType: ObjectAccess, keyReq: "innerKey"},
				{accessType: ObjectAccess, keyReq: "anotherValue"},
			},
		},
		{
			[]rune("$[0].item1[2].coolKey.neatValue[16]"),
			[]queryToken{
				{accessType: ArrayAccess, indexReq: 0},
				{accessType: ObjectAccess, keyReq: "item1"},
				{accessType: ArrayAccess, indexReq: 2},
				{accessType: ObjectAccess, keyReq: "coolKey"},
				{accessType: ObjectAccess, keyReq: "neatValue"},
				{accessType: ArrayAccess, indexReq: 16},
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
				t.Fatalf("Expected access type of %d, got: %d", tt.expectedToken[i].accessType, tok.accessType)
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

func TestGetByPath(t *testing.T) {
	testJSON := `
{
	"data": {
		"users": [{
			"first_name": "bradford",
			"last_name": "human",
			"email": "brad@example.com",
			"confirmed": true,
			"allergies": null,
			"age": 30,
			"random_items": [true, { "dog_name": "ellie" }]
		}]
	}
}`
	tests := []struct {
		query          string
		expectedResult string
	}{
		{
			"$.data.users[0].first_name",
			"bradford",
		},
		{
			"$.data.users[0].confirmed",
			"true",
		},
		{
			"$.data.users[0].allergies",
			"null",
		},
		{
			"$.data.users[0].age",
			"30",
		},
		{
			"$.data.users[0].random_items",
			"[true, { \"dog_name\": \"ellie\" }]",
		},
		{
			"$.data.users[0].random_items[1]",
			"{ \"dog_name\": \"ellie\" }",
		},
	}

	for _, tt := range tests {
		c, err := NewFromString(testJSON)
		if err != nil {
			t.Fatalf("\nError creating client: %v\n", err)
		}

		result, err := c.GetByPath(tt.query)
		if err != nil {
			fmt.Println(err)
		}

		if result != tt.expectedResult {
			t.Fatalf("Expected result type of %s, got: %s", tt.expectedResult, result)
		}
	}
}
