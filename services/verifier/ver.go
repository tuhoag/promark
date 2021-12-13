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

	"github.com/gorilla/mux"

	"github.com/bwesterb/go-ristretto"
	redis "gopkg.in/redis.v4"
)

type VerifierCryptoParams struct {
	CamId string `json:"camId"`
	H     string `json:"h"`
	S     string `json:"s"`
}

type NewVerifierCryptoParamsRequest struct {
	CamId string `json:"camId"`
	H     string `json:"h"`
}

type VerificationRequest struct {
	CamId string `json:"camId"`
	R     string `json:"R"`
}

type VerificationResponse struct {
	CamId string `json:"camId"`
	H     string `json:"h"`
	R     string `json:"r"`
	S     string `json:"s"`
	Comm  string `json:"comm"`
}

type CampaignCustomerVerifierProof struct {
	CamId  string `json:"camId"`
	UserId string `json:"userId`
	H      string `json:"h"`
	R      string `json:"r"`
	S      string `json:"s"`
	Comm   string `json:"comm"`
}

type VerifierProofVerificationRequest struct {
	CamId string `json:"camId"`
	R     string `json:"r"`
}

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

type HNumber struct {
	ID     string `json:"id"`
	HValue string `json:"hvalue"`
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

var verifierID string

func main() {
	var err error
	// var verifierID string
	var port string
	verifierID = os.Getenv("CORE_PEER_ID")
	port = os.Getenv("VER_PORT")
	logName := verifierID + ".log"

	fmt.Println("Verifier ID: " + verifierID)
	fmt.Println("LogName: " + logName)
	f, err = os.Create(logName)

	if err != nil {
		panic(err)
	}

	// // initilization
	// id = "00"
	// sValue = getSecretNumber(id)

	// // n, _ := f.WriteString("Secret value is:" + string(sValue.S) + string("\n"))
	// // fmt.Println(n)
	fmt.Printf("Starting to listen on port: %s\n", port)
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homeHandler)
	myRouter.HandleFunc("/camp/{id}", GetCampaignCryptoParamsHandler)
	myRouter.HandleFunc("/camp", CreateVerifierCampaignCryptoParamsHandler).Methods("POST")
	myRouter.HandleFunc("/camp/{id}/proof/{userId}", CreateCustomerCampaignProofHandler).Methods("POST")
	myRouter.HandleFunc("/camp/{id}/proof/{userId}", GetCustomerCampaignProofHandler)
	myRouter.HandleFunc("/camp/{id}/verify", VerifyCustomerCampaignProofHandler).Methods("POST")
	myRouter.HandleFunc("/comm", computeComm)
	myRouter.HandleFunc("/verify", verifyComm)
	log.Fatal(http.ListenAndServe(":"+port, myRouter))

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Getting HOME")
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	verifierId := os.Getenv("CORE_PEER_ID")
	port := os.Getenv("VER_PORT")

	verifierURL := verifierId + ":" + port

	// n, _ := f.WriteString("Welcome to verifying service" + string("\n"))
	// fmt.Println(n)

	fmt.Fprintf(w, "VerifierURL: "+verifierURL)
}
func GetCustomerCampaignProofHandler(w http.ResponseWriter, r *http.Request) {
	f.WriteString("getCampaignCryptoParamsHandler() calling" + string("\n"))
	vars := mux.Vars(r)
	camId := vars["id"]
	userId := vars["userId"]

	vCryptoParams, err := GetVerifierCryptoParams(camId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subProof, err := GetCustomerCampaignProof(camId, userId, *vCryptoParams)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(subProof)
}

func VerifyCustomerCampaignProofHandler(w http.ResponseWriter, r *http.Request) {
	f.WriteString("VerifyCustomerCampaignProofHandler() calling\n")
	vars := mux.Vars(r)
	camId := vars["id"]

	f.WriteString("camId:" + camId + "\n")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err)
	}
	f.WriteString("body:" + string(body) + "\n")

	var request VerificationRequest

	// var cryptoParams CampaignCryptoParams
	err = json.Unmarshal(body, &request)
	if err != nil {
		println(err)
	}

	fmt.Println("request.R:" + request.R)
	f.WriteString("request.R:" + request.R + "\n")
	SendLog(verifierID+":request.R", request.R)

	// get crypto campaigns
	vCryptoParams, err := GetVerifierCryptoParams(camId)

	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		f.WriteString("ERROR: " + err.Error() + "\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if it is unexisted
	// return error
	if vCryptoParams == nil {
		http.Error(w, "vCryptoParams is not existed", http.StatusInternalServerError)
		return
	}

	fmt.Println("vCryptoParams.H:" + vCryptoParams.H)
	f.WriteString("vCryptoParams.H:" + vCryptoParams.H + "\n")
	fmt.Println("vCryptoParams.S:" + vCryptoParams.S)
	f.WriteString("vCryptoParams.S:" + vCryptoParams.S + "\n")
	SendLog(verifierID+":vCryptoParams.S", vCryptoParams.H)
	SendLog(verifierID+":vCryptoParams.S", vCryptoParams.S)
	// convert H
	hDec, _ := b64.StdEncoding.DecodeString(vCryptoParams.H)
	hPoint := convertBytesToPoint(hDec)

	// convert s
	sDec, _ := b64.StdEncoding.DecodeString(vCryptoParams.S)
	sScalar := convertBytesToScalar(sDec)

	// convert r
	rDec, _ := b64.StdEncoding.DecodeString(request.R)
	rScalar := convertBytesToScalar(rDec)

	// calculate commitment
	comm := commitTo(&hPoint, &rScalar, &sScalar)
	commEnc := b64.StdEncoding.EncodeToString(comm.Bytes())

	response := VerificationResponse{
		CamId: camId,
		S:     vCryptoParams.S,
		R:     request.R,
		Comm:  commEnc,
		H:     vCryptoParams.H,
	}

	fmt.Println("response.Comm:" + response.Comm)
	f.WriteString("response.Comm:" + response.Comm + "\n")

	json.NewEncoder(w).Encode(response)
}

func CreateVerifierCampaignCryptoParamsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err)
	}
	log.Println(string(body))

	var paramsRequest NewVerifierCryptoParamsRequest
	err = json.Unmarshal(body, &paramsRequest)

	_, err = SetVerifierCryptoParams(paramsRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	vCryptoParams, err := GetVerifierCryptoParams(paramsRequest.CamId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(vCryptoParams)
}

func GetCampaignCryptoParamsHandler(w http.ResponseWriter, r *http.Request) {
	// input: H value, camId
	n, err := f.WriteString("GetCampaignCryptoParamsHandler() calling" + string("\n"))
	fmt.Println(n)

	vars := mux.Vars(r)
	camId := vars["id"]

	// get from redis
	vCryptoParams, err := GetVerifierCryptoParams(camId)

	if err != nil {
		fmt.Errorf("Cannot get cryptoparams %s", err)
	}

	json.NewEncoder(w).Encode(vCryptoParams)
}

func GetVerifierCryptoParams(camId string) (*VerifierCryptoParams, error) {
	client := GetRedisConnection()

	var cryptoParams VerifierCryptoParams

	val, err := client.Get(camId).Result()
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		f.WriteString("ERROR: " + err.Error())
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &cryptoParams)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		f.WriteString("ERROR: " + err.Error())
		return nil, err
	}

	return &cryptoParams, nil
}

func SetVerifierCryptoParams(paramsRequest NewVerifierCryptoParamsRequest) (bool, error) {
	client := GetRedisConnection()

	var cryptoParams VerifierCryptoParams

	val, err := client.Get(paramsRequest.CamId).Result()
	err = json.Unmarshal([]byte(val), &cryptoParams)
	var s ristretto.Scalar
	if err != nil {
		// params are not existed
		fmt.Println(err)
		s.Rand()
		sBytes := s.Bytes()
		sEnc := b64.StdEncoding.EncodeToString(sBytes)

		cryptoParams = VerifierCryptoParams{
			CamId: paramsRequest.CamId,
			H:     paramsRequest.H,
			S:     sEnc,
		}

		jsonParam, err := json.Marshal(cryptoParams)
		if err != nil {
			return false, err
		}

		err = client.Set(cryptoParams.CamId, jsonParam, 0).Err()
		if err != nil {
			return false, err
		}

	} else {
		fmt.Printf("The VerifierCryptoParams is existed for id %s", cryptoParams.CamId)
		return false, nil
	}

	return true, nil
}

func CreateCustomerCampaignProofHandler(w http.ResponseWriter, r *http.Request) {
	// input: camId
	f.WriteString("CreateCustomerCampaignProofHandler() calling" + string("\n"))
	vars := mux.Vars(r)
	camId := vars["id"]
	userId := vars["userId"]

	// get cryptoParams in db
	vCryptoParams, err := GetVerifierCryptoParams(camId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if it is unexisted
	// return error
	if vCryptoParams == nil {
		http.Error(w, "vCryptoParams is not existed", http.StatusInternalServerError)
		return
	}

	// otherwise:
	subProof, err := SetCustomerCampaignProof(camId, userId, *vCryptoParams)

	if err != nil && subProof == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(subProof)
}

func GetCustomerCampaignProof(camId string, userId string, vCryptoParams VerifierCryptoParams) (*CampaignCustomerVerifierProof, error) {
	var subProof CampaignCustomerVerifierProof
	proofId := camId + ":" + userId
	client := GetRedisConnection()
	val, err := client.Get(proofId).Result()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &subProof)

	if err != nil {
		return nil, err
	}

	return &subProof, nil
}

