// Package astutil provides useful methods to work with ast when you intend to make a generator.
package astutil

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"strings"

	"golang.org/x/tools/go/loader"
)

// GetProgram load program of s a pkg path
func GetProgram(s string) *loader.Program {
	args := []string{"--", s}

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
	_, err := conf.FromArgs(args[1:], false)
	if err != nil {
		fmt.Println(err)
	}
	prog, err := conf.Load()
	if err != nil {
		log.Println(err)
	}

	return prog
}

// GetImportPath return the import path of an identifier.
func GetImportPath(p *loader.PackageInfo, name string) string {
	ret := ""
	for _, file := range p.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.ImportSpec:
				if x.Path != nil {
					if x.Name != nil && x.Name.Name == name {
						y := x.Path.Value[1:]
						y = y[:len(y)-1] //rm ""
						if y != name {
							ret = y
						}
					} else if strings.HasSuffix(x.Path.Value, "/"+name+"\"") {
						y := x.Path.Value[1:]
						y = y[:len(y)-1] //rm ""
						ret = y
					}
				}
			}
			return true
		})
	}
	return ret
}

// FindTypes searches given package for every struct types definition
func FindTypes(p *loader.PackageInfo) []string {
	foundTypes := []string{}
	for _, file := range p.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.TypeSpec:
				if _, ok := x.Type.(*ast.StructType); ok {
					foundTypes = append(foundTypes, x.Name.Name)
				}
			}
			return true
		})
	}
	return foundTypes
}

// FindMethods searches given package for every struct methods definition
func FindMethods(p *loader.PackageInfo) map[string][]*ast.FuncDecl {
	foundMethods := map[string][]*ast.FuncDecl{}
	for _, file := range p.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				if x.Recv != nil {
					aboutType := ReceiverType(x)
					if aboutType != "" {
						if _, ok := foundMethods[aboutType]; !ok {
							foundMethods[aboutType] = []*ast.FuncDecl{}
						}
						foundMethods[aboutType] = append(foundMethods[aboutType], x)
					}
				}
			}
			return true
		})
	}
	return foundMethods
}

// FindCtors searches given package for every ctors of given struct list.
func FindCtors(p *loader.PackageInfo, aboutTypes []string) map[string]*ast.FuncDecl {
	foundCtors := map[string]*ast.FuncDecl{}
	for _, file := range p.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				if x.Recv == nil {
					for _, t := range aboutTypes {
						if "New"+t == MethodName(x) {
							foundCtors[t] = x
						}
					}
				}
			}
			return true
		})
	}
	return foundCtors
}

// Print any node x to string
func Print(x interface{}) string {
	var b bytes.Buffer
	fset := token.NewFileSet()
	format.Node(&b, fset, x)
	return b.String()
}

// PrintPkg all files of a package to string.
func PrintPkg(p *loader.PackageInfo) string {
	var b bytes.Buffer
	for _, file := range p.Files {
		b.WriteString(Print(file))
	}
	return b.String()
}

// IsExported name.
func IsExported(m string) bool {
	return ast.IsExported(m)
}

// MethodName returns the name of given func
func MethodName(m *ast.FuncDecl) string {
	return m.Name.Name
}

// MethodReturnPointer returns true if the func returns a pointer.
func MethodReturnPointer(m *ast.FuncDecl) bool {
	if m.Type.Results != nil {
		for _, p := range m.Type.Results.List {
			if _, ok := p.Type.(*ast.StarExpr); ok {
				return true
			}
		}
	}
	return false
}

// MethodReturnError returns true if the last out param is of type error.
func MethodReturnError(m *ast.FuncDecl) bool {
	if m.Type.Results != nil {
		for _, p := range m.Type.Results.List {
			if x, ok := p.Type.(*ast.Ident); ok {
				if x.Name == "error" {
					return true
				}
			}
		}
	}
	return false
}

// MethodReturnTypes returns all types of the out signature.
func MethodReturnTypes(m *ast.FuncDecl) []string {
	var ret []string
	if m.Type.Results != nil {
		for _, p := range m.Type.Results.List {
			ret = append(ret, ToString(p.Type))
		}
	}
	return ret
}

var retVar int

// MethodReturnVars create a list of of unqiue variables for each param of out signature.
func MethodReturnVars(m *ast.FuncDecl) []string {
	var ret []string
	if m.Type.Results != nil {
		for range m.Type.Results.List {
			ret = append(ret, fmt.Sprintf("retVar%v", retVar))
			retVar++
		}
	}
	return ret
}

