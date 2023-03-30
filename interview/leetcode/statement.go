package leetcode

import "strconv"

type Play struct {
	name     string
	playType string
}

type Performance struct {
	playID   string
	audience int
	amount   int
	play     Play
}

type Invoice struct {
	customer     string
	performances []Performance
}

type statementData struct {
	customer           string
	performances       []Performance
	totalAmount        int
	totalVolumeCredits int
}

func htmlStatement(invoice Invoice) string {
	return renderHtml(createStatementData(invoice))
}

func renderHtml(data statementData) string {
	result := "<h1>Statement for " + data.customer + "</h1>\n"
	result += "<table>\n"
	result += "<tr><th>play</th><th>seats</th><th>cost</th></tr>"

	for _, perf := range data.performances {
		// print line for this order
		result += " <tr><td>" + perf.play.name + "</td><td>" + strconv.Itoa(perf.audience) + "</td>"
		result += "<td>" + usd(float64(perf.amount)) + "</td></tr>\n"
	}

	result += "</table>\n"
	result += "<p>Amount owed is <em>" + usd(float64(data.totalAmount)) + "</em></p>\n"
	result += "<p>You earned <em>" + strconv.Itoa(data.totalVolumeCredits) + "</em> credits</p>\n"

	return result
}

func statement(invoice Invoice) string {
	return renderPlainText(createStatementData(invoice))
}

func usd(amount float64) string {
	return "$" + strconv.Itoa(int(amount/100)) + ".00"
}
