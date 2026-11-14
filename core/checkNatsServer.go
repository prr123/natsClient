// checkNatsServer


package main

import (
    "fmt"
    "log"
    "net"
    "os"
    "errors"
//  "bufio"
//  "strings"

    "strconv"
    "time"

    util "github.com/prr123/utility/utilLib"

)


var dbg bool


func main() {


    numArgs := len(os.Args)

    flags:=[]string{"dbg","port", "serv", "tmout"}

    useStr := "/serv=adr /port=<portnum> [/tmout=sec] [/dbg]"
    helpStr := fmt.Sprintf("help: tcp Client\n")

    if numArgs > len(flags)+1 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: %s %s\n", os.Args[0], useStr)
        os.Exit(1)
    }

   if numArgs == 1 {
        fmt.Printf("usage is: %s\n", useStr)
        fmt.Printf("%s\n", helpStr)
        os.Exit(1)
    }

   if numArgs == 2 {
        if os.Args[1] == "help" {
            fmt.Printf("usage is: %s\n", useStr)
            fmt.Printf("%s\n", helpStr)
            os.Exit(1)
        }
    }


    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}

    dbg = false
    _, ok := flagMap["dbg"]
    if ok {dbg = true}

    sadrStr := ""
    sval, ok := flagMap["serv"]
    if !ok {
        log.Fatalf("error -- no /serv flag in cli!")
    } else {
        if sval.(string) == "none" {log.Fatalf("error -- no serve IP provided with /serve flag!")}
        sadrStr = sval.(string)
    }

    portStr := ""
    pval, ok := flagMap["port"]
    if !ok {
        log.Fatalf("error -- no /port flag in cli!")
    } else {
        if pval.(string) == "none" {log.Fatalf("error -- no port number provided with /port flag!")}
        portStr = pval.(string)
    }

    tmoutStr := "1"
    tval, ok := flagMap["tmout"]
    if ok {
        if tval.(string) == "none" {log.Fatalf("error -- no tmout value provided with /stmout flag!")}
        tmoutStr = tval.(string)
    }

    tmoutSec, err :=strconv.Atoi(tmoutStr)
    if err != nil {log.Fatalf("error-- timeout is not an integer: %v", err)}

    tmout := time.Duration(tmoutSec)*time.Second

    servAdr := sadrStr + ":" + portStr

    if dbg {
        fmt.Printf(" serv IP:  %s\n", sadrStr)
        fmt.Printf(" port:     %s\n", portStr)
        fmt.Printf(" serv Adr: %s\n", servAdr)
        fmt.Printf(" time out: %v\n", tmout)
    }

    sResp := make([]byte, 256, 1024)

    conn, err := net.Dial("tcp", servAdr)
    if err != nil {
        log.Fatalf("error -- Connection failed to connect: %v", err)
    }
    defer conn.Close()

    fmt.Printf("Connection with %s established!\n", servAdr)

    //check for server response
    n, err := conn.Read(sResp)
    if err != nil {
        var netErr net.Error
        if errors.As(err, &netErr) && netErr.Timeout() {
            // The connection succeeded, but the server sent no data before the timeout.
            // For protocols like HTTP, this is normal behavior.
            fmt.Printf(" no server response\n")
        }
        // Any other error (e.g., EOF if the server instantly disconnects)
        log.Fatalf("error --- first conn read: %v", err)
    } else {
        // Server sent data immediately after connection
        fmt.Printf("Initial Server Response [%d]: %s\n", n, sResp[:n])
    }

}
