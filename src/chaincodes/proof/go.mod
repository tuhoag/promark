module github.com/tuhoag/promark/chaincodes/proof

go 1.13

require (
	github.com/bwesterb/go-ristretto v1.2.0
	github.com/hyperledger/fabric v2.1.1+incompatible
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20200728190242-9b3ae92d8664 // indirect
	github.com/hyperledger/fabric-contract-api-go v1.1.1
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sykesm/zap-logfmt v0.0.4 // indirect
	go.uber.org/zap v1.19.1 // indirect
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20160220154919-db14e161995a // indirect
	internal/promark_utils v0.0.0
)

replace internal/promark_utils => ../../internal/promark_utils
