package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	//IsAdmin  string `json:"IsAdmin"`
	IsOnline string `json:"IsOnline"`
}

type BankChaincode struct {
}

var (
	paramLength      map[string]int
	paramLengthError map[string]string
)

func init() {
	paramLength = make(map[string]int)
	paramLengthError = make(map[string]string)

	paramLength["add"] = 2
	paramLengthError["add"] = "Incorrect arguments. Expecting an account name and a balance value."
	paramLength["create"] = 2
	paramLengthError["create"] = "Incorrect arguments. Expecting an unique account name and an initial balance value."
	paramLength["delete"] = 1
	paramLengthError["delete"] = "Incorrect arguments. Expecting an account being deleted."
	paramLength["get"] = 1
	paramLengthError["get"] = "Incorrect arguments. Expecting an account name."
	paramLength["query"] = 2
	paramLengthError["query"] = "Incorrect arguments. Expecting an objectType and an account."
	paramLength["reduce"] = 2
	paramLengthError["reduce"] = "Incorrect arguments. Expecting an account name and a balance value."
	paramLength["rollback"] = 3
	paramLengthError["rollback"] = "Incorrect arguments. Expecting a debit account, credit account and a transaction id."
	paramLength["transfer"] = 3
	paramLengthError["transfer"] = "Incorrect arguments. Expecting a debit account, a credit account and a value"
	paramLength["register"] = 3
	paramLengthError["register"] = "Incorrect arguments. Expecting a username"
	paramLength["alterPasswd"] = 3
	paramLengthError["alterPasswd"] = "Incorrect arguments. Expecting a username , password and an alterPasswd "
	paramLength["login"] = 2
	paramLengthError["login"] = "Incorrect arguments. Expecting a debit account, a credit account and a value"
	paramLength["loginOut"] = 2
	paramLengthError["loginOut"] = "Incorrect arguments. Expecting a debit account, a credit account and a value"
}

func (t *BankChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	function, args := stub.GetFunctionAndParameters()
	//the first argument is in the variable "function"
	if function != "init" {
		return shim.Error("The first parameter needs to be a string: \"init\"")
	}
	if len(args) != 0 {
		return shim.Error("Incorrect arguments. Expecting an account name and a balance value")
	}

	// Set up any variables or assets here by calling stub.PutState()
	// We store the key and the value on the ledger
	//
	//err := stub.PutState(args[0], []byte(args[1]))
	//if err != nil {
	//	return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	//}
	return shim.Success([]byte(fmt.Sprintf("Success to initialize!")))
}

