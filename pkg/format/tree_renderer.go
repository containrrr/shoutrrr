package format

type TreeRenderer interface {
	RenderTree(root *ContainerNode, scheme string) string
}
