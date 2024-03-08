CREATE DATABASE usersdb;
USE usersdb;
CREATE TABLE users
(id INT AUTO_INCREMENT PRIMARY KEY,
user_name VARCHAR(25) NOT NULL,
age INT NOT NULL);
CREATE TABLE friends
(id INT AUTO_INCREMENT PRIMARY KEY,
source_id INT NOT NULL,
target_id INT NOT NULL,
FOREIGN KEY (source_id) REFERENCES users(id) ON DELETE CASCADE);
-- Создание пользователя
DELIMITER //
CREATE PROCEDURE AddUser (IN NameUs VARCHAR(25), IN Age INT)
BEGIN
	INSERT INTO Users (user_name, age) VALUES (NameUs, Age);
    SELECT id FROM Users  ORDER BY id DESC LIMIT 1;
END//
DELIMITER ;
-- Проверка наличия пользователя по ID
DELIMITER //
CREATE PROCEDURE check_exists_user (IN IdUser VARCHAR(25))
BEGIN
	SELECT EXISTS(SELECT id FROM Users WHERE id = IdUser) AS item;
END//
DELIMITER ;
-- Проверка дружбы между двумя пользователями
DELIMITER //
CREATE PROCEDURE checkFriends (IN IdSoursUser VARCHAR(25), IN IdTargetUser VARCHAR(25))
BEGIN
	SELECT EXISTS(SELECT id FROM friends WHERE source_id = IdSoursUser and IdTargetUser) AS item;
END//
DELIMITER ;
-- Создание дружбы
DELIMITER //
CREATE PROCEDURE createFriendship (IN IdSoursUser VARCHAR(25), IN IdTargetUser VARCHAR(25))
BEGIN
	INSERT INTO friends (source_id, target_id) VALUES (idSoursUser, IdTargetUser),(IdTargetUser, idSoursUser);
END//
DELIMITER ;

-- эта процедура вернёт друзей, по указанному id пользователя
DELIMITER //
CREATE PROCEDURE GivUssFr (IN IdUser INT)
BEGIN
	SELECT users.user_name, users.age FROM users LEFT JOIN friends ON friends.source_id = users.id WHERE target_id = IdUser;
END//
DELIMITER ;
-- Удаление пользоваетеля и всех зависимостей дружбы
DELIMITER //
CREATE PROCEDURE deleteUs (IN IdUser INT)
BEGIN
	DELETE FROM users WHERE id = IdUser;
    DELETE FROM friends WHERE target_id = IdUser;
END//
DELIMITER ;
-- Обновление возраста у пользователя
DELIMITER //
CREATE PROCEDURE updateAge (IN IdUser INT, IN NewAge INT)
BEGIN
	UPDATE users SET age = NewAge WHERE id = IdUser;
END//
DELIMITER ;