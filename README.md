## gen getter setter,mongo incremental update
### wanted
- 生成getter 可选生成setter(为了给后面table和config做readonly)
- 生成mongo的marshal/unmarshal
- 生成setter的脏标记,实现mongo的增量更新

### todo
- [X] MMap的k不为uint64时候,更新脏标记的编译错误
- [ ] AST生成实现Model增加Clone方法的实现(因为同一个Model只能有1个父,方便多个点个用同一个值的情况)
- [ ] AST生成实现每个set都需要判断是否 value已经有父节点了
- [ ] AST生成实现mongo的bson.tag检查
- [ ] AST生成实现mongo序列化反序列化的代码生成（已手动编辑验证）
    - [X] 优化反序列化的filed为嵌套的model时候多了一次序列化反序列化的问题
- [ ] AST生成实现mongo增量更新的代码生成（已手动验证）
- [ ] 通过命令行来按需生成getter/setter/dirty/mongo
    - [ ] 增加指令-g,getter用作table和config的只读,一般默认导出
    - [ ] 增加指令-s,setter用于区分权限
    - [ ] 增加指令-m,mongo用于生成mongo对应的增量更新、序列化、反序列化（依赖-s）