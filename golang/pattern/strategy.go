package pattern

import "fmt"

type IStrategy interface {
	do(int, int) int
}

type add struct {

}

func (*add) do(a, b int) int {
	return a + b
}

type Operator struct {
	strategy IStrategy
}

func (o *Operator) setStrategy(strategy IStrategy) {
	o.strategy = strategy
}

func (o *Operator) calculate(a, b int) int {
	return o.strategy.do(a, b)
}

func test() {
	operator := Operator{}
	operator.setStrategy(&add{})
	result := operator.calculate(1, 3)
	fmt.Println(result)
}
