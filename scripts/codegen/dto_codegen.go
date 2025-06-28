// File: scripts/codegen/dto_codegen.go
// Usage: go run scripts/codegen/dto_codegen.go --model User
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"text/template"
)

var modelName string

func init() {
	// Parse the "--model" flag to target a specific struct
	flag.StringVar(&modelName, "model", "", "Name of the model struct to generate DTO for")
	flag.Parse()

	if modelName == "" {
		log.Fatal("Please specify a model using --model flag")
	}
}

func main() {
	// Recursively walk through the entity directory to find all Go files
	err := filepath.Walk("internal/domain/entity", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".go") || info.IsDir() {
			return nil
		}
		return parseFileAndGenerateDTO(path)
	})

	if err != nil {
		log.Fatalf("Failed to scan entity directory: %v", err)
	}
}

func parseFileAndGenerateDTO(filepath string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filepath, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	// Traverse the AST to find the matching struct by name
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			st, ok2 := ts.Type.(*ast.StructType)
			if !ok || !ok2 || ts.Name.Name != modelName {
				continue
			}
			generateDTOFile(ts.Name.Name, st)
		}
	}
	return nil
}

func generateDTOFile(structName string, structType *ast.StructType) {
	var fields []string

	// Fields to skip in DTOs (server-managed fields)
	skipFields := map[string]bool{
		"ID":        true,
		"CreatedAt": true,
		"UpdatedAt": true,
		"DeletedAt": true,
		// add more as needed
	}

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue // skip anonymous/embedded fields
		}
		name := field.Names[0].Name
		if !unicode.IsUpper(rune(name[0])) {
			continue // skip unexported fields
		}
		if skipFields[name] {
			continue // skip unwanted fields
		}

		// Type conversion
		typeStr := exprToString(field.Type)

		// Special substitution: []byte -> []string
		if typeStr == "[]byte" {
			typeStr = "[]string"
		}

		if typeStr == "time.Time" || typeStr == "ulid.ULID" {
			typeStr = "string"
		}

		// Special substitution: User, []Comment -> omit or string
		if typeStr == "User" || strings.HasPrefix(typeStr, "[]") {
			// you could skip entirely or default to string
			continue
		}

		// JSON tag with omitempty
		jsonTag := strings.ToLower(name[:1]) + name[1:]

		// Default: no validation
		validateTag := ""

		// Check for "gorm:\"not null\""
		if field.Tag != nil {
			tagValue := field.Tag.Value // e.g., "`gorm:\"not null\" json:\"something\"`"
			if strings.Contains(tagValue, "not null") {
				validateTag = `validate:"required"`
			}
		}

		// Build struct tag string
		tagParts := []string{fmt.Sprintf("json:\"%s\"", jsonTag)}
		if validateTag != "" {
			tagParts = append(tagParts, validateTag)
		}
		structTag := "`" + strings.Join(tagParts, " ") + "`"

		fields = append(fields, fmt.Sprintf("\t%s %s %s", name, typeStr, structTag))
	}

	filename := fmt.Sprintf("internal/app/http/dto/request/%s_request.go", strings.ToLower(structName))
	os.MkdirAll(filepath.Dir(filename), os.ModePerm)

	var buf bytes.Buffer
	tpl := template.Must(template.New("dto").Parse(`package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type {{.DTOName}}Request struct {
{{range .Fields}}{{.}}
{{end}}}

// Bind parses and validates the request body and returns an entity
func (req *{{.DTOName}}) Bind(c *fiber.Ctx, v *validator.Validate) error {
	// Parse request body into req
	if err := c.BodyParser(req); err != nil {
		return err
	}

	// Validate request struct
	if err := v.Struct(req); err != nil {
		return err
	}

	return nil
}
	`))
	tpl.Execute(&buf, map[string]any{
		"DTOName": fmt.Sprintf("Create%s", structName),
		"Fields":     fields,
	})

	err := os.WriteFile(filename, buf.Bytes(), 0644)
	if err != nil {
		log.Printf("❌Failed to write DTO file: %v", err)
	} else {
		fmt.Printf("✅ Generated DTO: %s\n", filename)
	}
}

// Converts AST expressions to string types like: int, string, []string, *User
func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	default:
		return "interface{}" // fallback for unknown types
	}
}
