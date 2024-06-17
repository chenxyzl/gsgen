## gsgen
- 用go的ast来生成struct对象的方法列表getter,setter,json.Marshal/Unmarshal,bson.Marshal/Unmarshal,String,Clone,Clean,支持增量mongo的更新

### how to run
- 编译: go build -o ./bin/gsgen ./gsgen_tools/main.go
- 进入目录: cd ./bin
- 执行: ./gsgen -d="../example/nest" -f=".model.go,.mod.go" -s -b
    - -d表示目录
    - -f表示文件后缀
    - -s表示导出setter[即true:Setter|false:setter]
    - -b表示生成bson的序列化和反序列化
    - 支持的命令参考请执行 ./gen -h

### 例子
- 添加model(参考./example/bson/example_bson.model.go)
    - 增加1个xxx.model.go
    - 在文件中新增1个struct
- 执行生成
    - 执行命令生成 ./gsgen -d="../example/nest" -f=".model.go,.mod.go" -s -b
- 使用参考
    - getter/setter 参考"./test/model_test.go"的函数getTestC
    - mongodb/bson序列化反序列化支持 参考"./test/model_test.go"的函数TestMongoLoadSave
    - mongodb/bson增量更新支持 参考"./test/model_test.go"的函数TestBuildBson