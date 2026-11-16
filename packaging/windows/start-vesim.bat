@echo off
setlocal
cd /d "%~dp0"

if not exist config mkdir config
if not exist config\config.yaml (
  copy /y config\config.example.yaml config\config.yaml >nul
)
if not exist data mkdir data

echo Starting VESIM on http://127.0.0.1:7575 ...
start "" "http://127.0.0.1:7575"
vesim.exe -c config\config.yaml

echo.
echo VESIM stopped. Press any key to close this window.
pause >nul
