package webgin

type IRoutes interface {
	Group(string, ...handerfunc) *RouterGroup
	Use(handerfunc) IRoutes
	GET(string, handerfunc) IRoutes
	POST(string, handerfunc) IRoutes
	PUT(string, handerfunc) IRoutes
	DELETE(string, handerfunc) IRoutes
	Static(string, string) IRoutes
	UseWebsocket(string, handerfunc) IRoutes
}
type RouterGroup struct {
	Handlers handers
	basePath string
	engine   *Engine
	root     bool
}

func (group *RouterGroup) Use(hf handerfunc) IRoutes {
	if group.Handlers == nil {
		group.Handlers = make(handers, 1)
	}
	group.Handlers = append(group.Handlers, hf)
	return group.returnObj()
}
func (group *RouterGroup) GET(path string, hf handerfunc) IRoutes {
	path = group.basePath + path
	return group.handle("GET", path, hf)
}
func (group *RouterGroup) POST(path string, hf handerfunc) IRoutes {
	path = group.basePath + path
	return group.handle("POST", path, hf)
}
func (group *RouterGroup) PUT(path string, hf handerfunc) IRoutes {
	path = group.basePath + path
	return group.handle("PUT", path, hf)
}
func (group *RouterGroup) DELETE(path string, hf handerfunc) IRoutes {
	path = group.basePath + path
	return group.handle("DELETE", path, hf)
}

func (group *RouterGroup) UseWebsocket(path string,hf handerfunc) IRoutes {
	Test()
	path=group.basePath+path
	group.handle("GET",path,hf)
	return group.returnObj()
}
func (group *RouterGroup) Group(path string, hf ...handerfunc) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combinHanders(hf),
		basePath: group.basePath + path,
		engine:   group.engine,
	}
}
func (group *RouterGroup) Static(path, filePath string) IRoutes {
	uri := path + "/:filename"
	hf := func(c *Context) {
		n := c.param["filename"]
		f := filePath + "\\" + n
		c.File(f)
	}
	group.GET(uri, hf)
	return group.returnObj()
}

func (group *RouterGroup) handle(meth, path string, hf ...handerfunc) IRoutes {
	handers := group.combinHanders(hf)
	group.engine.addRoute(meth, path, handers)
	return group.returnObj()
}
func (group *RouterGroup) returnObj() IRoutes {
	if group.root {
		return group.engine
	}
	return group
}
func (group *RouterGroup) combinHanders(hf handers) handers {
	l := len(group.Handlers) + len(hf)
	allhander := make(handers, l)
	copy(allhander, group.Handlers)
	copy(allhander[len(group.Handlers):], hf)
	return allhander
}
