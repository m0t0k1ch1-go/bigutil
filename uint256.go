package bigutil

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethhexutil "github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	maxUint256Bits  = 256
	maxUint256Bytes = 32
)

// Uint256 represents an unsigned 256-bit integer.
type Uint256 struct {
	x big.Int
}

// NewUint256 returns a new Uint256.
func NewUint256(x *big.Int) (Uint256, error) {
	var x256 Uint256
	if err := x256.setBigInt(x); err != nil {
		return Uint256{}, err
	}

	return x256, nil
}

// MustNewUint256 panics if the input is invalid.
func MustNewUint256(x *big.Int) Uint256 {
	x256, err := NewUint256(x)
	if err != nil {
		panic(err)
	}

	return x256
}

func (x256 *Uint256) setBigInt(x *big.Int) error {
	if x == nil {
		return errors.New("invalid big.Int: nil")
	}
	if x.Sign() < 0 {
		return errors.New("invalid big.Int: negative")
	}
	if x.BitLen() > maxUint256Bits {
		return fmt.Errorf("invalid big.Int: exceeds %d bits", maxUint256Bits)
	}

	x256.x.Set(x)

	return nil
}

// NewUint256FromHex returns a new Uint256 from a hex string.
// The string must have a 0x/0X prefix; leading zeros are allowed and ignored.
func NewUint256FromHex(s string) (Uint256, error) {
	var x256 Uint256
	if err := x256.setHex(s); err != nil {
		return Uint256{}, err
	}

	return x256, nil
}

// MustNewUint256FromHex panics if the input is invalid.
func MustNewUint256FromHex(s string) Uint256 {
	x256, err := NewUint256FromHex(s)
	if err != nil {
		panic(err)
	}

	return x256
}

func (x256 *Uint256) setHex(s string) error {
	if !strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X") {
		return errors.New("invalid hex string: missing 0x/0X prefix")
	}
	if s == "0x" || s == "0X" {
		return errors.New("invalid hex string: empty")
	}

	d := strings.TrimLeft(s[2:], "0")
	if len(d) == 0 {
		d = "0"
	}

	s = "0x" + d

	x, err := ethhexutil.DecodeBig(s)
	if err != nil {
		return fmt.Errorf("invalid hex string: %w", err)
	}

	return x256.setBigInt(x)
}

// NewUint256FromUint64 returns a new Uint256 from a uint64.
func NewUint256FromUint64(i uint64) Uint256 {
	var x256 Uint256
	x256.x.SetUint64(i)

	return x256
}

// BigInt returns a copy of the underlying big.Int.
func (x256 Uint256) BigInt() *big.Int {
	var x big.Int
	x.Set(&x256.x)

	return &x
}

// String implements fmt.Stringer.
// It returns a 0x-prefixed lowercase hex string with no leading zeros (zero is "0x0").
func (x256 Uint256) String() string {
	return "0x" + x256.x.Text(16)
}

// Value implements driver.Valuer.
// It returns a minimal big-endian []byte (never nil); zero is encoded as a single 0x00 byte.
func (x256 Uint256) Value() (driver.Value, error) {
	b := x256.x.Bytes()
	if len(b) == 0 {
		b = []byte{0x00}
	}

	return b, nil
}

// Scan implements sql.Scanner.
// It accepts a big-endian []byte (length 1-32).
func (x256 *Uint256) Scan(src any) error {
	if src == nil {
		return errors.New("invalid source: nil")
	}

	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("unsupported source type: %T", src)
	}
	if len(b) == 0 {
		return errors.New("invalid source: empty bytes")
	}
	if len(b) > maxUint256Bytes {
		return fmt.Errorf("invalid source: exceeds %d bytes", maxUint256Bytes)
	}

	var x big.Int
	x.SetBytes(b)

	return x256.setBigInt(&x)
}

// MarshalText implements encoding.TextMarshaler.
// It returns a 0x-prefixed lowercase hex string with no leading zeros (zero is "0x0").
func (x256 Uint256) MarshalText() ([]byte, error) {
	return []byte(x256.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It accepts either a 0x/0X-prefixed hex string or a non-negative decimal string.
func (x256 *Uint256) UnmarshalText(text []byte) error {
	s := strings.TrimSpace(string(text))

	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		return x256.setHex(s)
	}

	var x big.Int
	if err := x.UnmarshalText([]byte(s)); err != nil {
		return err
	}

	return x256.setBigInt(&x)
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts a JSON string (0x/0X-prefixed hex or non-negative decimal) or a JSON number (non-negative integer).
func (x256 *Uint256) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.New("invalid json value: empty")
	}
	if string(b) == "null" {
		return errors.New("invalid json value: null")
	}

	if len(b) >= 2 && b[0] == '"' && b[len(b)-1] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return fmt.Errorf("invalid json string: %w", err)
		}

		return x256.UnmarshalText([]byte(s))
	}

	return x256.UnmarshalText(b)
}
