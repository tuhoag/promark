module campaign

go 1.13

require (
	github.com/go-openapi/spec v0.20.6 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/gobuffalo/envy v1.10.1 // indirect
	github.com/gobuffalo/packd v1.0.1 // indirect
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20220131132609-1476cf1d3206
	github.com/hyperledger/fabric-contract-api-go v1.1.1
	github.com/hyperledger/fabric-protos-go v0.0.0-20220613214546-bf864f01d75e // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	golang.org/x/net v0.0.0-20220617184016-355a448f1bc9 // indirect
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c // indirect
	google.golang.org/genproto v0.0.0-20220617124728-180714bec0ad // indirect
	internal/promark_utils v0.0.0
)

replace internal/promark_utils => ../../internal/promark_utils

replace github.com/tuhoag/elliptic-curve-cryptography-go => ../../internal/elliptic-curve-cryptography-go
