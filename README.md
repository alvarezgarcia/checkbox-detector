# Checkbox Detector

A REST API that detects and classifies checkboxes in document images.

## Architecture

The system is split into two components:

- **Go API** (`api/`) - HTTP server that handles requests, manages job lifecycle, and invokes the detection engine
- **Python detection engine** (`vision/`) - OpenCV-based script that processes images, finds checkbox candidates, and classifies them as checked or unchecked

The Go server calls the Python script as a subprocess per request, passing the image path and receiving results as JSON. The reasoning behind this architectural decision is documented in `DECISIONS.md`.

## How it works

1. Image is uploaded via `POST /detect`
2. Go saves it to a temp file and invokes the Python script
3. Python preprocesses the image (grayscale, threshold, line removal), finds square contours matching checkbox dimensions, and classifies each as checked/unchecked based on fill ratio
4. Results are returned as JSON with bounding boxes and checked state

## Setup

**Python**
```bash
cd vision
python3 -m venv venv
venv/bin/pip install -r requirements.txt
```

**Go**
```bash
cd api
cp .env.sample .env
go build -o checkbox-detector .
./checkbox-detector
```

## API

Only `POST /detect` was required by the challenge. A `?debug=true` flag was added to extend it: when present, it stores the original, clean, and annotated images on disk and returns a `detection_job_id` to retrieve them later. `GET /detect/{id}` and `GET /detect/{id}/image` were then added to expose the JSON results and the annotated image via API. The other files stored in the job directory (original and clean images) are only accessible via direct filesystem access, intended for hypothetical manual review by an operator.

### `POST /detect`

Upload an image and detect checkboxes.

`?debug=true` stores original, clean, and annotated images in the job directory and returns a `detection_job_id`.

```bash
curl -X POST https://api.yourdomain.com/detect -F "image=@document.jpg"
```

```json
{
  "boxes": [
    { "bbox": [10, 20, 35, 45], "is_checked": true },
    { "bbox": [10, 60, 35, 85], "is_checked": false }
  ]
}
```

With `?debug=true`:

```bash
curl -X POST "https://api.yourdomain.com/detect?debug=true" -F "image=@document.jpg"
```

```json
{
  "boxes": [
    { "bbox": [10, 20, 35, 45], "is_checked": true },
    { "bbox": [10, 60, 35, 85], "is_checked": false }
  ],
  "debug": {
    "detection_job_id": "checkbox-job-1234567890"
  }
}
```

The `detection_job_id` can then be used to retrieve the results and annotated image via the `GET` endpoints.

### `GET /detect/{id}`

Retrieve results for a previous job. Requires `?debug=true` on the original request.

### `GET /detect/{id}/image`

Retrieve the annotated image for a previous job. Requires `?debug=true` on the original request.

## Configuration

`PYTHON_BIN` - Python interpreter to use. Defaults to `python3`.

`BUCKET_DIR` - Directory for storing job files. Defaults to the OS temp dir.

## Running tests

```bash
cd vision
./tests/run_tests.sh
```
