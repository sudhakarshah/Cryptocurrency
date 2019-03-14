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

	// TODO: open port to listen to other nodes in another thread
	go startListening(inbox, os.Args[4])
	connect_string := fmt.Sprintf("CONNECT %s %s %s\n", os.Args[4], ip, os.Args[3])
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
			nd := Node{}
			// send 

		case "TRANSACTION":
			// check if the transaction exists if so, continue

			// Insert transaction

			// TODO: Gossip the transaction
		case "DIE":
			os.Exit(3)
		case "QUIT":
			// Send everyone that 
		case "LEAVE":
			// find the leaving node's ip and name in the members and remove
		case  "PING":
		default:
			fmt.Printf("CANNOT PARSE MESSAGE. RECIEVED %s\n", tokens[0])
		}

	}

}
