package main

import (
	"fmt"
	"net"
	"os"
	"errors"
	"time"
	_"encoding/json"
	"math/rand"
	"strings"
	"bufio"
)

var DEBUG = true


func connect_to_intro(ip string, port string)(net.Conn, error){
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s",ip,port))
	if err != nil {
		return nil, err
	}
	return conn, nil
}


func get_my_ip() (string, error) {
// This function is from https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func getRequest(conn net.Conn)string {
	// Make a buffer to hold incoming data.
	status, _ := bufio.NewReader(conn).ReadString('\n')
	return strings.TrimSpace(status)
}

/* queuing messages from the service*/
func queueIntroRequest(inbox *Box, conn net.Conn){
	reader := bufio.NewReaderSize(conn, 200)
	for {
		s, _ := reader.ReadString('\n')
		//fmt.Printf("# Recieved %d from intro\n", len(s))
		if len(s) > 2{
			m := Msg{}
			if m.Parse(s) < 0{
				continue
			}
			inbox.enqueue(m)
			time.Sleep(10)
		}
	}
}

func printDebug(s string){
	t := int64(time.Now().Unix())
	fmt.Printf("[DEBUG]%d: %s\n",t,s)

}

// This is only for inter node communication
// TODO: Recieve Message and close connection
func listener(inbox * Box, in_con net.Conn){
	var m Msg
	s, err := bufio.NewReader(in_con).ReadString('\n')
	fmt.Printf("# recieved string %s\n", s)
	if err != nil{
		fmt.Println("#Error in listening")
		fmt.Printf("# %s", err)
		// Something went wrong
		return
	}
	m.Parse(s)
	m.PutSock(in_con)
	inbox.enqueue(m)
}

// TODO: Spawn listern threads for each connection
func startListening(inbox * Box, port string){
	//fmt.Println("Started Listening on " + port)
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		// handle error
		fmt.Printf("# [ERROR] %s", err)

	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Printf("# [ERROR] %s", err)
		}
		go listener(inbox, conn)
	}
}

func queueHB(inbox *Box){
	for{
		m := Msg{Type:"HB", Data:"HB\n"}
		time.Sleep(10 * time.Second)
		inbox.enqueue(m)
	}
}

func queueBlockSolveBeat(inbox *Box){
	for{
		m := Msg{Type:"BLOCK_SOLVE", Data:"BLOCK_SOLVE\n"}
		time.Sleep(1 * time.Second)
		inbox.enqueue(m)
	}
}

func countConnectedNodes(members map[string]*Node)(int) {
	count := 0
	for _, v := range members{
		if (v.isConnected) {
			count++
		}
	}
	return count
}

func getUnconnectedNode(members map[string]*Node)(string) {
	for k, v := range members{
		if (!v.isConnected) {
			return k
		}
	}
	return ""
}

func updateIncludedTransactions(includedTransactions map[string]int, b Block) {
	for i := 0;  i < len(b.Transactions); i++ {
		includedTransactions[b.Transactions[i]] = 1
	}
}

