"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
import random

from sqlalchemy import CHAR, TypeDecorator
from sqlalchemy.dialects.postgresql import UUID
import time


class StringUUID(TypeDecorator):
    impl = CHAR
    cache_ok = True

    def process_bind_param(self, value, dialect):
        if value is None:
            return value
        elif dialect.name == 'postgresql':
            return str(value)
        else:
            return value.hex

    def load_dialect_impl(self, dialect):
        if dialect.name == 'postgresql':
            return dialect.type_descriptor(UUID())
        else:
            return dialect.type_descriptor(CHAR(36))

    def process_result_value(self, value, dialect):
        if value is None:
            return value
        return str(value)


def generate_snowflake_id() -> int:
    timestamp = int(time.time() * 1000)  # Current timestamp in milliseconds
    node_id = random.randint(0, 1023)  # Random node identifier
    sequence = random.randint(0, 4095)  # Random sequence number
    snowflake_id = (timestamp << 22) | (node_id << 12) | sequence
    return snowflake_id
