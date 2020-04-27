package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bradford-hamilton/dora/pkg/dora"
)

// Original Benchmark values for first three benches and some summaries
// benchmarkGetSingleValueWithDora 7175 ns/op
// MemAllocs: 17440028
// MemBytes: 678406512

// benchmarkisGetSingleValueWithUnmarshalAndNoSchema 5234 ns/op
// MemAllocs: 12724137
// MemBytes: 790109456

// benchmarkisGetSingleValueWithUnmarshalAndSchema 2975 ns/op
// MemAllocs: 5744389
// MemBytes: 361514088

// Dora's stats against unmarhsalling into an unknown interface{}:
// 		- Around 1.5x slower
// 		- Uses slightly less MemBytes
// 		- Uses slightly more MemAllocs

// Dora's stats against unmarhsalling into a known shape (testJSON struct):
// 		- Around 2-3x slower
// 		- Around 2-3x more MemBytes
// 		- Around 2-3x more  MemAllocs

func main() {
	res := testing.Benchmark(benchmarkGetSingleValueWithDora)
	fmt.Println("benchmarkGetSingleValueWithDora")
	fmt.Printf("MemAllocs: %d\n", res.MemAllocs)
	fmt.Printf("MemBytes: %d\n", res.MemBytes)

	fmt.Print("\n")

	res = testing.Benchmark(benchmarkisGetSingleValueWithUnmarshalAndNoSchema)
	fmt.Println("benchmarkisGetSingleValueWithUnmarshalAndNoSchema")
	fmt.Printf("MemAllocs: %d\n", res.MemAllocs)
	fmt.Printf("MemBytes: %d\n", res.MemBytes)

	fmt.Print("\n")

	res = testing.Benchmark(benchmarkisGetSingleValueWithUnmarshalAndSchema)
	fmt.Println("benchmarkisGetSingleValueWithUnmarshalAndSchema")
	fmt.Printf("MemAllocs: %d\n", res.MemAllocs)
	fmt.Printf("MemBytes: %d\n", res.MemBytes)
}

var sink string

func benchmarkGetSingleValueWithDora(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := getSingleValueWithDora()
		sink = v
	}
}

func benchmarkisGetSingleValueWithUnmarshalAndSchema(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := getSingleValueWithUnmarshalAndSchema()
		sink = v
	}
}

func benchmarkisGetSingleValueWithUnmarshalAndNoSchema(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := getSingleValueWithUnmarshalNoSchema()
		sink = v
	}
}

func getSingleValueWithDora() string {
	c, _ := dora.NewFromString(testJSONObject)
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
	}
}`
