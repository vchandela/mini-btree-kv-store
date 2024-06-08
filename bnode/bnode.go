package bnode

import (
	"encoding/binary"
	"mini_btree_kv_store/constants"
)

type BNode struct {
	data []byte // can be dumped to the disk
}

// header
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node.data[0:2])
}

func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}

func (node BNode) setHeader(btype, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:2], btype)
	binary.LittleEndian.PutUint16(node.data[2:4], nkeys)
}

// pointers
func (node BNode) getPtr(idx uint16) uint64 {
	if idx > node.nkeys() {
		panic("idx > nkeys")
	}
	pos := constants.HEADER + idx*8
	return binary.LittleEndian.Uint64(node.data[pos : pos+8])
}

func (node BNode) setPtr(idx uint16, val uint64) {
	if idx > node.nkeys() {
		panic("idx > nkeys")
	}
	pos := constants.HEADER + idx*8
	binary.LittleEndian.PutUint64(node.data[pos:pos+8], val)
}

// offset list
func offsetPos(node BNode, idx uint16) uint16 {
	if idx < 1 || idx > node.nkeys() {
		panic("idx > nkeys")
	}
	// each offset is 2 bytes
	return constants.HEADER + 8*node.nkeys() + 2*(idx-1)
}

func (node BNode) getOffset(idx uint16) uint16 {
	// offset of the 1st kv pair is always 0, so it's not stored in the offset list
	if idx == 0 {
		return 0
	}
	pos := offsetPos(node, idx)
	return binary.LittleEndian.Uint16(node.data[pos : pos+2])
}

func (node BNode) setOffset(idx uint16, val uint16) {
	pos := offsetPos(node, idx)
	binary.LittleEndian.PutUint16(node.data[pos:pos+2], val)
}

// offset list is used to locate the nth KV pair quickly
// key-values
func (node BNode) kvPos(idx uint16) uint16 {
	if idx > node.nkeys() {
		panic("idx > nkeys")
	}
	return constants.HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

func (node BNode) getKey(idx uint16) []byte {
	if idx > node.nkeys() {
		panic("idx > nkeys")
	}
	pos := node.kvPos(idx)
	keyLen := binary.LittleEndian.Uint16(node.data[pos : pos+2])
	return node.data[pos+4 : pos+4+keyLen]
}

func (node BNode) getVal(idx uint16) []byte {
	if idx > node.nkeys() {
		panic("idx > nkeys")
	}
	pos := node.kvPos(idx)
	keyLen := binary.LittleEndian.Uint16(node.data[pos : pos+2])
	valLen := binary.LittleEndian.Uint16(node.data[pos+2 : pos+4])
	return node.data[pos+4+keyLen : pos+4+keyLen+valLen]
}

// last offset in offset list is used to determine size of node in bytes
func (node BNode) nodeSizeInBytes() uint16 {
	return node.kvPos(node.nkeys())
}