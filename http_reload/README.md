# http_reload

...

## Usage

```
go run example/server.go

curl 'http://localhost:8080/?duration=20s'

kill -USR2 [pid]

curl 'http://localhost:8080/?duration=0s'
```
