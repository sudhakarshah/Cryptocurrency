package main

import(
	"sync"
	"errors"
	"net"
	"encoding/json"
	"fmt"
	"time"
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
	if nd.jsonEncoder == nil{
		nd.jsonEncoder = json.NewEncoder(nd.Sock)
	}
	startTime := time.Now()
	if e := nd.jsonEncoder.Encode(m); e != nil {
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
	b, _ := json.Marshal(m)
	fmt.Printf("SEND %d %s %d %d\n",int64(time.Now().Unix()), m.GetType(), len(b),  time.Since(startTime)) // time, msg type, size, duration
	return 0
}


func (nd *Node)ListenToFriend(inbox *Box){
	dec := json.NewDecoder(nd.Sock)
	for {
		var m Msg
		if err := dec.Decode(&m); err != nil{
			nd.mux.Lock()
			nd.Attempts++
			fmt.Printf("# Listen attempt %d.\n", nd.Attempts)
			fmt.Printf("# %s\n", err)
			if nd.Attempts > 2{
				fmt.Println("# Could not recieve message properly")
				nd.mux.Unlock()
				nd.Sock.Close()
				return
			}
			nd.mux.Unlock()
			continue
		}else{
			inbox.enqueue(m)
		}
		time.Sleep(10)
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
	var err error = nil
	in.mux.Lock()
	if len(in.messages) != 0{
		output = in.messages[0]
		in.messages = in.messages[1:]
	} else {
		err = errors.New("The inbox is empty")
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
