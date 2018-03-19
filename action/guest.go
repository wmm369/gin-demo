package recognition

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Res struct {
	Res   int            `json:"res"`
	Score float32        `json:"score"`
}

// Guest 人脸1比1
func Guest(c *gin.Context) {

	memberID := c.PostForm("member_id")
	memberName := c.PostForm("member_name")
	image1 := c.PostForm("image1")
	image2 := c.PostForm("image2")
	if memberID == "" || memberName == "" || image1 == "" || image2 == "" {
		c.JSON(http.StatusOK, gin.H{"res": 1002, "message": "参数缺失", "ip": c.ClientIP()})
		return
	}
	data := map[string]interface{}{
		"ctype": 0,
		"img1":  image1,
		"img2":  image2}

	res := Res{}
	json.Unmarshal([]byte(getRequest(data)), &res)

	if res.Res == 0 {
		c.JSON(http.StatusOK, res)
		return
	}
	c.JSON(http.StatusOK, gin.H{"res": 1003, "message": "比对失败", "ip": c.ClientIP()})
}

func getRequest(data map[string]interface{}) string {

	b, _ := json.Marshal(data)
	client := &http.Client{
		Timeout: 2 * time.Second, //设置超时
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1/test", strings.NewReader(string(b)))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body)
}
