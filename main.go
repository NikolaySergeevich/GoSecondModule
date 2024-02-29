package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type User struct{
	Name string `json:"name"`
	Age string `json:"age"`
	Friends []*User `json:"friends"` 
}

type Band struct{
	count int
	Team map[int]*User
}

func main(){
	mux := http.NewServeMux()
	team := Band{0, make(map[int]*User)}
	mux.HandleFunc("/create", team.Create)
	mux.HandleFunc("/make_friends", team.Make_friends)
	// mux.HandleFunc("/friends/", team.Friends)//получаем друзей. В запросе указывается ID пользователя
	// mux.HandleFunc("/make_friends", team.Make_friends)
	// mux.HandleFunc("/user", team.UserDel)
	// mux.HandleFunc("/user_id", team.UpdateUser)

	http.ListenAndServe(":8080", mux)
}

func (b *Band) Create(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		content, er := ioutil.ReadAll(r.Body)
		if er != nil{
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(er.Error()))
			return
		}
		defer r.Body.Close()

		var u User
		if err := json.Unmarshal(content, &u); err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(er.Error()))
			return
		}
		b.count = b.count + 1
		b.Team[b.count] = &u

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Пользователь создан c ID: " + strconv.Itoa((b.count))))
	}
	w.WriteHeader(http.StatusBadGateway) 	
}

func (b *Band) Make_friends(w http.ResponseWriter, r *http.Request){
	fmt.Println(string(r.RequestURI))
	if r.Method == "PUT"{
		content , er := ioutil.ReadAll(r.Body)
		if er != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(er.Error()))
			return
		}
		defer r.Body.Close()

		data := make(map[string]int)
		if err := json.Unmarshal(content, &data); err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			// w.Write([]byte(err.Error()))
			w.Write([]byte("Это тут 1"))
			return
		} 
		// fmt.Println(data)
		source_id := data["source_id"]
		target_id := data["target_id"]
		// fmt.Println("Id source_id = " + strconv.Itoa(source_id))
		// fmt.Println("Id target_id = " + strconv.Itoa(target_id))
		userSource, ok := b.Team[source_id]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Нет пользователя с ID = "+ strconv.Itoa(source_id)))
			return
		}
		userTarget, ok := b.Team[target_id]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Нет пользователя с ID = "+ strconv.Itoa(target_id)))
			return
		}
		userTarget.Friends = append(userTarget.Friends, userSource)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Пользователь " +userSource.Name+ " теперь в друзьях у " +userTarget.Name))
	}
	w.WriteHeader(http.StatusBadGateway)
}