package converter

import (
	"os"
	"strconv"
	"testing"
)

func TestConverter(t *testing.T) {
	tests := []struct {
		name     string
		goSource string
		expected string
	}{
		{
			name: "Basic struct with JSON tags",
			goSource: `package test

			type User struct {
				Name     string    ` + "`json:\"name\"`" + `
				Age      int       ` + "`json:\"age\"`" + `
				Email    string    ` + "`json:\"email\"`" + `
				IsActive bool      ` + "`json:\"is_active\"`" + `
				Scores   []float64 ` + "`json:\"scores\"`" + `
				Metadata map[string]interface{} ` + "`json:\"metadata\"`" + `
			}
		`,
			expected: `interface User {
  name: string;
  age: number;
  email: string;
  is_active: boolean;
  scores: number[];
  metadata: Record<string, any>;
}`,
		},
		{
			name: "Struct without JSON tags",
			goSource: `package test

			type Product struct {
				ID    int
				Title string
				Price float64
			}
		`,
			expected: `interface Product {
  ID: number;
  Title: string;
  Price: number;
}`,
		},
		{
			name: "Struct with pointer fields",
			goSource: `package test

			type PointerDemo struct {
				Name *string ` + "`json:\"name\"`" + `
				Age  *int    ` + "`json:\"age\"`" + `
			}
		`,
			expected: `interface PointerDemo {
  name: string;
  age: number;
}`,
		},
		{
			name: "Nested structs",
			goSource: `package test

			type Address struct {
				Street string ` + "`json:\"street\"`" + `
				City   string ` + "`json:\"city\"`" + `
			}
			type Person struct {
				Name    string  ` + "`json:\"name\"`" + `
				Address Address ` + "`json:\"address\"`" + `
			}
		`,
			expected: `interface Address {
  street: string;
  city: string;
}

interface Person {
  name: string;
  address: Address;
}`,
		},
		{
			name: "Array and slice of structs",
			goSource: `package test

			type Item struct {
				Value int ` + "`json:\"value\"`" + `
			}
			type Cart struct {
				Items []Item ` + "`json:\"items\"`" + `
			}
		`,
			expected: `interface Item {
  value: number;
}

interface Cart {
  items: Item[];
}`,
		},
		{
			name: "Map fields",
			goSource: `package test

			type MapDemo struct {
				Labels map[string]string ` + "`json:\"labels\"`" + `
				Counts map[int]int      ` + "`json:\"counts\"`" + `
			}
		`,
			expected: `interface MapDemo {
  labels: Record<string, string>;
  counts: Record<number, number>;
}`,
		},
		{
			name: "Embedded structs",
			goSource: `package test

			type Base struct {
				ID int ` + "`json:\"id\"`" + `
			}
			type Derived struct {
				Base
				Name string ` + "`json:\"name\"`" + `
			}
		`,
			expected: `interface Base {
  id: number;
}

interface Derived {
  Base: Base;
  name: string;
}`,
		},
		{
			name: "Unsupported types",
			goSource: `package test

			type Weird struct {
				Ch chan int ` + "`json:\"ch\"`" + `
				Func func()  ` + "`json:\"func\"`" + `
			}
		`,
			expected: `interface Weird {
  ch: any;
  func: any;
}`,
		},
		{
			name: "Multiple structs in one file",
			goSource: `package test

			type A struct { X int ` + "`json:\"x\"`" + ` }
			type B struct { Y string ` + "`json:\"y\"`" + ` }
		`,
			expected: `interface A {
  x: number;
}

interface B {
  y: string;
}`,
		},
		{
			name: "Omit and dash JSON tags",
			goSource: `package test

			type OmitDemo struct {
				A int ` + "`json:\"a,omitempty\"`" + `
				B int ` + "`json:\"-\"`" + `
			}
		`,
			expected: `interface OmitDemo {
  a: number;
}`,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			file := "testfile_" + strconv.Itoa(i) + ".go"
			err := os.WriteFile(file, []byte(tc.goSource), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
			defer os.Remove(file)

			conv := New()
			result, err := conv.ConvertFile(file)
			if err != nil {
				t.Fatalf("ConvertFile failed: %v", err)
			}

			if result != tc.expected {
				t.Errorf("Test %q failed.\nExpected:\n%s\n\nGot:\n%s", tc.name, tc.expected, result)
			}
		})
	}
}
