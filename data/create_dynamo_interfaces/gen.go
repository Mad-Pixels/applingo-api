package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

type DynamoSchema struct {
	TableName        string           `json:"table_name"`
	HashKey          string           `json:"hash_key"`
	RangeKey         string           `json:"range_key"`
	Attributes       []Attribute      `json:"attributes"`
	CommonAttributes []Attribute      `json:"common_attributes"`
	SecondaryIndexes []SecondaryIndex `json:"secondary_indexes"`
}

type Attribute struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SecondaryIndex struct {
	Name             string   `json:"name"`
	HashKey          string   `json:"hash_key"`
	RangeKey         string   `json:"range_key"`
	ProjectionType   string   `json:"projection_type"`
	NonKeyAttributes []string `json:"non_key_attributes,omitempty"`
}

const codeTemplate = `
// Code generated by dynamo_dictionary_table.go. DO NOT EDIT.

package {{.PackageName}}

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	TableName = "{{.TableName}}"
	{{range .SecondaryIndexes}}
	Index{{.Name}} = "{{.Name}}"
	{{- end}}
)

var (
	AttributeNames = []string{
		{{- range .AllAttributes}}
		"{{.Name}}",
		{{- end}}
	}

	IndexProjections = map[string][]string{
		{{- range .SecondaryIndexes}}
		Index{{.Name}}: {
			{{- if eq .ProjectionType "ALL"}}
			{{- range $.AllAttributes}}
			"{{.Name}}",
			{{- end}}
			{{- else}}
			"{{.HashKey}}", {{if .RangeKey}}"{{.RangeKey}}",{{end}}
			{{- range .NonKeyAttributes}}
			"{{.}}",
			{{- end}}
			{{- end}}
		},
		{{- end}}
	}
)

type DynamoSchema struct {
	TableName        string
	HashKey          string
	RangeKey         string
	Attributes       []Attribute
	CommonAttributes []Attribute
	SecondaryIndexes []SecondaryIndex
}

type Attribute struct {
	Name string
	Type string
}

type SecondaryIndex struct {
	Name             string
	HashKey          string
	RangeKey         string
	ProjectionType   string
	NonKeyAttributes []string
}

// SchemaItem implement struct for "{{.TableName}}"
type SchemaItem struct {
	{{range .AllAttributes}}
	{{SafeName .Name | ToCamelCase}} {{TypeGo .Type}} ` + "`json:\"{{.Name}}\"`" + `
	{{end}}
}

var TableSchema = DynamoSchema{
	TableName: "{{.TableName}}",
	HashKey:   "{{.HashKey}}",
	RangeKey:  "{{.RangeKey}}",
	Attributes: []Attribute{
		{{- range .Attributes}}
		{Name: "{{.Name}}", Type: "{{.Type}}"},
		{{- end}}
	},
	CommonAttributes: []Attribute{
		{{- range .CommonAttributes}}
		{Name: "{{.Name}}", Type: "{{.Type}}"},
		{{- end}}
	},
	SecondaryIndexes: []SecondaryIndex{
		{{- range .SecondaryIndexes}}
		{
			Name:           "{{.Name}}",
			HashKey:        "{{.HashKey}}",
			RangeKey:       "{{.RangeKey}}",
			ProjectionType: "{{.ProjectionType}}",
			{{- if .NonKeyAttributes}}
			NonKeyAttributes: []string{
				{{- range .NonKeyAttributes}}
				"{{.}}",
				{{- end}}
			},
			{{- end}}
		},
		{{- end}}
	},
}

type QueryBuilder struct {
	IndexName       string
	KeyCondition    expression.KeyConditionBuilder
	FilterCondition expression.ConditionBuilder
	UsedKeys        map[string]bool
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		UsedKeys: make(map[string]bool),
	}
}

{{range .AllAttributes}}
func (qb *QueryBuilder) With{{SafeName .Name | ToCamelCase}}({{SafeName .Name | ToLowerCamelCase}} {{TypeGo .Type}}) *QueryBuilder {
	attrName := "{{.Name}}"
	{{- $attrName := .Name}}
	{{- $goType := TypeGo .Type}}
	{{- $safeName := SafeName .Name}}
	{{- range $.SecondaryIndexes}}
	{{- if eq .HashKey $attrName}}
	if qb.IndexName == "" {
		qb.IndexName = Index{{.Name}}
		qb.KeyCondition = expression.Key(attrName).Equal(expression.Value({{$safeName | ToLowerCamelCase}}))
		qb.UsedKeys[attrName] = true
		return qb
	}
	{{- end}}
	{{- end}}
	if !qb.UsedKeys[attrName] {
		cond := expression.Name(attrName).Equal(expression.Value({{$safeName | ToLowerCamelCase}}))
		if qb.FilterCondition.IsSet() {
			qb.FilterCondition = qb.FilterCondition.And(cond)
		} else {
			qb.FilterCondition = cond
		}
		qb.UsedKeys[attrName] = true
	}
	return qb
}
{{end}}

func (qb *QueryBuilder) Build() (string, expression.KeyConditionBuilder, expression.ConditionBuilder) {
	return qb.IndexName, qb.KeyCondition, qb.FilterCondition
}

// PutItem create an AttributeValues map for PutItem in DynamoDB
func PutItem(item SchemaItem) (map[string]types.AttributeValue, error) {
	attributeValues := make(map[string]types.AttributeValue)
	{{range .AllAttributes}}
	{{- $attrName := .Name}}
	{{- $safeName := SafeName .Name}}
	{{- $goType := TypeGo .Type}}
	{{- if eq .Type "N"}}
		attributeValues["{{.Name}}"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", item.{{$safeName | ToCamelCase}})}
	{{- else if eq .Type "B"}}
		attributeValues["{{.Name}}"] = &types.AttributeValueMemberBOOL{Value: item.{{$safeName | ToCamelCase}}}
	{{- else if eq .Type "S"}}
		if item.{{$safeName | ToCamelCase}} != "" {
			attributeValues["{{.Name}}"] = &types.AttributeValueMemberS{Value: item.{{$safeName | ToCamelCase}}}
		}
	{{- else}}
		
	{{- end}}
	{{end}}
	return attributeValues, nil
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
`

