package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func PlaySpeaker(pipe io.Reader) error {
	cmd := exec.Command("pw-play", "-")
	var stderr = new(bytes.Buffer)
	cmd.Stderr = stderr
	cmd.Stdin = pipe

	err := cmd.Start()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		os.Exit(1)
	}

	go func() {
		err := cmd.Wait()
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}()
	return nil
}
