## 模仿自 groupcache 的分布式缓存
支持多种算法（LFU LRU ARC）
使用 grpc 用作节点间通信
## 待完成
- 测试编写
- 淘汰回调函数注册
- 服务发现
- 自定义 grpc 的负载均衡
- 完善注释和测试
- 编写文档
- 完成参考部分
- 添加 TTL 机制
- 服务间安全通信
## 参考项目
- groupcache
- peanutcache
- mcache
