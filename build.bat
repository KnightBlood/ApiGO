@echo off
SETLOCAL

:: 设置项目名称
set "PROJECT_NAME=apigo"

:: 设置输出目录
set "OUT_DIR=dist"

:: 创建输出目录
echo Creating output directory...
if exist "%OUT_DIR%" (	rd /s /q "%OUT_DIR%"
)
mkdir "%OUT_DIR%"

:: 获取当前系统架构
for /f "tokens=2 delims==" %%a in ('wmic os get osarchitecture /value') do set ARCH=%%a
if "%ARCH%"=="64-bit" set "GOARCH=amd64"
if "%ARCH%"=="32-bit" set "GOARCH=386"

:: 获取Go版本
set "GO_VERSION=go%GOARCH%-%DATE:~0,4%-%DATE:~5,2%-%DATE:~8,2%"

:: 构建二进制文件
echo Building %PROJECT_NAME%...
go build -o "%OUT_DIR%\%PROJECT_NAME%.exe" -ldflags "-s -w -X 'main.Version=%GO_VERSION%'" cmd\apigo\main.go
if errorlevel 1 (
	echo Build failed
	exit /b 1
)

echo Build successful

:: 复制配置文件
echo Copying config files...
mkdir "%OUT_DIR%\config"
xcopy /Y config\*.yaml "%OUT_DIR%\config\" >nul
xcopy /Y config\*.yml "%OUT_DIR%\config\" >nul
xcopy /Y config\*.json "%OUT_DIR%\config\" >nul
xcopy /Y config\*.toml "%OUT_DIR%\config\" >nul

:: 复制测试用例
echo Copying test cases...
mkdir "%OUT_DIR%\testcases"
xcopy /Y /E testcases\*.yaml "%OUT_DIR%\testcases\" >nul
xcopy /Y /E testcases\*.yml "%OUT_DIR%\testcases\" >nul
xcopy /Y /E testcases\*.json "%OUT_DIR%\testcases\" >nul
xcopy /Y /E testcases\*.toml "%OUT_DIR%\testcases\" >nul

:: 创建报告目录
echo Creating report directories...
mkdir "%OUT_DIR%\reports\allure\results"

:: 创建运行脚本
echo Creating run script...
echo @echo off > "%OUT_DIR%\run.bat"
echo echo Running %PROJECT_NAME% tests... >> "%OUT_DIR%\run.bat"
echo echo. >> "%OUT_DIR%\run.bat"
echo set "PROJECT_DIR=%%~dp0%%" >> "%OUT_DIR%\run.bat"
echo %%PROJECT_DIR%%\%PROJECT_NAME%.exe ^^^
	--config "%%PROJECT_DIR%%\config" ^^^
	--tests "%%PROJECT_DIR%%\testcases" ^^^
	--reports "%%PROJECT_DIR%%\reports" ^^^
	--format allure >> "%OUT_DIR%\run.bat"
echo. >> "%OUT_DIR%\run.bat"
echo echo. >> "%OUT_DIR%\run.bat"
echo echo Allure results saved to reports\allure\results" >> "%OUT_DIR%\run.bat"
echo echo To generate HTML report: java -jar allure-commandline/bin/allure.jar generate reports\allure\results -o reports\allure\html --clean >> "%OUT_DIR%\run.bat"
echo. >> "%OUT_DIR%\run.bat"
echo echo Press any key to exit... >> "%OUT_DIR%\run.bat"
pause >nul