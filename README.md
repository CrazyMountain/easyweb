# easyweb
An easy implementation of web server.

## 简易服务端程序使用指南

### 程序启动方式

服务端口为：`localhost:4546`

```shell
# 程序运行需要四个参数：数据库主机地址、数据库用户名、数据库用户密码、数据库名称
> bin/main.exe <database_host> <database_username> <database_password> <database_name>
```

### 数据库

采用MySQL数据库，需要在所选数据库下手动创建`users`表和`stars`表，建表语句如下：

``` mysql
create table `users` (
    `user_id` int not null auto_increment,
    `user_name` varchar(24) not null,
    `user_password` varchar(24) not null,
    `create_time` timestamp not null,
    PRIMARY KEY (`user_id`)
);

create table `stars` (
    `fan_id` int not null,
    `star_id` int not null
);
```

### 功能介绍

#### 注册

``` http
GET http://localhost:4546/register?username=<username>&password=<password>
```

#### 登录

```http
// 登录过期时间为30分钟
GET http://localhost:4546/login?username=<username>&password=<password>
```

#### 退出

```http
GET http://localhost:4546/logout?username=<username>
```

#### 关注

```http
GET http://localhost:4546/star?fan=<username1>&star=<username2>
```

#### 取关

```http
GET http://localhost:4546/unstar?fan=<username1>&star=<username2>
```

#### 获取关注列表

```http
GET http://localhost:4546/stars?fan=<username>
```

#### 获取粉丝列表

```http
GET http://localhost:4546/fans?star=<username>
```

### 已知的不足

- 为方便快捷选用了GET的方式发送请求

- `session`缓存在内存里，可以考虑存到`redis`
- 没有考虑高并发和调度器

