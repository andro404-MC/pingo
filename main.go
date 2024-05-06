package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/prometheus-community/pro-bing"
	flag "github.com/spf13/pflag"
)

var (
	projectName string = "Pingo"

	firstTry bool = true

	isConnected bool = false
	trail       bool = false
	inPolybar   bool = false

	reCheckingDelay uint8
	reTryingDelay   uint8
	timeout         uint8

	red      string
	green    string
	blue     string
	none     string
	pingIP   string
	noneStop bool

	mode string = "term"
	// term     : full terminal text
	// termMin  : (default) minimal terminal text
	// ico      : icon mode
	// notify   : notification mode
)

const (
	attempts  uint8  = 3
	IPdefault string = "google.com"

	redB   string = "\033[38;5;9m"
	greenB string = "\033[38;5;10m"
	blueB  string = "\033[38;5;33m"
	noneB  string = "\033[0m"

	redP   string = "%{F#c55555}"
	greenP string = "%{F#a9c474}"
	blueP  string = "%{F#0b79fe}"
	noneP  string = "%{F-}"
)

func main() {
	args()
	env()
	colorSetup()
	leftAttempts := attempts

	for leftAttempts != 0 {
		// try
		isConnected = pinger()
		if isConnected {
			if firstTry && !noneStop {
				break
			}
			// Connected but need recheck
			if !noneStop {
				recheckMSG()
				leftAttempts--
			} else {
				conectedMSG()
			}

			time.Sleep(time.Duration(reCheckingDelay) * time.Second)
		} else {
			// not Connected
			notconectedMSG()
			leftAttempts = attempts
			time.Sleep(time.Duration(reTryingDelay) * time.Second)
		}
		firstTry = false
	}
	// all GOOD :)
	stableMSG()
}

func sendNotification(amount, duration int, text string, isCritical bool) {
	cmd := []string{projectName, text, fmt.Sprintf("--expire-time=%d", duration)}
	if isCritical {
		cmd = append(cmd, "--urgency=critical")
	}

	for i := 1; i <= amount; i++ {
		exec.Command("notify-send", cmd...).Run()
	}
}

func pinger() bool {
	resultChannel := make(chan bool)

	go func() {
		defer close(resultChannel)

		pinger, err := probing.NewPinger("google.com")
		if err != nil {
			resultChannel <- false
			return
		}

		pinger.Count = 1
		err = pinger.Run()
		if err != nil {
			resultChannel <- false
			return
		}

		resultChannel <- true
	}()

	select {
	case result := <-resultChannel:
		return result
	case <-time.After(time.Millisecond * time.Duration(timeout)):
		return false
	}
}

func stableMSG() {
	switch mode {
	case "termMin":
		RLL()
		fmt.Printf(" %s✓%s - Conection is stable \n", green, none)
	case "ico":
		RLL()
		fmt.Printf(" %s●%s \n", green, none)
	case "notify":
		sendNotification(1, 5000, "Conection is stable", false)
	}
}

func recheckMSG() {
	switch mode {
	case "termMin":
		RLL()
		fmt.Printf(" %s◎%s - Connected re-checking ", blue, none)
	case "ico":
		RLL()
		fmt.Printf(" %s◎%s ", blue, none)
	case "notify":
		sendNotification(1, 5000, "Connected re-checking", false)
	}
}

func notconectedMSG() {
	switch mode {
	case "termMin":
		RLL()
		fmt.Printf(" %s●%s - Not connected ", red, none)
	case "ico":
		RLL()
		fmt.Printf(" %s●%s ", red, none)
	case "notify":
		sendNotification(1, 1000, "Not connected", true)
	}
}

func conectedMSG() {
	switch mode {
	case "termMin":
		RLL()
		fmt.Printf(" %s●%s - Connected ", green, none)
	case "ico":
		RLL()
		fmt.Printf(" %s●%s ", green, none)
	case "notify":
		sendNotification(1, 1000, "Not connected", true)
	}
}

func RLL() {
	// I KNOW THIS IS SO FUCKING DUMB
	if trail {
		fmt.Printf("\r                       ")
		fmt.Printf("\r")
	} else if !firstTry {
		fmt.Print("\n")
	}
}

func args() {
	flag.StringVarP(&mode, "mode", "m", "termMin", "Mode (term, termMin, ico, notify)")
	flag.BoolVarP(&noneStop, "nonestop", "n", false,
		"turn on noneStop (use '$ killall pingo' for stop)")
	flag.BoolVarP(&trail, "trail", "t", false, "trail (replace last line)")
	flag.BoolVarP(&inPolybar, "polybar", "p", false, "polybar colors")
	flag.Uint8Var(&timeout, "timeout", 200, "ping timeout in Milliseconds")
	flag.Uint8Var(&reCheckingDelay, "recheck-delay", 8, "delay between rechecks in Seconds")
	flag.Uint8Var(&reTryingDelay, "retry-delay", 1, "delay between retrys in Seconds")
	flag.Parse()
}

func env() {
	pingIP = os.Getenv("PINGOIP")
	if pingIP == "" {
		pingIP = IPdefault
	}
}

func colorSetup() {
	if inPolybar {
		red, blue, green, none = redP, blueP, greenP, noneP
	} else {
		red, blue, green, none = redB, blueB, greenB, noneB
	}
}
