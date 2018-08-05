package main

import (
	"fmt"
	"log"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {

	lidarLibTest()
}

// reading Garmin LidarLite data
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
