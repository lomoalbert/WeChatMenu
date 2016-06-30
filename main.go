package main

import (
	"fmt"

	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/menu"
	"encoding/json"
	"io/ioutil"
)

var menuClient *menu.Client

func main() {
	var name string
	var option string
	keyconfig := make(map[string]string)
	for {
		fmt.Println("请输入主文件名")
		fmt.Scanf("%s", &name)
		data, err := readfile(name + ".key.txt")
		if !check(err) {
			fmt.Println("未找到文件:" + name + ".key.txt")
			err = writefile(name + ".key.txt",[]byte("{\r\n\t\"AppID\":\"\",\r\n\t\"AppSecret\":\"\"\r\n}"))
                        if check(err){
				fmt.Println("已创建空文件,请填写公众号appid及appsecret:" + name + ".key.txt")
			}
			continue
		}
		err = json.Unmarshal(data, &keyconfig)
		if !check(err) {
			fmt.Println("数据格式错误:" + name + ".key.txt")
			continue
		}
		fmt.Println(keyconfig)
		break
	}
	initclient(keyconfig["AppID"], keyconfig["AppSecret"])
	for {
		fmt.Println("[1]下载当前菜单配置文件")
		fmt.Println("[2]上传新菜单配置文件")
		fmt.Println("[3]退出")
		fmt.Scanf("%s", &option)
		switch option {
		case "1":
			mymenu, _, err := menuClient.GetMenu();
			if !check(err) {
				continue
			}
			menujson, _ := json.MarshalIndent(mymenu, "", "\t")
			err = writefile(name + ".menu.txt", menujson)
			if !check(err) {
				continue
			}
		case "2":
			data, err := readfile(name + ".menu.txt")
			if !check(err) {
				continue
			}
			newmenu := new(menu.Menu)
			err = json.Unmarshal(data, newmenu)
			if !check(err) {
				continue
			}
			err = menuClient.CreateMenu(*newmenu)
			if !check(err) {
				continue
			}
		case "3":
			return
		default:
			continue
		}
	}
}

func check(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func readfile(filename string) ([]byte, error) {
	fmt.Println("Reading ", filename)
	data, err := ioutil.ReadFile(filename)
	return data, err
}

func writefile(filename string, data []byte) error {
	fmt.Println("Save to ", filename)
	err := ioutil.WriteFile(filename, data, 0644)
	return err
}

func initclient(AppID, AppSecret string) {
	var AccessTokenServer = mp.NewDefaultAccessTokenServer(AppID, AppSecret, nil) // 一個應用只能有一個實例
	var mpClient = mp.NewClient(AccessTokenServer, nil)
	menuClient = (*menu.Client)(mpClient)
}