# sina-spider
爬取新浪股票首页要闻，自动推送到邮箱，需要配置服务端
## 0.2 版
- 订阅后自动推送（一次 20 条）
## 0.1 版本
- 需要自行配置数据库表，账号密码都存储在 MySQL 中；
- 默认一次推送 20 条消息，可以自行修改；
- 验证码由 Redis 缓存 5 分钟后过期，需要重新发送；
- 使用 go 标准日志库记录错误；
- 服务器与客户端均在 Windows 环境下测试；
- 配置信息在 infomation.go 文件中
## TODO
- 订阅后定点推送
- 更多板块内容（证券、期货、黄金等）