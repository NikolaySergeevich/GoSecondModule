package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
)

/*
Метод проверяет есть и пользователь с определённым id
*/
func Check_exists_user(database *sql.DB, id int) bool {
	rowsUsers, err := database.Query(fmt.Sprintf("CALL check_exists_user (%s)", strconv.Itoa(id)))
	if err != nil {
		log.Println(err)
	}
	defer rowsUsers.Close()
	res := 0
	for rowsUsers.Next() {
		err := rowsUsers.Scan(&res)
		if err != nil {
			return false
		}
	}
	return res == 1
}

/*
Метод проверяет являются ли друзьями два пользователя
*/
func CheckFriends(database *sql.DB, idSours, idTarget int) bool {
	rowsUsers, err := database.Query(fmt.Sprintf("CALL checkFriends (%s, %s)", strconv.Itoa(idSours), strconv.Itoa(idTarget)))
	if err != nil {
		log.Println(err)
	}
	defer rowsUsers.Close()
	res := 0
	for rowsUsers.Next() {
		err := rowsUsers.Scan(&res)
		if err != nil {
			return false
		}
	}
	return res == 1
}

/*
Метод делает друзьями двух пользователей
*/
func CreateFriendship(database *sql.DB, idSours, idTarget int) error {
	rowsUsers, err := database.Query(fmt.Sprintf("CALL createFriendship (%s, %s)", strconv.Itoa(idSours), strconv.Itoa(idTarget)))
	if err != nil {
		log.Println(err)
	}
	defer rowsUsers.Close()
	if CheckFriends(database, idSours, idTarget) { //проверяет получилось ли подружить
		return nil
	}
	return errors.New("не получилось подружить")
}

/*
Метод удаляет пользователя
*/
func DeleteUs(database *sql.DB, id int) error{
	row, err := database.Query(fmt.Sprintf("CALL deleteUs (%s)", strconv.Itoa(id)))
	if err != nil{
		log.Println(err)
		return errors.New("не смог удалить пользователя")
	}
	defer row.Close()
	return nil
}

/*
Метод обновляет возраст
*/
func UpdateAge(database *sql.DB, id, newAge int) error{
	row, err := database.Query(fmt.Sprintf("CALL updateAge (%s, %s)", strconv.Itoa(id), strconv.Itoa(newAge)))
	if err != nil{
		log.Println(err)
		return errors.New("не смог обновить возраст")
	}
	defer row.Close()
	return nil
}
