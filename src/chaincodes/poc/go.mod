module github.com/tuhoag/promark/chaincodes/poc

go 1.16

require (
	github.com/go-openapi/spec v0.20.6 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/gobuffalo/envy v1.10.1 // indirect
	github.com/gobuffalo/packd v1.0.1 // indirect
	github.com/hyperledger/fabric v2.1.1+incompatible
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20220131132609-1476cf1d3206 // indirect
	github.com/hyperledger/fabric-contract-api-go v1.1.1
	github.com/hyperledger/fabric-protos-go v0.0.0-20220613214546-bf864f01d75e // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/miekg/pkcs11 v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/sykesm/zap-logfmt v0.0.4 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/net v0.0.0-20220630215102-69896b714898 // indirect
	golang.org/x/sys v0.0.0-20220704084225-05e143d24a9e // indirect
	google.golang.org/genproto v0.0.0-20220630174209-ad1d48641aa7 // indirect
	internal/promark_utils v0.0.0
)

replace internal/promark_utils => ../../internal/promark_utils

replace github.com/tuhoag/elliptic-curve-cryptography-go => ../../internal/elliptic-curve-cryptography-go
