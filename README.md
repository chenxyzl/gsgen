## gen_tools
- 用go的ast来生成struct对象的方法列表getter,setter,json.Marshal/Unmarshal,bson.Marshal/Unmarshal,String,clone,支持增量mongo的更新

### how to run
- 编译: go build -o ./gen ../tools/gen/main.go
- 执行: ./gen -d="../model" -f=".model.go,.mod.go" -s -b
    - -d表示目录
    - -f表示文件后缀
    - -s表示生成setter
    - -b表示生成bson的序列化和反序列化
    - 支持的命令参考请执行 ./gen -h

### 例子
- 添加model(参考example.model.go)
    - 增加1个xxx.model.go
    - 在文件中新增1个struct
- 生成对应的gen.go[getter,setter,clone,json序列化反序列化]和bson.go[bson序列化反序列化,增量更新]
    - 执行命令生成 ./gen -d="../model" -f=".model.go,.mod.go" -s -b
- 使用参考
    - getter/setter 参考model_test.go的函数getTestC
    - mongodb/bson序列化反序列化支持 参考model_test.go的函数TestMongoLoadSave
    - mongodb/bson增量更新支持 参考model_test.go的函数TestBuildDirty