package bigutil_test

import (
	"database/sql/driver"
	"encoding/json"
	"math/big"
	"testing"

	ethmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/bigutil/v3"
)

func TestNewUint256(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *big.Int
		}{
			{
				"negative",
				big.NewInt(-1),
			},
			{
				"too large",
				new(big.Int).Add(ethmath.MaxBig256, big.NewInt(1)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				_, err := bigutil.NewUint256(tc.in)
				require.Error(t, err)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *big.Int
			out  bigutil.Uint256
		}{
			{
				"zero",
				big.NewInt(0),
				bigutil.NewUint256FromUint64(0),
			},
			{
				"max",
				ethmath.MaxBig256,
				bigutil.MustNewUint256(ethmath.MaxBig256),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				x256, err := bigutil.NewUint256(tc.in)
				require.NoError(t, err)

				require.Zero(t, x256.BigInt().Cmp(tc.out.BigInt()))
			})
		}
	})
}

func TestNewUint256FromHex(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   string
		}{
			{
				"0x",
				"0x",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				_, err := bigutil.NewUint256FromHex(tc.in)
				require.Error(t, err)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   string
			out  bigutil.Uint256
		}{
			{
				"zero",
				"0x0",
				bigutil.NewUint256FromUint64(0),
			},
			{
				"zero (with leading zero digits)",
				"0x0000000000000000000000000000000000000000000000000000000000000000",
				bigutil.NewUint256FromUint64(0),
			},
			{
				"one",
				"0x1",
				bigutil.NewUint256FromUint64(1),
			},
			{
				"one (with leading zero digits)",
				"0x0000000000000000000000000000000000000000000000000000000000000001",
				bigutil.NewUint256FromUint64(1),
			},
			{
				"max",
				"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
				bigutil.MustNewUint256(ethmath.MaxBig256),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				x256, err := bigutil.NewUint256FromHex(tc.in)
				require.NoError(t, err)

				require.Zero(t, x256.BigInt().Cmp(tc.out.BigInt()))
			})
		}
	})
}

func TestUint256Value(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   bigutil.Uint256
			out  driver.Value
		}{
			{
				"zero value",
				bigutil.Uint256{},
				[]byte{0x0},
			},
			{
				"zero",
				bigutil.NewUint256FromUint64(0),
				[]byte{0x0},
			},
			{
				"max",
				bigutil.MustNewUint256(ethmath.MaxBig256),
				[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := tc.in.Value()
				require.NoError(t, err)

				require.Equal(t, tc.out, v)
			})
		}
	})
}

func TestUint256Scan(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
		}{
			{
				"nil",
				nil,
			},
			{
				"int64",
				int64(0),
			},
			{
				"uint64",
				uint64(0),
			},
			{
				"empty bytes",
				[]byte{},
			},
			{
				"too long bytes",
				[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00},
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var x256 bigutil.Uint256
				{
					err := x256.Scan(tc.in)
					require.Error(t, err)
				}
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			out  bigutil.Uint256
		}{
			{
				"zero",
				[]byte{0x0},
				bigutil.NewUint256FromUint64(0),
			},
			{
				"max",
				[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				bigutil.MustNewUint256(ethmath.MaxBig256),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var x256 bigutil.Uint256
				{
					err := x256.Scan(tc.in)
					require.NoError(t, err)
				}

				require.Zero(t, x256.BigInt().Cmp(tc.out.BigInt()))
			})
		}
	})
}

func TestUint256MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   bigutil.Uint256
			out  []byte
		}{
			{
				"zero value",
				bigutil.Uint256{},
				[]byte(`"0x0"`),
			},
			{
				"zero",
				bigutil.NewUint256FromUint64(0),
				[]byte(`"0x0"`),
			},
			{
				"max",
				bigutil.MustNewUint256(ethmath.MaxBig256),
				[]byte(`"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"`),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := json.Marshal(tc.in)
				require.NoError(t, err)

				require.Equal(t, tc.out, b)
			})
		}
	})
}

func TestUint256UnmarshalJSON(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
		}{
			{
				"0x",
				[]byte(`"0x"`),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var x256 bigutil.Uint256
				{
					err := json.Unmarshal(tc.in, &x256)
					require.Error(t, err)
				}
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  bigutil.Uint256
		}{
			{
				"zero (hexadecimal string)",
				[]byte(`"0x0"`),
				bigutil.NewUint256FromUint64(0),
			},
			{
				"zero (hexadecimal string with leading zero digits)",
				[]byte(`"0x0000000000000000000000000000000000000000000000000000000000000000"`),
				bigutil.NewUint256FromUint64(0),
			},
			{
				"one (hexadecimal string)",
				[]byte(`"0x1"`),
				bigutil.NewUint256FromUint64(1),
			},
			{
				"one (hexadecimal string with leading zero digits)",
				[]byte(`"0x0000000000000000000000000000000000000000000000000000000000000001"`),
				bigutil.NewUint256FromUint64(1),
			},
			{
				"max (hexadecimal string)",
				[]byte(`"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"`),
				bigutil.MustNewUint256(ethmath.MaxBig256),
			},
			{
				"zero (decimal string)",
				[]byte(`"0"`),
				bigutil.NewUint256FromUint64(0),
			},
			{
				"max (decimal string)",
				[]byte(`"115792089237316195423570985008687907853269984665640564039457584007913129639935"`),
				bigutil.MustNewUint256(ethmath.MaxBig256),
			},
			{
				"zero (number)",
				[]byte(`0`),
				bigutil.NewUint256FromUint64(0),
			},
			{
				"max (number)",
				[]byte(`115792089237316195423570985008687907853269984665640564039457584007913129639935`),
				bigutil.MustNewUint256(ethmath.MaxBig256),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var x256 bigutil.Uint256
				{
					err := json.Unmarshal(tc.in, &x256)
					require.NoError(t, err)
				}

				require.Zero(t, x256.BigInt().Cmp(tc.out.BigInt()))
			})
		}
	})
}
