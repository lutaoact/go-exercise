# 使用说明

```
# 下载protoc的地址 https://github.com/protocolbuffers/protobuf/releases/
# 下载的zip文件包含include目录
mv include/google /usr/local/include

cd tutorial
protoc -I/usr/local/include -I. --go_out=. addressbook.proto

cd ../
go build add_person.go
go build list_people.go

# 增加记录，如果文件不存在，会新建
./add_person addressbook.data

# 读取记录
./list_person addressbook.data
```
