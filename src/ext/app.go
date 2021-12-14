package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bwesterb/go-ristretto"
	"github.com/gorilla/mux"
	redis "gopkg.in/redis.v4"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
)

var f *os.File

type CampaignCryptoRequest struct {
	CamId        string
	CustomerId   string
	NumVerifiers int
}

type CampaignCryptoParams struct {
	CamID string `json:camId`
	H     string `json:"h"`
	// R1 [][]byte `json:"r1"`
	// R2 []byte `json:"r2"`
}

func main() {
	var err error
	f, err = os.Create("extlog")

	if err != nil {
		panic(err)
	}

	var port string
	port = os.Getenv("API_PORT")

	fmt.Printf("Starting to listen on port: %s\n", port)
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homeHandler)
	myRouter.HandleFunc("/camp", createCampaignCryptoParamsHandler).Methods("POST")
	myRouter.HandleFunc("/camp/{id}", getCampaignCryptoParamsHandler)

	// myRouter.HandleFunc("/camp", createVerifierCampaignCryptoParamsHandler).Methods("POST")
	// myRouter.HandleFunc("/computeCommitment", computeCommitmentHandler)
	// myRouter.HandleFunc("/comm", computeComm)
	// myRouter.HandleFunc("/verify", verifyComm)
	// log.Fatal(http.ListenAndServe(":"+port, myRouter))

	// http.HandleFunc("/camp", createCampaignCryptoParamsHandler)
	// http.HandleFunc("/usercamp", userParamsGeneratorHandler)
	// http.HandleFunc("/data", getAllDataHandler)

	// redisConnect()
	log.Fatal(http.ListenAndServe(":"+port, myRouter))
}

func GetRedisConnection() *redis.Client {
	// pool := redis.ConnectionPool(host="127.0.0.1", port=6379, db=0)
	// client := redis.StrictRedis(connection_pool=pool)
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		PoolSize: 10000,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Errorf("ERROR: %s", err)
		f.WriteString("ERROR: " + err.Error())

		return nil
	}
	fmt.Println("pong:" + string(pong))
	f.WriteString("pong:" + string(pong) + "\n")
	return client
}

func getAllDataHandler(w http.ResponseWriter, r *http.Request) {
	client := GetRedisConnection()

	// keys, err := redis.Strings(cn.Do("KEYS", "*"))
	var cursor uint64
	var n int
	for {
		var keys []string
		var err error
		keys, cursor, err = client.Scan(cursor, "", 10).Result()
		if err != nil {
			panic(err)
		}
		n += len(keys)
		if cursor == 0 {
			break
		}

		for _, key := range keys {
			val, _ := client.Get(key).Result()
			fmt.Fprintf(w, key+":"+val)

			// values = append(values, val)
		}

	}

	// return keys, values

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	port := os.Getenv("API_PORT")
	// fmt.Println("in generateParams hString: ", hParam)
	fmt.Fprintf(w, "crypto-service.promark.com:"+port)
	f.WriteString("home")
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

func getCampaignCryptoParamsHandler(w http.ResponseWriter, r *http.Request) {
	n, _ := f.WriteString("getCampaignCryptoParamsHandler() calling" + string("\n"))
	fmt.Println(n)

	vars := mux.Vars(r)
	camId := vars["id"]

	cryptoParams, err := getCampaignCryptoParams(camId)

	if err != nil {
		fmt.Errorf("Cannot get cryptoparams %s", err)
	}

	json.NewEncoder(w).Encode(cryptoParams)
}

func createCampaignCryptoParamsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err)
	}
	log.Println(string(body))

	f.WriteString("request params: " + string(body) + string("\n"))

	var request CampaignCryptoRequest
	// var cryptoParams CampaignCryptoParams
	err = json.Unmarshal(body, &request)
	if err != nil {
		println(err)
	}

	log.Println("campaign is:", request.CamId)

	// generate and set param
	createCampaignCryptoParams(request.CamId)

	//get param from db
	cryptoParams, err := getCampaignCryptoParams(request.CamId)

	//temporary return
	param, err := json.Marshal(&cryptoParams)

	// for log
	f.WriteString(string(param))

	// fmt.Println("wrote to file:", n)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(param))
}

func userParamsGeneratorHandler(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		println(err)
	}
	log.Println(string(body))

	f.WriteString("request params: " + string(body) + string("\n"))

	var request CampaignCryptoRequest
	// var cryptoParams CampaignCryptoParams
	err = json.Unmarshal(body, &request)
	if err != nil {
		println(err)
	}

	log.Println("campaign is:", request.CamId, request.NumVerifiers)

	// generate and set param
	createCampaignCryptoParams(request.CamId)
	//get param from db
	cryptoParams, err := getCampaignCryptoParams(request.CamId)

	//temporary return
	param, err := json.Marshal(&cryptoParams)

	// for log
	f.WriteString(string(param))

	// fmt.Println("wrote to file:", n)

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
func createCampaignCryptoParams(camId string) (*CampaignCryptoParams, error) {
	// var r1, r2 ristretto.Scalar
	// var r ristretto.Scalar
	// var rArr [][]byte

	client := GetRedisConnection()

	//generate campaign param
	_, err := client.Get(camId).Result()
	var cryptoParams CampaignCryptoParams

	if err != nil {
		fmt.Println(err)

		H := generateH()
		hBytes := H.Bytes()
		hEnc := b64.StdEncoding.EncodeToString(hBytes)
		// hString := convertPointToString(H)
		fmt.Println("hString:.\n", hEnc)

		//generate "n" number of R for for n verifiers
		// for i := 0; i < no; i++ {
		// 	r.Rand()
		// 	rBytes := r.Bytes()
		// 	fmt.Println("r:.\n", rBytes)

		// 	rArr = append(rArr, rBytes)

		// 	// r1.Rand()
		// 	// // r1String := string(r1.Bytes())
		// 	// r1Bytes := r1.Bytes()
		// 	// fmt.Println("r1:.\n", r1Bytes)

		// 	// r2.Rand()
		// 	// // r2String := string(r2.Bytes())
		// 	// r2Bytes := r2.Bytes()
		// 	// fmt.Println("r1:.\n", r2Bytes)
		// }

		// jsonParam, err := json.Marshal(campaign_param{H: hBytes, R1: r1Bytes, R2: r2Bytes})
		// if err != nil {
		// 	fmt.Println(err)
		// }

		cryptoParams = CampaignCryptoParams{
			CamID: camId,
			H:     hEnc,
		}

		jsonParam, err := json.Marshal(cryptoParams)
		if err != nil {
			return nil, err
		}

		//store to redis db
		err = client.Set(camId, jsonParam, 0).Err()
		if err != nil {
			return nil, err
		}

	} else {
		fmt.Println("The campaign already existed.\n")
	}

	return &cryptoParams, nil
}

func getCampaignCryptoParams(camId string) (*CampaignCryptoParams, error) {
	client := GetRedisConnection()

	val, err := client.Get(camId).Result()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(val)

	var campaign CampaignCryptoParams
	err = json.Unmarshal([]byte(val), &campaign)
	if err != nil {
		println(err)
		return nil, err
	}

	//test the value of H
	// H1 := convertBytesToPoint(campaign.H)
	// fmt.Println("H1 point:", H1)

	return &campaign, nil
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
