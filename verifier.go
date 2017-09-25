package main

type Verifier struct {
    root Buffer
}

func MakeVerifier(root Buffer) (verifier Verifier) {
    verifier.root = root
    return
}

func (verifier *Verifier) verify(proof Proof) bool {

    if(len(proof.nodes) == 0) {
        return false
    }

    if(proof.prooftype.add) {
        path := Hash(proof.key)
        cursor := verifier.root

        depth := 0
        for ; depth < len(proof.nodes); depth++ {
            node := proof.nodes[depth]

            if(!BufferEqual(node.label, cursor)) {
                return false
            }

            if(!(node.valid())) {
                return false
            }

            if(path.bit(uint64(depth))) {
                cursor = node.right
            } else {
                cursor = node.left
            }
        }

        if(!BufferEqual(proof.leaf.label, cursor) || !(proof.leaf.valid())) {
            return false
        }

        if(BufferEqual(proof.leaf.key, proof.key)) {
            if(proof.success == false) {
                return true
            } else {
                return false
            }
        }

        if(proof.success == false) {
            return false
        }

        if(len(proof.leaf.key) == 0 && proof.leaf.value == 0) {
            leaf := MakeLeaf(proof.key, proof.value)
            proof.setleaf(&leaf)
        } else {
            leafpath := Hash(proof.leaf.key)

            for {
                node := MakeNode()
                node.value = proof.leaf.value
                proof.addnode(&node)

                if(leafpath.bit(uint64(depth)) != path.bit(uint64(depth))) {
                    break
                }

                depth++
            }

            if(leafpath.bit(uint64(depth))) {
                proof.nodes[depth].right = proof.leaf.label
            } else {
                proof.nodes[depth].left = proof.leaf.label
            }

            depth++

            leaf := MakeLeaf(proof.key, proof.value)
            proof.setleaf(&leaf)
        }

        for {
            depth--;
            proof.nodes[depth].value += proof.value

            var child Buffer

            if(depth == len(proof.nodes) - 1) {
                child = proof.leaf.label
            } else {
                child = proof.nodes[depth + 1].label
            }

            if(path.bit(uint64(depth))) {
                proof.nodes[depth].right = child
            } else {
                proof.nodes[depth].left = child
            }

            proof.nodes[depth].update()

            if(depth == 0) {
                break
            }
        }

        verifier.root = proof.nodes[0].label

        return true
    } else if(proof.prooftype.get) {
        path := Hash(proof.key)
        cursor := verifier.root

        depth := 0
        for ; depth < len(proof.nodes); depth++ {
            node := proof.nodes[depth]

            if(!BufferEqual(node.label, cursor)) {
                return false
            }

            if(!(node.valid())) {
                return false
            }

            if(path.bit(uint64(depth))) {
                cursor = node.right
            } else {
                cursor = node.left
            }
        }

        if(!BufferEqual(proof.leaf.label, cursor) || !(proof.leaf.valid())) {
            return false
        }

        if(BufferEqual(proof.leaf.key, proof.key)) {
            if(proof.success == true && proof.value == proof.leaf.value) {
                return true
            } else {
                return false
            }
        } else {
            if(proof.success == false) {
                return true
            } else {
                return false
            }
        }
    } else if(proof.prooftype.update) {
        path := Hash(proof.key)
        cursor := verifier.root

        depth := 0
        for ; depth < len(proof.nodes); depth++ {
            node := proof.nodes[depth]

            if(!BufferEqual(node.label, cursor)) {
                return false
            }

            if(!(node.valid())) {
                return false
            }

            if(path.bit(uint64(depth))) {
                cursor = node.right
            } else {
                cursor = node.left
            }
        }

        if(!BufferEqual(proof.leaf.label, cursor) || !(proof.leaf.valid())) {
            return false
        }

        if(BufferEqual(proof.leaf.key, proof.key)) {
            if(proof.success == false) {
                return false
            }

            oldvalue := proof.leaf.value

            leaf := MakeLeaf(proof.key, proof.value)
            proof.setleaf(&leaf)

            for {
                depth--;

                proof.nodes[depth].value += proof.value
                proof.nodes[depth].value -= oldvalue

                var child Buffer

                if(depth == len(proof.nodes) - 1) {
                    child = proof.leaf.label
                } else {
                    child = proof.nodes[depth + 1].label
                }

                if(path.bit(uint64(depth))) {
                    proof.nodes[depth].right = child
                } else {
                    proof.nodes[depth].left = child
                }

                proof.nodes[depth].update()

                if(depth == 0) {
                    break
                }
            }

            verifier.root = proof.nodes[0].label
            return true
        } else {
            if(proof.success == false) {
                return true
            } else {
                return false
            }
        }
    } else if(proof.prooftype.delete) {
        path := Hash(proof.key)
        cursor := verifier.root

        depth := 0
        for ; depth < len(proof.nodes); depth++ {
            node := proof.nodes[depth]

            if(!BufferEqual(node.label, cursor)) {
                return false
            }

            if(!(node.valid())) {
                return false
            }

            if(path.bit(uint64(depth))) {
                cursor = node.right
            } else {
                cursor = node.left
            }
        }

        if(!BufferEqual(proof.leaf.label, cursor) || !(proof.leaf.valid())) {
            return false
        }

        if(BufferEqual(proof.leaf.key, proof.key)) {
            if(proof.success == false || proof.value != proof.leaf.value) {
                return false
            }

            placeholder := MakePlaceholder()

            for {
                depth--;

                proof.nodes[depth].value -= proof.value

                var child Buffer

                if(depth == len(proof.nodes) - 1) {
                    child = placeholder.label
                } else {
                    child = proof.nodes[depth + 1].label
                }

                if(path.bit(uint64(depth))) {
                    proof.nodes[depth].right = child
                } else {
                    proof.nodes[depth].left = child
                }

                if(depth > 0 && BufferEqual(proof.nodes[depth].left, placeholder.label) && BufferEqual(proof.nodes[depth].right, placeholder.label)) {
                    proof.nodes[depth].label = placeholder.label
                } else {
                    proof.nodes[depth].update()
                }

                if(depth == 0) {
                    break
                }
            }

            verifier.root = proof.nodes[0].label
            return true
        } else {
            if(proof.success == false) {
                return true
            } else {
                return false
            }
        }
    } else if(proof.prooftype.cumulative) {
        cursor := verifier.root

        if(!BufferEqual(proof.nodes[0].label, cursor)) {
            return false
        }

        if(!(proof.nodes[0].valid())) {
            return false
        }

        if(proof.nodes[0].value <= proof.cumulative) {
            return !(proof.success)
        }

        depth := 0

        for ; depth < len(proof.nodes); depth++ {
            node := proof.nodes[depth]

            if(depth < len(proof.nodes) - 1) {
                var left bool

                if(BufferEqual(node.left, proof.nodes[depth + 1].label)) {
                    left = true
                } else if(!BufferEqual(node.right, proof.nodes[depth + 1].label)) {
                    return false
                }

                if(!(proof.nodes[depth + 1].valid())) {
                    return false
                }

                if(left) {
                    if(proof.cumulative >= proof.nodes[depth + 1].value) {
                        return false
                    }
                } else {
                    if(proof.cumulative < (node.value - proof.nodes[depth + 1].value)) {
                        return false
                    }

                    proof.cumulative -= (node.value - proof.nodes[depth + 1].value)
                }
            } else {
                var left bool

                if(BufferEqual(node.left, proof.leaf.label)) {
                    left = true
                } else if(!BufferEqual(node.right, proof.leaf.label)) {
                    return false
                }

                if(!(proof.leaf.valid())) {
                    return false
                }

                if(left) {
                    if(proof.cumulative >= proof.leaf.value) {
                        return false
                    }
                } else {
                    if(proof.cumulative < (node.value - proof.leaf.value)) {
                        return false
                    }
                }

                return true
            }
        }

        return true
    }

    return false
}
