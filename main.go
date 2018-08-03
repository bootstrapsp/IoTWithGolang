package main

import (
	"fmt"
	"log"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {

	lidarLibTest()
}

// reading Garmin LidarLite data successfully
func lidarLibTest() {

	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")

	d := i2c.NewLIDARLiteDriver(firmataAdaptor)

	work := func() {
		gobot.Every(1*time.Second, func() {
			dist, err := d.Distance()

			if err != nil {
				log.Fatalln("failed to get dist")
			}
			fmt.Println("Fetching the dist", dist, "cms")
		})
	}

	lidarRobot := gobot.NewRobot("lidarBot", []gobot.Connection{firmataAdaptor},
		[]gobot.Device{d}, work)

	lidarRobot.Start()

}

// testing old lib for cur

func testMPL3115A2Lib() {

	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")

	d := i2c.NewMPL115A2Driver(firmataAdaptor, i2c.WithBus(0), i2c.WithAddress(0x60))

	work := func() {
		gobot.Every(1*time.Second, func() {
			currTemp, err := d.Temperature()
			if err != nil {
				log.Fatalln("failed to get temp")
			}
			fmt.Println("Fetching the temp", currTemp)
		})
	}

	pressureRobot := gobot.NewRobot("pressureBot", []gobot.Connection{firmataAdaptor},
		[]gobot.Device{d}, work)

	pressureRobot.Start()
	d.GetAddressOrDefault(0x00)
	fmt.Println(d)

}

// sample for testing Arduino connectivity by blinking LED

func ledBlinkerWithGPIO() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	led := gpio.NewLedDriver(firmataAdaptor, "13")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led},
		work,
	)

	robot.Start()

}
