package parser

import (
	"container/list"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"sync"
)

type OpertorStack struct {
	dll   *list.List
	mutex sync.Mutex
}

type Result struct {
	shouldAppend            bool
	isAllIntegrtedIdsExists bool
}

func NewStack() *OpertorStack {
	return &OpertorStack{mutex: sync.Mutex{}, dll: list.New()}
}

func (s *OpertorStack) Push(x interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.dll.PushBack(x)
}

func (s *OpertorStack) Pop() interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.dll.Len() == 0 {
		return nil
	}
	tail := s.dll.Back()
	val := tail.Value
	s.dll.Remove(tail)
	return val
}

func Evaluate(expr string) bool {
	fmt.Printf("Given Exprssion : %v\n", expr)

	e, err := parser.ParseExpr(expr)

	if err != nil {
		fmt.Errorf("Failed to parse the expression : %v", e)
	}
	b := &Result{false, false}
	v := &visitor{operStack: NewStack(), tree_height: 0, childStack: NewStack(), resultStack: NewStack(), containsAllIds: b}
	ast.Walk(v, e)
	// fmt.Printf("all id :%v \n", v.containsAllIds.isAllIntegrtedIdsExists)
	// fmt.Printf("res id :%v \n", v.containsAllIds.shouldAppend)
	return (v.containsAllIds.shouldAppend || !v.containsAllIds.isAllIntegrtedIdsExists)
}

// func printOperators(opStack *OpertorStack) {
// 	for opStack.dll.Len() > 0 {
// 		x := opStack.Pop().(*ast.BinaryExpr)
// 		// fmt.Printf("op : %v  left : %v   right : %v \n", x.Op, x.X, x.Y)
// 		fmt.Printf("op %v", x.Op)

// 	}
// }

// func printChildren(opStack *OpertorStack) {
// 	for opStack.dll.Len() > 0 {
// 		fmt.Printf("c :%v \n", opStack.Pop().(*ast.Ident).Name)
// 	}
// }

type visitor struct {
	operStack      *OpertorStack
	tree_height    int
	childStack     *OpertorStack
	resultStack    *OpertorStack
	containsAllIds *Result
}

func (V visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch x := n.(type) {
	case *ast.Ident:
		{
			fmt.Printf("%s%s : %d\n", strings.Repeat("\t", V.tree_height), x.Name, V.tree_height)
			if V.childStack.dll.Len() == 1 {

				leftChild := V.childStack.Pop().(*ast.Ident)
				rightChild := x

				evaluaeExpression(&V, leftChild, rightChild)

				if V.resultStack.dll.Len() == 2 {
					leftChild = V.resultStack.Pop().(*ast.Ident)
					rightChild = V.resultStack.Pop().(*ast.Ident)
					evaluaeExpression(&V, leftChild, rightChild)
				}
			} else {
				V.childStack.Push(x)
			}
			if !V.containsAllIds.isAllIntegrtedIdsExists {
				V.containsAllIds.isAllIntegrtedIdsExists = isAllIntegrtedIds(x.Name)
			}
		}
	case *ast.BinaryExpr:
		{
			fmt.Printf("%s%s : %d\n", strings.Repeat("\t", int(V.tree_height)), x.Op, V.tree_height)
			V.operStack.Push(x)
		}
	}
	V.tree_height = V.tree_height + 1
	return &V
}

func evaluaeExpression(V *visitor, leftChild *ast.Ident, rightChild *ast.Ident) {
	operator := V.operStack.Pop().(*ast.BinaryExpr)
	// fmt.Printf("evalating %v %v %v \n", leftChild.Name, operator.Op, rightChild.Name)
	if isAllIntegrtedIds(rightChild.Name) && isAllIntegrtedIds(leftChild.Name) {
		V.containsAllIds.shouldAppend = false
	}
	if isAllIntegrtedIds(rightChild.Name) || isAllIntegrtedIds(leftChild.Name) {
		var allItgChild *ast.Ident
		var nonAllItgChild *ast.Ident

		if isAllIntegrtedIds(rightChild.Name) {
			allItgChild = rightChild
			nonAllItgChild = leftChild
		} else {
			nonAllItgChild = rightChild
			allItgChild = leftChild
		}
		if operator.Op != token.AND {
			V.containsAllIds.shouldAppend = true
			V.resultStack.Push(nonAllItgChild)
		} else {
			V.resultStack.Push(allItgChild)
			V.containsAllIds.shouldAppend = false
		}

	} else {
		V.resultStack.Push(rightChild)
	}
}

func isAllIntegrtedIds(name string) bool {
	return "b" == name
}
