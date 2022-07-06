package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {

	// new server virtual MIC
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		fmt.Println("Wait client: ")
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("New client: ")
		clientPlayInVirtMic(conn)
	}
}

func clientPlayInVirtMic(conn net.Conn) {
	var device = "/tmp/andruscodexmic"

	defer conn.Close()
	file, err := syscall.Open(device, syscall.O_RDWR|syscall.O_CLOEXEC|syscall.O_NONBLOCK, 0644)
	if err != nil {
		panic(err)
	}

	for true {
		// read buff audio ffmpeg
		var buff = make([]byte, 256)
		n, err := conn.Read(buff)

		if err != nil {
			fmt.Println("Error read audio, ", err)
			break
		}

		// write audio in mic
		_, err = syscall.Write(file, buff[:n])
		if err != nil {
			if strings.Contains(err.Error(), "temporarily unavailable") {
				continue
			}

			fmt.Println("Error write audio, ", err)
			break
		}

	}

	return

	// Create virtual MIC client
	file, err = syscall.Open(device, syscall.O_RDWR|syscall.O_CLOEXEC, 0644)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(file)
	defer conn.Close()

	// convert audio in LINUX pipe
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-fflags",
		"nobuffer",
		// Input
		"-f",
		"s16le",
		"-ar",
		"44100",
		"-ac", "2",
		"-i",
		"pipe:0",
		// output
		"-f",
		"s16le",
		"-codec", "copy",
		"pipe:1",
	)

	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	writeDin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	readDout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	go func() {
		io.Copy(writeDin, conn)
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()
	go func() {
		var loop = true
		defer readDout.Close()
		for loop {

			// read buff audio ffmpeg
			var buff = make([]byte, 512)
			n, err := readDout.Read(buff)
			if err != nil {
				fmt.Println("Error read audio, ", err)
				break
			}

			// write audio in mic
			_, err = syscall.Write(file, buff[:n])
			if err != nil {
				fmt.Println("Error write audio, ", err)
				break
			}
		}
		log.Println("Write close")

	}()

	cmd.Run()

}

// func clientPlayInAlsa(conn net.Conn) {

// 	defer conn.Close()

// 	cmd := exec.Command(
// 		"ffmpeg",
// 		"-fflags",
// 		"nobuffer",
// 		// Input
// 		"-f",
// 		"s16le",
// 		"-ar",
// 		"41k",
// 		"-ac", "1",
// 		"-i",
// 		"pipe:0",
// 		// output
// 		"-f",
// 		"alsa",
// 		"default",
// 	)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	pipeReader, pipeWriter := io.Pipe()
// 	cmd.Stdin = pipeReader

// 	go func() {
// 		defer pipeWriter.Close()
// 		io.Copy(pipeWriter, conn)
// 	}()

// 	fmt.Println("Termino: ", cmd.Run())
// }
