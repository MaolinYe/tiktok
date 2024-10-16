# tiktok

【抖声短视频】【极简版抖音】功能：用户注册登录、推送视频、投稿视频、点赞视频、评论视频、取消点赞、删除评论；查看用户信息：投稿列表、点赞列表；查看视频信息：标题、作者、是否已点赞、点赞数、评论数、评论列表

采用kitex微服务架构，使用hertz接收HTTP请求，由jwt中间件进行鉴权，rpc客户端进行参数检验，rpc服务端调用dal层的gorm进行业务处理

设置MySQL主从数据库，读写分离，写主库，读分发从库，使用redis更新点赞数和评论数，定时同步MySQL

使用消息队列rabbitmq进行流量削峰，将高并发的数据库新增评论点赞记录请求排队处理

使用thrift作为rpc通信协议，使用etcd进行rpc服务注册与发现，使用minio存储资源提供限时访问url，使用viper读取配置文件，使用bcrypt加密用户密码，使用log写入日志

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

### Redis

存储更新频繁的点赞和评论

### rabbitmq

流量削峰

### minio

存储对象，提供限时访问url

### BCrypt

使用随机盐和工作因子进行blowfish对称加密，用于存储密码

### viper

管理和读取配置文件

### FFmpeg

获取视频封面

### 数据库

![1728978368971](image/README/1728978368971.png)

#### 读写分离

设置主从数据库，写在主库，读分发从库

在主库上把数据更改记录到二进制日志中（Binary Log）中，这些记录称为二进制日志事件

从库通过IO线程将主库上的日志复制到自己的中继日志（Relay Log）中

从库通过SQL线程读取中继日志中的事件，将其重放到自己数据上

##### 数据一致性

**缓存标记法**

写操作记录业务id，将预估的主从同步时间作为过期时间，读操作判断业务被标记，有就读主库，没有就读从库

**半同步复制**

主库在提交事务时等待至少一个从库确认已经接收到事务

**GTID**

Global Transaction Identifier，全局事务标识符

为每个事务分配一个唯一的ID，主从复制由原来的追踪日志改为追踪事务

### 问题

#### redis同步MySQL

不订阅过期、不获取全部键

用一个集合把需要同步的视频id装起来，定时任务中遍历集合，如果key不存在了，就从集合中删除，存在的同步MySQL

go没有集合，遍历过程中删除元素有可能导致迭代器失效，使用哈希map代替集合，key用视频id，value为空的结构体

集合同步的时候会对视频id进行删除，而点赞评论的时候会往集合里面添加视频id，此处需要对集合加锁同步，会出现点赞评论需要等待同步完成，高并发下每次点赞评论都加锁不可行

采用redis分库，一个数据库存储点赞数，一个数据库存储评论数，同步的时候用scan批量获取所有有效key出来同步MySQL

#### redis判断键存在后过期

先更新过期时间，更新失败说明不存在，新增

#### redis判断键存在新增键加锁

不加锁会出现多个并发读取数据库初始化数据的问题，如果在判断键存在的时候就加锁会带来性能问题

参考单例模式实现，使用双重校验锁，判断键存在的时候不加锁，键不存在的情况加锁，再次判断键是否存在，不存在就创建

这样如果键存在，那么不会阻塞，如果键不存在，只会创建一次键

#### 如何回滚

业务中的某个步骤失败如何更好的撤销回滚

#### go对象赋值

db是空的

```go
var db *gorm.db
db,err:=open
```

### 启动

启动MySQL主数据库

启动MySQL从数据库

启动etcd

```powershell
.\etcd.exe
```

启动redis数据库

```powershell
.\redis-server.exe
```

启动rabbitmq

```powershell
.\rabbitmq-server.bat
```

配置minio访问地址，启动minio

```shell
.\minio.exe server .\data\
```

启动服务器和各个服务

### 资料

[Kitex | CloudWeGo](https://www.cloudwego.io/zh/docs/kitex)

[Hertz | CloudWeGo](https://www.cloudwego.io/zh/docs/hertz/)
