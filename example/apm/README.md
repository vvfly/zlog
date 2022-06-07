自定义日志输出到apm, 默认输出error级别日志。

### 使用方法
- 申请项目的apm资源，导入相应环境变量
- 启动filebeat
- 执行main.go
```shell
go run main.go
```
- 调用api
```shell
curl http://localhost:8080/zlog
```