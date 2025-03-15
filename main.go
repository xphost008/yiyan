package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const mysqlUsername = "root"
const mysqlPassword = "66543986"
const mysqlDatabase = "yiyandata"
const IP = "127.0.0.1"

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type Yiyan struct {
	Id          int    `json:"id"`
	Content     string `json:"content"`
	Submitter   string `json:"submitter"`
	Source      string `json:"source"`
	Author      string `json:"author"`
	Classifiers string `json:"classifiers"`
}
type LikeRecord struct {
	YiyanId int    `json:"yiyan_id"`
	UserId  string `json:"user_id"`
}

func (likeRecord LikeRecord) TableName() string {
	return "like_record"
}
func (yiyan Yiyan) TableName() string {
	return "yiyan"
}
func (user User) TableName() string {
	return "users"
}

type YiyanResult struct {
	Id          int    `json:"id"`
	Content     string `json:"content"`
	Submitter   string `json:"submitter"`
	Source      string `json:"source"`
	Author      string `json:"author"`
	Classifiers string `json:"classifiers"`
	Likes       int    `json:"likes"`
	IsLiked     bool   `json:"is_liked"`
}

func getYiyanMost() *YiyanResult {
	var yiyanId int
	gdb.Table("like_record").
		Select("yiyan_id").
		Group("yiyan_id").
		Order("COUNT(*) DESC").
		Limit(1).
		Pluck("yiyan_id", &yiyanId)
	var yiyanResult *YiyanResult
	gdb.Table("yiyan").Where("id = ?", yiyanId).First(&yiyanResult)
	return yiyanResult
}
func getYiyanLike(yiyanId int) int {
	var likes int64
	gdb.Table("like_record").Where("yiyan_id = ?", yiyanId).Count(&likes)
	return int(likes)
}
func isMyYiyanToLike(yiyanId int, userId string) bool {
	var submitter string
	gdb.Table("yiyan").Select("submitter").Where("id = ?", yiyanId).Pluck("submitter", &submitter)
	return submitter == userId
}
func getIsLike(yiyanId int, userId string) bool {
	if userId == "" {
		return false
	}
	var record []*LikeRecord
	gdb.Table("like_record").Find(&record)
	for _, recordItem := range record {
		if recordItem.YiyanId == yiyanId && recordItem.UserId == userId {
			return true
		}
	}
	return false
}
func getRandomOne(c *gin.Context) {
	stuId, errId := c.Cookie("id")
	if errId != nil {
		stuId = ""
	}
	var yiyanResult *YiyanResult
	gdb.Table("yiyan").Order("RAND()").First(&yiyanResult)
	yiyanResult.IsLiked = getIsLike(yiyanResult.Id, stuId)
	yiyanResult.Likes = getYiyanLike(yiyanResult.Id)
	c.JSON(200, yiyanResult)
}
func getMost(c *gin.Context) {
	stuId, errId := c.Cookie("id")
	if errId != nil {
		stuId = ""
	}
	yiyanGet := getYiyanMost()
	yiyanGet.IsLiked = getIsLike(yiyanGet.Id, stuId)
	yiyanGet.Likes = getYiyanLike(yiyanGet.Id)
	c.JSON(200, yiyanGet)
}
func getAll(c *gin.Context) {
	stuId, errId := c.Cookie("id")
	if errId != nil {
		stuId = ""
	}
	var yiyanAll []*YiyanResult
	gdb.Table("yiyan").Find(&yiyanAll)
	for _, yiyan := range yiyanAll {
		yiyan.IsLiked = getIsLike(yiyan.Id, stuId)
		yiyan.Likes = getYiyanLike(yiyan.Id)
	}
	c.JSON(200, yiyanAll)
}
func getMy(c *gin.Context) {
	stuId, errId := c.Cookie("id")
	_, errUsername := c.Cookie("username")
	_, errPassword := c.Cookie("password")
	if errId != nil || errUsername != nil || errPassword != nil {
		c.SetCookie("id", "", -1, "/", "", false, true)
		c.SetCookie("username", "", -1, "/", "", false, true)
		c.SetCookie("password", "", -1, "/", "", false, true)
		c.Redirect(302, "/login")
		return
	}
	var result []*YiyanResult
	gdb.Table("yiyan").Where("submitter = ?", stuId).Find(&result)
	for _, resultItem := range result {
		resultItem.IsLiked = getIsLike(resultItem.Id, stuId)
		resultItem.Likes = getYiyanLike(resultItem.Id)
	}
	c.JSON(200, result)
}
func like(c *gin.Context) {
	stuId, errId := c.Cookie("id")
	_, errUsername := c.Cookie("username")
	_, errPassword := c.Cookie("password")
	if errId != nil || errUsername != nil || errPassword != nil {
		c.SetCookie("id", "", -1, "/", "", false, true)
		c.SetCookie("username", "", -1, "/", "", false, true)
		c.SetCookie("password", "", -1, "/", "", false, true)
		c.String(200, "nologin")
		return
	}
	var lk map[string]interface{}
	if err := c.ShouldBindJSON(&lk); err != nil {
		c.String(200, "format error!")
		return
	}
	if lk["is_liked"].(bool) {
		gdb.Table("like_record").Unscoped().Where("user_id = ? AND yiyan_id = ?", stuId, int(lk["id"].(float64))).Delete(nil)
	} else {
		if isMyYiyanToLike(int(lk["id"].(float64)), stuId) {
			c.String(200, "my")
			return
		}
		gdb.Table("like_record").Create(&LikeRecord{YiyanId: int(lk["id"].(float64)), UserId: stuId})
	}
	c.String(200, "ok")
}
func submit(c *gin.Context) {
	stuId, errId := c.Cookie("id")
	_, errUsername := c.Cookie("username")
	_, errPassword := c.Cookie("password")
	if errId != nil || errUsername != nil || errPassword != nil {
		c.SetCookie("id", "", -1, "/", "", false, true)
		c.SetCookie("username", "", -1, "/", "", false, true)
		c.SetCookie("password", "", -1, "/", "", false, true)
		c.Redirect(302, "/login")
		return
	}
	content := c.PostForm("content")
	source := c.PostForm("source")
	author := c.PostForm("author")
	classifiers := c.PostForm("classifiers")
	fmt.Println(content, source, author, classifiers)
	if content == "" || source == "" || author == "" || classifiers == "" {
		c.Redirect(302, "/submit")
		return
	}
	gdb.Table("yiyan").Create(&Yiyan{
		Content:     content,
		Source:      source,
		Author:      author,
		Classifiers: classifiers,
		Submitter:   stuId,
	})
	if gdb.Error != nil {
		c.Redirect(302, "/submit")
		return
	}
	c.Redirect(302, "/list")
}
func login(c *gin.Context) {
	stuId := c.PostForm("student_id")
	password := c.PostForm("password")
	stuIdC, errId := c.Cookie("id")
	usernameC, errUsername := c.Cookie("username")
	_, errPassword := c.Cookie("password")
	if errId != nil || errUsername != nil || errPassword != nil {
		var user User
		gdb.Where("id = ? AND password = ?", stuId, password).First(&user)
		if user.Username == "" || user.Password == "" || user.Id == "" {
			c.Redirect(302, "/login")
			return
		}
		c.SetCookie("id", user.Id, 21600, "/", "", false, true)
		c.SetCookie("username", user.Username, 21600, "/", "", false, true)
		c.SetCookie("password", user.Password, 21600, "/", "", false, true)
		c.Redirect(302, "/")
	} else if stuId == stuIdC {
		gdb.Table("users").Where("id = ? AND username = ?", stuIdC, usernameC).Update("password", password)
		c.SetCookie("id", stuIdC, 21600, "/", "", false, true)
		c.SetCookie("username", usernameC, 21600, "/", "", false, true)
		c.SetCookie("password", password, 21600, "/", "", false, true)
		c.Redirect(302, "/")
	} else {
		var user User
		gdb.Where("id = ? AND password = ?", stuId, password).First(&user)
		if user.Username == "" || user.Password == "" || user.Id == "" {
			c.Redirect(302, "/login")
			return
		}
		c.SetCookie("id", user.Id, 21600, "/", "", false, true)
		c.SetCookie("username", user.Username, 21600, "/", "", false, true)
		c.SetCookie("password", user.Password, 21600, "/", "", false, true)
		c.Redirect(302, "/")
	}
}