func (t *BankChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	// check the param length
	if expectedLen, ok := paramLength[fn]; !ok {
		return shim.Error("Undefined function")
	} else if expectedLen != len(args) {
		errStr, _ := paramLengthError[fn]
		return shim.Error(errStr)
	}

	var result string
	var err error

	switch fn {
	case "get":
		result, err = get(stub, args)
	case "add":
		result, err = add(stub, args)
	case "reduce":
		result, err = reduce(stub, args)
	case "create":
		result, err = create(stub, args)
	case "delete":
		result, err = deleteAcc(stub, args)
	case "transfer":
		result, err = transfer(stub, args)
	case "query":
		result, err = query(stub, args)
	case "rollback":
		result, err = rollback(stub, args)
	case "register":
		result, err = register(stub, args)
	case "alterPasswd":
		result, err = alterPasswd(stub, args)
	case "login":
		result, err = login(stub, args)
	case "loginOut":
		result, err = loginOut(stub, args)
	default:
		return shim.Error("You do not have authority to get access to this function!")
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

//根据username得到value（user对象）
func getUsrInfo(stub shim.ChaincodeStubInterface, username string) (User, bool) {
	var usr User
	b, err := stub.GetState(username)
	if err != nil {
		return usr, false
	}

	if b == nil {

		return usr, false
	}

	err = json.Unmarshal(b, &usr)
	if err != nil {
		return usr, false
	}

	return usr, true
}

//store the user
func putUsr(stub shim.ChaincodeStubInterface, user User) ([]byte, bool) {
	b, err := json.Marshal(user)
	if err != nil {
		return nil, false
	}

	//save the status of user
	err = stub.PutState(user.Username, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

//key is username ,User is value
//args[0] is username
func register(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var usr User

	usr.Username = args[0]
	usr.Password = args[1]
	usr.IsOnline = args[2]

	//err := json.Unmarshal([]byte(args[0]), &usr)
	//if err != nil {
	//	return "", fmt.Errorf("Deserialize failed!!")
	//}

	//duplicate checking
	_, exist := getUsrInfo(stub, usr.Username)
	if exist {
		return "", fmt.Errorf("The user to register already exists")
	}

	_, bl := putUsr(stub, usr)
	if !bl {
		return "", fmt.Errorf("Failed to register user")
	}

	err := stub.SetEvent(args[1], []byte{})
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return fmt.Sprintf("Success to register %s!", args[0]), nil
}

//
//update user information
//args:userObject
//need two param
//func  updateUsr(stub shim.ChaincodeStubInterface, args []string) (string, bool) {
//	var info User
//	err := json.Unmarshal([]byte(args[0]), &info)
//	if err != nil {
//		return "", false
//	}
//
//	result, bl := getUsrInfo(stub, info.Username)
//	if !bl {
//		return "", false
//	}
//
//	result.IsOnline = info.IsOnline
//	result.Password = info.Password
//
//	_, bl = putUsr(stub, result)
//	if !bl {
//		return "", false
//	}
//
//	return "", true
//
//}

//Update permission
//args[0] username
//args[1] is isadmin （you can only input yes or no）
//func (t *BankChaincode) updatePerm(stub shim.ChaincodeStubInterface, args []string) (string, error) {
//	if args[0] != "yes" && args[0] != "no" {
//		return "", fmt.Errorf("You can only input yes or no!")
//	}
//
//	var info User
//	info.Username = args[0]
//	result, bl := getUsrInfo(stub, info.Username)
//	if !bl {
//		return "", fmt.Errorf("Get userInfo failed!")
//	}
//	result.IsAdmin = args[1]
//
//	_, bl = putUsr(stub, result)
//	if !bl {
//		return "", fmt.Errorf("Put userObject failed!")
//	}
//
//	return fmt.Sprintf("Update the permission success!!"), nil
//}

//args[0] username
//args[1] oldpassword
//args[2] newpassword
func alterPasswd(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var info User
	info.Username = args[0]
	result, bl := getUsrInfo(stub, info.Username)
	if !bl {
		return "", fmt.Errorf("Get userInfo failed!")
	}
	if args[1] != result.Password {
		return "", fmt.Errorf("Wrong password! You should input the true passwprd!")
	}

	if args[1] == args[2] {
		return "", fmt.Errorf("Wrong password! You should alter the passwprd!")
	}

	result.Password = args[2]

	_, bl = putUsr(stub, result)
	if !bl {
		return "", fmt.Errorf("Put userObject failed!")
	}

	return fmt.Sprintf("Alter the password success!!"), nil
}

//args[0] is username
//args[1] is password
func login(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var info User
	info.Username = args[0]
	result, bl := getUsrInfo(stub, info.Username)
	if !bl {
		return "", fmt.Errorf("Get userInfo failed!")
	}
	if args[1] != result.Password {
		return "", fmt.Errorf("Wrong password! You should input the true passwprd!")
	}
	if result.IsOnline == "Yes" {

		return "", fmt.Errorf("You're logged in !")
	}
	result.IsOnline = "Yes"

	_, bl = putUsr(stub, result)
	if !bl {
		return "", fmt.Errorf("Put userObject failed!")
	}
	return fmt.Sprintf("Log in success!!"), nil
}

//args[0] is username
//args[1] is password
func loginOut(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var info User
	info.Username = args[0]
	result, bl := getUsrInfo(stub, info.Username)
	if !bl {
		return "", fmt.Errorf("Get userInfo failed!")
	}
	if args[1] != result.Password {
		return "", fmt.Errorf("Wrong password! You should input the true passwprd!")
	}
	if result.IsOnline == "No" {
		return "", fmt.Errorf("You're not logged in!!")
	}
	result.IsOnline = "No"
	_, bl = putUsr(stub, result)
	if !bl {
		return "", fmt.Errorf("Put userObject failed!")
	}

	return fmt.Sprintf("Log out success!!"), nil
}

// Get returns the value of the specified asset key
// When we need to query the remaining balance, we use this function.
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	// get the account information from the database.
	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}

	return fmt.Sprintf(" Account: %s; Balance: %s", args[0], string(value)), nil
}

// args[0] represents account, args[1] represents money.
// Add specific number of money to the specific account.
func add(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	valueTemp, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}

	intArgs1, err := strconv.Atoi(args[1])
	if err != nil {
		return "", fmt.Errorf("Atoi fail! With Error: %s", err)
	}

	intValueTemp, err := strconv.Atoi(string(valueTemp))
	if err != nil {
		return "", fmt.Errorf("Atoi fail! With Error: %s", err)
	}

	err = stub.PutState(args[0], []byte(strconv.Itoa(intArgs1+intValueTemp)))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s with error: %s", args[0], err)
	}

	return fmt.Sprintf("Add is success! Account: %s; Remaining balance is: %d", args[0], intValueTemp+intArgs1), nil

}

