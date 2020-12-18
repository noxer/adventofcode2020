package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var errClosingBracket = errors.New("unexpected closing bracket")

func main() {
	exprs, err := readExprs("input.txt")
	if err != nil {
		fmt.Printf("Fehler beim Einlesen der Aufgaben: %s\n", err)
		os.Exit(1)
	}

	sum := 0

	for _, expr := range exprs {
		sum += expr.AttemptSwap().Eval()
	}

	fmt.Printf("Ergebnis: %d\n", sum)
}

// Expr ...
type Expr struct {
	Op    byte
	Left  *Expr
	Right *Expr
	Const int
}

// IsLeaf ...
func (e *Expr) IsLeaf() bool {
	return e.Left == nil && e.Right == nil
}

// Eval ...
func (e *Expr) Eval() int {
	if e.IsLeaf() {
		return e.Const
	}

	switch e.Op {
	case '+':
		return e.Left.Eval() + e.Right.Eval()

	case '*':
		return e.Left.Eval() * e.Right.Eval()

	case '(':
		return e.Left.Eval()
	}

	panic("unknown operation " + string(e.Op))
}

// AttemptSwap ...
func (e *Expr) AttemptSwap() *Expr {
	if e.Left != nil {
		e.Left = e.Left.AttemptSwap()
	}
	if e.Right != nil {
		e.Right = e.Right.AttemptSwap()
	}

	if e.Op != '+' {
		return e
	}

	if e.Left != nil && e.Left.Op == '*' {
		newParent := e.Left
		e.Left = newParent.Right
		newParent.Right = e
		return newParent
	}

	if e.Right != nil && e.Right.Op == '*' {
		newParent := e.Right
		e.Right = newParent.Left
		newParent.Left = e
		return newParent
	}

	return e
}

//       (+)
//       / \
//     (+)  3
//     1 2

//       (*)
//       / \
//      1  (+)
//         2 3

func (e *Expr) String() string {
	if e.IsLeaf() {
		return strconv.Itoa(e.Const)
	}

	if e.Op == '(' {
		return fmt.Sprintf("(%s)", e.Left)
	}

	return fmt.Sprintf("{%s %c %s}", e.Left, e.Op, e.Right)
}

func parseString(line string) (*Expr, error) {
	buf := bufio.NewReader(strings.NewReader(line))
	return parse(buf)
}

func parse(buf *bufio.Reader) (*Expr, error) {
	left, err := consumeWhitespace(buf)
	if err != nil {
		return nil, err
	}

	if left >= '0' && left <= '9' {
		leftExpr := &Expr{Const: int(left - '0')}
		return parseFull(leftExpr, buf)
	}

	if left == '(' {
		leftExpr, err := parse(buf)
		if err != errClosingBracket {
			return nil, errors.New("missing closing bracket")
		}

		return parseFull(&Expr{Op: '(', Left: leftExpr}, buf)
	}

	return nil, fmt.Errorf("unexpected character %c", left)
}

func parseFull(left *Expr, buf *bufio.Reader) (*Expr, error) {
	var err error

	for {
		left, err = parseExpr(left, buf)
		if err != nil {
			return left, err
		}
	}
}

func parseExpr(left *Expr, buf *bufio.Reader) (*Expr, error) {
	op, err := consumeWhitespace(buf)
	if err != nil {
		return left, err
	}

	if op == ')' {
		return left, errClosingBracket
	}

	expr := &Expr{
		Op:   op,
		Left: left,
	}

	right, err := consumeWhitespace(buf)
	if err != nil {
		return expr, err
	}

	if right >= '0' && right <= '9' {
		expr.Right = &Expr{Const: int(right - '0')}
		return expr, nil
	}

	// this is an "("
	bracket, err := parse(buf)
	if err != errClosingBracket {
		return nil, errors.New("missing closing bracket")
	}
	expr.Right = &Expr{Op: '(', Left: bracket}

	return expr, nil
}

func consumeWhitespace(buf *bufio.Reader) (byte, error) {
	for {
		b, err := buf.ReadByte()
		if b == ' ' {
			continue
		}

		return b, err
	}
}

func readExprs(name string) ([]*Expr, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var exprs []*Expr
	for s.Scan() {
		expr, err := parseString(s.Text())
		if err != nil && err != io.EOF {
			fmt.Printf("Fehler beim Parsen von %s: %s\n", s.Text(), err)
			continue
		}

		exprs = append(exprs, expr)
	}

	return exprs, nil
}
