package main

import(
	"sync"
	"net"
)

type Node struct{
	IntroductionConnection net.Conn
	Messages Box
	Pred NList
	Succ NList
	//FingerTable 
	//HashTable 
	quit bool
	mux sync.Mutex
}


/*
func (nd *Node) FindSuccessor(id int){

}

func (nd *Node) FindPredecessor(id int){

}

func (nd *Node) Join(nn *Node){

}

func (nd *Node) InitFingerTable(nn *Node){

}

func (nd *Node) UpdateOthers(){

}

func (nd *Node) UpdateFingerTable(node_name string, i int){

}

func (nd *Node) Gossip(){

}

func (nd *Node) HandleMessage(s string){

}
*/



