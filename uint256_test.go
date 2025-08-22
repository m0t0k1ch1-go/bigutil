package bigutil_test

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/bigutil/v3"
)

var (
	maxUint256 = new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
)

func TestUint256(t *testing.T) {
	var x256 bigutil.Uint256
	require.Implements(t, (*fmt.Stringer)(nil), &x256)
	require.Implements(t, (*driver.Valuer)(nil), &x256)
	require.Implements(t, (*sql.Scanner)(nil), &x256)
	require.Implements(t, (*encoding.TextMarshaler)(nil), &x256)
	require.Implements(t, (*graphql.Marshaler)(nil), &x256)
	require.Implements(t, (*encoding.TextUnmarshaler)(nil), &x256)
	require.Implements(t, (*json.Unmarshaler)(nil), &x256)
	require.Implements(t, (*graphql.Unmarshaler)(nil), &x256)
}

func TestNewUint256(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *big.Int
			want string
		}{
			{
				"nil",
				nil,
				"invalid big.Int: nil",
			},
			{
				"negative",
				big.NewInt(-1),
				"invalid big.Int: negative",
			},
			{
				"exceeds 256 bits",
				new(big.Int).Add(maxUint256, big.NewInt(1)),
				"invalid big.Int: exceeds 256 bits",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				_, err := bigutil.NewUint256(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *big.Int
			want string
		}{
			{
				"zero",
				big.NewInt(0),
				"0x0",
			},
			{
				"one",
				big.NewInt(1),
				"0x1",
			},
			{
				"max",
				new(big.Int).Set(maxUint256),
				"0x" + strings.Repeat("f", 64),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				x256, err := bigutil.NewUint256(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, x256.String())

				tc.in.SetInt64(-1)

				require.Equal(t, tc.want, x256.String())
			})
		}
	})
}

func TestMustNewUint256(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *big.Int
			want string
		}{
			{
				"nil",
				nil,
				"invalid big.Int: nil",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				require.PanicsWithError(t, tc.want, func() {
					bigutil.MustNewUint256(tc.in)
				})
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *big.Int
			want string
		}{
			{
				"zero",
				big.NewInt(0),
				"0x0",
			},
		}

		for _, tc := range tcs {
			x256 := bigutil.MustNewUint256(tc.in)
			require.Equal(t, tc.want, x256.String())
		}
	})
}

func TestNewUint256FromHex(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   string
			want string
		}{
			{
				"empty",
				"",
				"invalid hex string: empty",
			},
			{
				"missing 0x/0X prefix",
				"0",
				"invalid hex string: missing 0x/0X prefix",
			},
			{
				"missing hex digits after 0x prefix",
				"0x",
				"invalid hex string: missing hex digits after 0x/0X prefix",
			},
			{
				"missing hex digits after 0X prefix",
				"0X",
				"invalid hex string: missing hex digits after 0x/0X prefix",
			},
			{
				"contains non-hex characters",
				"0xg",
				"invalid hex string",
			},
			{
				"exceeds 256 bits",
				"0x1" + strings.Repeat("0", 64),
				"invalid big.Int: exceeds 256 bits",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				_, err := bigutil.NewUint256FromHex(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   string
			want string
		}{
			{
				"0x-prefixed zero",
				"0x0",
				"0x0",
			},
			{
				"0X-prefixed zero",
				"0X0",
				"0x0",
			},
			{
				"0x-prefixed zero with leading zeros",
				"0x" + strings.Repeat("0", 64),
				"0x0",
			},
			{
				"0X-prefixed zero with leading zeros",
				"0X" + strings.Repeat("0", 64),
				"0x0",
			},
			{
				"0x-prefixed one with leading zeros",
				"0x" + strings.Repeat("0", 63) + "1",
				"0x1",
			},
			{
				"0X-prefixed one with leading zeros",
				"0X" + strings.Repeat("0", 63) + "1",
				"0x1",
			},
			{
				"0x-prefixed mixedcase max",
				"0x" + strings.Repeat("fF", 32),
				"0x" + strings.Repeat("f", 64),
			},
			{
				"0X-prefixed mixedcase max",
				"0X" + strings.Repeat("fF", 32),
				"0x" + strings.Repeat("f", 64),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				x256, err := bigutil.NewUint256FromHex(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, x256.String())
			})
		}
	})
}

