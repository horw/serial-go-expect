package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"go.bug.st/serial"
)

func main() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Found the following ports:", ports)

	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open("/dev/ttyUSB0", mode)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	expect := make(chan string)
	go Monitor(port, expect, &wg)
	expect <- "Hello world"
	expect <- "Restarting in 7 seconds"
	close(expect)
	wg.Wait()
}

func Monitor(port serial.Port, expect chan string, wg *sync.WaitGroup) {
	buff := make([]byte, 100)
	for expect_str := range expect {

		expect_buff := make([]byte, 0, 10000)
		fmt.Printf("\nStarting to catch the string: %s \n", expect_str)
		wg.Add(1)
		for {
			n, err := port.Read(buff)
			if err != nil {
				log.Fatal(err)
				break
			}
			if n == 0 {
				fmt.Println("\nEOF")
				break
			}
			fmt.Printf("%v", string(buff[:n]))
			expect_buff = append(expect_buff, buff...)
			if strings.Contains(string(expect_buff), expect_str) {
				break
			}
		}
		wg.Done()
	}
	wg.Done()

}
