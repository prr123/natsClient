package main

import (
	"fmt"
	"log"
    "os"
//	"time"
	"sync"
//	"bufio"
//    "strings"
//	"strconv"
//	"context"

    util "github.com/prr123/utility/utilLib"

	"github.com/nats-io/nats.go"
//	"github.com/nats-io/nats.go/jetstream"
)

var dbg bool

func main() {

    numarg := len(os.Args)
    flags:=[]string{"dbg", "serv", "port", "topic"}

    useStr := "/serv=<server ip> [/port=<4222>] /topic=<topic> [/dbg]"
    helpStr := "nats client sub"

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
    if !ok {
        log.Fatalf("error -- no ip flag provided!\n")
    } else {
        if ipval.(string) == "none" {log.Fatalf("error -- no server ip destination provided!\n")}
        ipDest = ipval.(string)
    }

	ipDestPort := "4222"
    pval, ok := flagMap["port"]
    if ok {
        if pval.(string) == "none" {log.Fatalf("error -- no port destination provided!\n")}
        ipDestPort = pval.(string)
    }

    subTopic := ""
    tval, ok := flagMap["topic"]
    if !ok {
        log.Fatalf("error -- no topic flag provided!\n")
    } else {
        if tval.(string) == "none" {log.Fatalf("error -- no topic provided!\n")}
        subTopic = tval.(string)
    }

    if dbg {
        fmt.Printf("dbg -- server: %s\n", ipDest)
        fmt.Printf("dbg -- port:   %s\n", ipDestPort)
        fmt.Printf("dbg -- topic:  %s\n", subTopic)
    }

	dest := ipDest + ":" + ipDestPort
//	nc, err := nats.Connect(nats.DefaultURL)
	nc, err := nats.Connect(dest)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Drain()

	log.Printf("connection established! Listening on Topic: %s\n", subTopic)

	var wg sync.WaitGroup
	wg.Add(1) // We want to wait for 1 callback cycle to hit


	// 1. Asynchronous Subscriber
	_, err = nc.Subscribe(subTopic, func(m *nats.Msg) {
		log.Printf("Received alert: %s\n", string(m.Data))
		if string(m.Data) == "exit" {
			fmt.Println("Goodbye!")
			wg.Done()
		}
	})

	if err != nil {log.Fatalf("error -- subscription failed: %v\n", err)}

	wg.Wait()

	log.Println("Initiating graceful drain shutdown...")
}
