@echo off
echo 开始初始化数据库...
echo.

echo 编译数据库初始化程序...
go build -o init_db.exe scripts/init_db.go
if %errorlevel% neq 0 (
    echo 编译失败！
    pause
    exit /b 1
)

echo 运行数据库初始化...
init_db.exe
if %errorlevel% neq 0 (
    echo 数据库初始化失败！
    pause
    exit /b 1
)

echo.
echo 数据库初始化完成！
echo 清理临时文件...
del init_db.exe

echo.
echo 现在可以重启应用程序测试微信登录功能
pause
