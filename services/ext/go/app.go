package main

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/bwesterb/go-ristretto"
	redis "gopkg.in/redis.v4"
)

// The prime order of the base point is 2^252 + 27742317777372353535851937790883648493.
var n25519, _ = new(big.Int).SetString("7237005577332262213973186563042994240857116359379907606001950938285454250989", 10)

func main() {
	http.HandleFunc("/", homeHandler)
	//http.HandleFunc("/post", postHandler)
	redisConnect()

	http.ListenAndServe(":5000", nil)
}

func redisConnect() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println("pong:::::", pong, err)

	// Store to db
	err = client.Set("name", "Elliot", 0).Err()
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		fmt.Println(err)
	}

	//read from db
	val, err := client.Get("name").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// hParam := generateParams()

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// fmt.Println("in generateParams hString: ", hParam)
	fmt.Fprintf(w, "Hello")
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "POST request successful")
	name := r.FormValue("username")

	fmt.Fprintf(w, "username = %s\n", name)
}

func generateParams() string {
	fmt.Println(string("in generateParams function"))

	H := generateH()
	fmt.Println("in generateParams H: \n", H)

	hBytes := H.Bytes()
	hString := string(hBytes)
	fmt.Println("in convertPointtoBytes H in bytes: \n", hBytes)
	//fmt.Println("in convertPointtoBytes H to string: ", H.String())
	//fmt.Println("in generateParams hString: ", hString)

	for i := 0; i < len(hString); i++ {
		fmt.Printf("%x ", hString[i])
	}

	fmt.Sprintf("Type H: %T\n", H)
	fmt.Sprintf("Type hBytes: %T\n", hBytes)

	return hString
}

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
