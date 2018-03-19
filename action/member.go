package recognition

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
)

//Persion 数据
type Persion struct {
	Name   string `bson:"name" `
	Phone  string `bson:"phone" `
	Age    int    `bson:"age"`
	Status int    `bson:"status"`
}

//User 数据信息
type User struct {
	Name   string `json:"name" `
	Phone  string `json:"phone" `
	Age    int    `json:"age"`
	Status int    `json:"status"`
}

//Users 数据集合
type Users struct {
	Users []User
}

//User 数据信息
type ImageOne struct {
	InfoID      string      `json:"infoId"  bson:"infoId"`
	Feature     bson.Binary `json:"roi_image"  bson:"roi_image"`
	Status      int         `json:"status" bson:"status"`
}

//Insert 插入一个
func Insert(c *gin.Context) {
	name := c.PostForm("name")
	age := c.DefaultPostForm("age", "17")
	phone := c.PostForm("phone")

	if name == "" || phone == "" {
		c.JSON(http.StatusOK, gin.H{"res": 1002, "message": "参数缺失", "ip": c.ClientIP()})
		return
	}
	data := map[string]interface{}{}
	mongo, err := mgo.Dial("127.0.0.1")
	defer mongo.Close()
	if err != nil {

		data["res"] = 1431
		data["message"] = "数据库异常"
		c.JSON(http.StatusOK, data)
		return
	}

	client := mongo.DB("testting").C("test_db")

	ageInt, _ := strconv.Atoi(age)
	info := Persion{
		Name:   name,
		Age:    ageInt,
		Status: 0,
		Phone:  phone}
	err = client.Insert(&info)
	if err != nil {
		data["res"] = 1002
		data["message"] = "插入数据失败"
		c.JSON(http.StatusOK, data)
		return
	}

	data["res"] = 0
	data["message"] = "插入数据成功"
	c.JSON(http.StatusOK, data)
}

//FindOne 查询一个
func FindOne(c *gin.Context) {
	name := c.PostForm("name")

	if name == "" {
		c.JSON(http.StatusOK, gin.H{"res": 1002, "message": "参数缺失", "ip": c.ClientIP()})
		return
	}
	data := map[string]interface{}{}
	mongo, err := mgo.Dial("127.0.0.1")
	defer mongo.Close()
	if err != nil {

		data["res"] = 1431
		data["message"] = "数据库异常"
		c.JSON(http.StatusOK, data)
		return
	}

	client := mongo.DB("testting").C("test_db")
	user := User{}
	err = client.Find(bson.M{"name": name}).One(&user)
	if err != nil {
		data["res"] = 1003
		data["message"] = "无用户"
		data["name"] = name
		data["data"] = user
		c.JSON(http.StatusOK, data)
		return
	}
	data["res"] = 0
	data["data"] = user
	c.JSON(http.StatusOK, data)
}

//Find 查询所有
func Find(c *gin.Context) {
	name := c.PostForm("name")

	if name == "" {
		c.JSON(http.StatusOK, gin.H{"res": 1002, "message": "参数缺失", "ip": c.ClientIP()})
		return
	}
	data := map[string]interface{}{}
	mongo, err := mgo.Dial("127.0.0.1")
	defer mongo.Close()
	if err != nil {

		data["res"] = 1431
		data["message"] = "数据库异常"
		c.JSON(http.StatusOK, data)
		return
	}

	client := mongo.DB("testting").C("test_db")
	user := User{}
	total, err := client.Find(bson.M{"name": name}).Count()
	iter := client.Find(bson.M{"name": name}).Sort("_id").Skip(1).Limit(15).Iter()

	if err != nil {
		data["res"] = 1003
		data["message"] = "无用户"
		data["name"] = name
		data["data"] = user
		c.JSON(http.StatusOK, data)
		return
	}

	var users Users
	for iter.Next(&user) {
		users.Users = append(users.Users, user)
	}
	if err := iter.Close(); err != nil {
		data["res"] = 1431
		data["message"] = "数据库异常"
		c.JSON(http.StatusOK, data)
		return
	}

	data["res"] = 0
	data["total"] = total
	data["data"] = users.Users
	c.JSON(http.StatusOK, data)
}

//Update 更新数据
func Update(c *gin.Context) {
	name := c.PostForm("name")
	status := c.DefaultPostForm("status", "0")
	if name == "" {
		c.JSON(http.StatusOK, gin.H{"res": 1002, "message": "参数缺失", "ip": c.ClientIP()})
		return
	}
	data := map[string]interface{}{}
	mongo, err := mgo.Dial("127.0.0.1")
	defer mongo.Close()
	if err != nil {

		data["res"] = 1431
		data["message"] = "数据库异常"
		c.JSON(http.StatusOK, data)
		return
	}

	client := mongo.DB("testting").C("test_db")
	statusInt, err := strconv.Atoi(status)
	//更新一个
	err = client.Update(bson.M{"name": name}, bson.M{"$set": bson.M{"status": statusInt}})

	_, err = client.UpdateAll(bson.M{"name": name}, bson.M{"$set": bson.M{"status": statusInt}})

	if err != nil {
		data["res"] = 1431
		data["message"] = "数据库异常"
		c.JSON(http.StatusOK, data)
		return
	}
	data["res"] = 0
	data["message"] = "更新成功"
	c.JSON(http.StatusOK, data)
}
func Image(c *gin.Context) {
	name := c.PostForm("name")

	if name == "" {
		c.JSON(http.StatusOK, gin.H{"res": 1002, "message": "参数缺失", "ip": c.ClientIP()})
		return
	}
	data := map[string]interface{}{}
	mongo, err := mgo.Dial("127.0.0.1:27017")
	defer mongo.Close()
	if err != nil {

		data["res"] = 1431
		data["message"] = "数据库异常"
		c.JSON(http.StatusOK, data)
		return
	}

	client := mongo.DB("testting").C("t_image")
	user := ImageOne{}
	err = client.Find(bson.M{"_id": bson.ObjectIdHex("5a0ada89f661c52f10428f70")}).One(&user)
	if err != nil {
		data["res"] = 1003
		data["message"] = "无用户"
		data["name"] = name
		data["data"] = user
		c.JSON(http.StatusOK, data)
		return
	}
	c.Header("Content-type","image/jpg")
	c.String(200, string(user.Feature.Data))
}
