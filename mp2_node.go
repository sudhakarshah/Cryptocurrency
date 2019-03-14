package main

import(
	"sync"
	"errors"
	"net"
)

type Node struct{
	Name string
	Ip string
	Port string
	LastActive int
}

type (nd *Node) Ping(friends []Node, conn *net.Conn){

}


type Box struct{
	messages []Msg
	mux sync.mutex
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


