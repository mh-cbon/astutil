package astutil

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestHasEllipse(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("package t\n")
	buf.WriteString(`func t(s ...string){}`)

	fset := token.NewFileSet()
	x, err := parser.ParseFile(fset, "nop.go", &buf, 0)
	if err != nil {
		t.Error(err)
	}
	y := x.Decls[0].(*ast.FuncDecl)
	want := true
	got := MethodHasEllipse(y)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestNotHasEllipse(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("package t\n")
	buf.WriteString(`func t(){}`)

	fset := token.NewFileSet()
	x, err := parser.ParseFile(fset, "nop.go", &buf, 0)
	if err != nil {
		t.Error(err)
	}
	y := x.Decls[0].(*ast.FuncDecl)
	want := false
	got := MethodHasEllipse(y)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestMethodParamNamesInvokation(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("package t\n")
	buf.WriteString(`func t(s ...string){}`)

	fset := token.NewFileSet()
	x, err := parser.ParseFile(fset, "nop.go", &buf, 0)
	if err != nil {
		t.Error(err)
	}
	y := x.Decls[0].(*ast.FuncDecl)
	want := "s..."
	got := MethodParamNamesInvokation(y, true)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestMethodReturnPointer(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("package t\n")
	buf.WriteString(`func t() *y {}`)

	fset := token.NewFileSet()
	x, err := parser.ParseFile(fset, "nop.go", &buf, 0)
	if err != nil {
		t.Error(err)
	}
	y := x.Decls[0].(*ast.FuncDecl)
	want := true
	got := MethodReturnPointer(y)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}

func TestNotMethodReturnPointer(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("package t\n")
	buf.WriteString(`func t() y {}`)

	fset := token.NewFileSet()
	x, err := parser.ParseFile(fset, "nop.go", &buf, 0)
	if err != nil {
		t.Error(err)
	}
	y := x.Decls[0].(*ast.FuncDecl)
	want := false
	got := MethodReturnPointer(y)
	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}
