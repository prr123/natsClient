// natsCleintpub

package main

import (
	"fmt"
	"log"
    "os"
//	"time"
//	"bufio"
//  "strings"
//	"strconv"
	"context"
//	"errors"

    util "github.com/prr123/utility/utilLib"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var dbg bool

func main() {

    numarg := len(os.Args)
    flags:=[]string{"dbg", "serv", "port", "stream"}

    useStr := "/stream=<stream> [/serv=<server ip>] [/port=<4222>] [/dbg]"
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

    streamNam := "all"
    stval, ok := flagMap["stream"]
    if ok {
        if stval.(string) == "none" {log.Fatalf("error -- no topic provided!\n")}
        streamNam = stval.(string)
    }

    if dbg {
        fmt.Printf("serverip: %s\n", ipDest)
        fmt.Printf("port:     %s\n", ipDestPort)
        fmt.Printf("stream:   %s\n", streamNam)
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

//	hasJet, apiLvl := nc.ConnectedServerJetStream()
	hasJet, _ := nc.ConnectedServerJetStream()

	if !hasJet {log.Fatalf("error -- jetstream not enabled!")}

	log.Printf("info -- JetStream enabled: %v\n", hasJet)

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("error -- cannot get jet stream: %v", err)
	}

	ctx := context.Background()
	if streamNam == "all" {
		// list stream names
		names := js.StreamNames(ctx)
		for name := range names.Name() {
			fmt.Printf("stream name: %s\n", name)
		}
		if names.Err() != nil {
    		log.Fatalf("error -- obtain js stream names: %v", names.Err())
		}
	} else {
		log.Printf("name check to come!\n")
	}

}


