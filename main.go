package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gosuri/uilive"
	"github.com/labstack/gommon/color"
)

var writer = uilive.New()

var config = struct {
	Host     string
	Password string
	DB       int64
	Key      string
	Client   *redis.Client
}{}

func parseArgs(args []string) {
	if len(args) > 4 {
		color.Red("Usage: ./RedStats [REDIS-HOST] [REDIS-PASSWORD] [REDIS-DB] [REDIS-KEY]")
		os.Exit(1)
	} else if len(args) == 4 {
		// Redis host
		config.Host = args[0]

		// Redis Password
		config.Password = args[1]

		// Redis DB
		if _, err := strconv.ParseInt(args[2], 10, 64); err == nil {
			config.DB, _ = strconv.ParseInt(args[2], 10, 64)
		} else {
			color.Red("Bad argument for DB: " + args[2])
			color.Red("Usage: ./RedStats [REDIS-HOST] [REDIS-PASSWORD] [REDIS-DB] [REDIS-KEY]")
			os.Exit(1)
		}

		// Redis Key
		config.Key = args[3]
	} else {
		color.Red("Only " + string(len(args)) + " specified.")
		color.Red("Usage: ./RedStats [REDIS-HOST] [REDIS-PASSWORD] [REDIS-DB] [REDIS-KEY]")
		os.Exit(1)
	}
}

func main() {
	parseArgs(os.Args[1:])

	// Connect to Redis server
	config.Client = redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Password: config.Password,
		DB:       int(config.DB),
	})

	// Check connection
	pong, err := config.Client.Ping().Result()
	if pong != "PONG" {
		fmt.Println("Unable to connect to Redis DB.")
		log.Fatal(err)
		os.Exit(1)
	}

	writer.Start()

	var nbIDsSecond, totalOps int64
	var opsSecond []int64
	var i = 0
	nbIDsSecond = 0
	for {
		nbIDsFirst, _ := config.Client.SCard(config.Key).Result()

		if i != 0 {
			opsSecond = append(opsSecond, nbIDsFirst-nbIDsSecond)
		} else {
			opsSecond = append(opsSecond, 0)
			i = 1
		}

		// Get current time
		timeNow := time.Now().Format(time.RFC850)

		// Average ops/s
		for _, ops := range opsSecond {
			totalOps += ops
		}
		opsPerSecond := totalOps / int64(len(opsSecond))

		fmt.Fprintln(writer, color.Green("[✔] [")+
			color.Yellow(timeNow)+
			color.Green("]\n[✔]")+
			color.Yellow(" -> ")+
			color.Green("Number of IDs: ")+
			color.Yellow(strconv.FormatInt(nbIDsFirst, 10))+
			color.Green("\n[✔]")+
			color.Yellow(" -> ")+
			color.Green("Average ops/s: ")+
			color.Yellow(strconv.FormatInt(opsPerSecond, 10)))

		nbIDsSecond, _ = config.Client.SCard(config.Key).Result()

		// Sleep 1 sc
		time.Sleep(1 * time.Second)
	}

	writer.Stop()
}
