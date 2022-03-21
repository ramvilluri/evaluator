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

	node, err := parser.ParseExpr(expr)

	if err != nil {
		fmt.Errorf("Failed to parse the expression : %v", node)
	}
	b := &Result{false, false}
	v := &visitor{
		operStack:          NewStack(),
		tree_height:        0,
		depthToChildrenMap: make(map[int]*ast.Ident),
		resultStackMap:     make(map[int]*OpertorStack),
		result:             b,
	}
	ast.Walk(v, node)

	if v.operStack.dll.Len() != 0 {
		operands := []*ast.Ident{}
		for k := range v.resultStackMap {
			for op := v.resultStackMap[k].dll.Front(); op != nil; op = op.Next() {
				operands = append(operands, op.Value.(*Operand).oper)
			}
		}

		if len(v.depthToChildrenMap) != 0 {
			for k := range v.depthToChildrenMap {
				// for op := v.childStackMap[k].dll.Front(); op != nil; op = op.Next() {
				if v.depthToChildrenMap[k] != nil {
					operands = append(operands, v.depthToChildrenMap[k])
				}
			}
		}

		evaluaeExpression(v, operands[0], operands[1])
	}
	return (v.result.shouldAppend || !v.result.isAllIntegrtedIdsExists)
}

type visitor struct {
	operStack          *OpertorStack
	tree_height        int
	result             *Result
	depthToChildrenMap map[int]*ast.Ident
	resultStackMap     map[int]*OpertorStack
}

func (V visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch x := n.(type) {
	case *ast.Ident:
		{
			fmt.Printf("%s%s : %d\n", strings.Repeat("\t", V.tree_height), x.Name, V.tree_height)
			leftChild := V.depthToChildrenMap[V.tree_height]
			if leftChild != nil {
				rightChild := x
				delete(V.depthToChildrenMap, V.tree_height)
				evaluaeExpression(&V, leftChild, rightChild)
				resStack := getStack(V.resultStackMap, V.tree_height)
				if resStack.dll.Len() == 2 {
					leftChild = resStack.Pop().(*Operand).oper
					rightChild = resStack.Pop().(*Operand).oper
					delete(V.resultStackMap, V.tree_height)
					evaluaeExpression(&V, leftChild, rightChild)
				}
			} else {
				V.depthToChildrenMap[V.tree_height] = x
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
		resStack.Push(&Operand{oper: rightChild, depth: V.tree_height})
	} else if isAllIntegrtedIds(rightChild.Name) || isAllIntegrtedIds(leftChild.Name) {

		allItgChild, nonAllItgChild := findAllIntegratedId(leftChild, rightChild)

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

func findAllIntegratedId(leftChild *ast.Ident, rightChild *ast.Ident) (*ast.Ident, *ast.Ident) {
	if isAllIntegrtedIds(leftChild.Name) {
		return leftChild, rightChild
	}

	return rightChild, leftChild

}

func isAllIntegrtedIds(name string) bool {
	return "b" == name
}
