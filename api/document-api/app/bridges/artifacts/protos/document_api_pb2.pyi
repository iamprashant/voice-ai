from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class IndexKnowledgeDocumentRequest(_message.Message):
    __slots__ = ("knowledgeId", "knowledgeDocumentId", "indexType")
    KNOWLEDGEID_FIELD_NUMBER: _ClassVar[int]
    KNOWLEDGEDOCUMENTID_FIELD_NUMBER: _ClassVar[int]
    INDEXTYPE_FIELD_NUMBER: _ClassVar[int]
    knowledgeId: int
    knowledgeDocumentId: _containers.RepeatedScalarFieldContainer[int]
    indexType: str
    def __init__(self, knowledgeId: _Optional[int] = ..., knowledgeDocumentId: _Optional[_Iterable[int]] = ..., indexType: _Optional[str] = ...) -> None: ...

class IndexKnowledgeDocumentResponse(_message.Message):
    __slots__ = ("code", "success")
    CODE_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    code: int
    success: bool
    def __init__(self, code: _Optional[int] = ..., success: bool = ...) -> None: ...
