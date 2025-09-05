package core

import (
	"math/big"
	"strings"

	"github.com/gabrielmrtt/taski/pkg/encodingutils"
	"github.com/google/uuid"
)

type Identity struct {
	Public   string
	Internal uuid.UUID
}

type Timestamps struct {
	CreatedAt *int64
	UpdatedAt *int64
}

func NewIdentity(publicPrefix string) Identity {
	uuid := uuid.New()
	length := new(big.Int).SetBytes(uuid[:])

	return Identity{
		Public:   publicPrefix + "_" + encodingutils.GenerateEncodedBase62String(length),
		Internal: uuid,
	}
}

func NewIdentityFromInternal(internalId uuid.UUID, publicPrefix string) Identity {
	length := new(big.Int).SetBytes(internalId[:])

	return Identity{
		Public:   publicPrefix + "_" + encodingutils.GenerateEncodedBase62String(length),
		Internal: internalId,
	}
}

func NewIdentityFromPublic(publicId string) Identity {
	parts := strings.Split(publicId, "_")
	if len(parts) != 2 {
		return Identity{}
	}

	encodedId := parts[1]

	num, err := encodingutils.DecodeEncodedBase62String(encodedId)
	if err != nil {
		return Identity{}
	}

	bytes := num.Bytes()

	if len(bytes) < 16 {
		padded := make([]byte, 16)
		copy(padded[16-len(bytes):], bytes)
		bytes = padded
	}

	var internalId uuid.UUID
	copy(internalId[:], bytes)

	return Identity{
		Public:   publicId,
		Internal: internalId,
	}
}
