package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
)

const (
	maxVal = 100
)

type Expression struct {
	val1       int
	val2       int
	op         string
	answers    string
	expression string

	m sync.Mutex
}

func NewExpression() *Expression {
	e := Expression{}
	e.generateExpression()
	return &e
}

func (e *Expression) generateExpression() {
	e.m.Lock()
	defer e.m.Unlock()

	e.val1 = e.generateValue()
	e.val2 = e.generateValue()
	e.op = e.generateOperation()
	res, err := e.getAnswer()

	if err != nil {
		log.Fatal(err)
	}

	e.answers = res

	e.expression = fmt.Sprintf("%d %s %d = ?", e.val1, e.op, e.val2)
}

func (e *Expression) generateValue() int {
	return rand.Intn(maxVal)
}

func (e *Expression) generateOperation() string {
	availableOperations := []string{"+", "-", "*"}

	idx := rand.Intn(len(availableOperations))

	return availableOperations[idx]
}

func (e *Expression) getAnswer() (string, error) {
	switch e.op {
	case "+":
		return fmt.Sprint(e.val1 + e.val2), nil
	case "-":
		return fmt.Sprint(e.val1 - e.val2), nil
	case "*":
		return fmt.Sprint(e.val1 * e.val2), nil

	default:
		return "", fmt.Errorf("operation %s does not supported", e.op)
	}
}

func (e *Expression) GetAnswer() string {
	e.m.Lock()
	defer e.m.Unlock()

	return e.answers
}
