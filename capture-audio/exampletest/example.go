package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func main() {
	os.Remove("musica.wav")
	fmt.Println("Init Record")
	cap, err := CaptureSpeaker()
	if err != nil {
		return
	}

	time.Sleep(50 * time.Second)
	cap.Cmd.Process.Kill()

	fmt.Println("finish")
	time.Sleep(5 * time.Second)
}

func CaptureSpeaker() (*recordModule, error) {

	//var stdin = new(bytes.Buffer)
	cmd := exec.Command("ffmpeg",
		"-f", "pulse",
		"-i", "alsa_output.pci-0000_00_1b.0.analog-stereo.monitor",
		//"salida.wav",
		"-f", "wav",
		"pipe:1",
	)
	pipe, _ := cmd.StdoutPipe()
	var stderr = new(bytes.Buffer)
	var stdin = new(bytes.Buffer)
	//cmd.Stdout = out
	cmd.Stderr = stderr
	cmd.Stdin = stdin

	err := cmd.Start()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		os.Exit(1)

	}

	go func() {
		file, _ := os.Create("musica.wav")
		scanerr := bufio.NewScanner(pipe)
		scanerr.Split(bufio.ScanBytes)
		for scanerr.Scan() {
			file.Write(scanerr.Bytes())
		}
	}()

	go func() {
		err := cmd.Wait()
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}()

	// reader, writer := io.Pipe()
	// cmd.Stdout = writer
	//err := cmd.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return &recordModule{
		Cmd: cmd,
	}, nil
}

type recordModule struct {
	Buffer io.Reader
	Cmd    *exec.Cmd
}
