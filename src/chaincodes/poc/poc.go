package main

import (
	// "bufio"
	"encoding/json"
	// "errors"
	// "log"
	// "net"
	// "sync"

	// "errors"
	"fmt"

	// "log"

	// "strconv"
	// "strings"

	putils "internal/promark_utils"

	// "github.com/bwesterb/go-ristretto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/util"
	eutils "github.com/tuhoag/elliptic-curve-cryptography-go/utils"
)

var LOG_MODE = "test"

type PoCSmartContract struct {
	contractapi.Contract
}

// func (s *PoCSmartContract) InitDefaultCryptoParams(ctx contractapi.TransactionContextInterface) (*putils.PoCProof, error) {

// }

func (s *PoCSmartContract) GeneratePoCProof(ctx contractapi.TransactionContextInterface, camId string, userId string) (*putils.PoCProof, error) {
	putils.SendLog("GeneratePoCProof", "", LOG_MODE)
	putils.SendLog("camId:", camId, LOG_MODE)
	putils.SendLog("userId", userId, LOG_MODE)

	campaignChaincodeArgs := util.ToChaincodeArgs("GetCampaignById", camId)
	response := ctx.GetStub().InvokeChaincode("campaign", campaignChaincodeArgs, "mychannel")

	putils.SendLog("response.Payload", string(response.Payload), LOG_MODE)
	putils.SendLog("response.Status", string(response.Status), LOG_MODE)
	putils.SendLog("response.message", string(response.Message), LOG_MODE)
	// putils.SendLog("response.error", string(response.Error))
	// putils.SendLog("response.message is nil", strconv.FormatBool(response.Message == ""))

	if response.Message != "" {
		return nil, fmt.Errorf(response.Message)
	}

	var campaign putils.Campaign

	err := json.Unmarshal([]byte(response.Payload), &campaign)
	if err != nil {
		return nil, err
	}

	proof, err := putils.GenerateProofFromVerifiers(&campaign, userId)

	if err != nil {
		return nil, err
	}

	return proof, nil
}

func (s *PoCSmartContract) GeneratePoCAndTPoCProof(ctx contractapi.TransactionContextInterface, camId string, userId string, numTPoCs int) (*putils.PoCAndTPoCProofs, error) {
	putils.SendLog("GeneratePoCProof", "", LOG_MODE)
	putils.SendLog("camId:", camId, LOG_MODE)
	putils.SendLog("userId", userId, LOG_MODE)

	poc, err := s.GeneratePoCProof(ctx, camId, userId)
	if err != nil {
		return nil, err
	}

	numVerifiers := len(poc.Rs)
	var tpocs = make([]putils.TPoCProof, numTPoCs)

	for i := 0; i < numTPoCs; i++ {
		commPoint, err := eutils.ConvertStringToPoint(poc.Comm)
		if err != nil {
			return nil, err
		}

		subComms := eutils.Split(commPoint, numVerifiers)

		var tcomms = make([]string, numVerifiers)
		var hashes = make([]string, numVerifiers)
		var key string

		for j := 0; j < numVerifiers; j++ {
			tcomms[j] = eutils.ConvertPointToString(subComms[j])
		}

		tpoc := putils.TPoCProof{
			TComms: tcomms,
			TRs:    poc.Rs,
			Hashes: hashes,
			Key:    key,
		}
		tpocs[i] = tpoc
	}

	result := putils.PoCAndTPoCProofs{
		PoC:   *poc,
		TPoCs: tpocs,
	}
	return &result, nil
}
