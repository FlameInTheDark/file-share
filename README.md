# File Share

Minimalistic file share service

## Used:

* PostgreSQL
* MinIO S3

# Installation

Requires `docker` and `docker-compose`.

Use `make install` to build and run the service.

# API

**GET** `/api/v1/file/:file_id` - Get download url

`curl --location --request GET 'http://localhost:8080/api/v1/file/:file_id'`

```json
{
  "url": "download_url"
}
```

---
**POST** `/api/v1/file` - Get upload URL and `file_id` to download it

`curl --location --request POST 'http://localhost:8080/api/v1/file' \
--form 'name="image.png"'`

Argument | Value
--- | ---
`name` | String - file name `image.png`

```json
{
  "url": "upload_url",
  "id": "file_id"
}
```

---
**GET** `/api/v1/file/:file_id/statistics` - Get file statistics

`curl --location --request GET 'http://localhost:8080/api/v1/file/:file_id/statistics'`

```json
{
  "downloads": 0
}
```
