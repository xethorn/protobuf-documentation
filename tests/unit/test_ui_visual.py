import os
import pytest
from playwright.sync_api import sync_playwright
from proto_doc.compiler import parse_proto_files
from proto_doc.generator import HTMLGenerator


def test_ui_visual_verification(tmp_path):
    # 1. Parse sample proto and options proto
    doc_data = parse_proto_files(["tests/fixtures/sample.proto", "tests/fixtures/options.proto"], include_paths=["."])

    # 2. Render HTML with custom primary color (#10b981) and custom title
    config = {
        "title": "Square Refunds API",
        "primary_color": "#10b981",
        "theme": "light",
    }
    generator = HTMLGenerator(config=config)
    html_content = generator.render(doc_data)

    html_file = tmp_path / "docs.html"
    html_file.write_text(html_content, encoding="utf-8")

    screenshot_dir = tmp_path / "screenshots"
    screenshot_dir.mkdir(exist_ok=True)

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page(viewport={"width": 1280, "height": 800})

        # Open generated HTML page
        page.goto(f"file://{html_file.absolute()}")

        # Verify page title
        assert "Square Refunds API" in page.title()

        # Verify sidebar elements exist
        assert page.is_visible("text=RefundsService")
        assert page.is_visible("text=RefundPayment")
        assert page.is_visible("text=PaymentRefund")

        # Take light theme screenshot
        light_screenshot = screenshot_dir / "docs_light.png"
        page.screenshot(path=str(light_screenshot), full_page=True)

        # Test theme toggle to dark
        page.click(".theme-toggle")
        page.wait_for_timeout(200)
        dark_screenshot = screenshot_dir / "docs_dark.png"
        page.screenshot(path=str(dark_screenshot), full_page=True)

        # Test Cmd + K search modal popup
        page.keyboard.press("Meta+k")
        page.wait_for_timeout(200)
        assert page.is_visible("#search-modal.open")

        # Type search query "Refund"
        page.fill("#search-input", "Refund")
        page.wait_for_timeout(200)

        search_screenshot = screenshot_dir / "search_modal.png"
        page.screenshot(path=str(search_screenshot))

        browser.close()

    print("Playwright UI verification successful!")
