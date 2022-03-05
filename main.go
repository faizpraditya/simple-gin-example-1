package main

import (
	"fmt"
	"log"
	"path/filepath"

	// http.StatusBadRequest
	"net/http"

	"github.com/gin-gonic/gin"
)

const UPLOAD_PATH = "/home/academy/golang-intermediate/simple_gin/static/"

func main() {
	// ReleaseMode (buat ngilangin gitdebug)
	gin.SetMode(gin.ReleaseMode)
	// By default
	// gin.SetMode(gin.DebugMode)
	// Secara default, gin.New() belum menggunakan middleware apapun
	// router := gin.Default()
	// By default sudah menggunakan gin.Logger() sama recovery
	router := gin.Default()
	// router := gin.New()
	// router.Use(gin.Logger())
	// Custom, nampilin path (alamat) dan method yang dipanggil
	// Bisa dipisah di fungsi sendiri
	// CustomLogger(router)
	// Handle panic, udah defer di dalam recovery
	// Bisa customRecovery
	// router.Use(gin.Recovery())
	// APP_PORT := os.Getenv("APP_PORT")

	router.GET("/:name/:age", hello)
	// router.POST("/", hello2)
	// router.GET("/", GetCustomerPaging)
	router.GET("/customer/:paging", GetCustomerPaging)
	// router.POST("/signin", login)
	// router.POST("/signin", gin.Logger(), login)

	// Router group
	// http://localhost:8181/client/signin
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "message": "PAGE_ROUTE_NOT_FOUND"})
	})
	client := router.Group("/client")
	// Explicit scope: https://medium.com/golangspec/scopes-in-go-a6042bb4298c
	// Inner block in main
	// Biar readabilitynya bisa bagus
	// Tescok ga bisa di akses di luar scope
	{
		// tescok := 123
		client.POST("/signin", login)
		client.POST("/upload", Upload)
		client.GET("")
	}
	// fmt.Println(tescok)

	// Default 8080
	router.Run(":8181")
	// router.Run(APP_PORT)
}

func CustomLogger(router *gin.Engine) {
	// Kalau mau dilempar ke handler tiap method, yang dilempar gin.Logger (gin.handlerfunction), buka router.Use
	router.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("%s %s\n", params.Path, params.Method)
	}))
}

type Credential struct {
	// kalau binding, masukin tagnya untuk form
	Username string `form:"uname" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func login(c *gin.Context) {
	// Contoh panic (selasa)
	// panic("panic in login")

	// username := c.PostForm("username")
	// password := c.PostForm("password")
	// c.FormFile("photo")

	var cred Credential
	err := c.ShouldBindJSON(&cred)
	if err != nil {
		// "net/http"
		// c.String(http.StatusBadRequest, err.Error())
		// Nerima sama, tapi bisa nerima apa aja karena dia interface
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest})
		return
	}
	// c.ShouldBind(&cred)
	// c.String(200, cred.Username+" "+cred.Password)
	// Pakai JSON
	// c.JSON(http.StatusOK, gin.H{"name ": cred.Username, "password ": cred.Password})

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "data": gin.H{"name ": cred.Username, "password ": cred.Password}})
}

func hello(c *gin.Context) {
	name := c.Param("name")
	fmt.Println(name)
	c.String(200, "hello "+name)
}

// func hello2(c *gin.Context) {
// 	c.String(200, "hello")
// }

// localhost:8181/?order=ASC&page=1&limit=1
// localhost:8181/customer/faiz?order=ASC&page=1&limit=1
func GetCustomerPaging(c *gin.Context) {
	paging := c.Param("paging")
	order := c.Query("order")
	page := c.Query("page")
	limit := c.Query("limit")
	c.String(200, paging+" "+order+" "+page+" "+limit)
}

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	filename := filepath.Base(file.Filename)
	log.Println(filename)
	if err := c.SaveUploadedFile(file, UPLOAD_PATH+filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
}
