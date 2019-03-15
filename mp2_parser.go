package main

import(
	"strings"
	"fmt"
	"strconv"
)


type Msg struct{
	data map[string] string
	Friends map[string] Node
	HashTable map[string] Msg
}


func (m *Msg)Parse(s string){
	m.data = make(map[string]string)
	tokens := strings.Split(strings.TrimSpace(s), " ")
	m.data["type"] = tokens[0]
	switch tokens[0] {
	case "CONNECT":
		m.parseConnect(tokens...)
	case "INTRODUCE":
		m.parseIntroduce(tokens...)
	case "TRANSACTION":
		m.parseTransaction(tokens...)
	case "LEAVE":
		m.parseConnect(tokens...)
	case "JOIN":
		m.parseConnect(tokens...)
	default:
		fmt.Printf("CANNOT PARSE MESSAGE. RECIEVED %s\n", tokens[0])
	}
}

func (m *Msg)FormatPing(friends map[string]Node){
	m.data = make(map[string]string)
	m.data["type"]="PING"
	m.Friends = friends
}

func (m *Msg)FormatInit(friends map[string]Node, hashtable map[string]Msg, s string){
	m.data = make(map[string]string)
	tokens := strings.Split(strings.TrimSpace(s), " ")
	m.parseConnect(tokens...)
	m.data["type"]="INIT"
	m.Friends = friends
	m.HashTable = hashtable
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
	m.data["time"] = tokens[1]
	m.data["id"] = tokens[2]
	m.data["source"] = tokens[3]
	m.data["dest"] = tokens[4]
	m.data["amount"] = tokens[5]
}
func (m * Msg)GetName()string{
	return m.data["name"]
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
func (m * Msg)SetIp(ip string){
	m.data["ip"] = ip
}
func (m * Msg)HasIp()bool{
	_, ok := m.data["ip"]
	return ok
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
