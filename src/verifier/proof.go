package main

import (
	// "bufio"
	// "context"
	// "path"

	// // "crypto/rsa"
	// b64 "encoding/base64"
	"encoding/json"
	// "errors"

	// // "errors"
	// "fmt"
	// // "log"
	// "io/ioutil"
	// "net"
	// "os"

	// "context"
	"fmt"

	ristretto "github.com/bwesterb/go-ristretto"
	redis "github.com/go-redis/redis/v8"
	pedersen "github.com/tuhoag/elliptic-curve-cryptography-go/pedersen"
	eutils "github.com/tuhoag/elliptic-curve-cryptography-go/utils"
	putils "internal/promark_utils"
	// redis "github.com/go-redis/redis/v8"
	// // "strings"
	// // "log"
	// "path/filepath"
)

func GenerateOrGetPoC(camId string, userId string) (*putils.PoCProof, error) {
	// watch key
	key := fmt.Sprintf("%s:%s", camId, userId)

	// get from the database
	maxRetries := 1000

	var pocProof putils.PoCProof

	txf := func(tx *redis.Tx) error {
		value, err := tx.Get(ctx, key).Result()
		if err != nil {
			return err
		}

		if err == redis.Nil {
			// generate R and store
			var r ristretto.Scalar
			r.Rand()
			C := pedersen.CommitTo(&H, &r, &s)

			rStr := eutils.ConvertScalarToString(&r)
			CStr := eutils.ConvertPointToString(C)

			pocProof.R = rStr
			pocProof.Comm = CStr

			jsonParam, err := json.Marshal(pocProof)
			if err != nil {
				return err
			}

			fmt.Println("Converted subproof to JSON:" + string(jsonParam))

			// Operation is commited only if the watched keys remain unchanged.
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Set(ctx, key, jsonParam, 0)

				fmt.Println("Added subproof to db")
				return nil
			})

			return err
		}

		err = json.Unmarshal([]byte(value), &pocProof)

		return err
	}

	rdb := putils.GetRedisConnection()

	for i := 0; i < maxRetries; i++ {
		err := rdb.Watch(ctx, txf, key)
		if err == nil {
			// success
			break

		}

		if err == redis.TxFailedErr {
			fmt.Println("There are some modifications on subproof")
			continue
		}

		return nil, err
	}

	return &pocProof, nil
}
