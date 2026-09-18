package bigutil

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

var (
	_ fmt.Stringer             = Uint256{}
	_ driver.Valuer            = Uint256{}
	_ sql.Scanner              = &Uint256{}
	_ encoding.TextMarshaler   = Uint256{}
	_ json.MarshalerTo         = Uint256{}
	_ json.Marshaler           = Uint256{}
	_ graphql.Marshaler        = Uint256{}
	_ encoding.TextUnmarshaler = &Uint256{}
	_ json.UnmarshalerFrom     = &Uint256{}
	_ json.Unmarshaler         = &Uint256{}
	_ graphql.Unmarshaler      = &Uint256{}
)

// Uint256 represents an unsigned 256-bit integer.
type Uint256 struct {
	x big.Int
}

// NewUint256 returns a new [Uint256] from a [big.Int].
func NewUint256(x *big.Int) (Uint256, error) {
	var x256 Uint256
	if err := x256.setBigInt(x); err != nil {
		return Uint256{}, err
	}

	return x256, nil
}

// MustNewUint256 is like [NewUint256] but panics if the input is invalid.
func MustNewUint256(x *big.Int) Uint256 {
	x256, err := NewUint256(x)
	if err != nil {
		panic(err)
	}

	return x256
}

func (x256 *Uint256) setBigInt(x *big.Int) error {
	if x == nil {
		return errors.New("invalid big integer: nil")
	}
	if x.Sign() < 0 {
		return errors.New("invalid big integer: negative")
	}
	if x.BitLen() > 256 {
		return errors.New("invalid big integer: exceeds 256 bits")
	}

	x256.x.Set(x)

	return nil
}

// NewUint256FromHex returns a new [Uint256] from a hexadecimal integer string.
// The string must have a 0x/0X prefix and must not be signed; leading zeros are allowed and ignored.
func NewUint256FromHex(s string) (Uint256, error) {
	var x256 Uint256
	if err := x256.setHex(s); err != nil {
		return Uint256{}, err
	}

	return x256, nil
}

// MustNewUint256FromHex is like [NewUint256FromHex] but panics if the input is invalid.
func MustNewUint256FromHex(s string) Uint256 {
	x256, err := NewUint256FromHex(s)
	if err != nil {
		panic(err)
	}

	return x256
}

func (x256 *Uint256) setHex(s string) error {
	if len(s) == 0 {
		return errors.New("invalid hexadecimal string: empty")
	}
	if !strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X") {
		return errors.New("invalid hexadecimal string: missing 0x/0X prefix")
	}
	if s == "0x" || s == "0X" {
		return errors.New("invalid hexadecimal string: missing hexadecimal digits after 0x/0X prefix")
	}
	if s[2] == '+' || s[2] == '-' {
		return errors.New("invalid hexadecimal string: must not be signed")
	}

	var x big.Int
	if _, ok := x.SetString(s[2:], 16); !ok {
		return errors.New("invalid hexadecimal string: must contain only hexadecimal digits")
	}

	return x256.setBigInt(&x)
}

func (x256 *Uint256) setDecimal(s string) error {
	if len(s) == 0 {
		return errors.New("invalid decimal string: empty")
	}
	if s[0] == '+' || s[0] == '-' {
		return errors.New("invalid decimal string: must not be signed")
	}

	var x big.Int
	if _, ok := x.SetString(s, 10); !ok {
		return errors.New("invalid decimal string: must contain only decimal digits")
	}

	return x256.setBigInt(&x)
}

// NewUint256FromUint64 returns a new [Uint256] from a uint64.
func NewUint256FromUint64(i uint64) Uint256 {
	var x256 Uint256
	x256.x.SetUint64(i)

	return x256
}

// BigInt returns a copy of the underlying [big.Int].
func (x256 Uint256) BigInt() *big.Int {
	var x big.Int
	x.Set(&x256.x)

	return &x
}

