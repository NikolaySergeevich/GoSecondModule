package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	Name    string  `json:"name"`
	Age     string  `json:"age"`
	Friends []*User `json:"friends"`
}

func (u *User) toString() (res string) {
	res = u.Name + ", " + u.Age + ", друзья: "
	for _, f := range u.Friends {
		res = res + f.Name + "; "
	}
	res = res + "\n"
	return
}

type Band struct {
	count int
	Team  map[int]*User
}

func main() {
	mux := http.NewServeMux()
	team := Band{0, make(map[int]*User)}
	mux.HandleFunc("/create", team.Create)
	mux.HandleFunc("/make_friends", team.Make_friends)
	mux.HandleFunc("/friends/", team.Friends) //получаем друзей. В запросе указывается ID пользователя
	mux.HandleFunc("/user", team.UserDel)
	// mux.HandleFunc("/user_id", team.UpdateUser)
	mux.HandleFunc("/get_all", team.GetAll)

	http.ListenAndServe(":8080", mux)
}

func (b *Band) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, er := ioutil.ReadAll(r.Body)
		if er != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(er.Error()))
			return
		}
		defer r.Body.Close()

		var u User
		if err := json.Unmarshal(content, &u); err != nil {
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

func (b *Band) Make_friends(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		content, er := ioutil.ReadAll(r.Body)
		if er != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(er.Error()))
			return
		}
		defer r.Body.Close()

		data := make(map[string]int)
		if err := json.Unmarshal(content, &data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// w.Write([]byte(err.Error()))
			w.Write([]byte("Это тут 1"))
			return
		}
		// fmt.Println(data)
		source_id := data["source_id"]
		target_id := data["target_id"]
		userSource, ok := b.Team[source_id]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Нет пользователя с ID = " + strconv.Itoa(source_id)))
			return
		}
		userTarget, ok := b.Team[target_id]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Нет пользователя с ID = " + strconv.Itoa(target_id)))
			return
		}
		for _, v := range userTarget.Friends {
			if v == userSource {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Пользователи " + userSource.Name + " и " + userTarget.Name + " уже друзья"))
				return
			}
		}
		userTarget.Friends = append(userTarget.Friends, userSource)
		userSource.Friends = append(userSource.Friends, userTarget)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Пользователи " + userSource.Name + " и " + userTarget.Name + " подружились"))
	}
	w.WriteHeader(http.StatusBadGateway)
}

func (b *Band) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		res := ""
		for _, v := range b.Team {
			res = res + v.toString()
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
	}
	w.WriteHeader(http.StatusBadGateway)
}

func (b *Band) Friends(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sl := strings.Split(r.RequestURI, "/")
		id, er := strconv.Atoi(sl[2])
		if er != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("id пользователя должно быть число"))
			return
		}
		user, ok := b.Team[id]
		if !ok {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Пользователя с id = " + strconv.Itoa(id) + " нет."))
			return
		}
		var res string
		if len(user.Friends) == 0 {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("У пользователя с id = " + strconv.Itoa(id) + " нет друзей"))
			return
		}
		for _, v := range user.Friends {
			res = res + v.toString()
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Друзья " + user.Name + ":\n" + res))
		return
	}
	w.WriteHeader(http.StatusBadGateway)
}

func (b *Band) UserDel(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		content, er := ioutil.ReadAll(r.Body)
		if er != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(er.Error()))
			return
		}
		defer r.Body.Close()

		data := make(map[string]int)
		if err := json.Unmarshal(content, &data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		id := data["target_id"]
		us, ok := b.Team[id]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Пользователя с id = " + strconv.Itoa(id) + " нет."))
		}
		for _, v := range us.Friends {
			for ind, f := range v.Friends {
				if f == us {
					if len(v.Friends) == 1 {
						v.Friends = v.Friends[:len(v.Friends)-1]
					} else {
						v.Friends[ind] = v.Friends[len(v.Friends)-1]
						v.Friends = v.Friends[:len(v.Friends)-1]
					}
				}
			}
		}
		delete(b.Team, id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Пользователь удалён"))
	}
	w.WriteHeader(http.StatusBadGateway)
}
