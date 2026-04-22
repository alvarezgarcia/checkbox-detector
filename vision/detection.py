import cv2
import numpy as np

MIN_BBOX_AREA = 200
MAX_BBOX_AREA = 10000
MIN_ASPECT_RATIO = 0.8
MAX_ASPECT_RATIO = 1.3
MIN_POLYGON_SIDES = 4
MAX_POLYGON_SIDES = 6
SIZE_MEDIAN_TOLERANCE = 0.3
CHECKED_FILL_RATIO_THRESHOLD = 0.15
CHECKED_INTERIOR_MARGIN = 0.2


def find_checkbox_candidates(clean):
    contours, _ = cv2.findContours(clean, cv2.RETR_LIST, cv2.CHAIN_APPROX_SIMPLE)

    candidates = []
    for cnt in contours:
        x, y, w, h = cv2.boundingRect(cnt)
        bbox_area = w * h
        aspect_ratio = w / float(h)

        if not (MIN_BBOX_AREA < bbox_area < MAX_BBOX_AREA):
            continue
        if not (MIN_ASPECT_RATIO < aspect_ratio < MAX_ASPECT_RATIO):
            continue

        approx = cv2.approxPolyDP(cnt, 0.02 * cv2.arcLength(cnt, True), True)
        if not (MIN_POLYGON_SIDES <= len(approx) <= MAX_POLYGON_SIDES):
            continue

        candidates.append((cnt, x, y, w, h))

    if not candidates:
        return []

    widths = [w for (_, _, _, w, _) in candidates]
    heights = [h for (_, _, _, _, h) in candidates]
    med_w = np.median(widths)
    med_h = np.median(heights)

    return [
        (cnt, x, y, w, h)
        for cnt, x, y, w, h in candidates
        if abs(w - med_w) < med_w * SIZE_MEDIAN_TOLERANCE and abs(h - med_h) < med_h * SIZE_MEDIAN_TOLERANCE
    ]


def is_checked(gray, x, y, w, h):
    margin = int(min(w, h) * CHECKED_INTERIOR_MARGIN)
    roi = gray[y + margin:y + h - margin, x + margin:x + w - margin]
    if roi.size == 0:
        return False
    _, roi_bin = cv2.threshold(roi, 0, 255, cv2.THRESH_BINARY_INV | cv2.THRESH_OTSU)
    fill_ratio = cv2.countNonZero(roi_bin) / roi_bin.size
    return bool(fill_ratio > CHECKED_FILL_RATIO_THRESHOLD)
