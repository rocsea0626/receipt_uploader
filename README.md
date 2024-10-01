# Receipt Uploader
This project provides an image uploader service for receipt scanning, enabling users to upload photos of their receipts, store them in different resolutions/sizes, and download them in an appropriate size. Both original and resized images are stored by the system.


## Quick start
Clone the repository:

```bash

git clone https://github.com/rocsea0626/receipt-uploader.git
cd receipt-uploader

go mod download

make run # run server in release mode
make dev-run # run server in debug mode

# Example request to uplaod a receipt
curl -X POST http://localhost:8080/receipts -F "receipt=@test_image.jpg" -H "username_token: token_guo"

response: 
{
  "receiptId": "4179e13020ad43bab4d8867338f0f048"
}

# Example request to download original image
curl -H "username_token: token_guo" -o download.jpg http://localhost:8080/receipts/{receipId}

# Example request to download different sizes of image
curl -H "username_token: token_guo" -o download-small.jpg http://localhost:8080/receipts/{receipId}\?size\=small
curl -H "username_token: token_guo" -o download-medium.jpg http://localhost:8080/receipts/{receipId}\?size\=medium
curl -H "username_token: token_guo" -o download-large.jpg http://localhost:8080/receipts/{receipId}\?size\=large

```

## Testing

```bash
make test # run unit and integration test
make test-verbose # run unit and integration test in verbose mode
make stress-test # run stress test
```

Temporary files will be created then deleted when tests complete:
- unit test: all unit test cases are defined within each module's folder.
- integration test: defined in `main_test.go` and it starts a server on localhost.
- stress test: defined in `stress_test.go` and multiple requests are send to a localhost server at the same time.



## System design and specifications
All requests must have `username_token` attached in the header. All images are stored under `receipts` folder.

### Uploading of receipt
  - Endpoint: Handled by request of `POST /receipts`
  - Each original upload of receipts is stored under `receipts/config.UPLOADS_DIR/` folder, named as `username#uuid-without-dash.jpg`.
  - A `201` response is immediately sent to client once receipt is stored at server so that user does not need to wait for resizing to complete

### Resizing of image
  - All images are named with uuid without "-" and resized images are suffixed by size, i.e., `4179e13020ad43bab4d8867338f0f048_small.jpg` and stored under `receipts/config.DIR_RESIZED/{username}` folder
  - Each original receipt is converted into 3 different sizes: small, medium and large.
  - Resized images are proportionally scaled to maintain original aspect ratio.
  - To prevent server being overwhelmed by large number of requests, a worker goroutine keeps running continuously in background to scan `receipts/config.UPLOADS_DIR` folder with an interval of `config.Interval` and resize the first uploaded receipt found in that folder. Currently, `config.Interval=1` which means it can resize 1 image per second maximum. However, this approach does not garauntee FIFO order of uploaded images.
  - Once resizing completes, original receipt will be moved to `receipts/config.DIR_RESIZED/{username}` folder
  
### Downloading of receipt 
- To get images with different size: `GET /api/receipts/{receiptId}?size=small|medium|large`
- To get image with original size: `GET /api/receipts/{receiptId}`

### Error Handling
- If image resizing failed, original uploaded receipt will be kept in `config.DIR_UPLOADS` folder.
- Internal system error messages are hidden from clients. Only standard http error messages defined in `constants` module are sent to clients.
- System should not crash because of any runtime error.

### Access control (Stretch):
- Access control is managed by `username_token` and `receiptId` and they are used to locate images in `receipts/config.DIR_RESIZED/{username}` folder
- If an user tries to download someone else's image, `404` response will be sent. `403` is not used for security reason.


## Data validatation:
### Uploading of receipts:
  - `POST /api/receipts`
  - Maximum size of upload is 10MB and minimum resolution is 600x800
  - Only `.jpg` format is supported for the sake of simplicity
  - Http error codes:
```
| Error code | Error case                                 |
|------------|--------------------------------------------|
| 400        | invalid input, size > 10MB                 |
| 400        | invalid input, width < 600                  |
| 400        | invalid input, height < 800                |
| 400        | invalid image format, format=png           |
| 405        | not allowed method to a endpoint           |
| 500        | internal server error                      |
| 201        | receipt is stored successfully             |
```