// String implements [fmt.Stringer].
// It encodes x256 as a 0x-prefixed lowercase hexadecimal string with no leading zeros (zero is "0x0").
func (x256 Uint256) String() string {
	return "0x" + x256.x.Text(16)
}

// Value implements [driver.Valuer].
// It encodes x256 as a minimal big-endian []byte (never nil); zero becomes a single 0x00 byte.
func (x256 Uint256) Value() (driver.Value, error) {
	b := x256.x.Bytes()
	if len(b) == 0 {
		b = []byte{0x00}
	}

	return b, nil
}

// Scan implements [sql.Scanner].
// It decodes a big-endian []byte (length 1-32) into x256.
func (x256 *Uint256) Scan(src any) error {
	if src == nil {
		return errors.New("unsupported source: nil")
	}

	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("unsupported source type: %T", src)
	}
	if len(b) == 0 {
		return errors.New("invalid source: empty bytes")
	}

	var x big.Int
	x.SetBytes(b)

	return x256.setBigInt(&x)
}

// MarshalText implements [encoding.TextMarshaler].
// It encodes x256 as a 0x-prefixed lowercase hexadecimal string with no leading zeros (zero is "0x0").
func (x256 Uint256) MarshalText() ([]byte, error) {
	return []byte(x256.String()), nil
}

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes x256 as a quoted 0x-prefixed lowercase hexadecimal string with no leading zeros (zero is "0x0") and writes it to enc.
func (x256 Uint256) MarshalJSONTo(enc *jsontext.Encoder) error {
	b, _ := x256.MarshalText()

	return json.MarshalEncode(enc, string(b))
}

// MarshalJSON implements [json.Marshaler].
// It is like [Uint256.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (x256 Uint256) MarshalJSON() ([]byte, error) {
	return json.Marshal(x256)
}

// MarshalGQL implements [graphql.Marshaler].
// It encodes x256 as a quoted 0x-prefixed lowercase hexadecimal string with no leading zeros (zero is "0x0") and writes it to w.
func (x256 Uint256) MarshalGQL(w io.Writer) {
	_, _ = io.WriteString(w, strconv.Quote(x256.String()))
}

// UnmarshalText implements [encoding.TextUnmarshaler].
// It decodes one of the following values into x256:
//   - hexadecimal integer string (see [NewUint256FromHex])
//   - decimal integer string (see [NewUint256FromDecimal])
func (x256 *Uint256) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return errors.New("invalid string: empty")
	}

	s := string(text)

	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		return x256.setHex(s)
	}

	return x256.setDecimal(s)
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes one of the following values from dec into x256:
//   - quoted hexadecimal integer string (see [NewUint256FromHex])
//   - quoted decimal integer string (see [NewUint256FromDecimal])
//   - unquoted decimal integer string (see [NewUint256FromDecimal])
func (x256 *Uint256) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch k := dec.PeekKind(); k {
	case jsontext.KindString:
		var s string
		if err := json.UnmarshalDecode(dec, &s); err != nil {
			return fmt.Errorf("invalid string: %w", err)
		}

		return x256.UnmarshalText([]byte(s))

	case jsontext.KindNumber:
		v, err := dec.ReadValue()
		if err != nil {
			return fmt.Errorf("failed to read value: %w", err)
		}

		return x256.setDecimal(string(v))

	default:
		return fmt.Errorf("unsupported json token kind: %v", k)
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [Uint256.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (x256 *Uint256) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, x256)
}

// UnmarshalGQL implements [graphql.Unmarshaler].
// It decodes one of the following values into x256:
//   - hexadecimal integer string (see [NewUint256FromHex])
//   - decimal integer string (see [NewUint256FromDecimal])
func (x256 *Uint256) UnmarshalGQL(v any) error {
	if v == nil {
		return errors.New("unsupported value: nil")
	}

	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("unsupported value type: %T", v)
	}

	return x256.UnmarshalText([]byte(s))
}
