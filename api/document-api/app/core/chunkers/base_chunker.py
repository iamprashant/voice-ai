"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

from abc import ABC, abstractmethod
from typing import Any, List

from pydantic.v1 import Extra
from semantic_router.encoders.base import BaseEncoder

from app.core.chunkers import Chunk
from app.core.splitter import BaseSplitter


class BaseChunker(ABC):
    name: str
    encoder: BaseEncoder
    splitter: BaseSplitter

    class Config:
        extra = Extra.allow

    def __init__(self, name, encoder, splitter):
        self.name = name
        self.encoder = encoder
        self.splitter = splitter

    def __call__(self, docs: List[str]) -> List[List[Chunk]]:
        raise NotImplementedError("Subclasses must implement this method")

    def _split(self, doc: str) -> List[str]:
        return self.splitter(doc)

    @abstractmethod
    def _chunk(self, splits: List[Any]) -> List[Chunk]:
        raise NotImplementedError("Subclasses must implement this method")
