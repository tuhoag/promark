package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"

	"github.com/bwesterb/go-ristretto"
)

type campaign_param struct {
	id string `json:"id"`
	H  []byte `json:"hvalue"`
	R  string `json:"r"`
}

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":5002", nil)
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

	v = rand.Int63n(100)

	err = json.Unmarshal(body, &campParam)
	if err != nil {
		println(err)
	}
	log.Println("campaign is:", convertBytesToPoint(campParam.H), campParam.R)

	tem := big.NewInt(v)

	//get the value of H
	H := convertBytesToPoint(campParam.H)
	r = convertStringToScalar(campParam.R)
	fmt.Println("H1 point:", H)

	comm = commitTo(&H, &r, V.SetBigInt(tem))
	fmt.Print("comm: \n", comm)

	s := comm.String()

	fmt.Fprintf(w, string(s))
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

	// fmt.Println("in convertStringToPoint s: \n", s)

	tmp := []byte(s)
	// fmt.Println("in convertStringToPoint tmp len: \n", tmp)

	copy(hBytes[:32], tmp[:32])

	H.SetBytes(&hBytes)
	// fmt.Println("in convertStringToPoint H reverted: ", H)

	return H
}

func convertPointToString(H ristretto.Point) string {
	// fmt.Println(string("in convertPointToString function"))

	// H := generateH()
	fmt.Println("H: ", H)

	hBytes := H.Bytes()
	hString := string(hBytes)
	//fmt.Println("in convertPointtoBytes H in bytes: ", hBytes)
	//fmt.Println("in convertPointtoBytes H to string: ", H.String())
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
