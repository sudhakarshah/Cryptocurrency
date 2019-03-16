package main

import(
	"sync"
	"errors"
)

type Node struct{
	Name string
	Ip string
	Port string
	LastActive int64
	Attemps int
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
