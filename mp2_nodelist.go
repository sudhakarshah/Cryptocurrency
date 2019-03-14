package main

import(
	"net"
	"fmt"
	"sync"
	"crypto/sha1"
)


type NodeInfo struct{
	Hash string
	Name string
	Ip string
	Port string
	Conn *net.Conn
}

type NList struct{
	Members []NodeInfo
	Mux sync.Mutex
}

func (nl *NList) Add(name string, ip string, port string, conn *net.Conn){
	hash = ShaHash(fmt.Sprintf("%s:%s",ip, port))
	n := NodeInfo{Hash:hash, Name:name, Ip:ip, Port:port, Conn:conn}
	nl.Members = append(nl.Members, n)
}

func (nl *NList) RemoveByHash(hash string){
	for i, n := range nl.Members {
		if n.Hash == hash {
			nl.Members = append(nl.Members[:i], nl.Members[i+1:]...)
			break
		}
	}
}

func ShaHash(s string)string{
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x",bs)
}
