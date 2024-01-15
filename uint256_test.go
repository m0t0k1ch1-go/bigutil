package bigutil_test

import (
	"database/sql/driver"
	"encoding/json"
	"math/big"
	"testing"

	ethmath "github.com/ethereum/go-ethereum/common/math"

	"github.com/m0t0k1ch1-go/bigutil/v2"
	"github.com/m0t0k1ch1-go/bigutil/v2/internal/testutil"
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
				bigutil.MustBigIntToUint256(big.NewInt(0)),
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
				if err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, v)
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
				bigutil.MustBigIntToUint256(big.NewInt(0)),
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
				if err := i.Scan(tc.in); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, i.BigInt().Cmp(tc.out.BigInt()), 0)
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
				bigutil.MustBigIntToUint256(big.NewInt(0)),
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
				if err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, b)
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
				"min (hex)",
				[]byte(`"0x0"`),
				bigutil.MustBigIntToUint256(big.NewInt(0)),
			},
			{
				"max (hex)",
				[]byte(`"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"`),
				bigutil.MustBigIntToUint256(ethmath.MaxBig256),
			},
			{
				"min (decimal)",
				[]byte(`"0"`),
				bigutil.MustBigIntToUint256(big.NewInt(0)),
			},
			{
				"max (decimal)",
				[]byte(`"115792089237316195423570985008687907853269984665640564039457584007913129639935"`),
				bigutil.MustBigIntToUint256(ethmath.MaxBig256),
			},
			{
				"min (number)",
				[]byte(`0`),
				bigutil.MustBigIntToUint256(big.NewInt(0)),
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
				if err := json.Unmarshal(tc.in, &i); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, i.BigInt().Cmp(tc.out.BigInt()), 0)
			})
		}
	})
}
