## gen_tools
- 用go的ast来生成struct对象的方法列表getter,setter,json.Marshal/Unmarshal,bson.Marshal/Unmarshal,String,clone,支持增量mongo的更新

### todo
- [X] MMap的k不为uint64时候,更新脏标记的编译错误
- [X] AST生成实现Model增加Clone方法的实现(因为同一个Model只能有1个父,方便多个点个用同一个值的情况)
- [X] AST生成实现每个set都需要判断是否 value已经有父节点了
- [X] AST生成实现bson的bson.tag检查
- [X] AST生成实现bson序列化反序列化的代码生成
    - [X] 优化反序列化的filed为嵌套的model时候多了一次序列化反序列化的问题
- [X] AST生成实现Model的json序列化和反序列化--类bson
- [X] AST生成实现Model的string
- [X] AST生成实现bson增量更新的代码生成
- [X] 通过命令行来按需生成getter/setter/dirty/bson
    - [X] 增加指令getter用作table和config的只读,默认导出不可配置
    - [X] 增加指令-s,setter用于区分权限
    - [X] 增加指令-b,用于生成bson对应的增量更新、序列化、反序列化（依赖-s）

### how to run
- 编译: go build -o ./gen ../tools/gen/main.go
- 执行: ./gen -d="../model" -f=".model.go,.mod.go" -s -b
    - -d表示目录
    - -f表示文件后缀
    - -s表示生成setter
    - -b表示生成bson的序列化和反序列化
    - 支持的命令参考请执行 ./gen -h
