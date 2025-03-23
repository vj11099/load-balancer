package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Println("")
	fmt.Println("X-Real-IP: ", req.Header.Get("X-Real-IP"))
	fmt.Println("X-Forwarded-For: ", req.Header.Get("X-Forwarded-For"))
	fmt.Println("RemoteAddr: ", req.RemoteAddr)
	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Port: ", port)
	fmt.Println("Host: ", req.Host)
	fmt.Println("UserAgent: ", req.UserAgent())
	userIP := net.ParseIP(ip)
	fmt.Println("userIP: ", userIP)
	time.Sleep(time.Second * 3)
	fmt.Fprintf(w, "hello")
}

func main() {
	http.HandleFunc("/", hello)

	err := godotenv.Load()
	if err != nil {
		panic("Error loading the .env file")
	}

	ports := strings.Fields(os.Getenv("PORTS"))

	// replicate the behaviour of servers running on multiple ports
	for _, i := range ports {
		fmt.Println("Listening at port ", i)
		go http.ListenAndServe(":"+i, nil)
	}

	select {}
}
