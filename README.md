# tiktok

【抖声】【极简版抖音】

### 架构

采用微服务架构，由hertz接收http请求，转交给kitex的rpc客户端处理，rpc客户端去etcd寻找对应的rpc服务，rpc服务调用dao层的数据库处理

### 技术栈

#### hertz

高性能HTTP框架，基于多路复用和协程处理

#### kitex

高性能微服务rpc框架

##### 为什么需要微服务

高内聚低耦合，易开发、易维护、易扩展

构建分布式服务，高并发

进程隔离提高容错

微服务可以使用不同的编程语言开发，根据需要使用不同的数据库

rpc客户端服务端分离，客户端可以对请求做预处理，减轻服务端压力

#### thrift

Facebook开发的rpc通信协议，通过IDL（Interface Definition Language，接口定义语言）规定rpc调用通信的接口和数据类型，支持二进制序列化和其他协议序列化，可以生成rpc代码框架

#### protobuf

Google开发的二进制序列化机制

#### etcd

分布式的kv数据库，用于服务注册与发现

etcd 中的 etc 取自 unix 系统的/etc目录，再加上一个d代表分布式就组成了 etcd，在 unix 系统中 /etc 目录用于存储系统的配置数据，单从名字看 etcd 可用于存储分布式系统的配置数据，有时候也把 etcd 简单理解为分布式 /etc 配置目录

##### 为什么需要服务注册与发现呢

服务地址可以动态变化，不需要硬编码

可以将请求均匀的分配到各个服务上，负载均衡

记录服务状态，容错

#### jwt

JSON Web Token，签名加密的json对象

##### 为什么需要jwt

安全，鉴权，登录状态过期

### 问题

#### 视频点赞不更新点赞数

视频流获取包含了视频播放地址和视频信息，点赞不会改变已经获取的视频信息，也不能重新拉取视频流（视频流会更新）

视频信息单独请求获取，点赞后重新请求视频信息

### 资料

[Kitex | CloudWeGo](https://www.cloudwego.io/zh/docs/kitex)

[Hertz | CloudWeGo](https://www.cloudwego.io/zh/docs/hertz/)
