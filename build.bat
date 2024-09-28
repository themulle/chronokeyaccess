@echo off
REM Set Go environment variables for cross-compilation targeting Raspberry Pi Zero (ARMv6)
set GOARCH=arm
set GOOS=linux
set GOARM=6
set CGO_ENABLED=0

REM Optional: Add Go binary to the path if not already set
REM set PATH=C:\Go\bin;%PATH%

REM Compile the Go application
go build -o bin/wgreader_rpi .\cmd\wgreader\main.go
go build -o bin/dooropener_rpi .\cmd\dooropener\main.go