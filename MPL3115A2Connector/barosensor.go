package mpl3115a2connector

import (
	"fmt"
	"log"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/firmata"
)

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
