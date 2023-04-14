package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

var proposal = Information{}

func load_blockchain_json(config string) bool {
	if _, err := os.Stat(config + ".json"); err == nil {
		// File exists
		// Open the JSON file
		jsonFile, err := os.Open(config + ".json")
		if err != nil {
			fmt.Println(err)
		}
		// Read the JSON file
		byteValue, _ := ioutil.ReadAll(jsonFile)
		// Unmarshal the JSON file into the blockchain
		json.Unmarshal(byteValue, &proposal)
		// Close the JSON file
		jsonFile.Close()
		return true
	} else {
		// File does not exist
		return false
	}
}

// **DEBUG OPTION** Create the JSON file for new vehicles
func create_config_json(config_name string) {
	fmt.Println("Creating " + config_name + ".json...")
	proposal := Information{}
	proposal.Vehicle_ID = fmt.Sprint(rand.Intn(100))
	proposal.Vehicle_Manu = "Volkswagen"
	proposal.Vehicle_Model = "Golf"
	proposal.Vehicle_State = "Active"
	proposal.Owner_ID = fmt.Sprint(rand.Intn(100))
	proposal.Passenger_Count = fmt.Sprint(rand.Intn(5))
	proposal.Region = "North-America"
	proposal.X_pos = fmt.Sprint(rand.Intn(100))
	proposal.Y_pos = fmt.Sprint(rand.Intn(100))
	proposal.V_Speed = fmt.Sprint(rand.Intn(100))

	// Create the JSON file
	jsonFile, err := os.Create(config_name + ".json")
	if err != nil {
		fmt.Println(err)
	}
	// Marshal the JSON file
	jsonData, err := json.MarshalIndent(proposal, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	// Write the JSON file
	jsonFile.Write(jsonData)
	// Close the JSON file
	jsonFile.Close()
}

func main() {
	//create_config_json("Insert_Vehicle")

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Please enter a vehicle name (e.g. nissan/toyota/volkswagen/volvo)")
		os.Exit(1)
	}
	// Load the selected JSON file
	load_blockchain_json(args[0])
	// Configure the connection
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}
	// Connect to the server
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Starting to send data...")
	// Send the data in an infinite loop
	for {
		time.Sleep(1 * time.Second)
		// Generate random data for testing values expected to change
		proposal.X_pos = fmt.Sprint(rand.Intn(100))
		proposal.Y_pos = fmt.Sprint(rand.Intn(100))
		proposal.V_Speed = fmt.Sprint(rand.Intn(100))
		// Convert the struct to a string and strip the brackets
		senddata := strings.Trim(fmt.Sprint(proposal), "{}")
		// Send the data as byte array
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

		// Doesn't matter with infinite loop
		defer conn.Close()

	}
}
