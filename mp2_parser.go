package main

import(
	"strings"
	"fmt"
	"strconv"
	"os"
)


type Msg struct{
	Type string
	Name string
	Ip string
	Port int
	Source int
	Dest int
	Amount int
	TID string
	TimeStamp string
	Friends []Node
	HashTable []Msg
}


func (m *Msg)Parse(s string){
	tokens := strings.Split(strings.TrimSpace(s), " ")
	if len(tokens) < 1{
		return
	}
	m.Type = tokens[0]
	switch tokens[0] {
	case "CONNECT":
		m.parseStandard(tokens...)
	case "INTRODUCE":
		m.parseStandard(tokens...)
	case "TRANSACTION":
		m.parseTransaction(tokens...)
	case "LEAVE":
		m.parseStandard(tokens...)
	case "JOIN":
		m.parseStandard(tokens...)
	case "DIE":
		return
	default:
		fmt.Printf("CANNOT PARSE MESSAGE. RECIEVED %s\n", tokens[0])
	}
}

func (m *Msg)FormatPing(friends map[string]Node, s string){
	tokens := strings.Split(strings.TrimSpace(s), " ")
	m.parseStandard(tokens...)
	m.Type ="PING"
	// Convert map to slice
	for _,v := range friends{
		m.Friends = append(m.Friends, v)
	}
}

func (m *Msg)GetFriends()map[string]Node{
	output := make(map[string]Node)
	for _, v := range m.Friends{
		k := fmt.Sprintf("%s:%s:%s",v.Name,v.Ip,v.Port)
		output[k] = v
	}
	return output
}
func (m *Msg)GetHashTable()map[string]Msg{
	output := make(map[string]Msg)
	for _, v := range m.HashTable{
		output[v.GetTID()] = v
	}
	return output
}

func (m *Msg)FormatInit(friends map[string]Node, hashtable map[string]Msg, s string){
	tokens := strings.Split(strings.TrimSpace(s), " ")
	m.parseStandard(tokens...)
	m.Type = "INIT"
	// Convert map to slice
	for _,v := range friends{
		m.Friends = append(m.Friends, v)
	}
	for _,v := range hashtable{
		m.HashTable = append(m.HashTable, v)
	}
}

func (m *Msg)parseStandard(tokens ...string){
	m.Type = tokens[0]
	m.Name = tokens[1]
	m.Ip = tokens[2]
	port, err := strconv.Atoi(tokens[3])
	if err != nil {
		fmt.Printf("CANNOT PARSE PORT. Got %x with length %d\n",tokens[3], len(tokens[3]))
		os.Exit(6)
	}
	m.Port = port
}
func (m *Msg)parseTransaction(tokens ...string){
	m.Type = tokens[0]
	m.TimeStamp = tokens[1]
	m.TID = tokens[2]
	source, err := strconv.Atoi(tokens[3])
	if err != nil {
		fmt.Println("CANNOT PARSE SOURCE. Got " + tokens[3])
		os.Exit(6)
	}
	m.Source = source
	dest, err := strconv.Atoi(tokens[4])
	if err != nil {
		fmt.Println("CANNOT PARSE Dest. Got " + tokens[4])
		os.Exit(6)
	}
	m.Dest = dest
	amount, err := strconv.Atoi(tokens[5])
	if err != nil {
		fmt.Println("CANNOT PARSE Amount. Got " + tokens[5])
		os.Exit(6)
	}
	m.Amount = amount
}
func (m * Msg)GetName()string{
	return m.Name
}

func (m * Msg)GetType()string{
	return m.Type
}

func (m * Msg)GetPort()string{
	return strconv.Itoa(m.Port)
}

func (m * Msg)GetIp()string{
	return m.Ip
}
func (m * Msg)SetIp(ip string){
	m.Ip = ip
}
func (m * Msg)HasIp()bool{
	if len(m.Ip) > 0{
		return true
	}
	return false
}
func (m * Msg)GetTimestamp()float64{
	ts, err := strconv.ParseFloat(m.TimeStamp, 64)
	if err != nil {
		return -1
	}
	return ts
}
func (m * Msg)GetTID()string{
	return m.TID
}
func (m * Msg)GetSource()int{
	return m.Source
}
func (m * Msg)GetDest()int{
	return m.Dest
}
func (m * Msg)GetAmount()int{
	return m.Amount
}
