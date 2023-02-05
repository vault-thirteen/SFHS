# SFHS
## Static Files HTTP Server

An HTTP server for serving static files.  
Static files are taken from the SFRODB database, which must be started separately.

### HTTP Status Codes
200 – Successful data retrieval.  
400 – Client has requested wrong data (UID is bad or file does not exist).  
500 – Server error has occurred.  

## Building
Use the `build.bat` script included with the source code.

## Startup Parameters

### Server
`server.exe <path-to-configuration-file>`

Example:  
`server.exe "settings.dat"`

## Settings
Format of the settings' file can be learned by studying the source code.
