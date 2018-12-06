package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
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

	var nbIDsSecond, idsSecond int64
	nbIDsSecond = 0
	for {
		nbIDsFirst, _ := config.Client.SCard(config.Key).Result()
		infoStats, _ := config.Client.Info("stats").Result()
		infoServer, _ := config.Client.Info("server").Result()
		infoClients, _ := config.Client.Info("clients").Result()
		infoMemory, _ := config.Client.Info("memory").Result()
		//infoPersistence, _ := config.Client.Info("persistence").Result()

		stats := strings.Split(infoStats, "\n")
		server := strings.Split(infoServer, "\n")
		clients := strings.Split(infoClients, "\n")
		memory := strings.Split(infoMemory, "\n")
		//persistence := strings.Split(infoPersistence, "\n")
		//hddUsageRaw := strings.Split(persistence[17], ":")[1]
		//hddUsage, _ := strconv.ParseInt(hddUsageRaw, 10, 64)
		ramUsed := strings.Split(memory[2], ":")[1]
		opsSecond := strings.Split(stats[3], ":")[1]
		instantInput := strings.Split(stats[6], ":")[1]
		instantOutput := strings.Split(stats[7], ":")[1]
		uptimeRaw := strings.Split(server[14], ":")[1]
		uptime, _ := strconv.Atoi(uptimeRaw[:len(uptimeRaw)-1])
		connectedClients := strings.Split(clients[1], ":")[1]

		// Get current time
		timeNow := time.Now().Format(time.RFC850)

		idsSecond = nbIDsFirst - nbIDsSecond

		fmt.Fprintln(writer, color.Green("[✔] [")+
			color.Yellow(timeNow)+
			color.Green("]\n[✔]")+
			color.Yellow(" -> ")+
			color.Green("Uptime: ")+
			color.Yellow(secondsToHuman(uptime))+
			color.Green("\n[✔]")+
			color.Yellow(" -> ")+
			color.Green("Number of IDs: ")+
			color.Yellow(humanize.Comma(nbIDsFirst))+
			color.Green("\n[✔]")+
			color.Yellow(" -> ")+
			color.Green("Connected clients: ")+
			color.Yellow(connectedClients)+
			color.Green("\n[✔]")+
			color.Yellow(" -> ")+
			color.Green("IDs/s: ")+
			color.Yellow(strconv.FormatInt(idsSecond, 10))+
			color.Green("\n[✔]")+
			color.Yellow(" -> ")+
			color.Green("Ops/s: ")+
			color.Yellow(opsSecond)+
			color.Green("\n[✔]")+
			color.Yellow(" -> ")+
			color.Green("Instantaneous input: ")+
			color.Yellow(instantInput[:len(instantInput)-1])+
			color.Green(" kbps")+
			color.Green("\n[✔]")+
			color.Yellow(" -> ")+
			color.Green("Instantaneous output: ")+
			color.Yellow(instantOutput[:len(instantOutput)-1])+
			color.Green(" kbps")+
			color.Green("\n[✔]")+
			color.Yellow(" -> ")+
			color.Green("RAM used: ")+
			color.Yellow(ramUsed))
		//color.Green("\n[✔]")+
		//color.Yellow(" -> ")+
		//color.Green("DB size: ")+
		//color.Yellow(hddUsage))

		nbIDsSecond, _ = config.Client.SCard(config.Key).Result()

		// Sleep 1 sc
		time.Sleep(1 * time.Second)
	}
}
