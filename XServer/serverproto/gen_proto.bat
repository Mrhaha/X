@echo off
REM 获取当前脚本所在的目录路径
SET basepath=%~dp0
cd /D %basepath%

REM 遍历根目录下的所有子文件夹（这些文件夹包含 .proto 文件）
for /F %%i in ('dir /AD /B') do (
    REM 如果当前文件夹中有 .proto 文件
    if exist "%%i\*.proto" (
        REM 输出正在编译的信息
        ECHO Compiling Proto files in folder "%%i"

        REM 生成 .pb.go 文件到当前文件夹
        .\protoc.exe --go_out=. --go_opt=paths=source_relative --plugin=protoc-gen-go.exe "%%i/*.proto"

        REM 生成 .grpc.pb.go 文件到当前文件夹
        .\protoc.exe --go-grpc_out=. --go-grpc_opt=paths=source_relative --plugin=protoc-gen-go.exe "%%i/*.proto"

        ECHO Compilation finished for "%%i"
    )
)

ECHO All proto files have been compiled.
pause
