module github.com/tuhoag/promark/services/ext

go 1.13

require (
	github.com/bwesterb/go-ristretto v1.2.0
	github.com/gorilla/mux v1.8.0 // indirect
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20160220154919-db14e161995a // indirect
	internal/promark_utils v0.0.0
)

replace internal/promark_utils => ../internal/promark_utils
