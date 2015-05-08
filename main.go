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

package main

import (
	"net/http"
	"os"
	"path"

	"github.com/AlexanderThaller/logger"
	"github.com/juju/errgo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"

	"github.com/AlexanderThaller/epimetheus/src/collectors"
)

const (
	Name = "epimetheus"
)

func main() {
	l := logger.New(Name, "main")

	err := configure()
	if err != nil {
		alert_and_debug(l, err, "can not configure application")
		os.Exit(1)
	}

	err = startCollectors()
	if err != nil {
		alert_and_debug(l, err, "can not start collectors")
		os.Exit(1)
	}

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(viper.GetString("ListenOn"), nil)
}

func alert_and_debug(l logger.Logger, err error, message string) {
	l.Alert(message, ": ", err.Error())
	l.Debug(message, ": ", errgo.Details(err))
}

func configure() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(path.Join("/etc/", Name))
	viper.AddConfigPath(path.Join("/usr/local/etc/", Name))
	viper.AddConfigPath(path.Join("$HOME", Name))

	err := viper.ReadInConfig()
	if err != nil {
		return errgo.Notef(err, "can not access config file")
	}

	viper.SetDefault("ListenOn", ":8080")

	return nil
}

func startCollectors() error {
	err := collectors.Load()
	if err != nil {
		return errgo.Notef(err, "can not start load collector")
	}

	return nil
}
