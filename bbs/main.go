package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"amis-x/atascii"
	"amis-x/bbscommon" // Adjust the import path as needed
)

const (
	ASCII_MODE = iota
	ATASCII_MODE
	ANSI_MODE
)

var (
	mode       = ASCII_MODE // Default mode
	db         *sql.DB
	config     bbscommon.Config
	configPath = "config.json"
)

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./bbs.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	bbscommon.CreateTables(db)
	config, err = bbscommon.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

func main() {
	initSystem()
	mainMenuLoop()
}

func initSystem() {
	setupModem()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection established.")

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	fmt.Fprint(writer, "Enter your Name > ")
	writer.Flush()
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if len(name) < 3 {
		fmt.Fprintln(writer, "Name too short. Disconnecting...")
		writer.Flush()
		return
	}

	userID, err := bbscommon.GetUserID(db, name)
	if err != nil {
		fmt.Fprintln(writer, "User not found, please register.")
		writer.Flush()
		registerUser(writer, reader, name)
	} else {
		fmt.Fprint(writer, "Enter your password > ")
		writer.Flush()
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)
		if !bbscommon.CheckPassword(db, userID, password) {
			fmt.Fprintln(writer, "Invalid password, try again.")
			writer.Flush()
			return
		}
	}

	fmt.Fprint(writer, "From City, State > ")
	writer.Flush()
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)
	if len(address) < 3 {
		fmt.Fprintln(writer, "Address too short. Disconnecting...")
		writer.Flush()
		return
	}

	fmt.Fprintf(writer, "You are %s, calling from %s. CORRECT (Y/N)? ", name, address)
	writer.Flush()
	confirmation, _ := reader.ReadString('\n')
	confirmation = strings.TrimSpace(confirmation)
	if !strings.EqualFold(confirmation, "Y") {
		fmt.Fprintln(writer, "Details not confirmed. Disconnecting...")
		writer.Flush()
		return
	}

	saveCallDetails(userID)
	printWelcome(writer)
}

func setupModem() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", config.PortNumber))
	if err != nil {
		log.Fatalf("Failed to set up modem: %v", err)
	}
	defer ln.Close()

	fmt.Println("Listening for incoming connections...")

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("Failed to accept connection: %v", err)
				continue
			}
			go handleConnection(conn)
		}
	}()
}

func mainMenuLoop() {
	var userName string

	for {
		fmt.Print("Command?>")
		command := readChar()
		switch command {
		case 'A':
			postMessage(userName)
		case 'B':
			viewMessages()
		case 'W':
			DisplayAtasciiFile("welcome.ata")
		case 'T':
			toggleMode()
		case 'Q':
			logOff()
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}

func readChar() rune {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Fatalf("Failed to read character: %v", err)
	}
	return char
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read line: %v", err)
	}
	return strings.TrimSpace(line)
}

func yesNo() bool {
	answer := readChar()
	return answer == 'Y' || answer == 'y'
}

func registerUser(writer *bufio.Writer, reader *bufio.Reader, name string) {
	fmt.Fprint(writer, "Enter a password > ")
	writer.Flush()
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Fprint(writer, "From City, State > ")
	writer.Flush()
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)

	err := bbscommon.RegisterUser(db, name, password, address)
	if err != nil {
		fmt.Fprintf(writer, "Failed to register user: %v\n", err)
		writer.Flush()
		return
	}
	fmt.Fprintln(writer, "Registration successful.")
	writer.Flush()
}

func saveCallDetails(userID int) {
	query := "INSERT INTO messages (user_id, content) VALUES (?, ?)"
	_, err := db.Exec(query, userID, "User logged in")
	if err != nil {
		log.Printf("Failed to save call details: %v", err)
	}
}

func printWelcome(writer *bufio.Writer) {
	fmt.Fprintln(writer, "Welcome to the BBS!")
	writer.Flush()
	DisplayFile("welcome.ata")
	fmt.Fprintln(writer, "Enjoy your stay.")
	writer.Flush()
}

func postMessage(userName string) {
	fmt.Print("Enter the message base name: ")
	messageBaseName := readLine()

	messageBaseID, err := bbscommon.GetMessageBaseID(db, messageBaseName)
	if err != nil {
		fmt.Println("Message base not found.")
		return
	}

	fmt.Print("Enter your message: ")
	message := readLine()

	userID, err := bbscommon.GetUserID(db, userName)
	if err != nil {
		log.Fatalf("Failed to get user ID: %v", err)
	}

	query := "INSERT INTO messages (user_id, message_base_id, content) VALUES (?, ?, ?)"
	_, err = db.Exec(query, userID, messageBaseID, message)
	if err != nil {
		log.Fatalf("Failed to post message: %v", err)
	}
	fmt.Println("Message posted.")
}

func viewMessages() {
	fmt.Print("Enter the message base name: ")
	messageBaseName := readLine()

	messageBaseID, err := bbscommon.GetMessageBaseID(db, messageBaseName)
	if err != nil {
		fmt.Println("Message base not found.")
		return
	}

	messages, err := bbscommon.ViewMessages(db, messageBaseID)
	if err != nil {
		log.Fatalf("Failed to retrieve messages: %v", err)
	}

	for _, msg := range messages {
		fmt.Println(msg)
	}
}

func logOff() {
	fmt.Println("Any Comments? (Y/N)")
	if yesNo() {
		fmt.Println("Enter comments")
		comment := readLine()
		fmt.Printf("Comment saved: %s\n", comment)
	}
	fmt.Println("Thanks for calling, please call again...")
}

func toggleMode() {
	mode = (mode + 1) % 3
	switch mode {
	case ASCII_MODE:
		fmt.Println("Switched to ASCII mode")
	case ATASCII_MODE:
		fmt.Println("Switched to ATASCII mode")
	case ANSI_MODE:
		fmt.Println("Switched to ANSI mode")
	}
}

func DisplayFile(filename string) {
	switch mode {
	case ASCII_MODE:
		data, err := os.ReadFile(filename)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		fmt.Println(string(data))
	case ATASCII_MODE:
		DisplayAtasciiFile(filename)
	case ANSI_MODE:
		// Implement ANSI file handling if needed
	}
}

func TranslateAtasciiToUnicode(data []byte) string {
	var result []rune
	for _, b := range data {
		if unicodeChar, exists := atascii.AtasciiToUnicode[b]; exists {
			result = append(result, unicodeChar)
		} else {
			result = append(result, '?')
		}
	}
	return string(result)
}

func DisplayAtasciiFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	translated := TranslateAtasciiToUnicode(data)
	fmt.Println(translated)
}