// args[0] represents account, args[1] represents money.
// Reduce specific number of money to the specific account.
func reduce(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	// Get the account from the worldstate database.
	valueTemp, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	// change the argument into integer.
	intArgs1, err := strconv.Atoi(args[1])
	if err != nil {
		return "", fmt.Errorf("Atoi fail! With Error: %s", err)
	}
	//
	intValueTemp, err := strconv.Atoi(string(valueTemp))
	if err != nil {
		return "", fmt.Errorf("Atoi fail! With Error: %s", err)
	}

	if intArgs1 > intValueTemp {
		return "", fmt.Errorf("The balance in %s's account is not enough to reduce!", args[0])
	}

	err = stub.PutState(args[0], []byte(strconv.Itoa(intValueTemp-intArgs1)))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s;  With Error: %s", args[0], err)
	}

	return fmt.Sprintf("Reduce is success! Account: %s; Remaining balance is: %d", args[0], intValueTemp-intArgs1), nil

}

// The function of this module is to create an account of ledger
// args[0] means the account ID
// args[1] means the account initial value.
func create(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var name []byte
	name, err := stub.GetState(args[0])
	if name != nil {
		return "", fmt.Errorf(fmt.Sprintf("The account has already existed!"))
	}
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Failed to get access to asset: %s; With error: %s", args[0], err))
	}

	// Set up any variables or assets here by calling stub.PutState()
	// We store the key and the value on the ledger
	err = stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Failed to create asset: %s; With Error: %s", args[0], err))
	}

	return fmt.Sprintf("Create account: %s  is success!", args[0]), nil

}

// delete an account of ledger.
// args[0] represents the account ID.
func deleteAcc(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	// delete the account.
	err := stub.DelState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to delete asset: %s with error: %s", args[0], err)
	}
	return fmt.Sprintf("Delete is success! Account: %s", args[0]), nil
}

// args[0] represents the debit account
// args[1] represents the credit account
// args[2] represents the money.
// transfer the money from the debit account to the credit account.
func transfer(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	//reduce money from the debit account.
	var argsD []string = make([]string, 2)
	argsD[0] = args[0]
	argsD[1] = args[2]
	_, err := reduce(stub, argsD)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Reduce debit account failed!"))
	}

	//add money to the cebit account.
	var argsC []string = make([]string, 2)
	argsC[0] = args[1]
	argsC[1] = args[2]
	_, err = add(stub, argsC)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Add credit account failed!"))
	}
	// store the transfer record into the database
	// "out" means the money go out from one's account,
	// so the organization of the key-value pair is:
	// Key is a composite key, its sequence is ["out"debit account] [credit account] [uuid] [time]
	// value is the amount of money been transfered.
	_, err = createHistoryKey(stub, args, "out")
	if err != nil {
		return "", fmt.Errorf("Create history records failed! with error: %s", err)
	}

	// store the transfer record into the database
	// "in" means the money go into one's account,
	// so the organization of the key-value pair is:
	// Key is a composite key, its sequence is ["in"credit account] [debit account] [uuid] [time]
	// value is the amount of money been transfered.
	_, err = createHistoryKey(stub, args, "in")
	if err != nil {
		return "", fmt.Errorf("Create history records failed! with error: %s", err)
	}

	return fmt.Sprintf("Transfer is success!"), nil
}

// create history transferring records
// "out" means the money go out from one's account,
// "in" means the money go into one's account,
// both "out" and "in" is tags, they emphasize on going out or in records
func createHistoryKey(stub shim.ChaincodeStubInterface, args []string, first string) (string, error) {
	// get the time of the transaction been finished.
	FormatTime, err := stub.GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Get transaction timestamp failed!"))
	}
	tm := time.Unix(FormatTime.Seconds, 0)

	// if we need to create an "out" record
	// the organization of the key-value pair is:
	// Key is a composite key, its sequence is ["out"debit account] [credit account] [uuid] [time]
	// value is the amount of money been transfered.
	if first == "out" {
		historyKey, err := stub.CreateCompositeKey(first, []string{
			args[0], "->", args[1],
			"\t", stub.GetTxID(),
			"\t", tm.Format("Mon Jan 2 15:04:05 +0800 UTC 2006"),
		})
		if err != nil {
			return "", fmt.Errorf("Create historyKey failed! With error: %s", err)
		}

		err = stub.PutState(historyKey, []byte(args[2]))
		if err != nil {
			return "", fmt.Errorf("Store transfer information failed! With error: %s", err)
		}

	} else if first == "in" {
		// so the organization of the key-value pair is:
		// Key is a composite key, its sequence is ["in"credit account] [debit account] [uuid] [time]
		// value is the amount of money been transfered.
		historyKey, err := stub.CreateCompositeKey(first, []string{
			args[1], "<-", args[0],
			"\t", stub.GetTxID(),
			"\t", tm.Format("Mon Jan 2 15:04:05 +0800 UTC 2006"),
		})
		if err != nil {
			return "", fmt.Errorf("Create historyKey failed! With error: %s", err)
		}

		err = stub.PutState(historyKey, []byte(args[2]))
		if err != nil {
			return "", fmt.Errorf("Store transfer information failed! With error: %s", err)
		}
	}

	return fmt.Sprintf("Insert records success!"), nil
}

