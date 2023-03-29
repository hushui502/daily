package main

import "fmt"

func main() {
	invoice := Invoice{
		customer: "BigCo",
		performances: []Performance{
			{
				playID:   "hamlet",
				audience: 55,
			},
			{
				playID:   "as-like",
				audience: 35,
			},
			{
				playID:   "othello",
				audience: 40,
			},
		},
	}

	fmt.Println(statement(invoice))
}
