package service

import (
	"fmt"
	"github.com/imroc/req/v3"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestQB(t *testing.T) {
	client := req.C().
		SetTimeout(5 * time.Second)
	if true {
		client = client.DevMode()
	}

	param := map[string]interface{}{
		"username": "admin",
		"password": "adminadmin",
	}

	res, err := client.R().SetFormDataAnyType(param).Post("http://localhost:9999/api/v2/auth/login")
	if err != nil {
		panic(err)
		return
	}
	if res.IsError() {
		panic(res.Error())
		return
	}
	if res.IsSuccess() {
		fmt.Println(res.String())
	}

	get, err := client.R().Get("http://localhost:9999/api/v2/log/main")
	if err != nil {
		panic(err)
		return
	}
	if get.IsError() && get.StatusCode == http.StatusForbidden {
		log.Println("未登录")
		return
	}
	if get.IsSuccess() {
		fmt.Println("获取成功")
		fmt.Println(get.String())
	}

}
