package format

// TreeRenderer renders a ContainerNode tree into a string
type TreeRenderer interface {
	RenderTree(root *ContainerNode, scheme string) string
}
