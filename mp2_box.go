package main

import(
	"sync"
)

type Box struct{
	messages []string
	mux sync.Mutex
}

func (in*Box) enqueue(m string){
	in.mux.Lock()
	in.messages = append(in.messages, m)
	in.mux.Unlock()
}

func (in*Box) pop()(string,error){
	var output string
	var err error = nil
	in.mux.Lock()
	if len(in.messages) != 0{
		output = in.messages[0]
		in.messages = in.messages[1:]
	} else {
		err = errors.New("The inbox is empty")
	}
	in.mux.Unlock()
	return output, err
}


type FingerTable struct{
	[]*NodeInfo
}

func (ft *FingerTable)Add(){
	
}
func (ft *FingerTable)Remove(){
	
}

type HashTable struct{
	data map[string] Msg
}

func (ht *HashTable) AddTransaction(Msg){
	hash := ShaHash(Msg.GetTID())
	ht.data[hash] = Msg
}
func (ht *HashTable) RemoveTransaction(Msg){
	hash := ShaHash(Msg.GetTID())
	delete(ht.data, hash)
}
func (ht *HashTable) FindTransaction(Msg)Msg{
	hash := ShaHash(Msg.GetTID())
	return ht.data[hash]
}
