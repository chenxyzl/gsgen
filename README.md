## gen getter setter,mongo incremental update
### wanted
- 生成getter 可选生成setter(为了给后面table和config做readonly)
- 生成mongo的marshal/unmarshal
- 生成setter的脏标记,实现mongo的增量更新

### todo
- [ ] 每个set都需要判断是否 value已经有父节点了
- [ ] MMap的k不为uint64时候,更新脏标记的编译错误
- [ ] mongo的bson.tag检查
- [ ] mongo序列化反序列化的代码生成
- [ ] mongo增量更新的代码生成
- [ ] 通过命令行来按需生成getter/setter/dirty/mongo
    - [ ] getter用作table和config的只读,一般默认导出
    - [ ] setter用于区分权限
    - [ ] dirty用作与增量更新
    - [ ] mongo用于生成mongo对应的序列化和反序列化