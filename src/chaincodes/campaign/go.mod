module campaign

go 1.13

require (
	github.com/bwesterb/go-ristretto v1.2.0 // indirect
	github.com/hyperledger/fabric-contract-api-go v1.1.1
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20160220154919-db14e161995a // indirect
	internal/promark_utils v0.0.0
)

replace internal/promark_utils => ../../internal/promark_utils
