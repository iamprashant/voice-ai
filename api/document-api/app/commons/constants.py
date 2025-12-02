"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
import re
from enum import Enum


class Product(str, Enum):
    """
    Product Enum
    """

    UNKNOWN = "unknown"
    IOS = "ios"
    ANDROID = "android"
    WEB = "web"
    MOTIFY = "motify"
    LUMBERJACK = "lumberjack"


AGENT_RE = re.compile(
    r"^(?P<tier>client|server)\/(?P<product>ios|android|web|motify|lumberjack)\/(?P<version>\d+\.\d+\.\d+(\-r\d+)?)\/(?P<os>.+)$"
)
WEB_USER_AGENTS = {
    "ios": re.compile(r"^mozilla\/5.0 \(ip(hone|od|ad); .+ os (?P<os>\w+) like mac"),
    "android": re.compile(r"^mozilla\/5\.0 \(linux; (?P<product>android) (?P<os>.+);"),
}
LOMOTIF_AGENT_HEADER = "HTTP_X_Lomotif_Agent"
CF_IP_COUNTRY = "HTTP_CF_IPCOUNTRY"
CF_CONNECTING_IP = "HTTP_CF_CONNECTING_IP"
X_COUNTRY_CODE = "HTTP_X_COUNTRY_CODE"
ACCEPT_LANGUAGE = "HTTP_ACCEPT_LANGUAGE"
X_USER_ID = "HTTP_X_USER_ID"
USER_AGENT = "HTTP_USER_AGENT"
