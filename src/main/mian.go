package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/xid"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	dbMaxConn     = 100 // 数据库最大连接数
	dbIdleConn    = 10  // 数据库最大空闲连接数
	usernameField = "username"
	passwordField = "password"
	fanField      = "fan"
	starField     = "star"
)

var sm sessionManager
var db *sql.DB

func init() {
	sm.sessions = make(map[string]*session)
}

// 初始化数据库
func initDB(host string, user string, password string, database string) {
	source := user + ":" + password + "@tcp(" + host + ":3306)/" + database + "?charset=utf8"
	var err error
	db, err = sql.Open("mysql", source)
	checkErr(err)
	db.SetMaxOpenConns(dbMaxConn)
	db.SetMaxIdleConns(dbIdleConn)
	err = db.Ping()
	checkErr(err)
}

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Illegal arguments: ", strings.Join(os.Args, " "))
		fmt.Println("Usage: main.exe <database_host> <database_username> <database_password> <database_name>")
		os.Exit(-1)
	}
	// 初始化数据库，端口默认3306
	initDB(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
	defer db.Close()
	// http服务
	mux := http.NewServeMux()
	mux.HandleFunc("/test", myHandler)
	mux.HandleFunc("/register", register)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/star", star)
	mux.HandleFunc("/unstar", unStar)
	mux.HandleFunc("/stars", getStarsList)
	mux.HandleFunc("/fans", getFansList)
	// 启动服务
	fmt.Println("Listening with port 4546...")
	http.ListenAndServe(":4546", mux)
}

// session
////////////////////

type sessionManager struct {
	lock     sync.Mutex
	sessions map[string]*session
}

type session struct {
	id         string
	username   string
	expire     time.Duration // 过期时间，单位为s
	lastAccess time.Time
}

func (sm *sessionManager) set(k string, v *session) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	if len(k) == 0 || v == nil {
		fmt.Println("Invalid parameters: ", k, v)
		return
	}
	sm.sessions[k] = v
}

func (sm *sessionManager) get(k string) *session {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	v, ok := sm.sessions[k]
	if ok {
		return v
	}
	return nil
}

func (sm *sessionManager) remove(k string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	delete(sm.sessions, k)
}

// user
////////////////////

type user struct {
	username   string
	password   string
	createTime time.Time
}

// http handler
////////////////////

// just for test
func myHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("we get it"))
}

// 注册
func register(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// 校验参数
	if !validateQuery2(&query, usernameField, passwordField) {
		w.Write([]byte("Query with Illegal fields: " + r.URL.RawQuery + ". Register failed."))
		return
	}
	// 校验用户名称是否重复，用户名和密码的合法性在前端通过正则校验
	username := query[usernameField][0]
	if isUserExists(username) {
		w.Write([]byte("Username already exists. Registration failed."))
		return
	}
	// 注册用户入库
	user := user{username, query[passwordField][0], time.Now()}
	id := addUser(user.username, user.password, user.createTime.Format("2006-01-02 15:04:05"))
	w.Write([]byte("Congratulations! Register successfully~"))
	fmt.Println("User added with id: " + strconv.Itoa(int(id)) + ".")
}

// 登录
func login(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// 校验参数
	if !validateQuery2(&query, usernameField, passwordField) {
		w.Write([]byte("Query with Illegal fields: " + r.URL.RawQuery + ". Login failed."))
		return
	}

	username := query[usernameField][0]
	// 校验用户是否已经登录且未过期
	s := sm.get(username)
	if nil != s {
		// session未过期
		if time.Now().Before(s.lastAccess.Add(s.expire)) {
			w.Write([]byte("Already login."))
			return
		}
		// session过期
		sm.remove(username)
	}
	// 校验用户名
	if !isUserExists(username) {
		w.Write([]byte("Username does not exists. Login failed."))
		return
	}
	// 校验密码
	password := query[passwordField][0]
	if !checkPassword(username, password) {
		w.Write([]byte("Password incorrect. Login failed."))
		return
	}

	// 新建session，默认过期时间30分钟
	uuid := xid.New()
	sm.set(username, &session{uuid.String(), username, 1800 * time.Second, time.Now()})
	w.Write([]byte("Login successfully. Welcome " + username + "!"))
	fmt.Println("User " + username + " login.")
}

