package main

import (
	"flag"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"log"
	"os"
	"time"
)

const (
	chipName                = "gpiochip0"
	lightSensorPin          = rpi.GPIO23
	buzzerPin               = rpi.GPIO24
	delayBeforeNotification = 5 * time.Second
)

var (
	debugLogger = log.New(os.Stdout, "[DEBUG  ] ", log.LstdFlags)
	infoLogger  = log.New(os.Stdout, "[INFO   ] ", log.LstdFlags)
	errorLogger = log.New(os.Stderr, "[ERROR  ] ", log.LstdFlags)
)

func main() {
	flag.Parse()

	infoLogger.Printf("Binding chip %q...", chipName)
	c, err := gpiod.NewChip(chipName)
	if err != nil {
		errorLogger.Fatalf("\nfailed to bind chip %q: %v", chipName, err)
	}
	defer c.Close()

	infoLogger.Printf("Requesting buzzer on line %q...", buzzerPin)
	buzzer, err := c.RequestLine(buzzerPin, gpiod.AsOutput(1))
	if err != nil {
		errorLogger.Fatalf("\nfailed to request line %q: %v", buzzerPin, err)
	}
	defer buzzer.Close()
	buzzer.SetValue(0)

	infoLogger.Printf("Requesting light sensor on line %q...", lightSensorPin)
	lightSensor, err := c.RequestLine(lightSensorPin)
	if err != nil {
		errorLogger.Fatalf("\nfailed to request line %q: %v", lightSensorPin, err)
	}
	defer lightSensor.Close()

	deny := make(chan struct{}, 1)
	notifying := false
	for {
		v, err := lightSensor.Value()
		lightOn := v == 0
		if err != nil {
			errorLogger.Printf("failed to retrieve light sensor value: %v", err)
		} else {
			if lightOn {
				debugLogger.Println("\u263c")
				if !notifying {
					emptyChannel(deny)
					infoLogger.Printf("Preparing to send notification...\n")
					notifying = true
					go fireNotificationIn(delayBeforeNotification, &deny, &notifying, buzzer)
				}
			} else {
				buzzer.SetValue(0)
				debugLogger.Println("\u263e")
				if len(deny) == 0 {
					deny <- struct{}{}
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func emptyChannel(ch chan struct{}) {
	for len(ch) > 0 {
		<-ch // emptying denial
	}
}

func fireNotificationIn(duration time.Duration, deny *chan struct{}, notifying *bool, buzzer *gpiod.Line) {
	select {
	case <-*deny:
		infoLogger.Printf("cancelling notification\n")
		*notifying = false
	case <-time.After(duration):
		infoLogger.Printf("notify\n")
		buzzer.SetValue(1)
		*notifying = false
	}
}
