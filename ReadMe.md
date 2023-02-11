# SFHS
## Static Files HTTP Server

An **HTTP** server for serving static files.  
Static files are taken from the **SFRODB** database.  
The server serves files of only one file type (extension).

### HTTP Status Codes
200 – Successful data retrieval.  
400 – Client has requested wrong data (UID is bad or file does not exist).  
500 – Server error has occurred.  

## Architecture
**HTTP** protocol is used for serving incoming requests.  

Static files are taken from an external database, the **SFRODB**.  
The server uses a pool of clients to connect to the **SFRODB** database.  
For more information about it, see the following link:  
https://github.com/vault-thirteen/SFRODB

## Building
Use the `build.bat` script included with the source code.

## Startup Parameters

### Server
`server.exe <path-to-configuration-file>`

Example:  
`server.exe "settings.dat"`

## Settings
Format of the settings' file can be learned by studying the source code.

## Performance

Performance test of the combination of **SFHS** together with **SFRODB** made 
in Apache JMeter may be found in the `test` folder. Quite a decent hardware 
shows about 22 kRPS in HTTPS mode and about 23 kRPS in HTTP mode, while test 
file size was about 1kB.
