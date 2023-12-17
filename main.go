package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Hash_Table struct {
	Table []*Doubly_Connected_Table
	Size  int
	Cout  int
}

type Pair struct {
	Key   string
	Value string
}

type Node_Table struct {
	Data     Pair
	Next     *Node_Table
	Previous *Node_Table
}

type Doubly_Connected_Table struct {
	Lenght int
	Head   *Node_Table
	Tail   *Node_Table
}

func NewHashTable(size int) *Hash_Table {
	table := make([]*Doubly_Connected_Table, size)
	for i := range table {
		table[i] = &Doubly_Connected_Table{}
	}
	return &Hash_Table{
		Table: table,
		Size:  size,
		Cout:  0,
	}
}

func (ht *Hash_Table) Hash(key string) int {
	key_int := 0
	for _, symbol := range key {
		key_int += int(symbol)
	}
	return key_int % ht.Size
}

func (ht *Hash_Table) Hset(conn net.Conn, key string, value string) {
	val := &Pair{Key: key, Value: value}
	hash := ht.Hash(val.Key)
	if ht.Table[hash].Lenght < 20 {
		ht.Table[hash].dpush(conn, *val)
		ht.Cout++
		return
	} else {
		oldSize := ht.Size
		newHT := NewHashTable(oldSize * 2)
		for i := 0; i < oldSize; i++ {
			currentNode := ht.Table[i].Head
			for currentNode != nil {
				new_hash := newHT.Hash(currentNode.Data.Key)
				newHT.Table[new_hash].dpush(conn, currentNode.Data)
				ht.Cout++
				currentNode = currentNode.Next
			}
		}
		*ht = *newHT
		new_hash := ht.Hash(val.Key)
		ht.Table[new_hash].dpush(conn, *val)
		ht.Cout++
		return
	}
}

func (ht *Hash_Table) Hget(conn net.Conn, key string) {
	hash := ht.Hash(key)
	currentNode := ht.Table[hash].Head
	for currentNode != nil {
		if currentNode.Data.Key == key {
			conn.Write([]byte(currentNode.Data.Value))
			return
		}
		currentNode = currentNode.Next
	}
	conn.Write([]byte("error"))
}

func (pair *Pair) ppush(key string, value string) {
	pair.Key = key
	pair.Value = value
}

func (doubly_connected *Doubly_Connected_Table) dpush(conn net.Conn, val Pair) {
	node_hesh := &Node_Table{Data: val}
	if doubly_connected.Head == nil {
		doubly_connected.Head = node_hesh
		doubly_connected.Tail = node_hesh
	} else {
		currentNode := doubly_connected.Head
		for currentNode != nil {
			if currentNode.Data.Key == val.Key {
				conn.Write([]byte("error"))
			}
			currentNode = currentNode.Next
		}
		doubly_connected.Tail.Next = node_hesh
		node_hesh.Previous = doubly_connected.Tail
		doubly_connected.Tail = node_hesh
	}
	doubly_connected.Lenght++
}

func hash(key string) int {
	int_key := 0                    // переменная для хранения суммы кодов символов ключа
	for i := 0; i < len(key); i++ { // проходим по всем символам ключа
		int_key += int(key[i]) // добавляем код символа к сумме
	}
	return int_key % 512 // возвращаем остаток от деления суммы на размер массива таблицы
}

func Gorutina(mutex *sync.Mutex, conn net.Conn, hTable *Hash_Table) {
	for {
		defer conn.Close()

		reader := bufio.NewReader(conn)
		var err error

		commands := [2]string{"", ""}
		for i := 0; i < 2; i++ {
			commands[i], err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading data:", err)
				return
			}
			commands[i] = strings.ReplaceAll(commands[i], "\n", "")
		}

		switch commands[0] {
		case "post":
			if commands[1] == "" {
				conn.Write([]byte("error"))
				return
			}
			mutex.Lock()
			base := "http://127.0.0.1:8080/" + strconv.Itoa(hTable.Cout+1)
			hTable.Hset(conn, base, commands[1])
			mutex.Unlock()
			conn.Write([]byte(base))

		case "post4":
			if commands[1] == "" {
				conn.Write([]byte("error"))
				return
			}
			mutex.Lock()
			base := strconv.Itoa(hTable.Cout+1)
			hTable.Hset(conn, base, commands[1])
			mutex.Unlock()
			conn.Write([]byte(base))

		case "get":
			if commands[1] == "" {
				conn.Write([]byte("error"))
				return
			} else {
				mutex.Lock()
				hTable.Hget(conn, commands[1])
				mutex.Unlock()
			}
		default:
			conn.Write([]byte("error"))
		}
	}
}

func main() {
	var mutex sync.Mutex
	hTable := NewHashTable(512)
	address := "127.0.0.1:6379"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Сервер запущен и готов к приему подключений")
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println("Error accepting connection:", err)
					continue
				}
				go Gorutina(&mutex, conn, hTable)
			}
		}()
	}
	wg.Wait()
}
