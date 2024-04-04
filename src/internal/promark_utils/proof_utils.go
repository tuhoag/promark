package promark_utils

import (
	"bufio"
	// b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	// "math/big"
	"net"

	ristretto "github.com/bwesterb/go-ristretto"
	eutils "github.com/tuhoag/elliptic-curve-cryptography-go/utils"
	// putils "internal/promark_utils"
)

func VerifyPoCSocket(campaign *Campaign, proof *PoCProof) (bool, error) {
	// fmt.Printf("proof.Rs: %s\n", proof.Rs)

	var C ristretto.Point
	C.SetZero()

	r, err := eutils.ConvertStringToScalar(proof.R)
	if err != nil {
		return false, err
	}

	rs := eutils.SplitScalar(r, proof.NumVerifiers)
	for i, verifierURL := range campaign.VerifierURLs {
		// call verifier to compute sub commitment
		r1Enc := eutils.ConvertScalarToString(rs[i])

		CiStr, err := RequestVerification(campaign.Id, r1Enc, verifierURL)
		if err != nil {
			return false, err
		}

		Ci, err := eutils.ConvertStringToPoint(CiStr)
		if err != nil {
			return false, err
		}

		C.Add(&C, Ci)
	}

	comm, err := eutils.ConvertStringToPoint(proof.Comm)
	if err != nil {
		return false, err
	}

	// putils.SendLog("proof.Com", proof.Comm, LOG_MODE)
	// putils.SendLog("calculated Com", b64.StdEncoding.EncodeToString(C.Bytes()), LOG_MODE)
	if C.Equals(comm) {
		return true, nil
	} else {
		return false, nil
	}
}

func RequestVerification(camId string, r string, url string) (string, error) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		fmt.Println("Error connecting:" + err.Error())
		return "", errors.New("ERROR:" + err.Error())
	}

	requestArgs := VerificationRequest{
		CamId: camId,
		R:     r,
	}

	jsonArgs, err := json.Marshal(requestArgs)

	SendRequest(conn, "commit-nocam", string(jsonArgs))
	// wait for response
	// wait for response
	responseStr, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		// sendLog("Error connecting:", err.Error())
		log.Println("Error after creating:", err.Error())
		return "", errors.New("Error  after creating:" + err.Error())
	}
	fmt.Println("Reiceived From: " + url + "-Response:" + responseStr)

	response, err := ParseResponse(responseStr)

	if err != nil {
		return "", errors.New("Error:" + err.Error())
	}

	var verificationResponse VerificationResponse
	err = json.Unmarshal([]byte(response.Data), &verificationResponse)

	if err != nil {
		fmt.Println("error: " + err.Error())
		return "", err
	}

	// putils.SendLog("verificationResponse.H:", verificationResponse.H, LOG_MODE)
	// putils.SendLog("verificationResponse.s:", verificationResponse.S, LOG_MODE)
	// putils.SendLog("verificationResponse.r:", verificationResponse.R, LOG_MODE)
	// putils.SendLog("verificationResponse.Comm:", verificationResponse.Comm, LOG_MODE)

	return verificationResponse.Comm, nil
}

func VerifyTPoCSocket(campaign *Campaign, proof *TPoCProof) (bool, error) {
	fmt.Printf("proof.TRs: %s\n", proof.TRs)

	var C, C2 ristretto.Point
	C.SetZero()
	C2.SetZero()

	for i, verifierURL := range campaign.VerifierURLs {
		// call verifier to compute sub commitment
		CiStr, err := RequestVerification(campaign.Id, proof.TRs[i], verifierURL)
		if err != nil {
			return false, err
		}

		Ci, err := eutils.ConvertStringToPoint(CiStr)
		if err != nil {
			return false, err
		}

		C2i, err := eutils.ConvertStringToPoint(proof.TComms[i])
		if err != nil {
			return false, err
		}

		C.Add(&C, Ci)
		C2.Add(&C2, C2i)
	}

	if C.Equals(&C2) {
		return true, nil
	} else {
		return false, nil
	}
}
