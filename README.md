# Distributed Property Graph

纯Go实现的分布式属性图服务，提供图/Schema管理、顶点边增量写入、批量导入、属性和邻接索引、只读快照、受限遍历查询以及可取消的图算法任务。默认内存模式包含多分片模拟器，不依赖外部数据库即可验证。

## 启动

`HTTP_ADDR=:28090 go run ./cmd/server`

使用`go run ./cmd/simulator -base http://127.0.0.1:28090`执行创建图、发布Schema、写入顶点/边、创建快照、查询路径和算法任务流程。服务提供`/healthz`、`/readyz`和`/metrics`。

目录包含`cmd`、`internal`、`api`、`configs`、`migrations`、`deploy`和`scripts`。各领域按domain/application/adapter/infrastructure分层，存储、分片、时钟、查询解析器、任务调度和节点通信使用接口注入。

非测试Go源码不少于2600行；测试、生成代码、依赖、锁文件、SQL、配置、脚本和构建产物不计入。
