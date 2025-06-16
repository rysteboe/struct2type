package converter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// Converter handles the conversion of Go structs to TypeScript types
type Converter struct {
	// Map of Go types to TypeScript types
	typeMap map[string]string
}

// New creates a new Converter instance
func New() *Converter {
	return &Converter{
		typeMap: map[string]string{
			"string":    "string",
			"int":       "number",
			"int8":      "number",
			"int16":     "number",
			"int32":     "number",
			"int64":     "number",
			"uint":      "number",
			"uint8":     "number",
			"uint16":    "number",
			"uint32":    "number",
			"uint64":    "number",
			"float32":   "number",
			"float64":   "number",
			"bool":      "boolean",
			"[]byte":    "string",
			"time.Time": "string",
		},
	}
}

// ConvertFile converts a Go source file to TypeScript types
func (c *Converter) ConvertFile(filePath string) (string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %w", err)
	}

	var types []string
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			tsType := c.convertStruct(typeSpec.Name.Name, structType)
			types = append(types, tsType)
		}
	}

	return strings.Join(types, "\n\n"), nil
}

// convertStruct converts a Go struct to a TypeScript interface
func (c *Converter) convertStruct(name string, structType *ast.StructType) string {
	var fields []string

	for _, field := range structType.Fields.List {
		var fieldName string
		// Embedded field (anonymous)
		if len(field.Names) == 0 {
			// Use the type name as the field name
			switch t := field.Type.(type) {
			case *ast.Ident:
				fieldName = t.Name
			case *ast.SelectorExpr:
				// For imported types (e.g., pkg.Type)
				fieldName = t.Sel.Name
			default:
				fieldName = "embedded"
			}
		} else {
			fieldName = field.Names[0].Name
		}

		fieldType := c.getTypeScriptType(field.Type)

		// Handle json tags if present
		if field.Tag != nil {
			tag := field.Tag.Value
			if strings.Contains(tag, "json:") {
				jsonTag := strings.Split(strings.Split(tag, "json:\"")[1], "\"")[0]
				if jsonTag == "-" {
					continue // skip fields with json:"-"
				}
				if jsonTag != "" {
					// Remove omitempty if present
					fieldName = strings.Split(jsonTag, ",")[0]
				}
			}
		}

		fields = append(fields, fmt.Sprintf("  %s: %s;", fieldName, fieldType))
	}

	return fmt.Sprintf("interface %s {\n%s\n}", name, strings.Join(fields, "\n"))
}

// getTypeScriptType converts a Go type to its TypeScript equivalent
func (c *Converter) getTypeScriptType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		if tsType, ok := c.typeMap[t.Name]; ok {
			return tsType
		}
		return t.Name
	case *ast.ArrayType:
		elemType := c.getTypeScriptType(t.Elt)
		return fmt.Sprintf("%s[]", elemType)
	case *ast.MapType:
		keyType := c.getTypeScriptType(t.Key)
		valueType := c.getTypeScriptType(t.Value)
		return fmt.Sprintf("Record<%s, %s>", keyType, valueType)
	case *ast.StarExpr:
		return c.getTypeScriptType(t.X)
	default:
		return "any"
	}
}
