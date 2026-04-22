import json
import os
import sys
import tempfile
from io import StringIO
from unittest.mock import patch

import cv2
import numpy as np
import pytest

sys.path.insert(0, os.path.join(os.path.dirname(__file__), ".."))

from detection import find_checkbox_candidates
from image import load_image
from job import save_job
from main import detect_checkboxes


def make_checkbox_image(size=25, count=1, spacing=50):
    img = np.zeros((200, 200), dtype=np.uint8)
    for i in range(count):
        x = 20 + i * spacing
        cv2.rectangle(img, (x, 20), (x + size, 20 + size), 255, -1)
    return img


class TestFindCheckboxCandidates:
    def test_returns_candidates_for_square_contours(self):
        clean = make_checkbox_image(size=25)
        candidates = find_checkbox_candidates(clean)
        assert len(candidates) == 1

    def test_returns_empty_for_blank_image(self):
        clean = np.zeros((200, 200), dtype=np.uint8)
        candidates = find_checkbox_candidates(clean)
        assert candidates == []

    def test_filters_out_too_small_contours(self):
        clean = make_checkbox_image(size=5)
        candidates = find_checkbox_candidates(clean)
        assert candidates == []

    def test_filters_out_too_large_contours(self):
        clean = make_checkbox_image(size=120)
        candidates = find_checkbox_candidates(clean)
        assert candidates == []


class TestLoadImage:
    def test_raises_when_file_does_not_exist(self):
        with pytest.raises(ValueError, match="Could not load image"):
            load_image("/nonexistent/path/image.jpg")


class TestSaveJob:
    def test_result_json_has_correct_format(self):
        boxes = [{"bbox": [10, 20, 35, 45], "is_checked": True}]
        with tempfile.TemporaryDirectory() as job_dir:
            img = np.zeros((100, 100, 3), dtype=np.uint8)
            gray = np.zeros((100, 100), dtype=np.uint8)
            save_job(job_dir, "image.jpg", img, gray, gray, [], boxes)

            with open(os.path.join(job_dir, "result.json")) as f:
                result = json.load(f)

            assert "boxes" in result
            assert isinstance(result["boxes"], list)
            assert result["boxes"][0]["bbox"] == [10, 20, 35, 45]
            assert result["boxes"][0]["is_checked"] is True


class TestStdout:
    def test_outputs_valid_json_to_stdout(self):
        with tempfile.NamedTemporaryFile(suffix=".jpg", delete=False) as tmp:
            img = np.zeros((100, 100, 3), dtype=np.uint8)
            cv2.imwrite(tmp.name, img)
            try:
                with patch("sys.stdout", new=StringIO()) as mock_stdout:
                    detect_checkboxes(tmp.name)
                    output = json.loads(mock_stdout.getvalue())
                    assert "boxes" in output
                    assert isinstance(output["boxes"], list)
            finally:
                os.unlink(tmp.name)
