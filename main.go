package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 数据库连接配置
const (
	mysqlUsername = "root"
	mysqlPassword = "******"
	mysqlDatabase = "yiyandata"
	IP            = "127.0.0.1"
)

// 用户模型
type User struct {
	Id       string `json:"id"`       // 用户ID
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// 一言模型
type Yiyan struct {
	Id          int    `json:"id"`          // 一言ID
	Content     string `json:"content"`     // 内容
	Submitter   string `json:"submitter"`   // 提交者ID
	Source      string `json:"source"`      // 来源
	Author      string `json:"author"`      // 作者
	Classifiers string `json:"classifiers"` // 分类
	Likes       int    `json:"likes"`       // 点赞数
}

// 点赞记录模型
type LikeRecord struct {
	YiyanId int    `json:"yiyan_id"` // 一言ID
	UserId  string `json:"user_id"`  // 用户ID
}

// 指定表名
func (LikeRecord) TableName() string {
	return "like_record"
}

func (Yiyan) TableName() string {
	return "yiyan"
}

func (User) TableName() string {
	return "users"
}

// 一言结果模型（包含是否点赞字段）
type YiyanResult struct {
	Id          int    `json:"id"`          // 一言ID
	Content     string `json:"content"`     // 内容
	Submitter   string `json:"submitter"`   // 提交者ID
	Source      string `json:"source"`      // 来源
	Author      string `json:"author"`      // 作者
	Classifiers string `json:"classifiers"` // 分类
	Likes       int    `json:"likes"`       // 点赞数
	IsLiked     bool   `json:"is_liked"`    // 是否已点赞
}

var gdb *gorm.DB

// 中间件：检查登录状态
func checkLogin(c *gin.Context) {
	stuId, errId := c.Cookie("id")
	_, errUsername := c.Cookie("username")
	_, errPassword := c.Cookie("password")
	if errId != nil || errUsername != nil || errPassword != nil {
		c.SetCookie("id", "", -1, "/", "", false, true)
		c.SetCookie("username", "", -1, "/", "", false, true)
		c.SetCookie("password", "", -1, "/", "", false, true)
		c.Redirect(302, "/login")
		c.Abort()
		return
	}
	c.Set("stuId", stuId)
	c.Next()
}

func main() {
	// 初始化Gin
	g := gin.Default()
	fmt.Println("已创建Gin~")
	// 连接MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUsername, mysqlPassword, IP, mysqlDatabase)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	gdb = db
	fmt.Println("已连接MySQL~")
	// 设置模板和静态资源
	g.Delims("[[", "]]")
	g.LoadHTMLGlob("templates/*")
	g.Static("/static", "./static")
	// 设置路由
	g.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	g.GET("/list", func(c *gin.Context) {
		c.HTML(200, "list.html", nil)
	})
	g.GET("/submit", checkLogin, func(c *gin.Context) {
		c.HTML(200, "submit.html", nil)
	})
	g.GET("/my", checkLogin, func(c *gin.Context) {
		c.HTML(200, "my.html", nil)
	})
	g.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	g.GET("/privacy", func(c *gin.Context) {
		c.HTML(200, "privacy.html", nil)
	})
	fmt.Println("已加载HTML静态资源~")
	// 设置API路由
	api := g.Group("/yiyan")
	{
		api.GET("/get_random_one", getRandomOne)
		api.GET("/get_all", getAll)
		api.GET("/get_my", checkLogin, getMy)
		api.GET("get_most", getMost)
		api.POST("/like", checkLogin, like)
		api.POST("/submit", checkLogin, submit)
		api.POST("/login", login)
	}
	fmt.Println("已创建动态链接~")
	// 启动服务
	if err := g.Run(":8080"); err != nil {
		panic(err)
	}
	fmt.Println("已结束web服务~")
}

// 判断是否是自己的作品
func isMyYiyanToLike(yiyanId int, userId string) bool {
	var submitter string
	gdb.Model(&Yiyan{}).Select("submitter").Where("id = ?", yiyanId).Scan(&submitter)
	return submitter == userId
}

// 获取是否点赞
func getIsLike(yiyanId int, userId string) bool {
	if userId == "" {
		return false
	}
	var count int64
	gdb.Model(&LikeRecord{}).Where("yiyan_id = ? AND user_id = ?", yiyanId, userId).Count(&count)
	return count > 0
}

// 获取随机一言
func getRandomOne(c *gin.Context) {
	var yiyanResult YiyanResult
	gdb.Model(&Yiyan{}).Order("RAND()").First(&yiyanResult)
	stuId, _ := c.Cookie("id")
	yiyanResult.IsLiked = getIsLike(yiyanResult.Id, stuId)
	c.JSON(200, yiyanResult)
}

// 获取点赞最多的一言
func getMost(c *gin.Context) {
	var yiyanGet YiyanResult
	gdb.Model(&Yiyan{}).Order("likes DESC").First(&yiyanGet)
	stuId, _ := c.Cookie("id")
	yiyanGet.IsLiked = getIsLike(yiyanGet.Id, stuId)
	c.JSON(200, yiyanGet)
}

// 获取所有一言
func getAll(c *gin.Context) {
	var yiyanAll []YiyanResult
	gdb.Model(&Yiyan{}).Find(&yiyanAll)
	stuId, _ := c.Cookie("id")
	for _, yiyan := range yiyanAll {
		yiyan.IsLiked = getIsLike(yiyan.Id, stuId)
	}
	c.JSON(200, yiyanAll)
