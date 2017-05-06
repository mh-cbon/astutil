package astutil

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"testing"

	"golang.org/x/tools/go/loader"
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
	want := "r string, v *pointer, u []slice, y ...string"
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

func TestMethodReturnTypes(t *testing.T) {
	y := getFuncDecl(`func t(r string, v *pointer, y ...string) (y, error) {}`)
	want := []string{"y", "error"}
	got := MethodReturnTypes(y)
	if want[0] != got[0] {
		t.Errorf("want %v got %v", want[0], got[0])
	}
	if want[1] != got[1] {
		t.Errorf("want %v got %v", want[1], got[1])
	}
}

func TestMethodReturnTypes2(t *testing.T) {
	y := getFuncDecl(`func t(r string, v *pointer, y ...string) (*y, []error) {}`)
	want := []string{"*y", "[]error"}
	got := MethodReturnTypes(y)
	if want[0] != got[0] {
		t.Errorf("want %v got %v", want[0], got[0])
	}
	if want[1] != got[1] {
		t.Errorf("want %v got %v", want[1], got[1])
	}
}

func TestNotMethodReturnTypes(t *testing.T) {
	y := getFuncDecl(`func t(r string, v *pointer, y ...string) {}`)
	got := MethodReturnTypes(y)
	if len(got) > 0 {
		t.Errorf("want %v got %v", 0, got)
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

func TestGetComment(t *testing.T) {
	prog := getProgramFromString(`// the comment.
// with two lines.
type T struct{k []string}`)
	pkg := prog.Package("thepackagename")
	s := GetStruct(pkg, "T")
	got := GetComment(prog, s.Pos())
	want := `the comment.
with two lines.
`
	if want != got {
		t.Errorf("want=%q got=%q", want, got)
	}
}

func TestGetAnnotations(t *testing.T) {
	prog := getProgramFromString(`//T is a type.
// @annotation is an annotation.
type T struct{k []string}`)
	pkg := prog.Package("thepackagename")
	s := GetStruct(pkg, "T")
	comment := GetComment(prog, s.Pos())
	got := GetAnnotations(comment, "@")
	want := "is an annotation."
	key := "annotation"
	if annot, ok := got[key]; ok {
		if annot != want {
			t.Errorf("want=%q got=%q", want, annot)

		}
	} else {
		t.Errorf("key %q not found", key)
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

func getProgramFromString(s string) *loader.Program {
	var buf bytes.Buffer
	buf.WriteString("package thepackagename\n\n")
	buf.WriteString(s)

	var conf loader.Config
	conf.ParserMode = parser.ParseComments
	// those 3 might change later. not sure.
	conf.TypeChecker.IgnoreFuncBodies = true
	conf.TypeChecker.DisableUnusedImportCheck = true
	conf.TypeChecker.Error = func(err error) {
		log.Println(err)
	}
	// this really matters otherise its a pain to generate a partial program.
	conf.AllowErrors = true
	file, err := conf.ParseFile(os.Getenv("gopath")+"/src/y/t.go", &buf)
	if err != nil {
		fmt.Println(err)
	}
	conf.CreateFromFiles("thepackagename", file)
	prog, err := conf.Load()
	if err != nil {
		log.Println(err)
	}
	return prog
}
