package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	id "github.com/s7techlab/cckit/identity"
	cckit "github.com/s7techlab/cckit/testing"
)

const stubsert1 = `-----BEGIN CERTIFICATE-----
MIICNjCCAd2gAwIBAgIRAMnf9/dmV9RvCCVw9pZQUfUwCgYIKoZIzj0EAwIwgYEx
CzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4g
RnJhbmNpc2NvMRkwFwYDVQQKExBvcmcxLmV4YW1wbGUuY29tMQwwCgYDVQQLEwND
T1AxHDAaBgNVBAMTE2NhLm9yZzEuZXhhbXBsZS5jb20wHhcNMTcxMTEyMTM0MTEx
WhcNMjcxMTEwMTM0MTExWjBpMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZv
cm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEMMAoGA1UECxMDQ09QMR8wHQYD
VQQDExZwZWVyMC5vcmcxLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0D
AQcDQgAEZ8S4V71OBJpyMIVZdwYdFXAckItrpvSrCf0HQg40WW9XSoOOO76I+Umf
EkmTlIJXP7/AyRRSRU38oI8Ivtu4M6NNMEswDgYDVR0PAQH/BAQDAgeAMAwGA1Ud
EwEB/wQCMAAwKwYDVR0jBCQwIoAginORIhnPEFZUhXm6eWBkm7K7Zc8R4/z7LW4H
ossDlCswCgYIKoZIzj0EAwIDRwAwRAIgVikIUZzgfuFsGLQHWJUVJCU7pDaETkaz
PzFgsCiLxUACICgzJYlW7nvZxP7b6tbeu3t8mrhMXQs956mD4+BoKuNI
-----END CERTIFICATE-----
`

const stubsert2 = `-----BEGIN CERTIFICATE-----
MIIC6zCCAlSgAwIBAgIJAM8WlT9qF9WvMA0GCSqGSIb3DQEBCwUAMIGMMQswCQYD
VQQGEwJSVTETMBEGA1UECAwKU29tZS1TdGF0ZTENMAsGA1UEBwwEQ2l0eTESMBAG
A1UECgwJT3JnMi5UZXN0MQswCQYDVQQLDAJDQTEcMBoGA1UEAwwTb3JnMi5zbXBs
ZS5mb3IudGVzdDEaMBgGCSqGSIb3DQEJARYLb3JnQHRlc3QucnUwHhcNMTkwNDEy
MTIxNjA2WhcNMjAwNDExMTIxNjA2WjCBjDELMAkGA1UEBhMCUlUxEzARBgNVBAgM
ClNvbWUtU3RhdGUxDTALBgNVBAcMBENpdHkxEjAQBgNVBAoMCU9yZzIuVGVzdDEL
MAkGA1UECwwCQ0ExHDAaBgNVBAMME29yZzIuc21wbGUuZm9yLnRlc3QxGjAYBgkq
hkiG9w0BCQEWC29yZ0B0ZXN0LnJ1MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKB
gQDdeKQ/fZpBrCPp7OnxW5FNOQD7m3fa7snnm+x5Ujk7UkVAeci2rUexEq4tHK9x
Fw15nf4JABmfywRYGJUBTeGQbIOCMQIbxj8C3o2cQn1kd7ZetaO87Q0XH0NSMLsX
SoZQ5qvGcsUsm71p1RKH1ta6XT8ds/N4EXS1/sAo8WNBnQIDAQABo1MwUTAdBgNV
HQ4EFgQUUMp3l+JTn2zbyhZCEp1rLwb0gmkwHwYDVR0jBBgwFoAUUMp3l+JTn2zb
yhZCEp1rLwb0gmkwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOBgQAj
YbW6yn8IwYHnskyM1viuHVkEodgfs1Y8fzfPhSNNvmfhaWDIq4onYXC/IUfPOMe3
fIqwsFOw7loQx6qHSbfjvSi9A2eLIt5xJUgYFSwzTynDFcIDRc9W/fu0I7uii0DQ
KDUz5Sahkjtevdcrkt9fCSoiPc+KKP6f3AlETJd8Ww==
-----END CERTIFICATE-----`

