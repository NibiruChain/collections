package collections

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// StringKeyEncoder can be used to encode string keys.
	StringKeyEncoder KeyEncoder[string] = stringKey{}
	// AccAddressKeyEncoder can be used to encode sdk.AccAddress keys.
	AccAddressKeyEncoder KeyEncoder[sdk.AccAddress] = accAddressKey{}
	// TimeKeyEncoder can be used to encode time.Time keys.
	TimeKeyEncoder KeyEncoder[time.Time] = timeKey{}
	// Uint64KeyEncoder can be used to encode uint64 keys.
	Uint64KeyEncoder KeyEncoder[uint64] = uint64Key{}
	// ValAddressKeyEncoder can be used to encode sdk.ValAddress keys.
	ValAddressKeyEncoder KeyEncoder[sdk.ValAddress] = valAddressKeyEncoder{}
	// ConsAddressKeyEncoder can be used to encode sdk.ConsAddress keys.
	ConsAddressKeyEncoder KeyEncoder[sdk.ConsAddress] = consAddressKeyEncoder{}
	// SdkDecKeyEncoder can be used to encode sdk.Dec keys.
	SdkDecKeyEncoder KeyEncoder[sdk.Dec] = sdkDecKeyEncoder{}
)

// testing sentinel errors
var (
	errInvalidKey = errors.New("collections: invalid key")

	errInvalidStringKeySize       = errors.New("collections: invalid string key bytes buffer size")
	errInvalidStringKeyNullChar   = errors.New("collections: invalid string key contains null character")
	errInvalidStringKeyNoNullChar = errors.New("collections: invalid string key bytes buffer is not null terminated")

	errInvalidUint64KeySize = errors.New("collections: invalid uint64 key bytes buffer size")
)

type stringKey struct{}

func (stringKey) Encode(s string) ([]byte, error) {
	if err := validString(s); err != nil {
		return nil, err
	}
	return append([]byte(s), 0), nil // null terminate it for safe prefixing
}

func (stringKey) Decode(b []byte) (int, string, error) {
	l := len(b)
	if l < 2 {
		return 0, "", errInvalidStringKeySize
	}
	for i, c := range b {
		if c == 0 {
			return i + 1, string(b[:i]), nil
		}
	}
	return 0, "", errInvalidStringKeyNoNullChar
}

type uint64Key struct{}

func (uint64Key) Stringify(u uint64) string { return strconv.FormatUint(u, 10) }
func (uint64Key) Encode(u uint64) ([]byte, error) {
	return sdk.Uint64ToBigEndian(u), nil
}
func (uint64Key) Decode(b []byte) (int, uint64, error) {
	if len(b) != 8 {
		return 0, 0, errInvalidUint64KeySize
	}
	return 8, sdk.BigEndianToUint64(b), nil
}

type timeKey struct{}

func (timeKey) Stringify(t time.Time) string       { return t.String() }
func (timeKey) Encode(t time.Time) ([]byte, error) { return sdk.FormatTimeBytes(t), nil }
func (timeKey) Decode(b []byte) (int, time.Time, error) {
	t, err := sdk.ParseTimeBytes(b)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("%w: TimeKey: %s", errInvalidKey, err)
	}
	return len(b), t, nil
}

type accAddressKey struct{}

func (accAddressKey) Stringify(addr sdk.AccAddress) string { return addr.String() }
func (accAddressKey) Encode(addr sdk.AccAddress) ([]byte, error) {
	return StringKeyEncoder.Encode(addr.String())
}
func (accAddressKey) Decode(b []byte) (int, sdk.AccAddress, error) {
	i, s, err := StringKeyEncoder.Decode(b)
	if err != nil {
		return 0, nil, fmt.Errorf("%w: AccAddressKey", err)
	}
	return i, sdk.MustAccAddressFromBech32(s), nil
}

type valAddressKeyEncoder struct{}

func (v valAddressKeyEncoder) Encode(key sdk.ValAddress) ([]byte, error) {
	return StringKeyEncoder.Encode(key.String())
}
func (v valAddressKeyEncoder) Decode(b []byte) (int, sdk.ValAddress, error) {
	r, s, err := StringKeyEncoder.Decode(b)
	if err != nil {
		return 0, nil, fmt.Errorf("%w: ValAddressKey", err)
	}
	valAddr, err := sdk.ValAddressFromBech32(s)
	if err != nil {
		return 0, nil, fmt.Errorf("%w: ValAddressKey: %s", errInvalidKey, err)
	}
	return r, valAddr, nil
}
func (v valAddressKeyEncoder) Stringify(key sdk.ValAddress) string { return key.String() }

func (stringKey) Stringify(s string) string {
	return s
}

func validString(s string) error {
	for i, c := range s {
		if c == 0 {
			return fmt.Errorf("%w at index %d: %s", errInvalidStringKeyNullChar, i, s)
		}
	}
	return nil
}

type consAddressKeyEncoder struct{}

func (consAddressKeyEncoder) Encode(key sdk.ConsAddress) ([]byte, error) {
	return StringKeyEncoder.Encode(key.String())
}
func (consAddressKeyEncoder) Decode(b []byte) (int, sdk.ConsAddress, error) {
	r, s, err := StringKeyEncoder.Decode(b)
	if err != nil {
		return 0, nil, fmt.Errorf("%w: ConsAddressKey", err)
	}
	consAddr, err := sdk.ConsAddressFromBech32(s)
	if err != nil {
		return 0, nil, fmt.Errorf("%w: ConsAddressKey: %s", errInvalidKey, err)
	}
	return r, consAddr, nil
}
func (consAddressKeyEncoder) Stringify(key sdk.ConsAddress) string { return key.String() }

type sdkDecKeyEncoder struct{}

func (sdkDecKeyEncoder) Stringify(key sdk.Dec) string { return key.String() }

func (sdkDecKeyEncoder) Encode(key sdk.Dec) ([]byte, error) {
	bz, err := key.Marshal()
	if err != nil {
		return nil, fmt.Errorf("%w: DecKey: %s", errInvalidKey, err)
	}
	return bz, nil
}
func (sdkDecKeyEncoder) Decode(b []byte) (int, sdk.Dec, error) {
	var dec sdk.Dec
	if err := dec.Unmarshal(b); err != nil {
		return 0, dec, fmt.Errorf("%w: DecKey bytes: %s", errInvalidKey, err)
	}

	return len(b), dec, nil
}
