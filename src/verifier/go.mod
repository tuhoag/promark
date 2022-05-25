module github.com/tuhoag/promark/services/verifier

go 1.13

require (
	github.com/bwesterb/go-ristretto v1.2.1
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-redis/redis/v8 v8.11.5
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/onsi/ginkgo/v2 v2.1.4 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/tuhoag/elliptic-curve-cryptography-go v0.0.4
	golang.org/x/net v0.0.0-20220524220425-1d687d428aca // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	internal/promark_utils v0.0.0
)

replace internal/promark_utils => ../internal/promark_utils
