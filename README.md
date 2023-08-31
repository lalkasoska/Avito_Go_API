# Avito Go API

An API for managing segments, user history, and more.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Usage Examples](#usage-examples)

## Features

- Add, delete, and reassign segments
- Retrieve users' segments
- Generate user history reports
- Easy-to-use JSON API

## Getting Started

To get started with the Avito Go API, follow these steps:

1. Clone this repository: `git clone https://github.com/yourusername/avito-go-api.git`
2. Run `docker compose up`
3. Access the API at `http://localhost:8080` in docker interface.

## API Endpoints

| Endpoint             | Description                  |
|----------------------|------------------------------|
| `POST /`             | Add a new segment            |
| `DELETE /`           | Delete a segment             |
| `PUT /`              | Reassign user's segments   |
| `GET /`              | Get user's segments        |
| `GET /get_history`   | Generate a user history file |
| `GET /report`        | Download latest history file |

## Usage Examples

### Adding a Segment

To add a new segment, send a POST request to the following endpoint:

```bash
curl -X POST http://localhost:8080/ -d '{"name": "NewSegment"}'
```
### Deleting a Segment

To delete a segment, send a DELETE request to the following endpoint:

```bash
curl -X DELETE http://localhost:8080/ -d '{"name": "SegmentToDelete"}'
```

### Getting User's Segments

To retrieve segments, send a GET request to the following endpoint:

```bash
curl -X GET http://localhost:8080/ -d '{"userId": 1}'
```

### Reassigning User's Segments
To reassign segments, send a PUT request to the following endpoint:

```bash
curl -X PUT http://localhost:8080/ -d
{
  "segmentsToAdd": ["Segment1", "Segment2"],
  "segmentsToRemove": ["Segment3"],
  "userId": 123
}
```

### Getting User History
To update history report file to a user's history, send a GET request to:
```bash
curl -X GET http://localhost:8080/get_history -d
{
  "userId": 1,
  "year": 2023,
  "month": 8
}
```

### Downloading Latest User History

To download latest history report file , send a GET request to:
```bash
curl http://localhost:8080/report
```
