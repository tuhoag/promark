package main

import (
	// "bufio"
	"encoding/json"
	"strconv"
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
	// eutils "github.com/tuhoag/elliptic-curve-cryptography-go/utils"
)

var LOG_MODE = "test"
var DEVICE_VALIDATION_MODE = "device"
var ALL_VALIDATION_MODE = "all"

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

	proof, err := putils.GeneratePoCProofFromVerifiers(&campaign)

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

	result, err := putils.GenerateTPoCs(poc, numTPoCs)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PoCSmartContract) GenerateTPoCProofs(ctx contractapi.TransactionContextInterface, camId string, CStr string, rStr string, numVerifiers int, numTPoCs int) (*putils.PoCAndTPoCProofs, error) {
	putils.SendLog("GeneratePoCProof", "", LOG_MODE)
	putils.SendLog("camId:", camId, LOG_MODE)
	fmt.Printf("camId:%s - C:%s - r:%s - numV:%d - numT:%d\n", camId, CStr, rStr, numVerifiers, numTPoCs)

	poc := putils.PoCProof{
		Comm:         CStr,
		R:            rStr,
		NumVerifiers: numVerifiers,
	}

	result, err := putils.GenerateTPoCs(&poc, numTPoCs)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func CallGetCampaignById(ctx contractapi.TransactionContextInterface, camId string) (*putils.Campaign, error) {
	campaignChaincodeArgs := util.ToChaincodeArgs("GetCampaignById", camId)
	response := ctx.GetStub().InvokeChaincode("campaign", campaignChaincodeArgs, "mychannel")

	putils.SendLog("response.Payload", string(response.Payload), LOG_MODE)
	putils.SendLog("response.Status", string(response.Status), LOG_MODE)
	putils.SendLog("response.message", string(response.Message), LOG_MODE)

	if response.Message != "" {
		return nil, fmt.Errorf(response.Message)
	}

	var campaign putils.Campaign

	err := json.Unmarshal([]byte(response.Payload), &campaign)
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}

func CallFindTokenTransactions(ctx contractapi.TransactionContextInterface, startTimeStr string, endTimeStr string) ([]putils.CustomerCampaignTokenTransaction, error) {
	chaincodeArgs := util.ToChaincodeArgs("FindTokenTransactionsByTimestamps", startTimeStr, endTimeStr)
	response := ctx.GetStub().InvokeChaincode("proof", chaincodeArgs, "mychannel")

	putils.SendLog("response.Payload", string(response.Payload), LOG_MODE)
	putils.SendLog("response.Status", string(response.Status), LOG_MODE)
	putils.SendLog("response.message", string(response.Message), LOG_MODE)

	if response.Message != "" {
		return nil, fmt.Errorf(response.Message)
	}

	var trans []putils.CustomerCampaignTokenTransaction

	err := json.Unmarshal([]byte(response.Payload), &trans)
	if err != nil {
		return nil, err
	}

	return trans, nil
}

func (s *PoCSmartContract) SimulateFindTokenTransactionsByCampaignId(ctx contractapi.TransactionContextInterface, camId string, mode string, limit int) ([]int, error) {
	var counts = []int{0, 0}

	if mode != DEVICE_VALIDATION_MODE && mode != ALL_VALIDATION_MODE {
		return counts, fmt.Errorf("mode '%s' is unsupported", mode)
	}

	campaign, err := CallGetCampaignById(ctx, camId)

	if err != nil {
		return counts, err
	}

	var rawtrans []putils.CustomerCampaignTokenTransaction
	rawtrans, err = CallFindTokenTransactions(ctx, strconv.FormatInt(campaign.StartTime, 10), strconv.FormatInt(campaign.EndTime, 10))

	// queryString := fmt.Sprintf(`{"selector":{"addedTime":{"$gte": %v,"$lte": %v}}}`, campaign.StartTime, campaign.EndTime)
	// resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return counts, err
	}
	// defer resultsIterator.Close()

	// var trans []*putils.CustomerCampaignTokenTransaction
	var count = 0
	var validCount = 0
	var numTrans = len(rawtrans)

	for count < limit {
		var tokenTransaction putils.CustomerCampaignTokenTransaction

		var tokenIdx int = count % numTrans
		tokenTransaction = rawtrans[tokenIdx]

		// verify
		var tpocs []putils.TPoCProof

		if mode == DEVICE_VALIDATION_MODE {
			tpocs = append(tpocs, tokenTransaction.DeviceTPoC)
		} else if mode == ALL_VALIDATION_MODE {
			tpocs = append(tpocs, tokenTransaction.DeviceTPoC)
			tpocs = append(tpocs, tokenTransaction.CustomerTPoC)
		}

		validity := false
		for _, tpoc := range tpocs {
			validity, err = putils.VerifyTPoCSocket(campaign, &tpoc)

			if err != nil {
				continue
			}
		}

		if validity {
			validCount += 1
		}

		count += 1
	}

	counts[0] = count
	counts[1] = validCount

	return counts, nil
}
