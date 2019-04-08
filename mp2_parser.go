package main

import(
	"strings"
	"fmt"
	"strconv"
	"os"
	"net"
	"math"
	"crypto/sha256"
	"sort"
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
	QuesHash string
	SolHash string
}

type Block struct{
	PrevHash string
	Hash string
	Solution string
	Transactions []string
	Length int
	Accounts map[int]int
}

var PREV_HASH_IND = 1
var CURR_HASH_IND = 2
var SOLUTION_IND = 3
var LENGTH_IND = 4
var TRANSACTION_IND = 5
var ACCOUNT_IND = 6
var BLOCK = "BLOCK"
var SOLVED = "SOLVED"

func transactionsToString(ts []string)string{
	return strings.Join(ts, ",")
}
func accountsToString(acc map[int]int)string{
	var accList []string
	for i, k := range acc{
		account := strconv.Itoa(i)
		amount := strconv.Itoa(k)
		accList = append(accList, fmt.Sprintf("%s:%s",account,amount))
	}
	return strings.Join(accList, ",")
}

// outputs serialized sequence of transactions strings and updated account info
func serializeTransactions(th map[string]Msg, old_acc map[int]int)([]string, map[int]int){
	var msgList []Msg
	for _, v:=range th{
		msgList = append(msgList, v)
	}

	sort.SliceStable(msgList, func(i, j int)bool{
		return msgList[i].GetTimestamp() < msgList[i].GetTimestamp()
	})

	account := old_acc
	var output []string

	// check to see if the transactions are consistent
	for _, v := range msgList{
		source := v.GetSource()
		dest := v.GetDest()
		amount := v.GetAmount()

		s_money := account[source]
		if s_money - amount < 0{
			continue
		}

		account[source] -= amount
		account[dest] += amount
		output = append(output, v.Data)
	}
	return output, account
}

// Convert Block to Msg
func (b *Block)FormatMsg()Msg{
	tokens := make([]string, 7)
	tokens[0] = BLOCK
	tokens[PREV_HASH_IND] = b.PrevHash
	tokens[CURR_HASH_IND] = b.Hash
	tokens[LENGTH_IND] = strconv.Itoa(b.Length)
	tokens[SOLUTION_IND] = b.PrevHash
	tokens[TRANSACTION_IND] = transactionsToString(b.Transactions)
	tokens[ACCOUNT_IND] = accountsToString(b.Accounts)
	s := strings.Join(tokens, " ")
	s += "\n"
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


func (b *Block)addTrans(s string){
	b.Transactions = append(b.Transactions, s)
}

func (b *Block)TransactionCount()int{
	return len(b.Transactions)
}


// Used to convert string to Msg to queue to inbox
func (m *Msg)ParseBlock(s string){
	m.Data = s
	m.Type = "BLOCK"
}
func (m *Msg)ParseSolved(tokens ...string){
	m.Type = SOLVED
	m.QuesHash = tokens[1]
	m.SolHash = tokens[2]
}
// Convert Msg to Block
func (m *Msg)FormatBlock()Block{
	tokens := strings.Split(strings.TrimSpace(m.Data), " ")
	length, _ := strconv.Atoi(tokens[LENGTH_IND])
	var accounts map[int]int
	accounts = make(map[int]int)
	for _, v := range strings.Split(tokens[ACCOUNT_IND],","){
		tk := strings.Split(v,":")
		account, err := strconv.Atoi(tk[0])
		if err != nil{
			fmt.Printf("#ERROR NOT ABLE TO PARSE ACCOUNT %s\n", err)
		}
		amount, err := strconv.Atoi(tk[1])
		if err != nil{
			fmt.Printf("#ERROR NOT ABLE TO PARSE AMOUNT  %s\n", err)
		}
		accounts[account] = amount
	}
	return Block{PrevHash:tokens[PREV_HASH_IND], Hash:tokens[CURR_HASH_IND], Solution:tokens[SOLUTION_IND], Length:length, Transactions:strings.Split(tokens[TRANSACTION_IND], ","), Accounts:accounts}
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
	case "SOLVED":
		m.ParseSolved(tokens...)
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
func (m * Msg)GetData()string{
	return m.Data
}
func (b * Block)TrasactionCount()int{
	return len(b.Transactions)
}
