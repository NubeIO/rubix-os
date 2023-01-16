package smod

import (
	"encoding/binary"
	"fmt"
	"math"
)

func uint16ToBytes(endianness Endianness, in uint16) (out []byte) {
	out = make([]byte, 2)
	switch endianness {
	case BigEndian:
		binary.BigEndian.PutUint16(out, in)
	case LittleEndian:
		binary.LittleEndian.PutUint16(out, in)
	}

	return
}

func uint16sToBytes(endianness Endianness, in []uint16) (out []byte) {
	for i := range in {
		out = append(out, uint16ToBytes(endianness, in[i])...)
	}

	return
}

func bytesToUint16(endianness Endianness, in []byte) (out uint16) {
	if math.Mod(float64(len(in)), 2) != 0 {
		fmt.Println("MODBUS CONVERSION ERROR bytesToUint16(): length of input []byte is less than 2")
		return uint16(0)
	}
	switch endianness {
	case BigEndian:
		out = binary.BigEndian.Uint16(in)
	case LittleEndian:
		out = binary.LittleEndian.Uint16(in)
	}

	return
}

func bytesToUint16s(endianness Endianness, in []byte) (out []uint16) {
	if math.Mod(float64(len(in)), 2) != 0 {
		fmt.Println("MODBUS CONVERSION ERROR bytesToUint16s(): length of input []byte is less than 2")
		return []uint16{uint16(0)}
	}
	for i := 0; i < len(in); i += 2 {
		out = append(out, bytesToUint16(endianness, in[i:i+2]))
	}

	return
}

func bytesToInt16s(endianness Endianness, in []byte) (out []int16) {
	if math.Mod(float64(len(in)), 2) != 0 {
		fmt.Println("MODBUS CONVERSION ERROR bytesToInt16s(): length of input []byte is less than 2")
		return []int16{int16(0)}
	}
	for i := 0; i < len(in); i += 2 {
		out = append(out, int16(bytesToUint16(endianness, in[i:i+2])))
	}

	return
}

func bytesToUint32s(endianness Endianness, wordOrder WordOrder, in []byte) (out []uint32) {
	var u32 uint32

	if math.Mod(float64(len(in)), 4) != 0 {
		fmt.Println("MODBUS CONVERSION ERROR bytesToUint32s(): length of input []byte is less than 4")
		return []uint32{uint32(0)}
	}
	for i := 0; i < len(in); i += 4 {
		switch endianness {
		case BigEndian:
			if wordOrder == HighWordFirst {
				u32 = binary.BigEndian.Uint32(in[i : i+4])
			} else {
				u32 = binary.BigEndian.Uint32(
					[]byte{in[i+2], in[i+3], in[i+0], in[i+1]}) // This line panics if the length of the in ([]byte) is too short.
			}
		case LittleEndian:
			if wordOrder == LowWordFirst {
				u32 = binary.LittleEndian.Uint32(in[i : i+4])
			} else {
				u32 = binary.LittleEndian.Uint32(
					[]byte{in[i+2], in[i+3], in[i+0], in[i+1]}) // This line panics if the length of the in ([]byte) is too short.
			}
		}

		out = append(out, u32)
	}

	return
}

func bytesToInt32s(endianness Endianness, wordOrder WordOrder, in []byte) (out []int32) {
	var i32 int32

	if math.Mod(float64(len(in)), 4) != 0 {
		fmt.Println("MODBUS CONVERSION ERROR bytesToInt32s(): length of input []byte is less than 4")
		return []int32{int32(0)}
	}
	for i := 0; i < len(in); i += 4 {
		switch endianness {
		case BigEndian:
			if wordOrder == HighWordFirst {
				i32 = int32(binary.BigEndian.Uint32(in[i : i+4]))
			} else {
				i32 = int32(binary.BigEndian.Uint32(
					[]byte{in[i+2], in[i+3], in[i+0], in[i+1]})) // This line panics if the length of the in ([]byte) is too short.
			}
		case LittleEndian:
			if wordOrder == LowWordFirst {
				i32 = int32(binary.LittleEndian.Uint32(in[i : i+4]))
			} else {
				i32 = int32(binary.LittleEndian.Uint32(
					[]byte{in[i+2], in[i+3], in[i+0], in[i+1]})) // This line panics if the length of the in ([]byte) is too short.
			}
		}

		out = append(out, i32)
	}

	return
}

func uint32ToBytes(endianness Endianness, wordOrder WordOrder, in uint32) (out []byte) {
	out = make([]byte, 4)

	switch endianness {
	case BigEndian:
		binary.BigEndian.PutUint32(out, in)

		// swap words if needed
		if wordOrder == LowWordFirst {
			out[0], out[1], out[2], out[3] = out[2], out[3], out[0], out[1]
		}
	case LittleEndian:
		binary.LittleEndian.PutUint32(out, in)

		// swap words if needed
		if wordOrder == HighWordFirst {
			out[0], out[1], out[2], out[3] = out[2], out[3], out[0], out[1]
		}
	}

	return
}

