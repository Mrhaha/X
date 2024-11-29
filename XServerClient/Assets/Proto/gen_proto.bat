@echo off
REM 获取当前脚本所在的目录路径
SET basepath=%~dp0
SET targetPath=%basepath%..\..\
cd /D %basepath%

REM 遍历根目录下的所有子文件夹（这些文件夹包含 .proto 文件）
for /F %%i in ('dir /AD /B') do (
    REM 如果当前文件夹中有 .proto 文件
    if exist "%%i\*.proto" (
        REM 输出正在编译的信息
        ECHO Compiling Proto files in folder "%%i"

        REM 生成 .pb.go 文件到当前文件夹
        %targetPath%\Packages\Grpc.Tools.2.46.6\tools\windows_x64\protoc.exe  --csharp_out=%%i  "%%i/*.proto"
    
        REM %targetPath%\Packages\Grpc.Tools.2.46.6\tools\windows_x64\protoc.exe  --grpc_out=%%i   --plugin=protoc-gen-grpc=%targetPath%\Packages\Grpc.Tools.2.46.6\tools\windows_x64\grpc_csharp_plugin.exe   "%%i/*.proto"

        ECHO Compilation finished for "%%i"
    )
)

ECHO All proto files have been compiled.
pause
