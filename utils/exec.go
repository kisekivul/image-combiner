package utils

import (
	"io"
	"os/exec"
)

func Exec(cmds []string) (outByte []byte, errByte []byte, err error) {
	cmd := exec.Command(cmds[0])
	if len(cmds) > 1 {
		cmd.Args = append(cmd.Args, cmds[1:]...)
	}

	var (
		stdout, stderr io.ReadCloser
	)

	if stdout, err = cmd.StdoutPipe(); err != nil {
		return
	}
	defer stdout.Close()

	if stderr, err = cmd.StderrPipe(); err != nil {
		return
	}
	defer stderr.Close()

	if err = cmd.Start(); err != nil {
		return
	}

	if outByte, err = io.ReadAll(stdout); err != nil {
		return
	}

	if errByte, err = io.ReadAll(stderr); err != nil {
		return
	}
	err = cmd.Wait()
	return
}
