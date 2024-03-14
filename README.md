## gen getter setter,mongo incremental update
### wanted
- 生成getter 可选生成setter(为了给后面table和config做readonly)
- 生成mongo的marshal/unmarshal
- 生成setter的脏标记,实现mongo的增量更新

### todo
- [ ] MMap的k不为uint64时候,更新脏标记的编译错误
- [ ] 每个set都需要判断是否 value已经有父节点了
- [ ] mongo的bson.tag检查
- [ ] mongo序列化反序列化的代码生成（已手动编辑验证）
    - [ ] 优化反序列化的filed为嵌套的model时候多了一次序列化反序列化的问题
- [ ] mongo增量更新的代码生成
- [ ] 通过命令行来按需生成getter/setter/dirty/mongo
    - [ ] 增加指令-g,getter用作table和config的只读,一般默认导出
    - [ ] 增加指令-s,setter用于区分权限
    - [ ] 增加指令-m,mongo用于生成mongo对应的序列化和反序列化（依赖-s）
    - [ ] 增加指令-d,dirty用作与增量更新（依赖-s）