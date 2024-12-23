package bigutil_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	ethmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/bigutil/v2"
)

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
				"min",
				bigutil.Uint64ToUint256(0),
				[]byte{0x0},
			},
			{
				"max",
				bigutil.MustBigIntToUint256(ethmath.MaxBig256),
				[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := tc.in.Value()
				require.Nil(t, err)

				require.Equal(t, tc.out, v)
			})
		}
	})
}

func TestUint256Scan(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			out  bigutil.Uint256
		}{
			{
				"min",
				[]byte{0x0},
				bigutil.Uint64ToUint256(0),
			},
			{
				"max",
				[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				bigutil.MustBigIntToUint256(ethmath.MaxBig256),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var i bigutil.Uint256
				require.Nil(t, i.Scan(tc.in))

				require.Zero(t, i.BigInt().Cmp(tc.out.BigInt()))
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
				"min",
				bigutil.Uint64ToUint256(0),
				[]byte(`"0x0"`),
			},
			{
				"max",
				bigutil.MustBigIntToUint256(ethmath.MaxBig256),
				[]byte(`"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"`),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := json.Marshal(tc.in)
				require.Nil(t, err)

				require.Equal(t, tc.out, b)
			})
		}
	})
}

func TestUint256UnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  bigutil.Uint256
		}{
			{
				"min (hexadecimal string)",
				[]byte(`"0x0"`),
				bigutil.Uint64ToUint256(0),
			},
			{
				"min (hexadecimal string with leading zero digits)",
				[]byte(`"0x0000000000000000000000000000000000000000000000000000000000000000"`),
				bigutil.Uint64ToUint256(0),
			},
			{
				"one (hexadecimal string)",
				[]byte(`"0x1"`),
				bigutil.Uint64ToUint256(1),
			},
			{
				"one (hexadecimal string with leading zero digits)",
				[]byte(`"0x0000000000000000000000000000000000000000000000000000000000000001"`),
				bigutil.Uint64ToUint256(1),
			},
			{
				"max (hexadecimal string)",
				[]byte(`"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"`),
				bigutil.MustBigIntToUint256(ethmath.MaxBig256),
			},
			{
				"min (decimal string)",
				[]byte(`"0"`),
				bigutil.Uint64ToUint256(0),
			},
			{
				"max (decimal string)",
				[]byte(`"115792089237316195423570985008687907853269984665640564039457584007913129639935"`),
				bigutil.MustBigIntToUint256(ethmath.MaxBig256),
			},
			{
				"min (number)",
				[]byte(`0`),
				bigutil.Uint64ToUint256(0),
			},
			{
				"max (number)",
				[]byte(`115792089237316195423570985008687907853269984665640564039457584007913129639935`),
				bigutil.MustBigIntToUint256(ethmath.MaxBig256),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var i bigutil.Uint256
				require.Nil(t, json.Unmarshal(tc.in, &i))

				require.Zero(t, i.BigInt().Cmp(tc.out.BigInt()))
			})
		}
	})
}
