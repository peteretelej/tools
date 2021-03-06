package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var (
		logFile = flag.String("logfile", "log.log", "log file to write")
		timeout = flag.Duration("timeout", time.Hour, "timeout after which daemon stops")
	)
	flag.Parse()

	f, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	t := time.Now()

	s := fmt.Sprintf("\n==== RESTARTED (%s) ====\n", time.Now().Format(time.RFC850))
	f.WriteString(s)

	for i := 1; ; i++ {
		time.Sleep(time.Second * 10)

		s := " ping "
		if i%10 == 0 {
			s += "\n"
		}

		if _, err := f.WriteString(s); err != nil {
			log.Fatal(err)
		}

		if time.Now().Sub(t) > *timeout {
			log.Printf("Timeout reached %s", *timeout)
			break
		}
	}
}
