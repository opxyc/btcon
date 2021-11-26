package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/opxyc/goutils/sch"
)

type List []string

func (l *List) Contains(item string) bool {
	for _, i := range *l {
		if i == item {
			return true
		}
	}
	return false
}

var devicesList = List{
	"XX:XX:XX:XX:XX:XX", "FC:58:FA:4E:46:F0",
}

func main() {
	connected := false
	conCount := 0
	log.SetFlags(log.Llongfile)

	prevOpSink := ""
	ch := make(chan struct{})
	ch2 := make(chan struct{})
	sch.PingAfter(context.Background(), time.Duration(time.Second*5), ch)
	sch.PingAfter(context.Background(), time.Duration(time.Second*5), ch2)
	go func() {
		for {
			<-ch2

			cmd := exec.Command("pacmd", "stat")
			op, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("%v", err)
			}
			opStr := string(op)
			s := strings.Split(opStr, "\n")
			for _, d := range s {
				if strings.Contains(d, "Default sink name:") {
					if strings.Compare(prevOpSink, d) != 0 {
						fmt.Println("connected set to false")
						connected = false
					} else if conCount > 1 {
						fmt.Println("connected set to true")
						connected = true
						conCount = 0
					}
					prevOpSink = d
				}
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			<-ch
			if connected {
				continue
			}

			cmd := exec.Command("bluetoothctl", "devices")
			op, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("%v", err)
			}
			opStr := string(op)
			s := strings.Split(opStr, "Device ")
			for _, d := range s {
				d = strings.Trim(d, "\n")
				s := strings.Split(d, " ")
				if len(s) < 2 {
					continue
				}
				mac := s[0]

				if devicesList.Contains(mac) {
					connectCmd := exec.Command("bluetoothctl", "connect", mac)
					op, _ := connectCmd.CombinedOutput()
					fmt.Print(string(op))
					conCount += 1
				}
			}
		}
	}()

	wg.Wait()
}
