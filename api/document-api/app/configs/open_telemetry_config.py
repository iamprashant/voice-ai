"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
from typing import List

from pydantic import BaseModel


class OpenTelemetryConfig(BaseModel):
    """
    OpenTelemetry configuration template
    """

    enable: bool

    # Debug
    debug: bool = False

    # flag for ssl/tls
    insecure: bool = True

    # ignore url
    ignore_urls: List[str] = ["/readiness/", "/healthz/"]

    # for setting log attribute
    enable_log_tracing: bool = False

    class Config:

        # For secret key
        case_sensitive = True
        env_file_encoding = "utf-8"
