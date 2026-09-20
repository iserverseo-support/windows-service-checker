//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	exitRunning  = 0
	exitStopped  = 1
	exitNotFound = 2
	exitError    = 3
)

// checkService queries a Windows service using sc.exe.
func checkService(serviceName string) int {
	if strings.TrimSpace(serviceName) == "" {
		fmt.Println("ERROR: Service name cannot be empty")
		return exitError
	}

	cmd := exec.Command("sc.exe", "query", serviceName)
	output, err := cmd.CombinedOutput()

	result := string(output)

	if err != nil {
		if strings.Contains(strings.ToLower(result), "1060") {
			fmt.Printf("%s|NOT_FOUND|%d\n", serviceName, exitNotFound)
			return exitNotFound
		}

		fmt.Printf("%s|ERROR|%d\n", serviceName, exitError)
		return exitError
	}

	switch {
	case strings.Contains(result, "RUNNING"):
		fmt.Printf("%s|RUNNING|%d\n", serviceName, exitRunning)
		return exitRunning

	case strings.Contains(result, "STOPPED"):
		fmt.Printf("%s|STOPPED|%d\n", serviceName, exitStopped)
		return exitStopped

	default:
		fmt.Printf("%s|UNKNOWN|%d\n", serviceName, exitError)
		return exitError
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: windows-service-checker <service-name> [service-name...]")
		os.Exit(exitError)
	}

	fmt.Println("SERVICE|STATUS|EXIT_CODE")

	overallExit := exitRunning

	for _, serviceName := range os.Args[1:] {
		exitCode := checkService(serviceName)

		if exitCode > overallExit {
			overallExit = exitCode
		}
	}

	os.Exit(overallExit)
}
