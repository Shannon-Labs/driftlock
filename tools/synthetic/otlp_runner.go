package main

import (
	"errors"
	"fmt"
	"log"
)

func runOTLPBatch(generator *OTLPGenerator, normalCount, anomalousCount int) error {
	var errs []error

	if normalCount > 0 {
		if err := generator.SendToCollector(generator.GenerateNormalLogs(normalCount)); err != nil {
			errs = append(errs, fmt.Errorf("normal logs: %w", err))
		} else {
			log.Printf("Sent %d normal logs to collector", normalCount)
		}

		if err := generator.SendToCollector(generator.GenerateNormalMetrics(normalCount)); err != nil {
			errs = append(errs, fmt.Errorf("normal metrics: %w", err))
		} else {
			log.Printf("Sent %d normal metrics to collector", normalCount)
		}
	}

	if anomalousCount > 0 {
		if err := generator.SendToCollector(generator.GenerateAnomalousLogs(anomalousCount)); err != nil {
			errs = append(errs, fmt.Errorf("anomalous logs: %w", err))
		} else {
			log.Printf("Sent %d anomalous logs to collector", anomalousCount)
		}

		if err := generator.SendToCollector(generator.GenerateAnomalousMetrics(anomalousCount)); err != nil {
			errs = append(errs, fmt.Errorf("anomalous metrics: %w", err))
		} else {
			log.Printf("Sent %d anomalous metrics to collector", anomalousCount)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