func main() {
	rootDir, err := filepath.Abs(".")
	if err != nil {
		fmt.Printf("Cannot find root directory: %v\n", err)
		return
	}
	tmplDir := filepath.Join(rootDir, ".tmpl")

	files, err := os.ReadDir(tmplDir)
	if err != nil {
		fmt.Printf("Read template directory failed %s: %v\n", tmplDir, err)
		return
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			processSchemaFile(filepath.Join(tmplDir, file.Name()), rootDir)
		}
	}
}

func processSchemaFile(jsonPath, rootDir string) {
	jsonFile, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Printf("Failed read json %s: %v\n", jsonPath, err)
		return
	}

	var schema DynamoSchema
	err = json.Unmarshal(jsonFile, &schema)
	if err != nil {
		fmt.Printf("Parse json failed %s: %v\n", jsonPath, err)
		return
	}

	packageName := "gen_" + strings.ReplaceAll(schema.TableName, "-", "_")
	packageDir := filepath.Join(rootDir, packageName)

	if err := os.MkdirAll(packageDir, os.ModePerm); err != nil {
		fmt.Printf("Create dictionary failed %s: %v\n", packageDir, err)
		return
	}
	outputPath := filepath.Join(packageDir, packageName+".go")

	funcMap := template.FuncMap{
		"ToCamelCase":      toCamelCase,
		"ToLowerCamelCase": toLowerCamelCase,
		"SafeName":         safeName,
		"TypeGo":           typeGo,
	}
	allAttributes := append(schema.Attributes, schema.CommonAttributes...)

	schemaMap := map[string]interface{}{
		"PackageName":      packageName,
		"TableName":        schema.TableName,
		"HashKey":          schema.HashKey,
		"RangeKey":         schema.RangeKey,
		"Attributes":       schema.Attributes,
		"CommonAttributes": schema.CommonAttributes,
		"AllAttributes":    allAttributes,
		"SecondaryIndexes": schema.SecondaryIndexes,
	}

	tmpl, err := template.New("schema").Funcs(funcMap).Parse(codeTemplate)
	if err != nil {
		fmt.Printf("Parse template failed: %v\n", err)
		return
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Creation file error %s: %v\n", outputPath, err)
		return
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, schemaMap)
	if err != nil {
		fmt.Printf("Template execution failed %s: %v\n", outputPath, err)
		return
	}
	fmt.Printf("%s sucessful generated!\n", schema.TableName)
}

func toCamelCase(s string) string {
	var result string
	capitalizeNext := true
	for _, r := range s {
		if r == '_' || r == '-' {
			capitalizeNext = true
		} else if capitalizeNext {
			result += string(unicode.ToUpper(r))
			capitalizeNext = false
		} else {
			result += string(r)
		}
	}
	return result
}

func toLowerCamelCase(s string) string {
	if s == "" {
		return ""
	}
	s = toCamelCase(s)
	return strings.ToLower(s[:1]) + s[1:]
}

var reservedWords = map[string]bool{
	"break":       true,
	"default":     true,
	"func":        true,
	"interface":   true,
	"select":      true,
	"case":        true,
	"defer":       true,
	"go":          true,
	"map":         true,
	"struct":      true,
	"chan":        true,
	"else":        true,
	"goto":        true,
	"package":     true,
	"switch":      true,
	"const":       true,
	"fallthrough": true,
	"if":          true,
	"range":       true,
	"type":        true,
	"continue":    true,
	"for":         true,
	"import":      true,
	"return":      true,
	"var":         true,
}

func safeName(s string) string {
	if reservedWords[s] {
		return s + "_"
	}
	return s
}

func typeGo(dynamoType string) string {
	switch dynamoType {
	case "S":
		return "string"
	case "N":
		return "int"
	case "B":
		return "bool"
	default:
		return "interface{}"
	}
}
