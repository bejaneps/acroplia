# Acroplia test task

In this task, I have to write a CLI tool that will interact with Acroplia API and Web interface. Furthermore, I have to cover it with automated and manual tests.

**Parts to be done:** 
1. Login page
2. Creating textpad
3. Sending a private message

**Parts that are done:**
1. Login page
2. Creating textpad
3. Sending a private message

# Build and Install

0. Run Selenium server in separate terminal: ```java -jar selenium-standalone-server.jar```
1. **When running tests using gradle, you have to kill process of selenium-standalone-server. Using linux: ```ps aux | grep java; kill {process_id}```. Using windows: use task manager.**

## Using go

1. Download and install [golang](https://golang.org/dl/)
2. Build executable: ```go build main.go```
3. Run executable: ```./acroplia --help```
4. Run tests: ```go test -v ./...```

## Using gradle

1. Build executable: ```gradle goBuild```
2. Run executable: ```./acroplia --help```
3. Run tests: 
```
gradle testAll - to run all tests
gradle testLoginAPI - to run tests for login through API
gradle testLoginWeb - to run tests for login through Web
gradle testTextpadAPI - to run tests for creating textpad through API
gradle testTextpadWeb - to run tests for creating textpad through Web
gradle testMessageAPI - to run tests for sending message through API
gradle testMessageWeb - to run tests for sending message through Web (currently server sends 500 Internal Server Error even when message is sent)
``` 

**NOTE:** apparently it takes plenty of time to build executable, I recommend to use native Go way to build executable.

# Project structure

**cmd** - contains cli tool's root code

**config** - contains configuration files, for now just config.toml

**drivers** - contains Selenium drivers, for now Chrome and Firefox are tested

**gradle** - contains gradle related config vars

**internal** - contains code for connecting to API or Web interface

**internal/crud_test** - contains manual tests for API

**internal/services_test** - contains manual tests for Selenium

**tests** - contains automated tests and test cases

# Program flags

1. ```--log, -l {filename|filepath}``` - store program logs in separate file, by default all logs are written to stdout (can be .json)
2. ```--debug, -d``` - turn on debugging mode, will print more info about program's steps
3. ```--email {your_email}``` - email in Acroplia used for login
4. ```--password {your_password}``` - password in Acroplia used for login, **required** for login by: username, email or phone
5. ```--username {your_username}``` - username in Acroplia used for login (through API only)
6. ```--phone {phone_number}``` - phone in Acroplia used for login
7. ```--output, -o {filename|filepath}``` - store response body from API calls in separate file, by default response body is written to stdout
8. ```--selenium-port {number}``` - port on which Selenium standalone server is listening
9.  ```--selenium-browser {browser_name}``` - browser to be used by Selenium
10. ```--selenium-options {string[]...}``` - additional options for Selenium browser (like --headless and etc)
11. ```--help, -h``` - to view additional flags for each command

**NOTE:** in case if you don't prefer using flags, you can add all info to _config/config.toml_ file.
