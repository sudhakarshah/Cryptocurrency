package main

import(
	"sync"
	"errors"
	"net"
	"encoding/json"
	"fmt"
	"time"
	"bufio"
)

type Node struct{
	Name string
	Ip string
	Port string
	LastActive int64
	Attempts int
	Sock net.Conn
	jsonEncoder *json.Encoder
	mux sync.Mutex
}

func (nd *Node)SendJson(m Msg)int{
	startTime := time.Now()
	b, err := fmt.Fprintf(nd.Sock, m.Data)
	if err != nil {
		// TODO: json failed to send
		nd.mux.Lock()
		if nd.Attempts > 2{
			nd.mux.Unlock()
			nd.Sock.Close()
			return -1
		}
		return 1
	}else{
		nd.mux.Lock()
		nd.Attempts = 0
		nd.mux.Unlock()
	}
	fmt.Printf("SEND %d %s %d %d\n",int64(time.Now().Unix()), m.GetType(), b,  time.Since(startTime)) // time, msg type, size, duration
	return 0
}


func (nd *Node)ListenToFriend(inbox *Box){
	reader := bufio.NewReader(nd.Sock)
	for {
		s, err := reader.ReadString('\n')
		//fmt.Printf("# Recieved %s from %s\n", s, nd.Name)
		time.Sleep(1 * time.Millisecond)
		if err != nil {
			// fmt.Printf("# Failed listening to node %s\n", nd.Name)
			// fmt.Printf("# ERROR: %s\n", err)
			return
		}
		if len(s) > 0{
			var m Msg
			m.Parse(s)
			inbox.enqueue(m)
		}
	}
}


type Box struct{
	messages []Msg
	mux sync.Mutex
}

func (in*Box) enqueue(m Msg){
	in.mux.Lock()
	in.messages = append(in.messages, m)
	in.mux.Unlock()
}

func (in*Box) pop()(Msg,error){
	var output Msg
	var err error
	in.mux.Lock()
	if len(in.messages) != 0{
		output = in.messages[0]
		in.messages = append(in.messages[:0], in.messages[1:]...)
	} else {
		err = errors.New("The inbox is empty")
		in.mux.Unlock()
		return output, err
	}
	in.mux.Unlock()
	return output, err
}

type Log struct{
	Time int64
	Event string // Either send or transmit
	Type string // type of msg being sent
	Duration int64
	MembersCount int
	TransactionCount int
}
