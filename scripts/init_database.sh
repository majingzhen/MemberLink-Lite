#!/bin/bash

echo "开始初始化数据库..."
echo

echo "编译数据库初始化程序..."
go build -o init_db scripts/init_db.go
if [ $? -ne 0 ]; then
    echo "编译失败！"
    exit 1
fi

echo "运行数据库初始化..."
./init_db
if [ $? -ne 0 ]; then
    echo "数据库初始化失败！"
    exit 1
fi

echo
echo "数据库初始化完成！"
echo "清理临时文件..."
rm -f init_db

echo
echo "现在可以重启应用程序测试微信登录功能"
