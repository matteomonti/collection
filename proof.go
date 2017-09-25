package main

type NodeStep struct {
    label Buffer
    value uint64
    left Buffer
    right Buffer
}

func (step *NodeStep) valid() bool {
    return BufferEqual(step.label, Hash(BufferFromBool(false), BufferFromUint64(step.value), step.left, step.right))
}

func (step *NodeStep) update() {
    step.label = Hash(BufferFromBool(false), BufferFromUint64(step.value), step.left, step.right)
}

func MakeNodeStep(node *Node) (nodestep NodeStep) {
    nodestep.label = node.label
    nodestep.value = node.value
    nodestep.left = node.left.label
    nodestep.right = node.right.label

    return
}

type LeafStep struct {
    label Buffer
    value uint64
    key Buffer
}

func (step *LeafStep) valid() bool {
    return BufferEqual(step.label, Hash(BufferFromBool(true), BufferFromUint64(step.value), step.key))
}

func MakeLeafStep(node *Node) (leafstep LeafStep) {
    leafstep.label = node.label
    leafstep.value = node.value
    leafstep.key = node.key

    return
}

type Proof struct {
    prooftype struct {
        add bool
        get bool
        update bool
        delete bool
        cumulative bool
    }

    key Buffer
    value uint64
    cumulative uint64
    success bool

    nodes []NodeStep
    leaf LeafStep
}

func MakeAddProof() (proof Proof) {
    proof.prooftype.add = true
    return
}

func MakeGetProof() (proof Proof) {
    proof.prooftype.get = true
    return
}

func MakeUpdateProof() (proof Proof) {
    proof.prooftype.update = true
    return
}

func MakeDeleteProof() (proof Proof) {
    proof.prooftype.delete = true
    return
}

func MakeCumulativeProof(value uint64) (proof Proof) {
    proof.prooftype.cumulative = true
    proof.cumulative = value
    return
}

func (proof *Proof) addnode(node *Node) {
    proof.nodes = append(proof.nodes, MakeNodeStep(node))
}

func (proof *Proof) setleaf(node *Node) {
    proof.leaf = MakeLeafStep(node)
}
