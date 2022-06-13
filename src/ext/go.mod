module github.com/tuhoag/promark/services/ext

go 1.13

require (
	github.com/bwesterb/go-ristretto v1.2.1
	github.com/tuhoag/elliptic-curve-cryptography-go v0.0.4
	internal/promark_utils v0.0.0
)

replace internal/promark_utils => ../internal/promark_utils

replace github.com/tuhoag/elliptic-curve-cryptography-go => ../internal/elliptic-curve-cryptography-go
