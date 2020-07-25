package main


import "webgin/webgin"

func main()  {
	h:=webgin.Defaul()
	r:=h.Group("/user")
	r.Use(setwho)
	r.GET("/123/:id/:name",getwho )
	r.POST("/123",postwho)
	r.Static("/pic","C:\\Users\\浩瀚宇\\Documents\\Tencent Files\\735268835\\FileRecv\\MobileFile\\Image")
	h.UseWebsocket("/ws/:id/:towho", func(c *webgin.Context) {
		id:=c.Param("id")
		towho:=c.Param("towho")
		webgin.Updatewebsocket(c,id,towho)
	})
	h.Run("localhost:8080")
}
func setwho(ctx *webgin.Context)  {
	ctx.Set("name","shazi")
}
func getwho(context *webgin.Context) {
	id:=context.Param("id")
	name:=context.Param("name")
	context.JSON(webgin.H{
		"id":id,
		"name":name,
	})
}
func postwho(ctx *webgin.Context)  {
	id:=ctx.PostForm("id")
	name:=ctx.Get("name")
	ctx.JSON(webgin.H{
		"id":id,
		"name":name,
	})
}


