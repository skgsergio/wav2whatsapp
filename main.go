// Copyright (c) 2021 Sergio Conde &lt;skgsergio@gmail.com&gt;
//
// This program is free software: you can redistribute it and/or modify it under
// the terms of the GNU General Public License as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any later
// version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
// PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with
// this program. If not, see <https://www.gnu.org/licenses/>.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gen2brain/dlgs"
)

func convertToMonoOpusOGG(input, output string) ([]byte, []byte, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("opusenc", "--downmix-mono", input, output)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

func init() {
	// Quick and dirty hack for adding homebrew path in macOS APP
	os.Setenv("PATH", fmt.Sprintf("%s:/usr/local/bin", os.Getenv("PATH")))
}

func main() {
	// File selection dialog
	fileName, ok, err := dlgs.File(
		"Select audio file",
		"*.wav *.flac *.aif *.aiff *.raw *.pcm",
		false,
	)
	if err != nil {
		log.Panic(err)
	}

	if !ok {
		_, err = dlgs.Error("Error", "No file selected.")
		if err != nil {
			log.Panic(err)
		}
		os.Exit(1)
	}

	// Run conversion, preserve filename changing extension
	stdout, stderr, err := convertToMonoOpusOGG(
		fileName,
		fmt.Sprintf("%s.ogg", strings.TrimSuffix(fileName, filepath.Ext(fileName))),
	)
	if err != nil {
		errMsg := err.Error() + "\n\nPATH: " + os.Getenv("PATH")

		if _, ok := err.(*exec.ExitError); ok {
			errMsg = fmt.Sprintf("%s\n%s", stdout, stderr)
		}

		_, err = dlgs.Error("Opusenc runtime error", errMsg)
		if err != nil {
			log.Panic(err)
		}
		os.Exit(1)
	}

	// Display opusenc output
	_, err = dlgs.Info("File converted", fmt.Sprintf("%s\n%s", stdout, stderr))
	if err != nil {
		log.Panic(err)
	}

	// Open file explorer to the file directory
	err = open(filepath.Dir(fileName))
	if err != nil {
		log.Panic(err)
	}
}