// MethodParamNames reutrns the list of variable in the in signature.
func MethodParamNames(m *ast.FuncDecl) string {
	var ret []string
	for _, p := range m.Type.Params.List {
		ret = append(ret, p.Names[0].Name)
	}
	return strings.Join(ret, ", ")
}

// MethodParamTypes reutrns the list of variable type in the in signature.
func MethodParamTypes(m *ast.FuncDecl) string {
	var ret []string
	for _, p := range m.Type.Params.List {
		ret = append(ret, ToString(p.Type))
	}
	return strings.Join(ret, ", ")
}

// MethodParamNamesInvokation reutrns the list of variable in the in signature as an invokation.
// If withEllipse is true, the last argument gets uses with the ellipse token.
func MethodParamNamesInvokation(m *ast.FuncDecl, withEllipse bool) string {
	var ret []string
	for _, p := range m.Type.Params.List {
		ret = append(ret, p.Names[0].Name)
	}
	if withEllipse && len(ret) > 0 {
		ret[len(ret)-1] += "..."
	}
	return strings.Join(ret, ", ")
}

// MethodHasEllipse returns true if last param has ellipse.
func MethodHasEllipse(m *ast.FuncDecl) bool {
	l := m.Type.Params.List
	if len(l) > 0 {
		_, ok := l[len(l)-1].Type.(*ast.Ellipsis)
		return ok
	}
	return false
}

// MethodParams returns the in signature.
func MethodParams(m *ast.FuncDecl) string {
	var ret []string
	for _, p := range m.Type.Params.List {
		c := p.Names[0].Name + " "
		c += ToString(p.Type)
		ret = append(ret, c)
	}
	return strings.Join(ret, ", ")
}

// MethodParamsToProps returns the in signature as property list.
func MethodParamsToProps(m *ast.FuncDecl) string {
	var ret []string
	for _, p := range m.Type.Params.List {
		c := p.Names[0].Name + " "
		switch i := p.Type.(type) {
		case *ast.Ellipsis:
			c += "[]" + ToString(i.Elt)
		default:
			c += ToString(p.Type)
		}
		ret = append(ret, c)
	}
	return strings.Join(ret, "\n")
}

// GetSignatureImportIdentifiers extract import identifers from the method signature.
func GetSignatureImportIdentifiers(m *ast.FuncDecl) []string {
	ret := []string{}
	paramsType := MethodParamTypes(m)
	for _, p := range strings.Split(paramsType, ", ") {
		p = strings.TrimSpace(p)
		x := strings.Split(p, ".")
		if len(x) > 1 {
			y := GetUnpointedType(x[0])
			y = GetUnslicedType(y) //todo: can do better.
			if len(y) > 0 {
				ret = append(ret, y)
			}
		}
	}
	returnsType := MethodReturnTypes(m)
	for _, p := range returnsType {
		p = strings.TrimSpace(p)
		x := strings.Split(p, ".")
		if len(x) > 1 {
			y := GetUnpointedType(x[0])
			y = GetUnslicedType(y) //todo: can do better.
			if len(y) > 0 {
				ret = append(ret, y)
			}
		}
	}
	return ret
}

// SetReceiverName sets the receiver variable name of a method.
func SetReceiverName(m *ast.FuncDecl, name string) {
	m.Recv.List[0].Names[0].Name = name
}

// SetReceiverPointer makes sure the receiver type is a pointer.
func SetReceiverPointer(m *ast.FuncDecl, pointer bool) {
	if y, ok := m.Recv.List[0].Type.(*ast.StarExpr); ok {
		if pointer == false {
			m.Recv.List[0].Type = y.X
		}
	} else if u, ok := m.Recv.List[0].Type.(*ast.Ident); ok {
		if pointer {
			m.Recv.List[0].Type = &ast.StarExpr{X: u}
		}
	}
}

// SetReceiverTypeName sets the type of the receiver.
func SetReceiverTypeName(x *ast.FuncDecl, name string) {
	if y, ok := x.Recv.List[0].Type.(*ast.StarExpr); ok {
		y.X.(*ast.Ident).Name = name
	} else if u, ok := x.Recv.List[0].Type.(*ast.Ident); ok {
		u.Name = name
	}
}

// ReceiverName returns the receiver variable name.
func ReceiverName(m *ast.FuncDecl) string {
	return m.Recv.List[0].Names[0].Name
}

