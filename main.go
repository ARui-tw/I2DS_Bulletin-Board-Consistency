package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var (
	config_file = flag.String("config_file", "server_config.json", "The server config file")
)

type Config struct {
	Type    string `json:"type"`
	Primary int    `json:"primary"`
	Child   []int  `json:"child"`
}

func main() {
	flag.Parse()

	jsonFile, err := os.Open(*config_file)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config

	json.Unmarshal(byteValue, &config)

	if config.Type == "PBP" {
		// start primary server
		cmd := exec.Command("go", "run", "BBC_primary_server/main.go", "-config_file", *config_file, "-port", strconv.Itoa(config.Primary))
		cmdReader, _ := cmd.StderrPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
			return
		}

		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				fmt.Printf("\t > %s\n", scanner.Text())
			}
		}()

		cmd.Start()

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
			return
		}
		time.Sleep(1 * time.Second)

		// start child server
		for _, child := range config.Child {
			cmd := exec.Command("go", "run", "BBC_server/main.go", "-addr", fmt.Sprintf("localhost:%d", config.Primary), "-port", strconv.Itoa(child))
			cmdReader, _ := cmd.StderrPipe()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
				return
			}

			scanner := bufio.NewScanner(cmdReader)
			go func() {
				for scanner.Scan() {
					fmt.Printf("\t > %s\n", scanner.Text())
				}
			}()

			cmd.Start()

			if err != nil {
				fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
				return
			}

			time.Sleep(1 * time.Second)
		}
	}

	// sleep forever
	<-make(chan int)
}
