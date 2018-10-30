package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

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

func main() {

	//testValue1 := Person{"Thomas", 5983, 37}
	//	testValue2 := Person{"Eliott", 4798, 24}
	//testValue3 := Person{"Alice", 2346, 28}

	//msgAdd := []byte("ADD")
	msgGet := []byte("GET")

	//serializedTestValue1, _ := Serialize(&testValue1)
	//serializedTestValue2, _ := Serialize(&testValue2)
	//serializedTestValue3, _ := Serialize(&testValue3)

	conn, err := net.Dial("tcp", "localhost:6666")
	if err != nil {
		fmt.Println(err)
	}
	/*
		packet := append(msgAdd, serializedTestValue1...)
		fmt.Println(string(packet))

		_, err = conn.Write(packet)
		if err != nil {
			fmt.Println(err)
			return
		}
		resp := make([]byte, 64)
		bufio.NewReader(conn).Read(resp)
		respout, _ := Deserialize(resp)
		fmt.Println(string(resp))
		fmt.Println(respout)
		conn.Close()

		conn, _ = net.Dial("tcp", "localhost:6666")
	*/
	getpacket := append(msgGet, []byte("Thomas")...)

	conn.Write(getpacket)

	getresp := make([]byte, 64)
	bufio.NewReader(conn).Read(getresp)
	degetresp, _ := Deserialize(getresp)
	fmt.Println(degetresp)

	/*
		packet2 := append(msgAdd, serializedTestValue2...)
		conn.Write(packet2)
		resp2 := make([]byte, 64)
		bufio.NewReader(conn).Read(resp2)
		fmt.Println(string(resp2))

		packet3 := append(msgAdd, serializedTestValue3...)
		conn.Write(packet3)
		resp3 := make([]byte, 64)
		bufio.NewReader(conn).Read(resp3)
		fmt.Println(string(resp3))
	*/
}
