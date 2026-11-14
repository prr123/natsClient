// natsCleintpub

package main

import (
	"fmt"
	"log"
    "os"
	"time"
//	"bufio"
//  "strings"
//	"strconv"
	"context"
	"errors"

    util "github.com/prr123/utility/utilLib"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var dbg bool

func main() {

    numarg := len(os.Args)
    flags:=[]string{"dbg", "serv", "port"}

    useStr := "/serv=<server ip> [/port=<4222>] [/dbg]"
    helpStr := "nats server check for jetstream"

    if numarg > len(flags) +1 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: %s %s\n", os.Args[0], useStr)
        os.Exit(-1)
    }

    if numarg == 1 || (numarg > 1 && os.Args[1] == "help") {
        fmt.Printf("help: %s\n", helpStr)
        fmt.Printf("usage is: %s %s\n", os.Args[0], useStr)
        os.Exit(1)
    }

    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}

    dbg = false
    _, ok := flagMap["dbg"]
    if ok {dbg = true}

    ipDest := "173.249.39.249"
    ipval, ok := flagMap["serv"]
    if ok {
//        log.Fatalf("error -- no server flag provided!\n")
//    } else {
        if ipval.(string) == "none" {log.Fatalf("error -- no server ip destination provided!\n")}
        ipDest = ipval.(string)
    }

	ipDestPort := "4222"
    pval, ok := flagMap["port"]
    if ok {
        if pval.(string) == "none" {log.Fatalf("error -- no port destination provided!\n")}
        ipDestPort = pval.(string)
    }
/*
    msg := ""
    tval, ok := flagMap["msg"]
    if !ok {
        log.Fatalf("error -- no msg flag provided!\n")
    } else {
        if tval.(string) == "none" {log.Fatalf("error -- no msg name provided!\n")}
        msg = tval.(string)
    }

    subTopic := ""
    stval, ok := flagMap["topic"]
    if !ok {
        log.Fatalf("error -- no topic flag provided!\n")
    } else {
        if stval.(string) == "none" {log.Fatalf("error -- no topic provided!\n")}
        subTopic = stval.(string)
    }
*/
    if dbg {
        fmt.Printf("serverip: %s\n", ipDest)
        fmt.Printf("port:     %s\n", ipDestPort)
//        fmt.Printf("topic:    %s\n", subTopic)
//        fmt.Printf("msg:      %s\n", msg)
    }

	dest := ipDest + ":" + ipDestPort
//	nc, err := nats.Connect(nats.DefaultURL)
	nc, err := nats.Connect(dest)
	if err != nil {
		log.Fatalf("error -- Failed to connect to NATS: %v", err)
	}
	defer nc.Drain()

	log.Printf("info -- nats connection established!\n")

	enabled, err := isJetStreamEnabled(nc)
	if err != nil {
		log.Fatalf("error -- checking JetStream status: %v", err)
	}

	log.Printf("info -- JetStream enabled: %v\n", enabled)
}


func isJetStreamEnabled(nc *nats.Conn) (bool, error) {

	js, err := jetstream.New(nc)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// AccountInfo queries the server and returns specific errors if disabled
	_, err = js.AccountInfo(ctx)
	if err != nil {
		if errors.Is(err, jetstream.ErrJetStreamNotEnabled) ||
			errors.Is(err, jetstream.ErrJetStreamNotEnabledForAccount) ||
			errors.Is(err, nats.ErrNoResponders) {
			return false, nil
		}
		// Return other errors (e.g., timeouts or permission issues)
		return false, err
	}

	return true, nil
}
