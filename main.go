package main

import (
	"flag"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"log"
	"os"
	"time"
)

var (
	debugLogger = log.New(os.Stdout, "[DEBUG  ] ", log.LstdFlags)
	infoLogger  = log.New(os.Stdout, "[INFO   ] ", log.LstdFlags)
	errorLogger = log.New(os.Stderr, "[ERROR  ] ", log.LstdFlags)
)

const (
	chipName = "gpiochip0"
)

var (
	delayBeforeNotification = flag.Duration("notification-delay", 5*time.Second, "delay between the moment the light is detected and the notification is sent")
	lightSensorPin          = flag.Int("light-sensor-pin", rpi.GPIO23, "GPIO pin number on which the light sensor is plugged")
	buzzerPin               = flag.Int("buzzer-pin", rpi.GPIO24, "GPIO pin number on which the buzzer is plugged")
)

func main() {
	flag.Parse()

	c, err := gpiod.NewChip(chipName)
	if err != nil {
		errorLogger.Fatalf("failed to bind chip %q: %v", chipName, err)
	}

	n := NewNotifier(c, *buzzerPin)

	lightSensor, err := c.RequestLine(*lightSensorPin)
	if err != nil {
		errorLogger.Fatalf("failed to request line %d: %v", lightSensorPin, err)
	}

	for {
		v, err := lightSensor.Value()
		lightOn := v == 0
		if err != nil {
			errorLogger.Printf("failed to retrieve light sensor value: %v", err)
		} else {
			if lightOn {
				debugLogger.Println("\u263c")
				if !n.isNotifying {
					n.Prepare()
					go n.fireNotificationIn(*delayBeforeNotification)
				}
			} else {
				debugLogger.Println("\u263e")
				n.FireCancellation()
			}
		}
		time.Sleep(1 * time.Second)
	}
}
