// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
	"context"

	"go.opentelemetry.io/otel/metric"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gmetric"
)

// localCounterPerformer is an implementer for interface CounterPerformer.
type localCounterPerformer struct {
	metric.Float64ObservableCounter
	baseObservePerformer iBaseObservePerformer
}

// newCounterPerformer creates and returns a CounterPerformer that truly takes action to implement Counter.
func newCounterPerformer(meter metric.Meter, config gmetric.CounterConfig) (gmetric.CounterPerformer, error) {
	var (
		baseObservePerformer = newBaseObservePerformer(config.MetricConfig)
		options              = []metric.Float64ObservableCounterOption{
			metric.WithDescription(config.Help),
			metric.WithUnit(config.Unit),
		}
	)
	if !hasGlobalCallbackMetricSet.Contains(config.MetricKey()) {
		callback := metric.WithFloat64Callback(func(ctx context.Context, observer metric.Float64Observer) error {
			if config.Callback != nil {
				result, err := config.Callback(ctx)
				if err != nil {
					return gerror.WrapCodef(
						gcode.CodeOperationFailed,
						err,
						`callback failed for metric "%s"`, config.Name,
					)
				}
				if result != nil {
					observer.Observe(
						result.Value,
						baseObservePerformer.MergeAttributesToObserveOptions(result.Attributes)...,
					)
				}
				return nil
			}

			observer.Observe(
				baseObservePerformer.GetValue(),
				baseObservePerformer.GetObserveOptions()...,
			)
			return nil
		})
		options = append(options, callback)
	}
	counter, err := meter.Float64ObservableCounter(config.Name, options...)
	if err != nil {
		return nil, gerror.WrapCodef(
			gcode.CodeInternalError,
			err,
			`create Float64ObservableCounter failed with config: %+v`,
			config,
		)
	}
	return &localCounterPerformer{
		Float64ObservableCounter: counter,
		baseObservePerformer:     baseObservePerformer,
	}, nil
}

// Inc increments the counter by 1. Use Add to increment it by arbitrary
// non-negative values.
func (l *localCounterPerformer) Inc(option ...gmetric.Option) {
	l.Add(1, option...)
}

// Add adds the given value to the counter. It panics if the value is < 0.
func (l *localCounterPerformer) Add(increment float64, option ...gmetric.Option) {
	l.baseObservePerformer.AddValue(increment)
	l.baseObservePerformer.SetObserveOptionsByOption(option...)
}

// RemoveCallback removes the callback when global callback is defined on metric.
func (l *localCounterPerformer) RemoveCallback() {}
