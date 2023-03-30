package leetcode

import "strconv"

func createStatementData(invoice Invoice) statementData {
	var statementData statementData
	statementData.customer = invoice.customer
	enrichPerformance(statementData, invoice)

	statementData.totalVolumeCredits = totalVolumeCredits(statementData)
	statementData.totalAmount = int(totalAmount(statementData))

	return statementData
}

func renderPlainText(data statementData) string {
	result := "Statement for " + data.customer + "\n"

	for _, perf := range data.performances {
		// print line for this order
		result += " " + perf.play.name + ": " + usd(float64(perf.amount)) + " (" + strconv.Itoa(perf.audience) + " seats)\n"
	}

	result += "Amount owed is " + usd(float64(data.totalAmount)) + "\n"
	result += "You earned " + strconv.Itoa(data.totalVolumeCredits) + " credits\n"

	return result
}

func totalAmount(data statementData) float64 {
	var result float64
	for _, perf := range data.performances {
		result += amountFor(perf)
	}

	return result
}

func enrichPerformance(data statementData, invoice Invoice) {
	for i, perf := range invoice.performances {
		data.performances[i].play = playFor(perf)
		data.performances[i].audience = perf.audience
		data.performances[i].amount = int(amountFor(perf))
	}
}

func totalVolumeCredits(data statementData) int {
	var result int
	for _, perf := range data.performances {
		// add volume credits
		result += volumeCreditsFor(perf)
	}

	return result
}

func volumeCreditsFor(perf Performance) int {
	volumeCredits := 0
	volumeCredits += int(perf.audience)
	// add extra credit for every ten comedy attendees
	if "comedy" == playFor(perf).playType {
		volumeCredits += int(perf.audience / 5)
	}

	return volumeCredits
}

func playFor(aPerformance Performance) Play {
	return plays[aPerformance.playID]
}

func amountFor(aPerformance Performance) float64 {
	result := 0.0

	switch playFor(aPerformance).playType {
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
		panic("unknown type: " + playFor(aPerformance).playType)
	}

	return result
}
