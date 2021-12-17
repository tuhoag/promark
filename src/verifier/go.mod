module github.com/tuhoag/promark/services/verifier

go 1.13

require (
	github.com/bwesterb/go-ristretto v1.2.0
	github.com/garyburd/redigo v1.6.3 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.17.0 // indirect
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20160220154919-db14e161995a // indirect
	gopkg.in/redis.v4 v4.2.4
	internal/promark_utils v0.0.0
)

replace internal/promark_utils => ../internal/promark_utils
