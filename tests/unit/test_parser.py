"""
Unit tests for Project 1: Protobuf to JSON parser.
"""

import os
import json
import pytest
from proto_doc.compiler import parse_proto_files


def test_parse_sample_proto():
    sample_file = "tests/fixtures/sample.proto"
    result = parse_proto_files([sample_file], include_paths=["."])

    assert "files" in result
    assert result["version"] == "1.0"

    # Locate sample.proto in files list
    sample_doc = None
    for f in result["files"]:
        if f["name"] == sample_file:
            sample_doc = f
            break

    assert sample_doc is not None
    assert sample_doc["package"] == "square.payments"
    assert sample_doc["syntax"] == "proto3"

    # Verify Services
    assert len(sample_doc["services"]) == 1
    service = sample_doc["services"][0]
    assert service["name"] == "RefundsService"
    assert "Service managing payment refunds" in service["comments"]["description"]

    # Verify Methods
    assert len(service["methods"]) == 2
    method_refund = service["methods"][0]
    assert method_refund["name"] == "RefundPayment"
    assert method_refund["input_type"] == "square.payments.RefundPaymentRequest"
    assert method_refund["output_type"] == "square.payments.RefundPaymentResponse"
    assert "Refunds a payment." in method_refund["comments"]["description"]

    # Verify Messages
    message_names = [m["name"] for m in sample_doc["messages"]]
    assert "PaymentRefund" in message_names
    assert "RefundPaymentRequest" in message_names
    assert "RefundPaymentResponse" in message_names

    # Find PaymentRefund message
    payment_refund = next(m for m in sample_doc["messages"] if m["name"] == "PaymentRefund")
    assert payment_refund["full_name"] == "square.payments.PaymentRefund"
    assert "PaymentRefund represents a refund of a payment." in payment_refund["comments"]["description"]

    # Check fields of PaymentRefund
    fields = payment_refund["fields"]
    id_field = next(f for f in fields if f["name"] == "id")
    assert id_field["type"] == "string"
    assert id_field["number"] == 1
    assert "Unique identifier for the refund." in id_field["comments"]["description"]

    # Check oneof membership
    card_last_four = next(f for f in fields if f["name"] == "card_last_four")
    assert card_last_four["oneof_name"] == "destination"

    # Check Enums
    assert len(sample_doc["enums"]) == 1
    status_enum = sample_doc["enums"][0]
    assert status_enum["name"] == "PaymentRefundStatus"
    assert len(status_enum["values"]) == 4
    pending_val = next(v for v in status_enum["values"] if v["name"] == "PENDING")
    assert pending_val["number"] == 1
    assert "Refund request is currently pending processing." in pending_val["comments"]["description"]


def test_parse_options_extensions():
    options_file = "tests/fixtures/options.proto"
    result = parse_proto_files([options_file], include_paths=["."])

    options_doc = next(f for f in result["files"] if f["name"] == options_file)
    assert options_doc["package"] == "square.options"

    # Verify extensions
    extensions = options_doc["extensions"]
    assert len(extensions) == 2

    is_experimental = next(e for e in extensions if e["name"] == "is_experimental")
    assert is_experimental["extendee"] == "google.protobuf.MethodOptions"
    assert is_experimental["type"] == "bool"
    assert "Description comment for experimental extension" in is_experimental["comments"]["description"]
