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
	S  []byte `json:"s"`
}

type CommValue struct {
	ID   string
	comm []byte
}

type Cam struct {
	ID string
	No int
}

type CommParam struct {
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
var sValue SecretNumber

// var (
// 	extURL        = "http://external.promark.com:5000"
// 	camRequestURL = extURL + "/camp"
// )

func main() {
	var err error
	var fname string
	var port string
	fname = os.Getenv("VER_NAME")
	port = os.Getenv("VER_PORT")
	fmt.Println(fname)
	f, err = os.Create(fname)

	if err != nil {
		panic(err)
	}

	// initilization
	id = "00"
	sValue = getSecretNumber(id)

	// n, _ := f.WriteString("Secret value is:" + string(sValue.S) + string("\n"))
	// fmt.Println(n)

	http.HandleFunc("/", home)
	http.HandleFunc("/comm", computeComm)
	http.ListenAndServe(":"+port, nil)
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

	n, _ := f.WriteString("Welcome to verifying service" + string("\n"))
	fmt.Println(n)

	fmt.Fprintf(w, "Hello")
}

func computeComm(rw http.ResponseWriter, req *http.Request) {
	var r ristretto.Scalar
	var V ristretto.Scalar
	var comm ristretto.Point

	n, err := f.WriteString("computeComm() calling" + string("\n"))
	fmt.Println(n)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		println(err)
	}
	log.Println(string(body))

	var commParam CommParam
	err = json.Unmarshal(body, &commParam)

	hDec, _ := b64.StdEncoding.DecodeString(commParam.Hbyte)
	tmp := convertBytesToPoint(hDec)

	rDec, _ := b64.StdEncoding.DecodeString(commParam.R)
	r = convertBytesToScalar(rDec)
	// r.Rand()
	rstring := string(body)

	n, err = f.WriteString("body:" + rstring + string("\n"))

	fmt.Println(n)

	//compute the comm value
	V = convertBytesToScalar(sValue.S)

	//get the value of H
	comm = commitTo(&tmp, &r, &V)
	cEnc := b64.StdEncoding.EncodeToString(comm.Bytes())

	n1, err := f.WriteString("cEnc:" + cEnc + string("\n"))
	fmt.Println(n1)

	fmt.Fprintf(rw, cEnc)
}

func getSecretNumber(id string) SecretNumber {
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

		returnValue = SecretNumber{ID: id, S: s.Bytes()}

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
	fmt.Println("in convertBytesToPoint result:", result)

	return H
}

func convertStringToScalar(s string) ristretto.Scalar {

	var r ristretto.Scalar
	var rBytes [32]byte

	tmp := []byte(s)
	copy(rBytes[:32], tmp[:])

	result := r.SetBytes(&rBytes)
	fmt.Println("in convertStringToScalar result:", result)

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
