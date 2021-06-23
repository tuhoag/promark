package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/bwesterb/go-ristretto"
	redis "gopkg.in/redis.v4"
)

var f *os.File

type Campaign struct {
	ID string `json:"id"`
	No int    `json:"ver"`
	H  byte   `json:`
}

type campaign_request struct {
	ID string `json:"id"`
	No int    `json:"no"`
}

type campaign_param struct {
	H  []byte `json:"hvalue"`
	R1 []byte `json:"r1"`
	R2 []byte `json:"r2"`
}

func main() {
	var err error
	f, err = os.Create("extlog")

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", homeHandler)
	//http.HandleFunc("/post", postHandler)
	http.HandleFunc("/camp", campaignParams)

	// redisConnect()

	http.ListenAndServe(":5000", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

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

func campaignParams(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		println(err)
	}
	log.Println(string(body))

	var campaign campaign_request
	var campParam campaign_param
	err = json.Unmarshal(body, &campaign)
	if err != nil {
		println(err)
	}
	log.Println("campaign is:", campaign.ID, campaign.No)

	// generate and set param
	setParam(campaign.ID, campaign.No)
	//get param from db
	campParam = getParam(campaign.ID)

	//temporary return
	param, err := json.Marshal(campParam)

	// for log
	n, err := f.WriteString(string(param))

	fmt.Println("wrote to file:", n)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(rw, string(param))
}

func redisConnect() {
	// client := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })

	// pong, err := client.Ping().Result()
	// fmt.Println("pong:::::", pong, err)

	// // Store to db
	// err = client.Set("name", "Elliot", 0).Err()
	// // if there has been an error setting the value
	// // handle the error
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// //read from db
	// val, err := client.Get("name").Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(val)
}

// Campaign function part
func setParam(id string, no int) {
	var r1, r2 ristretto.Scalar

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	//generate campaign param
	_, err := client.Get(id).Result()

	if err != nil {
		fmt.Println(err)

		H := generateH()
		hBytes := H.Bytes()
		// hString := convertPointToString(H)
		fmt.Println("hString:.\n", hBytes)

		r1.Rand()
		// r1String := string(r1.Bytes())
		r1Bytes := r1.Bytes()
		fmt.Println("r1:.\n", r1Bytes)

		r2.Rand()
		// r2String := string(r2.Bytes())
		r2Bytes := r2.Bytes()
		fmt.Println("r1:.\n", r2Bytes)

		jsonParam, err := json.Marshal(campaign_param{H: hBytes, R1: r1Bytes, R2: r2Bytes})
		if err != nil {
			fmt.Println(err)
		}

		//store to redis db
		err = client.Set(id, jsonParam, 0).Err()
		if err != nil {
			fmt.Println(err)
		}

	} else {
		fmt.Println("The campaign already existed.\n")
	}
}

func getParam(id string) campaign_param {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	val, err := client.Get(id).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)

	var campaign campaign_param
	err = json.Unmarshal([]byte(val), &campaign)
	if err != nil {
		println(err)
	}

	//test the value of H
	H1 := convertBytesToPoint(campaign.H)
	fmt.Println("H1 point:", H1)

	return campaign
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

//////////// Pedersen function
var n25519, _ = new(big.Int).SetString("7237005577332262213973186563042994240857116359379907606001950938285454250989", 10)

// Generate a random point on the curve
func generateH() ristretto.Point {
	var random ristretto.Scalar
	var H ristretto.Point
	random.Rand()
	H.ScalarMultBase(&random)

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
