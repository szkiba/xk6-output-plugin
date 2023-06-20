// SPDX-FileCopyrightText: 2023 Iván Szkiba
//
// SPDX-License-Identifier: MIT

//go:build mage
// +build mage

package main

import (
	"path/filepath"
	"strings"

	"github.com/magefile/mage/sh"
	"github.com/princjef/mageutil/bintool"
	"github.com/princjef/mageutil/shellcmd"
)

var Default = All

var linter = bintool.Must(bintool.New(
	"golangci-lint{{.BinExt}}",
	"1.51.1",
	"https://github.com/golangci/golangci-lint/releases/download/v{{.Version}}/golangci-lint-{{.Version}}-{{.GOOS}}-{{.GOARCH}}{{.ArchiveExt}}",
))

func Lint() error {
	if err := linter.Ensure(); err != nil {
		return err
	}

	return linter.Command(`run`).Run()
}

func Test() error {
	return shellcmd.Command(`go test -count 1 -coverprofile=coverage.txt ./...`).Run()
}

func Build() error {
	return shellcmd.Command(`xk6 build --with github.com/szkiba/xk6-output-plugin=.`).Run()
}

func Example() error {
	return shellcmd.Command(`go build -ldflags="-s -w" -o xk6-output-plugin-example ./examples/example-go`).Run()
}

func Coverage() error {
	return shellcmd.Command(`go tool cover -html=coverage.txt`).Run()
}

func glob(patterns ...string) (string, error) {
	buff := new(strings.Builder)

	for _, p := range patterns {
		m, err := filepath.Glob(p)
		if err != nil {
			return "", err
		}

		_, err = buff.WriteString(strings.Join(m, " ") + " ")
		if err != nil {
			return "", err
		}
	}

	return buff.String(), nil
}

func License() error {
	all, err := glob(
		"*.go",
		"*/*.go",
		".*.yml",
		".gitignore",
		"*/.gitignore",
		".github/workflows/*",
	)
	if err != nil {
		return err
	}

	return shellcmd.Command(
		`reuse annotate --copyright "Iván Szkiba" --merge-copyrights --license MIT --skip-unrecognised ` + all,
	).Run()
}

func Clean() error {
	sh.Rm("magefiles/bin")
	sh.Rm("coverage.txt")
	sh.Rm("bin")
	sh.Rm("k6")

	return nil
}

func All() error {
	if err := Lint(); err != nil {
		return err
	}

	return Build()
}

func xk6run(arg string) shellcmd.Command {
	return shellcmd.Command("xk6 run --quiet --no-summary --no-usage-report " + arg)
}

func Run() error {
	return xk6run(`-d 1m --out plugin=./example script.js`).Run()
}
