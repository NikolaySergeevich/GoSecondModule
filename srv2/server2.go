package main
// curl -X POST -H 'Content-Type: application/json' -d '{"name":"Nikola","age":"26","friends":[]}' http://localhost:9000/create
// curl -X PUT -H 'Content-Type: application/json' -d '{"source_id": 1, "target_id": 2}' http://localhost:9000/make_friends
// curl -X GET  http://localhost:9000/friends/1
// curl -X PUT -H 'Content-Type: application/json' -d '{"new age":33}' http://localhost:9000/1
// curl -X GET  http://localhost:9000/get_all
// curl -X DELETE -H 'Content-Type: application/json' -d '{"target_id": 2}' http://localhost:9000/user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hom/db"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

type User struct {
	Name    string  `json:"name"`
	Age     string  `json:"age"`
	Friends []*User `json:"friends"`
}

const TEXT_HELP string = "Эта программа позволяет добавлять пользователей в базу, подружить пользователей между собой, удалять его из базы,\n" +
	"посмотреть на друзей конкретного пользователя, обновить возраст пользователю.\n" +
	"Запросы, которыми вы можете воспользоваться:\n" +
	"/create - добавить пользователя\n" +
	"/make_friends - подружить двоих пользователей\n" +
	"/user - удаляет пользователя\n" +
	"/friends/user_id - получить друзей указанного пользователя\n" +
	"/user_id - обновляет возраст пользователя\n" +
	"/help - просмотреть доступные запросы\n" +
	"/get_all - посмотреть всех пользователей\n"

func (u *User) toString() (res string) {
	res = u.Name + ", " + u.Age + ", друзья: "
	for _, f := range u.Friends {
		res = res + f.Name + "; "
	}
	res = res + "\n"
	return
}

type Band struct {
	// count int
	Team map[int]*User
}

func main() {
	db, err := sql.Open("mysql", "root:I#dohom38@/usersDB")

	if err != nil {
		log.Println(err)
		log.Println("Не подключился к БД")
	}
	database = db
	defer db.Close()

	mux := http.NewServeMux()
	// team := Band{make(map[int]*User)}
	mux.HandleFunc("/create", Create)
	mux.HandleFunc("/make_friends", Make_friends)
	mux.HandleFunc("/friends/", Friends)
	mux.HandleFunc("/user", UserDel)
	mux.HandleFunc("/", UpdateUser)
	mux.HandleFunc("/get_all", GetAll)
	mux.HandleFunc("/help", Help)

	fmt.Println("Сервер № 2 запущен")
	http.ListenAndServe(":8082", mux)
}

/*
Обработчик создаёт нового пользователя
*/
func Create(w http.ResponseWriter, r *http.Request) {
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
		queryString := fmt.Sprintf("CALL AddUser('%s', %s)", u.Name, u.Age)
		rowsUsers, err := database.Query(queryString)
		if err != nil {
			log.Println(err)
		}
		defer rowsUsers.Close()
		id := 0
		for rowsUsers.Next() {
			err := rowsUsers.Scan(&id)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Пользователь создан c ID: " + strconv.Itoa((id))))
	}
	w.WriteHeader(http.StatusBadGateway)
}

/*
Обработчик обрабатывает PUT запрос. Проверяет есть ли заявленные пользователи и добавляет их друг другу в друзья
*/
func Make_friends(w http.ResponseWriter, r *http.Request) {
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
			w.Write([]byte(err.Error()))
			return
		}
		// fmt.Println(data)
		source_id := data["source_id"]
		target_id := data["target_id"]
		if !db.Check_exists_user(database, source_id) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Нет пользователя с ID = " + strconv.Itoa(source_id)))
			return
		}
		if !db.Check_exists_user(database, target_id) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Нет пользователя с ID = " + strconv.Itoa(target_id)))
			return
		}
		if db.CheckFriends(database, source_id, target_id) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Пользователи c ID " + strconv.Itoa(source_id) + " и " + strconv.Itoa(target_id) + " уже друзья"))
			return
		}
		err := db.CreateFriendship(database, source_id, target_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Пользователи " + strconv.Itoa(source_id) + " и " + strconv.Itoa(target_id) + " подружились"))
	}
	w.WriteHeader(http.StatusBadGateway)
}

