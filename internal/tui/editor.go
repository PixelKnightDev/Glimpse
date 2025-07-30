/*
 * Glimpse - Editor Integration Module
 * Author: Pratyush Yadav <pratyushyadav0106@gmail.com>
 * Description: Cross-platform editor integration for opening files at specific lines
 * License: MIT License
 */
package tui

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenFileInEditor(filename string, line int) error {
	var args []string
	var cmd string

	switch runtime.GOOS {
	case "darwin":
		if isCommandAvailable("code") {
			cmd = "code"
			args = []string{"-g", fmt.Sprintf("%s:%d", filename, line)}
		} else {
			cmd = "open"
			args = []string{filename}
		}
	case "linux":
		if isCommandAvailable("code") {
			cmd = "code"
			args = []string{"-g", fmt.Sprintf("%s:%d", filename, line)}
		} else {
			cmd = "xdg-open"
			args = []string{filename}
		}
	case "windows":
		if isCommandAvailable("code") {
			cmd = "code"
			args = []string{"-g", fmt.Sprintf("%s:%d", filename, line)}
		} else {
			cmd = "notepad"
			args = []string{filename}
		}
	default:
		return fmt.Errorf("unsupported OS")
	}

	go func() {
		exec.Command(cmd, args...).Start()
	}()

	return nil
}

func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
