// natsClientCheckStreamV2
// programm that checks whether jetstream is enables
// v2 uses ConnectedSercerJetStream
//

package main

import (
	"fmt"
	"log"
    "os"
//	"time"
//	"bufio"
//  "strings"
//	"strconv"
//	"context"
//	"errors"

    util "github.com/prr123/utility/utilLib"

	"github.com/nats-io/nats.go"
//	"github.com/nats-io/nats.go/jetstream"
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

    if dbg {
        fmt.Printf("dbg -- serverip: %s\n", ipDest)
        fmt.Printf("dbg -- port:     %s\n", ipDestPort)
    }

	dest := ipDest + ":" + ipDestPort

	nc, err := nats.Connect(dest)
	if err != nil {
		log.Fatalf("error -- Failed to connect to NATS: %v", err)
	}
	defer nc.Drain()

	log.Printf("info -- nats connection established!\n")

	hasJet, apiLvl := nc.ConnectedServerJetStream()
	fmt.Printf("info -- has Jetstream:  %t\n", hasJet)
	fmt.Printf("info -- api level::     %d\n", apiLvl)

}


