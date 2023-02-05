:: This script builds the server.
@ECHO OFF

SET build_dir=_build_
SET examples_dir=example
SET server_dir=server
SET settings_file=settings.dat
SET server_starter_script=start-server.bat

MKDIR "%build_dir%"

:: Build the server.
CD "%examples_dir%\%server_dir%"
go build
MOVE "%server_dir%.exe" ".\..\..\%build_dir%\"
CD ".\..\..\"

:: Copy some additional files for the server.
COPY "%examples_dir%\%server_dir%\%settings_file%" "%build_dir%\"
COPY "%examples_dir%\%server_dir%\%server_starter_script%" "%build_dir%\"
