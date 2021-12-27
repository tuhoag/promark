module github.com/tuhoag/promark/services/verifier

go 1.13

require (
	github.com/bwesterb/go-ristretto v1.2.0
	github.com/go-redis/redis/v8 v8.11.4
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.17.0 // indirect
	internal/promark_utils v0.0.0
)

replace internal/promark_utils => ../internal/promark_utils
