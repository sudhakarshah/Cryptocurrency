package main

import(
	"strings"
	"fmt"
	"strconv"
)


type Msg struct{
	data map[string] string
}


func (m *Msg)Parse(s string){
	tokens := strings.Split(strings.TrimSpace(s), " ")
	switch tokens[0] {
	case "CONNECT":
		m.parseConnect(tokens...)
	case "INTRODUCE":
		m.parseIntroduce(tokens...)
	case "TRANSACTION":
		m.parseTransaction(tokens...)
	case "DIE":
		m.data["type"] = tokens[0]
	case "QUIT":
		m.data["type"] = tokens[0]
	default:
		fmt.Printf("CANNOT PARSE MESSAGE. RECIEVED %s\n", tokens[0])
	}
}

func (m *Msg)parseConnect(tokens ...string){
	m.data["type"] = tokens[0]
	m.data["name"] = tokens[1]
	m.data["ip"] = tokens[2]
	m.data["port"] = tokens[3]
}
func (m *Msg)parseIntroduce(tokens ...string){
	m.data["type"] = tokens[0]
	m.data["name"] = tokens[1]
	m.data["ip"] = tokens[2]
	m.data["port"] = tokens[3]
}
func (m *Msg)parseTransaction(tokens ...string){
	m.data["type"] = tokens[0]
	m.data["name"] = tokens[1]
	m.data["ip"] = tokens[2]
	m.data["port"] = tokens[3]
}

func (m * Msg)GetType()string{
	return m.data["type"]
}

func (m * Msg)GetPort()string{
	return m.data["port"]
}

func (m * Msg)GetIp()string{
	return m.data["ip"]
}
func (m * Msg)GetTimestamp()float64{
	ts, err := strconv.ParseFloat(m.data["timestamp"], 64)
	if err != nil {
		return -1
	}
	return ts
}
func (m * Msg)GetTID()string{
	return m.data["id"]
}
func (m * Msg)GetSource()int{
	s, err := strconv.Atoi(m.data["source"])
	if err != nil {
		return -1
	}
	return s
}
func (m * Msg)GetDest()int{
	s, err := strconv.Atoi(m.data["dest"])
	if err != nil {
		return -1
	}
	return s
}
func (m * Msg)GetAmount()int{
	s, err := strconv.Atoi(m.data["amount"])
	if err != nil {
		return -1
	}
	return s
}
