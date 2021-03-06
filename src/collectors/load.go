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
	"time"

	"github.com/AlexanderThaller/epimetheus/src/data"
	"github.com/AlexanderThaller/logger"
	"github.com/juju/errgo"
	"github.com/prometheus/client_golang/prometheus"
)

func Load() error {
	l := logger.New("collectors", "load")

	load1 := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "load_1",
		Help: "Load average over the last minute",
	})

	load5 := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "load_5",
		Help: "Load average over the last five minutes",
	})

	load15 := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "load_15",
		Help: "Load average over the last 15 minutes",
	})

	prometheus.MustRegister(load1)
	prometheus.MustRegister(load5)
	prometheus.MustRegister(load15)

	go func() {
		for {
			values, err := data.Cpu()
			if err != nil {
				l.Warning(errgo.Notef(err, "can not get load values"))
				time.Sleep(time.Second * 5)
				continue
			}

			load1.Set(values["load.01"])
			load5.Set(values["load.05"])
			load15.Set(values["load.15"])

			time.Sleep(time.Second * 5)
		}
	}()

	return nil
}
