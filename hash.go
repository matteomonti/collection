package main

import "crypto/sha256"

func Hash(item Serializable, items... Serializable) Buffer {
    digest := sha256.New()

    digest.Write(BufferFromUint64(uint64(len(item.serialize()))).serialize())
    digest.Write(item.serialize())

    for _, item := range items {
        digest.Write(BufferFromUint64(uint64(len(item.serialize()))).serialize())
        digest.Write(item.serialize());
    }

    return digest.Sum(nil)
}
