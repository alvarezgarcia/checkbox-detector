import json
import os

import cv2

from detection import is_checked


def save_job(job_dir, filepath, img, clean, gray, checkboxes, boxes):
    _, ext = os.path.splitext(filepath)
    annotated = img.copy()
    for _, x, y, w, h in checkboxes:
        color = (0, 255, 0) if is_checked(gray, x, y, w, h) else (0, 0, 255)
        cv2.rectangle(annotated, (x, y), (x + w, y + h), color, 2)
    cv2.imwrite(os.path.join(job_dir, f"original{ext}"), img)
    cv2.imwrite(os.path.join(job_dir, f"clean{ext}"), clean)
    cv2.imwrite(os.path.join(job_dir, f"annotated{ext}"), annotated)
    with open(os.path.join(job_dir, "result.json"), "w") as f:
        json.dump({"boxes": boxes}, f)
