"""
Copyright (c) 2023-2025 RapidaAI
Author: Prashant Srivastav <prashant@rapida.ai>

Licensed under GPL-2.0 with Rapida Additional Terms.
See LICENSE.md for details or contact sales@rapida.ai for commercial use.
"""

from abc import ABC, abstractmethod
from typing import List, Optional, Union

from pydantic import BaseModel, Field

from app.exceptions.authentication_exception import InvalidAuthorizationTokenException


class Account(BaseModel):
    id: int
    name: str
    email: str


class Token(BaseModel):
    id: int
    token: str
    tokenType: str


class OrganizationRole(BaseModel):
    id: int
    organizationId: int
    role: str
    organizationName: str


class ProjectRole(BaseModel):
    id: int
    projectId: int
    role: str
    projectName: str


class User(ABC, BaseModel):

    @abstractmethod
    def user_id(self):
        raise NotImplementedError("illegal authenticated user")

    @abstractmethod
    def project_id(self):
        raise NotImplementedError("illegal authenticated user")

    @abstractmethod
    def organization_id(self):
        raise NotImplementedError("illegal authenticated user")


class AuthenticatedUser(User):
    user: Account
    token: Token
    organizationRole: OrganizationRole
    projectRoles: List[ProjectRole]
    currentProject: Optional[ProjectRole] = Field(None)

    def select_project(self, project_id: str) -> Optional[ProjectRole]:
        for project in self.projectRoles:
            if project.projectId == int(project_id):
                self.currentProject = project
                return project
        return None

    @property
    def user_id(self) -> int:
        return self.user.id

    @property
    def project_id(self) -> Union[int, None]:
        if not self.currentProject:
            return None
        return self.currentProject.projectId

    @property
    def organization_id(self) -> int:
        return self.organizationRole.organizationId


class InternalAuthenticatedUser(User):
    userId: int
    projectId: int
    organizationId: int

    @property
    def user_id(self) -> int:
        return self.userId

    @property
    def project_id(self) -> Union[int, None]:
        return self.projectId

    @property
    def organization_id(self) -> int:
        return self.organizationId


class AnonymousUser(User):
    @property
    def user_id(self) -> int:
        raise InvalidAuthorizationTokenException(
            "anonymous user doen't have any attribute."
        )

    @property
    def project_id(self) -> Union[int, None]:
        raise InvalidAuthorizationTokenException(
            "anonymous user doen't have any attribute."
        )

    @property
    def organization_id(self) -> int:
        raise InvalidAuthorizationTokenException(
            "anonymous user doen't have any attribute."
        )
