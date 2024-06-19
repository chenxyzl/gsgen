## gsgen
- 用go的ast来生成struct对象的方法列表getter,setter,json.Marshal/Unmarshal,bson.Marshal/Unmarshal,String,Clone,Clean,支持增量mongo的更新

### 怎么使用/how to run
- 安装gsgen_tools: go install github.com/chenxyzl/gsgen/gsgen_tools@latest
- 执行: gsgen_tools -d="./example/nest" -f=".model.go,.mod.go" -s -b
    - -d表示目录
    - -f表示文件后缀
    - -s表示导出setter[即true:Setter|false:setter]
    - -b表示生成bson的序列化和反序列化
    - 支持的命令参考请执行 ./gsgen_tools -h