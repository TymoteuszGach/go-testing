package main

import (
	"fmt"
)

type GCFCalculator interface{
	Calculate(number1 int64, number2 int64) int64
}

type GCFRequestProcessor struct{
	dbAdapter DatabaseAdapter
	calculator GCFCalculator
}

func NewGCFRequestProcessor(dbAdapter DatabaseAdapter, calculator GCFCalculator) *GCFRequestProcessor {
	return &GCFRequestProcessor{dbAdapter: dbAdapter, calculator: calculator}
}

func (processor *GCFRequestProcessor) Process(request GCFRequest) error {
	gcf := processor.calculator.Calculate(request.Number1, request.Number2)

	gcfResult := GCFResult{
		Number1: request.Number1,
		Number2: request.Number2,
		GCF:     gcf,
	}

	err := processor.dbAdapter.SaveGCFResult(gcfResult)
	if err != nil {
		return fmt.Errorf("cannot save result to database: %v", err)
	}

	return nil
}