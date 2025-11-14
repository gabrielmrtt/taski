package core

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
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
	CreatedAt *DateTime `json:"createdAt"`
	UpdatedAt *DateTime `json:"updatedAt"`
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
		field := InvalidInputErrorField{
			Field: "name",
			Error: "name cannot be empty",
		}
		return NewInvalidInputError("name cannot be empty", []InvalidInputErrorField{field})
	}

	if len(n.Value) < 3 || len(n.Value) > 255 {
		field := InvalidInputErrorField{
			Field: "name",
			Error: "name must be between 3 and 255 characters",
		}
		return NewInvalidInputError("name must be between 3 and 255 characters", []InvalidInputErrorField{field})
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
		field := InvalidInputErrorField{
			Field: "description",
			Error: "description cannot be empty",
		}
		return NewInvalidInputError("description cannot be empty", []InvalidInputErrorField{field})
	}

	if len(d.Value) > 510 {
		field := InvalidInputErrorField{
			Field: "description",
			Error: "description must be less than 510 characters",
		}
		return NewInvalidInputError("description must be less than 510 characters", []InvalidInputErrorField{field})
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
		field := InvalidInputErrorField{
			Field: "color",
			Error: "color cannot be empty",
		}
		return NewInvalidInputError("color cannot be empty", []InvalidInputErrorField{field})
	}

	if len(c.Value) != 7 || c.Value[0] != '#' {
		field := InvalidInputErrorField{
			Field: "color",
			Error: "color must be a valid hex color",
		}
		return NewInvalidInputError("color must be a valid hex color", []InvalidInputErrorField{field})
	}

	return nil
}

type DateTime struct {
	Value int64
}

/*
NewDateTime creates a new DateTime object with the current epoch time.
*/
func NewDateTime() DateTime {
	return DateTime{Value: datetimeutils.EpochNow()}
}

/*
NewDateTimeFromEpoch creates a new DateTime object from an epoch time.
*/
func NewDateTimeFromEpoch(value int64) (DateTime, error) {
	if value < 0 {
		return DateTime{}, NewInvalidInputError("invalid date timestamp", []InvalidInputErrorField{
			{
				Field: "datetime",
				Error: "datetime cannot be negative",
			},
		})
	}

	return DateTime{Value: value}, nil
}

/*
NewDateTimeFromRFC3339 creates a new DateTime object from an RFC 3339 string.
*/
func NewDateTimeFromRFC3339(value string) (DateTime, error) {
	if !datetimeutils.IsValidRFC3339(value) {
		return DateTime{}, NewInvalidInputError("invalid date timestamp", []InvalidInputErrorField{
			{
				Field: "datetime",
				Error: "datetime is not a valid RFC 3339 date",
			},
		})
	}

	return DateTime{Value: datetimeutils.RFC3339ToEpoch(value)}, nil
}

/*
ToEpoch converts the DateTime object to an epoch time.
*/
func (d DateTime) ToEpoch() int64 {
	return d.Value
}

/*
ToRFC3339 converts the DateTime object to an RFC 3339 string.
*/
func (d DateTime) ToRFC3339() string {
	return datetimeutils.EpochToRFC3339(d.Value)
}

/*
Equals checks if two DateTime objects are equal.
*/
func (d DateTime) Equals(_d DateTime) bool {
	return d.Value == _d.Value
}

/*
IsBefore checks if the DateTime object is before another DateTime object.
*/
func (d DateTime) IsBefore(_d DateTime) bool {
	return d.Value < _d.Value
}

/*
IsAfter checks if the DateTime object is after another DateTime object.
*/
func (d DateTime) IsAfter(_d DateTime) bool {
	return d.Value > _d.Value
}

/*
IsBeforeOrEqual checks if the DateTime object is before or equal to another DateTime object.
*/
func (d DateTime) IsBeforeOrEqual(_d DateTime) bool {
	return d.Value <= _d.Value
}

/*
IsAfterOrEqual checks if the DateTime object is after or equal to another DateTime object.
*/
func (d DateTime) IsAfterOrEqual(_d DateTime) bool {
	return d.Value >= _d.Value
}
