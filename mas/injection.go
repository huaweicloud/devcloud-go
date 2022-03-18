/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2021.
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License.  You may obtain a copy of the
 * License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 *
 */

package mas

import (
	"math/rand"
	"time"
)

// InjectionProperties chaos configuration
type InjectionProperties struct {
	Active         bool            `yaml:"active"`
	Duration       int             `yaml:"duration"`
	Interval       int             `yaml:"interval"`
	Percentage     int             `yaml:"percentage"`
	DelayInjection *DelayInjection `yaml:"delayInjection"`
	ErrorInjection *ErrorInjection `yaml:"errorInjection"`
}

// DelayInjection delay configuration
type DelayInjection struct {
	Active     bool `yaml:"active"`
	Percentage int  `yaml:"percentage"`
	TimeMs     int  `yaml:"timeMs"`
	JitterMs   int  `yaml:"jitterMs"`
}

// NewDelayInjection sda
func NewDelayInjection(active bool, percentage, timeMs, jitterMs int) *DelayInjection {
	return &DelayInjection{
		Active:     active,
		Percentage: percentage,
		TimeMs:     timeMs,
		JitterMs:   jitterMs,
	}
}

func (d *DelayInjection) checkActive() (int, bool) {
	if d.Active {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		if r.Intn(100) <= d.Percentage {
			jitterMs := r.Intn(2*d.JitterMs+1) - d.JitterMs
			if d.TimeMs+jitterMs < 0 {
				return 0, true
			}
			return d.TimeMs + jitterMs, true
		}
	}
	return 0, false
}

// InjectionError error details
type InjectionError struct {
	Err        error
	Percentage int
}

// ErrorInjection error configuration
type ErrorInjection struct {
	Active     bool `yaml:"active"`
	Percentage int  `yaml:"percentage"`
	errs       []*InjectionError
}

func NewErrorInjection(active bool, percentage int) *ErrorInjection {
	return &ErrorInjection{
		Active:     active,
		Percentage: percentage,
		errs:       make([]*InjectionError, 0),
	}
}

func (e *ErrorInjection) checkActive() (error, bool) {
	if e.Active {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		if r.Intn(100) <= e.Percentage {
			return e.errs[r.Intn(len(e.errs))].Err, true
		}
	}
	return nil, false
}

// InjectionDuration ingestion period details
type InjectionDuration struct {
	duration    int
	total       int
	startTimeMs int64
}

func NewInjectionDuration(duration, total int) *InjectionDuration {
	return &InjectionDuration{
		duration:    duration,
		total:       total,
		startTimeMs: time.Now().Unix(),
	}
}

func (i *InjectionDuration) checkActive() bool {
	return (time.Now().Unix()-i.startTimeMs)%int64(i.total) <= int64(i.duration)
}

// InjectionManagement chaos injection details
type InjectionManagement struct {
	active            bool
	injectionDuration *InjectionDuration
	percentage        int
	delayInjection    *DelayInjection
	errorInjection    *ErrorInjection
}

func CompliancePercentage(percentage int) int {
	if percentage < 0 {
		return 0
	}
	if percentage > 100 {
		return 100
	}
	return percentage
}

func NewInjectionManagement(chaos *InjectionProperties) *InjectionManagement {
	return &InjectionManagement{
		active:            chaos.Active,
		injectionDuration: NewInjectionDuration(chaos.Duration, chaos.Interval),
		percentage:        CompliancePercentage(chaos.Percentage),
		delayInjection: NewDelayInjection(chaos.DelayInjection.Active,
			CompliancePercentage(chaos.DelayInjection.Percentage),
			chaos.DelayInjection.TimeMs,
			chaos.DelayInjection.JitterMs),
		errorInjection: NewErrorInjection(chaos.ErrorInjection.Active,
			CompliancePercentage(chaos.ErrorInjection.Percentage)),
	}
}

func (i *InjectionManagement) SetError(errs []error) {
	if i.errorInjection.Active {
		if errs != nil && len(errs) > 0 {
			for _, err := range errs {
				i.errorInjection.errs = append(i.errorInjection.errs,
					&InjectionError{err, i.errorInjection.Percentage})
			}
		}
	}
}

func (i *InjectionManagement) AddError(errs []*InjectionError) {
	if i.errorInjection.Active {
		if errs != nil && len(errs) > 0 {
			for _, err := range errs {
				err.Percentage = CompliancePercentage(err.Percentage)
			}
			i.errorInjection.errs = append(i.errorInjection.errs, errs...)
		}
	}
}

// Inject chaos injection triggering
func (i *InjectionManagement) Inject() error {
	if i == nil {
		return nil
	}
	if !i.active {
		return nil
	}
	if !i.injectionDuration.checkActive() {
		return nil
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if r.Intn(100) <= i.percentage {
		if err, active := i.errorInjection.checkActive(); active {
			return err
		}
		if delay, active := i.delayInjection.checkActive(); active {
			time.Sleep(time.Millisecond * time.Duration(delay))
		}
	}
	return nil
}
