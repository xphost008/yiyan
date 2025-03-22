package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
)

const (
	mysqlUsername = "root"
	mysqlPassword = "66543986"
	mysqlDatabase = "yiyandata"
	IP            = "127.0.0.1"
	Port          = "3306"
)

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

func isValidUserId(userId string) bool {
	if userId == "" {
		return false
	}
	var id = ""
	gdb.Table("users").Where("id = ?", userId).Pluck("id", &id)
	return id != ""
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
func getYiyanMost() *YiyanResult {
	var yiyanId int
	gdb.Table("like_record").
		Select("yiyan_id").
		Group("yiyan_id").
		Order("COUNT(*) DESC").
		Limit(1).
		Pluck("yiyan_id", &yiyanId)
	if yiyanId == 0 {
		yiyanId = 1
	}
	var yiyanResult *YiyanResult
	gdb.Table("yiyan").Where("id = ?", yiyanId).First(&yiyanResult)
	return yiyanResult
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
	if errId != nil || errUsername != nil || errPassword != nil || !isValidUserId(stuId) {
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
	if errId != nil || errUsername != nil || errPassword != nil || !isValidUserId(stuId) {
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
	yId := int(lk["id"].(float64))
	if yId == 0 {
		c.String(200, "id is 0")
		return
	}
	if lk["is_liked"].(bool) {
		gdb.Table("like_record").Unscoped().Where("user_id = ? AND yiyan_id = ?", stuId, yId).Delete(nil)
	} else {
		if isMyYiyanToLike(yId, stuId) {
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
	if errId != nil || errUsername != nil || errPassword != nil || !isValidUserId(stuId) {
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
	if errId != nil || errUsername != nil || errPassword != nil || !isValidUserId(stuId) {
		var user User
		gdb.Where("id = ? AND password = ?", stuId, password).First(&user)
		if user.Username == "" || user.Password == "" || user.Id == "" || !isValidUserId(user.Id) {
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
func addUser(c *gin.Context) {
	if !isAdmin(c) {
		c.String(404, "404 page not found")
		return
	}
	stuId := c.PostForm("id")
	username := c.PostForm("name")
	if stuId == "" || username == "" {
		c.Redirect(302, "/admin")
		return
	}
	gdb.Table("users").Create(&User{
		Id:       stuId,
		Username: username,
		Password: stuId,
	})
	c.Redirect(302, "/admin")
}
func deleteUser(c *gin.Context) {
	if !isAdmin(c) {
		c.String(404, "404 page not found")
		return
	}
	stuId := c.PostForm("id")
	if stuId == "" {
		c.Redirect(302, "/admin")
	}
	gdb.Table("like_record").Unscoped().Where("user_id = ?", stuId).Delete(nil)
	gdb.Table("yiyan").Unscoped().Where("submitter = ?", stuId).Delete(nil)
	gdb.Table("users").Unscoped().Where("id = ?", stuId).Delete(nil)
	c.Redirect(302, "/admin")
}
func deleteYiyan(c *gin.Context) {
	if !isAdmin(c) {
		c.String(404, "404 page not found")
		return
	}
	yiyanId := c.PostForm("id")
	if yiyanId == "" {
		c.Redirect(302, "/admin")
	}
	gdb.Table("yiyan").Unscoped().Where("id = ?", yiyanId).Delete(nil)
	c.Redirect(302, "/admin")
}
func getUserInfo(c *gin.Context) {
	if !isAdmin(c) {
		c.String(404, "404 page not found")
		return
	}
	var lk map[string]interface{}
	if err := c.ShouldBindJSON(&lk); err != nil {
		c.String(200, "format error!")
		return
	}
	yId, _ := strconv.Atoi(lk["id"].(string))
	var userId string
	var user User
	gdb.Table("yiyan").Select("submitter").Where("id = ?", yId).Pluck("submitter", &userId)
	gdb.Table("users").Where("id = ?", userId).First(&user)
	if user.Id == "" || user.Username == "" || user.Password == "" {
		c.String(200, "no")
	} else {
		c.JSON(200, user)
	}
}
func isAdmin(c *gin.Context) bool {
	stuId, errId := c.Cookie("id")
	username, errUsername := c.Cookie("username")
	password, errPassword := c.Cookie("password")
	if errId != nil || errUsername != nil || errPassword != nil {
		return false
	}
	var user *User
	gdb.Where("id = ? AND username = ? AND password = ?", stuId, username, password).First(&user)
	if user.Username == "" || user.Password == "" || user.Id == "" || !isValidUserId(user.Id) {
		return false
	}
	var realId string
	gdb.Table("admin").Select("id").Where("id = ?", stuId).Pluck("id", &realId)
	if stuId == realId {
		return true
	}
	return false
}

var gdb *gorm.DB

// 操控数据库、操控网页Cookie、Session存储用户信息等。
func main() {
	g := gin.Default()
	fmt.Println("已创建Gin~")
	dsn := mysqlUsername + ":" + mysqlPassword + "@tcp(" + IP + ":" + Port + ")/" + mysqlDatabase + "?charset=utf8mb4&parseTime=True&loc=Local"
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
		stuId, errStuId := c.Cookie("id")
		_, errUsername := c.Cookie("username")
		_, errPassword := c.Cookie("password")
		if errStuId != nil || errUsername != nil || errPassword != nil || !isValidUserId(stuId) {
			c.SetCookie("id", "", -1, "/", "", false, true)
			c.SetCookie("username", "", -1, "/", "", false, true)
			c.SetCookie("password", "", -1, "/", "", false, true)
			c.Redirect(302, "/login")
		} else {
			c.HTML(200, "submit.html", nil)
		}
	})
	g.GET("/my", func(c *gin.Context) {
		stuId, errStuId := c.Cookie("id")
		_, errUsername := c.Cookie("username")
		_, errPassword := c.Cookie("password")
		if errStuId != nil || errUsername != nil || errPassword != nil || !isValidUserId(stuId) {
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
	g.GET("/admin", func(c *gin.Context) {
		if isAdmin(c) {
			c.HTML(200, "admin.html", nil)
		} else {
			c.String(404, "404 page not found")
		}
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
		yiyan.POST("/addUser", addUser)
		yiyan.POST("/deleteUser", deleteUser)
		yiyan.POST("/deleteYiyan", deleteYiyan)
		yiyan.POST("/getUserInfo", getUserInfo)
	}
	fmt.Println("已创建动态链接~")
	errGin := g.Run(":8080")
	if errGin != nil {
		panic(errGin)
	}
	fmt.Println("已结束web服务~")
}
