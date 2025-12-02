"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
import logging
from logging.config import dictConfig
from typing import Dict


def get_logger(name: str, log_config: Dict) -> logging.Logger:
    """
    get logger
    :return: logger
    """
    # set logger config
    dictConfig(log_config)
    return logging.getLogger(name)
