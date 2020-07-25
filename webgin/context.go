package webgin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"strings"
	"sync"
)

var message=make( chan []byte)
var conns map[string]*myconn
var star=make(chan int)
var one sync.Once

type myconn struct {
	towho string
	conn  *websocket.Conn
}
type Context struct {
	w    http.ResponseWriter
	r    *http.Request
	keys map[string]interface{}
	sync.RWMutex
	formparam  map[string]string
	queryparam map[string]string
	myjson     []byte
	param      map[string]string
}

type H map[string]interface{}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	log.Println("接收到请求", r.Method, r.RequestURI)
	c := &Context{
		w:          w,
		r:          r,
		keys:       make(map[string]interface{}),
		formparam:  nil,
		queryparam: nil,
		myjson:     nil,
	}
	return c
}
func (c *Context) Bind(goal interface{}) (string, error) {
	info := c.r.Header.Get("Content-Type")
	if strings.HasPrefix(info, "json") {
		c.BindJSON(goal)
		return "", nil
	}
	if strings.HasPrefix(info, "multipart") || strings.HasPrefix(info, "x-www-form-urlencoded") {
		return c.PostForm(goal.(string)), nil
	}
	err := errors.New("暂不支持该格式")
	return "", err
}

func (c *Context) PostForm(key string) string {
	if c.formparam == nil {
		c.formparam = form(c.r)
	}
	return c.formparam[key]
}
func form(r *http.Request) (m map[string]string) {
	m = make(map[string]string)

	err := r.ParseForm()
	if dealerr(err) {
		return
	}
	if len(r.PostForm) < 1 {
		err = r.ParseMultipartForm(1024)
		if dealerr(err) {
			return
		}
		for i, i2 := range r.MultipartForm.Value {
			m[i] = i2[0]
		}
		return
	}
	for i, i2 := range r.PostForm {
		m[i] = i2[0]
	}
	return
}

func (c *Context) Query(key string) string {
	if c.queryparam == nil {
		c.queryparam = parseQuery(c.r)
	}
	return c.queryparam[key]
}
func parseQuery(r *http.Request) (param map[string]string) {
	param = make(map[string]string)
	uri := r.RequestURI
	uris := strings.Split(uri, "?")
	if len(uris) == 1 {
		return
	}
	thing := strings.Split(uris[len(uris)-1], "&")
	for _, i2 := range thing {
		num := strings.Split(i2, "=")
		if len(num) != 2 {
			return
		}
		param[num[0]] = num[1]
	}
	return
}

func (c *Context) BindJSON(obj interface{}) {
	if c.myjson == nil {
		c.myjson = getjson(c.r)
	}
	err := json.Unmarshal(c.myjson, obj)
	if dealerr(err) {
		return
	}
}
func getjson(r *http.Request) []byte {
	myjson, err := ioutil.ReadAll(r.Body)
	if dealerr(err) {
		return nil
	}
	return myjson
}

func (c *Context) Param(key string) (value string) {
	c.RLock()
	defer c.RUnlock()
	value, ok := c.param[key]
	if !ok {
		return ""
	}
	return
}

func (c *Context) Set(key string, value interface{}) {
	c.Lock()
	c.keys[key] = value
	c.Unlock()
}
func (c *Context) Get(key string) (value interface{}) {
	c.RLock()
	defer c.RUnlock()
	value, ok := c.keys[key]
	if !ok {
		return nil
	}
	return
}

func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.w, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 0,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}
func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.r.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (c *Context) JSON(obj interface{}) {
	a, err := json.Marshal(obj)
	if dealerr(err) {
		return
	}
	c.w.Write(a)
}
func (c *Context) File(path string) {
	http.ServeFile(c.w, c.r, path)
}

func dealerr(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

func Updatewebsocket(c *Context, id string, towho string) {
	if towho == "" {
		towho = "all"
	}
	if conns == nil {
		conns = make(map[string]*myconn)
	}
	u := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := u.Upgrade(c.w, c.r, nil)
	if dealerr(err) {
		return
	}
	co := &myconn{
		towho: towho,
		conn:  conn,
	}
	fmt.Println(co)
	conns[id]=co
	one.Do(func() {
		close(star)
	})
}

func Getmess() {
	i := 0
	<-star
	for {
		fmt.Println(i)
		i++
		for i, i2 := range conns {
			//if i2 == nil {
			//	continue
			//}
			_, mess, err := i2.conn.ReadMessage()
			if dealerr(err) {
				return
			}
			x := map[string]string{
				"id":    i,
				"towho": i2.towho,
				"mess":  string(mess),
			}
			js, err := json.Marshal(x)
			if dealerr(err) {
				return
			}
			message <- js
		}
	}
}
func Sendmess() {
	for {
		mess := <-message
		ma := make(map[string]string)
		_ = json.Unmarshal(mess, &ma)
		for i, i2 := range conns {
			if ma["towho"] == "all" || ma["towho"] == i {
				err := i2.conn.WriteMessage(websocket.TextMessage, mess)
				if dealerr(err) {
					delete(conns,i)
				}
			}
		}
	}
}
func Test() {
	go Getmess()
	go Sendmess()
}
