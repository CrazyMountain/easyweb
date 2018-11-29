# easyweb 

A simple implementation of web server with GO. 



## Note

- 测试前需要开启数据库服务，并保证配置文件里所配置的数据库存在 

- 把测试文件放到根目录是为了避免初始化过程中配置文件路径的不兼容：`go test -v`命令在执行时也会执行每个包下面的init函数，若测试文件路径与主程序不一致，会导致初始化的过程中因找不到配置文件而报错。
- 已知的不足
  - 暂时没有处理404页面
  - 没有采用token的方式处理登录
  - 没有抽象出interface并通过mock进行测试



## Modules support

- gin: `github.com/gin-gonic/gin`
- gorm: `github.com/jinzhu/gorm`
- go-ini: `github.com/go-ini/ini`
- goconvey: `github.com/smartystreets/goconvey`
- monkey: `bou.ke/monkey`
- mysql: `github.com/go-sql-driver/mysql`



## Api

- Sign up

  ``` http
  Method: POST
  Target: /v1/users
  Body: username=<username>&password=<password>
  ```

- Sign in

  ``` http
  Method: POST
  Target: /v1/sessions
  Body: username=<username>&password=<password>
  ```

- Sign out

  ``` http
  Method: DELETE(with cookie=session_id)
  Target: /v1/sessions
  ```

- Follow

  ``` http
  Method: POST
  Target: /v1/users/<username_of_fan>/follows/<username_of_star>
  ```

- Unfollow

  ``` http
  Method: DELETE
  Target: /v1/users/<username_of_fan>/follows/<username_of_star>
  ```

- Get fans

  ``` http
  Method: GET
  Target: /v1/users/<username_of_fan>/follows/0
  ```

- Get stars

  ``` http
  Method: GET
  Target: /v1/users/<username_of_fan>/follows/1
  ```


## Configuration

This project deals configurations with `go-ini`. The configuration file `${PROJECT_HOME}/confserver.ini` contains 3 parts which applied for `server`, `database` and `gin` , description as follows.

```properties
# configurations of easyweb

[server]
port = 4546 # whb server port
session_expire_duration = 5 # session duration after sign in (unit minute)

[database]
type = mysql		# database type, default mysql
host = localhost	# address of database server, default localhost
port = 3306			# port of database server, default 3306
username = root		# username for sign in
password = password	# password for sign in
database = easyweb	# database name
max_open_conn = 100	# max number of opened connections
max_idle_conn = 10	# max nmber of idle connections

[gin]
mode = debug		# runtime mode of gin
```

## Unit Test

You can use command `go test -v`  to see test results, or  use command `goconvey` to show test results in your browser.

