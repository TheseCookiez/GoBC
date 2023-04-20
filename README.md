# GoBC
Golang Blockchain created as part of a study project

Branches:
    master: Initial currency version
    AV_TCP: Final networked version with AV data
    Autonomous: In-dev AV-TCP branch
    Development: In-dev master branch

# How to Use:
  - Download / git clone
  - To run the blockchain backend:
          Run "go run bc.go tcp" to run using the Go runtime or compile to executable and run with the "tcp" flag.
  - To run the AV client:
          Run "go run client.go $Config_Name" or compile to executable. $Config_Name = "nissan/toyota" or other JSON config file
          
# Behaviour
The blockchain backend should start printing out incoming data, blockchain validation status and the elapsed time to create new blocks.   
