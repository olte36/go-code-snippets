package algorithms

func BreadthFirstTravese(root *TreeNode, fn func(*TreeNode) bool) {
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if curr != nil {
			if fn(curr) {
				return
			}
			queue = append(queue, curr.Leaves...)
		}
	}
}
