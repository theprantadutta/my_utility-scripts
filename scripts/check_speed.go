package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/showwin/speedtest-go/speedtest"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var speedTestClient = speedtest.New()

	// Handle interruption signal (Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		<-interrupt
		color.HiMagenta("\nExiting...")
		wg.Done()
		os.Exit(0)
	}()

	color.HiCyan("Testing download and upload speeds...")

	serverList, _ := speedTestClient.FetchServers()
	targets, _ := serverList.FindServer([]int{})

	for _, s := range targets {
		color.HiBlue(fmt.Sprintf("Testing server %s...\n", s.Name))

		// Please make sure your host can access this test server,
		// otherwise you will get an error.
		// It is recommended to replace a server at this time
		color.HiGreen("Checking Latency...")
		err := s.PingTest(nil)
		if err != nil {
			color.HiRed("Error testing latency:", err)
			return
		}

		color.HiGreen(fmt.Sprintf("Latency: %s", s.Latency))

		color.HiCyan("Checking Download Speed...")
		err = s.DownloadTest()
		if err != nil {
			color.HiRed("Error testing download speed:", err)
			return
		}
		color.HiCyan(fmt.Sprintf("Download Speed: %.2f Mbps", s.DLSpeed))

		color.HiCyan("Checking Upload Speed...")
		err = s.UploadTest()
		if err != nil {
			color.HiRed("Error testing upload speed:", err)
			return
		}
		color.HiCyan(fmt.Sprintf("Upload Speed: %.2f Mbps", s.ULSpeed))

		color.HiCyan(fmt.Sprintf("| %-15s | %-15s | %-15s | %-15s |\n", "Server Name", "Latency", "Download Speed", "Upload Speed"))
		color.HiCyan(fmt.Sprintf("| %-15s | %-15s | %-15s | %-15s |\n", s.Name, s.Latency, fmt.Sprintf("%.2f Mbps", s.DLSpeed), fmt.Sprintf("%.2f Mbps", s.ULSpeed)))

		// Reset counter
		s.Context.Reset()
	}

	color.HiMagenta("Press Ctrl+C to exit.")
	wg.Wait()
}
