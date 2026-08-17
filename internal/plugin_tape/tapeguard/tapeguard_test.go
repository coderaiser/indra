package tapeguard

import (
	"go/ast"
	"testing"

	"coderaiser/indra/types"
)

func fixture(imports ...string) *ast.File {
	file := &ast.File{}
	for _, p := range imports {
		file.Imports = append(file.Imports, &ast.ImportSpec{
			Path: &ast.BasicLit{Value: p},
		})
	}
	return file
}

func TestImported(t *testing.T) {
	cases := []struct {
		name string
		path types.Path
		want bool
	}{
		{
			name: "tape imported",
			path: types.Path{Stack: []ast.Node{&ast.Ident{}, fixture(`"github.com/coderaiser/go-tape"`)}},
			want: true,
		},
		{
			name: "tape not imported",
			path: types.Path{Stack: []ast.Node{&ast.Ident{}, fixture(`"fmt"`)}},
			want: false,
		},
		{
			name: "no file in stack",
			path: types.Path{Stack: []ast.Node{&ast.Ident{}}},
			want: false,
		},
		{
			name: "empty stack",
			path: types.Path{},
			want: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Imported(nil, c.path)
			if got != c.want {
				t.Fatalf("Imported = %v, want %v", got, c.want)
			}
		})
	}
}
