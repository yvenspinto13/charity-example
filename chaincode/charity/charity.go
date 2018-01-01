package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the donation structure, with 4 properties.  Structure tags are used by encoding/json library
type Donation struct {
	Donor	string `json:"donor"`
	Amount	string `json:"amount"`
	Date	string `json:"date"`
	Cause	string `json:"cause"`
}

/*
 * The Init method is called when the Smart Contract "fabdonation" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabdonation"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryDonation" {
		return s.queryDonation(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createDonation" {
		return s.createDonation(APIstub, args)
	} else if function == "queryAllDonations" {
		return s.queryAllDonations(APIstub)
	} else if function == "totalDonationAmount" {
		return s.totalDonationAmount(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryDonation(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	donationAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(donationAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	donations := []Donation{
		Donation{Donor: "Yvens Pinto", Amount: "30000", Date: "13/12/2017", Cause: "23ed Birthday"},
		Donation{Donor: "Lerisa Gomes", Amount: "20000", Date: "04/11/2017", Cause: "1st Salary"},
		Donation{Donor: "Vishal Robertson", Amount: "10000", Date: "24/08/2017", Cause: "Anniversary"},
		Donation{Donor: "Asif Muhamad", Amount: "60000", Date: "01/01/2017", Cause: "New Year"},
		Donation{Donor: "Elon Musk", Amount: "30000", Date: "31/05/2017", Cause: "Company Bonus"},
	}

	i := 0
	for i < len(donations) {
		fmt.Println("i is ", i)
		donationAsBytes, _ := json.Marshal(donations[i])
		APIstub.PutState("DONATION"+strconv.Itoa(i), donationAsBytes)
		fmt.Println("Added", donations[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createDonation(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var donation = Donation{Donor: args[1], Amount: args[2], Date: args[3], Cause: args[4]}

	donationAsBytes, _ := json.Marshal(donation)
	APIstub.PutState(args[0], donationAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllDonations(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "DONATION0"
	endKey := "DONATION999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")	
	fmt.Printf("- queryAllDonations:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) totalDonationAmount(APIstub shim.ChaincodeStubInterface) sc.Response {
	var totalAmount int64 = 0
	startKey := "DONATION0"
	endKey := "DONATION999"
	
	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("Total Donations: [")

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
	
		if err != nil {
			return shim.Error(err.Error())
		}

		stringRow := string(queryResponse.Value)
		stringSlice1 := strings.Split(stringRow, ",")
		stringSlice2 := strings.Split(stringSlice1[0], ":")
		stringSlice3 := strings.Split(stringSlice2[1], "\"")
		stringAmount := stringSlice3[1]

		i64, err := strconv.ParseInt(stringAmount, 10, 64)
		if err != nil {
			return shim.Error(err.Error())
		}
		totalAmount = totalAmount + i64

	}
	buffer.WriteString(strconv.FormatInt(totalAmount, 10))
	buffer.WriteString("]")

	fmt.Printf("- queryAllDonations:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}