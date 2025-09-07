# Video Downloader REST API (yt-dlp Wrapper)

A simple Go REST API to download, list, serve, and delete videos using
[yt-dlp](https://github.com/yt-dlp/yt-dlp). Includes a Dockerfile for easy
deployment.

## Features

- **POST /videos**: Download a video from a URL (YouTube, Twitter, Instagram,
  etc.)
- **GET /videos**: List all downloaded videos with their creation date
- **GET /videos/{filename}**: Download/stream a video file by filename
- **DELETE /videos/{filename}**: Delete a specific video file
- **DELETE /videos?before=UNIX_TIMESTAMP**: Delete all videos older than the
  given timestamp

## Requirements

- Go 1.19+
- yt-dlp (installed in Docker container)

## Usage

### Build and Run Locally

1. Install Go and yt-dlp
2. Build:
   ```sh
   go build -o video-downloader main.go
   ```
3. Run:
   ```sh
   ./video-downloader
   ```

### Run with Docker

1. Build the Docker image:
   ```sh
   docker build -t video-downloader .
   ```
2. Run the container:
   ```sh
   docker run -p 8080:8080 -v $(pwd)/videos:/app/videos video-downloader
   ```

## API Endpoints

### Download Video

```
POST /videos
Content-Type: application/json
{
  "url": "<video_url>"
}
```

Response:

```
{
  "filename": "<filename>",
  "description": "<description>"
}
```

### List Videos

```
GET /videos
```

Response:

```
[
  { "filename": "abc123.mp4", "created": "2025-09-07 12:34:56" },
  ...
]
```

### Get Video File

```
GET /videos/{filename}
```

### Delete Video File

```
DELETE /videos/{filename}
```

### Delete Videos Older Than Timestamp

```
DELETE /videos?before=UNIX_TIMESTAMP
```

Response:

```
{
  "deleted": ["abc123.mp4", ...]
}
```

## Notes

- Twitter, X, and Instagram URLs will use the video description; other platforms
  use the title.
- All URLs in descriptions are removed.
- Videos are stored in the `videos` directory.

## License

MIT
