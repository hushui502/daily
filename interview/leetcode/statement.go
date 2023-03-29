package main

import "strconv"

type Play struct {
	name     string
	playType string
}

type Performance struct {
	playID   string
	audience int
}

type Invoice struct {
	customer     string
	performances []Performance
}

var (
	plays = map[string]Play{
		"hamlet":    {"Hamlet", "tragedy"},
		"as-like":   {"As You Like It", "comedy"},
		"othello":   {"Othello", "tragedy"},
		"king lear": {"King Lear", "tragedy"},
	}
)

func statement(invoice Invoice) string {
	var totalAmount float64
	var volumeCredits int
	result := "Statement for " + invoice.customer + ""

	for _, perf := range invoice.performances {
		thisAmount := amountFor(perf, playFor(perf))

		// add volume credits
		volumeCredits += max(perf.audience-30, 0)
		// add extra credit for every ten comedy attendees
		if "comedy" == playFor(perf).playType {
			volumeCredits += int(perf.audience / 5)
		}
		// print line for this order
		result += " " + playFor(perf).name + ": " + format(thisAmount/100) + " (" + strconv.Itoa(perf.audience) + " seats)"
		totalAmount += thisAmount
	}

	result += "Amount owed is " + format(totalAmount/100)
	result += "You earned " + strconv.Itoa(volumeCredits) + " credits"
	return result
}

func format(amount float64) string {
	return "$" + strconv.Itoa(int(amount/100)) + ".00"
}

func playFor(aPerformance Performance) Play {
	return plays[aPerformance.playID]
}

func amountFor(aPerformance Performance, play Play) float64 {
	result := 0.0

	switch play.playType {
	case "tragedy":
		result = 40000
		if aPerformance.audience > 30 {
			result += 1000 * float64(aPerformance.audience-30)
		}
	case "comedy":
		result = 30000
		if aPerformance.audience > 20 {
			result += 10000 + 500*float64(aPerformance.audience-20)
		}
		result += 300 * float64(aPerformance.audience)
	default:
		panic("unknown type: " + play.playType)
	}

	return result
}
