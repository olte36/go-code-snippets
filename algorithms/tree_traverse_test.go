package algorithms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeTraverse(t *testing.T) {

	testCases := []struct {
		name       string
		traverseFn func(*TreeNode, func(*TreeNode) bool)
		want       []string
	}{
		{
			name:       "BreadthFirstTravese",
			traverseFn: BreadthFirstTravese,
			want:       []string{"A", "B", "C", "D", "E", "F", "G"},
		},
		{
			name:       "DepthFirstSearchPreOrderRecursive",
			traverseFn: DepthFirstTraversePreOrderRecursive,
			want:       []string{"A", "B", "E", "F", "C", "D", "G"},
		},
		{
			name:       "DepthFirstSearchPostOrderRecursive",
			traverseFn: DepthFirstTraversePostOrderRecursive,
			want:       []string{"E", "F", "B", "C", "G", "D", "A"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			root := createTestTree()
			var got []string
			tc.traverseFn(root, func(n *TreeNode) bool {
				got = append(got, n.Val)
				return false
			})

			assert.EqualValues(t, tc.want, got)
		})
	}
}

func createTestTree() *TreeNode {
	//	      A
	//	   /  |  \
	//	  B   C   D
	//	/ | \     \
	// E  .  F     G
	return &TreeNode{
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
}
