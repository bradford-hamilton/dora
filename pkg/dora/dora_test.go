package dora

import (
	"encoding/json"
	"fmt"
	"testing"
)

const TestJSON = `
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
	},
	"codes": [200, 201, 400, 403, 404.567],
	"superNest": {
		"inner1": {
			"inner2": {
				"inner3": {
					"inner4": [{ "inner5": { "inner6": "neato" } }]
				}
			}
		}
	},
    "date": "04/19/2020",
    "enabled": true,
	"PI": 3.1415,
	"disabled": false
}`

func TestScanQueryTokens(t *testing.T) {
	tests := [...]struct {
		input         []byte
		expectedToken []queryToken
	}{
		{
			input: []byte("$.item1[2].innerKey"),
			expectedToken: []queryToken{
				{accessType: ObjectAccess, key: "item1"},
				{accessType: ArrayAccess, index: 2},
				{accessType: ObjectAccess, key: "innerKey"},
			},
		},
		{
			input: []byte("$[25].item3"),
			expectedToken: []queryToken{
				{accessType: ArrayAccess, index: 25},
				{accessType: ObjectAccess, key: "item3"},
			},
		},
		{
			input: []byte("$[7].item4.innerKey"),
			expectedToken: []queryToken{
				{accessType: ArrayAccess, index: 7},
				{accessType: ObjectAccess, key: "item4"},
				{accessType: ObjectAccess, key: "innerKey"},
			},
		},
		{
			input: []byte("$.item1[2].innerKey.anotherValue"),
			expectedToken: []queryToken{
				{accessType: ObjectAccess, key: "item1"},
				{accessType: ArrayAccess, index: 2},
				{accessType: ObjectAccess, key: "innerKey"},
				{accessType: ObjectAccess, key: "anotherValue"},
			},
		},
		{
			input: []byte("$[0].item1[2].coolKey.neatValue[16]"),
			expectedToken: []queryToken{
				{accessType: ArrayAccess, index: 0},
				{accessType: ObjectAccess, key: "item1"},
				{accessType: ArrayAccess, index: 2},
				{accessType: ObjectAccess, key: "coolKey"},
				{accessType: ObjectAccess, key: "neatValue"},
				{accessType: ArrayAccess, index: 16},
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
			if tok.key != tt.expectedToken[i].key {
				t.Fatalf("Expected key of %s, got: %s", tt.expectedToken[i].key, tok.key)
			}
			if tok.index != tt.expectedToken[i].index {
				t.Fatalf("Expected index of %d, got: %d", tt.expectedToken[i].index, tok.index)
			}
		}
	}
}

func TestClient_GetString(t *testing.T) {
	tests := [...]struct {
		query          string
		expectedResult string
	}{
		{
			query:          "$.data.users[0].first_name",
			expectedResult: "bradford",
		},
		{
			query:          "$.data.users[0].confirmed",
			expectedResult: "true",
		},
		{
			query:          "$.data.users[0].allergies",
			expectedResult: "null",
		},
		{
			query:          "$.data.users[0].age",
			expectedResult: "30.000000",
		},
		{
			query:          "$.data.users[0].random_items",
			expectedResult: "[true, { \"dog_name\": \"ellie\" }]",
		},
		{
			query:          "$.data.users[0].random_items[1]",
			expectedResult: "{ \"dog_name\": \"ellie\" }",
		},
		{
			query:          "$.codes",
			expectedResult: "[200, 201, 400, 403, 404.567]",
		},
		{
			query:          "$.codes[1]",
			expectedResult: "201.000000",
		},
		{
			query:          "$.superNest.inner1.inner2.inner3.inner4[0].inner5.inner6",
			expectedResult: "neato",
		},
		{
			query:          "$.date",
			expectedResult: "04/19/2020",
		},
	}

	for _, tt := range tests {
		c, err := NewFromString(TestJSON)
		if err != nil {
			t.Fatalf("\nError creating client: %v\n", err)
		}

		result, err := c.GetString(tt.query)
		if err != nil {
			fmt.Println(err)
		}

		if result != tt.expectedResult {
			t.Fatalf("Expected result type of %s, got: %s", tt.expectedResult, result)
		}
	}
}

func TestClient_GetBool(t *testing.T) {
	tests := [...]struct {
		query          string
		expectedResult bool
	}{
		{
			query:          "$.enabled",
			expectedResult: true,
		},
		{
			query:          "$.disabled",
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		c, err := NewFromString(TestJSON)
		if err != nil {
			t.Fatalf("\nError creating client: %v\n", err)
		}

		result, err := c.GetBool(tt.query)
		if err != nil {
			fmt.Println(err)
		}

		if result != tt.expectedResult {
			t.Fatalf("Expected result type of %T, got: %T", tt.expectedResult, result)
		}
	}
}

func TestClient_GetFloat64(t *testing.T) {
	tests := [...]struct {
		query          string
		expectedResult float64
	}{
		{
			query:          "$.PI",
			expectedResult: 3.1415,
		},
		{
			query:          "$.codes[1]",
			expectedResult: 201.000000,
		},
		{
			query:          "$.codes[4]",
			expectedResult: 404.567,
		},
	}
	for _, tt := range tests {
		c, err := NewFromString(TestJSON)
		if err != nil {
			t.Fatalf("\nError creating client: %v\n", err)
		}

		result, err := c.GetFloat64(tt.query)
		if err != nil {
			fmt.Println(err)
		}

		if result != tt.expectedResult {
			t.Fatalf("Expected result type of %f, got: %f", tt.expectedResult, result)
		}
	}
}

// func TestClient_SetString(t *testing.T) {
// 	tests := [...]struct {
// 		query     string
// 		newString string
// 	}{
// 		{
// 			query:     "$.data.users[0].first_name",
// 			newString: "randy",
// 		},
// 		{
// 			query:     "$.data.users[0].confirmed",
// 			newString: "false",
// 		},
// 		{
// 			query:     "$.data.users[0].allergies",
// 			newString: "not null",
// 		},
// 		{
// 			query:     "$.data.users[0].age",
// 			newString: "39.001",
// 		},
// 		{
// 			query:     "$.data.users[0].random_items",
// 			newString: "[false, { \"dog_name\": \"missy\" }]",
// 		},
// 		{
// 			query:     "$.data.users[0].random_items[1]",
// 			newString: "{ \"dog_name\": \"missy\" }",
// 		},
// 		{
// 			query:     "$.codes",
// 			newString: "[400, 401, 403, 404, 100.567]",
// 		},
// 		{
// 			query:     "$.codes[1]",
// 			newString: "321.04",
// 		},
// 		{
// 			query:     "$.superNest.inner1.inner2.inner3.inner4[0].inner5.inner6",
// 			newString: "burrito",
// 		},
// 		{
// 			query:     "$.date",
// 			newString: "04/20/2020",
// 		},
// 	}

// 	for _, tt := range tests {
// 		c, err := NewFromString(TestJSON)
// 		if err != nil {
// 			t.Fatalf("\nError creating client: %v\n", err)
// 		}

// 		if err := c.SetString(tt.query, tt.newString); err != nil {
// 			t.Fatalf("\nError getting string: %v\n", err)
// 		}

// 		s, err := c.GetString(tt.query)
// 		if err != nil {
// 			t.Fatalf("\nError getting string after SetString: %v\n", err)
// 		}
// 		if s != tt.newString {
// 			t.Fatalf("Expected set to overwrite the old value with %s, but it did not. Got: %s", tt.newString, s)
// 		}
// 	}
// }

// func TestClient_SetBool(t *testing.T) {
// 	tests := [...]struct {
// 		query   string
// 		newBool bool
// 	}{
// 		{
// 			query:   "$.disabled",
// 			newBool: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		c, err := NewFromString(TestJSON)
// 		if err != nil {
// 			t.Fatalf("\nError creating client: %v\n", err)
// 		}

// 		if err := c.SetBool(tt.query, tt.newBool); err != nil {
// 			t.Fatalf("\nError getting bool: %v\n", err)
// 		}

// 		b, err := c.GetBool(tt.query)
// 		if err != nil {
// 			t.Fatalf("\nError getting string after SetString: %v\n", err)
// 		}
// 		if b != tt.newBool {
// 			t.Fatalf("Expected set to overwrite the old value with %t, but it did not. Got: %t", tt.newBool, b)
// 		}
// 	}
// }

// func TestClient_SetFloat64(t *testing.T) {
// 	tests := [...]struct {
// 		query    string
// 		newFloat float64
// 	}{
// 		{
// 			query:    "$.PI",
// 			newFloat: 5.018,
// 		},
// 	}
// 	for _, tt := range tests {
// 		c, err := NewFromString(TestJSON)
// 		if err != nil {
// 			t.Fatalf("\nError creating client: %v\n", err)
// 		}

// 		if err = c.SetFloat64(tt.query, tt.newFloat); err != nil {
// 			t.Fatalf("\nError setting float64: %v\n", err)
// 		}

// 		f, err := c.GetFloat64(tt.query)
// 		if err != nil {
// 			t.Fatalf("\nError getting float64 after SetFloat64: %v\n", err)
// 		}
// 		if f != tt.newFloat {
// 			t.Fatalf("Expected set to overwrite the old value with %f, but it did not. Got: %f", tt.newFloat, f)
// 		}
// 	}
// }

// -------------------------------------- Benchmarks -------------------------------------- //
var sink string

func BenchmarkGetSingleValueWithDora(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := getSingleValueWithDora()
		sink = v
	}
}

func BenchmarkIsGetSingleValueWithUnmarshalAndSchema(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := getSingleValueWithUnmarshalAndSchema()
		sink = v
	}
}

