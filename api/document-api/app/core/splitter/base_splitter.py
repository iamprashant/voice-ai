"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
from typing import List

from pydantic.v1 import BaseModel, Extra


class BaseSplitter(BaseModel):
    class Config:
        extra = Extra.allow

    def __call__(self, doc: str) -> List[str]:
        raise NotImplementedError("Subclasses must implement this method")