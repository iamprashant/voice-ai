"""
Tests for vector_factory.py

Covers:
- P3 fix: Vector._knowledge_document must annotate with the ORM KnowledgeDocument,
  not the proto-generated one from knowledge_api_pb2
- Vector can be instantiated with an ORM KnowledgeDocument mock
"""
import pytest
from unittest.mock import MagicMock

from app.core.rag.datasource.vdb.vector_factory import Vector
from app.models.knowledge_model import KnowledgeDocument as OrmKnowledgeDocument


class TestVectorFactoryKnowledgeDocumentImport:
    """P3: annotation must use ORM model not proto KnowledgeDocument."""

    def test_knowledge_document_annotation_is_orm_model(self):
        annotation = Vector.__annotations__.get("_knowledge_document")
        assert annotation is OrmKnowledgeDocument, (
            "Vector._knowledge_document annotation must be the SQLAlchemy ORM model "
            "from app.models.knowledge_model, not the proto-generated class."
        )

    def test_knowledge_document_is_not_proto_type(self):
        try:
            from app.bridges.artifacts.protos.knowledge_api_pb2 import (
                KnowledgeDocument as ProtoKnowledgeDocument,
            )
        except ImportError:
            pytest.skip("Proto KnowledgeDocument not importable")

        annotation = Vector.__annotations__.get("_knowledge_document")
        assert annotation is not ProtoKnowledgeDocument, (
            "Vector._knowledge_document must use the ORM model, not the proto-generated class. "
            "Fix: import from app.models.knowledge_model instead of knowledge_api_pb2."
        )

    def test_orm_knowledge_document_is_sqlalchemy_model(self):
        """Verify the ORM model has expected SQLAlchemy columns."""
        assert hasattr(OrmKnowledgeDocument, "__tablename__")
        assert OrmKnowledgeDocument.__tablename__ == "knowledge_documents"
        assert hasattr(OrmKnowledgeDocument, "id")
        assert hasattr(OrmKnowledgeDocument, "knowledge_id")


class TestVectorInstantiation:

    def test_vector_stores_knowledge_document(self):
        from app.core.embedding.embedding import Embedder
        from app.core.rag.datasource.vdb.vector_base import BaseVector

        orm_doc = MagicMock(spec=OrmKnowledgeDocument)
        embedder = MagicMock(spec=Embedder)
        processor = MagicMock(spec=BaseVector)

        v = Vector(
            knowledge_document=orm_doc,
            embedder=embedder,
            vector_processor=processor,
        )

        assert v._knowledge_document is orm_doc

    def test_vector_stores_embedder(self):
        from app.core.embedding.embedding import Embedder
        from app.core.rag.datasource.vdb.vector_base import BaseVector

        orm_doc = MagicMock(spec=OrmKnowledgeDocument)
        embedder = MagicMock(spec=Embedder)
        processor = MagicMock(spec=BaseVector)

        v = Vector(
            knowledge_document=orm_doc,
            embedder=embedder,
            vector_processor=processor,
        )

        assert v._embeddings is embedder

    def test_vector_stores_processor(self):
        from app.core.embedding.embedding import Embedder
        from app.core.rag.datasource.vdb.vector_base import BaseVector

        orm_doc = MagicMock(spec=OrmKnowledgeDocument)
        embedder = MagicMock(spec=Embedder)
        processor = MagicMock(spec=BaseVector)

        v = Vector(
            knowledge_document=orm_doc,
            embedder=embedder,
            vector_processor=processor,
        )

        assert v._vector_processor is processor
