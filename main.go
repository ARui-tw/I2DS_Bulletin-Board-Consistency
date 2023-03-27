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

	log "github.com/sirupsen/logrus"
)

var (
	config_file = flag.String("config_file", "server_config.json", "The server config file")
)

type Config struct {
	Type    string `json:"type"`
	Primary int    `json:"primary"`
	Child   []int  `json:"child"`
	NR      int    `json:"NR"`
	NW      int    `json:"NW"`
}

func main() {
	flag.Parse()

	jsonFile, err := os.Open(*config_file)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened json config")
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
			cmd := exec.Command("go", "run", "BBC_child_server/main.go", "-addr", fmt.Sprintf("localhost:%d", config.Primary), "-port", strconv.Itoa(child))
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
	} else if config.Type == "quorum" {
		if config.NR+config.NW <= len(config.Child)+1 || config.NW <= (len(config.Child)+1)/2 {
			log.Fatal("Invalid config")
		}

		// start coordinator server
		cmd := exec.Command("go", "run", "BBC_coordinator_server/main.go", "-config_file", *config_file, "-port", strconv.Itoa(config.Primary))
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
			cmd := exec.Command("go", "run", "BBC_quorum_server/main.go", "-addr", fmt.Sprintf("localhost:%d", config.Primary), "-port", strconv.Itoa(child))
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
	} else {
		log.Fatal("Invalid config type")
	}

	log.Info("All servers started")

	// sleep forever
	<-make(chan int)
}