var gdb *gorm.DB

// 操控数据库、操控网页Cookie、Session存储用户信息等。
func main() {
	g := gin.Default()
	fmt.Println("已创建Gin~")
	dsn := mysqlUsername + ":" + mysqlPassword + "@tcp(127.0.0.1:3306)/" + mysqlDatabase + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, errMySQL := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if errMySQL != nil {
		panic(errMySQL)
	}
	gdb = db
	fmt.Println("已连接MySQL~")
	g.Delims("[[", "]]")
	g.LoadHTMLGlob("templates/*")
	g.Static("/static", "./static")
	g.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	g.GET("/list", func(c *gin.Context) {
		c.HTML(200, "list.html", nil)
	})
	g.GET("/submit", func(c *gin.Context) {
		_, errStuId := c.Cookie("id")
		_, errUsername := c.Cookie("username")
		_, errPassword := c.Cookie("password")
		if errStuId != nil || errUsername != nil || errPassword != nil {
			c.SetCookie("id", "", -1, "/", "", false, true)
			c.SetCookie("username", "", -1, "/", "", false, true)
			c.SetCookie("password", "", -1, "/", "", false, true)
			c.Redirect(302, "/login")
		} else {
			c.HTML(200, "submit.html", nil)
		}
	})
	g.GET("/my", func(c *gin.Context) {
		_, errStuId := c.Cookie("id")
		_, errUsername := c.Cookie("username")
		_, errPassword := c.Cookie("password")
		if errStuId != nil || errUsername != nil || errPassword != nil {
			c.SetCookie("id", "", -1, "/", "", false, true)
			c.SetCookie("username", "", -1, "/", "", false, true)
			c.SetCookie("password", "", -1, "/", "", false, true)
			c.Redirect(302, "/login")
		} else {
			c.HTML(200, "my.html", nil)
		}
	})
	g.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	g.GET("/privacy", func(c *gin.Context) {
		c.HTML(200, "privacy.html", nil)
	})
	fmt.Println("已加载HTML静态资源~")
	yiyan := g.Group("/yiyan")
	{
		yiyan.GET("/get_random_one", getRandomOne)
		yiyan.GET("/get_all", getAll)
		yiyan.GET("/get_my", getMy)
		yiyan.GET("get_most", getMost)
		yiyan.POST("/like", like)
		yiyan.POST("/submit", submit)
		yiyan.POST("/login", login)
	}
	fmt.Println("已创建动态链接~")
	errGin := g.Run(":8080")
	if errGin != nil {
		panic(errGin)
	}
	fmt.Println("已结束web服务~")
}
