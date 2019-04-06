package main

import(
	"strings"
	"fmt"
	"strconv"
	"os"
	"net"
	"math"
	"crypto/sha256"
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
	Sock net.Conn
	Data string
}

type Block struct{
	PrevHash string
	Hash string
	Solution string
	Transactions []string
}

var PREV_HASH_IND = 1
var CURR_HASH_IND = 2
var SOLUTION_IND = 3
var TRANSACTION_IND = 4
var BLOCK = "BLOCK"

// Convert Block to Msg
func (b *Block)FormatMsg()Msg{
	tokens := make([]string, 4)
	tokens[0] = BLOCK
	tokens[PREV_HASH_IND] = b.PrevHash
	tokens[CURR_HASH_IND] = b.Hash
	tokens[SOLUTION_IND] = b.PrevHash
	tokens = append(tokens, b.Transactions...)
	s := strings.Join(tokens, " ")
	return Msg{Type:BLOCK, Data:s}
}

func (b *Block)FormatVerify()string{
	return fmt.Sprintf("VERIFY %s %s\n", b.Hash, b.Solution)
}

func (b *Block)FormatSolve()string{
	if b.Hash == ""{
		b.generateHash()
	}
	return fmt.Sprintf("SOLVE %s\n", b.Hash)
}

func (b *Block)generateHash(){
	// Format input string for hash
	format := append(b.Transactions, b.PrevHash)
	s := strings.Join(format, " ")
	b.Hash = fmt.Sprintf("%x",sha256.Sum256([]byte(s)))
}

// Used to convert string to Msg to queue to inbox
func (m *Msg)ParseBlock(s string){
	m.Data = s
	m.Type = "BLOCK"

}

// Convert Msg to Block
func (m *Msg)FormatBlock()Block{
	tokens := strings.Split(strings.TrimSpace(m.Data), " ")
	return Block{PrevHash:tokens[PREV_HASH_IND], Hash:tokens[CURR_HASH_IND], Solution:tokens[SOLUTION_IND], Transactions:tokens[TRANSACTION_IND:]}
}

func (m *Msg)Parse(s string)int{
	m.Data = s
	tokens := strings.Split(strings.TrimSpace(s), " ")
	if len(tokens) < 1 || len(s) < 1{
		return -1
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
		return 1
	case BLOCK:
		m.ParseBlock(s)
	default:
		fmt.Printf("CANNOT PARSE STRING. RECIEVED %s\n", tokens[0])
		return -1
	}
	return 1
}

func (m *Msg)PutSock(conn net.Conn){
	m.Sock = conn
}

// creates a ping then sequence of introduce messages
func FormatPing(friends map[string]*Node)[]Msg{
	var output []Msg
	for _,v := range friends{
		m := Msg{}
		m.Parse(fmt.Sprintf("INTRODUCE %s %s %s\n", v.Name, v.Ip, v.Port))
		output = append(output, m)
	}
	return output
}


func FormatInit(friends map[string]*Node, hashtable map[string]Msg, fc int)[]Msg{
	var output []Msg
	for _,v := range friends{
		m := Msg{}
		m.Parse(fmt.Sprintf("INTRODUCE %s %s %s\n", v.Name, v.Ip, v.Port))
		output = append(output, m)
	}
	chance := int(math.Ceil(float64(fc)/2.0))
	if chance == 0{chance = 1}
	for _,v := range hashtable{
		//if rand.Intn(chance) == 0 {
		output = append(output, v)
		//}
	}
	return output
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
