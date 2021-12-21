package main

import (
	"fmt"
	// http.StatusBadRequest
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.GET("/:name/:age", hello)
	// router.POST("/", hello2)
	// router.GET("/", GetCustomerPaging)
	router.GET("/customer/:paging", GetCustomerPaging)
	router.POST("/signin", login)

	// Default 8080
	router.Run(":8181")
}

type Credential struct {
	// kalau binding, masukin tagnya untuk form
	Username string `form:"uname" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func login(c *gin.Context) {
	// username := c.PostForm("username")
	// password := c.PostForm("password")
	var cred Credential
	err := c.ShouldBindJSON(&cred)
	if err != nil {
		// "net/http"
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	// c.ShouldBind(&cred)
	c.String(200, cred.Username+" "+cred.Password)
}

func hello(c *gin.Context) {
	name := c.Param("name")
	fmt.Println(name)
	c.String(200, "hello "+name)
}

func hello2(c *gin.Context) {
	c.String(200, "hello")
}

// localhost:8181/?order=ASC&page=1&limit=1
// localhost:8181/customer/faiz?order=ASC&page=1&limit=1
func GetCustomerPaging(c *gin.Context) {
	paging := c.Param("paging")
	order := c.Query("order")
	page := c.Query("page")
	limit := c.Query("limit")
	c.String(200, paging+" "+order+" "+page+" "+limit)
}
