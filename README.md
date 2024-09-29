# LARVIS Poker
This project provides an image uploader service for receipt scanning, enabling users to upload photos of their receipts, store them in different resolutions, and download them in an appropriate size. The service functions using a REST API and stores images in the local filesystem.


## Quick start:
Clone the repository:

```bash

git clone https://github.com/yourusername/receipt-uploader.git
cd receipt-uploader

go mod download

make start
```

## System design and specifications
All images are stored under `images` folder. 

- Uploading of receipt 
  - Handled by request of `POST /receipts`
  - Each request must have 1 attribute in header: `username_token`
  - System does not check duplicate receipts. In reality this can be handled by `TransactionID` provided by payment service provider
  - Original upload of receipts are stored under `images/config.DIR_UPLOAD`
  - Converted images are stored under `images/config.DIR_GENERATED/username`
  - Each original are converted into 3 different sizes: small, medium and large. large is the original size of uploaded receipt
  - The images will be proportionally scaled according to their original dimensions to maintain aspect ratio.
  - All images are named with uuid without "-" and converted images are suffixed by size, i.e., small, medium, large
- Processing of image
  - Only sending response when storing and resizing receipts are completed. In reality, this can be optimized one job queue to store uploaded receipt and another job queue to resize images. However, it may become confusing for users. Optimization can be implemented depends on use case.
  - Rate limiting: 


- Downloading of receipt 
  - Handled by request of `GET /api/receipts/{receiptId}?size=small|medium|large`
  - Each request must have 1 attribute in header: `username_token`
  - `size` is required, otherwise request will be rejected


## Data validatation:
- Uploading of receipt 
  - Maximum size of upload is 10MB and minimum resolution is 600x800. System rejects any image failed such requirements
  - Only `.jpg` file is supported for the sake of simplicity

- User Token Validation
Each user token can only contain lowercase letters, digits, and underscores. Validation rules for the `user_token` are listed below:

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
| **Edge Cases**             |                    |                                                       |


## Access control:
Access control of generated images are managed by `username_token` and `receiptId`. If an user tries to download someone else's image, `404` response will be sent. `403` is not used for security reason.

## API endpoints:

- `POST /api/receipts`
Request Body:
```bash
{
  "userId": "string",  // optional for permissions system
  "image": "file"      // binary file upload of the receipt image
}
```
Response:
```bash
{
    "receiptId": "string",
}
```

400 Bad Request: Invalid image format
```bash
{
    "error": "Invalid image format"
}

```
403 Forbidden: User does not have permission to access this receipt (if permissions system is implemented)
```bash
{
    "error": "Access denied"
}
```

500 Internal Server Error: Unexpected error
```bash
{
    "error": "internal server error"
}
```

- `GET /api/receipts/{receiptId}?size=small|medium|large`

Response:
200 OK: Retrieved successfully, returns the image with the specified size. 
```bash
{
    Header
}
```

404 Not Found: Receipt does not exist
```bash
{
    "error": "Receipt not found"
}
```

403 Forbidden: User does not have permission to access this receipt (if permissions system is implemented)
```bash
{
    "error": "Access denied"
}
```
## Stretch
- assume each upload of receipt has a transcationID
- assume there is access_token in upload and download

## Testing
- unit test
- integration test
- stress test
## Project structure

## System security
- Simple access control has been implemented to demonstrate the idea that user can only access receipts uploaded by her/himself
- Encryption when uploading image, in reality system should use https for communication