func TestMustNewUint256FromHex(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		tcs := []struct {
			name string
			in   string
			want string
		}{
			{
				"empty",
				"",
				"invalid hex string: empty",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				require.PanicsWithError(t, tc.want, func() {
					bigutil.MustNewUint256FromHex(tc.in)
				})
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   string
			want string
		}{
			{
				"0x-prefixed zero",
				"0x0",
				"0x0",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				x256 := bigutil.MustNewUint256FromHex(tc.in)
				require.Equal(t, tc.want, x256.String())
			})
		}
	})
}

func TestUint256_BigInt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   bigutil.Uint256
			want string
		}{
			{
				"zero",
				bigutil.NewUint256FromUint64(0),
				"0x0",
			},
			{
				"one",
				bigutil.NewUint256FromUint64(1),
				"0x1",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				x := tc.in.BigInt()
				require.Equal(t, tc.want, "0x"+x.Text(16))

				x.SetInt64(-1)

				s := tc.in.String()
				require.Equal(t, tc.want, s)
			})
		}
	})
}

func TestUint256_Value(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   bigutil.Uint256
			want driver.Value
		}{
			{
				"zero value",
				bigutil.Uint256{},
				[]byte{0x00},
			},
			{
				"zero",
				bigutil.NewUint256FromUint64(0),
				[]byte{0x00},
			},
			{
				"one",
				bigutil.NewUint256FromUint64(1),
				[]byte{0x1},
			},
			{
				"max",
				bigutil.MustNewUint256(maxUint256),
				bytes.Repeat([]byte{0xff}, 32),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := tc.in.Value()
				require.NoError(t, err)
				require.Equal(t, tc.want, v)
			})
		}
	})
}

func TestUint256_Scan(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"nil",
				nil,
				"invalid source: nil",
			},
			{
				"int64",
				int64(0),
				"unsupported source type: int64",
			},
			{
				"string",
				"0x0",
				"unsupported source type: string",
			},
			{
				"[]byte: empty",
				[]byte{},
				"invalid source: empty []byte",
			},
			{
				"[]byte: exceeds 256 bits",
				append([]byte{0x01}, bytes.Repeat([]byte{0x00}, 32)...),
				"invalid source: invalid big.Int: exceeds 256 bits",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var x256 bigutil.Uint256
				err := x256.Scan(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"[]byte: zero",
				[]byte{0x00},
				"0x0",
			},
			{
				"[]byte: zero with leading zeros",
				bytes.Repeat([]byte{0x00}, 32),
				"0x0",
			},
			{
				"[]byte: one",
				[]byte{0x01},
				"0x1",
			},
			{
				"[]byte: one with leading zeros",
				append(bytes.Repeat([]byte{0x00}, 31), 0x01),
				"0x1",
			},
			{
				"[]byte: max",
				bytes.Repeat([]byte{0xff}, 32),
				"0x" + strings.Repeat("f", 64),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var x256 bigutil.Uint256
				err := x256.Scan(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, x256.String())
			})
		}
	})
}

func TestUint256_MarshalText(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   bigutil.Uint256
			out  []byte
		}{
			{
				"zero value",
				bigutil.Uint256{},
				[]byte("0x0"),
			},
			{
				"zero",
				bigutil.NewUint256FromUint64(0),
				[]byte("0x0"),
			},
			{
				"one",
				bigutil.NewUint256FromUint64(1),
				[]byte("0x1"),
			},
			{
				"max",
				bigutil.MustNewUint256(maxUint256),
				[]byte("0x" + strings.Repeat("f", 64)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := tc.in.MarshalText()
				require.NoError(t, err)
				require.Equal(t, tc.out, b)
			})
		}
	})
}

