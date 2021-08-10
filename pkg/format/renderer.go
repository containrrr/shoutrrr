package format

// Renderer renders a supplied node tree to a string
type Renderer interface {
	// RenderTree renders a ContainerNode tree into a ansi-colored console string
	RenderTree(root *ContainerNode, scheme string) string
}

func testRenderTee(r Renderer, v interface{}) string {
	return r.RenderTree(getRootNode(v), "mock")
}
