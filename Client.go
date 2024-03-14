package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

type selfassignmentTicketC struct {
	TicketID    int    `json:"TicketID"`
	ClientID    string `json:"ClientID"`
	TicketTaken bool   `json:"ticketTaken"`
}

func main() {
	var tempClientID string
	fmt.Print("Client started...\nPlease enter ClientID: ")
	fmt.Scan(&tempClientID)
	fmt.Printf("\nYour ClientID: %v\n", tempClientID)

	var connection net.Conn
	con, _ := net.Dial("tcp", ":7777")
	connection = con

	var isClientConnected bool
	isClientConnected = false
	d := json.NewDecoder(connection)
	d.Decode(&isClientConnected)
	if isClientConnected {
		fmt.Println("You are now connected with the Server!")
	}

	var ticketListMainFunc []selfassignmentTicketC
	var endProgram bool
	endProgram = false

	go func(tempClientID string) {
		var input string
		for {
			fmt.Scanln(&input)
			switch input {
			case "q":
				fmt.Println("\nClient closed...")
				e := json.NewEncoder(connection)
				e.Encode(-123456789)
				connection.Close()
				endProgram = true
			default:
				selfTicketRequest, _ := strconv.Atoi(input)

				var ticketList []selfassignmentTicketC
				ticketList = ticketListMainFunc

				if selfTicketRequest > 0 && selfTicketRequest <= len(ticketList) {
					if ticketList[selfTicketRequest-1].TicketTaken == false {
						e := json.NewEncoder(connection)
						e.Encode(selfTicketRequest)
					} else {
						if ticketList[selfTicketRequest-1].ClientID == tempClientID {
							fmt.Println("It´s already your Ticket")
						} else {
							fmt.Println("Sorry, This Ticket is already taken!")
						}
					}
				} else if input != "" {
					if len(ticketList) == 0 {
						fmt.Println("No Tickets available")
					} else {
						fmt.Printf("Please enter a number, thats is between(exclusive) 0 and %v!\n", len(ticketList)+1)
					}
				}
			}
		}
	}(tempClientID)

	go func(tempClientID string) {

		e := json.NewEncoder(connection)
		e.Encode(tempClientID)

		d := json.NewDecoder(connection)
		for {
			var ticketList []selfassignmentTicketC
			d.Decode(&ticketList)
			ticketListMainFunc = ticketList
			consoleOutputTicketListC(ticketList)
		}
	}(tempClientID)

	for {
		if endProgram == true {
			return
		}
	}

}

func consoleOutputTicketListC(ticketList []selfassignmentTicketC) {
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
	fmt.Println("number: selfassignment, q: quit")
}
