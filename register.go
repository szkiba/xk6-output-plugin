// SPDX-FileCopyrightText: 2021 - 2023 Iván Szkiba
//
// SPDX-License-Identifier: MIT

package outputplugin

import (
	"github.com/szkiba/xk6-output-plugin/extension"
	"go.k6.io/k6/output"
)

// Register the extensions on module initialization.
func init() { // nolint:gochecknoinits
	output.RegisterExtension("plugin", New)
}

func New(params output.Params) (output.Output, error) { //nolint:ireturn
	return extension.New(params)
}
