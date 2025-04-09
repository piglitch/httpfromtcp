# HTTP Server Implementation

A lightweight HTTP/1.1 server built from scratch using Go.

## Features

- HTTP/1.1 compliant server implementation  
- Supports GET requests  
- Proper header parsing and generation  
- TCP socket connection management  
- Error handling with appropriate status codes

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/http-server.git
cd http-server

# Run the server
go run ./cmd/httpserver/main.go
```
## Usage
You can use any HTTP client to make requests. Here's an example using Node.js:
Using Node.js http module:
```bash
const http = require('http');

http.get('http://localhost:42069/hello', (res) => {
  let data = '';

  res.on('data', chunk => {
    data += chunk;
  });

  res.on('end', () => {
    console.log('Response:', data);
  });
}).on('error', err => {
  console.error('Error:', err.message);
});
```

## Then test it using one of the Node.js examples above.

## Expected output:
```bash
Hello, World!
```

## Implementation Details
This server is written in Go using low-level TCP socket programming and manual HTTP/1.1 protocol handling. Key components include:

- TCP connection management

- HTTP request parsing

- Response formatting

- Basic routing system

## Limitations
- Currently only supports GET requests

- No support for HTTPS

- Limited to basic content types (e.g., text/plain, text/html)

## Contributing
Pull requests are welcome! For major changes, please open an issue first to discuss what you'd like to change or improve.


