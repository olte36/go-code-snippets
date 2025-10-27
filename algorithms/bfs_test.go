package algorithms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBreadthFirstTravese(t *testing.T) {
	root := &TreeNode{
		Val: "A",
		Leaves: []*TreeNode{
			// left
			{
				Val: "B",
				Leaves: []*TreeNode{
					// left
					{
						Val: "E",
					},
					// middle
					nil,
					// right
					{
						Val: "F",
					},
				},
			},
			// middle
			{
				Val: "C",
			},
			// right
			{
				Val: "D",
				Leaves: []*TreeNode{
					// left
					nil,
					// middle
					nil,
					// right
					{
						Val: "G",
					},
				},
			},
		},
	}

	var got []string
	BreadthFirstTravese(root, func(n *TreeNode) bool {
		got = append(got, n.Val)
		return false
	})

	want := []string{"A", "B", "C", "D", "E", "F", "G"}
	assert.EqualValues(t, want, got)
}
