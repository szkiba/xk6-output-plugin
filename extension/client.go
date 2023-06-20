// SPDX-FileCopyrightText: 2023 Iván Szkiba
//
// SPDX-License-Identifier: MIT

package extension

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"github.com/szkiba/xk6-output-plugin-go/output"
)

const pluginCommandPrefix = "xk6-output-plugin-"

func newPlugin(cmd string, logger logrus.FieldLogger) (output.Output, error) {
	args, err := shlex.Split(cmd)
	if err != nil {
		return nil, err
	}

	dir, file := filepath.Split(args[0])
	if !strings.HasPrefix(file, pluginCommandPrefix) {
		file = pluginCommandPrefix + file
	}

	args[0] = filepath.Join(dir, file)

	if dir == "./" {
		args[0] = dir + args[0]
	}

	hclogger := wrapLogger(logger.WithField("plugin", strings.TrimPrefix(filepath.Base(args[0]), pluginCommandPrefix)))

	if c, err := exec.LookPath(pluginCommandPrefix + args[0]); err == nil {
		args[0] = c
	}

	client := plugin.NewClient(&plugin.ClientConfig{ // nolint:exhaustruct
		HandshakeConfig:  output.Handshake,
		Plugins:          output.PluginMap,
		Cmd:              exec.Command(args[0], args[1:]...), // nolint:gosec
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		Logger:           hclogger,
	})

	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("output")
	if err != nil {
		return nil, err
	}

	if out, ok := raw.(output.Output); ok {
		return out, nil
	}

	return nil, ErrInvalidPlugin
}

var ErrInvalidPlugin = errors.New("invalid plugin")
