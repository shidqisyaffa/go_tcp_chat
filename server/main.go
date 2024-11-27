package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	clients   = make(map[net.Conn]string)
	clientsMu sync.Mutex
)

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	conn.Write([]byte("Masukkan nama Anda: "))
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	if name == "" || strings.HasPrefix(name, "GET") {
		fmt.Println("Koneksi tidak valid, menutup.")
		return
	}

	clientsMu.Lock()
	clients[conn] = name
	clientsMu.Unlock()

	fmt.Printf("Klien %s terhubung.\n", name)
	broadcast(fmt.Sprintf("%s bergabung dalam obrolan.\n", name), conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Klien %s terputus. Error: %v\n", name, err)
			break
		}

		message = strings.TrimSpace(message)
		fmt.Printf("Pesan diterima dari %s: %s\n", name, message)
		broadcast(fmt.Sprintf("%s: %s\n", name, message), conn)
	}

	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
	broadcast(fmt.Sprintf("%s keluar dari obrolan.\n", name), nil)
}

func broadcast(message string, sender net.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn := range clients {
		if conn != sender {
			conn.Write([]byte(message + "\n"))
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:playit")
	if err != nil {
		fmt.Println("Gagal memulai server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server berjalan di port 9000...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Gagal menerima koneksi:", err)
			continue
		}
		go handleClient(conn)
	}
}