### Downloading of receipts:
  - `GET /api/receipts/{receiptId}?size=small|medium|large`
  - query parameter size can only be smalle, medium or large.
  - if no size is provided, original size image will be returned
  - Http error codes:
```
| Error code | Error case                                 |
|------------|--------------------------------------------|
| 400        | invalid query parameter value, ?size=xl    |
| 400        | invalid query parameter key, ?resolution=small |
| 404        | not found, access control failed or not found  |
| 405        | not allowed method to a endpoint           |
| 500        | internal server error                      |
| 200        | success                                    |
| 200        | success, size is empty, ?size               |
| 200        | success, size is empty, ?size=              |
```

### User token
Each user token can only contain lowercase letters, digits, and underscores. Validation scenarios for the `user_token` are listed below:
```
| Validation Check           | Value              | Description                                           |
|----------------------------|--------------------|-------------------------------------------------------|
| **Valid**                  | user_token         | Example of a valid token.                             |
|                            | user123            | Contains letters and digits, valid.                   |
|                            | _username          | Starts with an underscore, still valid.               |
|                            | user_name123       | Contains letters, underscores, and digits, valid.     |
| **Invalid**                | INVALID_TOKEN      | Contains uppercase letters, invalid.                  |
|                            | username!          | Contains a special character (!), invalid.            |
|                            | user name          | Contains a space, invalid.                            |
|                            | user@name          | Contains a special character (@), invalid.           |
```

## Project structure
```
.
├── Makefile
├── README.md
├── internal
│   ├── constants
│   │   └── constants.go
│   ├── handlers
│   │   ├── download_receipt.go
│   │   ├── download_receipt_test.go
│   │   ├── health.go
│   │   ├── upload_receipt.go
│   │   └── upload_receipt_test.go
│   ├── http_utils
│   │   ├── http_utils.go
│   │   └── http_utils_test.go
│   ├── image_worker
│   │   ├── types.go
│   │   ├── worker.go
│   │   └── worker_test.go
│   ├── images
│   │   ├── images.go
│   │   ├── images_test.go
│   │   ├── mock
│   │   │   └── images_mock.go
│   │   └── types.go
│   ├── logging
│   │   └── logging.go
│   ├── middlewares
│   │   ├── auth.go
│   │   └── auth_test.go
│   ├── models
│   │   ├── configs
│   │   │   └── configs.go
│   │   ├── http_requests
│   │   │   └── http_requests.go
│   │   ├── http_responses
│   │   │   └── http_responses.go
│   │   └── image_meta
│   │       ├── image_meta.go
│   │       └── image_meta_test.go
│   ├── test_utils
│   │   └── test_utils.go
│   └── utils
│       └── utils.go
├── main.go
├── main_test.go
├── stress_test.go
└── test_image.jpg
```

- `main.go` is entry point of application
- `main_test.go` defines all integration test cases
- `stress_test.go` defines all stress test cases
- `test_image.jpg` test image used in stress test
- `internal/handlers/` defines logic of a handler for each endpoint
- `internal/http_utils/` utility functions for http request
- `internal/image_worker/` defines logic of image_worker which keeping running in backgroud all the time
- `internal/models/image_meta` a data object contains metainfo of a image file, such as path, username, receiptId
- `internal/utils/` contains definition of utility functions
- `internal/images/` defines logics of image resizing

## Implementation concerns:

This section addresses improvement can be done in real world case

### System security
- Access control: JWT token with OAuth 2.0, user-role based ACM and database can be utilized to achieve more comprehensive access control
- Data encryption: all servers should be configured to use HTTPS for secure communication, and data storage must be encrypted to protect sensitive information.
- Duplicate receipts: system should check if same receipt for same transcation has been uploaded more than once.
### Performance
- Mechanism needs to implemented to prevent system being crashed by large number of requests. `JobQueue` or `EventBus` can be utilized for this.
- Caching, Rate limiting and request throttling should be inplemented as well
- If resizing failed, mechanism need to be implemented to notify client, such as sending an email.
- Use load balancer and containers to allow to scale up horizontally
### API backward compability
- API endpoints need to be versioned