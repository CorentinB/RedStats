package main

import (
	"os"
	"strconv"

	"github.com/labstack/gommon/color"
)

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
