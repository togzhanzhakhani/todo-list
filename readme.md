# Todo List API

## Introduction

This project is a RESTful API for a Todo List microservice built in Go. It provides endpoints for creating, updating, deleting, marking as done, and listing tasks.

## Installation

1. Clone the repository:

```sh
git clone https://github.com/togzhanzhakhani/todo-list.git
cd todo-list
```

2. Build the project:

```sh
make build
```

3. Run the server:

```sh
make run
```

## Deployed on RENDER:

Base URL: https://todo-list.onrender.com

## API Endpoints
### Create a New Task
#### URL: /api/todo-list/tasks
#### Method: POST
#### Content-Type: application/json
### Request Body:

```sh
{
  "title": "Buy a book",
  "activeAt": "2023-08-04"
}
```
### Response:

```sh
{
  "id": "task-id"
}
```

### Update an Existing Task
#### URL: /api/todo-list/tasks/{ID}
#### Method: PUT
#### Content-Type: application/json
### Request Body:

```sh
{
  "title": "Buy a high-performance applications book",
  "activeAt": "2023-08-05"
}
```

### Response: 204 No Content
### Delete a Task
#### URL: /api/todo-list/tasks/{ID}
#### Method: DELETE
### Response: 204 No Content
### Mark a Task as Done
#### URL: /api/todo-list/tasks/{ID}/done
#### Method: PUT
### Response: 204 No Content
### List Tasks by Status
#### URL: /api/todo-list/tasks?status={status}
#### Method: GET
### Response:

```sh
[
  {
    "id": "task-id",
    "title": "Buy a high-performance applications book",
    "activeAt": "2023-08-05"
  },
  ...
]
```
#### Query Parameters:
#### status (optional): active or done. Default is active.
