package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pin/tftp"
)

var storagePath = ""

func isWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

func checkpath(path string) {
	_, err := os.Stat(path)
	colorRed := string("\033[31m")
	colorReset := string("\033[0m")
	if isWindows() {
		colorRed = ""
		colorReset = ""
	}
	if os.IsNotExist(err) {
		log.Fatal(colorRed + "Folder ( " + path + " ) does not exist." + colorReset)
	}
}

// readHandler is called when client starts file download from server
func readHandler(filename string, rf io.ReaderFrom) error {
	fmt.Printf("Sending File - %s : ", filename)
	filename = filepath.Join(storagePath, filename)
	file, err := os.Open(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	n, err := rf.ReadFrom(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes sent\n", n)
	return nil
}

// writeHandler is called when client starts file upload to server
func writeHandler(filename string, wt io.WriterTo) error {
	fmt.Printf("Receiving File - %s : ", filename)
	filename = filepath.Join(storagePath, filename)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	n, err := wt.WriteTo(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes received\n", n)
	return nil
}

func main() {
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	colorWhite := "\033[37m"
	if isWindows() {
		colorRed = ""
		colorGreen = ""
		colorReset = ""
		colorWhite = ""
	}

	tempDir, err := os.Getwd()
	if err != nil {
		fmt.Printf(string(colorRed))
		fmt.Println(err)
		fmt.Printf(string(colorReset))
	}

	var addr = flag.String("addr", "<ALL>", "IP address")
	var port = flag.String("port", "69", "UDP port for the tftp server")
	var dir = flag.String("dir", tempDir, "The Directory to service with the TFTP Server")
	flag.Parse()

	fmt.Print("\n\ntftp-server is using the following options\n\n")
	fmt.Printf("%s-addr=%s%s\n", string(colorWhite), string(colorGreen), *addr)
	fmt.Printf("%s-port=%sUDP/%s\n", string(colorWhite), string(colorGreen), *port)
	fmt.Printf("%s-dir=%s%s\n", string(colorWhite), string(colorGreen), *dir)
	checkpath(*dir)
	storagePath = *dir
	fmt.Printf(string(colorReset))
	s := tftp.NewServer(readHandler, writeHandler)
	s.SetTimeout(5 * time.Second)
	// break up the input ip and port and add it to the listen and serve method
	fmt.Println("tftp-server is " + string(colorGreen) + "running" + string(colorReset))
	if *addr == "<ALL>" {
		err = s.ListenAndServe(":69")
	} else {
		err = s.ListenAndServe(*addr + ":" + *port)
	}

	if err != nil {
		fmt.Printf(string(colorRed))
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		fmt.Printf(string(colorReset))
		os.Exit(1)
	}
}
