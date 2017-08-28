package gkgen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateHexValidationCode_String(t *testing.T) {
	hv := NewHexValidator()
	e := ExampleStruct{}
	et := reflect.TypeOf(e)
	field, _ := et.FieldByName("HexString")

	code, err := hv.Generate(et, field, []string{})
	require.NoError(t, err)
	code = strings.Replace(strings.TrimSpace(code), "\t", "", -1)
	require.Equal(t, "if err := gokay.IsHex(&s.HexString); err != nil {\nerrorsHexString = append(errorsHexString, err)\n}", code)
}

func TestGenerateHexValidationCode_StringPtr(t *testing.T) {
	hv := NewHexValidator()
	e := ExampleStruct{}
	et := reflect.TypeOf(e)
	field, _ := et.FieldByName("HexStringPtr")
	code, err := hv.Generate(et, field, []string{})
	require.NoError(t, err)
	code = strings.Replace(strings.TrimSpace(code), "\t", "", -1)
	require.Equal(t, "if err := gokay.IsHex(s.HexStringPtr); err != nil {\nerrorsHexStringPtr = append(errorsHexStringPtr, err)\n}", code)
}

func TestHexValidatorASTGenerate_string(t *testing.T) {
	src := fmt.Sprintf(`
package main

type TestStruct struct {
	HexString string %s
}
`, "`valid:\"Hex\"`")

	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}
	hv := NewHexValidator()
	code, err := hv.ASTGenerate(f.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[0], nil)
	require.NoError(t, err)
	code = strings.Replace(strings.TrimSpace(code), "\t", "", -1)
	require.Equal(t, "if err := gokay.IsHex(&s.HexString); err != nil {\nerrorsHexString = append(errorsHexString, err)\n}", code)
}
