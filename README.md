## gsgen
- 用go的ast来生成struct对象的方法列表getter,setter,json.Marshal/Unmarshal,bson.Marshal/Unmarshal,String,Clone,Clean,支持增量mongo的更新

### 怎么使用/how to run
- 安装gsgen_tools: go install github.com/chenxyzl/gsgen/gsgen_tools@latest
- 执行: gsgen_tools -d="./example/nest" -f=".model.go,.mod.go" -s -b -a="// test head annotations" -i="github.com/chenxyzl/gsgen/example/common,github.com/chenxyzl/gsgen/example/common.Common"
    - -d 表示目录
    - -f 表示文件后缀
    - -s 可选,表示导出setter[即true:Setter|false:setter]
    - -b 可选,表示生成bson的序列化和反序列化
    - -a 可选,数组,表示追加在头部的注释(一般用于给生成的文件添加额外的提示信息等)
    - -i 可选,数组,表示忽略检查的外部包,也可以指定外部包的特定类型(用于导入外部包类型,注:外部包对象的的安全无法保证，增量更新都不支持，建议外部包对象也用本工具生成)
    - 支持的命令参考请执行 ./gsgen_tools -h
- 参考: Makefile