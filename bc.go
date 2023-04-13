package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

/*
TODO:
- Add a function to check vehicle position and speed is within the bounds (maybe)
- Add a function to check if the vehicle changes region
- Add a function to check if the vehicle changes occupancy
- Add a function to check if the vehicle changes state
- Add a function to check if the vehicle changes static information (Manufacturer, Model, ID, etc.)
*/

var blockchain []block

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

type html_data struct {
	Title string `json:"title"`
	Data  string `json:"data"`
}

type Information struct {
	Vehicle_ID      int    `json:"Vehicle_ID"`    // Vehicle ID
	Vehicle_Manu    string `json:"Vehicle_Manu"`  // Manufacturer
	Vehicle_Model   string `json:"Vehicle_Model"` // Model
	Vehicle_State   string `json:"Vehicle_State"` // State (Active, Inactive, etc.) / (Driving, Parked, etc.)
	Owner_ID        int    `json:"Owner_ID"`      // Owner ID
	Passenger_Count int    `json:"Occupancy"`     // Occupancy
	Region          string `json:"Region"`        // Region
	X_pos           int    `json:"X_pos"`         // X Position
	Y_pos           int    `json:"Y_pos"`         // Y Position
	V_Speed         int    `json:"V_Speed"`       // Vehicle Speed (km/h)
}

type hash_info struct {
	Hash         string `json:"hash"`
	Previoushash string `json:"previous_hash"`
	Proof        string `json:"proof"`
}

type block struct {
	Index     int         `json:"index"`
	Timestamp string      `json:"timestamp"`
	Data      Information `json:"data"`
	Hash      hash_info   `json:"hash"`
}

func create_block(proposal block, add_to_chain int) {
	// Declare variables
	var new_block block
	var chain_lenght = len(blockchain)

	/*
		Verify that the Vehcile_ID and owner_ID are valid
		Verify that the Vehicle_Manu and Vehicle_Model has not changed
	*/
	switch verify_data(proposal) {
	case true:
		fmt.Println("Data is valid")
	default:
		fmt.Println("Data verification failed, block not created")
		break
	}
	// Assign the rest of the values
	new_block.Data = Information{
		proposal.Data.Vehicle_ID,
		proposal.Data.Vehicle_Manu,
		proposal.Data.Vehicle_Model,
		proposal.Data.Vehicle_State,
		proposal.Data.Owner_ID,
		proposal.Data.Passenger_Count,
		proposal.Data.Region,
		proposal.Data.X_pos,
		proposal.Data.Y_pos,
		proposal.Data.V_Speed,
	}

	// Assign values to Index, Timestamp, Data, Previoushash, and Hash
	new_block.Index = chain_lenght + 1
	new_block.Timestamp = time.Now().Format("01/02/2006 15:04:05")
	new_block.Hash.Previoushash = blockchain[chain_lenght-1].Hash.Hash

	// Calculate the Hash of the new block
	hashinfo := better_hash(new_block)
	new_block.Hash.Hash = hashinfo.Hash
	new_block.Hash.Proof = hashinfo.Proof

	// **DEBUG OPTION TO NOT APPEND**
	// Append the new block to the blockchain and save the blockchain to a JSON file
	if add_to_chain == 1 {
		blockchain = append(blockchain, new_block)
		save_blockchain_json()
	} else {
		fmt.Println("Block not added to chain")
		fmt.Println(new_block)
	}
	// If validation fails, remove the last block from the blockchain
	if !validate_blockchain() {
		blockchain = blockchain[:len(blockchain)-1]
	}
}

