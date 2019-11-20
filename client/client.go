package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var timeout string

// RootCmd asdf
var RootCmd = &cobra.Command{
	Use:   "info",
	Short: "Telnet-client",
}

// TelnetClientCmd init
var TelnetClientCmd = &cobra.Command{
	Use:   "go-telnet",
	Short: "Run telnet client",
	Run: func(cmd *cobra.Command, args []string) {
		var serverAddr strings.Builder
		switch len(args) {
		case 1:
			serverAddr.WriteString(args[0])
		case 2:
			serverAddr.WriteString(args[0])
			serverAddr.WriteString(":")
			serverAddr.WriteString(args[1])
		default:
			log.Fatalf("wrong server params")
		}

		runes := []rune(timeout)
		if len(runes) < 2 {
			log.Fatalf("wrong timeout params")
		}
		time, err := strconv.Atoi(string(runes[:len(runes)-1]))
		if err != nil {
			log.Fatalf("can not parse timeout params")
		}

		fmt.Printf("start %s, %d\n", serverAddr.String(), time)
		server(serverAddr.String(), time)

	},
}

func main() {
	TelnetClientCmd.Flags().StringVar(&timeout, "timeout", "10s", "default 10s")
	RootCmd.AddCommand()
	RootCmd.AddCommand(TelnetClientCmd)
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func readRoutine(ctx context.Context, cancel context.CancelFunc, scanner chan (string)) {

OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		case text := <-scanner:
			log.Printf("From server: %s", text)
		}
	}
	log.Printf("Finished readRoutine")
}

func writeRoutine(ctx context.Context, cancel context.CancelFunc, conn net.Conn, scanner chan (string)) {
OUTER3:
	for {
		select {
		case <-ctx.Done():
			cancel()
			break OUTER3
		case text := <-scanner:
			log.Printf("To server %v\n", text)

			conn.Write([]byte(fmt.Sprintf("%s\n", text)))
		}
	}
	log.Printf("Finished writeRoutine")
}

func server(serverAddr string, timeout int) {
	dialer := &net.Dialer{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	scannerIn := make(chan string)
	scannerOut := make(chan string)

	conn, err := dialer.DialContext(ctx, "tcp", serverAddr)

	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}

	wg := sync.WaitGroup{}

	go func() {
		scanner := bufio.NewScanner(conn)
	OUTER2:
		for {
			if !scanner.Scan() {
				cancel()
				break OUTER2
			}
			str := scanner.Text()
			scannerIn <- str
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdout)
	OUTER2:
		for {
			if !scanner.Scan() {
				cancel()
				break OUTER2
			}
			str := scanner.Text()
			scannerOut <- str
		}
	}()

	wg.Add(1)
	go func() {
		readRoutine(ctx, cancel, scannerIn)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		writeRoutine(ctx, cancel, conn, scannerOut)
		wg.Done()
	}()

	wg.Wait()
	conn.Close()
}
