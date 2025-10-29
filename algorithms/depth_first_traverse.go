package algorithms

func DepthFirstTraversePreOrderRecursive(root *TreeNode, fn func(*TreeNode) bool) {
	if root == nil {
		return
	}
	if fn(root) {
		return
	}
	for _, leaf := range root.Leaves {
		DepthFirstTraversePreOrderRecursive(leaf, fn)
	}
}

func DepthFirstTraversePostOrderRecursive(root *TreeNode, fn func(*TreeNode) bool) {
	if root == nil {
		return
	}
	for _, leaf := range root.Leaves {
		DepthFirstTraversePostOrderRecursive(leaf, fn)
	}
	fn(root)
}
