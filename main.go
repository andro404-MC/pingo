package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/prometheus-community/pro-bing"
	flag "github.com/spf13/pflag"
)

var (
	projectName string = "Pingo"

	firstTry    bool = true
	isConnected bool = false

	mode string = "term"
	// term     : (default) full terminal text
	// termMin  : minimal terminal text
	// ico      : icon mode
	// notify   : notification mode

	noTrail   bool
	inPolybar bool
	noneStop  bool

	reCheckingDelay uint
	reTryingDelay   uint
	timeout         uint

	leftAttempts uint8

	red    string
	green  string
	blue   string
	none   string
	pingIP string
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
	leftAttempts = attempts

	for leftAttempts != 0 {
		// try
		isConnected, stats := pinger()
		if isConnected {
			if firstTry && !noneStop {
				break
			}
			// Connected but need recheck
			if !noneStop {
				recheckMSG(stats)
				leftAttempts--
			} else {
				conectedMSG(stats)
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

func pinger() (bool, *probing.Statistics) {
	pinger, err := probing.NewPinger(pingIP)
	if err != nil {
		return false, &probing.Statistics{}
	}

	pinger.Timeout = time.Millisecond * time.Duration(timeout)
	pinger.Count = 1

	err = pinger.Run()
	if err != nil {
		return false, &probing.Statistics{}
	}

	stats := pinger.Statistics()
	return stats.PacketsRecv >= 1, &*stats
}

func stableMSG() {
	switch mode {
	case "term":
		RLL()
		fmt.Printf(
			" %s✓%s - Conection is stable - %d/%d \n",
			green, none,
			attempts, attempts,
		)
		break
	case "termMin":
		RLL()
		fmt.Printf(" %s✓%s - Conection is stable \n", green, none)
		break
	case "ico":
		RLL()
		fmt.Printf(" %s●%s \n", green, none)
		break
	case "notify":
		sendNotification(1, 5000, "Conection is stable", false)
		break
	}
}

func recheckMSG(stats *probing.Statistics) {
	switch mode {
	case "term":
		RLL()
		fmt.Printf(
			" %s◎%s - Connected re-checking - %dms - %d/%d ",
			blue, none,
			int64(stats.AvgRtt.Milliseconds()),
			attempts-leftAttempts, attempts,
		)
		break
	case "termMin":
		RLL()
		fmt.Printf(" %s◎%s - Connected re-checking ", blue, none)
		break
	case "ico":
		RLL()
		fmt.Printf(" %s◎%s ", blue, none)
		break
	case "notify":
		sendNotification(1, 5000, "Connected re-checking", false)
		break
	}
}

func notconectedMSG() {
	switch mode {
	case "term":
		RLL()
		fmt.Printf(
			" %s●%s - Not connected - %sReset%s ",
			red, none, red, none,
		)
		break
	case "termMin":
		RLL()
		fmt.Printf(" %s●%s - Not connected ", red, none)
		break
	case "ico":
		RLL()
		fmt.Printf(" %s●%s ", red, none)
		break
	case "notify":
		sendNotification(1, 1000, "Not connected", true)
		break
	}
}

func conectedMSG(stats *probing.Statistics) {
	switch mode {
	case "term":
		RLL()
		fmt.Printf(
			" %s●%s - Connected - %dms ",
			green, none,
			int64(stats.AvgRtt.Milliseconds()),
		)
		break
	case "termMin":
		RLL()
		fmt.Printf(" %s●%s - Connected ", green, none)
		break
	case "ico":
		RLL()
		fmt.Printf(" %s●%s ", green, none)
		break
	case "notify":
		sendNotification(1, 1000, "Not connected", true)
		break
	}
}

func RLL() {
	// TODO better implementation
	// I KNOW THIS IS SO FUCKING DUMB but its working
	if !noTrail {
		fmt.Printf("\r" + strings.Repeat(" ", 50))
		fmt.Printf("\r")
	} else if !firstTry {
		fmt.Print("\n")
	}
}

func args() {
	flag.StringVarP(&mode, "mode", "m", "term", "Mode (term, termMin, ico, notify)")
	flag.BoolVarP(&noneStop, "nonestop", "n", false,
		"turn on noneStop (use '$ killall pingo' for stop)")
	flag.BoolVar(&noTrail, "no-trail", false, "no trail (no replacing last line)")
	flag.BoolVarP(&inPolybar, "polybar", "p", false, "polybar colors")
	flag.UintVar(&timeout, "timeout", 200, "ping timeout in Milliseconds")
	flag.UintVar(&reCheckingDelay, "recheck-delay", 8, "delay between rechecks in Seconds")
	flag.UintVar(&reTryingDelay, "retry-delay", 1, "delay between retrys in Seconds")
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