func updateAccount(accounts map[int]int, m Msg)int{
	isPossible := false
	if m.Source == 0 {
		isPossible = true
	}
	balance, ok := accounts[m.Source]
	if ok && balance >= m.Amount {
		accounts[m.Source] -= m.Amount
		isPossible = true
	}

	if isPossible {
		// if dest account doesnt exists
		if _, ok := accounts[m.Dest]; !ok {
			accounts[m.Dest] = 0
		}
		accounts[m.Dest] += m.Amount
		return 0
	}
	return -1

}
func main(){
	// Expects 3 arguments: ip, port,
	if len(os.Args) != 5 {
		fmt.Println("Expected 3 arguments: Intro ip, Intro port, Local Listening Port, Name")
		return
	}
	ip, err := get_my_ip()

	if err != nil{
		fmt.Println("Could not get local ip")
		fmt.Println(err)
	}

	inbox := Box{}
	var members map[string]*Node
	var hashtable map[string]Msg
	var accounts map[int]int
	var tempAccounts map[int]int
	var includedTransactions map[string]int
	//var connected_members map[string]*Node
	members = make(map[string]*Node)
	hashtable = make(map[string]Msg)
	accounts = make(map[int]int)
	includedTransactions = make(map[string]int)

	currentSolvingBlock := Block{}
	name := os.Args[4]
	port := os.Args[3]
	// queue of blocks in the current chain
	chain := Box{}
	required_connection := 4;
	// TODO: open port to listen to other nodes in another thread
	go startListening(&inbox, port)
	connect_string := fmt.Sprintf("CONNECT %s %s %s\n", name, ip, port)
	connService, err := connect_to_intro(os.Args[1], os.Args[2])
	if err != nil{
		//fmt.Println(err)
	}
	fmt.Fprintf(connService, connect_string)
	go queueIntroRequest(&inbox, connService)
	go queueHB(&inbox)
	go queueBlockSolveBeat(&inbox)
	// gossip flipper
	//send := true

	// handle requests
	for {
		m, err := inbox.pop()
		// sleeping if no message in inbox
		if err != nil{
			time.Sleep(10)
			continue
		}
		fmt.Printf("RECIEVED %d %s %d %d %d\n",int64(time.Now().Unix()), m.GetType(), len(m.Data), len(members), len(hashtable) )
		//fmt.Printf("Members Count: %d\nTransaction Count: %d\n", len(members), len(hashtable))
		switch m.GetType() {
		case "INTRODUCE":
			target_id := fmt.Sprintf("%s:%s:%s",m.GetName(),m.GetIp(),m.GetPort())
			my_id := fmt.Sprintf("%s:%s:%s",name,ip,port)
			// already known node or itself
			if _, ok := members[target_id];ok || target_id == my_id {
				continue
			}

			// less than 4 connected
			if (countConnectedNodes(members) < required_connection) {
				// send JOIN MESSAGE then INIT message
				join := Msg{}
				join.Parse(fmt.Sprintf("JOIN %s %s %s\n", name, ip, port))
				conn, err := net.Dial("tcp", m.GetIp()+":"+m.GetPort())

				if err != nil {
					fmt.Printf("# Cannot connect to the introduced node %s\n", m.GetName())
					continue
				}

				var nd Node
				nd = Node{Name:m.GetName(), Ip:m.GetIp(), Port:m.GetPort(), LastActive:int64(time.Now().Unix()), Sock:conn, Attempts:0,isConnected:true}
				for _, nd := range members{
					nd.SendJson(m)
				}
				// send
				if nd.SendJson(join) == 0{
					members[target_id] = &nd
					go nd.ListenToFriend(&inbox)
				}
			} else {
				var nd Node
				nd = Node{Name:m.GetName(), Ip:m.GetIp(), Port:m.GetPort(), LastActive:int64(time.Now().Unix()), Sock:nil, Attempts:0,isConnected:false}
				members[target_id] = &nd
			}

		case "JOIN":
			// send INTRODUCE and TRANACTION MESSAGES
			target_id := fmt.Sprintf("%s:%s:%s",m.GetName(),m.GetIp(),m.GetPort())
			if _, ok := members[target_id];ok{
				continue
			}
			ping := FormatPing(members)
			nd := Node{Name:m.GetName(), Ip:m.GetIp(), Port:m.GetPort(), LastActive:int64(time.Now().Unix()), Sock:m.Sock, Attempts:0, isConnected:true}
			go nd.ListenToFriend(&inbox)
			members[fmt.Sprintf("%s:%s:%s",m.GetName(),m.GetIp(),m.GetPort())] = &nd
			for _, msg := range ping{
				nd.SendJson(msg)
			}

		case "TRANSACTION":
			// check if the transaction exists if so, continue

			if _, ok := hashtable[m.GetTID()]; ok{
				//fmt.Println("TRANSACTION EXISTS")
				continue
			}

			// Insert transaction
			hashtable[m.GetTID()] = m
			fmt.Printf("UPDATE %d %s %s\n",int64(time.Now().Unix()), m.GetType(), m.GetTID() ) // time, msg type, size, member_count, transaction_count

			var removeList []string

			for k, v := range members{
				if rand.Intn(3) != 0 {
					continue
				}
				if v.SendJson(m) != 0 {
					fmt.Printf("# Could not send message to %s\n", v.Name)
					removeList = append(removeList, k)
				}
			}
			for _, k := range removeList{
				//fmt.Printf("Removing %s from members\n", k)
				delete(members, k)
			}
		case "DIE":
			os.Exit(3)
		case "QUIT":
			quit := Msg{}
			quit.Parse(fmt.Sprintf("LEAVE %s %s %s\n", name, ip, port))
			// Dump every remaining message to all members
			os.Exit(3)
		case "LEAVE":
			key := fmt.Sprintf("%s:%s:%s",m.GetName(),m.GetIp(),m.GetPort())
			if _, ok := members[key]; !ok {
				continue
			}
			delete(members, key)
			// If so, then forward the message to others
		case "HB":
			// Array of introduce
			ping := FormatPing(members)
			var removeList []string
			// Ping everyone in the contacts
			for k, v := range members{
				if rand.Intn(10) != 0{
					continue
				}
				if v.Ip == ip && v.Port == port {
					continue
				}
				for i, p := range ping{
					if i > 3{
						break
					}
					if v.SendJson(p) != 0 {
						removeList = append(removeList, k)
					}else{
						//fmt.Printf("Pinged %s\n", k)
					}
				}
			}
			// If ping failed, remove from the contacts
			for _, k := range removeList{
				fmt.Printf("# Removing %s from members\n", k)
				delete(members, k)
			}

			// If number of connections not sufficient, connect to more nodes
			count := countConnectedNodes(members)
			// looping the required number of times
			for i := count; i < required_connection; i++ {
				target_id := getUnconnectedNode(members)
				if (target_id == ""){
					break
				}
				// actually a nd pointer
				nd := members[target_id]
				// send JOIN MESSAGE then INIT message
				join := Msg{}
				join.Parse(fmt.Sprintf("JOIN %s %s %s\n", name, ip, port))
				conn, err := net.Dial("tcp", nd.Ip+":"+nd.Port)

				if err != nil {
					fmt.Printf("# Cannot connect to the introduced node %s\n", nd.Name)
					continue
				}
				nd.Sock = conn
				nd.isConnected = true
				if nd.SendJson(join) == 0 {
					members[target_id] = nd
					go nd.ListenToFriend(&inbox)
				}
			}
		case "SOLVED":
			if m.QuesHash != currentSolvingBlock.Hash {
				continue
			}

			currentSolvingBlock.Solution = m.SolHash
			chain.addBlock(currentSolvingBlock)
			updateIncludedTransactions(includedTransactions, currentSolvingBlock);
			accounts = currentSolvingBlock.Accounts
			// undo the solve request sent to the service


			// Broadcasting the new block to everyone
			msg := currentSolvingBlock.FormatMsg()
			fmt.Println("SEND_BLOCK %d %s\n", int64(time.Now().Unix()), currentSolvingBlock.Hash)
			for _, v := range members{
				if v.SendJson(msg) != 0 {
					fmt.Printf("# Could not send block message to %s\n", v.Name)
					//removeList = append(removeList, k)
				}
			}
			currentSolvingBlock = Block{Hash:""}



		case "BLOCK_SOLVE":
			// looping through all messages i have
			newBlock := Block{}

			tempAccounts = make(map[int]int)
			for k,v := range accounts {
			  tempAccounts[k] = v
			}


			for k,v := range hashtable {
				// if it has not been included in blocks before
				if _, ok := includedTransactions[k]; !ok {
					if updateAccount(tempAccounts, v) == -1 {
						continue
					}
					newBlock.addTrans(k)
				}
			}
			newBlock.Accounts = tempAccounts
			lastBlock, err  := chain.peepBack()
			// no block in this chain
			if err != nil {
				newBlock.PrevHash = "0"
				newBlock.Length = 1
			} else {
				newBlock.PrevHash = lastBlock.Hash
				newBlock.Length = lastBlock.Length + 1
			}

			// sending solve message to the service
			fmt.Fprintf(connService, newBlock.FormatSolve())
			currentSolvingBlock = newBlock


		case "BLOCK":
			b := m.FormatBlock()
			fmt.Printf("RECIEVED_BLOCK %d %s\n", int64(time.Now().Unix()), b.Hash)
			// len := len(chain.messages)
			lastBlock, err  := chain.peepBack()

			updateChain := false
			// if no block then accept whatever you get and add it to ur chain
			if err != nil {
				updateChain = true
			} else if (b.Length > lastBlock.Length) || (b.Length == lastBlock.Length && (b.TransactionCount() > lastBlock.TransactionCount() || b.Hash > lastBlock.Hash)) {
				// replacing last block with new block if new length larger. TODO: need to add tie breaking strategy
				updateChain = true
			}
			ansestorLen := -1
			if updateChain {
				// on success scenario handled, that is prev block present in the queue
				ansestorLen = chain.addBlock(b)
				updateIncludedTransactions(includedTransactions, b);
				accounts = b.Accounts
				// undo the solve request sent to the service
				currentSolvingBlock = Block{Hash:""}

				fmt.Printf("ACCEPTED %d %s %s\n", int64(time.Now().Unix()), b.Hash, strings.Join(b.Transactions, ","))

				// gossiping the block to other nodes
				for _, v := range members{
					if v.SendJson(m) != 0 {
						fmt.Printf("# Could not send block message to %s\n", v.Name)
						//removeList = append(removeList, k)
					}
				}
			}

			if (b.Length == lastBlock.Length) {
				fmt.Printf("CHAIN_SPLIT %d %s %d\n", int64(time.Now().Unix()), b.Hash, ansestorLen)
			}



		default:
			fmt.Printf("# CANNOT PARSE MESSAGE. RECIEVED %+v\n",m )
		}

	}

}
