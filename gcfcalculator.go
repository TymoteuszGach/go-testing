package main

type EuclideanGCFCalculator struct{}

func NewEuclideanGCFCalculator() *EuclideanGCFCalculator {
	return &EuclideanGCFCalculator{}
}

func (calculator *EuclideanGCFCalculator) Calculate(number1 int64, number2 int64) int64 {
	if number1 == 0 {
		return number2
	}
	return calculator.Calculate(number2%number1, number1)
}