const stubsert3 = `-----BEGIN CERTIFICATE-----
MIIC1zCCAkCgAwIBAgIJALT0qq+TvmLeMA0GCSqGSIb3DQEBCwUAMIGCMQswCQYD
VQQGEwJTRzELMAkGA1UECAwCU0cxCzAJBgNVBAcMAlNHMRYwFAYDVQQKDA1vcmcu
MjAxOS5yZWFsMQswCQYDVQQLDAJPVTEZMBcGA1UEAwwQc2VydmVyLmZxZG4ubmFt
ZTEZMBcGCSqGSIb3DQEJARYKb3RAdGVzdC5ydTAeFw0xOTA0MTIxMzA3MDhaFw0y
MDA0MTExMzA3MDhaMIGCMQswCQYDVQQGEwJTRzELMAkGA1UECAwCU0cxCzAJBgNV
BAcMAlNHMRYwFAYDVQQKDA1vcmcuMjAxOS5yZWFsMQswCQYDVQQLDAJPVTEZMBcG
A1UEAwwQc2VydmVyLmZxZG4ubmFtZTEZMBcGCSqGSIb3DQEJARYKb3RAdGVzdC5y
dTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEA1Tk+5udnKPYKxH7YykqMOwzT
iJs3b7o3K4M/XOVL+NIo3x9db0NhFRJAlniLFf6r/OB6uIHTEIBYZwlOn56O1gMJ
ZNiX/FOHcD3hkONKaNW6wmIwqOEWJQa/lPf7a9IW8XWbGdJ6yZms9/eNCqdEXxKo
scqW2piivkSQql9AF50CAwEAAaNTMFEwHQYDVR0OBBYEFMtHEEhfDzEfvW4EoGSM
oBwf6gB6MB8GA1UdIwQYMBaAFMtHEEhfDzEfvW4EoGSMoBwf6gB6MA8GA1UdEwEB
/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADgYEAOO0x734kh2wyAZcefL8ErPJDgl5E
n3gGW6mZdmmBtRSzf3pnBcUAHGJR1WP1huXeCZ55Bpxf+P2kSqSOXYslvBi16Rm/
8n3qqqIKsCXHA8kikAxhO5g32mN+CsFc3hS9cwmi3vJKqLGuXGTL2y4a+ZKNMD92
l2MWdBDM1G+TL1Y=
-----END CERTIFICATE-----
`

func checkInit(t *testing.T, stub *cckit.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Payload))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *cckit.MockStub, name string) {
	bytes, _ := stub.GetState(name)

	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	fmt.Printf("\nCheckState key:   %s  value:    %s  \n", name, string(bytes))

}