// Generate hash and proof of work for the block
func better_hash(block block) hash_info {
	var Hash = sha256.New()
	// Add all the block information to a byte array
	var hash_data = []byte(
		fmt.Sprint(block.Index) +
			block.Timestamp +
			fmt.Sprint(block.Data.Vehicle_ID) +
			block.Data.Vehicle_Manu +
			block.Data.Vehicle_Model +
			block.Data.Vehicle_State +
			fmt.Sprint(block.Data.Owner_ID) +
			fmt.Sprint(block.Data.Passenger_Count) +
			block.Data.Region +
			fmt.Sprint(block.Data.X_pos) +
			fmt.Sprint(block.Data.Y_pos) +
			fmt.Sprint(block.Data.V_Speed) +
			block.Hash.Previoushash,
	)
	var proof = 0
	var n bool = false
	// Iterate through the hash until the first 3 characters are 0 and return the hash and proof
	for !n {
		proof += 1
		Hash.Write(hash_data)
		Hash.Write([]byte(fmt.Sprint(proof)))
		if fmt.Sprintf("%x", Hash.Sum(nil))[:3] == "000" {
			n = true
		}
	}
	return hash_info{fmt.Sprintf("%x", Hash.Sum(nil)), block.Hash.Previoushash, fmt.Sprint(proof)}
}

// Validate the blockchain
func validate_blockchain() bool {
	// Iterate through the blockchain
	for i := 1; i < len(blockchain); i++ {
		var current_block = blockchain[i]
		var previous_block = blockchain[i-1]
		// Check if the current block's hash is correct given the current block's data, hash and proof
		test := better_hash(current_block)
		if test.Hash == current_block.Hash.Hash {
			// fmt.Println("\nValidating Hash Success")
			if test.Proof == current_block.Hash.Proof {
			} else {
				fmt.Println("\nValidating Proof Failed With ", test.Proof, " and ", current_block.Hash.Proof)
				return false
			}
		} else {
			fmt.Println("\nValidating Hash Failed", better_hash(current_block).Hash, "\n\n", current_block.Hash.Hash)
			return false
		}
		// Check if the current block's previous hash is equal to the previous block's hash
		if current_block.Hash.Previoushash != previous_block.Hash.Hash {
			fmt.Println("\nValidating Previous Hash Failed")
			return false
		}
	}
	fmt.Println("\nValidation Successfull")
	return true
}

// Calculate the balance of a sender and verify that the sender has enough funds to send the Amount
func verify_data(proposal block) bool {
	// Iterate through the blockchain
	for i := 0; i < len(blockchain); i++ {
		// Check if the Vehicle_ID and Owner_ID are valid
		if blockchain[i].Data.Vehicle_ID == proposal.Data.Vehicle_ID {
			if blockchain[i].Data.Owner_ID == proposal.Data.Owner_ID {
				// Check if the Vehicle_Manu and Vehicle_Model have not changed
				if blockchain[i].Data.Vehicle_Manu != proposal.Data.Vehicle_Manu {
					fmt.Println("Vehicle Manufacturer Mismatch")
					return false
				}
				if blockchain[i].Data.Vehicle_Model != proposal.Data.Vehicle_Model {
					fmt.Println("Vehicle Model Mismatch")
					return false
				}
			} else {
				fmt.Println("Vehicle Owner Mismatch")
				return false
			}
		}
	}
	return true
}

// Create the genesis block
func genesis_block() {
	var genesis block
	genesis.Index = 1
	genesis.Timestamp = "01/01/2018"
	genesis.Data = Information{0, "Genesis", "Genesis", "Genesis", 0, 0, "Genesis", 0, 0, 0}
	genesis.Hash.Proof = "0"
	genesis.Hash.Hash = "gabagol"
	genesis.Hash.Previoushash = "shaboingus"
	blockchain = append(blockchain, genesis)
}

// Save the current blockchain as a JSON file
func save_blockchain_json() {
	// Create a new JSON file
	os.Create("blockchain.json")
	// Marshal the blockchain into a JSON file
	json, err := json.MarshalIndent(blockchain, "", " ")
	if err != nil {
		fmt.Println(err.Error())
	}
	// Write the JSON file
	os.WriteFile("blockchain.json", []byte(json), 0644)
}

// Load the blockchain from a JSON file
func load_blockchain_json() bool {
	if _, err := os.Stat("blockchain.json"); err == nil {
		// File exists
		// Open the JSON file
		jsonFile, err := os.Open("blockchain.json")
		if err != nil {
			fmt.Println(err)
		}
		// Read the JSON file
		byteValue, _ := ioutil.ReadAll(jsonFile)
		// Unmarshal the JSON file into the blockchain
		json.Unmarshal(byteValue, &blockchain)
		// Close the JSON file
		jsonFile.Close()
		return true
	} else {
		// File does not exist
		return false
	}
}

