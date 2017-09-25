package main

type Serializable interface {
    serialize() []byte
}