func SetCustomerCampaignProof(camId string, userId string, vCryptoParams VerifierCryptoParams) (*CampaignCustomerVerifierProof, error) {
	var subProof CampaignCustomerVerifierProof
	proofId := camId + ":" + userId
	client := GetRedisConnection()
	val, err := client.Get(proofId).Result()
	err = json.Unmarshal([]byte(val), &subProof)

	if err != nil {
		// params are not existed
		fmt.Println(err)
		// generate R
		var rScalar ristretto.Scalar
		rScalar.Rand()

		// convert H
		hDec, _ := b64.StdEncoding.DecodeString(vCryptoParams.H)
		hPoint := convertBytesToPoint(hDec)

		// convert s
		sDec, _ := b64.StdEncoding.DecodeString(vCryptoParams.S)
		sScalar := convertBytesToScalar(sDec)

		// calculate commitment
		comm := commitTo(&hPoint, &rScalar, &sScalar)

		// return R, Com
		rEnc := b64.StdEncoding.EncodeToString(rScalar.Bytes())
		commEnc := b64.StdEncoding.EncodeToString(comm.Bytes())
		subProof = CampaignCustomerVerifierProof{
			CamId:  camId,
			UserId: userId,
			H:      vCryptoParams.H,
			R:      rEnc,
			S:      vCryptoParams.S,
			Comm:   commEnc,
		}

		jsonParam, err := json.Marshal(subProof)
		if err != nil {
			return nil, err
		}

		err = client.Set(proofId, jsonParam, 0).Err()
		if err != nil {
			return nil, err
		}

	} else {
		fmt.Printf("The CampaignCustomerVerifierProof is existed for id %s", proofId)
		return &subProof, fmt.Errorf("subProof with id %s is existed", proofId)
	}

	return &subProof, nil
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

	// store hValue to redid db
	stored := storeHValue(commParam.ID, commParam.Hbyte)
	fmt.Println(stored)

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

func verifyComm(rw http.ResponseWriter, req *http.Request) {
	var r ristretto.Scalar
	var V ristretto.Scalar
	var comm ristretto.Point
	var hValue string

	n, err := f.WriteString("computeComm() calling" + string("\n"))
	fmt.Println(n)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		println(err)
	}
	log.Println(string(body))

	var commParam CommParam
	err = json.Unmarshal(body, &commParam)

	// get Hvalue from database
	hValue = getHValue(commParam.ID)
	hDec, _ := b64.StdEncoding.DecodeString(hValue)
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

func storeHValue(id string, h string) bool {
	var returnValue bool
	var hValue HNumber

	// connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	val, err := client.Get(id).Result()
	err = json.Unmarshal([]byte(val), &hValue)

	if err != nil {
		fmt.Println(err)
		hValue = HNumber{ID: id, HValue: h}

		jsonParam, err := json.Marshal(hValue)
		if err != nil {
			fmt.Println(err)
		}

		err = client.Set(id, jsonParam, 0).Err()
		if err != nil {
			fmt.Println(err)
			returnValue = false
		}

		returnValue = true
	} else {
		fmt.Print("The Hvalue is existed for id")
		returnValue = true
	}

	return returnValue
}

func getHValue(id string) string {
	var returnHValue string
	var H HNumber
	// connect to Redis

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	val, err := client.Get(id).Result()
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal([]byte(val), &H)
	if err != nil {
		fmt.Println(err)
	}

	returnHValue = H.HValue

	return returnHValue
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
func commitTo(H *ristretto.Point, r *ristretto.Scalar, x *ristretto.Scalar) ristretto.Point {
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

type DebugLog struct {
	Name  string
	Value string
}

var logURL = "http://logs.promark.com:5003/log"

func SendLog(name, message string) {
	logmessage := DebugLog{name, message}
	jsonLog, err := json.Marshal(logmessage)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
	}

	logRequest := string(jsonLog)

	c := &http.Client{}

	reqJSON, err := http.NewRequest("POST", logURL, strings.NewReader(logRequest))
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
