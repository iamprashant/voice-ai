"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
from app.utils.date import is_date


def test_is_date():
    """
    Testing is_date
    :return:
    """
    assert is_date("1990-12-1")
    assert is_date("2005/3")
    assert is_date("Jan 19, 1990")
    assert not is_date("today is 2019-03-27")
    assert is_date("today is 2019-03-27", fuzzy=True)
    assert is_date("Monday at 12:01am")
    assert not is_date("xyz_not_a_date")
    assert not is_date("yesterday")
