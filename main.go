package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"io"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

//var PATHDIR string = "/home/diego/Downloads"
//jeff.otoni@gmail.com
//1234

type login struct {
	Status  string `json:"status"`
	Msg     string `json:"msg"`
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

type userLogin struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func main() {

	if len(os.Args) != 4 {

		fmt.Println("Entre com os parametros")
		fmt.Println(" -path <diretorio onde estara os arquivos>")
		fmt.Println(" -user <jeff.otoni@gmail.com>")
		fmt.Println(" -senha <1234>")
		os.Exit(0)
	}

	pathdir := os.Args[1] //Receive path of the files.
	user := os.Args[2]    //Receive user for login.
	pass := os.Args[3]    //Receive pass for login.

	apiUrl := "https://fileserver.s3apis.com/"
	resource := "/v1/user/login"

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	userjson := &userLogin{User: user, Password: pass}

	bjson, err := json.Marshal(userjson)
	if err != nil {
		log.Fatal("Erro ao fazer Marshal!")
		return
	}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(bjson)) // URL-encoded payload

	r.Header.Add("X-Key", "ZmlsZXNlcnZlcjIwMThnb2xhbmdiaA==")
	r.Header.Add("Content-Type", "application/json")

	resp, _ := client.Do(r)
	defer resp.Body.Close()

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

		// pegando o retorno
		fmt.Println(string(ret.Msg))
		fmt.Println(string(ret.Status))
		fmt.Println(string(ret.Expires))
		fmt.Println(string(ret.Token))

		// pode fazer..
		if ret.Status == "ok" && ret.Token != "" {
			if !verifyIsDir(pathdir) {
				log.Println("Doe's not an directoris valid or not exists.")
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
			u2, _ := url.ParseRequestURI(apiUrl)
			u2.Path = resource
			urlStr := u2.String()

			for _, file := range files {
				if !file.IsDir() {
					pathdirNowFile := pathdir + "/" + file.Name()
					postFile(ret.Token, pathdirNowFile, urlStr)
				}
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

func postFile(Token, filename string, targetUrl string) error {

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	client := &http.Client{}
	//postData := make([]byte, 100)
	resp, err := http.NewRequest("POST", targetUrl, bodyBuf)
	if err != nil {
		return err
	}
	resp.Header.Add("Content-Type", contentType)
	resp.Header.Add("Authorization", "Bearer "+Token)
	resp2, err := client.Do(resp)

	defer resp2.Body.Close()

	fmt.Println(resp2)

	return nil
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()

	if err != nil {

		return nil, err
	}

	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}

	// copy...
	//io.Copy(part, file)

	part.Write(fileContents)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", uri, body)
}
