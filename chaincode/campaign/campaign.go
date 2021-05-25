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

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	ID    string `json:"ID"`
	Value int    `json:value`
	Owner string `json:"owner"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "asset1", Value: 5, Owner: "Tomoko"},
		{ID: "asset2", Value: 5, Owner: "Brad"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
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
