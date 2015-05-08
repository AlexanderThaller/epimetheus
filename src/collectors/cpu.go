// The MIT License (MIT)
//
// Copyright (c) 2015 Alexander Thaller
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package collectors

import (
	"strings"
	"time"

	"github.com/AlexanderThaller/logger"
	"github.com/juju/errgo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
)

func Cpu() error {
	collectors, err := cpu_register()
	if err != nil {
		return errgo.Notef(err, "can not register collectors for cpu")
	}

	go cpu_worker(collectors)

	return nil
}

func cpu_register() (map[string]prometheus.Gauge, error) {
	l := logger.New("collectors", "cpu", "register")

	stats, err := cpu.CPUTimes(false)
	if err != nil {
		return nil, errgo.New("can not get cpu stats")
	}

	l.Trace("Cpu Info: ", stats)

	collectors := make(map[string]prometheus.Gauge)

	for _, stat := range stats {
		l.Trace("Stat: ", stat)
		cpu := strings.Replace(stat.CPU, "-", "", -1)

		{
			name := "cpu_" + cpu + "_user"
			collector := prometheus.NewGauge(prometheus.GaugeOpts{
				Name: name,
				Help: "Current temperature of the CPU.",
			})

			collectors[name] = collector
			prometheus.MustRegister(collector)
			collector.Set(stat.User)
		}

		{
			name := "cpu_" + cpu + "_system"
			collector := prometheus.NewGauge(prometheus.GaugeOpts{
				Name: name,
				Help: "Current temperature of the CPU.",
			})

			collectors[name] = collector
			prometheus.MustRegister(collector)
			collector.Set(stat.System)
		}

		{
			name := "cpu_" + cpu + "_idle"
			collector := prometheus.NewGauge(prometheus.GaugeOpts{
				Name: name,
				Help: "Current temperature of the CPU.",
			})

			collectors[name] = collector
			prometheus.MustRegister(collector)
			collector.Set(stat.Idle)
		}

		{
			name := "cpu_" + cpu + "_nice"
			collector := prometheus.NewGauge(prometheus.GaugeOpts{
				Name: name,
				Help: "Current temperature of the CPU.",
			})

			collectors[name] = collector
			prometheus.MustRegister(collector)
			collector.Set(stat.Nice)
		}
	}

	return collectors, nil
}

func cpu_worker(collectors map[string]prometheus.Gauge) {
	l := logger.New("collectors", "cpu", "worker")

	for {
		stats, err := cpu.CPUTimes(false)
		if err != nil {
			l.Warning(errgo.New("can not get cpu stats"))
			continue
		}

		for _, stat := range stats {
			l.Trace("Stat: ", stat)
			cpu := strings.Replace(stat.CPU, "-", "", -1)

			{
				name := "cpu_" + cpu + "_user"
				collector := collectors[name]
				collector.Set(stat.User)
			}

			{
				name := "cpu_" + cpu + "_system"
				collector := collectors[name]
				collector.Set(stat.System)
			}

			{
				name := "cpu_" + cpu + "_idle"
				collector := collectors[name]
				collector.Set(stat.Idle)
			}

			{
				name := "cpu_" + cpu + "_nice"
				collector := collectors[name]
				collector.Set(stat.Nice)
			}
		}
		time.Sleep(time.Second * 5)
	}
}
