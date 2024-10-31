package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

type Serializable interface {
	Serialize() string
}

func HashSha256[T Serializable](item T) string {
	hash := sha256.Sum256([]byte(item.Serialize()))

	return hex.EncodeToString(hash[:])

}
