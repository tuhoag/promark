module campaign

go 1.13

require (
	github.com/bwesterb/go-ristretto v1.2.0
	github.com/hyperledger/fabric-contract-api-go v1.1.1
	internal/promark_utils v0.0.0
)
replace (
    internal/promark_utils => ../../internal/promark_utils
)
