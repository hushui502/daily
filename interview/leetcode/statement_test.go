package main

import (
	"reflect"
	"testing"
)

func Test_amountFor(t *testing.T) {
	type args struct {
		aPerformance Performance
		play         Play
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				aPerformance: Performance{
					playID:   "hamlet",
					audience: 55,
				},
				play: Play{
					name:     "Hamlet",
					playType: "tragedy",
				},
			},
			want: 65000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := amountFor(tt.args.aPerformance, tt.args.play); got != tt.want {
				t.Errorf("amountFor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_format(t *testing.T) {
	type args struct {
		amount float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				amount: 65000,
			},
			want: "$650.00",
		},
		{
			name: "test2",
			args: args{},
			want: "$0.00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := format(tt.args.amount); got != tt.want {
				t.Errorf("format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_playFor(t *testing.T) {
	type args struct {
		aPerformance Performance
	}
	tests := []struct {
		name string
		args args
		want Play
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				aPerformance: Performance{
					playID:   "hamlet",
					audience: 55,
				},
			},
			want: Play{
				name:     "Hamlet",
				playType: "tragedy",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := playFor(tt.args.aPerformance); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("playFor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_statement(t *testing.T) {
	type args struct {
		invoice Invoice
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				invoice: Invoice{
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
				},
			},
			want: `Statement for BigCo
 Hamlet: $650.00 (55 seats)
 As You Like It: $580.00 (35 seats)
 Othello: $500.00 (40 seats)
Amount owed is $1,730.00
You earned 47 credits
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := statement(tt.args.invoice); got != tt.want {
				t.Errorf("statement() = %v, want %v", got, tt.want)
			}
		})
	}
}
