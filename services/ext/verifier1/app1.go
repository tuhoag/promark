package main

import (
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

type SecretNumber struct {
	ID string `json:"id"`
	S  string `json:"s"`
}

type CommValue struct {
	ID   string
	comm []byte
}

var id string

var f *os.File

func main() {
	var err error
	f, err = os.Create("ver1log")

	// initilization
	id = "00"
	getSecretNumber()

	if err != nil {
		panic(err)
	}
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

func computeComm(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		println(err)
	}
	log.Println(string(body))

	var campParam campaign_param
	var r, V ristretto.Scalar
	var v int64
	var comm ristretto.Point

	// generate COMM value - start

	err = json.Unmarshal(body, &campParam)
	if err != nil {
		println(err)
	}

	v = rand.Int63n(100)

	tem := big.NewInt(v)

	//get the value of H
	H := convertBytesToPoint(campParam.H)
	r = convertStringToScalar(campParam.R)
	n, err := f.WriteString(campParam.R)

	fmt.Println("H point:", H.Bytes())
	// n, err := f.WriteString(string(H.Bytes())

	comm = commitTo(&H, &r, V.SetBigInt(tem))

	commBytes := comm.Bytes()

	// for log
	n1, err := f.WriteString(string(commBytes))

	fmt.Println("wrote to file:", n, n1)

	if err != nil {
		panic(err)
	}

	cTemp := CommValue{campParam.id, commBytes}

	jsonData, err := json.Marshal(cTemp)

	// generate COMM - end
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(jsonData))
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