// 退出
func logout(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// 校验参数
	if !validateQuery1(&query, usernameField) {
		w.Write([]byte("Query with Illegal fields: " + r.URL.RawQuery + ". Logout cancel."))
		return
	}

	username := query[usernameField][0]
	// 校验用户名
	if !isUserExists(username) {
		w.Write([]byte("Username does not exists. Logout failed."))
		return
	}
	if sm.get(username) == nil {
		w.Write([]byte("User " + username + " did not login. Logout failed."))
		return
	}
	sm.remove(username)
	w.Write([]byte("Logout successfully."))
	fmt.Println("User " + username + " logout.")
}

// 关注
func star(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// 校验参数
	if !validateQuery2(&query, fanField, starField) {
		w.Write([]byte("Query with Illegal fields: " + r.URL.RawQuery + ". Operation failed."))
		return
	}

	fan := query[fanField][0]
	star := query[starField][0]
	// 校验用户是否存在或者是否已经关注
	if !isUserExists(fan) {
		w.Write([]byte("User " + fan + " does not exists. Operation failed."))
		return
	}
	if !isUserExists(star) {
		w.Write([]byte("User " + star + " does not exists. Operation failed."))
		return
	}
	if fan == star {
		w.Write([]byte("User " + fan + " can not be a fan of user " + star + ". Operation failed."))
		return
	}
	if isStarExists(fan, star) {
		w.Write([]byte("User " + fan + " is already a fan of user " + star + ". Operation failed."))
		return
	}
	// 添加关注
	addStar(fan, star)
	w.Write([]byte("User " + fan + " has been a fan of user " + star + " successfully."))
	fmt.Println("User " + fan + " is a fan of user " + star + " now.")
}

// 取关
func unStar(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// 校验参数
	if !validateQuery2(&query, fanField, starField) {
		w.Write([]byte("Query with Illegal fields: " + r.URL.RawQuery + ". Operation failed."))
		return
	}

	fan := query[fanField][0]
	star := query[starField][0]
	// 校验用户是否存在或者不存在关注
	if !isUserExists(fan) {
		w.Write([]byte("User " + fan + " does not exists. Operation failed."))
		return
	}
	if !isUserExists(star) {
		w.Write([]byte("User " + star + " does not exists. Operation failed."))
		return
	}
	if !isStarExists(fan, star) {
		w.Write([]byte("User " + fan + " is not a fan of user " + star + ". Operation failed."))
		return
	}
	// 取消关注
	removeStar(fan, star)
	w.Write([]byte("User " + fan + " is not a fan of user " + star + " anymore."))
	fmt.Println("User " + fan + " is not a fan of user " + star + " anymore.")
}

// 获取关注列表
func getStarsList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// 校验参数
	if !validateQuery1(&query, fanField) {
		w.Write([]byte("Query with Illegal fields: " + r.URL.RawQuery + ". Logout cancel."))
		return
	}
	// 获取关注列表
	fan := query[fanField][0]
	if !isUserExists(fan) {
		w.Write([]byte("User " + fan + " does not exists. Operation failed."))
		return
	}
	stars := *getStars(fan)
	res := strings.Join(stars, " ")
	w.Write([]byte("User " + fan + "'s stars:" + res))
}

// 获取粉丝列表
func getFansList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// 校验参数
	if !validateQuery1(&query, starField) {
		w.Write([]byte("Query with Illegal fields: " + r.URL.RawQuery + ". Logout cancel."))
		return
	}
	// 获取粉丝列表
	star := query[starField][0]
	if !isUserExists(star) {
		w.Write([]byte("User " + star + " does not exists. Operation failed."))
		return
	}
	fans := *getFans(star)
	res := strings.Join(fans, " ")
	w.Write([]byte("User " + star + "'s fans:" + res))
}

