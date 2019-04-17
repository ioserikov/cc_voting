package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	id "github.com/s7techlab/cckit/identity" // s7 techlab MIT license
)

const consttrimstring = "{}[]"      // symbols for clearing
const constcsvseparator = rune(';') // separator for csv

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct{}

// struct for voters in chain.  VoteID - key for this value
type votelist struct {
	repourl string
	EndDate string // TODO - Date. TxTime is needed to
	voters  []string
}

// struct for voting
type voterslist struct {
	voter   string // owner of vote - it.s certname is in key
	vote    string // vote
	comment string // comment to voter
}

// struct for report of vote
type votereport struct {
	voteID     string       // Hash of vote
	repourl    string       // link to repo with add data to vote
	voteresult string       // result according to votes
	votestr    []voterslist // strings of votes
}

// struct for csv fle columns
type csvrow struct { // struct for string of csv file
	voteid      string
	voterepo    string
	voteenddate string
	voter       string
	vote        string
	result      string
	comment     string
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init - run on instantiate \ update
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// to us  - fix metadata of Initialization
	// certcreattor, _ := id.FromStub(stub)

	// invoker, _ := toBytes(certcreattor.Cert.Issuer.CommonName)
	// key := "Init of chaincode"
	// stub.PutState(key, invoker)

	return shim.Success(nil)

}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	if strings.ToLower(function) == "votestart" { //vote starting
		return t.votestart(stub, args)
	} else if strings.ToLower(function) == "voteend" { // Not Implemented yet
		return t.votend(stub, args)
	} else if strings.ToLower(function) == "voteresult" { // only result returning
		return t.voteresult(stub, args)
	} else if strings.ToLower(function) == "votehistory" { //Not Implemented Yet
		return t.votehistory(stub, args)
	} else if strings.ToLower(function) == "vote" { // vote and save to ledger
		return t.vote(stub, args)
	} else if strings.ToLower(function) == "voteresultcsv" { // vote and save to ledger
		return t.voteresultcsv(stub, args)
	}

	return shim.Success(nil)
}

// func inside Invoke Routing
func (t *SimpleChaincode) votestart(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	argnum := len(args)

	if argnum < 4 {
		return shim.Error("Too less arguments")
	}

	voteid := args[0] // Hash ID of voting

	votelist := votelistfromargs(args) // parsing to struct from args - for convience

	if err := checkuniqvoters(votelist.voters); err != nil {
		return shim.Error(" duplicate voters")
	}

	value, _ := toBytes(*votelist)      //struct as []byte
	err := stub.PutState(voteid, value) // simple save in ledger - check for uniq

	if err != nil {
		return shim.Error("some err")
	}

	return shim.Success(nil)

} // votestart

func (t *SimpleChaincode) vote(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	voteID := args[0]

	if _, err := stub.GetState(voteID); err != nil {
		return shim.Error(" no such voting")
	}

	certident, _ := id.FromStub(stub) // identity of ivoker
	voter := certident.Cert.Issuer.CommonName
	votekey := voteID + "|" + voter // key of our chaininput tx

	vote := args[1]
	if !checkvotes(vote) { // vote is to be one of {yes, no, neutral}
		return shim.Error("can't determine vote")
	}

	comment := args[2]

	votestr := voterslist{
		voter, // duplicate of voter ()
		vote,
		comment,
	}

	value, _ := toBytes(votestr) // own marshaller

	//fmt.Printf("\n %v votestr: ", votestr)
	stub.PutState(votekey, value)

	return shim.Success(nil)
} // vote

