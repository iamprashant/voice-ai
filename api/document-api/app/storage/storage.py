"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

from collections.abc import Generator
from typing import Union

from app.configs.storage_config import AssetStoreConfig
from app.storage.file_storage.base_storage import BaseStorage
from app.storage.file_storage.local_storage import LocalStorage
from app.storage.file_storage.s3_storage import S3Storage


class Storage:
    storage_runner: BaseStorage

    def __init__(self, config: AssetStoreConfig):
        storage_type = config.storage_type
        if storage_type == "s3":
            self.storage_runner = S3Storage(config=config)
        else:
            self.storage_runner = LocalStorage(config=config)

    def save(self, filename, data):
        self.storage_runner.save(filename, data)

    def load(self, filename: str, stream: bool = False) -> Union[bytes, Generator]:
        if stream:
            return self.load_stream(filename)
        else:
            return self.load_once(filename)

    def load_once(self, filename: str) -> bytes:
        return self.storage_runner.load_once(filename)

    def load_stream(self, filename: str) -> Generator:
        return self.storage_runner.load_stream(filename)

    def download(self, filename, target_filepath):
        self.storage_runner.download(filename, target_filepath)

    def exists(self, filename):
        return self.storage_runner.exists(filename)

    def delete(self, filename):
        return self.storage_runner.delete(filename)

    async def disconnect(self):
        self.storage_runner = None

    async def is_connected(self) -> bool:
        return True
