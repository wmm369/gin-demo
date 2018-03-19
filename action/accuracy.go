package recognition

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
)

//Per 查询结果
type Per struct {
	ID    bson.ObjectId `json:"id"    bson:"_id"`
	Name  string        `json:"name " `
	Phone string        `json:"phone" `
	Age   int           `json:"age"`
}

//Men 集合
type Men struct {
	Per []Per
}

//Accuracy 人脸1比n
func Accuracy(c *gin.Context) {
	devCode := c.PostForm("dev_code")
	ctype := c.DefaultPostForm("ctype", "0")
	memberID := c.PostForm("member_id")
	memberName := c.PostForm("member_name")
	if devCode == "" || memberID == "" || memberName == "" {
		c.JSON(http.StatusOK, gin.H{"res": 1002, "message": "参数缺失", "ip": c.ClientIP()})
		return
	}
	ctypes, _ := strconv.Atoi(ctype)
	if ctypes == 1 {

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

	client := mongo.DB("testting").C("gin_db")
	res, err := client.Find(bson.M{"status": 1}).Count()
	result := Per{}
	iter := client.Find(bson.M{"status": 1}).Sort("_id").Iter()
	var personAll Men
	for iter.Next(&result) {
		personAll.Per = append(personAll.Per, result)
	}

	if err := iter.Close(); err != nil {
		data["res"] = 1431
		data["message"] = "数据库异常"
		c.JSON(http.StatusOK, data)
		return
	}
	mongo.Close()
	data["res"] = 0
	data["total"] = res
	data["data"] = personAll.Per
	c.JSON(http.StatusOK, data)
}
