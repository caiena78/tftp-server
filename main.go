package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pin/tftp"
)

var storagePath = ""

func checkpath(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Fatal("Folder ( " + path + " ) does not exist.")
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
	tempDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	var addr = flag.String("addr", "<ALL>", "IP address")
	var port = flag.String("port", "69", "UDP port for the tftp server")
	var dir = flag.String("dir", tempDir, "The Directory to service with the TFTP Server")
	flag.Parse()
	fmt.Print("\n\ntftp-server is using the following options\n\n")
	fmt.Printf("-addr=%s\n", *addr)
	fmt.Printf("-port=UDP/%s\n", *port)
	fmt.Printf("-dir=%s\n", *dir)
	checkpath(*dir)
	storagePath = *dir

	s := tftp.NewServer(readHandler, writeHandler)
	s.SetTimeout(5 * time.Second)
	// break up the input ip and port and add it to the listen and serve method
	fmt.Println("tftp-server is running")
	if *addr == "<ALL>" {
		err = s.ListenAndServe(":69")
	} else {
		err = s.ListenAndServe(*addr + ":" + *port)
	}

	if err != nil {
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		os.Exit(1)
	}
}
