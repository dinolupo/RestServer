# RestServer

Sample Rest API server based on Go with Asynchronous capabilities.

If you need a local REST/JSON API mock server then maybe this utility is for you.

## Compile 

> Linux / WSL2

```
go build -o RestServer RestServer.go
```

> Windows Cross Compiling on WSL2

```sh
GOOS=windows GOARCH=386 go build -o RestServer.exe RestServer.go
```

## Command Line Options

```sh
./RestServer -help

Usage of ./RestServer:
  -dir string
        Static files directory (default "static")
  -port int
        Listen on defined port. (default 9195)
```

## Sample Asynchronous API Call (no callback)

**Method:**
POST, GET 

**Path:**
 /sleep/{seconds}/{filename}

**Path Parameters**

|       Name                 |      Description          | 
| -------------------------  | ------------------------- | 
| seconds  | seconds to sleep before callback function is called | 
| filename | name of the file in the root directory to serve immediately as response | 

**Request Body**

A JSON object with the following properties:

|       Name                 |      Description          | 
| -------------------------  | ------------------------- | 
| url  | POST callback url to be called after sleep parameter is expired | 
| body | body to send alongwith the callback request | 


**Result**

> POST

Returns immediately the "filename" in the same directory of the executable.
Then after *seconds* sleep is expired, a POST request to *url* with *body* is made.

Sample Run

create a file `sample.json` in the same folder of the executable:

```json
{
    "isValid": false,
    "reason": [
        {"code": "0", "description": "OK"}
    ]
}
```

Then launch the following command:

> sample terminal (same can be done with Postman)

```sh
curl --location --request POST 'http://localhost:9195/sleep/5/sample.json' \
--header 'Content-Type: text/plain' \
--data-raw '{
   "url": "http://dummy.restapiexample.com/api/v1/create",
   "body": 	{"name":"test","salary":"123","age":"23"}
}'
```

> console output (see the 5 seconds delay in the response log)
```log
./RestServer
2020/12/02 23:46:23 Start listening on port :9195...
url    :  http://dummy.restapiexample.com/api/v1/create
body   :  map[age:23 name:test salary:123]
2020/12/02 23:47:04 Waiting 5 seconds.
2020/12/02 23:47:09 Callback is being called:
2020/12/02 23:47:09 URL:
http://dummy.restapiexample.com/api/v1/create
2020/12/02 23:47:09 Body:
{"age":"23","name":"test","salary":"123"}
2020/12/02 23:47:10 {"status":"success","data":{"age":"23","name":"test","salary":"123","id":8530},"message":"Successfully! Record has been added."}
```

### Static Content

Static content are served in the *static* subdirectory (or use command line parameter for different folder)

### Dummy APIs

**Method:**
GET,POST,PUT,DELETE

**Path:**
/api/v1

