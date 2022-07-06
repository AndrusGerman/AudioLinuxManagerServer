package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func CaptureSpeaker(device string) (*recordModule, error) {
	cmd := exec.Command("ffmpeg",
		"-fflags",
		"nobuffer",
		// input
		"-f", "pulse",
		"-i", device,
		// output
		"-f", "wav",
		"pipe:1",
	)
	pipe, _ := cmd.StdoutPipe()
	var stderr = new(bytes.Buffer)
	var stdin = new(bytes.Buffer)
	cmd.Stderr = stderr
	cmd.Stdin = stdin

	err := cmd.Start()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		os.Exit(1)
	}

	go func() {
		err := cmd.Wait()
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}()

	return &recordModule{
		cmd:    cmd,
		Buffer: pipe,
	}, nil
}

type recordModule struct {
	Buffer io.ReadCloser
	cmd    *exec.Cmd
}

func (ctx *recordModule) Close() {
	if ctx.cmd.Process != nil {
		go ctx.cmd.Process.Kill()
	}
	go ctx.Buffer.Close()
}