/*
Обработчик обработывает GET запрос и отправляет клиенту всех пользователей
*/
func GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		rowsUsers, err := database.Query("select * from Users")
		if err != nil {
			log.Println(err)
		}
		defer rowsUsers.Close()
		var id int
		users := []User{}
		for rowsUsers.Next() {
			u := User{}
			err := rowsUsers.Scan(&id, &u.Name, &u.Age)
			if err != nil {
				fmt.Println(err)
				continue
			}
			rowsFriends, err := database.Query("select target_id FROM Friends WHERE source_id = " + strconv.Itoa(id))
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer rowsFriends.Close()
			for rowsFriends.Next() {
				idFriend := 0
				friend := User{}
				err := rowsFriends.Scan(&idFriend)
				rowsFriend, er := database.Query("select * FROM Users WHERE id = " + strconv.Itoa(idFriend))
				if er != nil {
					fmt.Println(err)
					continue
				}
				defer rowsFriend.Close()
				for rowsFriend.Next() {
					err := rowsFriend.Scan(&idFriend, &friend.Name, friend.Age)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
				u.Friends = append(u.Friends, &friend)
			}
			users = append(users, u)
		}
		res := ""
		for _, v := range users {
			res = res + v.toString()
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
	}
	w.WriteHeader(http.StatusBadGateway)
}

/*
Обработчик обрабатывает GET запрос и отправляет клиенту друзей того пользователя, которого он запросил
*/
func Friends(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sl := strings.Split(r.RequestURI, "/")
		id, er := strconv.Atoi(sl[2])
		if er != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("id пользователя должно быть число"))
			return
		}
		if !db.Check_exists_user(database, id) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Пользователя с ID = " + strconv.Itoa(id) + " нет."))
			return
		}
		rowsFr, err := database.Query(fmt.Sprintf("CALL givUssFr (%s)", strconv.Itoa(id)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		res := []*User{}
		for rowsFr.Next() {
			u := &User{}
			err := rowsFr.Scan(&u.Name, &u.Age)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Не смог прочитать USER"))
				return
			}
			res = append(res, u)
		}
		var resStr string
		if len(res) == 0{
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("У пользователя с id = " + strconv.Itoa(id) + " нет друзей"))
			return
		}
		
		for _, v := range res {
			resStr = resStr + v.toString()
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Друзья указанного пользователя:\n" + resStr))
		return
	}
	w.WriteHeader(http.StatusBadGateway)
}

func Help(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(TEXT_HELP))
	}
	w.WriteHeader(http.StatusBadGateway)
}

/*
Обработчик обрабатывает DELETE запрос и удаляет указанного в теле запроса пользователя
*/
func UserDel(w http.ResponseWriter, r *http.Request) {
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
		if !db.Check_exists_user(database, id){
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Пользователя с ID %s нет", strconv.Itoa(id))))
			return
		}
		if err := db.DeleteUs(database, id); err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Пользователь удалён"))
	}
	w.WriteHeader(http.StatusBadGateway)
}

/*
Обработчик обрабатывает UPDATE запрос и обновляет возраст пользователя. ID пользователя указан в запросе. Новый возраст
лежит в в теле запроса.
*/
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		spl := strings.Split(r.RequestURI, "/")
		userId := spl[1]
		id, err := strconv.Atoi(userId)
		if err != nil{
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Некоректно введён ID"))
			return
		}
		if !db.Check_exists_user(database, id){
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Пользователя с  ID " +strconv.Itoa(id)+ " нет"))
			return
		}
		textByte, er := ioutil.ReadAll(r.Body)
		if er != nil{
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(er.Error()))
			return
		}
		defer r.Body.Close()
		textInf := make(map[string]int)
		if err := json.Unmarshal(textByte, &textInf); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(er.Error()))
			return
		}
		newAge, ok := textInf["new age"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("не правильно передан json"))
			return
		}
		if err := db.UpdateAge(database, id, newAge); err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Возраст пользователя успешно обновлён"))
	}
	w.WriteHeader(http.StatusBadGateway)
}
