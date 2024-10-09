# tiktok

【抖声短视频】【极简版抖音】用户可以进行注册登录

采用kitex微服务架构，使用hertz接收HTTP请求，由jwt中间件进行鉴权，rpc客户端进行参数检验，rpc服务端调用dal层的gorm进行业务处理

使用etcd进行rpc服务注册与发现，使用minio存储资源提供限时访问url，使用thrift作为rpc通信协议，使用viper读取配置文件，使用bcrypt加密用户密码

## 架构

采用微服务架构，由hertz接收http请求，转交给kitex的rpc客户端处理，rpc客户端去etcd寻找对应的rpc服务，rpc服务调用dao层的数据库处理

## 技术栈

### hertz

高性能HTTP框架，基于多路复用和协程处理

### kitex

高性能微服务rpc框架

#### 为什么需要微服务

高内聚低耦合，易开发、易维护、易扩展

构建分布式服务，高并发

进程隔离提高容错

微服务可以使用不同的编程语言开发，根据需要使用不同的数据库

rpc客户端服务端分离，客户端可以对请求做预处理，减轻服务端压力

### thrift

Facebook开发的rpc通信协议，通过IDL（Interface Definition Language，接口定义语言）规定rpc调用通信的接口和数据类型，支持二进制序列化和其他协议序列化，可以生成rpc代码框架

### protobuf

Google开发的二进制序列化机制

### etcd

分布式的kv数据库，用于服务注册与发现

etcd 中的 etc 取自 unix 系统的/etc目录，再加上一个d代表分布式就组成了 etcd，在 unix 系统中 /etc 目录用于存储系统的配置数据，单从名字看 etcd 可用于存储分布式系统的配置数据，有时候也把 etcd 简单理解为分布式 /etc 配置目录

#### 为什么需要服务注册与发现呢

服务地址可以动态变化，不需要硬编码

可以将请求均匀的分配到各个服务上，负载均衡

记录服务状态，容错

### jwt

JSON Web Token，签名加密的json对象

内容：头部+载荷+签名

头部：加密算法

载荷：用户状态信息

签名：头部和载荷的签名

#### 为什么需要jwt

安全，鉴权，登录状态过期

### minio

存储对象，提供限时访问url

### BCrypt

使用随机盐和工作因子进行blowfish对称加密，用于存储密码

### viper

管理和读取配置文件

### FFmpeg

获取视频封面

### 数据库

#### 读写分离

设置主从数据库，写在主库，读分发从库

在主库上把数据更改记录到二进制日志中（Binary Log）中，这些记录称为二进制日志事件

从库通过IO线程将主库上的日志复制到自己的中继日志（Relay Log）中

从库通过SQL线程读取中继日志中的事件，将其重放到自己数据上

##### 数据一致性

**缓存标记法**

写操作记录业务id，将预估的主从同步时间作为过期时间，读操作判断业务被标记，有就读主库，没有就读从库

### 问题

#### 视频点赞不更新点赞数

视频流获取包含了视频播放地址和视频信息，点赞不会改变已经获取的视频信息，也不能重新拉取视频流（视频流会更新）

视频信息单独请求获取，点赞后重新请求视频信息

#### 如何回滚

业务中的某个步骤失败如何更好的撤销回滚

#### go对象赋值

db是空的

```go
var db *gorm.db
db,err:=open
```

### 资料

[Kitex | CloudWeGo](https://www.cloudwego.io/zh/docs/kitex)

[Hertz | CloudWeGo](https://www.cloudwego.io/zh/docs/hertz/)
