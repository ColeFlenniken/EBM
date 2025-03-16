 
@echo off
echo %CD%
if "%~1"=="" (
    echo "Please pass a parameter"
    exit /b
)
IF "%~1"=="view" (
    C:\Users\f8col\OneDrive\Desktop\Projects\EBM\src\src.exe
    exit /b
)
IF "%~1"=="add" (
    echo "JIO"
    IF "%~2"=="" (
        echo Missing bookmark name
        exit /b
    )
    echo %~2,%CD% >> "C:\Users\f8col\OneDrive\Desktop\Projects\EBM\src\bm.csv" 
    exit /b
)
  
set "key=%~1"
set "value="

REM Loop through each line in the CSV file
for /f "tokens=1-2 delims=," %%A in (C:\Users\f8col\OneDrive\Desktop\Projects\EBM\src\bm.csv) do (
    if "%%A"=="%key%" (
        set "value=%%B"
        goto :found
    )
)
goto :nfound

:found
    if defined value (
        cd %value%
    ) else (
        echo Key not found.
    )
:nfound
exit /b
