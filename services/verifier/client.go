//curl -X POST -H 'Content-Type: application/json' -d "{\"Test\": \"that\"}" http://localhost:8081/test

package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/bwesterb/go-ristretto"
)

var (
	homeURL        = "http://0.0.0.0:8081"
	jsonURL        = homeURL + "/test"
	camRequestURL  = homeURL + "/camp"
	commRequestURL = homeURL + "/comm"
	testConvertURL = homeURL + "/convert"
	ver1URL        = "http://0.0.0.0:5001/comm"
)

var param campaign_param

type Data struct {
	Test string
}

type Cam struct {
	ID string
	No int
}

type TestConvert struct {
	ID    string
	Hbyte string
	r     string
}

type campaign_param struct {
	H  []byte `json:"hvalue"`
	R1 string `json:"r1"`
	R2 string `json:r2`
}

type ResultConvert struct {
	ID string
	C  []byte
}

func main() {
	fmt.Println("client app running")
	// postTest()
	// getTest()
	requestCamParams()

	// data := "abc123!?$*&()'-=@~"
	// sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	// fmt.Println(sEnc)
	// sDec, _ := b64.StdEncoding.DecodeString(sEnc)
	// fmt.Println(string(sDec))
	fmt.Println()

	// commRequest()
}

func getTest() {
	response, err := http.Get(homeURL)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	d1 := []byte("hello\n")
	err = ioutil.WriteFile("dat1", d1, 0644)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("dat2")
	n2, err := f.Write(d1)

	// n3, err := f.WriteString("writes\n")
	fmt.Println("wrote to file:", n2)

	fmt.Println(string(responseData))
}

func requestCamParams() {
	c := &http.Client{}

	message := Cam{"id4", 2}
	//jsonData := `{"Test":"that"}`

	jsonData, err := json.Marshal(message)

	request := string(jsonData)

	reqJSON, err := http.NewRequest("POST", camRequestURL, strings.NewReader(request))
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return
	}

	respJSON, err := c.Do(reqJSON)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return
	}
	defer respJSON.Body.Close()

	data, err := ioutil.ReadAll(respJSON.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return
	}

	err = json.Unmarshal([]byte(data), &param)
	if err != nil {
		println(err)
	}

	fmt.Println("return data all:", string(param.H))

	testConvertBytesToPoint()
}

func testConvertBytesToPoint() {
	c := &http.Client{}

	hBytes := param.H
	//encode data
	hEnc := b64.StdEncoding.EncodeToString(hBytes)
	rEnc := b64.StdEncoding.EncodeToString([]byte(param.R1))

	message := TestConvert{"id3", hEnc, rEnc}

	jsonData, err := json.Marshal(message)

	request := string(jsonData)

	reqJSON, err := http.NewRequest("POST", ver1URL, strings.NewReader(request))
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return
	}

	respJSON, err := c.Do(reqJSON)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return
	}
	defer respJSON.Body.Close()

	data, err := ioutil.ReadAll(respJSON.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return
	}

	// var commResult ResultConvert

	// err = json.Unmarshal([]byte(data), &commResult)

	// H := convertBytesToPoint(commResult.C)

	fmt.Println("testConvertBytesToPoint result", string(data))
}

func commRequest() {
	c := &http.Client{}

	// message := Cam{"id4", 2}

	jsonData, err := json.Marshal(param)

	request := string(jsonData)

	reqJSON, err := http.NewRequest("POST", commRequestURL, strings.NewReader(request))
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return
	}

	respJSON, err := c.Do(reqJSON)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return
	}
	defer respJSON.Body.Close()

	data, err := ioutil.ReadAll(respJSON.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return
	}

	fmt.Println("return data all:", string(data))
}

func postTest() {
	c := &http.Client{}

	// create H value
	// var H ristretto.Point
	H := generateH()

	hString := convertPointToString(H)
	fmt.Println("the value of hString is:", hString)

	message := Data{"id1"}

	jsonData, err := json.Marshal(message)

	//jsonData := `{"Test":"that"}`

	request := string(jsonData)

	reqJSON, err := http.NewRequest("POST", jsonURL, strings.NewReader(request))
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return
	}

	respJSON, err := c.Do(reqJSON)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return
	}
	defer respJSON.Body.Close()

	data, err := ioutil.ReadAll(respJSON.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return
	}

	var H1 ristretto.Point

	fmt.Printf("read resp.Body successfully:\n%v\n", string(data))
	// fmt.Println("resp data in bytes: ", data)

	H1 = convertStringToPoint(string(data))
	fmt.Println("in data to H1 point: ", &H1)
}

// Pedersen part
// The prime order of the base point is 2^252 + 27742317777372353535851937790883648493.
var n25519, _ = new(big.Int).SetString("7237005577332262213973186563042994240857116359379907606001950938285454250989", 10)

// Commit to a value x
// H - Random secondary point on the curve
// r - Private key used as blinding factor
// x - The value (number of tokens)
func commitTo(H *ristretto.Point, r, x *ristretto.Scalar) ristretto.Point {
	//ec.g.mul(r).add(H.mul(x));
	var result, rPoint, transferPoint ristretto.Point
	rPoint.ScalarMultBase(r)
	transferPoint.ScalarMult(H, x)
	result.Add(&rPoint, &transferPoint)
	return result
}

// Generate a random point on the curve
func generateH() ristretto.Point {
	var random ristretto.Scalar
	var H ristretto.Point
	random.Rand()
	H.ScalarMultBase(&random)

	return H
}

func convertPointToString(H ristretto.Point) string {

	// H := generateH()
	fmt.Println("in convertPointToString H: ", H)

	hBytes := H.Bytes()
	fmt.Println("in convertPointToString hBytes: ", hBytes)

	hString := string(hBytes)
	//fmt.Println("in convertPointToString hString: ", hString)

	return hString
}

func convertStringToPoint(s string) ristretto.Point {

	var H ristretto.Point
	var hBytes [32]byte

	fmt.Println("in convertStringToPoint s: ", s)

	tmp := []byte(s)
	fmt.Println("in convertStringToPoint len(tmp): ", len(tmp))

	fmt.Println("in convertStringToPoint s: ", s)
	fmt.Println("in convertStringToPoint tmp: ", tmp)

	copy(hBytes[:32], tmp[:])

	result := H.SetBytes(&hBytes)

	fmt.Println("in convertStringToPoint Hbytes: ", tmp)
	fmt.Println("in convertStringToPoint result:", result)
	// fmt.Println("in convertStringToPoint H reverted: ", H)

	return H
}

func convertBytesToPoint(b []byte) ristretto.Point {

	var H ristretto.Point
	var hBytes [32]byte

	copy(hBytes[:32], b[:])

	result := H.SetBytes(&hBytes)
	fmt.Println("in convertStringToPoint result:", result)

	return H
}
