package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	enqueue = `http://127.0.0.1:8080`
	count   = `http://127.0.0.1:8080/count`
	exit    = "q\n"
	user    = "user\n"
	robot   = "robot\n"
	counter = "reader\n"
)

var totalRPM int

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter text: \n")
		text, _ := reader.ReadString('\n')
		switch text {
		case exit:
			fmt.Println("exiting program...")
			os.Exit(1)
		case robot:
			rpm := rand.Intn(400) + 100
			totalRPM += rpm
			fmt.Printf("total RPS: %d\n", totalRPM/60)
			go createNewRobot(rpm)
		case user:
			rpm := rand.Intn(99)
			totalRPM += rpm
			fmt.Printf("total RPS: %d\n", totalRPM/60)
			go createNewRobot(rpm)
		case counter:
			rpm := rand.Intn(99)
			totalRPM += rpm
			fmt.Printf("total RPS: %d\n", totalRPM/60)
			go createNewReader(rpm)
		default:
		}
	}
}

func createNewRobot(rpm int) {
	delayms := int(float64(60) / float64(rpm) * 1000)
	userID := strconv.Itoa(rand.Int())
	req, err := http.NewRequest("GET", enqueue + fmt.Sprintf("/?user_id=%s", userID), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := http.Client{}
	for {
		_, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Duration(delayms) * time.Millisecond)
	}
}
func createNewReader(rpm int) {
	delayms := int(float64(60) / float64(rpm) * 1000)
	req, err := http.NewRequest("GET", count, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := http.Client{}
	for {
		_, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Duration(delayms) * time.Millisecond)
	}
}

