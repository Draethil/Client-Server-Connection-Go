package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type selfassignmentTicket struct {
	TicketID    int    `json:"TicketID"`
	ClientID    string `json:"ClientID"`
	TicketTaken bool   `json:"ticketTaken"`
}

type client struct {
	clientID   string
	connection net.Conn
}

func main() {
	var ticketList []selfassignmentTicket
	var clientList []client
	go func() {
		listener, err := net.Listen("tcp", ":7777")
		if err != nil {
			log.Fatalln(err)
		}

		for {
			connection, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
			}

			var messageThisClientConnected bool
			messageThisClientConnected = true
			e := json.NewEncoder(connection)
			e.Encode(messageThisClientConnected)

			d := json.NewDecoder(connection)
			var tempClient client
			tempClient.connection = connection
			d.Decode(&tempClient.clientID)
			clientList = append(clientList, tempClient)

			for k := 0; k < len(clientList); k++ {
				if clientList[k].clientID == tempClient.clientID {
					e := json.NewEncoder(clientList[k].connection)
					e.Encode(ticketList)
				}
			}

			go func(tempClient client) {
				for {
					var TicketID int
					d := json.NewDecoder(tempClient.connection)
					d.Decode(&TicketID)

					if TicketID == -123456789 {
						fmt.Printf("\nClient: '%v' disconnected!\n", tempClient.clientID)
						consoleOutputTicketList(ticketList)
					}

					for k := 0; k < len(ticketList); k++ {
						if ticketList[k].TicketID == TicketID && ticketList[k].ClientID == "" {
							ticketList[k].ClientID = tempClient.clientID
							ticketList[k].TicketTaken = true

							for l := 0; l < len(clientList); l++ {
								e := json.NewEncoder(clientList[l].connection)
								e.Encode(ticketList)
							}

							consoleOutputTicketList(ticketList)
						}
					}
				}
			}(tempClient)
			fmt.Printf("\nNew Client connected: %v\n", tempClient.clientID)
			consoleOutputTicketList(ticketList)
		}
	}()

	fmt.Println("Server started...\nNo Tickets available")
	fmt.Println("n: new ticket, q: quit")
	var ticketCounter int
	ticketCounter = 1
	for {
		var input string
		fmt.Scan(&input)

		switch input {
		case "n":
			var tempTicket selfassignmentTicket
			tempTicket.TicketTaken = false
			tempTicket.TicketID = ticketCounter
			ticketCounter++
			ticketList = append(ticketList, tempTicket)
			for k := 0; k < len(clientList); k++ {
				e := json.NewEncoder(clientList[k].connection)
				e.Encode(ticketList)
			}
			consoleOutputTicketList(ticketList)
		case "q":
			fmt.Println("\nServer closed...")
			return
		default:
			fmt.Println("Wrong Input")
		}
	}
}

func consoleOutputTicketList(ticketList []selfassignmentTicket) {
	var length = len(ticketList)
	if 0 >= length {
		fmt.Println("\nNo Tickets available")
	} else {
		fmt.Println("\nTickets:")
		for k := 0; k < length; k++ {
			var temp = ticketList[k]
			if temp.ClientID == "" {
				fmt.Printf("%v: ticket%v (Not assigned yet)\n", k+1, temp.TicketID)
			} else {
				fmt.Printf("%v: ticket%v (%v)\n", k+1, temp.TicketID, temp.ClientID)
			}
		}
	}
	fmt.Println("n: new ticket, q: quit")
}
