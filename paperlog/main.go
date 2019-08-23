package main

import (
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"os"
)

var (
	message   = flag.String("message", "", "The message to log")
	isError   = flag.Bool("error", false, "Log message as an ERROR")
	isWarning = flag.Bool("warn", false, "Log message as a WARNING")

	// optional
	addr    = flag.String("addr", "", "The Papertrail address")
	appname = flag.String("appname", "", "Your application name")
)

// Papertrail configs
var (
	Addr    = os.Getenv("PAPERLOG_ADDR")
	AppName = os.Getenv("PAPERLOG_APPNAME")
)

func main() {
	flag.Parse()
	log.SetPrefix("[Paperlog] ")
	if Addr == "" {
		Addr = *addr
	}
	if AppName == "" {
		AppName = *appname
	}
	if Addr == "" || AppName == "" {
		log.Fatal("Missing Papertrail configs from environment or arguments, please see Paperlog README")
	}
	if *message == "" {
		log.Fatal("No message specified")
	}
	if err := run(*message); err != nil {
		log.Fatalf("ERROR: failed to log: %v\n\t%s", err, *message)
	}
}

func run(message string) error {
	w, err := syslog.Dial("udp", Addr, syslog.LOG_EMERG|syslog.LOG_KERN, AppName)
	if err != nil {
		return fmt.Errorf("failed to dial syslog: %s", Addr)
	}

	switch {
	case *isError:
		err = w.Err(message)
	case *isWarning:
		err = w.Warning(message)
	default:
		err = w.Info(message)
	}

	if errC := w.Close(); errC != nil {
		log.Printf("failed to close syslog connection: %v", errC)
	}

	if err != nil {
		return fmt.Errorf("Failed to write message to Papertrail: %v", err)
	}

	return nil
}
