package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func unmarshal(data []byte) interface{} {
	dataItems := strings.Split(string(data), "\r\n")
	fmt.Println(dataItems)
	var res interface{}
	if dataItems[0][0] == '+' {
		return dataItems[0][:len(dataItems[0])-1]
	} else if dataItems[0][0] == '$' { //todo check the size to determine empty or null value
		return dataItems[0][1]
	}
	res = []string{}
	for i := 0; i < len(dataItems); i++ {
		if dataItems[i] == "" || len(dataItems[i]) > 0 && (dataItems[i][0] == '$' || dataItems[i][0] == '*') {
			continue
		}
		res = append(res.([]string), dataItems[i])
	}

	return res

}
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
		fmt.Println("Connection established", conn)
	}
}
func wordWriter(w string) string {
	return w + "\r\n"
}
func marshal(d interface{}) []byte {
	switch v := d.(type) {
	case []string:
		//todo: implement marshal

		fmt.Println("int:", v)
	case string:
		return []byte("+" + wordWriter(d.(string)))
	default:
		fmt.Println("unknown")
	}
	return []byte{}
}

func commander(v string, in interface{}) interface{} {
	if strings.ToLower(v) == "ping" {
		ret := strings.ToUpper("pong")
		return ret
	} else if strings.ToLower(v) == "echo" {
		return in.([]string)[0]
	}
	return "unsupported"
}
func handleCommands(command interface{}) []byte {
	var res interface{}
	v, k := command.(string)
	switch {
	case k:
		res = commander(v, nil)
	case len(command.([]string)) == 1:
		d, _ := command.([]string)
		res = commander(d[0], nil)
	case len(command.([]string)) > 1:
		d, _ := command.([]string)
		res = commander(d[0], d[1:])
	default:
		res = "notsupported"
	}
	return marshal(res)
}

func handleConnection(conn net.Conn) {
	for conn != nil {
		buf := make([]byte, 1024)
		len, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Error reading: %#v\n", err)
			return
		}
		fmt.Printf("Message received: %s\n", string(buf[:len]))
		fmt.Printf("Message received from the unmarshal: %s\n", unmarshal(buf[:len]))

		conn.Write(handleCommands(unmarshal(buf[:len])))
	}

}
