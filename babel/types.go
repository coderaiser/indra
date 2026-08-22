package babel

import "go/ast"

// IsIdent mirrors types.isIdentifier.
func IsIdent(node ast.Node) bool {
	_, ok := node.(*ast.Ident)
	return ok
}

// IsCallExpr mirrors types.isCallExpression.
func IsCallExpr(node ast.Node) bool {
	_, ok := node.(*ast.CallExpr)
	return ok
}

// IsSelector mirrors types.isMemberExpression.
func IsSelector(node ast.Node) bool {
	_, ok := node.(*ast.SelectorExpr)
	return ok
}

// IsCompositeLit mirrors types.isArrayExpression / isObjectExpression
// (Go has one node kind for both slice and struct literals).
func IsCompositeLit(node ast.Node) bool {
	_, ok := node.(*ast.CompositeLit)
	return ok
}

// IsArrayExpr reports whether node is a composite literal of a slice type.
func IsArrayExpr(node ast.Node) bool {
	lit, ok := node.(*ast.CompositeLit)
	if !ok {
		return false
	}
	if lit.Type == nil {
		return false
	}
	_, isSlice := lit.Type.(*ast.ArrayType)
	return isSlice
}

// IsObjectExpr reports whether node is a composite literal of a struct type.
func IsObjectExpr(node ast.Node) bool {
	lit, ok := node.(*ast.CompositeLit)
	if !ok {
		return false
	}
	if lit.Type == nil {
		return false
	}
	switch lit.Type.(type) {
	case *ast.Ident, *ast.SelectorExpr:
		return true
	}
	return false
}

// IsFuncLit mirrors types.isFunction (FuncLit case).
func IsFuncLit(node ast.Node) bool {
	_, ok := node.(*ast.FuncLit)
	return ok
}

// IsBasicLit mirrors types.isStringLiteral / isNumericLiteral /
// isBooleanLiteral.
func IsBasicLit(node ast.Node) bool {
	_, ok := node.(*ast.BasicLit)
	return ok
}

// IsStatement mirrors types.isStatement.
func IsStatement(node ast.Node) bool {
	_, ok := node.(ast.Stmt)
	return ok
}

// IsFile mirrors types.isFile.
func IsFile(node ast.Node) bool {
	_, ok := node.(*ast.File)
	return ok
}

// IsBoolLit mirrors putout types.isBooleanLiteral(node, {value}).
// Reports whether node is an identifier with name "true" or "false"
// and its value matches val.
func IsBoolLit(node ast.Node, val bool) bool {
	ident, ok := node.(*ast.Ident)
	if !ok {
		return false
	}
	if ident.Name == "true" {
		return val
	}
	if ident.Name == "false" {
		return !val
	}
	return false
}
