package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"strings"
)

// CoffeMachine is our model
type CoffeeMachine struct {
}

//
// ELEMENTARY FUNCTIONS (Init, Invoke, main)
//

// Init is called at initialization
func (t *CoffeeMachine) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke is called when chaincode is started
// This is a pattern used in all Hyperledger projects. We define, which Methods are allowed here.
func (t *CoffeeMachine) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	// Extract function and args from transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error

	//Our 4 external methods are defined here:
	fnmap := map[string]func(stubInterface shim.ChaincodeStubInterface, strings []string) (string, error){
		"cleanMachine": cleanMachine,
		"refillCoffee": refillCoffee,
		"drawCoffee":   drawCoffee,
		"storeUser":    storeUser,
	}

	function, ok := fnmap[fn]

	if ok == false {
		// no such function
		return shim.Error("Can't find requested function.")
	}

	// function gets called here
	result, err = function(stub, args)

	// something went wrong
	if err != nil {
		return shim.Error(err.Error())
	}

	// success, return result
	return shim.Success([]byte(result))

}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(CoffeeMachine)); err != nil {
		fmt.Printf("Error starting CoffeeMachine chaincode: %s", err)
	}
}


//
// USER FUNCTIONS
//

// Adding a new user
func storeUser(stubInterface shim.ChaincodeStubInterface, args []string) (string, error) {

	name := args[0]
	setCounter(stubInterface, 1)
	counter := getCounter(stubInterface) +1
	var userkey strings.Builder
	userkey.WriteString("user" + strconv.FormatInt(counter, 10))

	err := stubInterface.PutState(userkey.String(), []byte(name))

	return "User "+userkey.String()+" stored", err

}

// List our users
func getListOfUsers(stubInterface shim.ChaincodeStubInterface) []string {

	iterator, err := stubInterface.GetStateByPartialCompositeKey("user", []string{"user"})
	if err != nil {
		// return empty slice
		return []string{}
	}
	defer iterator.Close()

	var userlist = make([]string, 5)
	var i int

	for i = 0; iterator.HasNext(); i++ {
		responseRange, err := iterator.Next()
		if err != nil {
			// error? Skip this one
			continue
		}
		userlist = append(userlist, string(responseRange.Value))
	}

	return userlist

}

// Find user to clean or refill
func getDutyUser(stubInterface shim.ChaincodeStubInterface) string {
	// gets the current count from the blockchain
	users := getListOfUsers(stubInterface)
	counter := getCounter(stubInterface)
	return users[counter]
}

// Get User counter (see, who´s next)
func getCounter(stubInterface shim.ChaincodeStubInterface) int64 {
	// gets the current user count from the blockchain
	var raw []byte
	var counter int64
	var err error

	// pull current count from the blockchain
	raw, err = stubInterface.GetState("userCounter")

	if err != nil {
		// something went wrong, we return a bogus value
		// not good, but short
		return -1
	}

	counter, err = binary.ReadVarint(bytes.NewBuffer(raw))

	if err != nil {
		// and something went wrong
		return -1
	}

	return counter
}

// Increment user counter
func setCounter(stubInterface shim.ChaincodeStubInterface, increment int) {

	if increment == 0 {
		increment = 1
	}

	lengthofuserlist := len(getListOfUsers(stubInterface))
	// count elements in list
	counter := (getCounter(stubInterface) + int64(increment)) % int64(lengthofuserlist)

	newcounter := []byte(strconv.FormatInt(counter, 10))

	stubInterface.PutState("userCounter", newcounter)

}

//
// COFFEE AND MACHINE FUNCTIONS
//

// Get coffeeLevel
func getCoffeeLevel(stubInterface shim.ChaincodeStubInterface) int64 {
	var raw []byte
	var v int64
	var ok error

	raw, ok = stubInterface.GetState("coffeeLevel")
	if ok != nil {
		return 0
	}

	v, ok = binary.ReadVarint(bytes.NewBuffer(raw))
	if ok == nil {
		return v
	} else {
		return 0
	}
}

// Get dirtLevel
func getDirtLevel(stubInterface shim.ChaincodeStubInterface) int64 {
	var raw []byte
	var v int64
	var ok error

	raw, ok = stubInterface.GetState("dirtLevel")
	if ok != nil {
		return 0
	}

	v, ok = binary.ReadVarint(bytes.NewBuffer(raw))
	if ok == nil {
		return v
	} else {
		return 0
	}
}

// Set coffeeLevel
func setCoffeeLevel(stubInterface shim.ChaincodeStubInterface, coffeeLevel int64) (error) {
	var err error

	if coffeeLevel <0 {
		return fmt.Errorf("Machine empty.")
	}

	// convert int64 to []byte
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, coffeeLevel)
	b := buf[:n]

	err = stubInterface.PutState("coffeeLevel", b)

	if err != nil {
		return fmt.Errorf("An error occured while saving the coffeeLevel.")
	}

	return nil
}

// Set dirtLevel
func setDirtLevel(stubInterface shim.ChaincodeStubInterface, dirtLevel int64) error {
	var err error

	if dirtLevel <0 {
		return fmt.Errorf("Machine dirty.")
	}

	// convert int64 to []byte
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, dirtLevel)
	b := buf[:n]

	err = stubInterface.PutState("dirtLevel", b)

	if err != nil {
		return fmt.Errorf("An error occured while saving the dirtLevel.")
	}

	return nil
}

// Method to be called externally: clean the machine (set dirtCounter to 60)
func cleanMachine(stubInterface shim.ChaincodeStubInterface, args []string) (string, error) {

	// write new dirtLevel into blockchain
	return "Machine cleaned. dirtLevel set to 60", setDirtLevel(stubInterface, 60)

}

// Method to be called externally: refill the machine (set coffeeCounter to 25)
func refillCoffee(stubInterface shim.ChaincodeStubInterface, args []string) (string, error) {

	// write new CoffeeLevel into blockchain
	return "Machine refilled. coffeeLevel set to 25", setCoffeeLevel(stubInterface, 25)
}


// Draw a coffee. Contains our main business logic
func drawCoffee(stubInterface shim.ChaincodeStubInterface, args []string) (string, error) {
	// draws x coffees from machine
	 if len(args) != 1 {
	   return "", fmt.Errorf("Incorrect arguments. Expecting a value")
	 }
   
	 var cups int
	 var user string
	 var err error
   
	 cups, err = strconv.Atoi(args[0])
   
	 if err != nil {
	   return "", 
	   fmt.Errorf("Expecting an integer value.")
	 }
   
	 // do the math
	 coffeeLevel := getCoffeeLevel(stubInterface)
	 coffeeLevel -= int64(cups)
   
	 // store the coffeeLevel in the blockchain
	 err = setCoffeeLevel(stubInterface, coffeeLevel)
	 msg := "Not enough coffee in machine!"
   
	 if err == nil {
	   return "Here is your coffee!", nil
	 }else{
		user = getDutyUser(stubInterface)
		msg = msg  + "It is your job, " + user + "!"
		return msg, nil
	 }
   

   }