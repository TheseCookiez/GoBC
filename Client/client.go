package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

type Information struct {
	Vehicle_ID      string `json:"Vehicle_ID"`    // Vehicle ID
	Vehicle_Manu    string `json:"Vehicle_Manu"`  // Manufacturer
	Vehicle_Model   string `json:"Vehicle_Model"` // Model
	Vehicle_State   string `json:"Vehicle_State"` // State (Active, Inactive, etc.) / (Driving, Parked, etc.)
	Owner_ID        string `json:"Owner_ID"`      // Owner ID
	Passenger_Count string `json:"Occupancy"`     // Occupancy
	Region          string `json:"Region"`        // Region
	X_pos           string `json:"X_pos"`         // X Position
	Y_pos           string `json:"Y_pos"`         // Y Position
	V_Speed         string `json:"V_Speed"`       // Vehicle Speed (km/h)
}

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Enter a message to send to the server:")
	for {
		time.Sleep(1 * time.Second)
		proposal := Information{}
		proposal.Vehicle_ID = fmt.Sprint(rand.Intn(100))
		proposal.Vehicle_Manu = "Ford"
		proposal.Vehicle_Model = "F150"
		proposal.Vehicle_State = "Active"
		proposal.Owner_ID = fmt.Sprint(rand.Intn(100))
		proposal.Passenger_Count = fmt.Sprint(rand.Intn(100))
		proposal.Region = "North-America"
		proposal.X_pos = fmt.Sprint(rand.Intn(100))
		proposal.Y_pos = fmt.Sprint(rand.Intn(100))
		proposal.V_Speed = fmt.Sprint(rand.Intn(100))

		senddata := strings.Trim(fmt.Sprint(proposal), "{}")
		_, err = conn.Write([]byte(fmt.Sprint(senddata)))

		/*
			This is to send stdin instead of random data

			for {
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				_, err = conn.Write([]byte(fmt.Sprintf("%s\r\n", text))) // Send AV struct here instead of string
				if err != nil {
					println("Write failed:", err.Error())
					os.Exit(1)
				}
			}
		*/

		// Doesn't matter at this point
		defer conn.Close()

	}
}
