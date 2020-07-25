package webgin

import (
	"log"
	"net/http"
	"strings"
)

type handerfunc func(*Context)
type handers []handerfunc
type handermap map[string]handers
type methtree map[string]handermap
type Engine struct {
	RouterGroup
	tree methtree
}

func Defaul() *Engine {
	e := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "",
			root:     true,
		},
		tree: make(methtree, 7),
	}
	e.RouterGroup.engine = e
	return e
}

func (engine *Engine) Use(hf handerfunc) IRoutes {
	engine.RouterGroup.Use(hf)
	return engine
}
func (engine *Engine) addRoute(meth, path string, hs handers) {
	//if engine.tree==nil {
	//	engine.tree=make(methtree,7)
	//}
	log.Println(meth, path)
	myhm, ok := engine.tree[meth]
	if !ok {
		hm := make(handermap, 1)
		engine.tree[meth] = hm
		myhm = hm
	}
	if myhm[path] == nil {
		le := len(hs)
		hs := make(handers, le)
		myhm[path] = hs
	}
	for _, i2 := range hs {
		myhm[path] = append(engine.tree[meth][path], i2)
	}
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)
	meth := r.Method
	hm := engine.tree[meth]
	uri := r.RequestURI
	realuri := strings.Split(uri, "?")
	if len(realuri) < 1 {
		return
	}
	hs := hm[realuri[0]]
	if hs == nil {
		for i, hs := range hm {
			paths := strings.Split(i, ":")
			if len(paths) < 2 {
				continue
			}
			matchs := strings.Split(r.RequestURI, paths[0])
			if len(matchs) < 2 {
				continue
			}
			all := len(paths)
			value := strings.Split(matchs[1], "/")
			if len(matchs) < 2 || all-1 != len(value) {
				continue
			}
			param := make(map[string]string, all-1)
			for i := 1; i < all; i++ {
				key := strings.Split(paths[i], "/")
				param[key[0]] = value[i-1]
			}
			c.Lock()
			c.param = param
			c.Unlock()
			for _, i2 := range hs {
				if i2 != nil {
					i2(c)
				}
			}
			return
		}
		w.Write([]byte("404 not find"))
		return
	}
	for _, i2 := range hs {
		if i2 != nil {
			i2(c)
		}
	}
}
func (engine *Engine) Run(port string) {
	err := http.ListenAndServe(port, engine)
	if dealerr(err) {
		return
	}
}
