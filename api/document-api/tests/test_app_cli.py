"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
import subprocess
import sys

import pytest

PYTHON_VERSION = float(f"{sys.version_info.major}.{sys.version_info.minor}")


@pytest.mark.skipif(reason="Minor marker differences", condition=PYTHON_VERSION != 3.8)
def test_requirements_txt():
    """Validate that requirements.txt and requirements-dev.txt
    are up2date with Pipefile"""
    temp_output_dir = "tests/temp_output"
    req_test_file_path = f"{temp_output_dir}/test-requirements.txt"
    req_dev_test_file_path = f"{temp_output_dir}/test-requirements-dev.txt"

    subprocess.call(f"mkdir -p {temp_output_dir}", shell=True)
    subprocess.call(
        "pipenv requirements \
                                    > {}".format(
            req_test_file_path
        ),
        shell=True,
    )

    subprocess.call(
        "pipenv requirements --dev-only \
                                    > {}".format(
            req_dev_test_file_path
        ),
        shell=True,
    )

    with open("requirements.txt") as file:
        req_file = file.read()

    with open("requirements-dev.txt") as file:
        req_dev_file = file.read()

    with open(req_test_file_path) as file:
        req_test_file = file.read()

    with open(req_dev_test_file_path) as file:
        req_dev_test_file = file.read()

    subprocess.call(f"rm -rf {temp_output_dir}", shell=True)

    assert req_file == req_test_file

    assert req_dev_file == req_dev_test_file
