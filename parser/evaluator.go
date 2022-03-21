package parser

import (
	"container/list"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type OpertorStack struct {
	dll *list.List
}

type Result struct {
	shouldAppend            bool
	isAllIntegrtedIdsExists bool
}

type Operand struct {
	oper  *ast.Ident
	depth int
}

func NewStack() *OpertorStack {
	return &OpertorStack{dll: list.New()}
}

func (s *OpertorStack) Push(x interface{}) {
	s.dll.PushBack(x)
}

func (s *OpertorStack) Pop() interface{} {
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
	// v := &visitor{operStack: NewStack(), tree_height: 0, childStack: NewStack(), resultStack: NewStack(), result: b}
	v := &visitor{operStack: NewStack(), tree_height: 0, childStackMap: make(map[int]*OpertorStack), resultStackMap: make(map[int]*OpertorStack), result: b}
	ast.Walk(v, e)

	if v.operStack.dll.Len() != 0 {
		operands := []*ast.Ident{}
		for k := range v.resultStackMap {
			for op := v.resultStackMap[k].dll.Front(); op != nil; op = op.Next() {
				operands = append(operands, op.Value.(*Operand).oper)
			}
		}

		if len(v.childStackMap) != 0 {
			for k := range v.childStackMap {
				for op := v.childStackMap[k].dll.Front(); op != nil; op = op.Next() {
					operands = append(operands, op.Value.(*ast.Ident))
				}
			}
		}

		evaluaeExpression(v, operands[0], operands[1])
	}
	// fmt.Printf("all id :%v \n", v.containsAllIds.isAllIntegrtedIdsExists)
	// fmt.Printf("res id :%v \n", v.containsAllIds.shouldAppend)
	return (v.result.shouldAppend || !v.result.isAllIntegrtedIdsExists)
}

type visitor struct {
	operStack   *OpertorStack
	tree_height int
	// childStack     *OpertorStack
	// resultStack    *OpertorStack
	result         *Result
	childStackMap  map[int]*OpertorStack
	resultStackMap map[int]*OpertorStack
}

func (V visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch x := n.(type) {
	case *ast.Ident:
		{
			fmt.Printf("%s%s : %d\n", strings.Repeat("\t", V.tree_height), x.Name, V.tree_height)
			chStack := getStack(V.childStackMap, V.tree_height)
			if chStack.dll.Len() == 1 {

				leftChild := chStack.Pop().(*ast.Ident)
				rightChild := x
				delete(V.childStackMap, V.tree_height)
				evaluaeExpression(&V, leftChild, rightChild)
				resStack := getStack(V.resultStackMap, V.tree_height)
				if resStack.dll.Len() == 2 {
					leftChild = resStack.Pop().(*Operand).oper
					rightChild = resStack.Pop().(*Operand).oper
					delete(V.resultStackMap, V.tree_height)
					evaluaeExpression(&V, leftChild, rightChild)
				}
			} else {
				chStack.Push(x)
			}
			if !V.result.isAllIntegrtedIdsExists {
				V.result.isAllIntegrtedIdsExists = isAllIntegrtedIds(x.Name)
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

func getStack(stackMap map[int]*OpertorStack, depth int) *OpertorStack {
	if stackMap[depth] == nil {
		stackMap[depth] = NewStack()
	}
	return stackMap[depth]
}

func evaluaeExpression(V *visitor, leftChild *ast.Ident, rightChild *ast.Ident) {
	operator := V.operStack.Pop().(*ast.BinaryExpr)
	resStack := getStack(V.resultStackMap, V.tree_height-1)
	fmt.Printf("evaluating %v %v %v \n", leftChild.Name, operator.Op, rightChild.Name)
	if isAllIntegrtedIds(rightChild.Name) && isAllIntegrtedIds(leftChild.Name) {
		V.result.shouldAppend = false
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
			V.result.shouldAppend = true
			resStack.Push(&Operand{oper: nonAllItgChild, depth: V.tree_height})
		} else {
			resStack.Push(&Operand{oper: allItgChild, depth: V.tree_height})
			V.result.shouldAppend = false
		}

	} else {
		resStack.Push(&Operand{oper: rightChild, depth: V.tree_height})
	}
}

func isAllIntegrtedIds(name string) bool {
	return "b" == name
}