func BenchmarkIsGetSingleValueWithUnmarshalAndNoSchema(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := getSingleValueWithUnmarshalNoSchema()
		sink = v
	}
}

func getSingleValueWithDora() string {
	c, _ := NewFromString(testJSONObject)
	r, _ := c.GetString("$.item1[2].some.thing")
	return r
}

func getSingleValueWithUnmarshalAndSchema() string {
	type testJSON struct {
		Item1 []struct {
			Some struct {
				Thing string `json:"thing"`
			} `json:"some"`
		}
	}
	var tj testJSON
	json.Unmarshal([]byte(testJSONObject), &tj)
	return tj.Item1[2].Some.Thing
}

func getSingleValueWithUnmarshalNoSchema() string {
	var rootMap map[string]interface{}
	json.Unmarshal([]byte(testJSONObject), &rootMap)
	itemOne, _ := rootMap["item1"]
	switch val := itemOne.(type) {
	case []interface{}:
		obj := val[2].(map[string]interface{})
		obj2, _ := obj["some"].(map[string]interface{})
		thing, _ := obj2["thing"].(string)
		return thing
	}
	return ""
}

const testJSONObject = `{
	"item1": ["aryitem1", "aryitem2", {"some": {"thing": "coolObj"}}],
	"item2": "simplestringvalue",
	"item3": {
		"item4": {
			"item5": {
				"item6": ["thing1", 2],
				"item7": {"reallyinnerobjkey": {"is": "anobject"}}
			}
		}
	},
	"item4": 1.2345,
	"item5": true
}`