func bytesToFloat32s(endianness Endianness, wordOrder WordOrder, in []byte) (out []float32) {
	var u32s []uint32

	u32s = bytesToUint32s(endianness, wordOrder, in)

	for _, u32 := range u32s {
		out = append(out, math.Float32frombits(u32))
	}

	return
}

func float32ToBytes(endianness Endianness, wordOrder WordOrder, in float32) (out []byte) {
	out = uint32ToBytes(endianness, wordOrder, math.Float32bits(in))

	return
}

func bytesToUint64s(endianness Endianness, wordOrder WordOrder, in []byte) (out []uint64) {
	var u64 uint64

	if math.Mod(float64(len(in)), 8) != 0 {
		fmt.Println("MODBUS CONVERSION ERROR bytesToUint64s(): length of input []byte is less than 8")
		return []uint64{uint64(0)}
	}

	for i := 0; i < len(in); i += 8 {
		switch endianness {
		case BigEndian:
			if wordOrder == HighWordFirst {
				u64 = binary.BigEndian.Uint64(in[i : i+8])
			} else {
				u64 = binary.BigEndian.Uint64(
					[]byte{in[i+6], in[i+7], in[i+4], in[i+5],
						in[i+2], in[i+3], in[i+0], in[i+1]})
			}
		case LittleEndian:
			if wordOrder == LowWordFirst {
				u64 = binary.LittleEndian.Uint64(in[i : i+8])
			} else {
				u64 = binary.LittleEndian.Uint64(
					[]byte{in[i+6], in[i+7], in[i+4], in[i+5],
						in[i+2], in[i+3], in[i+0], in[i+1]})
			}
		}

		out = append(out, u64)
	}

	return
}

func bytesToInt64s(endianness Endianness, wordOrder WordOrder, in []byte) (out []int64) {
	var i64 int64

	if math.Mod(float64(len(in)), 8) != 0 {
		fmt.Println("MODBUS CONVERSION ERROR bytesToInt64s(): length of input []byte is less than 8")
		return []int64{int64(0)}
	}

	for i := 0; i < len(in); i += 8 {
		switch endianness {
		case BigEndian:
			if wordOrder == HighWordFirst {
				i64 = int64(binary.BigEndian.Uint64(in[i : i+8]))
			} else {
				i64 = int64(binary.BigEndian.Uint64(
					[]byte{in[i+6], in[i+7], in[i+4], in[i+5],
						in[i+2], in[i+3], in[i+0], in[i+1]}))
			}
		case LittleEndian:
			if wordOrder == LowWordFirst {
				i64 = int64(binary.LittleEndian.Uint64(in[i : i+8]))
			} else {
				i64 = int64(binary.LittleEndian.Uint64(
					[]byte{in[i+6], in[i+7], in[i+4], in[i+5],
						in[i+2], in[i+3], in[i+0], in[i+1]}))
			}
		}
		out = append(out, i64)
	}

	return
}

func uint64ToBytes(endianness Endianness, wordOrder WordOrder, in uint64) (out []byte) {
	out = make([]byte, 8)

	switch endianness {
	case BigEndian:
		binary.BigEndian.PutUint64(out, in)

		// swap words if needed
		if wordOrder == LowWordFirst {
			out[0], out[1], out[2], out[3], out[4], out[5], out[6], out[7] =
				out[6], out[7], out[4], out[5], out[2], out[3], out[0], out[1]
		}
	case LittleEndian:
		binary.LittleEndian.PutUint64(out, in)

		// swap words if needed
		if wordOrder == HighWordFirst {
			out[0], out[1], out[2], out[3], out[4], out[5], out[6], out[7] =
				out[6], out[7], out[4], out[5], out[2], out[3], out[0], out[1]
		}
	}

	return
}

func bytesToFloat64s(endianness Endianness, wordOrder WordOrder, in []byte) (out []float64) {
	var u64s []uint64

	u64s = bytesToUint64s(endianness, wordOrder, in)

	for _, u64 := range u64s {
		out = append(out, math.Float64frombits(u64))
	}

	return
}

func float64ToBytes(endianness Endianness, wordOrder WordOrder, in float64) (out []byte) {
	out = uint64ToBytes(endianness, wordOrder, math.Float64bits(in))

	return
}

func encodeBools(in []bool) (out []byte) {
	var byteCount uint
	var i uint

	byteCount = uint(len(in)) / 8
	if len(in)%8 != 0 {
		byteCount++
	}

	out = make([]byte, byteCount)
	for i = 0; i < uint(len(in)); i++ {
		if in[i] {
			out[i/8] |= 0x01 << (i % 8)
		}
	}

	return
}

func decodeBools(quantity uint16, in []byte) (out []bool) {
	var i uint
	for i = 0; i < uint(quantity); i++ {
		out = append(out, ((in[i/8]>>(i%8))&0x01) == 0x01)
	}

	return
}
