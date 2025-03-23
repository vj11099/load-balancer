package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	port  string = ""
	ports []string
)

func requestHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("")
	fmt.Println("Accept:", req.Header.Get("Accept"))
	fmt.Println("X-Real-IP: ", req.Header.Get("X-Real-IP"))
	fmt.Println("X-Forwarded-For: ", req.Header.Get("X-Forwarded-For"))
	fmt.Println("RemoteAddr: ", req.RemoteAddr)

	ip, _port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Port: ", _port)
	fmt.Println("Host: ", req.Host)
	fmt.Println("UserAgent: ", req.UserAgent())
	userIP := net.ParseIP(ip)
	fmt.Println("userIP: ", userIP)

	RoundRobinLb()
	url := "http://localhost:" + port + "/"
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	io.Copy(w, response.Body)

	// resp, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("\nResponse: ", string(resp))
}

func findNextPort() {
	for i := range ports {
		if port == ports[i] {
			port = ports[i+1]
			break
		}
	}
}

func RoundRobinLb() {
	switch true {
	case port == "" || port == ports[len(ports)-1]:
		port = ports[0]
	default:
		findNextPort()
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading the .env file")
	}

	ports = strings.Fields(os.Getenv("PORTS"))

	http.HandleFunc("/", requestHandler)
	fmt.Println("Listening at port 3001")
	http.ListenAndServe(":3001", nil)
}
