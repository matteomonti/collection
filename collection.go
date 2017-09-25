package main

type Collection struct {
    root Node
}

func MakeCollection() Collection {
    return Collection{MakeNode()}
}

func (collection *Collection) add(key Buffer, value uint64) Proof {
    proof := MakeAddProof()
    proof.key = key
    proof.value = value

    path := Hash(key)

    depth := uint64(0)
    leafset := false
    cursor := &(collection.root)

    var breadcrumbs [256]*Node
    breadcrumbs[depth] = cursor

    for {
        if(!leafset) {
            proof.addnode(cursor)
        }

        bit := path.bit(depth)

        if(bit) {
            cursor = cursor.right
        } else {
            cursor = cursor.left
        }

        depth++
        breadcrumbs[depth] = cursor

        if(cursor.placeholder) {
            if(!leafset) {
                proof.setleaf(cursor)
            }

            (*cursor) = MakeLeaf(key, value)
            break
        } else if(cursor.leaf) {
            if(!leafset) {
                leafset = true
                proof.setleaf(cursor)
            }

            if(BufferEqual(key, cursor.key)) {
                return proof
            }

            leaf := *cursor
            leafpath := Hash(leaf.key)
            leafbit := leafpath.bit(depth)

            (*cursor) = MakeNode()

            if(leafbit) {
                *(cursor.right) = leaf
            } else {
                *(cursor.left) = leaf
            }
        }
    }

    for {
        depth--
        cursor = breadcrumbs[depth]

        cursor.update()

        if(depth == 0) {
            break
        }
    }

    proof.success = true
    return proof
}

func (collection *Collection) get(key Buffer) Proof {
    proof := MakeGetProof()
    proof.key = key

    path := Hash(key)

    depth := uint64(0)
    cursor := &(collection.root)

    for {
        proof.addnode(cursor)

        bit := path.bit(depth)

        if(bit) {
            cursor = cursor.right
        } else {
            cursor = cursor.left
        }

        depth++

        if(cursor.leaf) {
            proof.setleaf(cursor)

            if(BufferEqual(cursor.key, key)) {
                proof.success = true
                proof.value = cursor.value
            }

            break
        }
    }

    return proof
}

func (collection *Collection) update(key Buffer, value uint64) Proof {
    proof := MakeUpdateProof()

    proof.key = key
    proof.value = value

    path := Hash(key)

    depth := uint64(0)
    cursor := &(collection.root)

    var breadcrumbs [256]*Node
    breadcrumbs[depth] = cursor

    for {
        proof.addnode(cursor)

        bit := path.bit(depth)

        if(bit) {
            cursor = cursor.right
        } else {
            cursor = cursor.left
        }

        depth++
        breadcrumbs[depth] = cursor

        if(cursor.leaf) {
            proof.setleaf(cursor)

            if(!BufferEqual(cursor.key, key)) {
                return proof
            }

            break
        }
    }

    proof.success = true
    (*cursor) = MakeLeaf(key, value)

    for {
        depth--
        cursor = breadcrumbs[depth]

        cursor.update()

        if(depth == 0) {
            break
        }
    }

    return proof
}

func (collection *Collection) delete(key Buffer) Proof {
    proof := MakeDeleteProof()

    proof.key = key

    path := Hash(key)

    depth := uint64(0)
    cursor := &(collection.root)

    var breadcrumbs [256]*Node
    breadcrumbs[depth] = cursor

    for {
        proof.addnode(cursor)

        bit := path.bit(depth)

        if(bit) {
            cursor = cursor.right
        } else {
            cursor = cursor.left
        }

        depth++
        breadcrumbs[depth] = cursor

        if(cursor.leaf) {
            proof.setleaf(cursor)

            if(!BufferEqual(cursor.key, key)) {
                return proof
            }

            break
        }
    }

    proof.success = true
    proof.value = cursor.value

    (*cursor) = MakePlaceholder()

    for {
        depth--
        cursor = breadcrumbs[depth]

        cursor.update()

        if(depth == 0) {
            break
        }

        if(cursor.left.placeholder && cursor.right.placeholder) {
            *(cursor) = MakePlaceholder()
        }
    }

    return proof
}

func (collection *Collection) cumulative(value uint64) Proof {
    proof := MakeCumulativeProof(value)

    cursor := &(collection.root)

    if(value >= collection.root.value) {
        proof.addnode(cursor)
        return proof
    }

    for {
        proof.addnode(cursor)

        if(value >= cursor.left.value) {
            value -= cursor.left.value
            cursor = cursor.right
        } else {
            cursor = cursor.left
        }

        if(cursor.leaf) {
            proof.setleaf(cursor)

            proof.success = true
            proof.key = cursor.key
            proof.value = cursor.value

            break
        }
    }

    return proof
}
