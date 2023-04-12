package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
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

func create_block(proposal block, add_to_chain int) block {
	// Declare variables
	var new_block block
	var new_block_Data Information
	var chain_lenght = len(blockchain)

	// Assign values to sender, Recipient, and Amount
	// Check if the sender has enough funds to send the Amount

	// if balance(sender) < amount {
	// 	// If the sender does not have enough funds, print an error message and
	// 	// return the new block without appending it to the blockchain
	// 	fmt.Printf("Insufficient funds! %s has %d and is trying to send %d\n", sender, balance(sender), amount)
	// 	return new_block
	// } else {
	// 	new_block_Data.Sender = sender
	// }

	// Assign the rest of the values
	new_block_Data.Vehicle_ID = proposal.Data.Vehicle_ID
	new_block_Data.Vehicle_Manu = proposal.Data.Vehicle_Manu
	new_block_Data.Vehicle_Model = proposal.Data.Vehicle_Model
	new_block_Data.Vehicle_State = proposal.Data.Vehicle_State
	new_block_Data.Owner_ID = proposal.Data.Owner_ID
	new_block_Data.Passenger_Count = proposal.Data.Passenger_Count
	new_block_Data.Region = proposal.Data.Region
	new_block_Data.X_pos = proposal.Data.X_pos
	new_block_Data.Y_pos = proposal.Data.Y_pos
	new_block_Data.V_Speed = proposal.Data.V_Speed

	// Assign values to Index, Timestamp, Data, Previoushash, and Hash
	new_block.Index = chain_lenght + 1
	new_block.Timestamp = time.Now().Format("01/02/2006 15:04:05")
	new_block.Data = new_block_Data
	//new_block.Proof = proof_of_work(blockchain[len(blockchain)-1])
	new_block.Hash.Previoushash = blockchain[chain_lenght-1].Hash.Hash

	// Calculate the Hash of the new block
	new_block.Hash.Hash = better_hash(new_block).Hash
	new_block.Hash.Proof = better_hash(new_block).Proof

	// Append the new block to the blockchain and save the blockchain to a JSON file
	if add_to_chain == 1 {
		blockchain = append(blockchain, new_block)
		save_blockchain_json()
	}
	//validate_blockchain()
	return new_block
}

func better_hash(block block) hash_info {
	var Hash = sha256.New()
	// Add all the block information to a byte array ********* FIX THIS *********
	var block_information = []byte(fmt.Sprint(block.Index) + block.Timestamp + fmt.Sprint(block.Data.Vehicle_ID) + block.Data.Vehicle_Manu + block.Data.Vehicle_Model + block.Data.Vehicle_State + fmt.Sprint(block.Data.Owner_ID) + fmt.Sprint(block.Data.Passenger_Count) + block.Data.Region + fmt.Sprint(block.Data.X_pos) + fmt.Sprint(block.Data.Y_pos) + fmt.Sprint(block.Data.V_Speed) + block.Hash.Previoushash)
	var proof = 0
	var n bool = false
	for !n {
		proof += 1
		Hash.Write(block_information)
		Hash.Write([]byte(fmt.Sprint(proof)))
		if fmt.Sprintf("%x", Hash.Sum(nil))[:3] == "000" {
			//fmt.Printf("\nValid proof found!: %d\n", proof)
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
		// fmt.Println("****************** Starting Validation of Block ", current_block.Index, "******************")
		// fmt.Println("\nValidating Current Block's Hash ", current_block.Hash.Hash, " to ", test.Hash)
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
		// fmt.Println("\nValidating Current Block's Previous Hash to Previous Block's Hash", current_block.Hash.Previoushash, " to ", previous_block.Hash.Hash)
		if current_block.Hash.Previoushash != previous_block.Hash.Hash {
			fmt.Println("\nValidating Previous Hash Failed")
			return false
		}
	}
	fmt.Println("\nValidation Successfull")
	return true
}

// Calculate the balance of a sender and verify that the sender has enough funds to send the Amount
// func balance(sender string) int {
// 	var balance = 0
// 	// Iterate through the blockchain
// 	for i := 0; i < len(blockchain); i++ {
// 		// Check if the sender is the sender or recipient of the given transaction
// 		// If the current sender is the sender, subtract the Amount from the balance
// 		if blockchain[i].Data.Sender == sender {
// 			balance -= blockchain[i].Data.Amount
// 		}
// 		// If the current sender is the recipient, add the Amount to the balance
// 		if blockchain[i].Data.Recipient == sender {
// 			balance += blockchain[i].Data.Amount
// 		}
// 	}
// 	return balance
// }

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
		proposal.Data.Region = "North America"
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

	// Create a few new blocks and add them to the blockchain

	// create_block("POOL", "Alice", 100, 1)
	// create_block("Alice", "Bob", 50, 1)
	// create_block("Bob", "Alice", 25, 1)
	// create_block("Steven", "James", 25, 1)
	// create_block("POOL", "Steven", 100, 1)
	// create_block("Steven", "James", 25, 1)
	// create_block("James", "Steven", 25, 1)
	// create_block("POOL", "James", 100, 1)
	// create_block("James", "Bob", 25, 1)
	// save_blockchain_json()

	// Check if the blockchain is valid

	if validate_blockchain() == true {
		fmt.Println("\nBlockchain is valid")
	} else {
		fmt.Println("\nBlockchain is invalid!")

	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bc_json, _ := json.MarshalIndent(blockchain, "", " ")
		str_bc_json := string(bc_json)
		content := html_data{Title: "Blockchain", Data: str_bc_json}
		t, _ := template.ParseFiles("template.html")
		t.Execute(w, content)
	})

	http.HandleFunc("/chain", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(blockchain)
	})

	http.HandleFunc("/valid", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if validate_blockchain() == true {
			w.Write([]byte("Blockchain is valid"))
		} else {
			w.Write([]byte("Blockchain is invalid!"))
		}
	})

	//fmt.Println("Listening on port 8080...")
	//log.Fatal(http.ListenAndServe(":8080", nil))
	cli_test()
	// var evilblock block
	// evilblock.Index = 1
	// evilblock.Timestamp = "01/01/2018"
	// evilblock.Data = Data{"Alice", "HAX0R", 10000000}
	// evilblock.Hash = "123"
	// evilblock.Previoushash = "0"
	// blockchain = append(blockchain[:2], evilblock)

}