// not implemented yet
func (t *SimpleChaincode) votend(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) voteresult(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	voteid := args[0] // Key of vote, also part of a key of report
	voteidrwkey := voteid + "|result"
	votebyte, _ := stub.GetState(voteid)    //metadata & voters
	votestruct := bytesToVoteList(votebyte) // unmarshal to struct tested!!!

	results := map[string]int{} // map fo counting our results
	votearr := []voterslist{}   // empty struct of {voter , vote, comment }

	for _, v := range votestruct.voters { // for every voter in our list of voters

		key := voteid + "|" + v // !!! key for vote is key + "|voter"
		val, _ := stub.GetState(key)
		if val != nil { // if we have record aboot voter in ledger
			strarr := strings.Split(fmt.Sprintf("\n voter %s", val), " ")
			voter := v
			voice := strings.ToLower(strarr[3])
			comment := strings.Join(strarr[4:], " ")
			was, _ := results[voice]

			votebuf := voterslist{
				"\nVoter: " + voter,
				"\nVoice: " + voice,
				"\nComment: " + comment,
			}
			votearr = append(votearr, votebuf)

			results[voice] = was + 1

		} //v != nil

	} //or i := range votestruct.voters

	result, err := resultfrommap(results)

	voterep := votereport{
		"\n Vote ID: " + voteid,
		"\nRepo Url: " + votestruct.repourl,
		"\nResolution of Voting: " + result + "\n",
		votearr,
	}

	banswer, _ := toBytes(voterep)
	stub.PutState(voteidrwkey, banswer) // log resultcall in chain for future

	if err != nil {
		return shim.Error("error counting result")
	}
	return shim.Success(banswer)

} //voteresult

func (t *SimpleChaincode) voteresultcsv(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var csvfile []csvrow // file is array of csv strings

	buf := new(bytes.Buffer)

	// csv writer options
	w := csv.NewWriter(buf)
	w.UseCRLF = true
	w.Comma = constcsvseparator

	// copy of voteresult to gater in csv

	voteid := args[0] // Key of vote, also part of a key of report
	//	voteidrwkey := voteid + "|result"
	votebyte, _ := stub.GetState(voteid)    //metadata & voters
	votestruct := bytesToVoteList(votebyte) // unmarshal to struct tested!!!

	results := map[string]int{} // map fo counting our results
	votearr := []voterslist{}   // empty struct of {voter , vote, comment }

	for _, v := range votestruct.voters { // for every voter in our list of voters

		key := voteid + "|" + v // !!! key for vote is key + "|voter"
		val, _ := stub.GetState(key)

		if val != nil { // if we have record about voter in ledger
			strarr := strings.Split(fmt.Sprintf("\n voter %s", val), " ")
			voter := v
			voice := strings.ToLower(strarr[3])
			comment := strings.Join(strarr[4:], " ")
			//	comment = string(comment[0:10])
			was, _ := results[voice]

			votebuf := voterslist{
				voter,
				voice,
				comment,
			}
			votearr = append(votearr, votebuf)

			results[voice] = was + 1

			// combine a csv string

			csvbuf := csvrow{
				voteid,
				votestruct.repourl,
				votestruct.EndDate,
				voter,
				voice,
				" ",
				comment,
			}

			csvfile = append(csvfile, csvbuf)
		} //v != nil

	} //or i := range votestruct.voters

	result, err := resultfrommap(results)

	for _, v := range csvfile {
		v.result = result

		strbuf := []string{}
		strbuf = append(strbuf, v.voteid)
		strbuf = append(strbuf, v.voterepo)
		strbuf = append(strbuf, v.voteenddate)
		strbuf = append(strbuf, v.voter)
		strbuf = append(strbuf, v.vote)
		strbuf = append(strbuf, v.result)
		strbuf = append(strbuf, v.comment)

		// clear from { } [ ]
		strbuf = clearfromtyped(strbuf)

		w.Write(strbuf) // append to inter
	}

	w.Flush() // write to buf

	banswer, _ := toBytes(buf)

	if err != nil {
		return shim.Error("error counting result")
	}
	return shim.Success(banswer)
} // voteresultcsv

//  not implmemented yet
func (t *SimpleChaincode) votehistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	key := args[0]
	resp, err := stub.GetState(key)

	if resp == nil || err != nil {
		return shim.Error("error getting history")

	}

	return shim.Success(resp)

}

//----------------- additional funcs

