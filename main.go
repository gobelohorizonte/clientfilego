package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

//var PATHDIR string = "/home/diego/Downloads"
//jeff.otoni@gmail.com
//1234

type login struct {
	statusLogin string `json:"status"`
	msg         string `json:"msg"`
	token       string `json:"token"`
}

func main() {

	pathdir := os.Args[1] //Receive path of the files.
	user := os.Args[2]    //Receive user for login.
	pass := os.Args[3]    //Receive pass for login.

	apiUrl := "https://fileserver.s3apis.com/"
	resource := "/v1/user/login"
	data := url.Values{}
	data.Add("password", pass)
	data.Add("user", user)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("X-Key", "ZmlsZXNlcnZlcjIwMThnb2xhbmdiaA==")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)

	//defer resp.Body.Close()

	if resp.StatusCode == 200 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Erro na leitura")
			return
		}
		ret := login{}
		err = json.Unmarshal(body, &ret)
		if err != nil {
			log.Println("Erro ao ler json")
		}

		fmt.Println(string(ret.msg))
		if ret.statusLogin == "" {

			if !verifyIsDir(pathdir) {
				log.Fatal("Doe's not an directoris valid or not exists.")
				return
			}
			/**
			read of the dir
			*/
			files, err := ioutil.ReadDir(pathdir)
			if err != nil {
				log.Fatal(err)
			}
			/**
			execute loop in files of the dir's
			*/
			resource = "/v1/file/upload"
			data := url.Values{}

			for _, file := range files {
				if !file.IsDir() {
					data.Add("file[]", file.Name())
				}
			}
			client = &http.Client{}
			r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
			r.Header.Add("Authorization", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiamVmZi5vdG9uaUBnbWFpbC5jb20iLCJ1aWQiOiI0ZjQ0NzgxMTAzODkwNjA5ZmY1MDdmNjIzMTdlOGExMGFiMDc4ZjFmIiwidWlkd2tzIjoiMDE2NTcxNmNiNzBmMWM1N2ZhMzhhZGY5MGI1Y2QyMTdmNDE1NTJhNyIsImV4cCI6MTUzNzkwMzczMywiaXNzIjoiand0IEZpbGVTZXJ2ZXIifQ.Sse1XOfFOr1rxjhJIugF2CEQfB6e1PX4XizBlhL_eRhgTq8WJ5gbHSY0ab22gyHBBpGyB9Pwb0mmpcbHgvfGcW4Xvt0RsGatXBEgkV8uKsVP1zTwQxH_zyAuLOMIicFPvw42lYiOEjNwsGP5ujjcaaNzVUCYkHnezPuDQVq1LNQ")
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

			resp, _ := client.Do(r)
			if resp.StatusCode == 200 {
				fmt.Printf("Post é: %+v\r\n", ret)
			}
		}

	} else {

	}
	/**

	curl -X POST https://fileserver.s3apis.com/v1/file/upload --form "file=@seuarquivo" -H "Authorization: Bearer <token>"
	verifica if is dir valid
	*/

}

func verifyIsDir(pathDir string) bool {
	if stat, err := os.Stat(pathDir); err == nil && stat.IsDir() {
		return true
	}
	return false
}

/*
func checkIsFile(pathFile string) bool {
	isFile, _ :=
	return true
}*/
