"""
HTML Documentation generator module for proto-doc.
Renders Jinja2 template with customizable themes, titles, and options.
"""

import os
import json
from typing import Dict, Any, Optional
from jinja2 import Environment, FileSystemLoader

DEFAULT_CONFIG = {
    "title": "API Documentation",
    "primary_color": "#006aff",
    "primary_hover_color": "#0051c7",
    "font_family": "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
    "theme": "light",
    "logo_url": "",
}


class HTMLGenerator:
    """Generates standalone HTML documentation from parsed JSON data."""

    def __init__(self, config: Optional[Dict[str, Any]] = None):
        self.config = {**DEFAULT_CONFIG, **(config or {})}
        template_dir = os.path.join(os.path.dirname(__file__), "templates")
        self.env = Environment(loader=FileSystemLoader(template_dir))

    def render(self, doc_data: Dict[str, Any]) -> str:
        template = self.env.get_template("doc_template.html")
        doc_data_json = json.dumps(doc_data)
        return template.render(
            doc_data=doc_data,
            doc_data_json=doc_data_json,
            config=self.config
        )
