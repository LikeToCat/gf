// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localMetricInstrument implements interface MetricInstrument.
type localMetricInstrument struct {
	name    string
	version string
}

// newMetricInstrument creates and returns a MetricInstrument.
func newMetricInstrument(name, version string) MetricInstrument {
	return &localMetricInstrument{
		name:    name,
		version: version,
	}
}

// Name returns the instrument name of the metric.
func (l *localMetricInstrument) Name() string {
	return l.name
}

// Version returns the instrument version of the metric.
func (l *localMetricInstrument) Version() string {
	return l.version
}
