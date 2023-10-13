package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// stack

type NodeS struct {
	data string
	next *NodeS
}

type Stack struct {
	head *NodeS
}

// queue

type NodeQ struct {
	data string
	next *NodeQ
}

type Queue struct {
	head *NodeQ
	tail *NodeQ
}

//hash-table

type NodeH struct {
	Key   string
	Value string
}

type HashTable struct {
	Table map[int][]NodeH
}

func hashFunc(a int, key string) int {
	sum1 := 0
	sum2 := 0
	for _, c := range key {
		sum1 += int(c)
	}

	for _, c := range key {
		sum2 += (int(c) % 2)
	}

	ans := sum1 + a*sum2
	return ans
}

//set

type NodeSH struct {
	Key   string
	Value string
}

type Set struct {
	Set map[int][]NodeSH
}

//---------------------------------

func (stack *Stack) pop() (string, error) {
	if stack.head == nil {
		return "", errors.New("пустой стек")
	} else {
		x := stack.head.data
		stack.head = stack.head.next
		return x, nil
	}
}

func (stack *Stack) push(val string) {
	if stack.head == nil {
		stack.head = new(NodeS)
		stack.head.data = val
	} else {
		new(NodeS).next = stack.head
		stack.head.data = val
	}
}

//queue

func (queue *Queue) pushQ(val string) {
	if queue.head == nil {
		x := new(NodeQ)
		queue.head = x
		queue.head.data = val
		queue.tail = x
		queue.tail.data = val
	} else {
		new(NodeQ).next = queue.head
		queue.head.data = val
	}
}

func (queue *Queue) popQ() (string, error) {
	if queue.head == nil {
		return "", errors.New("")
	} else {
		data := queue.head.data
		queue.head = queue.head.next
		return data, nil
	}
}

//hash-table

func (ht *HashTable) Add(key, value string) {
	index := hashFunc(1, key)
	node := NodeH{Key: key, Value: value}
	if ht.Table[index] != nil {
		for i := 2; i < 32; i++ {
			index := hashFunc(i, key)
			if ht.Table[index] == nil {
				ht.Table[index] = []NodeH{node}
			}
		}
	} else {
		ht.Table[index] = []NodeH{node}
	}
}

func (ht *HashTable) Get(key string) (string, bool) {
	index := hashFunc(1, key)
	if ht.Table[index] != nil {
		for _, node := range ht.Table[index] {
			if node.Key == key {
				return node.Value, true
			}
		}
	}
	return "", false
}

func (ht *HashTable) Delete(key string) bool {
	index := hashFunc(1, key)
	if ht.Table[index] != nil {
		for _, node := range ht.Table[index] {
			if node.Key == key {
				node.Key = "0"
				node.Value = "0"
				return true
			}
		}
	}
	return false
}

//set

func (st *Set) AddS(key, value string) {
	index := hashFunc(1, key)
	node := NodeSH{Key: key, Value: "1"}
	if st.Set[index] != nil {
		for i := 2; i < 32; i++ {
			index := hashFunc(i, key)
			if st.Set[index] == nil {
				st.Set[index] = []NodeSH{node}
			}
		}
	} else {
		st.Set[index] = []NodeSH{node}
	}
}

func (st *Set) GetS(key string) (string, bool) {
	index := hashFunc(1, key)
	if st.Set[index] != nil {
		for _, node := range st.Set[index] {
			if node.Key == key {
				return node.Value, true
			}
		}
	}
	return "", false
}

func (st *Set) DeleteS(key string) bool {
	index := hashFunc(1, key)
	if st.Set[index] != nil {
		for _, node := range st.Set[index] {
			if node.Key == key {
				node.Key = "0"
				node.Value = "0"
				return true
			}
		}
	}
	return false
}

//------------------------------------

func main() {
	ht := HashTable{Table: make(map[int][]NodeH)}
	st := Set{Set: make(map[int][]NodeSH)}
	stack := &Stack{}
	queue := &Queue{}

	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	x := strings.Split(line, "--file ")
	y := strings.Split(x[1], "-- query ")

	//file := y[0]
	que := strings.Split(y[1], " ")

	if que[0] == "SPUSH" {
		stack.push(que[1])
	} else if que[0] == "SPOP" {
		stack.pop()
	} else if que[0] == "QPUSH" {
		queue.pushQ(que[1])
	} else if que[0] == "QPOP" {
		queue.popQ()
	} else if que[0] == "HSET" {
		keyandval := strings.Split(que[1], " ")
		ht.Add(keyandval[0], keyandval[1])
	} else if que[0] == "HDEL" {
		keyandval := strings.Split(que[1], " ")
		ht.Delete(keyandval[0])
	} else if que[0] == "HGET" {
		keyandval := strings.Split(que[1], " ")
		ht.Get(keyandval[0])
	} else if que[0] == "SADD" {
		keyandval := strings.Split(que[1], " ")
		st.AddS(keyandval[0], keyandval[1])
	} else if que[0] == "SREM" {
		keyandval := strings.Split(que[1], " ")
		st.DeleteS(keyandval[0])
	} else if que[0] == "SISMEMBER" {
		keyandval := strings.Split(que[1], " ")
		st.GetS(keyandval[0])
	}

	fmt.Println("Вы ввели:", strings.TrimSpace(line))
}
