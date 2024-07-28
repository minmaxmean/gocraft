package datatypes

import (
	"fmt"
	"io"
)

const VARINT_SEGMENT_BITS = 0x7F
const VARINT_CONTINUE_BIT = 0x80

type VarInt int32

func (value *VarInt) ReadFrom(r io.Reader) (int64, error) {
	pos := 0
	nn := int64(0)
	*value = 0
	b := make([]byte, 1)
	for {
		n, err := r.Read(b)
		nn += int64(n)
		if err != nil {
			return nn, err
		}
		*value |= VarInt(b[0]&VARINT_SEGMENT_BITS) << pos
		if b[0]&VARINT_CONTINUE_BIT == 0 {
			break
		}
		pos += 7
		if pos >= 32 {
			return nn, fmt.Errorf("VarInt too big")
		}
	}
	return nn, nil
}

func (v VarInt) WriteTo(w io.Writer) (int64, error) {
	nn := int64(0)
	value := uint32(v)
	for {
		b := byte(value) & VARINT_SEGMENT_BITS
		value >>= 7
		if value != 0 {
			b |= VARINT_CONTINUE_BIT
		}

		n, err := w.Write([]byte{b})
		nn += int64(n)
		if err != nil {
			return nn, err
		}
		if value == 0 {
			break
		}
	}
	return nn, nil
}

const STRING_MAX_LEN = 32767

type String string

func (value *String) ReadFrom(r io.Reader) (int64, error) {
	var strlen VarInt
	nn, err := strlen.ReadFrom(r)
	if err != nil {
		return nn, err
	}
	if !(1 <= strlen && strlen <= STRING_MAX_LEN) {
		return nn, fmt.Errorf("incorrect string length: want %d max: %d", strlen, STRING_MAX_LEN)
	}
	s := make([]byte, strlen)
	n, err := io.ReadFull(r, s)
	nn += int64(n)
	if err != nil {
		return nn, err
	}
	*value = String(s)
	return nn, nil
}

func (value *String) WriteTo(w io.Writer) (int64, error) {
	strlen := VarInt(len(*value))
	if !(1 <= strlen && strlen <= STRING_MAX_LEN) {
		return 0, fmt.Errorf("incorrect string length: want %d max: %d", strlen, STRING_MAX_LEN)
	}
	nn, err := strlen.WriteTo(w)
	if err != nil {
		return nn, err
	}

	n, err := w.Write([]byte(*value))
	nn += int64(n)
	if err != nil {
		return nn, err
	}
	return nn, nil
}

type UShort int16

func (value *UShort) ReadFrom(r io.Reader) (int64, error) {
	b := make([]byte, 2)
	n, err := io.ReadFull(r, b)
	nn := int64(n)
	if err != nil {
		return nn, err
	}
	*value = UShort(b[0])<<8 | UShort(b[1])
	return nn, nil
}