// write orgs in votestr.voters field (slice)
func votelistfromargs(args []string) *votelist {

	votestr := new(votelist)

	votestr.repourl = args[1]
	votestr.EndDate = args[2]

	for i := 3; i < len(args); i++ {
		votestr.voters = append(votestr.voters, args[i])

	}

	return votestr
}

// check the vote is inside {Yes, No, Neutral }
func checkvotes(vote string) bool {
	res := false

	res = (strings.ToLower(vote) == "yes") || (strings.ToLower(vote) == "no") || (strings.ToLower(vote) == "neutral")

	return res

}

// []byte to votelist unmarshal
func bytesToVoteList(votestr []byte) *votelist { //tested

	//var buf bytes.Buffer
	var lvotelist votelist

	strslice := strings.Split(string(votestr), " ")
	strbuf := clearfromtyped(strslice[2:])

	lvotelist.repourl = strbuf[0]
	lvotelist.EndDate = strslice[1]

	for _, v := range strbuf {
		lvotelist.voters = append(lvotelist.voters, v)
	}
	//}

	return &lvotelist

} // bytesToVoteList

// ToBytes converts inteface{} (string, []byte , struct to ToByter interface to []byte for storing in state
// from s7techlab cckit + refactor
func toBytes(value interface{}) ([]byte, error) {
	if value == nil {
		return nil, nil
	}

	switch value.(type) {

	// first priority if value implements ToByter interface

	case proto.Message:
		return proto.Marshal(proto.Clone(value.(proto.Message)))
	case bool:
		return []byte(strconv.FormatBool(value.(bool))), nil
	case string:
		return []byte(value.(string)), nil
	case uint:
		return []byte(fmt.Sprint(value.(uint))), nil
	case int:
		return []byte(fmt.Sprint(value.(int))), nil
	case int32:
		return []byte(fmt.Sprint(value.(int32))), nil
	case []byte:
		return value.([]byte), nil

	default:
		valueType := reflect.TypeOf(value).Kind()

		switch valueType {
		case reflect.Ptr:
			fallthrough
		case reflect.Struct:

			return []byte(fmt.Sprintf("%v", value)), nil

			//return json.Marshal(value)
		case reflect.Array:
			return []byte(fmt.Sprintf("%v", value)), nil
		case reflect.Map:
			return []byte(fmt.Sprintf("%v", value)), nil
		case reflect.Slice:
			return []byte(fmt.Sprintf("%v", value)), nil
		//	return json.Marshal(value)
		// used when type based on string
		case reflect.String:
			return []byte(reflect.ValueOf(value).String()), nil

		default:
			return nil, fmt.Errorf(
				`toBytes converting supports ToByter interface,struct,array,slice,bool and string, current type is %s`,
				valueType)
		}

	}
}

//count resolution ov voting
func resultfrommap(res map[string]int) (string, error) {
	resolution := ""
	err := fmt.Errorf("Error during count results ")

	yesvotes, _ := res["yes"]
	novotes, _ := res["no"]
	delta := yesvotes - novotes
	switch {

	case delta == 0:
		{
			err = nil
			resolution = "equal"
		}
	case delta < 0:
		{
			err = nil
			resolution = "no"
		}
	case delta > 0:
		{
			err = nil
			resolution = "yes"
		}
	}

	return resolution, err
}

// check list of voters on duplicates absence
func checkuniqvoters(voters []string) error { //tested

	strsearch := strings.Join(voters, " ")
	uniq := true
	for _, v := range voters {

		if uniq = !(strings.Count(strsearch, v) > 1); uniq != true {
			break
		}

	}

	if uniq != true {
		return fmt.Errorf("There are duplicate voters in Your list") // negative
	}
	return nil // positive
}

// clear form symbols like { } []
func clearfromtyped(strsrc []string) []string {
	strdest := []string{}

	for _, v := range strsrc {
		str := strings.Trim(v, consttrimstring)
		strdest = append(strdest, str)
	}

	return strdest

}