func checkQuery(t *testing.T, stub *cckit.MockStub, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("query"), []byte(name)})

	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != value {
		fmt.Println("Query value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *cckit.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}

	if fmt.Sprintf("%s", args[0]) == "voteresult" || fmt.Sprintf("%s", args[0]) == "votehistory" || fmt.Sprintf("%s", args[0]) == "voteresultcsv" {
		//must parse result to struct
		bur := res.GetPayload()
		burs := string(bur)
		fmt.Printf("\nmessage %s\n", burs)

	}
}

func TestExample01_Init(t *testing.T) {

	fmt.Println("test init!!")
	sc := new(SimpleChaincode)
	stub := cckit.NewMockStub("crocc", sc)

	var bufbyte [][]byte // empty strings for init args - our case

	bsert := []byte(stubsert1)
	stub.MockCreator("MSP", bsert)

	checkInit(t, stub, bufbyte)
	fmt.Println("End TestExample01_Init")
}

func TestExample02_InvokeStart(t *testing.T) {

	fmt.Println("test invoke votestart")

	sc := new(SimpleChaincode)
	stub := cckit.NewMockStub("crocc", sc)
	stub.ClearCreatorAfterInvoke = false // for saving creator for testing
	bsert := []byte(stubsert1)
	stub.MockCreator("MSP", bsert)

	funcs := "votestart"
	bfuncs := []byte(funcs)

	IDVote := "VoteHash"
	bIDVote := []byte(IDVote)
	brepourl := []byte("https://git.repo")
	bEndDate := []byte("10.04.2019.10.00")

	Orgs := []string{
		"Org1",
		"Org2",
		"Org3",
	}

	//	orgbyte := []byte(strings.Join(Orgs, " "))

	buff := [][]byte{
		bfuncs,
		bIDVote,
		brepourl,
		bEndDate,
	}

	for _, v := range Orgs {
		curbyte := []byte(v)
		buff = append(buff, curbyte)
	}

	checkInvoke(t, stub, buff)

	checkState(t, stub, IDVote)

	fmt.Println("End of TestExample02_Invoke")
}

func TestExample4_VoteStart_Vote(t *testing.T) {

	fmt.Println("begin Test 4 Votestart and Vote By one test")

	sc := new(SimpleChaincode)
	stub := cckit.NewMockStub("crocc", sc)

	sert := []byte(stubsert1)
	stub.MockCreator("MSP1", sert)
	stub.ClearCreatorAfterInvoke = false

	certident, _ := id.FromStub(stub)
	creattor := certident.Cert.Issuer.CommonName

	bfuncs := []byte("votestart")
	IDVote := "VoteHash"
	bIDVote := []byte("VoteHash")
	brepourl := []byte("https://git.repo")
	bEndDate := []byte("10.04.2019.10.00")

	Orgs := []string{ // here we need to put real idenitites
		creattor,
		"ca.Org2.example.com",
		"ca.Org3.example.com",
	}

	//	orgbyte := []byte(strings.Join(Orgs, " "))

	// buff := [][]byte{
	// 	bfuncs,
	// 	bIDVote,
	// 	brepourl,
	// 	bEndDate,
	//
	// }

	buff := [][]byte{
		bfuncs,
		bIDVote,
		brepourl,
		bEndDate,
	}

	for _, v := range Orgs {
		curbyte := []byte(v)
		buff = append(buff, curbyte)
	}

	checkInvoke(t, stub, buff)
	checkState(t, stub, IDVote)

	bfuncs = []byte("Vote")
	bVoteID := []byte("VoteHash")
	bvote := []byte("No")
	bcomment := []byte("Because I can That's Why!!")

	buff = [][]byte{
		bfuncs,
		bVoteID,
		bvote,
		bcomment,
	}

	checkInvoke(t, stub, buff)
	votekey := "VoteHash|" + creattor
	checkState(t, stub, votekey)

}

func TestExample5_VoteStart_Vote_Result(t *testing.T) {

	fmt.Println("begin Test 5 Votestart and Vote and VoteResult By one test")

	sc := new(SimpleChaincode)
	stub := cckit.NewMockStub("crocc", sc)

	sert := []byte(stubsert1)
	stub.MockCreator("MSP1", sert)
	stub.ClearCreatorAfterInvoke = false

	certident, _ := id.FromStub(stub)
	creattor := certident.Cert.Issuer.CommonName

	bfuncs := []byte("votestart")
	IDVote := "VoteHash"
	bIDVote := []byte("VoteHash")
	brepourl := []byte("https://git.repo")
	bEndDate := []byte("10.04.2019.10.00")

	Orgs := []string{ // here we need to put real idenitites
		creattor,
		"ca.Org2.example.com",
		"ca.Org3.example.com",
	}
	//
	// orgbyte := []byte(strings.Join(Orgs, " "))
	//
	// buff := [][]byte{
	// 	bfuncs,
	// 	bIDVote,
	// 	brepourl,
	// 	bEndDate,
	// 	orgbyte,
	// }

	buff := [][]byte{
		bfuncs,
		bIDVote,
		brepourl,
		bEndDate,
	}

	for _, v := range Orgs {
		curbyte := []byte(v)
		buff = append(buff, curbyte)
	}

	checkInvoke(t, stub, buff)
	checkState(t, stub, IDVote)

	bfuncs = []byte("Vote")
	bVoteID := []byte("VoteHash")
	bvote := []byte("No")
	bcomment := []byte("Because I can That's Why!!")

	buff = [][]byte{
		bfuncs,
		bVoteID,
		bvote,
		bcomment,
	}

	checkInvoke(t, stub, buff)
	votekey := "VoteHash|" + creattor
	checkState(t, stub, votekey)

	bfuncs = []byte("voteresult")
	buff = [][]byte{
		bfuncs,
		bVoteID,
	}
	checkInvoke(t, stub, buff)

}

func TestExample6_VoteStart_Vote_Result_3orgs(t *testing.T) {
	fmt.Println("begin Test 6 3 Orgs")

	sc := new(SimpleChaincode)
	stub := cckit.NewMockStub("crocc", sc)

	sert2 := []byte(stubsert2)
	stub.MockCreator("MSP1", sert2)
	stub.ClearCreatorAfterInvoke = false

	certident2, _ := id.FromStub(stub)
	creator2 := certident2.Cert.Issuer.CommonName

	sert3 := []byte(stubsert3)
	stub.MockCreator("MSP2", sert3)

	certident3, _ := id.FromStub(stub)

	creator3 := certident3.Cert.Issuer.CommonName

	sert1 := []byte(stubsert1)
	stub.MockCreator("MSP3", sert1)

	certident, _ := id.FromStub(stub)
	creattor1 := certident.Cert.Issuer.CommonName

	bfuncs := []byte("votestart")
	IDVote := "VoteHash"
	bIDVote := []byte("VoteHash")
	brepourl := []byte("https://git.repo")
	bEndDate := []byte("10.04.2019.10.00")

	Orgs := []string{ // here we need to put real idenitites
		creattor1,
		creator2,
		creator3,
	}

	buff := [][]byte{
		bfuncs,
		bIDVote,
		brepourl,
		bEndDate,
	}

	for _, v := range Orgs {
		curbyte := []byte(v)
		buff = append(buff, curbyte)
	}

	checkInvoke(t, stub, buff)
	checkState(t, stub, IDVote)

	// #1 - Voice 1
	bfuncs = []byte("Vote")
	bVoteID := []byte("VoteHash")
	bvote := []byte("Yes")
	bcomment := []byte("Because I can That's Why!!")

	buff = [][]byte{
		bfuncs,
		bVoteID,
		bvote,
		bcomment,
	}

	checkInvoke(t, stub, buff)
	votekey := "VoteHash|" + creattor1
	checkState(t, stub, votekey)

	// #2
	stub.MockCreator("MSP1", sert2)

	bfuncs = []byte("Vote")
	bVoteID = []byte("VoteHash")
	bvote = []byte("Neutral")
	bcomment = []byte("Also I can!")

	buff = [][]byte{
		bfuncs,
		bVoteID,
		bvote,
		bcomment,
	}

	checkInvoke(t, stub, buff)
	votekey = "VoteHash|" + creator2
	checkState(t, stub, votekey)
	// # 3

	stub.MockCreator("MSP1", sert3)

	bfuncs = []byte("Vote")
	bVoteID = []byte("VoteHash")
	bvote = []byte("No")
	bcomment = []byte("And I Am")

	buff = [][]byte{
		bfuncs,
		bVoteID,
		bvote,
		bcomment,
	}

	checkInvoke(t, stub, buff)
	votekey = "VoteHash|" + creator3
	checkState(t, stub, votekey)

	bfuncs = []byte("voteresultcsv")
	buff = [][]byte{
		bfuncs,
		bVoteID,
	}
	checkInvoke(t, stub, buff)

	bfuncs = []byte("votehistory")
	buff = [][]byte{
		bfuncs,
		bVoteID,
	}
	checkInvoke(t, stub, buff)

}
