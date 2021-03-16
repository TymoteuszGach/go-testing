package main

import (
	"testing"
)

func TestEuclideanGCFCalculator_Calculate(t *testing.T) {
	type args struct {
		number1 int64
		number2 int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "two positive numbers",
			args: args{
				number1: 70,
				number2: 15,
			},
			want: 5,
		},
		{
			name: "first number is 1",
			args: args{
				number1: 1,
				number2: 14,
			},
			want: 1,
		},
		{
			name: "second number is 1",
			args: args{
				number1: 5,
				number2: 1,
			},
			want: 1,
		},
		{
			name: "first number is 0",
			args: args{
				number1: 0,
				number2: 4,
			},
			want: 4,
		},
		{
			name: "second number is 0",
			args: args{
				number1: 7,
				number2: 0,
			},
			want: 7,
		},
		{
			name: "both numbers are 0",
			args: args{
				number1: 0,
				number2: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := NewEuclideanGCFCalculator()
			if got := calculator.Calculate(tt.args.number1, tt.args.number2); got != tt.want {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