func TestUint256_MarshalGQL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   bigutil.Uint256
			want string
		}{
			{
				"zero value",
				bigutil.Uint256{},
				`"0x0"`,
			},
			{
				"zero",
				bigutil.NewUint256FromUint64(0),
				`"0x0"`,
			},
			{
				"one",
				bigutil.NewUint256FromUint64(1),
				`"0x1"`,
			},
			{
				"max",
				bigutil.MustNewUint256(maxUint256),
				`"0x` + strings.Repeat("f", 64) + `"`,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var buf bytes.Buffer
				tc.in.MarshalGQL(&buf)
				require.Equal(t, tc.want, buf.String())
			})
		}
	})
}

func TestUint256_UnmarshalJSON(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"empty",
				[]byte{},
				"invalid json value: empty",
			},
			{
				"null",
				[]byte(`null`),
				"invalid json value: null",
			},
			{
				"number: negative",
				[]byte(`-1`),
				"invalid json number: invalid big.Int: negative",
			},
			{
				"number: exceeds 256 bits",
				[]byte(`115792089237316195423570985008687907853269984665640564039457584007913129639936`),
				"invalid json number: invalid big.Int: exceeds 256 bits",
			},
			{
				"number: fractional",
				[]byte(`0.0`),
				"invalid json number",
			},
			{
				"number: exponential",
				[]byte(`0e0`),
				"invalid json number",
			},
			{
				"string: empty",
				[]byte(`""`),
				"invalid json string: invalid string: empty",
			},
			{
				"string: invalid decimal",
				[]byte(`"invalid"`),
				"invalid json string: invalid decimal string",
			},
			{
				"string: negative decimal",
				[]byte(`"-1"`),
				"invalid json string: invalid big.Int: negative",
			},
			{
				"string: missing hex digits after 0x prefix",
				[]byte(`"0x"`),
				"invalid json string: invalid hex string: missing hex digits after 0x/0X prefix",
			},
			{
				"string: missing hex digits after 0X prefix",
				[]byte(`"0X"`),
				"invalid json string: invalid hex string: missing hex digits after 0x/0X prefix",
			},
			{
				"string: hex contains invalid escape sequences",
				[]byte(`"0x\x"`),
				"invalid json string",
			},
			{
				"string: hex contains non-hex characters",
				[]byte(`"0xg"`),
				"invalid json string: invalid hex string",
			},
			{
				"string: hex exceeds 256 bits",
				[]byte(`"0x1` + strings.Repeat("0", 64) + `"`),
				"invalid json string: invalid big.Int: exceeds 256 bits",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var x256 bigutil.Uint256
				err := x256.UnmarshalJSON(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"number: zero",
				[]byte(`0`),
				"0x0",
			},
			{
				"number: one",
				[]byte(`1`),
				"0x1",
			},
			{
				"number: max",
				[]byte(`115792089237316195423570985008687907853269984665640564039457584007913129639935`),
				"0x" + strings.Repeat("f", 64),
			},
			{
				"string: decimal zero",
				[]byte(`"0"`),
				"0x0",
			},
			{
				"string: decimal one",
				[]byte(`"1"`),
				"0x1",
			},
			{
				"string: decimal max",
				[]byte(`"115792089237316195423570985008687907853269984665640564039457584007913129639935"`),
				"0x" + strings.Repeat("f", 64),
			},
			{
				"string: 0x-prefixed hex zero",
				[]byte(`"0x0"`),
				"0x0",
			},
			{
				"string: 0X-prefixed hex zero",
				[]byte(`"0X0"`),
				"0x0",
			},
			{
				"string: 0x-prefixed hex zero with leading zeros",
				[]byte(`"0x` + strings.Repeat("0", 64) + `"`),
				"0x0",
			},
			{
				"string: 0X-prefixed hex zero with leading zeros",
				[]byte(`"0X` + strings.Repeat("0", 64) + `"`),
				"0x0",
			},
			{
				"string: 0x-prefixed hex one",
				[]byte(`"0x1"`),
				"0x1",
			},
			{
				"string: 0X-prefixed hex one",
				[]byte(`"0X1"`),
				"0x1",
			},
			{
				"string: 0x-prefixed hex one with leading zeros",
				[]byte(`"0x` + strings.Repeat("0", 63) + `1"`),
				"0x1",
			},
			{
				"string: 0X-prefixed hex one with leading zeros",
				[]byte(`"0X` + strings.Repeat("0", 63) + `1"`),
				"0x1",
			},
			{
				"string: 0x-prefixed mixedcase hex max",
				[]byte(`"0x` + strings.Repeat("fF", 32) + `"`),
				"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			},
			{
				"string: 0X-prefixed mixedcase hex max",
				[]byte(`"0X` + strings.Repeat("fF", 32) + `"`),
				"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var x256 bigutil.Uint256
				err := x256.UnmarshalJSON(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, x256.String())
			})
		}
	})
}

