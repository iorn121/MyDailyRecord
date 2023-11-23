@echo off
setlocal

set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0

go build -o lambda main.go
if errorlevel 1 exit /b %errorlevel%

if exist lambda.zip (
    del /f lambda.zip
)

zip -r lambda.zip ./lambda
if errorlevel 1 exit /b %errorlevel%

endlocal