import json
import sys

from detection import find_checkbox_candidates, is_checked
from image import load_image, remove_lines
from job import save_job


def detect_checkboxes(filepath, job_dir=None):
    img, gray, thresh = load_image(filepath)

    clean = remove_lines(thresh)
    checkboxes = find_checkbox_candidates(clean)

    checkboxes.sort(key=lambda c: (c[2], c[1]))

    boxes = [
        {"bbox": [x, y, x + w, y + h], "is_checked": is_checked(gray, x, y, w, h)}
        for _, x, y, w, h in checkboxes
    ]

    if job_dir:
        save_job(job_dir, filepath, img, clean, gray, checkboxes, boxes)
    else:
        print(json.dumps({"boxes": boxes}))


def main():
    if len(sys.argv) < 2:
        print("Usage: python main.py <image_path> [job_dir]")
        sys.exit(1)

    job_dir = sys.argv[2] if len(sys.argv) >= 3 else None
    detect_checkboxes(sys.argv[1], job_dir)


if __name__ == "__main__":
    main()
