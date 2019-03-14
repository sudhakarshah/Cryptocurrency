package main

import (
	"fmt"
	"net"
	"os"
	"errors"
	"time"
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

func getRequest(conn *net.Conn)string {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 4096)
	// Read the incoming connection into the buffer.
	reqLen, err := (*conn).Read(buf)
	_ = reqLen
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	return string(buf)
}

func queueIntroRequest(inbox *Box, conn *net.Conn){
	for {
		s := getRequest(conn)
		m := Msg{}
		m.Parse(s)
		inbox.enqueue(m)
		time.Sleep(1)
	}
}

func printDebug(s string){
	t := int32(time.Now().Unix())
	fmt.Printf("[DEBUG]%d: %s\n",t,s)

}

// This is only for inter node communication
// TODO: Recieve Message and close connection
func listener(inbox * Box, in_con net.Conn){
	var m Msg
	dec := jsor.NewDecoder(*in_con)
	if err := dec.Decode(&m); err != nil{
		// Something went wrong
		return
	}
	inbox.enqueue(m)
}

// TODO: Spawn listern threads for each connection
func startListening(inbox * Box, port string){
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s",port))
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go listener(inbox, conn)
	}
}

func sendJson(ip string, port string, m Msg){
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		// TODO: cant dial
	}
	if e := json.NewEncoder(conn).Encode(m); e != nil {
		// TODO: json failed to send
	}
	conn.Close()
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
	members := []Node
	var hashtable map[string]Msg
	name := os.Args[4]
	port := os.Args[3]

	// TODO: open port to listen to other nodes in another thread
	go startListening(inbox, port)
	connect_string := fmt.Sprintf("CONNECT %s %s %s\n", name, ip, port)
	conn, err := connect_to_intro(os.Args[1], os.Args[2])
	fmt.Fprintf(conn, connect_string)
	go handleIntroRequest(&inbox, &conn)

	// handle requests
	for {
		m, err := inbox.pop()
		if err != nil{
			time.Sleep(1)
			continue
		}
		switch m.GetType() {
		case "INTRODUCE":
			// connect to the introduced node and send membership
			nd := Node{Name:m.GetName, Ip:m.GetIp(), Port:m.GetPort(), LastActive:-1}
			ping := Msg{}
			ping.FormatPingMessage(members)
			// send 
			sendJson(m.GetIp(), m.GetPort(), ping)

		case "TRANSACTION":
			// check if the transaction exists if so, continue

			if val, ok := hashtable[m.GetTID()]; ok{
				continue
			}

			// Insert transaction
			hashtable[m.GetTID()] = m

			// TODO: Gossip the transaction
			for i, nd := range members{
				coin := rand.Intn(1)
				if coin < 1{
					continue
				}
				sendJson(nd.Ip, nd.Port, m)
			}
		case "DIE":
			os.Exit(3)
		case "QUIT":
			quit := Msg{}
			quit.Parse(fmt.Sprintf("LEAVE %s %s %s", name, ip, port))
			// Send leave message to everyone
			for i, nd := range members{
				sendJson(nd.Ip, nd.Port, quit)
			}
			os.Exit(3)
		case "LEAVE":
			found := false
			// check if the leaving node is in  the members
			for i, nd := range members(
				if nd.IP == m.GetIP() && nd.Port == m.GetPort() && nd.Name == m.GetName(){
					found = true
					members = append(members[:i],members[i+1:]...)
					break
				}
			)
			if !found {
				continue
			}

			// If so, then forward the message to others
			for i, nd := range members{
				sendJson(nd.Ip, nd.Port, m)
			}
		case  "PING":
			pf := m.Friends
			for i, nd := range pf {
				
			}
		default:
			fmt.Printf("CANNOT PARSE MESSAGE. RECIEVED %s\n", tokens[0])
		}

	}

}
