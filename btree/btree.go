package btree

import (
	"mini_btree_kv_store/bnode"
	"mini_btree_kv_store/constants"
)

func init() {
	// max size of node with 1 kv-pair
	node1max := constants.HEADER + 8 + 2 + 4 + constants.BTREE_MAX_KEY_SIZE + constants.BTREE_MAX_VAL_SIZE
	if node1max > constants.BTREE_PAGE_SIZE {
		panic("nodelmax exceeds BTREE_PAGE_SIZE")
	}
}

type BTree struct {
	root uint64 // pointer (64-bit int referencing disk pages instead of in-memory nodes; we can't use in-memory pointers)
	// callbacks for managing on-disk pages
	get func(uint64) bnode.BNode // deference a pointer to get a disk page
	new func(bnode.BNode) uint64 // allocate a new disk page
	del func(uint64)             // deallocate a disk page
}
