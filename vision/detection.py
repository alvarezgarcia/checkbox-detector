import cv2
import numpy as np


def find_checkbox_candidates(clean):
    contours, _ = cv2.findContours(clean, cv2.RETR_LIST, cv2.CHAIN_APPROX_SIMPLE)

    candidates = []
    for cnt in contours:
        x, y, w, h = cv2.boundingRect(cnt)
        bbox_area = w * h
        aspect_ratio = w / float(h)

        if not (200 < bbox_area < 10000):
            continue
        if not (0.8 < aspect_ratio < 1.3):
            continue

        approx = cv2.approxPolyDP(cnt, 0.02 * cv2.arcLength(cnt, True), True)
        if not (4 <= len(approx) <= 6):
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
        if abs(w - med_w) < med_w * 0.3 and abs(h - med_h) < med_h * 0.3
    ]


def is_checked(gray, x, y, w, h, threshold=0.15):
    margin = int(min(w, h) * 0.2)
    roi = gray[y + margin:y + h - margin, x + margin:x + w - margin]
    if roi.size == 0:
        return False
    _, roi_bin = cv2.threshold(roi, 0, 255, cv2.THRESH_BINARY_INV | cv2.THRESH_OTSU)
    fill_ratio = cv2.countNonZero(roi_bin) / roi_bin.size
    return bool(fill_ratio > threshold)
