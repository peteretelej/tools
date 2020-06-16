package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	var (
		interval = flag.Duration("interval", time.Second*5, "Interval between active connections checks")
	)
	flag.Parse()
	log.SetPrefix("[active_conn] ")
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case _ = <-c:
			log.Printf("Shutting down..")
			if err := cleanup(); err != nil {
				log.Fatalf("cleanup failed: %v", err)
			}
			os.Exit(0)
		}
	}()
	if err := dbconn(); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	if err := checks(*interval); err != nil {
		log.Printf("checks failed: %v", err)
		if err := cleanup(); err != nil {
			log.Fatal(err)
		}
	}
}

func dbconn() error {
	for _, env := range []string{"DB_HOST", "DB_USER", "DB_PASS"} {
		if v := os.Getenv(env); v == "" {
			return fmt.Errorf("Missing environment variable: %s (please see .env.sample)", env)
		}
	}
	var err error
	db, err = sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
			os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"), os.Getenv("DB_NAME")),
	)
	if err != nil {
		return err
	}
	return db.Ping()
}

func cleanup() error {
	if err := db.Close(); err != nil {
		log.Printf("Failed to close db conn: %v", err)
	}
	return nil
}

func countConns() (int, error) {
	var count int
	var a string
	err := db.QueryRow("show status where `variable_name` = 'Threads_connected' ").Scan(&a, &count)
	if err != nil {
		return 0, err
	}
	return count, nil

}

func checks(interval time.Duration) error {
	if interval < time.Second {
		return fmt.Errorf("interval can't be less than 1 second")
	}
	for {
		c, err := countConns()
		if err != nil {
			return err
		}
		log.Printf("DB connections: %d", c)
		time.Sleep(interval)
	}
}
