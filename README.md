# Nats Clients

Family of programs to explore the nats server.  

## Core

### CheckNatsServer

Progam that is a tcp client. It connects with the nats server and displays the response.  

### natsClientPub

Program that publishes a message on a topic on the nats server.  

### natsClientSub

Program that subscribes to a message on a topic on the nats server.  

### natsClientSubQueue

Program that subscribes to a message in queue group.

### natsClientResp

Program that subscribes to a message on a topic on the nats server. The client responds with a response message to the publisher after receiving a message.

### natsClientReq

Program that publishes a message and requests a reply.  




## Jet Stream

### natsClientCheckStream

program that check whether jet stream is enable in the nats server.  


### natsClientCheckStreamV2

modern version of previous program  

### natsClientStream

program that lists all streams of a jetstream enabled nats server.  


