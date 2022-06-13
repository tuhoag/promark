module promark_utils

go 1.13

require (
	github.com/bwesterb/go-ristretto v1.2.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/tuhoag/elliptic-curve-cryptography-go v0.0.4
)

replace github.com/tuhoag/elliptic-curve-cryptography-go => ../elliptic-curve-cryptography-go
