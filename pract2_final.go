package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

// stack

type NodeS struct {
	data string
	next *NodeS
}

type Stack struct {
	Name string
	head *NodeS
}

// queue

type NodeQ struct {
	data string
	next *NodeQ
}

type Queue struct {
	Name string
	head *NodeQ
	tail *NodeQ
}

//hash-table

type NodeH struct {
	Key   string
	Value string
}

type HashTable struct {
	Name     string
	Table    []*NodeH
	Capacity int
}

func (ht *HashTable) hashFunc(a int, key string) int {
	sum1 := 0
	sum2 := 0
	for _, c := range key {
		sum1 += int(c)
	}

	for _, c := range key {
		sum2 += (int(c) % 2)
	}

	ans := (sum1 + a*sum2) % ht.Capacity
	return ans
}

func NewHashTable(capacity int, stname string) *HashTable {
	return &HashTable{
		Name:     stname,
		Table:    make([]*NodeH, capacity),
		Capacity: capacity,
	}
}

//set

type NodeSH struct {
	Key   string
	Value string
}

type Set struct {
	Name     string
	Table    []*NodeSH
	Capacity int
}

func (st Set) hashfuncSet(a int, key string) int {
	sum1 := 0
	sum2 := 0
	for _, c := range key {
		sum1 += int(c)
	}

	for _, c := range key {
		sum2 += (int(c) % 2)
	}

	ans := (sum1 + a*sum2) % st.Capacity
	return ans
}

func newSet(capacity int, stname string) *Set {
	return &Set{
		Name:     stname,
		Table:    make([]*NodeSH, capacity),
		Capacity: capacity,
	}
}

type DataStructure struct {
	name       string
	hashTables []HashTable //`json:"hashTables"`
	stacks     []Stack     //`json:"stacks"`
	queues     []Queue     //`json:"queues"`
	sets       []Set       //`json:"sets"`
}

type MainDataStructure struct {
	datastructures []DataStructure
	mutex          sync.Mutex
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
	newnode := &NodeS{data: val}
	if stack.head == nil {
		stack.head = newnode
		stack.head.data = val
	} else {
		newnode.next = stack.head
		stack.head = newnode
		stack.head.data = val
	}
}

//queue

