package bigutil

import (
	"database/sql/driver"
	"math/big"

	ethhexutil "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/samber/oops"
)

const (
	maxByteLength = 32
	maxBitLength  = maxByteLength * 8
)

// Uint256 is a wrapper for big.Int that represents uint256.
type Uint256 struct {
	x big.Int
}

// Uint64ToUint256 converts the given uint64 to Uint256.
func Uint64ToUint256(i uint64) Uint256 {
	return MustBigIntToUint256(new(big.Int).SetUint64(i))
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
	x256, err := HexToUint256(s)
	if err != nil {
		panic(err)
	}

	return x256
}

// BigIntToUint256 converts the given big.Int to Uint256.
func BigIntToUint256(x *big.Int) (Uint256, error) {
	x256 := Uint256{}
	{
		if err := x256.setBigInt(x); err != nil {
			return Uint256{}, err
		}
	}

	return x256, nil
}

// MustBigIntToUint256 converts the given big.Int to Uint256.
// It panics for invalid input.
func MustBigIntToUint256(x *big.Int) Uint256 {
	x256, err := BigIntToUint256(x)
	if err != nil {
		panic(err)
	}

	return x256
}

// BigInt returns the big.Int.
func (x256 Uint256) BigInt() *big.Int {
	return &x256.x
}

// String implements the fmt.Stringer interface.
func (x256 Uint256) String() string {
	return x256.string()
}

// Value implements the driver.Valuer interface.
func (x256 Uint256) Value() (driver.Value, error) {
	b := x256.x.Bytes()
	if len(b) == 0 {
		b = []byte{0x0}
	}

	return b, nil
}

// Scan implements the sql.Scanner interface.
func (x256 *Uint256) Scan(src any) error {
	if src == nil {
		return oops.New("src must not be nil")
	}

	b, ok := src.([]byte)
	if !ok {
		return oops.Errorf("unexpected src type: %T", src)
	}
	if len(b) == 0 {
		return oops.New("src must not be empty")
	}
	if len(b) > maxByteLength {
		return oops.Errorf("src must be less than or equal to %d bytes", maxByteLength)
	}

	x256.x.SetBytes(b)

	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (x256 Uint256) MarshalText() ([]byte, error) {
	return []byte(x256.string()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (x256 *Uint256) UnmarshalText(text []byte) error {
	if l := len(text); l >= 2 && text[0] == '0' && text[1] == 'x' {
		if l == 2 {
			return oops.New("must not be empty")
		}

		var textWithoutLeadingZeroDigits []byte
		{
			for idx, c := range text[2:] {
				if c == '0' {
					continue
				}

				textWithoutLeadingZeroDigits = append([]byte{'0', 'x'}, text[2+idx:]...)

				break
			}

			if len(textWithoutLeadingZeroDigits) == 0 {
				textWithoutLeadingZeroDigits = []byte{'0', 'x', '0'}
			}
		}

		x, err := ethhexutil.DecodeBig(string(textWithoutLeadingZeroDigits))
		if err != nil {
			return err
		}

		return x256.setBigInt(x)

	} else {
		x := new(big.Int)
		{
			if err := x.UnmarshalText(text); err != nil {
				return err
			}
		}

		return x256.setBigInt(x)
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (x256 *Uint256) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	return x256.UnmarshalText(b)
}

func (x256 Uint256) string() string {
	return ethhexutil.EncodeBig(&x256.x)
}

func (x256 *Uint256) setBigInt(x *big.Int) error {
	if x.Sign() < 0 {
		return oops.New("must be positive")
	}
	if x.BitLen() > maxBitLength {
		return oops.Errorf("must be less than or equal to %d bits", maxBitLength)
	}

	x256.x = *x

	return nil
}
