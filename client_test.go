/*
 * Copyright (C) 2017 Dgraph Labs, Inc. and Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

//go:generate efaceconv
//ec::[]interface{}:SLice
//ec::[]interface{}:SLice

package dgo_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lehajam/dgo"
	"github.com/stretchr/testify/assert"
)

type PersonOneFriendOneSchool struct {
	Uid      string     `json:"uid,omitempty"`
	Name     string     `json:"name,omitempty"`
	Age      int        `json:"age,omitempty"`
	Dob      *time.Time `json:"dob,omitempty"`
	Married  bool       `json:"married,omitempty"`
	Raw      []byte     `json:"raw_bytes",omitempty`
	Friend   Person     `json:"friend,omitempty"`
	Location loc        `json:"loc,omitempty"`
	School   School     `json:"school,omitempty"`
}

func TestUnmarshal_Errors(t *testing.T) {
	response := []byte(`{
		"name": "Alice",
		"age": 26,
		"friend": [{
			"name": "Bob",
			"age":  24
		}]
	}`)

	// First check that it does indeed fail with the standard json marshaler
	err := json.Unmarshal(response, &PersonOneFriendOneSchool{})
	assert.NotNil(t, err, "json.Unmarshal managed to unmarshal dgraph response correctly, this means dgo.Unmarshal might not be required anymore")

	// Then check that dgo Unmarshal method solves the problem
	err = dgo.Unmarshal(response, &PersonOneFriendOneSchool{})
	assert.Nil(t, err, "dgo.Unmarshal failed to unmarshal dgraph response")
}

func TestUnmarshal_Values(t *testing.T) {
	type unmarshalResult struct {
		Title    string
		Actual   []byte
		Expected []byte
	}

	cases := []unmarshalResult{
		{
			Title: "Array to struct",
			Actual: []byte(`{
				"name": "Alice",
				"age": 26,
				"friend": [{
					"name": "Bob",
					"age":  24
				}]
			}`),
			Expected: []byte(`{
				"name": "Alice",
				"age": 26,
				"friend": {
					"name": "Bob",
					"age":  24
				}
			}`),
		},
		{
			Title: "Empty array to struct",
			Actual: []byte(`{
				"name": "Alice",
				"age": 26,
				"friend": []
			}`),
			Expected: []byte(`{
				"name": "Alice",
				"age": 26
				}`),
		},
		{
			Title: "No array specified",
			Actual: []byte(`{
				"name": "Alice",
				"age": 26
			}`),
			Expected: []byte(`{
				"name": "Alice",
				"age": 26
			}`),
		},
	}

	for _, c := range cases {
		var actual, expected PersonOneFriendOneSchool
		t.Log(string(c.Title))

		dgo.Unmarshal(c.Actual, &actual)
		json.Unmarshal(c.Expected, &expected)

		t.Log(actual)
		t.Log(expected)
		assert.EqualValues(t, expected, actual, "dgo.Unmarshal did not produce expected object for test \""+c.Title+"\"")
	}
}

// base (dgo)
// BenchmarkUnmarshal-4   	   50000	     23921 ns/op	   11571 B/op	     139 allocs/op
// json
// BenchmarkUnmarshal-4   	  500000	      2815 ns/op	    1272 B/op	      11 allocs/op
// json-iterator
// BenchmarkUnmarshal-4   	  100000	     22116 ns/op	   11181 B/op	     134 allocs/op
func BenchmarkUnmarshal(b *testing.B) {
	response := []byte(`{
		"name": "Alice",
		"age": 26,
		"friend": [{
			"name": "Bob",
			"age":  24
		}]
	}`)

	for n := 0; n < b.N; n++ {
		dgo.Unmarshal(response, &PersonOneFriendOneSchool{})
	}
}
