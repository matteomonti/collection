package main

type Node struct {
    label Buffer
    value uint64

    leaf bool
    placeholder bool

    left *Node
    right *Node

    key Buffer
}

func (node *Node) update() {
    node.value = node.left.value + node.right.value
    node.label = Hash(BufferFromBool(node.leaf), BufferFromUint64(node.value), node.left.label, node.right.label)
}

func (node *Node) validate(path... bool) bool {
    if(node.leaf) {
        if(node.placeholder) {
            if(node.value != 0 || len(node.key) != 0) {
                return false
            }
        } else {
            leafpath := Hash(node.key)

            for depth, step := range path {
                if(step != leafpath.bit(uint64(depth))) {
                    return false
                }
            }
        }

        label := Hash(BufferFromBool(node.leaf), BufferFromUint64(node.value), node.key)

        if(!BufferEqual(label, node.label)) {
            return false
        } else {
            return true
        }
    } else {
        label := Hash(BufferFromBool(node.leaf), BufferFromUint64(node.value), node.left.label, node.right.label)

        if(!BufferEqual(label, node.label)) {
            return false
        } else {
            leftpath := make([]bool, len(path))
            rightpath := make([]bool, len(path))

            copy(leftpath, path)
            copy(rightpath, path)

            leftpath = append(leftpath, false)
            rightpath = append(rightpath, true)

            return node.left.validate(leftpath...) && node.right.validate(rightpath...)
        }
    }
}

func MakeLeaf(key Buffer, value uint64) (node Node) {
    node.leaf = true

    node.key = key
    node.value = value

    node.label = Hash(BufferFromBool(node.leaf), BufferFromUint64(node.value), node.key)
    return
}

func MakePlaceholder() Node {
    node := MakeLeaf(Buffer{}, 0)
    node.placeholder = true

    return node
}

func MakeNode() (node Node) {
    node.left = new(Node)
    *(node.left) = MakePlaceholder()

    node.right = new(Node)
    *(node.right) = MakePlaceholder()

    node.label = Hash(BufferFromBool(node.leaf), BufferFromUint64(node.value), node.left.label, node.right.label)
    return
}
