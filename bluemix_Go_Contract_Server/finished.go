package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

var AreaWeatherIndexStr = "_areaindex"    //name for the key/value that will store a list of all known Weathers
var ActiveStakeStr = "_openstake" //name for the key/value that will store all open stake
var UserIndexStr = "_userindex"   

type Weather struct {
	Name        string `json:"name"`        // rainy sunny cloudy
	Temperature int    `json:"temperature"` // unused
}

type Area struct {
	Name         string    `json:"name"` //the fieldtags are needed to keep case from bouncing around
	Address      string    `json:"address"`  
	Owner        string    `json:"owner"` //unused
	WeatherIndex []Weather `json:"weather_index"`
}

type User struct {
	Name string `json:"name"`     
	Coin int    `json:"Coin"`
}

type AnStake struct {      // if this prediction is corrected, the Owner get coin = Number * Rate
	Owner string `json:"owner"` //  this stake owner
	Timestamp     int64  `json:"timestamp"`     // when this stake entry into force
	Number        int    `json:"number"`        // Number of coin
	Rate          int    `json:"rate"`          // decide how many coins stake owner will get.
	State         string `json:"state"`         // wait active end
	Insurant      string `json:"insurant"`      // who is the area owner -- unused
}

type ActiveStake struct {
	AllStake []AnStake `json:"all_stake"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval))) 
	if err != nil {
		return nil, err
	}

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty) //marshal an emtpy array of strings to clear the index
	err = stub.PutState(AreaWeatherIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	err = stub.PutState(UserIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	var stakes ActiveStake
	jsonAsBytes, _ = json.Marshal(stakes) //clear the open trade struct
	err = stub.PutState(ActiveStakeStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}


func (t *SimpleChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "write" { //writes a value to the chaincode state
		return t.Write(stub, args)
	} else if function == "create_user" { //create a new user
		return t.create_user(stub, args)
	} else if function == "create_area" { //create a new area
		return t.create_area(stub, args)
	} else if function == "create_stake" { //create an open tstake
		return t.create_stake(stub, args)
		fmt.Println("create_stake")
	} else if function == "update_weather" { //update the weather of one area
		return t.update_weather(stub, args)
		fmt.Println("update weather")
	}
	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil //send it onward
}

// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// Create User - create a new User,
// ============================================================================================================================

func (t *SimpleChaincode) create_user(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1       2     3
	//  'name'   'money'
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	//input sanitation
	fmt.Println("- start create user")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}

	name := strings.ToLower(args[0])
	coin, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("2rd argument must be a numeric string")
	}

	//check if user already exists
	UserAsBytes, err := stub.GetState(name)
	if err != nil {
		return nil, errors.New("Failed to get marble name")
	}
	res := User{}
	json.Unmarshal(UserAsBytes, &res)
	if res.Name == name {
		fmt.Println("This user arleady exists: " + name)
		fmt.Println(res)
		return nil, errors.New("This user arleady exists") //all stop a user by this name exists
	}

	//build the user json string manually
	var user User
	user.Name = name
	user.Coin = coin
	UserAsBytes, err = json.Marshal(user)
	err = stub.PutState(name, UserAsBytes) //store user with id as key
	if err != nil {
		return nil, err
	}

	//get the user index
	UsersAsBytes, err := stub.GetState(UserIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get marble index")
	}
	var UserIndex []string
	json.Unmarshal(UsersAsBytes, &UserIndex) //un stringify it aka JSON.parse()

	//append
	UserIndex = append(UserIndex, name) //add user name to index list
	fmt.Println("! User index: ", UserIndex)
	jsonAsBytes, _ := json.Marshal(UserIndex)
	err = stub.PutState(UserIndexStr, jsonAsBytes) //store name of user

	fmt.Println("- end create User")
	return nil, nil
}

// ============================================================================================================================
// create area
// ============================================================================================================================
func (t *SimpleChaincode) create_area(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1       2     3                4
	//  'name'   'addre' 'own'  'weathername'  Temperature
	if len(args) <= 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting >=4")
	}
	stub.PutState("start create area", []byte(strings.ToLower(args[0])))
	//input sanitation
	fmt.Println("- start create area")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	newarea := Area{}
	name := strings.ToLower(args[0])
	newarea.Name = name
	newarea.Address = strings.ToLower(args[1])
	newarea.Owner = strings.ToLower(args[2])

	fmt.Println("- create new area")
	jsonAsBytes, _ := json.Marshal(newarea)
	err = stub.PutState("_debug1", jsonAsBytes)

	for i := 3; i < len(args); i++ { //create and append each willing trade
		Temperature, err := strconv.Atoi(args[i+1])
		if err != nil {
			msg := "is not a numeric string " + args[i+1]
			fmt.Println(msg)
			return nil, errors.New(msg)
		}

		Weather_now := Weather{}
		Weather_now.Name = args[i]
		Weather_now.Temperature = Temperature
		fmt.Println("! created weather: " + args[i])
		jsonAsBytes, _ = json.Marshal(Weather_now)
		err = stub.PutState("_debug2", jsonAsBytes)

		newarea.WeatherIndex = append(newarea.WeatherIndex, Weather_now)
		fmt.Println("! appended weather")
		i++
	}

	//check if area already exists
	AreaAsBytes, err := stub.GetState(name)
	if err != nil {
		return nil, errors.New("Failed to get area name")
	}

	res := Area{}
	json.Unmarshal(AreaAsBytes, &res)
	if res.Name == name {
		fmt.Println("This area arleady exists: " + name)
		fmt.Println(res)
		return nil, errors.New("This area arleady exists") //all stop cause this name exists
	}

	newareaAsBytes, _ := json.Marshal(newarea)
	stub.PutState(name, newareaAsBytes)
	//get the Area index
	AreaAsBytes, err = stub.GetState(AreaWeatherIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get marble index")
	}
	var AreaIndex []string

	json.Unmarshal(AreaAsBytes, &AreaIndex) //un stringify it aka JSON.parse()
	//append
	AreaIndex = append(AreaIndex, name) //add area name to index list
	AreaAsBytes, _ = json.Marshal(AreaIndex)
	err = stub.PutState(AreaWeatherIndexStr, AreaAsBytes) //store name of area

	fmt.Println("- end create User")
	return nil, nil
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// ============================================================================================================================
// Create Stake- create a new Stake
// ============================================================================================================================

func (t *SimpleChaincode) create_stake(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1       2     3
	//  'insurant'   'beneficial' 'Number' 'rate' 'state'
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	//input sanitation
	fmt.Println("- start create user")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}

	new_stake := AnStake{}
	new_stake.Insurant = strings.ToLower(args[0])
	new_stake.Owner = strings.ToLower(args[1])
	new_stake.Number, err = strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}
	new_stake.Rate, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("4rd argument must be a numeric string")
	}
	new_stake.State = strings.ToLower(args[4])
	new_stake.Timestamp = makeTimestamp()

	jsonAsBytes, _ := json.Marshal(new_stake)
	err = stub.PutState("_debug1", jsonAsBytes)

	//get the marble index
	StakeAsBytes, err := stub.GetState(ActiveStakeStr)
	if err != nil {
		return nil, errors.New("Failed to get marble index")
	}
	var Stakes ActiveStake
	json.Unmarshal(StakeAsBytes, &Stakes) //un stringify it aka JSON.parse()

	//append
	Stakes.AllStake = append(Stakes.AllStake, new_stake) //add marble name to index list
	StakeAsBytes, _ = json.Marshal(Stakes)
	err = stub.PutState(ActiveStakeStr, StakeAsBytes)

	fmt.Println("- end create User")
	return nil, nil
}

func (t *SimpleChaincode) update_weather(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1       2     3
	//  'area_name'   'weather type' 'Temperature'
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	//input sanitation
	fmt.Println("- start create user")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}

	Temperature, err := strconv.Atoi(args[2])
	if err != nil {
		msg := "is not a numeric string " + args[2]
		fmt.Println(msg)
		return nil, errors.New(msg)
	}
	Weather_now := Weather{}
	Weather_now.Name = args[1]
	Weather_now.Temperature = Temperature

	jsonAsBytes, _ := json.Marshal(Weather_now)
	err = stub.PutState("_debug2", jsonAsBytes)
	areaname := strings.ToLower(args[0])
	areaAsByte, err := stub.GetState(areaname)
	if err != nil {
		return nil, errors.New("area not exist")
	}

	var update_area Area
	json.Unmarshal(areaAsByte, &update_area)

	jsonAsBytes, _ = json.Marshal(update_area.WeatherIndex)

	err = stub.PutState("_debug4", jsonAsBytes)

	update_area.WeatherIndex = append(update_area.WeatherIndex, Weather_now)

	areaAsByte, err = json.Marshal(update_area)
	if err != nil {
		return nil, errors.New("area marshal fail")
	}
	stub.PutState(areaname, areaAsByte)

	err = stub.PutState("_debug3", areaAsByte)

	//check if terrible weather
	if len(update_area.WeatherIndex) >= 3 {
		var Stakes ActiveStake
		StakeAsBytes, err := stub.GetState(ActiveStakeStr)
		if err != nil {
			return nil, errors.New("stake get error")
		}
		json.Unmarshal(StakeAsBytes, &Stakes) //un stringify it aka JSON.parse()
		bad_count := 0
		wl := len(update_area.WeatherIndex)
		for i := wl - 3; i < wl; i++ {
			Weather := update_area.WeatherIndex[i]
			if Weather.Name == "rainy" {
				bad_count += 1
			}
			if Weather.Name == "sunny" {
				bad_count = 0
			}
		}
		if bad_count >= 3 {
			for i, val := range Stakes.AllStake {
				if val.State == "actived" && val.Insurant == areaname {
					Stakes.AllStake[i].State = "solved"
					benefit := val.Number * val.Rate
					username := val.Owner
					var lucky_dog User
					userAsByte, err := stub.GetState(username)
					if err != nil {
						return nil, errors.New("user don't exist")
					}
					json.Unmarshal(userAsByte, &lucky_dog)
					lucky_dog.Coin = lucky_dog.Coin + benefit
					userAsByte, err = json.Marshal(lucky_dog)
					stub.PutState(username, userAsByte)

				}
			}
		}

		StakeAsBytes, err = json.Marshal(Stakes)
		stub.PutState(ActiveStakeStr, StakeAsBytes)
	}

	return nil, nil
}

