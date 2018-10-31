package main

// Building a concurrent key-value store server

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// Connection Handlers

func handleConnection(conn net.Conn, db *DB) {

	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from " + remoteAddr)

	scanner := bufio.NewReader(conn)

	v := make([]byte, 128)
	for {

		_, err := scanner.Read(v)

		if err != nil {
			break
		}
		fmt.Println(len(v))
		handleMessage(db, v, conn)

	}

	fmt.Println("Client " + remoteAddr + "Disconnected")
}

func handleMessage(db *DB, msg []byte, conn net.Conn) {
	// handle Message takes the message scanned from the connection
	// and the connection so it can reply back

	switch {

	case string(msg[:3]) == "ADD":
		fmt.Println(">" + string(msg[:3]))
		fmt.Println(len(msg))
		go addHandler(db, msg)
		response := "SUCCESSFULLY ADDED"
		fmt.Fprintf(conn, response+"\n")
		fmt.Println("Done")
	case string(msg[:3]) == "DEL":
		fmt.Println(">" + string(msg[:2]))
		response := "Timestamp" + time.Now().String() + ": Successfully Deleted from store" + "\n"
		conn.Write([]byte(response))
	case string(msg[:3]) == "GET":
		fmt.Println(">" + string(msg[:3]))
		response, _ := getHandler(db, msg)
		responseStr, _ := Deserialize(response)
		fmt.Println(responseStr)
		conn.Write(response)
	default:
		conn.Write([]byte("Unrecognized Command"))

	}
}

// addHandler handles add messages
func addHandler(db *DB, msg []byte) {
	// MSG Format
	// First 3 bytes are the message
	// Rest n bytes are the Person struct
	fmt.Println(len(msg))
	serialized := msg[3:]
	//	serialized = append(serialized, msg[len(msg)-1])

	fmt.Println("ADD HANDLER")
	// deserialize the person to a struct
	//p := new(Person)

	p, err := Deserialize(serialized)
	if err != nil {
		fmt.Println("FAILED TO DESERIALIZE", err)
		return
	}
	fmt.Println("ADDED ", p)
	db.Put(p.Name, *p)
	return
}

// getHandler handles get messages
func getHandler(db *DB, msg []byte) ([]byte, error) {
	// skip first 3 bytes expect the next n bytes to be a key
	key := msg[3:]

	p, err := db.Get(string(key))

	if err != nil {
		return nil, err
	}
	return Serialize(&p)
}

// Person is a complex data
type Person struct {
	Name string
	ID   int
	Age  int
}

// Serialize turnes struct to bytes
func Serialize(p *Person) ([]byte, error) {
	var b bytes.Buffer

	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(p)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Deserialize turnes bytes to struct
func Deserialize(b []byte) (*Person, error) {

	p := new(Person)

	s := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(s)

	err := decoder.Decode(p)

	if err != nil {
		return &Person{}, err
	}

	return p, nil

}

// DB is an inmemeory threadsafe golang map
type DB struct {
	Lock  sync.RWMutex
	Store map[string]Person
}

// Get returns a value by key
func (db *DB) Get(key string) (Person, error) {

	var value Person
	// lock the mutex for read
	db.Lock.RLock()
	value = db.Store[key]
	// unlock the mutex
	db.Lock.RUnlock()
	// maps returns zero value for non present keys
	// zero value of an empty struct is the empty struct
	if value.Name != key {
		return Person{}, errors.New("Key Not Present")
	}

	return value, nil

}

// Put Stores a new value with a key
func (db *DB) Put(key string, value Person) {

	// Lock the mutex with the write-lock
	db.Lock.Lock()
	db.Store[key] = value
	// Unlock the mutex
	db.Lock.Unlock()

}

// Del deletes a key-value pair
func (db *DB) Del(key string) {
	db.Lock.Lock()
	delete(db.Store, key)
	db.Lock.Unlock()
}

// TCPServe starts the tcp server
func TCPServe(db *DB) {
	fmt.Println("Starting TCP Server")
	src := "localhost:6666"

	listener, _ := net.Listen("tcp", src)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error")
		}

		go handleConnection(conn, db)
	}
}
func main() {

	james := &Person{"James", 1234, 28}

	jamesAsBytes, err := Serialize(james)

	if err != nil {
		fmt.Printf("Failed to Serialize %v with err %s\n", james, err)
	}

	fmt.Println("James was successfully cerealized (turned to cereal)")

	jamesAsPerson, err := Deserialize(jamesAsBytes)

	if err != nil {
		fmt.Printf("Failed to Uncerealize james he's forever turned to cereals with error %s\n", err)
	}

	fmt.Printf("James was uncerealized to human successfully: %v\n", jamesAsPerson)
	database := new(DB)
	database.Store = make(map[string]Person, 10)
	TCPServe(database)
}
