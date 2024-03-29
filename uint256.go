package bigutil

import (
	"database/sql/driver"
	"encoding/json"
	"math/big"

	ethhexutil "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

const (
	maxByteLength = 32
	maxBitLength  = maxByteLength * 8
)

// Uint256 is a wrapper for big.Int that represents uint256.
type Uint256 struct {
	x big.Int
}

// HexToUint256 converts the given hex string to Uint256.
func HexToUint256(s string) (Uint256, error) {
	x, err := ethhexutil.DecodeBig(s)
	if err != nil {
		return Uint256{}, err
	}

	return BigIntToUint256(x)
}

// MustHexToUint256 converts the given hex string to Uint256.
// It panics for invalid input.
func MustHexToUint256(s string) Uint256 {
	i, err := HexToUint256(s)
	if err != nil {
		panic(err)
	}

	return i
}

// BigIntToUint256 converts the given big.Int to Uint256.
func BigIntToUint256(x *big.Int) (Uint256, error) {
	i := Uint256{}

	if err := i.setBigInt(x); err != nil {
		return Uint256{}, err
	}

	return i, nil
}

// MustBigIntToUint256 converts the given big.Int to Uint256.
// It panics for invalid input.
func MustBigIntToUint256(x *big.Int) Uint256 {
	i, err := BigIntToUint256(x)
	if err != nil {
		panic(err)
	}

	return i
}

// BigInt returns the big.Int.
func (i Uint256) BigInt() *big.Int {
	return &i.x
}

// String implements the fmt.Stringer interface.
func (i Uint256) String() string {
	return i.string()
}

// Value implements the driver.Valuer interface.
func (i Uint256) Value() (driver.Value, error) {
	b := i.x.Bytes()
	if len(b) == 0 {
		b = []byte{0x0}
	}

	return b, nil
}

// Scan implements the sql.Scanner interface.
func (i *Uint256) Scan(src any) error {
	if src == nil {
		return errors.New("src must not be nil")
	}

	b, ok := src.([]byte)
	if !ok {
		return errors.Errorf("unexpected src type: %T", src)
	}
	if len(b) == 0 {
		return errors.New("src must not be empty")
	}
	if len(b) > maxByteLength {
		return errors.Errorf("src must be less than or equal to %d bytes", maxByteLength)
	}

	i.x.SetBytes(b)

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (i Uint256) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.string())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (i *Uint256) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		if len(b) >= 4 && b[1] == '0' && b[2] == 'x' {
			var x ethhexutil.Big
			if err := json.Unmarshal(b, &x); err != nil {
				return err
			}

			return i.setBigInt(x.ToInt())
		} else {
			var x big.Int
			if err := json.Unmarshal(b[1:len(b)-1], &x); err != nil {
				return err
			}

			return i.setBigInt(&x)
		}
	} else {
		var x big.Int
		if err := json.Unmarshal(b, &x); err != nil {
			return err
		}

		return i.setBigInt(&x)
	}
}

func (i Uint256) string() string {
	return ethhexutil.EncodeBig(&i.x)
}

func (i *Uint256) setBigInt(x *big.Int) error {
	if x.Sign() < 0 {
		return errors.New("must be positive")
	}
	if x.BitLen() > maxBitLength {
		return errors.Errorf("must be less than or equal to %d bits", maxBitLength)
	}

	i.x = *x

	return nil
}
