"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
from typing import Optional
from semantic_router.encoders import BedrockEncoder


class AWSBedrockEncoder(BedrockEncoder):

    """
    The above function is a constructor that initializes an object with a model name and API key.
    
    :param model_name: The `model_name` parameter in the `__init__` method is a string that represents
    the name of the model. It is used to initialize an instance of the class with the specified model
    name
    :type model_name: str
    :param api_key: The `api_key` parameter is a string that represents an authentication key used to
    access an API. It is typically provided by the API provider to authenticate and authorize requests
    made by the client application
    :type api_key: str
    :return: The `super().__init__(name=model_name, api_key=api_key)` statement is returning the
    initialization of the parent class with the `model_name` and `api_key` parameters passed to it.
    """
    def __init__(self,  name: str, access_key_id: Optional[str], secret_access_key: Optional[str], region: Optional[str]):
        return super().__init__(name=name, access_key_id=access_key_id, secret_access_key=secret_access_key, region=region)
    