package campaign

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	// "net/url"

	"github.com/bwesterb/go-ristretto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Campaign struct {
	ID         string `json:"ID"`
	Name       string `json:"Name"`
	Advertiser string `json:"Advertiser"`
	Business   string `json:"Business"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	campaigns := []Campaign{
		{ID: "id1", Name: "campaign1", Advertiser: "Adv0", Business: "Bus0"},
		{ID: "id2", Name: "campaign2", Advertiser: "Adv0", Business: "Bus0"},
	}

	for _, campaign := range campaigns {
		campaignJSON, err := json.Marshal(campaign)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(campaign.ID, campaignJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Campaign, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	// close the resultsIterator when this function is finished
	defer resultsIterator.Close()

	var campaigns []*Campaign
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var campaign Campaign
		err = json.Unmarshal(queryResponse.Value, &campaign)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, &campaign)
	}

	return campaigns, nil
}

// Create a new campaign
func (s *SmartContract) CreateCampaign(ctx contractapi.TransactionContextInterface, id string, name string, advertiser string, business string) error {
	existing, err := ctx.GetStub().GetState(id)

	if err != nil {
		return errors.New("Unable to read the world state")
	}

	if existing != nil {
		return fmt.Errorf("Cannot create asset since its id %s is existed", id)
	}

	// try to call to web service: start
	response, err := http.Get("http://external.promark.com:5000/")

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	//end

	//message := string(responseData)
	//var H ristretto.Point

	message := convertStringToPoint(string(responseData))
	fmt.Println("test value of message", message)

	campaign := Campaign{
		ID:         id,
		Name:       name,
		Advertiser: advertiser,
		Business:   business,
	}

	campaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, campaignJSON)

	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) QueryLedgerById(ctx contractapi.TransactionContextInterface, id string) ([]*Campaign, error) {
	queryString := fmt.Sprintf(`{"selector":{"id":{"$lte": "%s"}}}`, id)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	var campaigns []*Campaign

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var campaign Campaign
		err = json.Unmarshal(queryResponse.Value, &campaign)
		if err != nil {
			return nil, err
		}

		campaigns = append(campaigns, &campaign)
	}

	resultsIterator.Close()

	response, err := http.Get("http://external.promark.com:5000/")
	if err != nil {
		fmt.Println(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	message := convertStringToPoint(string(responseData))
	fmt.Println("test value of message", message)

	return campaigns, nil
}

/////////////////// Pedersen functions //////////////////////////////////
func (s *SmartContract) testPedersen(ctx contractapi.TransactionContextInterface) ([]*Campaign, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}

	// close the resultsIterator when this function is finished
	defer resultsIterator.Close()

	var campaigns []*Campaign
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var campaign Campaign
		err = json.Unmarshal(queryResponse.Value, &campaign)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, &campaign)
	}

	// response, err := http.Get("http://external.promark.com:5000/")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// responseData, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// message := convertStringToPoint(string(responseData))
	// fmt.Println("test value of message", message)
	//
	return campaigns, nil
}

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

func convertStringToPoint(s string) ristretto.Point {

	var H ristretto.Point
	var hBytes [32]byte

	tmp := []byte(s)
	copy(hBytes[:32], tmp[:])

	H.SetBytes(&hBytes)
	fmt.Println("in convertPointtoBytes hString: ", H)

	return H
}

func convertBytesToPoint(s string) ristretto.Point {

	var H ristretto.Point
	var hBytes [32]byte

	//tmp := []byte(s)
	copy(hBytes[:32], s[:])

	H.SetBytes(&hBytes)
	fmt.Println("in convertPointtoBytes hString: ", H)

	return H
}
