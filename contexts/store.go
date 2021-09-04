package contexts



// Context structure
type Context struct {
	link *Node
}

// New - create a new Context
func New() *Context {
	c := &Context{newNode()}
	return c
}

// Set - add the variable in context.
func (c *Context) Set(key interface{}, value interface{}) {
	c.link = c.link.add(key, value)
}

// Get - get variable.
func (c *Context) Get(key interface{}) interface{} {
	return c.link.search(key)
}

// Fix - fix the state to create a branch.
func (c *Context) Fix() *Context {
	nc := &Context{newNode()}
	nc.link = c.link
	return nc
}

// Node structure
type Node struct {
	parent *Node
	key    interface{}
	value  interface{}
}

// newNode - create a new Node-struct
func newNode() *Node {
	n := &Node{}
	return n
}

// search - find variable in nodes.
func (n *Node) search(key interface{}) interface{} {
	if n.key == key {
		return n.value
	}
	if n.parent != nil {
		return n.parent.search(key)
	}
	return nil
}

// add - add a node to a new variable.
func (n *Node) add(key interface{}, value interface{}) *Node {
	x := newNode()
	x.parent = n
	x.key = key
	x.value = value
	return x
}