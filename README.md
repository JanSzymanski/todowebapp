# todowebapp

Todo application with web interface utilising todostorelib

## The progress

### API

Basic API server implemented - tested with curl
ToDo: implement concurrency to handle calls from multiple sources

### APIPOKER

To be implemented
Purpose: to 'fire' API request at the server for testing purposes

### MIDDLEWARE APP - http client

To be implemented
Purpose: handling webpage requests and communicating with remote API server and hosting http/template

### HTML tebsite

Simple HTML interface yet to be implemented

## Basic CURL commands to test the API

### GET

```bash
curl http://localhost:8080/todos
curl http://localhost:8080/todo/1
```

### POST

```bash
curl -X POST http://localhost:8080/addtodo -H "Content-Type: text/plain" -d "some todo message"
```

### PATCH

```bash
curl -X PATCH http://localhost:8080/chmtodo -H "Content-Type: application/json" -d '{"Id": 1, "Msg": "some new todo message"}'
curl -X PATCH http://localhost:8080/chstodo -H "Content-Type: application/json" -d '{"Id": 1, "Stat": "active"}'
```

### DELETE

```bash
curl -X DELETE http://localhost:8080/deltodo/1
```
