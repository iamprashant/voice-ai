"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
from contextvars import ContextVar
from enum import Enum
from typing import Optional


class ObservabilityType(Enum):
    """
    open-telemetry type
    """

    OPENTELEMETRY = 0

    ELASTICAPM = 1


# using for otel span and trace id
observability_type: ContextVar[Optional[ObservabilityType]] = ContextVar(
    "observability_type", default=None
)
trace_name: ContextVar[Optional[str]] = ContextVar(
    "trace_name", default="app.middlewares.open_telemetry_middleware"
)
meter_name: ContextVar[Optional[str]] = ContextVar(
    "meter_name", default="app.middlewares.open_telemetry_middleware"
)
trace_id: ContextVar[Optional[str]] = ContextVar("trace_id", default=None)
span_id: ContextVar[Optional[str]] = ContextVar("span_id", default=None)
service_name: ContextVar[Optional[str]] = ContextVar(
    "service_name", default="python_service"
)
