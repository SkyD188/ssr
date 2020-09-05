package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

var cookie []*http.Cookie

func main() {
	app := iris.New()
	app.Get("/get", func(c context.Context) {
		var s string
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < 9; i++ {
			s += strconv.Itoa(rand.Intn(10))
		}
		resp, err := http.Post("https://hashcloud.one/api/v1/passport/auth/register",
			"application/x-www-form-urlencoded",
			strings.NewReader("email=1"+s+"@qq.com&password="+s+""))
		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		cookie = resp.Cookies()

		c.WriteString(getSubUrl())
	})
	app.Run(iris.Addr(":9999"))
}

// get请求https://hashcloud.one/api/v1/user/getSubscribe
func getSubUrl() string {
	client := http.Client{}

	resp, err := http.NewRequest("GET", "https://hashcloud.one/api/v1/user/getSubscribe", nil)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(cookie); i++ {
		resp.AddCookie(cookie[i])
	}
	do, err := client.Do(resp)
	if err != nil {
		fmt.Println(err)
	}

	defer do.Body.Close()
	getBody, err := ioutil.ReadAll(do.Body)
	if err != nil {
		fmt.Println(err)
	}
	var s struct {
		Data struct {
			SubscribeUrl string `json:"subscribe_url"`
		}
	}
	_ = json.Unmarshal(getBody, &s)

	return s.Data.SubscribeUrl
}
