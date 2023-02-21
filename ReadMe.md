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

## Installation
`go install github.com/vault-thirteen/SFHS/cmd/server@latest`

## Startup Parameters

### Server
`server.exe <path-to-configuration-file>`  
`server.exe`  

Example:  
`server.exe "settings.txt"`  
`server.exe`  

**Notes**:  
If the path to a configuration file is omitted, the default one is used.  
Default name of the configuration file is `settings.txt`.

## Settings
Format of the settings' file for a server is quite simple. It uses line
breaks as a separator between parameters. Described below are meanings of each 
line.

1. Server's hostname.
2. Server's listen port.
3. Work mode: HTTP or HTTPS.
4. Path to the certificate file for the HTTPS work mode.
5. Path to the key file for the HTTPS work mode.
6. Hostname of the SFRODB database.
7. Main port of the SFRODB database.
8. Auxiliary port of the SFRODB database.
9. Size of the client pool for the SFRODB database.
10. File extension of served files.
11. MIME type of served files.
12. TTL of served files, i.e. value of the `max-age` field of the 
`Cache-Control` HTTP header.
13. Allowed origin for HTTP CORS, i.e. value of the 
`Access-Control-Allow-Origin` HTTP header.

**Notes**:
* File extension here is used as a normal extension with a dot (period) prefix, 
because Go language uses such format for file extensions. This is not good, but 
this is how Golang works.

## Performance

Performance test of the combination of **SFHS** together with **SFRODB** made 
in Apache JMeter may be found in the `test` folder. Quite a decent hardware 
shows about 22 kRPS in HTTPS mode and about 23 kRPS in HTTP mode, while test 
file size was about 1kB.