func TestUint256_UnmarshalGQL(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"nil",
				nil,
				"invalid graphql value: nil",
			},
			{
				"int",
				int(0),
				"unsupported graphql value type: int",
			},
			{
				"string: empty",
				"",
				"invalid graphql string: invalid string: empty",
			},
			{
				"string: invalid decimal",
				"invalid",
				"invalid graphql string: invalid decimal string",
			},
			{
				"string: negative decimal",
				"-1",
				"invalid graphql string: invalid big.Int: negative",
			},
			{
				"string: missing hex digits after 0x prefix",
				"0x",
				"invalid graphql string: invalid hex string: missing hex digits after 0x/0X prefix",
			},
			{
				"string: missing hex digits after 0X prefix",
				"0X",
				"invalid graphql string: invalid hex string: missing hex digits after 0x/0X prefix",
			},
			{
				"string: hex contains non-hex characters",
				"0xg",
				"invalid graphql string: invalid hex string",
			},
			{
				"string: hex exceeds 256 bits",
				"0x1" + strings.Repeat("0", 64),
				"invalid graphql string: invalid big.Int: exceeds 256 bits",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var x256 bigutil.Uint256
				err := x256.UnmarshalGQL(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"string: decimal zero",
				"0",
				"0x0",
			},
			{
				"string: decimal one",
				"1",
				"0x1",
			},
			{
				"string: decimal max",
				"115792089237316195423570985008687907853269984665640564039457584007913129639935",
				"0x" + strings.Repeat("f", 64),
			},
			{
				"string: 0x-prefixed hex zero",
				"0x0",
				"0x0",
			},
			{
				"string: 0X-prefixed hex zero",
				"0X0",
				"0x0",
			},
			{
				"string: 0x-prefixed hex zero with leading zeros",
				"0x" + strings.Repeat("0", 64),
				"0x0",
			},
			{
				"string: 0X-prefixed hex zero with leading zeros",
				"0X" + strings.Repeat("0", 64),
				"0x0",
			},
			{
				"string: 0x-prefixed hex one",
				"0x1",
				"0x1",
			},
			{
				"string: 0X-prefixed hex one",
				"0X1",
				"0x1",
			},
			{
				"string: 0x-prefixed hex one with leading zeros",
				"0x" + strings.Repeat("0", 63) + "1",
				"0x1",
			},
			{
				"string: 0X-prefixed hex one with leading zeros",
				"0X" + strings.Repeat("0", 63) + "1",
				"0x1",
			},
			{
				"string: 0x-prefixed mixedcase hex max",
				"0x" + strings.Repeat("fF", 32),
				"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			},
			{
				"string: 0X-prefixed mixedcase hex max",
				"0X" + strings.Repeat("fF", 32),
				"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var x256 bigutil.Uint256
				err := x256.UnmarshalGQL(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, x256.String())
			})
		}
	})
}
