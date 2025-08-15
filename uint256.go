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

// NewUint256 returns a new Uint256.
func NewUint256(x *big.Int) (Uint256, error) {
	var x256 Uint256
	{
		if err := x256.setBigInt(x); err != nil {
			return Uint256{}, err
		}
	}

	return x256, nil
}

// MustNewUint256 returns a new Uint256.
// It panics for invalid input.
func MustNewUint256(x *big.Int) Uint256 {
	x256, err := NewUint256(x)
	if err != nil {
		panic(err)
	}

	return x256
}

// NewUint256FromUint64 returns a new Uint256 from a uint64.
func NewUint256FromUint64(i uint64) Uint256 {
	return MustNewUint256(new(big.Int).SetUint64(i))
}

// NewUint256FromHex returns a new Uint256 from a hexadecimal string.
func NewUint256FromHex(s string) (Uint256, error) {
	var x256 Uint256
	{
		if err := x256.setHex(s); err != nil {
			return Uint256{}, err
		}
	}

	return x256, nil
}

// MustNewUint256FromHex returns a new Uint256 from a hexadecimal string.
// It panics for invalid input.
func MustNewUint256FromHex(s string) Uint256 {
	x256, err := NewUint256FromHex(s)
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
		return oops.New("src bytes must not be empty")
	}
	if len(b) > maxByteLength {
		return oops.Errorf("src bytes must be less than or equal to %d bytes", maxByteLength)
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
	if has0xPrefix(string(text)) {
		return x256.setHex(string(text))

	} else {
		x := new(big.Int)
		{
			if err := x.UnmarshalText(text); err != nil {
				return oops.Wrap(err)
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
		return oops.New("x must not be negative")
	}
	if x.BitLen() > maxBitLength {
		return oops.Errorf("x must be less than or equal to %d bits", maxBitLength)
	}

	x256.x = *x

	return nil
}

func (x256 *Uint256) setHex(s string) error {
	if !has0xPrefix(s) {
		return oops.New("s must have 0x prefix")
	}
	if is0x(s) {
		return oops.New("s must not be 0x")
	}

	var b []byte
	{
		b = append(b, '0', 'x')

		for idx, c := range s[2:] {
			if c == '0' {
				continue
			}

			b = append(b, s[2+idx:]...)

			break
		}

		if is0x(string(b)) {
			b = append(b, '0')
		}
	}

	x, err := ethhexutil.DecodeBig(string(b))
	if err != nil {
		return oops.Wrap(err)
	}

	return x256.setBigInt(x)
}
