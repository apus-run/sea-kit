package ginx

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
)

func GinHandler(r *gin.Engine) *gin.Engine {
	getHelloFun := func(c *Context) {
		type getForm struct {
			ID uint64 `form:"id" binding:"gt=0"`
		}

		id, _ := strconv.ParseUint(c.Context.Param("id"), 10, 0)
		form := &getForm{ID: id}
		err := c.ShouldBindQuery(form)
		if err != nil {
			fmt.Println(err)
			c.JSONE(http.StatusBadRequest, "参数错误", err)
			return
		}
		fmt.Printf("%+v\n", form)
		c.Success()
	}

	createHelloFun := func(c *Context) {
		type postForm struct {
			Name  string `json:"name" form:"name" binding:"required"`
			Age   int    `json:"age" form:"age" binding:"gte=0,lte=150"`
			Email string `json:"email" form:"email" binding:"email,required"`
		}
		form := &postForm{}
		err := c.ShouldBind(form)

		if err != nil {
			fmt.Println(err)
			c.JSONE(http.StatusBadRequest, "参数错误", err)
			return
		}
		fmt.Printf("%+v\n", form)
		c.Success()
	}

	deleteHelloFun := func(c *Context) {
		type deleteForm struct {
			IDS []uint64 `form:"ids" binding:"required,min=1"`
		}
		form := &deleteForm{}
		err := c.ShouldBind(form)
		if err != nil {
			fmt.Println(err)
			c.JSONE(http.StatusBadRequest, "参数错误", err)
			return
		}
		fmt.Printf("%+v\n", form)
		c.Success()
	}

	updateHelloFun := func(c *Context) {
		type updateForm struct {
			ID    uint64 `json:"id" form:"id" binding:"required,gt=0"`
			Age   int    `json:"age" form:"age" binding:"gte=0,lte=150"`
			Email string `json:"email" form:"email" binding:"email,required"`
		}
		form := &updateForm{}
		err := c.ShouldBind(form)
		if err != nil {
			fmt.Println(err)
			c.JSONE(http.StatusBadRequest, "参数错误", err)
			return
		}
		fmt.Printf("%+v\n", form)
		c.Success()
	}

	r.GET("/hello", Handle(getHelloFun))
	r.GET("/hello/:id", Handle(getHelloFun))
	r.POST("/hello", Handle(createHelloFun))
	r.DELETE("/hello", Handle(deleteHelloFun))
	r.PUT("/hello", Handle(updateHelloFun))

	return r
}

func TestValidator(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	handler := GinHandler(r)

	// run server using httptest
	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/hello").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		HasValue("msg", "参数错误")
	e.POST("/hello").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		HasValue("msg", "参数错误")

	// "/user/0"
	e.GET("/hello/{id}", 0).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		HasValue("msg", "参数错误")

	// "/user/0"
	e.GET("/hello/{id}").
		WithPath("id", 0).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		HasValue("msg", "参数错误")

	// "/user/0?sort=asc"
	e.GET("/hello/{id}", 0).WithQuery("sort", "asc").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		HasValue("msg", "参数错误")

	e.DELETE("/hello").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		HasValue("msg", "参数错误")
	e.PUT("/hello").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		HasValue("msg", "参数错误")

	// t.Logf("obj: %#v", obj)
}
