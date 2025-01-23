package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
)

// Define a map of color names to ANSI color codes
var colorCodes = map[string]string{
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"gray":    "\033[90m",
	"orange":  "\033[38;5;214m",
	"pink":    "\033[38;5;207m",
	"lime":    "\033[38;5;10m",
	"white":   "\033[37m", // Default color
}

// Function to colorize a string
func colorize(text, color string) string {
	// Get the ANSI code for the color; default to white if not found
	ansiCode, exists := colorCodes[color]
	if !exists {
		ansiCode = colorCodes["white"]
	}
	// Return the colored string with a reset at the end
	return fmt.Sprintf("%s%s\033[0m", ansiCode, text)
}

type Command struct {
	Label    string
	Color    string
	CmdStr   string
	ExecPath string
}

type ConfigFile struct {
	Commands []Command
}

func createConfigFile(config ConfigFile) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath := filepath.Join(workingDir, "config.xrun.json")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}

func readConfigFile() (ConfigFile, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return ConfigFile{}, err
	}

	filePath := filepath.Join(workingDir, "config.xrun.json")
	file, err := os.Open(filePath)
	if err != nil {
		return ConfigFile{}, err
	}
	defer file.Close()

	var config ConfigFile
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return ConfigFile{}, err
	}

	return config, nil
}

func runCommand(cmdStruct Command, wg *sync.WaitGroup) {
	defer wg.Done()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", fmt.Sprintf("cd %s && %s", cmdStruct.ExecPath, cmdStruct.CmdStr))
	} else {
		cmd = exec.Command("sh", "-c", fmt.Sprintf("cd %s && %s", cmdStruct.ExecPath, cmdStruct.CmdStr))
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}
	fmt.Println(colorize(fmt.Sprintf("Running command: %s", cmdStruct.Label), cmdStruct.Color))

	outputScanner := bufio.NewScanner(stdout)
	errorScanner := bufio.NewScanner(stderr)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for outputScanner.Scan() {
			fmt.Printf("%s %s\n", colorize(fmt.Sprintf("[%s]", cmdStruct.Label), cmdStruct.Color), outputScanner.Text())
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for errorScanner.Scan() {
			fmt.Printf("%s %s\n", colorize(fmt.Sprintf("[%s]", cmdStruct.Label), cmdStruct.Color), errorScanner.Text())
		}
	}()

}

func showCommandsInfo(config ConfigFile) {
	fmt.Println(colorize("Configured Commands:", "yellow"))
	for _, cmd := range config.Commands {
		fmt.Printf("%s  %s\n", colorize(cmd.Label+" -> ", cmd.Color), cmd.CmdStr)
	}
	fmt.Println(colorize("========================", "yellow"))
}

func showArtBanner() {
	banner := `
▄   ▄  ▄▄▄ █  ▐▌▄▄▄▄  
 ▀▄▀  █    ▀▄▄▞▘█   █ 
▄▀ ▀▄ █         █   █ 
                        
`
	fmt.Println(colorize(banner, "cyan"))
}
func main() {
	isInit := flag.Bool("init", false, "Create a new xrun configuration")
	flag.Parse()

	if *isInit {
		config := ConfigFile{
			Commands: []Command{
				{
					Label:    "echo",
					Color:    "green",
					CmdStr:   "echo 'Hello, World!'",
					ExecPath: ".",
				},
			},
		}
		createConfigFile(config)
		fmt.Println("Creating a new xrun configuration")
	} else {
		config, err := readConfigFile()
		if err != nil {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}

		showArtBanner()
		showCommandsInfo(config)

		var wg sync.WaitGroup

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

		for _, cmd := range config.Commands {
			wg.Add(1)
			go runCommand(cmd, &wg)
		}

		go func() {
			<-signalChan
			fmt.Println("\nReceived interrupt. Stopping all commands...")
			wg.Wait() // Ensure all goroutines finish
			os.Exit(0)
		}()

		wg.Wait()

	}

}