// query for the transferring history.
// args[0] represents the objectType, that is, "in" or "out"
// args[1] represents the account name
func query(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var PCKey []string = make([]string, 1)
	PCKey[0] = args[1]
	// intend to get the record of transferring
	it, err := stub.GetStateByPartialCompositeKey(args[0], PCKey)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Cannot get by partial composite key!"))
	}

	defer it.Close()
	// result contains all the appropriate results
	result := ""
	header := "AccountAccociation | ID | Time | Amount"
	if args[0] != "in" && args[0] != "out" {
		return "", fmt.Errorf(fmt.Sprintf("You have typed a wrong objectType!"))
	}

	for it.HasNext() {
		item, err := it.Next()
		if err != nil {
			return "", fmt.Errorf(fmt.Sprintf("Get next of iterator failed!"))
		}
		result = result + fmt.Sprintf("%s\t%s\n", item.GetKey()[(len(args[0])+1):], item.GetValue()) // omit "in" / "out"
	}

	if result == "" {
		return "", fmt.Errorf("Do not have any records!")
	} else {
		return fmt.Sprintf("Query success! The entries obeys <%s>:\n%s", header, result), nil
	}
}

// the supervisor can rollback the transferring operation
// args[0] represents debit account in transferring record
// args[1] represents credit account in transferring record
// args[2] represents transaction id in transferring record
func rollback(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	// get satisfied out record
	var PCKeyOut []string = make([]string, 1)
	PCKeyOut[0] = args[0]
	itOut, err := stub.GetStateByPartialCompositeKey("out", PCKeyOut)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Cannot get by partial composite key when get \"in\" record!"))
	}
	//get money value and delete "out" record
	err = itOut.Close()

	if err != nil {
		return "", fmt.Errorf("Close itIn failed!")
	}

	var money []byte
	if itOut.HasNext() == false {
		return "", fmt.Errorf(fmt.Sprintf("Database do not have such records! Please check you arguments!"))
	}
	for itOut.HasNext() {
		item, err := itOut.Next()
		if err != nil {
			return "", fmt.Errorf(fmt.Sprintf("Get next of iteratorOut failed!"))
		}
		// get attribute from composite key
		_, attrArray, err := stub.SplitCompositeKey(item.GetKey())
		if err != nil {
			return "", fmt.Errorf(fmt.Sprintf("Split composite key failed!"))
		}
		// compare the input hash code with the hash code stored in database
		IsThisOne := strings.Compare(attrArray[4], args[2])
		if IsThisOne == 0 {
			money = item.GetValue()

			err = stub.DelState(item.GetKey())

			if err != nil {
				return "", fmt.Errorf("Delete state failed!")
			}

			break
		}
	}
	// delete "in" record
	var PCKeyIn []string = make([]string, 1)
	PCKeyIn[0] = args[1]
	itIn, err := stub.GetStateByPartialCompositeKey("in", PCKeyIn)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Cannot get by partial composite key when get \"in\" record!"))
	}

	err = itIn.Close()

	if err != nil {
		return "", fmt.Errorf("Close itIn failed!")
	}

	if itIn.HasNext() == false {
		return "", fmt.Errorf(fmt.Sprintf("Database do not have such records! Please check you arguments!"))
	}
	for itIn.HasNext() {
		item, err := itIn.Next()
		if err != nil {
			return "", fmt.Errorf(fmt.Sprintf("Get next of iteratorIn failed!"))
		}
		// get attribute from composite key
		_, attrArray, err := stub.SplitCompositeKey(item.GetKey())
		if err != nil {
			return "", fmt.Errorf(fmt.Sprintf("Split composite key failed!"))
		}
		// compare the input hash code with the hash code stored in database
		IsThisOne := strings.Compare(attrArray[4], args[2])
		if IsThisOne == 0 {

			err = stub.DelState(item.GetKey())

			if err != nil {
				return "", fmt.Errorf("Delete state failed!")
			}

			break
		}
	}

	// Then we should put money back into debit account.
	//reduce money from the debit account.
	var argsD []string = make([]string, 2)
	argsD[0] = args[1]
	argsD[1] = string(money)
	_, err = reduce(stub, argsD)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Reduce debit account failed! With error: %s", err))
	}

	//add money to the cebit account.
	var argsC []string = make([]string, 2)
	argsC[0] = args[0]
	argsC[1] = string(money)
	_, err = add(stub, argsC)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Add cebit account failed! With error: %s", err))
	}

	return fmt.Sprintf("rollback Success!"), nil
}

func main() {
	if err := shim.Start(new(BankChaincode)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
