import cv2

LINE_REMOVAL_LENGTH = 80


def load_image(filepath):
    img = cv2.imread(filepath)
    if img is None:
        raise ValueError(f"Could not load image: {filepath}")
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    _, thresh = cv2.threshold(gray, 0, 255, cv2.THRESH_BINARY_INV | cv2.THRESH_OTSU)
    return img, gray, thresh


def remove_lines(thresh):
    kernel_h = cv2.getStructuringElement(cv2.MORPH_RECT, (LINE_REMOVAL_LENGTH, 1))
    kernel_v = cv2.getStructuringElement(cv2.MORPH_RECT, (1, LINE_REMOVAL_LENGTH))
    horizontal = cv2.morphologyEx(thresh, cv2.MORPH_OPEN, kernel_h)
    vertical = cv2.morphologyEx(thresh, cv2.MORPH_OPEN, kernel_v)
    clean = cv2.subtract(thresh, horizontal)
    clean = cv2.subtract(clean, vertical)
    return clean
