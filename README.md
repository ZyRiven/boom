# Boom-V1

<div align="center">
	<img src="http://121.36.4.48:39003/develop/Disk/attr/boom1.png" style="width: 200px">
    <p>
        <h1>Boom V1.0</h1>
    </p>
</div>

----------------------------------------------------------------

## 框架简介

* 数据库采用`gorm`，官方文档：https://gorm.io/zh_CN/docs/

## 特征

好用

## 配置文件

config/config.go

## Api接口

需先在Develop/develop.go下定义结构体并初始化，再新建go文件写接口，最后在route/route.go里注册路由

示例：

* 在route/route.go里添加

```go 
RegisterController(Develop.NewController(r.Group("/user")))
```

* 在Develop/develop.go下定义结构体并初始化，

```go     
type Controller struct {
group *duang.RouterGroup
}

func New(g *duang.RouterGroup) Controller {
return Controller{group: g}
}
```

* 在Develop下新建go文件

```go
func (cc Controller) List() {
cc.group.GET("/list", func (c *duang.Context) {
// 检测登陆
token := app.GetRequest(c).Header["Token"]
_, err := duang.GetUser(token)
if err != "" {
c.JSON(http.StatusUnauthorized, duang.H{
"code": 401,
"msg":  err,
})
return
}
c.JSON(200, duang.H{
"code": 200,
"data": result,
})
})
}

```

## 数据库

使用时先引用

``` go
import "gohello/duang"
```

数据库采用`gorm`,官方文档：https://gorm.io/zh_CN/docs/

### 查询单条记录`Pdo_get()`

参数为 **(表名，字段，条件)** ，当查询字段为空时，默认查全部字段

* 示例：

```go 
w := map[string]interface{}{
"id": 869934687064064,
"Limit": 2,
"Order": "id DESC",
}
columns := []string{
"name",
}
result := duang.Pdo_get("user", columns, w)
```

或

``` go
result := duang.Pdo_get("user", []string{"name"}, map[string]interface{}{"id": 1})
```

### 查询多条记录`Pdo_getall()`

当条件有`Limit`并为数组时，表示需要分页。返回值为 **(数据，总页数，总条数)**

* 示例：

``` go
w := map[string]interface{}{
    "id": 869934687064064, 
    "Limit": [2]int{1, 10},  // 第一页，每页10条
    "Order": "id DESC",
}
columns := []string{
    "name",
}
list, pageNum, total := duang.Pdo_getall("user", columns, w)
```

### 插入记录`Pdo_insert()`

参数为 **(表名，数据)** ，返回值为影响的记录数

* 示例：

``` go
data := map[string]interface{}{
    "dept_id": 333
}
add := duang.Pdo_insertall("role_dept", data)
```

或

``` go
data := []map[string]interface{}{
    {"dept_id": 333}, 
    {"dept_id": 444},
}
add := duang.Pdo_insertall("role_dept", data)
```

### 更新记录`Pdo_update()`

参数为 **(表名，数据，条件)** ，返回值为受影响条数

* 示例：

``` go
w := map[string]interface{}{
    "id": 869934687064064,
}
data := map[string]interface{}{
    "dept_id": 111,
}
up := duang.Pdo_update("role_dept", data, w)
```

### 删除记录`Pdo_delete()`

参数为 **(表名，条件)** ，返回值为受影响条数

* 示例：

``` go
del := duang.Pdo_delete("role_dept", map[string]interface{}{"id": 8})
```

### 查询条数`Pdo_count()`

参数为 **(表名，条件)** ，返回值为条数

* 示例：

``` go
cou := duang.Pdo_count("role_dept", map[string]interface{}{"id": 8})
```

## Token

将id作为参数生成token，token有效时间默认为8小时

``` go
//生成token
token, _ := duang.EnToken("id")
//解token
a, _ := DeToken(token)
if a != nil {
    fmt.Println(a.ID)
}
```

自定义token参数：

在duang/duangTools.go中修改结构体

``` go
type MyCustomClaims struct {
	ID string
	jwt.RegisteredClaims
}
```

修改token有效时间：

修改`EnToken`方法中`MyCustomClaims`结构体的`ExpiresAt`值

``` go
claims := MyCustomClaims{
    val, 
    jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)), //当前时间往后8小时
        Issuer: "test" //签名人
    },
}
```

## WebSocket

目录`duang/webSocket.go`，连接websocket后会收到一个id

基本请求参数：

``` js 
{
    "toId":"",              // 接收者的id
    "selfId":""              // 自己的id
}
```

## Modbus

``` go
import "duang/modbus"
```

RTU

``` go
modbus.ModbusRtuMain(1, 3, 0, 4) // 01 03 00 00 00 04 xx xx
```

TCP

``` go
modbus.ModbusTcpMain(1, 6, 0, 3) // [00 01 00 00 00 06]01 06 00 00 00 03
```

--------------------------------------------------------

## 其他

获取当前文件路径

``` go
_, filename, _, _ := runtime.Caller(0)
abPath := path.Dir(filename)
fmt.Println(abPath)
```

缓存

``` go
Data := "hello~~~"
ca := cache2go.Cache("myCache")
ca.Add("userDate", 5*time.Second, Data) // 存在5秒
value, err := ca.Value("userDate")
if err != nil {
	fmt.Println(err)
	return
}
fmt.Println(value)
```
