# struct2type

A command-line tool that converts Go structs to TypeScript types. It supports basic Go types, arrays, maps, and structs with JSON tags.

## Features

- Converts Go structs to TypeScript interfaces
- Supports JSON tags for field names
- Handles basic Go types (string, int, float, bool, etc.)
- Supports arrays and maps
- Preserves struct field names and types
- Command-line interface with output file option

## Installation

```bash
go install github.com/rysteboe/struct2type/cmd/struct2type@latest
```

This will download the latest version of your module and install the binary to your `GOPATH/bin` directory. Make sure your `GOPATH/bin` is in your system's `PATH` to run the tool from anywhere.

## Usage

Basic usage:
```bash
struct2type input.go
```

Save output to a file:
```bash
struct2type input.go -o output.ts
```

## Example

Input (Go):
```go
type User struct {
    Name     string    `json:"name"`
    Age      int       `json:"age"`
    Email    string    `json:"email"`
    IsActive bool      `json:"is_active"`
    Scores   []float64 `json:"scores"`
    Metadata map[string]interface{} `json:"metadata"`
}
```

Output (TypeScript):
```typescript
interface User {
  name: string;
  age: number;
  email: string;
  is_active: boolean;
  scores: number[];
  metadata: Record<string, any>;
}
```

## Development

1. Clone the repository:
```bash
git clone https://github.com/rysteboe/struct2type.git
cd struct2type
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests:
```bash
go test ./...
```

4. Build:
```bash
go build -o struct2type cmd/struct2type/main.go
```

## License

MIT License - see LICENSE file for details 