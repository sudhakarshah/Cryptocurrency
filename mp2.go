package main

import (
	"fmt"
	"net"
	"os"
	"errors"
)


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

func handleRequest(conn *net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 4096)
	// Read the incoming connection into the buffer.
	reqLen, err := (*conn).Read(buf)
	_ = reqLen
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println(string(buf))
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

	connect_string := fmt.Sprintf("CONNECT %s %s %s\n", os.Args[4], ip, os.Args[3])
	conn, err := connect_to_intro(os.Args[1], os.Args[2])
	fmt.Fprintf(conn, connect_string)

	// Handles incoming requests.
	for {
		handleRequest(&conn)
	}

}
