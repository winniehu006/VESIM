@echo off
setlocal EnableExtensions
cd /d "%~dp0"

if not exist config mkdir config
if not exist data mkdir data
if not exist vesim.exe (
  echo vesim.exe is missing next to this script.
  pause
  exit /b 1
)
if not exist config\config.yaml (
  if not exist config\config.example.yaml (
    echo config\config.example.yaml is missing, cannot create config\config.yaml.
    pause
    exit /b 1
  )
  copy /y config\config.example.yaml config\config.yaml >nul
)

rem Follow server.port from config\config.yaml; fall back to the default.
set "PORT=7575"
for /f "tokens=2 delims=:" %%p in ('findstr /b /c:"  port:" config\config.yaml') do (
  set "PORT=%%p"
  goto :port_found
)
:port_found
set "PORT=%PORT: =%"

echo Starting VESIM on http://localhost:%PORT% ...
rem Give the server a moment to bind before opening the browser. localhost also
rem reaches the IPv6 listener when another program holds 127.0.0.1:%PORT%.
start "" cmd /c "ping -n 4 127.0.0.1 >nul & start http://localhost:%PORT%"
vesim.exe -c config\config.yaml

echo.
echo VESIM stopped. Press any key to close this window.
pause >nul
