package main

import (
	"fmt"
	"net"
	"os"
	"errors"
	"time"
	"encoding/json"
	"math/rand"
	"strings"
	"bufio"
	"sync"
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

func queueIntroRequest(inbox *Box, conn net.Conn){
	for {
		s := strings.TrimRight(getRequest(conn), "\n")
		if len(s) > 0{
			m := Msg{}
			m.Parse(s)
			inbox.enqueue(m)
			time.Sleep(1)
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
	dec := json.NewDecoder(in_con)
	if err := dec.Decode(&m); err != nil{
		// Something went wrong
		return
	}
	//fmt.Printf("Recieved %s\n", m.GetType())
	inbox.enqueue(m)
}

// TODO: Spawn listern threads for each connection
func startListening(inbox * Box, port string){
	//fmt.Println("Started Listening on " + port)
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		// handle error
		//fmt.Printf("[ERROR] %s", err)

	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			//fmt.Printf("[ERROR] %s", err)
		}
		go listener(inbox, conn)
	}
}

func sendJson(ip string, port string, m Msg)int{
	//fmt.Printf("Sending %s to port %d\n", m.GetType(), len(port))
	conn, err := net.Dial("tcp", ip+":"+port)
	startTime := int64(time.Now().Unix())
	if err != nil {
		// TODO: cant dial
		//fmt.Printf("[ERROR] %s", err)
		return -1
	}
	if e := json.NewEncoder(conn).Encode(m); e != nil {
		// TODO: json failed to send
		conn.Close()
		//fmt.Printf("[ERROR JSON] %s", e)
		return -1
	}
	endTime := int64(time.Now().Unix())
	b, _ := json.Marshal(m)
	fmt.Printf("SEND %d %s %d %d\n",int64(time.Now().Unix()), m.GetType(), len(b),  endTime - startTime) // time, msg type, size, duration
	return 0
}

func enable_hb(hb *bool, hbmux *sync.Mutex){
	time.Sleep(1000 * time.Millisecond)
	(*hbmux).Lock()
	(*hb) = true
	(*hbmux).Unlock()
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
	var members map[string]Node
	var hashtable map[string]Msg
	members = make(map[string]Node)
	hashtable = make(map[string]Msg)
	name := os.Args[4]
	port := os.Args[3]

	heartbeat := false
	var hbmux sync.Mutex

	// TODO: open port to listen to other nodes in another thread
	go startListening(&inbox, port)
	connect_string := fmt.Sprintf("CONNECT %s %s %s\n", name, ip, port)
	conn, err := connect_to_intro(os.Args[1], os.Args[2])
	if err != nil{
		//fmt.Println(err)
	}
	fmt.Fprintf(conn, connect_string)
	go queueIntroRequest(&inbox, conn)
	go enable_hb(&heartbeat, &hbmux)

	// handle requests
	for {
		m, err := inbox.pop()
		if err != nil{
			time.Sleep(1)
			continue
		}
		b, _ := json.Marshal(m)
		fmt.Printf("RECIEVED %d %s %d %d %d\n",int64(time.Now().Unix()), m.GetType(), len(b), len(members), len(hashtable) ) // time, msg type, size, member_count, transaction_count
		//fmt.Printf("Members Count: %d\nTransaction Count: %d\n", len(members), len(hashtable))
		hbmux.Lock()
		if heartbeat{
			ping := Msg{}
			ping.FormatPing(members, fmt.Sprintf("PING %s %s %s\n", name, ip, port))
			var removeList []string
			// Ping everyone in the contacts
			for k, v := range members{
				if v.Ip == ip && v.Port == port {
					continue
				}
				if sendJson(v.Ip, v.Port, ping) != 0 {
					removeList = append(removeList, k)
				}else{
					//fmt.Printf("Pinged %s\n", k)
				}
			}
			// If ping failed, remove from the contacts
			for _, k := range removeList{
				//fmt.Printf("Removing %s from members\n", k)
				delete(members, k)
			}
			heartbeat = false

		}
		hbmux.Unlock()
		switch m.GetType() {
		case "INTRODUCE":
			// send JOIN MESSAGE then INIT message
			join := Msg{}
			join.Parse(fmt.Sprintf("JOIN %s %s %s", name, ip, port))

			// send 
			if sendJson(m.GetIp(), m.GetPort(), join) == 0{
				nd := Node{Name:m.GetName(), Ip:m.GetIp(), Port:m.GetPort(), LastActive:int64(time.Now().Unix())}
				members[fmt.Sprintf("%s:%s:%s",m.GetName(),m.GetIp(),m.GetPort())] = nd
			}
		case "JOIN":
			//fmt.Printf("Sending INIT msg to %s\n", m.GetName())
			// Send init message
			init := Msg{}
			init.FormatInit(members, hashtable, fmt.Sprintf("INIT %s %s %s", name, ip, port))
			if sendJson(m.GetIp(), m.GetPort(), init) == 0{
				//fmt.Printf("Sent INIT msg to %s\n", m.GetName())
				nd := Node{Name:m.GetName(), Ip:m.GetIp(), Port:m.GetPort(), LastActive:int64(time.Now().Unix())}
				members[fmt.Sprintf("%s:%s:%s",m.GetName(),m.GetIp(),m.GetPort())] = nd
			}
		case "INIT":
			nd := Node{Name:m.GetName(), Ip:m.GetIp(), Port:m.GetPort(), LastActive:int64(time.Now().Unix())}
			members[fmt.Sprintf("%s:%s:%s",m.GetName(),m.GetIp(),m.GetPort())] = nd
			for k, v := range m.GetFriends(){
				if _, ok := members[k]; ok{
					//fmt.Println("Friend Exists")
					continue
				}
				intro := Msg{}
				intro.Parse(fmt.Sprintf("INTRODUCE %s %s %s", v.Name, v.Ip, v.Port))
				inbox.enqueue(intro)
			}
			for k, v := range m.GetHashTable(){
				hashtable[k] = v
			}

		case "TRANSACTION":
			// check if the transaction exists if so, continue

			if _, ok := hashtable[m.GetTID()]; ok{
				//fmt.Println("TRANSACTION EXISTS")
				continue
			}

			// Insert transaction
			hashtable[m.GetTID()] = m

			i := 0

			var removeList []string

			for k, v := range members{
				coin := rand.Intn(1)
				if coin < 1 && i > 3{
					continue
				}
				if sendJson(v.Ip, v.Port, m) != 0 {
					removeList = append(removeList, k)
				}
				i += 1
			}
			for _, k := range removeList{
				//fmt.Printf("Removing %s from members\n", k)
				delete(members, k)
			}
		case "DIE":
			os.Exit(3)
		case "QUIT":
			quit := Msg{}
			quit.Parse(fmt.Sprintf("LEAVE %s %s %s", name, ip, port))
			// Send leave message to everyone
			for _, nd := range members{
				sendJson(nd.Ip, nd.Port, quit)
			}
			os.Exit(3)
		case "LEAVE":
			key := fmt.Sprintf("%s:%s:%s",m.GetName(),m.GetIp(),m.GetPort())
			if _, ok := members[key]; !ok {
				continue
			}
			delete(members, key)
			// If so, then forward the message to others
			/*
			for _, nd := range members{
				sendJson(nd.Ip, nd.Port, m)
			}
			*/
		case  "PING":
			senderId := fmt.Sprintf("%s:%s:%s",m.GetName(),m.GetIp(),m.GetPort())
			// update the last active 
			sender := members[senderId]
			sender.LastActive = int64(time.Now().Unix())
			members[senderId] = sender
			pf := m.GetFriends()
			//fmt.Printf("Received ping with friend list length of %d\n", len(pf) )
			for k, v := range pf {
				if v.Ip == ip && v.Port == port {
					continue
				}
				if _, ok := members[k]; ok {
					continue
				}
				intro := Msg{}
				introString := fmt.Sprintf("INTRODUCE %s %s %s\n", v.Name, v.Ip, v.Port)
				intro.Parse(introString)
				inbox.enqueue(intro)
			}
		default:
			//fmt.Printf("CANNOT PARSE MESSAGE. RECIEVED %+v\n",m )
		}

	}

}
