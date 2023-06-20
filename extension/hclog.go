// SPDX-FileCopyrightText: 2023 Iván Szkiba
//
// SPDX-License-Identifier: MIT

package extension

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/sirupsen/logrus"
)

const fieldPrefixJSON = "\t\n"

type sinkAdapterToLogrus struct {
	log logrus.FieldLogger
}

func (a *sinkAdapterToLogrus) Accept(name string, level hclog.Level, msg string, args ...interface{}) {
	logger := a.log

	last := len(args) - 1

	for idx := 0; idx < len(args); idx++ {
		if idx >= last {
			logger = logger.WithField(hclog.MissingKey, fmt.Sprint(args[idx]))

			continue
		}

		if s, ok := args[idx].(string); ok {
			if s == "timestamp" {
				idx++

				continue
			}
		}

		if s, ok := args[idx+1].(string); ok {
			if strings.HasPrefix(s, fieldPrefixJSON) {
				m := map[string]interface{}{}

				if err := json.Unmarshal([]byte(s[len(fieldPrefixJSON):]), &m); err == nil {
					logger = logger.WithField(fmt.Sprint(args[idx]), m)
					idx++

					continue
				}
			}
		}

		logger = logger.WithField(fmt.Sprint(args[idx]), fmt.Sprint(args[idx+1]))
		idx++
	}

	logfunc(logger, level)(msg)
}

func logfunc(logger logrus.FieldLogger, level hclog.Level) func(...interface{}) {
	switch level {
	case hclog.Debug, hclog.Trace:
		return logger.Debug
	case hclog.Warn:
		return logger.Warn
	case hclog.Error:
		return logger.Error
	case hclog.Info, hclog.DefaultLevel, hclog.NoLevel:
		return logger.Info
	case hclog.Off:
		return func(i ...interface{}) {}
	default:
		return logger.Info
	}
}

func wrapLogger(l logrus.FieldLogger) hclog.Logger {
	h := hclog.NewInterceptLogger(&hclog.LoggerOptions{Output: io.Discard}) // nolint:exhaustruct

	h.RegisterSink(&sinkAdapterToLogrus{log: l})

	return h
}
