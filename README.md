# LARVIS Poker
This project provides an image uploader service for receipt scanning, enabling users to upload photos of their receipts, store them in different resolutions, and download them in an appropriate size. The service functions using a REST API and stores images in the local filesystem.


## Quick start:
Clone the repository:

```bash

git clone https://github.com/yourusername/receipt-uploader.git
cd receipt-uploader

go mod download

go run main.go
```

## Data validatation:


## Image Resizing:
Images will be resized to fit standard resolutions:
- Small (e.g., 640x480)
- Medium (e.g., 1280x720)
- Large (e.g., 1920x1080)
The images will be proportionally scaled according to their original dimensions to maintain aspect ratio.

## API endpoints:

`POST /api/receipts`

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
    "message": "Receipt uploaded successfully",
    "receiptId": "string",
    "imageUrl": "string" // URL to access the uploaded image
}
```

400 Bad Request: Invalid image format
```bash
{
    "error": "Invalid image format"
}
```

500 Internal Server Error: Unexpected error
```bash
{
    "error": "An unexpected error occurred"
}
```

`GET /api/receipts/{receiptId}?resolution=small|medium|large`
Response:
200 OK: Retrieved successfully, returns the image with the specified or default resolution. If query parameter `resolution` is not specified, default value `small` is used

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


## Testing

## Project structure
