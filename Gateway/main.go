package main

import "gitlab.com/IIIS/backend/gateway/cmd"

func main() {
	cmd.Execute()
}

/*
import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	//"time"
)

const (
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "8070"
	CONN_TYPE = "tcp"
	CONN_BASE = "http://localhost:4040/gateway/access"
)

type meshData struct {
	addr string `json:"addr"`
	data string `json:"data"`
}

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	dlen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// `&myStoredVariable` is the address of the variable we want to store our
	// parsed data in

	var data map[string]interface{}

	if err := json.Unmarshal(buf[:dlen], &data); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", data)
	fmt.Printf("%s\n", buf[:dlen])

	if _, ok := data["addr"]; !ok {
		return
	}

	if _, ok := data["data"]; !ok {
		return
	}

	//err = json.Unmarshal(buf[:dlen], &data)
	enc := base64.StdEncoding.EncodeToString(buf[:dlen])

	go handlePost([]byte(enc))
	fmt.Println(string(enc))
	//fmt.Printf("%+v\n", data.addr)
	//time.Sleep(2 * time.Second)
	// Send a response back to person contacting us.
	conn.Write(buf[:dlen])
	// Close the connection when you're done with it.
	conn.Close()
}

func handlePost(data []byte) {
	//req, err := http.NewRequest("POST", CONN_BASE, bytes.NewBuffer(data))
	//req.Header.Set("X-Custom-Header", "myvalue")
	//req.Header.Set("Content-Type", "application/json")

	resp, err := http.Post(CONN_BASE, "base64", bytes.NewBuffer(data))

	if err != nil {
		fmt.Println(err)
		//panic(err)
	} else {

		defer resp.Body.Close()
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
}

*/
