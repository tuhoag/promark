package campaign

import (
	"encoding/json"
	"errors"
	"fmt"

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

func CreateCampaign(ctx contractapi.TransactionContextInterface, id string, name string, advertiser string, business string) error {
	existing, err := ctx.GetStub().GetState(id)

	if err != nil {
		return errors.New("Unable to read the world state.")
	}

	if existing != nil {
		return fmt.Errorf("Campaign id %s is existed.", id)
	}

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

	fmt.Printf("Campaign JSON: %s", campaignJSON)

	err = ctx.GetStub().PutState(id, campaignJSON)

	if err != nil {
		return err
	}

	return nil
}