// mysql operation
////////////////////

func addUser(userName string, password string, createTime string) int64 {
	sql := "insert users set user_name=?, user_password=?, create_time=?"
	stmt, err := db.Prepare(sql)
	checkErr(err)
	res, err := stmt.Exec(userName, password, createTime)
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	return id
}

// validate username and make sure it is unique
func isUserExists(username string) bool {
	sql := "select * from users where user_name=?"
	rows, err := db.Query(sql, username)
	checkErr(err)
	for rows.Next() {
		return true
	}
	return false
}

func checkPassword(username string, password string) bool {
	sql := "select user_password from users where user_name=?"
	rows, err := db.Query(sql, username)
	checkErr(err)
	for rows.Next() {
		var userPassword string
		err = rows.Scan(&userPassword)
		checkErr(err)
		return userPassword == password
	}
	return false
}

func isStarExists(fan string, star string) bool {
	sql := "select fan_id, star_id from stars where fan_id=(select user_id from users where user_name=?) and star_id=(select user_id from users where user_name=?)"
	rows, err := db.Query(sql, fan, star)
	checkErr(err)
	for rows.Next() {
		return true
	}
	return false
}

func addStar(fan string, star string) {
	sql := "insert into stars (fan_id, star_id) values ((select user_id from users where user_name=?), (select user_id from users where user_name=?))"
	stmt, err := db.Prepare(sql)
	checkErr(err)
	_, err = stmt.Exec(fan, star)
	checkErr(err)
}

func removeStar(fan string, star string) {
	sql := "delete from stars where fan_id=(select user_id from users where user_name=?) and star_id=(select user_id from users where user_name=?)"
	stmt, err := db.Prepare(sql)
	checkErr(err)
	_, err = stmt.Exec(fan, star)
	checkErr(err)
}

func getStars(fan string) *[]string {
	// 获取关注数量
	sql1 := "select count(star_id) from stars where fan_id=(select user_id from users where user_name=?)"
	rows1, err := db.Query(sql1, fan)
	checkErr(err)
	var num int
	for rows1.Next() {
		err = rows1.Scan(&num)
		checkErr(err)
	}
	res := make([]string, num)
	// 获取关注用户
	sql2 := "select user_name from users where user_id=any(select star_id from stars where fan_id=(select user_id from users where user_name=?))"
	rows2, err := db.Query(sql2, fan)
	checkErr(err)
	for rows2.Next() {
		var star string
		err = rows2.Scan(&star)
		checkErr(err)
		res = append(res, star)
	}
	return &res
}

func getFans(star string) *[]string {
	// 获取粉丝数量
	sql1 := "select count(fan_id) from stars where star_id=(select user_id from users where user_name=?)"
	rows1, err := db.Query(sql1, star)
	checkErr(err)
	var num int
	for rows1.Next() {
		err = rows1.Scan(&num)
		checkErr(err)
	}
	res := make([]string, num)
	// 获取粉丝列表
	sql2 := "select user_name from users where user_id=any(select fan_id from stars where star_id=(select user_id from users where user_name=?))"
	rows2, err := db.Query(sql2, star)
	checkErr(err)
	for rows2.Next() {
		var star string
		err = rows2.Scan(&star)
		checkErr(err)
		res = append(res, star)
	}
	return &res
}

// tools
////////////////////

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func contains(mapping map[string][]string, key string) bool {
	_, ok := mapping[key]
	return ok
}

func validateQuery(queryP *url.Values, filed string) bool {
	query := *queryP
	if !contains(query, filed) {
		return false
	} else if v := query[filed]; len(v) != 1 {
		return false
	}
	return true
}

func validateQuery1(queryP *url.Values, filed string) bool {
	query := *queryP
	if len(query) != 1 {
		return false
	}
	return validateQuery(queryP, filed)
}

func validateQuery2(queryP *url.Values, filed1 string, filed2 string) bool {
	query := *queryP
	if len(query) != 2 {
		return false
	}
	return validateQuery(queryP, filed1) && validateQuery(queryP, filed2)
}
