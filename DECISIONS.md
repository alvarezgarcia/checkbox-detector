# Decisions

## Approach

The system is split into two components: a Go HTTP API and a Python detection engine. Go handles the request lifecycle, file management, and response formatting. Python handles the image processing using OpenCV.

The detection pipeline works as follows:

1. Convert the image to grayscale and binarize it
2. Remove horizontal and vertical lines, which tend to produce false positives and make processing heavier given the spreadsheet-like layout of most forms
3. Find contours and filter candidates to identify checkbox-like squares
4. Apply a median-based size filter to discard outliers, since real checkboxes in a document are approximately the same size
5. For each candidate, crop the interior and compute a fill ratio to determine if it is checked

Checkboxes are sorted top-left to bottom-right before returning results.

## Tradeoffs

**Classic computer vision vs machine learning**

I have prior familiarity with OpenCV going back to its early days as a C library and later Python bindings, and briefly contributed to porting it as Node.js bindings.
Prbably an ML model would likely be more robust across different document styles, but classic CV requires no training data, no model serving infrastructure, and is fully explainable. The tradeoff is that it relies on tuned constants that may not generalize well to all document types.

**Filesystem storage vs database**

Even if this wasn't an original request of the challenge I wanted to give users an option to see and understand the recognition.
Storage is tied to the machine, no external dependencies.
For production use, this should be moved to object storage.

**Subprocess invocation vs offloading to Lambda**

The Python script is invoked as a subprocess per request.
This is simple and keeps the system self-contained, but it has startup cost. Under high concurrency will not behave properly.
An alternative would be to offload detection to a Lambda function, but wasn't done here in order to avoid infra complexity.

**Synchronous API vs async queue + webhook**

The current design is synchronous: the client uploads an image, waits while detection runs, and receives results in the same HTTP response. This was the original request and works well for fast images but ties up a connection for the duration of processing.
A production design would likely be async: the POST returns a job ID immediately, detection runs in a background worker consuming from a queue, and the client either polls the GET endpoint or provides a webhook URL to be notified when complete.
The debug job ID and GET endpoints already hint at this shape: the main missing piece is decoupling detection from the request lifecycle.

## Known limitations

- No authentication on any endpoint
- The detection constants (area bounds, aspect ratio, fill ratio threshold) were tuned for a specific document style and may not work well on documents with different checkbox sizes or layouts
- Scanned images with noise, skew, or low resolution will produce less reliable results
- The path to the Python script is hardcoded as a relative path and assumes the binary is run from the `api/` directory
