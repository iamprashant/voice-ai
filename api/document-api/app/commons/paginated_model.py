"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""
from typing import Any, Dict, Generic, List, Optional, TypeVar

from pydantic import Field, HttpUrl
from pydantic.generics import GenericModel

T = TypeVar("T")


class PaginatedModel(GenericModel, Generic[T]):
    """
    JResponse wrapper for paginated response provide template for paginated data
    - generate next and prev url
    - result count
    - result list
    """

    __term_key: str = "term"

    # over all result count
    count: int = Field(..., title="count")
    # url where the prev and next request will be redirected
    # hidden param used for calculating next and prev url
    paginated_url: HttpUrl = Field(..., title="paginated_url")
    # request params
    # usually a dict with query and offset, limit
    # search param to generate complete url with params
    param: Dict = Field(..., title="param")

    # pagination result will be always list of generic types
    results: Optional[List[T]] = Field(..., title="results")

    @property
    def previous(self) -> Optional[str]:
        """
        Generating previous link for paginated result
        :returns: url for prev page None for first page.
        :rtype: str
        """
        offset: int = self.param.get("offset")
        limit: int = self.param.get("limit")
        query: str = self.param.get(self.__term_key)

        # for prev page
        offset = offset - limit
        #  for prev link offset cannot go to negative
        if offset < 0:
            return None
        return f"{self.paginated_url}?{self.__term_key}={query}&offset={offset}&limit={limit}"

    @property
    def next(self) -> Optional[str]:
        """
        Generating next link for paginated result
        :returns: url for next page None for last page.
        :rtype: str
        """
        offset: int = self.param.get("offset")
        limit: int = self.param.get("limit")
        query: str = self.param.get(self.__term_key)

        offset = offset + limit
        #  for next link offset cannot go to more than the count
        if self.count - offset < 1:
            return None
        return f"{self.paginated_url}?{self.__term_key}={query}&offset={offset}&limit={limit}"

    @classmethod
    def get_properties(cls):
        # work around for serializing properties
        return [
            prop
            for prop in dir(cls)
            if isinstance(getattr(cls, prop), property)
            and prop not in ("__values__", "fields")
        ]

    def dict(self, *args, **kwargs) -> Any:
        # overriding dict method and adding properties to serialize
        self.__dict__.update(
            {prop: getattr(self, prop) for prop in self.get_properties()}
        )
        return super().dict(*args, **kwargs)


class UserListPaginatedModel(PaginatedModel, Generic[T]):
    @property
    def previous(self) -> Optional[str]:
        # no prev link for feed
        return None

    @property
    def next(self) -> Optional[str]:
        """
        Generating next link for paginated result
        :returns: url for next page None for last page.
        :rtype: str
        """
        offset: int = self.param.get("offset")
        limit: int = self.param.get("limit")

        offset = offset + limit
        #  for next link offset cannot go to more than the count
        if self.count - offset < 1:
            return None
        return f"{self.paginated_url}?offset={offset}&limit={limit}"


class SearchableChannelPaginatedModel(PaginatedModel, Generic[T]):
    __term_key: str = "term"

    @property
    def previous(self) -> Optional[str]:
        # no prev link for feed
        return None

    @property
    def next(self) -> Optional[str]:
        """
        Generating next link for paginated result
        :returns: url for next page None for last page.
        :rtype: str
        """
        lomotif_id: int = self.param.get("lomotif_id")
        offset: int = self.param.get("offset")
        limit: int = self.param.get("limit")
        query: str = self.param.get(self.__term_key)
        offset = offset + limit
        #  for next link offset cannot go to more than the count
        if self.count - offset < 1:
            return None
        return f"{self.paginated_url}?{self.__term_key}={query}&lomotif_id={lomotif_id}&offset={offset}&limit={limit}"
