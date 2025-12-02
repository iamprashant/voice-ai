"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
from dateutil.parser import parse


def is_date(string, fuzzy=False):
    """
    Return whether the string can be interpreted as a date.
    - https://stackoverflow.com/a/25341965/7120095
    :param string: str, string to check for date
    :param fuzzy: bool, ignore unknown tokens in string if True
    """

    try:
        parse(string, fuzzy=fuzzy)
        return True
    except ValueError:
        return False
