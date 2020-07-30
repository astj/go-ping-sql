package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

func main() {
	var sslCaCert string
	flag.StringVar(&sslCaCert, "ssl-ca", "", "SSL CA certfile path for mysql. If you want use this CA, append ?tls=custom to DSN.")
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		log.Fatalln("usage: go-ping-sql [mysql|postgres] <DSN>")
	}

	driver := args[0]
	dsn := ""
	if len(args) > 1 {
		dsn = args[1]
	}

	if driver == "mysql" && sslCaCert != "" {
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(sslCaCert)
		if err != nil {
			log.Fatal(err)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			log.Fatal("Failed to append PEM.")
		}
		mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs: rootCertPool,
		})
	}

	// XXX how dynamic!!!
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Successed to ping", driver, dsn)
	os.Exit(0)
}