func handleRequest(conn net.Conn) {
	data := block{}
	// incoming request
	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Received message:", string(buffer))

		words := strings.Fields(string(buffer))
		Vehicle_ID := words[0]
		Vehicle_Manu := words[1]
		Vehicle_Model := words[2]
		Vehicle_State := words[3]
		Owner_ID := words[4]
		Passenger_Count := words[5]
		Region := words[6]
		X_pos := words[7]
		Y_pos := words[8]
		V_Speed := words[9]

		data.Data.Vehicle_ID, _ = strconv.Atoi(Vehicle_ID)
		data.Data.Vehicle_Manu = Vehicle_Manu
		data.Data.Vehicle_Model = Vehicle_Model
		data.Data.Vehicle_State = Vehicle_State
		data.Data.Owner_ID, _ = strconv.Atoi(Owner_ID)
		data.Data.Passenger_Count, _ = strconv.Atoi(Passenger_Count)
		data.Data.Region = Region
		data.Data.X_pos, _ = strconv.Atoi(X_pos)
		data.Data.Y_pos, _ = strconv.Atoi(Y_pos)
		data.Data.V_Speed, _ = strconv.Atoi(V_Speed)

		create_block(data, 0)
	}
}

// CLI test function
func cli_test() {
	proposal := block{}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\nWelcome to the Blockchain CLI\nCreate new block? (y/n)")
	scanner.Scan()
	input := scanner.Text()
	if input == "y" {
		proposal.Data.Vehicle_ID = rand.Intn(100)
		proposal.Data.Vehicle_Manu = "Ford"
		proposal.Data.Vehicle_Model = "F150"
		proposal.Data.Vehicle_State = "Active"
		proposal.Data.Owner_ID = rand.Intn(100)
		proposal.Data.Passenger_Count = rand.Intn(5)
		proposal.Data.Region = "North-America"
		proposal.Data.X_pos = rand.Intn(100)
		proposal.Data.Y_pos = rand.Intn(100)
		proposal.Data.V_Speed = rand.Intn(100)

		fmt.Println("Creating new block with data:\n ", proposal.Data)
		create_block(proposal, 1)
	} else {
		fmt.Printf("\nExiting CLI %s", input)
	}

}

func main() {
	// Check if a blockchain already exists and load it
	// If no blockchain exists, create a new one with a new genesis block
	if load_blockchain_json() == false {
		genesis_block()
		save_blockchain_json()
		fmt.Println("\nNo blockchain found, creating a new one with a new genesis block")
	} else {
		fmt.Println("\nBlockchain found, loading it")
	}

	// Check if the blockchain is valid
	if validate_blockchain() == true {
		fmt.Println("\nBlockchain is valid")
	} else {
		fmt.Println("\nBlockchain is invalid!")
		os.Exit(1)
	}

	// HTTP Server Pages
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bc_json, _ := json.MarshalIndent(blockchain, "", " ")
		str_bc_json := string(bc_json)
		content := html_data{Title: "Blockchain", Data: str_bc_json}
		t, _ := template.ParseFiles("template.html")
		t.Execute(w, content)
	})
	// Display the blockchain
	http.HandleFunc("/chain", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(blockchain)
	})
	// Run validation
	http.HandleFunc("/valid", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if validate_blockchain() == true {
			w.Write([]byte("Blockchain is valid"))
		} else {
			w.Write([]byte("Blockchain is invalid!"))
		}
	})

	// Parse command line arguments
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("\nPlease specify a mode (cli/http/tcp) as an argument")
		os.Exit(1)
	}
	// Start the specified mode
	switch args[0] {
	case "cli":
		fmt.Println("Starting CLI...")
		cli_test()
	case "http":
		fmt.Println("Listening on port 8080...")
		log.Fatal(http.ListenAndServe(":8080", nil))
	case "tcp":
		listen, err := net.Listen(TYPE, HOST+":"+PORT)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		fmt.Println("Listening on host " + HOST + " and port " + PORT)
		// close listener
		defer listen.Close()
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			go handleRequest(conn)
		}
	}
}
