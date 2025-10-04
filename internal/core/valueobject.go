package core

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/gabrielmrtt/taski/pkg/encodingutils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Identity struct {
	Public   string
	Internal uuid.UUID
}

func NewIdentity(publicPrefix string) Identity {
	uuid := uuid.New()
	length := new(big.Int).SetBytes(uuid[:])

	return Identity{
		Public:   publicPrefix + "_" + encodingutils.GenerateEncodedBase62String(length),
		Internal: uuid,
	}
}

func NewIdentityWithoutPublic() Identity {
	return Identity{
		Public:   "",
		Internal: uuid.New(),
	}
}

func NewIdentityWithoutPublicFromInternal(internalId uuid.UUID) Identity {
	return Identity{
		Public:   "",
		Internal: internalId,
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
		fmt.Printf("INVALID IDENTITY FOR %s", publicId)
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

func (i Identity) Equals(_i Identity) bool {
	return i.Public == _i.Public && i.Internal == _i.Internal
}

func (i Identity) IsEmpty() bool {
	return i.Internal == uuid.Nil
}

type Timestamps struct {
	CreatedAt *int64 `json:"createdAt"`
	UpdatedAt *int64 `json:"updatedAt"`
}

type Name struct {
	Value string
}

func NewName(value string) (Name, error) {
	n := Name{Value: value}

	if err := n.Validate(); err != nil {
		return Name{}, err
	}

	return n, nil
}

func (n Name) Validate() error {
	if n.Value == "" {
		return errors.New("name cannot be empty")
	}

	if len(n.Value) < 3 || len(n.Value) > 255 {
		return errors.New("name must be between 3 and 255 characters")
	}

	return nil
}

func (n Name) String() string {
	return n.Value
}

func (n Name) Equals(_n Name) bool {
	return n.Value == _n.Value
}

func IsValidName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	_, err := NewName(name)

	return err == nil
}

type Description struct {
	Value string
}

func NewDescription(value string) (Description, error) {
	d := Description{Value: value}

	if err := d.Validate(); err != nil {
		return Description{}, err
	}

	return d, nil
}

func (d Description) Validate() error {
	if d.Value == "" {
		return errors.New("description cannot be empty")
	}

	if len(d.Value) > 510 {
		return errors.New("description must be less than 510 characters")
	}

	return nil
}

func (d Description) String() string {
	return d.Value
}

func (d Description) Equals(_d Description) bool {
	return d.Value == _d.Value
}

func IsValidDescription(fl validator.FieldLevel) bool {
	description := fl.Field().String()

	_, err := NewDescription(description)

	return err == nil
}

type Color struct {
	Value string
}

func NewColor(value string) (Color, error) {
	c := Color{Value: value}

	if err := c.Validate(); err != nil {
		return Color{}, err
	}

	return c, nil
}

func (c Color) Validate() error {
	if c.Value == "" {
		return errors.New("color cannot be empty")
	}

	if len(c.Value) != 7 || c.Value[0] != '#' {
		return errors.New("color must be a valid hex color")
	}

	return nil
}