func (queue *Queue) pushQ(val string) {
	newnode := &NodeQ{data: val}
	if queue.head == nil {
		queue.head = newnode
		queue.tail = newnode
	} else {
		queue.tail.next = newnode
		queue.tail = newnode
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
	index := ht.hashFunc(1, key)
	node := &NodeH{Key: key, Value: value}
	if ht.Table[index] == nil {
		ht.Table[index] = node
	} else {
		for i := 2; i < 32; i++ {
			index := ht.hashFunc(i, key)
			if ht.Table[index] == nil {
				ht.Table[index] = node
				break
			}
		}
	}
}

func (ht *HashTable) Get(key string) (string, bool) {
	index := ht.hashFunc(1, key)
	if ht.Table[index] != nil && ht.Table[index].Key == key {
		return ht.Table[index].Value, true
	}
	return "", false
}

func (ht *HashTable) Delete(key string) bool {
	index := ht.hashFunc(1, key)
	if ht.Table[index] != nil && ht.Table[index].Key == key {
		ht.Table[index].Key = "0"
		ht.Table[index].Value = "0"
		return true
	}
	return false
}

//set

func (st *Set) AddS(key, value string) {
	index := st.hashfuncSet(1, key)
	node := &NodeSH{Key: key, Value: value}
	if st.Table[index] == nil {
		st.Table[index] = node
	} else {
		for i := 2; i < 32; i++ {
			index := st.hashfuncSet(i, key)
			if st.Table[index] == nil {
				st.Table[index] = node
				break
			}
		}
	}
}

func (st *Set) GetS(key string) (string, bool) {
	index := st.hashfuncSet(1, key)
	if st.Table[index] != nil && st.Table[index].Key == key {
		return st.Table[index].Value, true
	}
	return "", false
}

func (st *Set) DeleteS(key string) bool {
	index := st.hashfuncSet(1, key)
	if st.Table[index] != nil && st.Table[index].Key == key {
		st.Table[index].Key = "0"
		st.Table[index].Value = "0"
		return true
	}
	return false
}

//func saveStructure(filePath string, structure DataStructure) {
// Кодируем структуру в формат JSON
//	jsonData, err := json.MarshalIndent(structure, "", "  ")
//	if err != nil {
//		return
//	}
// Записываем данные в файл
//	err = ioutil.WriteFile(filePath, jsonData, 0644)
//	if err != nil {
//		return
//	}
//}

//------------------------------------

var db MainDataStructure

func maincode(conn net.Conn) {
	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	line := string(buffer[:n])
	str := strings.ReplaceAll(line, "\"", "")
	str = strings.ReplaceAll(str, "'", "")

	if str == "quit" {
		break
	}

	var key string
	var val string

	elements := strings.Split(str, " ")
	filePtr := strings.TrimSpace(elements[1])
	command := strings.TrimSpace(elements[3])
	structName := strings.TrimSpace(elements[4])

	fmt.Println(structName, len(structName))

	if len(elements) == 7 {
		key = strings.TrimSpace(elements[5])
		val = strings.TrimSpace(elements[6])
	} else if len(elements) == 6 {
		key = strings.TrimSpace(elements[5])
		val = strings.TrimSpace(elements[5])
	}

	fmt.Println("Key - ", key, len(key))
	fmt.Println("Val - ", val, len(val))
	fmt.Println("Command - ", command, len(command))

	found := 0
	index := 0

	for i := range db.datastructures {
		if db.datastructures[i].name == filePtr {
			index = i
			found = 1
		}
	}
	if found == 0 {
		newStruct := DataStructure{name: filePtr}
		db.datastructures = append(db.datastructures, newStruct)
		index = len(db.datastructures) - 1
	}

	// Проверяем, существует ли файл
	//	if _, err := os.Stat(filePtr); os.IsNotExist(err) {
	// Файл не существует, создаем его
	//		file, err := os.Create(filePtr)
	//		if err != nil {
	//			return
	//		}
	//		file.Close()
	//		database = DataStructure{}
	//		saveStructure(filePtr, database)
	//	} else {
	// Файл существует, открываем его и читаем содержимое
	//		content, err := ioutil.ReadFile(filePtr)
	//		err = json.Unmarshal(content, &database)
	//		if err != nil {
	//			return
	//		}
	//	}

	// Обработка каждой входной команд

	db.mutex.Lock()

	if command == "SPUSH" {
		found := 0
		for i := range db.datastructures[index].stacks {
			if db.datastructures[index].stacks[i].Name == structName {
				db.datastructures[index].stacks[i].push(val)
				found = 1
			}

		}
		if found == 0 {
			newStack := Stack{Name: structName}
			newStack.push(val)
			db.datastructures[index].stacks = append(db.datastructures[index].stacks, newStack)
		}
	} else if command == "SPOP" {
		found := 0
		for i := range db.datastructures[index].stacks {
			if db.datastructures[index].stacks[i].Name == structName {
				outputString, error := db.datastructures[index].stacks[i].pop()
				if error != nil {
					fmt.Println(error)
				} else {
					found = 1
					response := []byte(outputString + "\n")
					_, err = conn.Write(response)
					if err != nil {
						fmt.Println("Error writing:", err.Error())
						return
					}
				}
			}

		}
		if found == 0 {
			fmt.Println("Stack does not exist")
		}
	} else if command == "QPUSH" {
		found := 0
		for i := range db.datastructures[index].queues {
			if db.datastructures[index].queues[i].Name == structName {
				db.datastructures[index].queues[i].pushQ(val)
				found = 1
			}

		}
		if found == 0 {
			newQueue := Queue{Name: structName}
			newQueue.pushQ(val)
			db.datastructures[index].queues = append(db.datastructures[index].queues, newQueue)
		}
	} else if command == "QPOP" {
		found := 0
		for i := range db.datastructures[index].queues {
			if db.datastructures[index].queues[i].Name == structName {
				outputString, error := db.datastructures[index].queues[i].popQ()
				if error != nil {
					fmt.Println(error)
				} else {
					found = 1
					response := []byte(outputString + "\n")
					_, err = conn.Write(response)
					if err != nil {
						fmt.Println("Error writing:", err.Error())
						return
					}
				}
			}

		}
		if found == 0 {
			fmt.Println("Stack does not exist")
		}
	} else if command == "HSET" {
		found := 0
		for i := range db.datastructures[index].hashTables {
			if db.datastructures[index].hashTables[i].Name == structName {
				db.datastructures[index].hashTables[i].Add(key, val)
				found = 1
			}
		}
		if found == 0 {
			newTable := NewHashTable(512, structName)
			newTable.Add(key, val)
			db.datastructures[index].hashTables = append(db.datastructures[index].hashTables, *newTable)
		}
	} else if command == "HDEL" {
		found := 0
		for i := range db.datastructures[index].hashTables {
			if db.datastructures[index].hashTables[i].Name == structName {
				outputString := db.datastructures[index].hashTables[i].Delete(key)
				out := strconv.FormatBool(outputString)
				response := []byte(out + "\n")
				_, err = conn.Write(response)
				if err != nil {
					fmt.Println("Error writing:", err.Error())
					return
				}
				found = 1
			}
		}
		if found == 0 {
			response := []byte("Stack does not exist")
			_, err = conn.Write(response)
			if err != nil {
				fmt.Println("Error writing:", err.Error())
				return
			}
		}
	} else if command == "HGET" {
		found := 0
		for i := range db.datastructures[index].hashTables {
			if db.datastructures[index].hashTables[i].Name == structName {
				outputString, error := db.datastructures[index].hashTables[i].Get(key)
				if error == false {
					fmt.Println(error)
				}
				found = 1

				response := []byte(outputString + "\n")
				_, err = conn.Write(response)
				if err != nil {
					fmt.Println("Error writing:", err.Error())
					return
				}
			}
		}
		if found == 0 {
			fmt.Println("Hashtable does not exist")
		}
	} else if command == "SADD" {
		found := 0
		for i := range db.datastructures[index].sets {
			if db.datastructures[index].sets[i].Name == structName {
				db.datastructures[index].sets[i].AddS(key, "1")
				found = 1
			}
		}
		if found == 0 {
			newSet := newSet(512, structName)
			newSet.Name = structName
			newSet.AddS(key, "1")
			db.datastructures[index].sets = append(db.datastructures[index].sets, *newSet)
		}
	} else if command == "SREM" {
		found := 0
		for i := range db.datastructures[index].sets {
			if db.datastructures[index].sets[i].Name == structName {
				outputString := db.datastructures[index].sets[i].DeleteS(key)
				out := strconv.FormatBool(outputString)
				response := []byte(out + "\n")
				_, err = conn.Write(response)
				if err != nil {
					fmt.Println("Error writing:", err.Error())
					return
				}
				found = 1
			}
		}
		if found == 0 {
			response := []byte("Stack does not exist")
			_, err = conn.Write(response)
			if err != nil {
				fmt.Println("Error writing:", err.Error())
				return
			}
		}
	} else if command == "SISMEMBER" {
		found := 0
		for i := range db.datastructures[index].sets {
			if db.datastructures[index].sets[i].Name == structName {
				outputString, error := db.datastructures[index].sets[i].GetS(key)
				if error == false {
					fmt.Println(error)
				}
				found = 1

				response := []byte(outputString + "\n")
				_, err = conn.Write(response)
				if err != nil {
					fmt.Println("Error writing:", err.Error())
					return
				}
			}
		}
		if found == 0 {
			fmt.Println("Hashtable does not exist")
		}
	}

	db.mutex.Unlock()

	fmt.Println(db.datastructures[index])
	//fmt.Println(db.datastructures[index])

	// Преобразовываем структуру в байтовый массив
	//	jsonData, err := json.MarshalIndent(database, "", "  ")
	//	if err != nil {
	//		fmt.Println("Ошибка при маршалинге структуры в JSON:", err)
	//		return
	//	}

	// Записываем байтовый массив в файл
	//	err = ioutil.WriteFile(filePtr, jsonData, 0644)
	//	if err != nil {
	//		fmt.Println("Ошибка при записи данных в файл:", err)
	//		return
	//	}
}

func main() {
	fmt.Println("Запуск сервера на порту 6379...")

	// Слушаем порт 6379
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err.Error())
		return
	}
	defer listener.Close()

	// Принимаем новые соединения и запускаем их обработку в отдельных горутинах
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err.Error())
			continue
		}
		go maincode(conn)
	}
}
