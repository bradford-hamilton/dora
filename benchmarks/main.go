package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bradford-hamilton/dora/pkg/dora"
)

// benchmarkGetSingleValueWithDora 6953 ns/op
// MemAllocs: 					   18475087
// MemBytes: 					   726808072

// benchmarkisGetSingleValueByUnmarshal 2845 ns/op
// MemAllocs: 					   6224052
// MemBytes: 					   391700928

// Right now dora seems to be around 2-3x slower then unmarshal

func main() {
	res := testing.Benchmark(benchmarkGetSingleValueWithDora)
	fmt.Printf("%s\n%#[1]v\n", res)
	fmt.Println("benchmarkGetSingleValueWithDora")
	fmt.Printf("MemAllocs: %d\n", res.MemAllocs)
	fmt.Printf("MemBytes: %d\n", res.MemBytes)

	res = testing.Benchmark(benchmarkisGetSingleValueByUnmarshal)
	fmt.Printf("%s\n%#[1]v\n", res)
	fmt.Println("benchmarkisGetSingleValueByUnmarshal")
	fmt.Printf("MemAllocs: %d\n", res.MemAllocs)
	fmt.Printf("MemBytes: %d\n", res.MemBytes)
}

var result string

func benchmarkGetSingleValueWithDora(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := getSingleValueWithDora()
		result = v
	}
}

func benchmarkisGetSingleValueByUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := getSingleValueWithUnmarshal()
		result = v
	}
}

func getSingleValueWithDora() string {
	c, _ := dora.NewFromString(testJSONObject)
	r, _ := c.Get("$.item1[2].some.thing")
	return r
}

func getSingleValueWithUnmarshal() string {
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
