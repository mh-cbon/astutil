package astutil

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestHasEllipse(t *testing.T) {
	y := getFuncDecl(`func t(s ...string){}`)
	want := true
	got := MethodHasEllipse(y)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestNotHasEllipse(t *testing.T) {
	y := getFuncDecl(`func t(){}`)
	want := false
	got := MethodHasEllipse(y)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestMethodParamNamesInvokation(t *testing.T) {
	y := getFuncDecl(`func t(s ...string){}`)
	want := "s..."
	got := MethodParamNamesInvokation(y, true)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestMethodReturnPointer(t *testing.T) {
	y := getFuncDecl(`func t() *y {}`)
	want := true
	got := MethodReturnPointer(y)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestNotMethodReturnPointer(t *testing.T) {
	y := getFuncDecl(`func t() y {}`)
	want := false
	got := MethodReturnPointer(y)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestMethodParams(t *testing.T) {
	y := getFuncDecl(`func t() y {}`)
	want := ""
	got := MethodParams(y)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestMethodParams2(t *testing.T) {
	y := getFuncDecl(`func t(r string, v *pointer, u []slice, y ...string) y {}`)
	want := "r string, v pointer, u []slice, y ...string"
	got := MethodParams(y)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestMethodParams3(t *testing.T) {
	y := getFuncDecl(`func t(r func()) y {}`)
	want := "r func()"
	got := MethodParams(y)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestMethodReturnError(t *testing.T) {
	y := getFuncDecl(`func t(r string, v *pointer, y ...string) (y, error) {}`)
	want := true
	got := MethodReturnError(y)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestNotMethodReturnError(t *testing.T) {
	y := getFuncDecl(`func t(r string, v *pointer, y ...string) y {}`)
	want := false
	got := MethodReturnError(y)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestNotMethodReturnError2(t *testing.T) {
	y := getFuncDecl(`func t(r string, v *pointer, y ...string) {}`)
	want := false
	got := MethodReturnError(y)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestStructProps(t *testing.T) {
	y := getStructDecl(`type t struct{k string}`)
	props := StructProps(y.Type.(*ast.StructType))
	iwant := 1
	igot := len(props)
	if iwant != igot {
		t.Errorf("want %v got %v", iwant, igot)
	}

	prop := props[0]
	want := "k"
	got := prop["name"]
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
	want = "string"
	got = prop["type"]
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
	want = ""
	got = prop["tag"]
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestStructProps2(t *testing.T) {
	y := getStructDecl(`type t struct{k []string}`)
	props := StructProps(y.Type.(*ast.StructType))
	iwant := 1
	igot := len(props)
	if iwant != igot {
		t.Errorf("want %v got %v", iwant, igot)
	}

	prop := props[0]
	want := "k"
	got := prop["name"]
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
	want = "[]string"
	got = prop["type"]
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
	want = ""
	got = prop["tag"]
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func getFuncDecl(s string) *ast.FuncDecl {
	var buf bytes.Buffer
	buf.WriteString("package t\n")
	buf.WriteString(s)

	fset := token.NewFileSet()
	x, err := parser.ParseFile(fset, "nop.go", &buf, 0)
	if err != nil {
		panic(err)
	}
	return x.Decls[0].(*ast.FuncDecl)
}

func getStructDecl(s string) *ast.TypeSpec {
	var buf bytes.Buffer
	buf.WriteString("package t\n")
	buf.WriteString(s)

	fset := token.NewFileSet()
	x, err := parser.ParseFile(fset, "nop.go", &buf, 0)
	if err != nil {
		panic(err)
	}
	return x.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec)
}
