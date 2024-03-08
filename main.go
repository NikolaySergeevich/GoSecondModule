package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	counter           int    = 0
	firstInstansHost  string = "http://localhost:8081"
	secondInstansHost string = "http://localhost:8082"
	// chanServ1 chan *http.Request
	// chanServ2 chan *http.Request
)

func main() {
	http.HandleFunc("/create", Create)
	http.HandleFunc("/make_friends", Put)
	http.HandleFunc("/friends/", Get)
	http.HandleFunc("/user", UserDel)
	http.HandleFunc("/", Put)
	http.HandleFunc("/get_all", Get)
	http.HandleFunc("/help", Get) //выставляем ручку
	fmt.Println("PROXY is runing")
	log.Fatalln(http.ListenAndServe("localhost:9000", nil)) //включаем сервер.
}
func Create(w http.ResponseWriter, r *http.Request) {
	defer func(){
		if arg := recover(); arg != nil{
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Скорее всего что-то с запросом"))
		}
	}()
	courentyHost := ""
	if r.Method == "POST" {
		if counter == 0 {
			courentyHost = firstInstansHost
			counter++
		} else {
			courentyHost = secondInstansHost
			counter--
		}
		textBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}
		defer r.Body.Close()
		text := string(textBytes)

		resp, err := http.Post(courentyHost+"/create", "Content-Type: application/json", bytes.NewBuffer([]byte(text)))
		if err != nil {
			log.Fatal(err)
		}

		textByteResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		if counter == 1 {
			fmt.Println("Отработал ПЕРВЫЙ сервер")
		} else {
			fmt.Println("Отработал ВТОРОЙ сервер")
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(textByteResp))
		return
	}
	w.WriteHeader(http.StatusBadGateway)
}

func Get(w http.ResponseWriter, r *http.Request) {
	courentyHost := ""
	if r.Method == "GET" {
		if counter == 0 {
			courentyHost = firstInstansHost
			counter++
		} else {
			courentyHost = secondInstansHost
			counter--
		}
		resp, err := http.Get(courentyHost + r.RequestURI)
		if err != nil {
			log.Fatal(err)
		}
		textByteResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		if counter == 1 {
			fmt.Println("Отработал ПЕРВЫЙ сервер")
		} else {
			fmt.Println("Отработал ВТОРОЙ сервер")
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(textByteResp))
		return
	}
	w.WriteHeader(http.StatusBadGateway)
}
func Put(w http.ResponseWriter, r *http.Request) {
	courentyHost := ""
	if r.Method == "PUT" {
		if counter == 0 {
			courentyHost = firstInstansHost
			counter++
		} else {
			courentyHost = secondInstansHost
			counter--
		}

		textByte, er := ioutil.ReadAll(r.Body)
		if er != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(er.Error()))
		}

		req, err := http.NewRequest(http.MethodPut, courentyHost+r.RequestURI, bytes.NewBuffer(textByte))
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		// Выполняем запрос
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		textByteResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		if counter == 1 {
			fmt.Println("Отработал ПЕРВЫЙ сервер")
		} else {
			fmt.Println("Отработал ВТОРОЙ сервер")
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(textByteResp))
		return
	}
	w.WriteHeader(http.StatusBadGateway)
}

func UserDel(w http.ResponseWriter, r *http.Request){
	courentyHost := ""
	if r.Method == "DELETE" {
		if counter == 0 {
			courentyHost = firstInstansHost
			counter++
		} else {
			courentyHost = secondInstansHost
			counter--
		}

		textByte, er := ioutil.ReadAll(r.Body)
		if er != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(er.Error()))
		}

		req, err := http.NewRequest("DELETE", courentyHost+"/user", bytes.NewBuffer(textByte))
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		// Выполняем запрос
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		textByteResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		if counter == 1 {
			fmt.Println("Отработал ПЕРВЫЙ сервер")
		} else {
			fmt.Println("Отработал ВТОРОЙ сервер")
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(textByteResp))
		return
	}
	w.WriteHeader(http.StatusBadGateway)
}