// ReceiverType returns the type of the receiver in a method.
func ReceiverType(x *ast.FuncDecl) string {
	ret := ""
	if y, ok := x.Recv.List[0].Type.(*ast.StarExpr); ok {
		ret = y.X.(*ast.Ident).Name
	} else if u, ok := x.Recv.List[0].Type.(*ast.Ident); ok {
		ret = u.Name
	}
	return ret
}

// IsAPointedType returns true for starType.
func IsAPointedType(t string) bool {
	return len(t) > 0 && t[0] == '*'
}

// GetUnpointedType always return the dereferenced type.
// A non pointer types is returned untouched.
func GetUnpointedType(t string) string {
	if IsAPointedType(t) {
		return t[1:]
	}
	return t
}

// GetPointedType always return the type prefixed with a *.
// A pointer types is returned untouched.
func GetPointedType(t string) string {
	if !IsAPointedType(t) {
		t = "*" + t
	}
	return t
}

// IsASlicedType returns true for sliceType.
func IsASlicedType(t string) bool {
	return len(t) > 1 && t[:1] == "[]"
}

// GetUnslicedType always return the unsliced type.
// A non pointer types is returned untouched.
func GetUnslicedType(t string) string {
	if IsASlicedType(t) {
		return t[2:]
	}
	return t
}

//go:generate lister basic_gen.go string:StringSlice
var basicTypes = NewStringSlice().Push(
	"string",
	"int",
	"uint",
	"int8",
	"uint8",
	"int16",
	"uint16",
	"int32",
	"uint32",
	"int64",
	"uint64",
	"float",
	"float64",
	"ufloat",
	"ufloat64",
)

// IsBasic return true when the given type is a basic string...
// The type is always dereferenced.
func IsBasic(t string) bool {
	if IsAPointedType(t) {
		t = t[1:]
	}
	//todo: must have a better way to do this.
	return basicTypes.Index(t) > -1
}

// GetStruct searches given package for a struct named s
func GetStruct(p *loader.PackageInfo, s string) *ast.StructType {
	var ret *ast.StructType
	for _, file := range p.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.TypeSpec:
				if y, ok := x.Type.(*ast.StructType); ok && x.Name.Name == s {
					ret = y
					return false
				}
			}
			return true
		})
	}
	return ret
}

// StructProps returns all props and their types of type s.
func StructProps(s *ast.StructType) []map[string]string {
	ret := []map[string]string{}
	for _, f := range s.Fields.List {
		t := ""
		if f.Tag != nil {
			t = ToString(f.Tag)
		}
		name := ""
		for _, n := range f.Names {
			name += ToString(n) + "."
		}

		ret = append(ret, map[string]string{
			"name": strings.TrimRight(name, "."),
			"type": ToString(f.Type),
			"tag":  t,
		})
	}
	return ret
}

// IsArrayType returns true when the given string is an []Array.
func IsArrayType(s string) bool {
	return len(s) > 0 && s[0] == '['
}

// IsStarType returns true when the given string is a *Star.
func IsStarType(s string) bool {
	return len(s) > 0 && s[0] == '*'
}

// ToString takes an ast.Node and print it to string.
// It does not accept nil values, please check it ahead.
func ToString(n interface{}) string {
	if n == nil {
		return ""
	}
	var buf bytes.Buffer
	fset := token.NewFileSet()
	printer.Fprint(&buf, fset, n)
	return buf.String()
}

// GetComment takes an ast.Node to find its comment.
func GetComment(prog *loader.Program, t token.Pos) string {
	_, path, _ := prog.PathEnclosingInterval(t, t)
	for _, n := range path {
		switch n := n.(type) {
		case *ast.GenDecl:
			return n.Doc.Text()
		case *ast.FuncDecl:
			return n.Doc.Text()
		}
	}
	return ""
}

// GetAnnotations extracts annotations of a comment.
// start is the symbol starting an annotation.
// If start is @,
// @annotationName annotationValue
func GetAnnotations(comment string, start string) map[string]string {
	ret := map[string]string{}
	lines := strings.Split(comment, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > len(start) && line[:len(start)] == start {
			line = line[len(start):]
			x := strings.Split(line, " ")
			if len(x) > 1 {
				name := x[0]
				value := strings.Join(x[1:], " ")
				ret[name] = value
			}
		}
	}
	return ret
}
