package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/soarpenguin/smsmail/g"
	"github.com/soarpenguin/smsmail/http"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	handleVersion(*version)
	handleHelp(*help)

	// global config
	g.ParseConfig(*cfg)

	// http
	http.Start()

	select {}
}

func handleVersion(displayVersion bool) {
	if displayVersion {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
}

func handleHelp(displayHelp bool) {
	if displayHelp {
		flag.Usage()
		os.Exit(0)
	}
}
