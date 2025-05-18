@echo off
set GOEXE=golang_ram_cpu_info__console.exe
set GUIEXE=golang_ram_cpu_info.exe

echo Tidying modules...
go mod tidy

echo Building console version...
go build -ldflags="-s -w" -buildvcs=false -o %GOEXE%
if %errorlevel% neq 0 (
    echo Failed to build console version
    exit /b %errorlevel%
)

echo Building GUI version (no terminal)...
go build -ldflags="-s -w -H=windowsgui" -buildvcs=false -o %GUIEXE%
if %errorlevel% neq 0 (
    echo Failed to build GUI version
    exit /b %errorlevel%
)

echo Compressing with UPX (if installed)...
where upx >nul 2>nul
if %errorlevel%==0 (
    upx --best --lzma %GOEXE%
    upx --best --lzma %GUIEXE%
) else (
    echo UPX not found, skipping compression.
)

echo Done!
dir %GOEXE%
dir %GUIEXE%
