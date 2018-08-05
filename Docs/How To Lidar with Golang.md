# Connecting and reading Garmin LidarLite Sensor with Gobot Framework

## GOBOT Framework

GOBOT is a framework written in Go programming language. Useful for connecting robotic components, vareity of hardware & IoT devices. 

Framework consists of 

* Robots -> Virtual entity representing rover, drones, sensors etc.
* Adaptors -> Allows connectivity to the hardware e.g. connection to Arduino is done using Firmata Adaptor, defining how to talk to it
* Drivers -> Defines specific functionality to support on specific hardware devices e.g. buttons, sensors, etc.
* API -> Provides RESTful API to query Robot status
  
There are additional core features of the framework that I recommend having a look esp. Events, Commands allowing Subscribing / Publishing events to the device for more refer to the [doc](https://gobot.io/documentation/)

There's already a long list of [Platforms](https://gobot.io/documentation/platforms/) for which the drivers and adaptors are available.
For this blog I will be working with Arduino + [Garmin LidarLite v3](https://buy.garmin.com/en-US/US/p/557294). There are cheaper versions available for distance measturement, however if you are looking for high performance, high precision optical distance measurement sensor, then this is it.

# Pre-requisite
* Install Go see [doc](https://golang.org/doc/install)
* [Install Gort](http://gort.io/)
* [Install Gobot](https://gobot.io/documentation/platforms/arduino/)
* [Wire-up LidarLite Sensor with Arduino](https://learn.sparkfun.com/tutorials/lidar-lite-v3-hookup-guide?_ga=2.113778690.440097749.1533500292-605393167.1533500292)

# How to connect
For our current setup I have Arduino connected to Ubuntu 16 over serial port, see [here](https://gobot.io/documentation/platforms/arduino/) if you are looking for a different platform.

For ubuntu you just need following 3 commands to connect and upload the firmata as our Adaptor to prepare Arduino for connectivity

```bash
// Look for the connected serial devices
$ gort scan serial

// install avrdude to upload firmata to the Arduino
$ gort arduino install

// uploading the firmata to the serial port found via first scan command, mine was found at /dev/ttyACM0
$ gort arduino upload firmata /dev/ttyACM0
```

## Reading Sensor data

Since there is a available driver for the LidarLite, I will be using it in the following Go code below in a file called main.go which connects and reads the sensor data.


For connecting and reading the sensor data its we need the driver, connection object & the taks / work that the robot is supposed to perform. 

### Adaptor

```Go
firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0") // this the port on which for me Arduino is connecting
```

### Driver

As previously mentioned that Gobot provides several drivers on of the them is LidarLite[LidarLite](https://gobot.io/documentation/drivers/lidar-lite/) we will be using this like so

```Go
d := i2c.NewLIDARLiteDriver(firmataAdaptor)
```

### Work

Now that we have the adaptor & the driver setup lets assign the work this robot needs to do, which is to read the distance

```Go
work := func() {
		gobot.Every(1*time.Second, func() {
			dist, err := d.Distance()

			if err != nil {
				log.Fatalln("failed to get dist")
			}
			fmt.Println("Fetching the dist", dist, "cms")
		})
	}
```

Notice the **Every** function provided by gobot to define that we want to perform certain action as the time lapses, here we are gathering the distance.

#### Note that the distance returned by the lidarLite sensor is in CMs & the max range for the sensor is 40m

#### Robot

Now we create the robot representing our entity which in this case is simple, its just the sensor itself

```Go
    lidarRobot := gobot.NewRobot("lidarBot", 
    []gobot.Connection{firmataAdaptor},
	[]gobot.Device{d}, work)
```

This defines the vitual representation of the entity and the driver + the work this robot needs to do.

Here's the complete code. Before running this pacakge make sure to build it as you likely will have to execute the runnable with sudo. To build simply navigate to the folder in the shell where the main.go exists and execute 

```bash
$ go build
```

This will create runnable file with the package name execute the same with sudo if needed like so 

```bash
$ sudo ./GarminLidarLite
```

And if everything done as required following ouput will appear with sensor readings printed out every second

```output
2018/08/05 22:46:54 Initializing connections...
2018/08/05 22:46:54 Initializing connection Firmata-634725A2E59CBD50 ...
2018/08/05 22:46:54 Initializing devices...
2018/08/05 22:46:54 Initializing device LIDARLite-5D4F0034ECE4D0EB ...
2018/08/05 22:46:54 Robot lidarBot initialized.
2018/08/05 22:46:54 Starting Robot lidarBot ...
2018/08/05 22:46:54 Starting connections...
2018/08/05 22:46:54 Starting connection Firmata-634725A2E59CBD50 on port /dev/ttyACM0...
2018/08/05 22:46:58 Starting devices...
2018/08/05 22:46:58 Starting device LIDARLite-5D4F0034ECE4D0EB...
2018/08/05 22:46:58 Starting work...
Fetching the dist 166 cms
Fetching the dist 165 cms
Fetching the dist 165 cms
```

Here's complete code for reference.

```Go
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

    lidarRobot := gobot.NewRobot("lidarBot", 
    []gobot.Connection{firmataAdaptor},
	[]gobot.Device{d}, work)

	lidarRobot.Start()

}

```

