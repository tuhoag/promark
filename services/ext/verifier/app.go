package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/bwesterb/go-ristretto"
	redis "gopkg.in/redis.v4"
)

type campaign_param struct {
	id string `json:"id"`
	H  []byte `json:"hvalue"`
	R  string `json:"r"`
}

type ext_param struct {
	H  []byte `json:"hvalue"`
	R1 string `json:"r1"`
	R2 string `json:"r2"`
}

type SecretNumber struct {
	ID string `json:"id"`
	S  string `json:"s"`
}

type CommValue struct {
	ID   string
	comm []byte
}

type Cam struct {
	ID string
	No int
}

type TestConvert struct {
	ID    string `json:"ID"`
	Hbyte string `json:"Hbyte"`
	R     string `json:"R"`
}

type ResultConvert struct {
	ID string
	C  []byte
}

var id string
var f *os.File
var (
	extURL        = "http://external.promark.com:5000"
	camRequestURL = extURL + "/camp"
)

func main() {
	var err error
	var fname string
	fname = os.Getenv("VER_NAME")
	fmt.Println(fname)
	f, err = os.Create(fname)

	if err != nil {
		panic(err)
	}

	// initilization
	id = "00"
	getSecretNumber()

	http.HandleFunc("/", home)
	http.HandleFunc("/comm", computeComm)
	http.ListenAndServe(":5001", nil)
}

func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello")
}

func computeComm2(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		println(err)
	}
	log.Println(string(body))

	var campParam campaign_param
	// var r, V ristretto.Scalar
	// var v int64
	// var comm ristretto.Point

	// generate COMM value - start
	err = json.Unmarshal(body, &campParam)
	if err != nil {
		println(err)
	}

	// v = rand.Int63n(100)

	// tem := big.NewInt(v)

	//get the value of H
	// H := convertBytesToPoint(campParam.H)
	// r = convertStringToScalar(campParam.R)
	n, err := f.WriteString(string(campParam.H))
	n1, err := f.WriteString(campParam.R)

	fmt.Println("wrote to file:", n, n1)

	// fmt.Println("H point:", H.Bytes())
	// n, err := f.WriteString(string(H.Bytes())

	// comm = commitTo(&H, &r, V.SetBigInt(tem))

	// commBytes := comm.Bytes()

	if err != nil {
		panic(err)
	}

	// cTemp := CommValue{campParam.id, commBytes}

	// jsonData, err := json.Marshal(cTemp)

	// // generate COMM - end
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Fprintf(w, string(jsonData))
	fmt.Fprintf(w, "OK")
}

func computeComm1(w http.ResponseWriter, req *http.Request) {
	var campParam campaign_param

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		println(err)
	}
	log.Println(string(body))

	err = json.Unmarshal(body, &campParam)
	if err != nil {
		println(err)
	}

	// For connect to ext service
	c := &http.Client{}
	message := Cam{campParam.id, 2}
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
	var camParamExt ext_param
	fmt.Println("return data all:", string(data))
	err = json.Unmarshal([]byte(data), &camParamExt)
	if err != nil {
		println(err)
	}
	// end calling to external service

	// comm compute
	var r ristretto.Scalar
	var V ristretto.Scalar

	r.Rand()
	v := rand.Int63n(100)
	tem := big.NewInt(v)
	tmp := convertBytesToPoint(camParamExt.H)
	n, err := f.WriteString(string(camParamExt.H))

	//get the value of H
	var comm ristretto.Point
	comm = commitTo(&tmp, &r, V.SetBigInt(tem))
	commBytes := comm.Bytes()

	cTemp := CommValue{campParam.id, commBytes}

	jsonData1, err := json.Marshal(cTemp)
	// end of computation

	n, err = f.WriteString(string(jsonData1))
	fmt.Println(n)

	fmt.Fprintf(w, string(jsonData1))
}

func computeComm(rw http.ResponseWriter, req *http.Request) {
	var r ristretto.Scalar
	var V ristretto.Scalar
	var comm ristretto.Point

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		println(err)
	}
	log.Println(string(body))

	var testConvert TestConvert
	err = json.Unmarshal(body, &testConvert)

	hDec, _ := b64.StdEncoding.DecodeString(testConvert.Hbyte)
	tmp := convertBytesToPoint(hDec)

	rDec, _ := b64.StdEncoding.DecodeString(testConvert.R)

	r = convertStringToScalar(string(rDec))
	// r.Rand()
	rstring := string(body)

	n, err := f.WriteString("body:" + rstring + string("\n"))
	n, err = f.WriteString("testConvert R: " + string(rDec) + string("\n"))

	fmt.Println(n)

	//compute the comm value

	v := rand.Int63n(100)
	tem := big.NewInt(v)

	//get the value of H
	comm = commitTo(&tmp, &r, V.SetBigInt(tem))
	cEnc := b64.StdEncoding.EncodeToString(comm.Bytes())

	// returnValue := ResultConvert{testConvert.ID, comm.Bytes()}

	// jsonD, err := json.Marshal(returnValue)

	// // n, err := f.WriteString(string(jsonD))
	// fmt.Println(n)
	n1, err := f.WriteString("cEnc:" + cEnc + string("\n"))
	fmt.Println(n1)

	fmt.Fprintf(rw, cEnc)
}

func getSecretNumber() SecretNumber {
	var V ristretto.Scalar
	var returnValue SecretNumber

	// connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	//generate campaign param
	val, err := client.Get(id).Result()
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal([]byte(val), &returnValue)

	tem := rand.Int63n(100)
	v := big.NewInt(tem)

	if err != nil {
		fmt.Println(err)
		s := V.SetBigInt(v)

		returnValue = SecretNumber{ID: id, S: s.String()}

		jsonParam, err := json.Marshal(returnValue)
		if err != nil {
			fmt.Println(err)
		}

		//store to redis db
		err = client.Set(id, jsonParam, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
	}
	return returnValue
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
	fmt.Println("generateH", H)
	return H
}

func convertStringToPoint(s string) ristretto.Point {

	var H ristretto.Point
	var hBytes [32]byte

	tmp := []byte(s)

	copy(hBytes[:32], tmp[:32])
	H.SetBytes(&hBytes)

	return H
}

func convertPointToString(H ristretto.Point) string {
	// H := generateH()
	fmt.Println("H: ", H)

	hBytes := H.Bytes()
	hString := string(hBytes)
	fmt.Println("in convertPointtoString hString: ", hString)

	return hString
}

func convertBytesToPoint(b []byte) ristretto.Point {

	var H ristretto.Point
	var hBytes [32]byte

	copy(hBytes[:32], b[:])

	result := H.SetBytes(&hBytes)
	fmt.Println("in convertStringToPoint result:", result)

	return H
}

func convertStringToScalar(s string) ristretto.Scalar {

	var r ristretto.Scalar
	var rBytes [32]byte

	tmp := []byte(s)
	copy(rBytes[:32], tmp[:])

	result := r.SetBytes(&rBytes)
	fmt.Println("in convertBytesToScalar result:", result)

	return r
}

func convertScalarToString(s ristretto.Scalar) string {
	sBytes := s.Bytes()
	sString := string(sBytes)

	return sString
}

func convertBytesToScalar(b []byte) ristretto.Scalar {
	var r ristretto.Scalar
	var rBytes [32]byte

	copy(rBytes[:32], b[:])

	result := r.SetBytes(&rBytes)
	fmt.Println("in convertBytesToScalar result:", result)

	return r
}
