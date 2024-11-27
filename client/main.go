package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)


func main() {
	serverAddr := "your ip playit"
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Gagal terhubung ke server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Terhubung ke server.")
	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	fmt.Print("Masukkan nama Anda: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	conn.Write([]byte(name + "\n"))

	// Goroutine untuk menerima pesan dari server
	go func() {
		for {
			message, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("Terputus dari server.")
				os.Exit(0)
			}
			fmt.Print(message)
		}
	}()


	for {
		fmt.Print("Ketik pesan: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		conn.Write([]byte(message + "\n"))
	}
}
