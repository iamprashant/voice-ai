"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

from typing import List
from pydantic import BaseModel


class IndexDocumentRequest(BaseModel):
    knowledgeDocumentId: List[int]
    knowledgeId: int
