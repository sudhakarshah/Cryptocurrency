package parser_mp2

import{
	"strings"
	"fmt"
	"strconv"
}

var TRANSACTION_COMMAND = "TRANSACTION"


// this package handles how the messages are formatted and parsed
type Msg interfase {
	Command string
	Loads()
	Dump() string
}

type NodeMsg struct{
	Command string
	Name string
	Ip string
	Port string
}

type TransactionMsg struct{
	Command string
	Timestamp float64
	Id string
	Source int
	Dest int
	Amount int
}

func (tm * TransactionMsg) parseTransaction(timestamp string, source string, dest string, amount string)(float64,int,int,int){
	ts := -1
	s := -1
	d := -1
	a := -1
	ts, err := strconv.ParseFloat(timestamp, 64)
	if err != nil {
		return ts, s, d, a, err
	}
	s, err := strconv.Atoi(source)
	if err != nil {
		return ts, s, d, a, err
	}
	d, err := strconv.Atoi(dest)
	if err != nil {
		return ts, s, d, a, err
	}
	a, err := strconv.Atoi(amount)
	if err != nil {
		return ts, s, d, a, err
	}
	return ts, s, d, a, nil
}

func (nm * NodeMsg)Loads(s string){
	tokens := strings.Split(strings.TrimSpace(s), " ")
	return Message{Command: tokens[0], Name: tokens[1], Ip: tokens[2], Port: tokens[3]}
}
func (tm * TransactionMsg)Loads(s string){
	tokens := strings.Split(strings.TrimSpace(s), " ")
	timestamp, source, dest, amount, err := parseTransaction(tokens[1], tokens[3], tokens[4], tokens[5])
	if err != nil{
		fmt.Println("ERROR IN PARSING")
	}
	return TransactionMsg{Command: }
}

func (nm * NodeMsg)Dump()string{
	return fmt.Sprintf("%s %s %s %s", nm.Command, nm.Name, nm.Ip, nm.Port)
}


