package routing

import (
	//"fmt"
	"sync"
	//"net"
	//"math/rand"
	//"log"
)

const (
	alpha = 3
	findsize = 16
)

type table struct {
	bucket	[]Node
	mutex	sync.Mutex
	selfID	Hash
}

type Hash int

type Node interface {
	address()	string
	ID()	Hash
}

type nodesByDistance struct {
	entries []Node
	target	Hash
}


func (t *table)	GetNodeAddress(targetID Hash) []Node {
	var (
		asked	= make(map[Hash]bool)
		result	*nodesByDistance
		seen	= make(map[Hash]bool)
		reply	= make(chan []Node, alpha)
		pendingQueries = 0
	)
	
	asked[t.selfID] = true
	
	t.mutex.Lock()
	result = t.closest(targetID, findsize)
	t.mutex.Unlock()
	
	for _,node := range result.entries {
		if node.ID() == targetID{
			return result.entries
		}
	}
	//fmt.Println(result.entries)
	for {
		for i := 0; i < len(result.entries) && pendingQueries < alpha; i++ {
			n := result.entries[i]
			if !asked[n.ID()] {
				asked[n.ID()] = true
				//seen[n.ID] = true
				pendingQueries++
				
				go t.findnode(n, targetID, reply)
			}
		}
		if pendingQueries == 0 {
			// we have asked all closest nodes, stop the search
			break
		}
		// wait for the next reply
		for _, n := range <-reply {
			if n != nil && !seen[n.ID()] {
				seen[n.ID()] = true
				result.push(n, findsize)
			}
		}
		pendingQueries--
	}
	
	return result.entries
	
}

func (t *table) findnode(n Node, targetID Hash, reply chan<- []Node) {
	
	//send and receive simulation
	nodes := make([]Node, 0, findsize)
	/*
	for i := 0; i < findsize; i++ {
		x:= rand.Intn(100)
		node := &Node{address:"ha",ID:Hash(x)}
		nodes = append(nodes , node)
	}
	*/
	reply <- nodes
}

func (t *table) closest(target Hash, nresults int) *nodesByDistance {
	closeset := &nodesByDistance{target: target}
	for _, n := range t.bucket {
		if n != nil {
			closeset.push(n, nresults)
		}
	}
	return closeset
}

func (h *nodesByDistance) push(n Node, maxElems int) {
	h.entries = append(h.entries,n)
	for i , node := range h.entries{
		if distance(node.ID() , h.target) > distance(n.ID() , h.target) {
			copy(h.entries[i+1:],h.entries[i:])
			h.entries[i] = n
			break
		}
	}
	if len(h.entries) > maxElems {
		h.entries = h.entries[:maxElems]
	}
}

func distance(a Hash, b Hash) int {
	if a > b {
		return int(a-b)
	}else {
		return int(b-a)
	}
}

/*

type node struct{
	addr	string
	id	Hash
}

func (n node) address() string{
	return n.addr
}

func (n node) ID() Hash{
	return n.id
}

func main() {
	var t table
	t.selfID = 41
	n1 := &node{addr:"na",id:43}
	n2 := &node{addr:"na",id:44}
	n3 := &node{addr:"na",id:45}
	n4 := &node{addr:"nb",id:46}
	n5 := &node{addr:"na",id:47}
	n6 := &node{addr:"na",id:48}
	n7 := &node{addr:"na",id:49}
	n8 := &node{addr:"na",id:40}
	var buckets []*node
	buckets = append(buckets,n1)
	buckets = append(buckets,n2)
	buckets = append(buckets,n3)
	buckets = append(buckets,n4)
	buckets = append(buckets,n5)
	buckets = append(buckets,n6)
	buckets = append(buckets,n7)
	buckets = append(buckets,n8)
	
	for _,m := range buckets{
		t.bucket = append(t.bucket,m)
	}
	
	nodes := t.closest(46,4)
	fmt.Println("closest test:closest(46,4)")
	for _,n := range nodes.entries{
		fmt.Println(n.ID())
	}
	
	
}
*/