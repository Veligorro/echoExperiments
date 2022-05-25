package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io"
	"log"
	"net/http"
	"time"
)

type User struct {
	UserName string `json:"name"`
	UserType string `json:"type"`
	UserId   int    `json:"id"`
}

// different way to get json from Request body
func addUser(c echo.Context) error {
	user := User{}
	defer c.Request().Body.Close()
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Print(err)
		return c.String(http.StatusInternalServerError, "")
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("this is your user %#v", user)
	return c.String(http.StatusOK, "We got your cat")
}

func addUserSec(c echo.Context) error {
	user := User{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
		return c.String(http.StatusBadRequest, "что-то пошло не так")

	}
	log.Printf("this is type %T your user2 %#v", user, user)
	return c.String(http.StatusOK, "add douchebag")
}

func addUserThird(c echo.Context) error {
	user := User{}
	err := c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "something wrong")
	}
	fmt.Printf("Value %#v and Type %T\n", user, user)
	return c.JSON(http.StatusOK, user)
}

func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "admin eat bucket of dick")
}

// пример кастомной middleware функции, в данном случае функция устанавливает добавляет дополнительное поле в заголовок
// можно добавлять свои заголовки .Header.Set("myHeader", "myHeaderInform")
func ChangeHeaderData(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderCookie, "OMMMM COOKIE !!!!")
		return h(c)
	}
}

// создание cookie
// для авторизации берем данные из строки запроса
func loginHeader(c echo.Context) error {
	username := c.QueryParam("name")
	password := c.QueryParam("password")
	// тут можно также создать проверку вводимых данных на соответствие с данными в БД
	if username == "admin" && password == "admin1234" {
		cookie := &http.Cookie{}
		cookie.Name = "UserCookieType"
		cookie.Value = "SomethingElse"
		cookie.Expires = time.Now().Add(72 * time.Hour)
		c.SetCookie(cookie)
		return c.String(http.StatusOK, "new COOKIE!!! NOM NOM NOM")
	}
	return c.String(http.StatusUnauthorized, "incorrect password or username")
}

func main() {
	e := echo.New()
	g := e.Group("/admin")       // в строку адреса добавляется значение группы
	cookie := e.Group("/cookie") // группируем запросы

	// показывает кастомный лог запросов к серверу при  любом обращении к адресу содержащему /admin/...
	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		//Format: `[${time_rfc3339} ${status} ${method} ${host} ${path} ]` + "\n",
		Format: `yaml:${time_rfc3339} ${status} ${method} ${host} ${path}` + "\n",
	}))

	// Валидация. для валидации стоит завести БД с данными о пользователях, например таблица (id,LogName, password, rightsGroup)
	// таким образом если LogName == true и связанный с ним password == true, то применяем опеределенную группу прав, например admin или user
	// пароль запрашивается при любом обращении к групе запросов в данном случае логин и пароль определен самой функцией
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "admin" && password == "admin1234" {
			return true, nil
		}
		return false, fmt.Errorf("incorrect password or username")
	}))
	g.Use(ChangeHeaderData)
	cookie.GET("/set_Cookie", loginHeader) // localhost:8080/coockie/set_Cookie
	g.GET("/main", mainAdmin)              // admin/main
	e.POST("/user", addUser)
	e.POST("/user3/", addUserThird)
	e.POST("/user2/", addUserSec)
	//e.GET("/cook", writeCookie)
	e.Start(":8080")
}
