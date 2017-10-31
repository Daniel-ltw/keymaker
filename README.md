# Keymaker
A Random API key generator/validator. 

This can sit in the private layer without public access, other applications talks to it. 

![Diagram to kind of explain how this could work](https://github.com/Daniel-ltw/Keymaker/blob/master/Keymaker.jpg?raw=true "Steps to how this could work")
Steps:
1. Every response from the web application, before serving back the response, request a key from keymaker
2. Keymaker returns with a key
3. Web application serves up the response with the generated key
4. Client has the API url and the key, makes the request to the API
5. API validates the key with keymaker 
   * if key valid, return the request
   * if key invalid, ignore the request
6. Client gets their response from the API, if key valid

### Go Version
Currently, it is using Go version 1.8

### Application Elements/Dependencies
* air - simple http server, route handler
* beego/orm - Database ORM
* go-sqlite3 - Go sqlite library (this data store could be swap out)
* go-flow - workflow for goroutine

### Application structure
1. keymaker.go - Everything is written in one file, could be refactored
   * init - initialize the application
   * main - the main application, triggers the application to run
   * keygetter - main function for the /get http request
   * keyvalidate - main function for the /validate http request
   * RandStringBytesMaskImprSrc - Random Hash Generator
   * newKeyWorker - workflow goroutine, this function is called with a goroutine and will follow the workflow of active, inactive, killed

### Workflow explanation for newKeyWorker
1. Active - This is when the random key is first generated, it gets stored in the database with the active status
2. Inactive - This is when 30 seconds have passed after it has been active, we disable the key and make it invalid
3. Killed - This is when 60 seconds have passed after it has been active, we remove the key from the database, cleaning things up

## Support
Please contact me with regards to implementation and details
