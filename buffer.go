package main

import "encoding/binary"

type Buffer []uint8

func (buffer Buffer) serialize() []byte {
    return buffer
}

func (buffer Buffer) bit(index uint64) bool {
    byteidx := index / 8
    bitidx := 7 - (index % 8)

    return ((buffer[byteidx] & (uint8(1) << bitidx)) != 0)
}

func (buffer Buffer) bits(count int) (bits []bool) {
    for i := 0; i < count; i++ {
        bits = append(bits, buffer.bit(uint64(i)))
    }

    return
}

func BufferFromBool(value bool) Buffer {
    if(value) {
        return Buffer{1}
    } else {
        return Buffer{0}
    }
}

func BufferFromUint32(value uint32) Buffer {
    bytes := make([]byte, 4)
    binary.BigEndian.PutUint32(bytes, value)

    return bytes
}

func BufferFromUint64(value uint64) Buffer {
    bytes := make([]byte, 8)
    binary.BigEndian.PutUint64(bytes, value)

    return bytes
}

func BufferFromString(value string) Buffer {
    return []byte(value)
}

func BufferEqual(lho, rho Buffer) bool {
    if(lho == nil && rho == nil) {
        return true
    }

    if(lho == nil || rho == nil) {
        return false
    }

    if(len(lho) != len(rho)) {
        return false
    }

    for i := range(lho) {
        if(lho[i] != rho[i]) {
            return false
        }
    }

    return true
}
