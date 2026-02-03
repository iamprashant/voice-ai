from google.protobuf import any_pb2 as _any_pb2
import app.bridges.artifacts.protos.common_pb2 as _common_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ConversationToolCall(_message.Message):
    __slots__ = ("id", "toolId", "name", "args", "time")
    class ArgsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: _any_pb2.Any
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...
    ID_FIELD_NUMBER: _ClassVar[int]
    TOOLID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    ARGS_FIELD_NUMBER: _ClassVar[int]
    TIME_FIELD_NUMBER: _ClassVar[int]
    id: str
    toolId: str
    name: str
    args: _containers.MessageMap[str, _any_pb2.Any]
    time: _timestamp_pb2.Timestamp
    def __init__(self, id: _Optional[str] = ..., toolId: _Optional[str] = ..., name: _Optional[str] = ..., args: _Optional[_Mapping[str, _any_pb2.Any]] = ..., time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ...) -> None: ...

class ConversationToolResult(_message.Message):
    __slots__ = ("id", "toolId", "name", "args", "success", "time")
    class ArgsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: _any_pb2.Any
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...
    ID_FIELD_NUMBER: _ClassVar[int]
    TOOLID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    ARGS_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    TIME_FIELD_NUMBER: _ClassVar[int]
    id: str
    toolId: str
    name: str
    args: _containers.MessageMap[str, _any_pb2.Any]
    success: bool
    time: _timestamp_pb2.Timestamp
    def __init__(self, id: _Optional[str] = ..., toolId: _Optional[str] = ..., name: _Optional[str] = ..., args: _Optional[_Mapping[str, _any_pb2.Any]] = ..., success: bool = ..., time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ...) -> None: ...

class ConversationDirective(_message.Message):
    __slots__ = ("id", "type", "args", "time")
    class DirectiveType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        DIRECTIVE_TYPE_UNSPECIFIED: _ClassVar[ConversationDirective.DirectiveType]
        END_CONVERSATION: _ClassVar[ConversationDirective.DirectiveType]
        TRANSFER_CONVERSATION: _ClassVar[ConversationDirective.DirectiveType]
        PAUSE_CONVERSATION: _ClassVar[ConversationDirective.DirectiveType]
        MUTE_CALLER: _ClassVar[ConversationDirective.DirectiveType]
        UNMUTE_CALLER: _ClassVar[ConversationDirective.DirectiveType]
    DIRECTIVE_TYPE_UNSPECIFIED: ConversationDirective.DirectiveType
    END_CONVERSATION: ConversationDirective.DirectiveType
    TRANSFER_CONVERSATION: ConversationDirective.DirectiveType
    PAUSE_CONVERSATION: ConversationDirective.DirectiveType
    MUTE_CALLER: ConversationDirective.DirectiveType
    UNMUTE_CALLER: ConversationDirective.DirectiveType
    class ArgsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: _any_pb2.Any
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...
    ID_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    ARGS_FIELD_NUMBER: _ClassVar[int]
    TIME_FIELD_NUMBER: _ClassVar[int]
    id: str
    type: ConversationDirective.DirectiveType
    args: _containers.MessageMap[str, _any_pb2.Any]
    time: _timestamp_pb2.Timestamp
    def __init__(self, id: _Optional[str] = ..., type: _Optional[_Union[ConversationDirective.DirectiveType, str]] = ..., args: _Optional[_Mapping[str, _any_pb2.Any]] = ..., time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ...) -> None: ...

class ConversationConfiguration(_message.Message):
    __slots__ = ("assistantConversationId", "assistant", "time", "metadata", "args", "options", "inputConfig", "outputConfig")
    class MetadataEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: _any_pb2.Any
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...
    class ArgsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: _any_pb2.Any
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...
    class OptionsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: _any_pb2.Any
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...
    ASSISTANTCONVERSATIONID_FIELD_NUMBER: _ClassVar[int]
    ASSISTANT_FIELD_NUMBER: _ClassVar[int]
    TIME_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    ARGS_FIELD_NUMBER: _ClassVar[int]
    OPTIONS_FIELD_NUMBER: _ClassVar[int]
    INPUTCONFIG_FIELD_NUMBER: _ClassVar[int]
    OUTPUTCONFIG_FIELD_NUMBER: _ClassVar[int]
    assistantConversationId: int
    assistant: _common_pb2.AssistantDefinition
    time: _timestamp_pb2.Timestamp
    metadata: _containers.MessageMap[str, _any_pb2.Any]
    args: _containers.MessageMap[str, _any_pb2.Any]
    options: _containers.MessageMap[str, _any_pb2.Any]
    inputConfig: StreamConfig
    outputConfig: StreamConfig
    def __init__(self, assistantConversationId: _Optional[int] = ..., assistant: _Optional[_Union[_common_pb2.AssistantDefinition, _Mapping]] = ..., time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., metadata: _Optional[_Mapping[str, _any_pb2.Any]] = ..., args: _Optional[_Mapping[str, _any_pb2.Any]] = ..., options: _Optional[_Mapping[str, _any_pb2.Any]] = ..., inputConfig: _Optional[_Union[StreamConfig, _Mapping]] = ..., outputConfig: _Optional[_Union[StreamConfig, _Mapping]] = ...) -> None: ...

class ICEServer(_message.Message):
    __slots__ = ("urls", "username", "credential")
    URLS_FIELD_NUMBER: _ClassVar[int]
    USERNAME_FIELD_NUMBER: _ClassVar[int]
    CREDENTIAL_FIELD_NUMBER: _ClassVar[int]
    urls: _containers.RepeatedScalarFieldContainer[str]
    username: str
    credential: str
    def __init__(self, urls: _Optional[_Iterable[str]] = ..., username: _Optional[str] = ..., credential: _Optional[str] = ...) -> None: ...

class ICECandidate(_message.Message):
    __slots__ = ("candidate", "sdpMid", "sdpMLineIndex", "usernameFragment")
    CANDIDATE_FIELD_NUMBER: _ClassVar[int]
    SDPMID_FIELD_NUMBER: _ClassVar[int]
    SDPMLINEINDEX_FIELD_NUMBER: _ClassVar[int]
    USERNAMEFRAGMENT_FIELD_NUMBER: _ClassVar[int]
    candidate: str
    sdpMid: str
    sdpMLineIndex: int
    usernameFragment: str
    def __init__(self, candidate: _Optional[str] = ..., sdpMid: _Optional[str] = ..., sdpMLineIndex: _Optional[int] = ..., usernameFragment: _Optional[str] = ...) -> None: ...

class WebRTCSDP(_message.Message):
    __slots__ = ("type", "sdp")
    class SDPType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        SDP_TYPE_UNSPECIFIED: _ClassVar[WebRTCSDP.SDPType]
        OFFER: _ClassVar[WebRTCSDP.SDPType]
        ANSWER: _ClassVar[WebRTCSDP.SDPType]
    SDP_TYPE_UNSPECIFIED: WebRTCSDP.SDPType
    OFFER: WebRTCSDP.SDPType
    ANSWER: WebRTCSDP.SDPType
    TYPE_FIELD_NUMBER: _ClassVar[int]
    SDP_FIELD_NUMBER: _ClassVar[int]
    type: WebRTCSDP.SDPType
    sdp: str
    def __init__(self, type: _Optional[_Union[WebRTCSDP.SDPType, str]] = ..., sdp: _Optional[str] = ...) -> None: ...

class ClientSignaling(_message.Message):
    __slots__ = ("sessionId", "sdp", "iceCandidate", "disconnect")
    SESSIONID_FIELD_NUMBER: _ClassVar[int]
    SDP_FIELD_NUMBER: _ClassVar[int]
    ICECANDIDATE_FIELD_NUMBER: _ClassVar[int]
    DISCONNECT_FIELD_NUMBER: _ClassVar[int]
    sessionId: str
    sdp: WebRTCSDP
    iceCandidate: ICECandidate
    disconnect: bool
    def __init__(self, sessionId: _Optional[str] = ..., sdp: _Optional[_Union[WebRTCSDP, _Mapping]] = ..., iceCandidate: _Optional[_Union[ICECandidate, _Mapping]] = ..., disconnect: bool = ...) -> None: ...

class ServerSignaling(_message.Message):
    __slots__ = ("sessionId", "config", "sdp", "iceCandidate", "ready", "clear", "error")
    SESSIONID_FIELD_NUMBER: _ClassVar[int]
    CONFIG_FIELD_NUMBER: _ClassVar[int]
    SDP_FIELD_NUMBER: _ClassVar[int]
    ICECANDIDATE_FIELD_NUMBER: _ClassVar[int]
    READY_FIELD_NUMBER: _ClassVar[int]
    CLEAR_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    sessionId: str
    config: WebRTCConfig
    sdp: WebRTCSDP
    iceCandidate: ICECandidate
    ready: bool
    clear: bool
    error: str
    def __init__(self, sessionId: _Optional[str] = ..., config: _Optional[_Union[WebRTCConfig, _Mapping]] = ..., sdp: _Optional[_Union[WebRTCSDP, _Mapping]] = ..., iceCandidate: _Optional[_Union[ICECandidate, _Mapping]] = ..., ready: bool = ..., clear: bool = ..., error: _Optional[str] = ...) -> None: ...

class WebRTCConfig(_message.Message):
    __slots__ = ("iceServers", "audioCodec", "sampleRate")
    ICESERVERS_FIELD_NUMBER: _ClassVar[int]
    AUDIOCODEC_FIELD_NUMBER: _ClassVar[int]
    SAMPLERATE_FIELD_NUMBER: _ClassVar[int]
    iceServers: _containers.RepeatedCompositeFieldContainer[ICEServer]
    audioCodec: str
    sampleRate: int
    def __init__(self, iceServers: _Optional[_Iterable[_Union[ICEServer, _Mapping]]] = ..., audioCodec: _Optional[str] = ..., sampleRate: _Optional[int] = ...) -> None: ...

class StreamConfig(_message.Message):
    __slots__ = ("audio", "text")
    AUDIO_FIELD_NUMBER: _ClassVar[int]
    TEXT_FIELD_NUMBER: _ClassVar[int]
    audio: AudioConfig
    text: TextConfig
    def __init__(self, audio: _Optional[_Union[AudioConfig, _Mapping]] = ..., text: _Optional[_Union[TextConfig, _Mapping]] = ...) -> None: ...

class AudioConfig(_message.Message):
    __slots__ = ("sampleRate", "audioFormat", "channels")
    class AudioFormat(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        LINEAR16: _ClassVar[AudioConfig.AudioFormat]
        MuLaw8: _ClassVar[AudioConfig.AudioFormat]
    LINEAR16: AudioConfig.AudioFormat
    MuLaw8: AudioConfig.AudioFormat
    SAMPLERATE_FIELD_NUMBER: _ClassVar[int]
    AUDIOFORMAT_FIELD_NUMBER: _ClassVar[int]
    CHANNELS_FIELD_NUMBER: _ClassVar[int]
    sampleRate: int
    audioFormat: AudioConfig.AudioFormat
    channels: int
    def __init__(self, sampleRate: _Optional[int] = ..., audioFormat: _Optional[_Union[AudioConfig.AudioFormat, str]] = ..., channels: _Optional[int] = ...) -> None: ...

class TextConfig(_message.Message):
    __slots__ = ("charset",)
    CHARSET_FIELD_NUMBER: _ClassVar[int]
    charset: str
    def __init__(self, charset: _Optional[str] = ...) -> None: ...

class ConversationInterruption(_message.Message):
    __slots__ = ("id", "type", "time")
    class InterruptionType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        INTERRUPTION_TYPE_UNSPECIFIED: _ClassVar[ConversationInterruption.InterruptionType]
        INTERRUPTION_TYPE_VAD: _ClassVar[ConversationInterruption.InterruptionType]
        INTERRUPTION_TYPE_WORD: _ClassVar[ConversationInterruption.InterruptionType]
    INTERRUPTION_TYPE_UNSPECIFIED: ConversationInterruption.InterruptionType
    INTERRUPTION_TYPE_VAD: ConversationInterruption.InterruptionType
    INTERRUPTION_TYPE_WORD: ConversationInterruption.InterruptionType
    ID_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    TIME_FIELD_NUMBER: _ClassVar[int]
    id: str
    type: ConversationInterruption.InterruptionType
    time: _timestamp_pb2.Timestamp
    def __init__(self, id: _Optional[str] = ..., type: _Optional[_Union[ConversationInterruption.InterruptionType, str]] = ..., time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ...) -> None: ...

class ConversationAssistantMessage(_message.Message):
    __slots__ = ("audio", "text", "id", "completed", "time")
    AUDIO_FIELD_NUMBER: _ClassVar[int]
    TEXT_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    COMPLETED_FIELD_NUMBER: _ClassVar[int]
    TIME_FIELD_NUMBER: _ClassVar[int]
    audio: bytes
    text: str
    id: str
    completed: bool
    time: _timestamp_pb2.Timestamp
    def __init__(self, audio: _Optional[bytes] = ..., text: _Optional[str] = ..., id: _Optional[str] = ..., completed: bool = ..., time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ...) -> None: ...

class ConversationUserMessage(_message.Message):
    __slots__ = ("audio", "text", "id", "completed", "time")
    AUDIO_FIELD_NUMBER: _ClassVar[int]
    TEXT_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    COMPLETED_FIELD_NUMBER: _ClassVar[int]
    TIME_FIELD_NUMBER: _ClassVar[int]
    audio: bytes
    text: str
    id: str
    completed: bool
    time: _timestamp_pb2.Timestamp
    def __init__(self, audio: _Optional[bytes] = ..., text: _Optional[str] = ..., id: _Optional[str] = ..., completed: bool = ..., time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ...) -> None: ...

class AssistantTalkInput(_message.Message):
    __slots__ = ("configuration", "message", "signaling")
    CONFIGURATION_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    SIGNALING_FIELD_NUMBER: _ClassVar[int]
    configuration: ConversationConfiguration
    message: ConversationUserMessage
    signaling: ClientSignaling
    def __init__(self, configuration: _Optional[_Union[ConversationConfiguration, _Mapping]] = ..., message: _Optional[_Union[ConversationUserMessage, _Mapping]] = ..., signaling: _Optional[_Union[ClientSignaling, _Mapping]] = ...) -> None: ...

class AssistantTalkOutput(_message.Message):
    __slots__ = ("code", "success", "configuration", "interruption", "user", "assistant", "tool", "toolResult", "directive", "error", "signaling")
    CODE_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CONFIGURATION_FIELD_NUMBER: _ClassVar[int]
    INTERRUPTION_FIELD_NUMBER: _ClassVar[int]
    USER_FIELD_NUMBER: _ClassVar[int]
    ASSISTANT_FIELD_NUMBER: _ClassVar[int]
    TOOL_FIELD_NUMBER: _ClassVar[int]
    TOOLRESULT_FIELD_NUMBER: _ClassVar[int]
    DIRECTIVE_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    SIGNALING_FIELD_NUMBER: _ClassVar[int]
    code: int
    success: bool
    configuration: ConversationConfiguration
    interruption: ConversationInterruption
    user: ConversationUserMessage
    assistant: ConversationAssistantMessage
    tool: ConversationToolCall
    toolResult: ConversationToolResult
    directive: ConversationDirective
    error: _common_pb2.Error
    signaling: ServerSignaling
    def __init__(self, code: _Optional[int] = ..., success: bool = ..., configuration: _Optional[_Union[ConversationConfiguration, _Mapping]] = ..., interruption: _Optional[_Union[ConversationInterruption, _Mapping]] = ..., user: _Optional[_Union[ConversationUserMessage, _Mapping]] = ..., assistant: _Optional[_Union[ConversationAssistantMessage, _Mapping]] = ..., tool: _Optional[_Union[ConversationToolCall, _Mapping]] = ..., toolResult: _Optional[_Union[ConversationToolResult, _Mapping]] = ..., directive: _Optional[_Union[ConversationDirective, _Mapping]] = ..., error: _Optional[_Union[_common_pb2.Error, _Mapping]] = ..., signaling: _Optional[_Union[ServerSignaling, _Mapping]] = ...) -> None: ...

class CreateMessageMetricRequest(_message.Message):
    __slots__ = ("assistantId", "assistantConversationId", "messageId", "metrics")
    ASSISTANTID_FIELD_NUMBER: _ClassVar[int]
    ASSISTANTCONVERSATIONID_FIELD_NUMBER: _ClassVar[int]
    MESSAGEID_FIELD_NUMBER: _ClassVar[int]
    METRICS_FIELD_NUMBER: _ClassVar[int]
    assistantId: int
    assistantConversationId: int
    messageId: str
    metrics: _containers.RepeatedCompositeFieldContainer[_common_pb2.Metric]
    def __init__(self, assistantId: _Optional[int] = ..., assistantConversationId: _Optional[int] = ..., messageId: _Optional[str] = ..., metrics: _Optional[_Iterable[_Union[_common_pb2.Metric, _Mapping]]] = ...) -> None: ...

class CreateMessageMetricResponse(_message.Message):
    __slots__ = ("code", "success", "data", "error")
    CODE_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    code: int
    success: bool
    data: _containers.RepeatedCompositeFieldContainer[_common_pb2.Metric]
    error: _common_pb2.Error
    def __init__(self, code: _Optional[int] = ..., success: bool = ..., data: _Optional[_Iterable[_Union[_common_pb2.Metric, _Mapping]]] = ..., error: _Optional[_Union[_common_pb2.Error, _Mapping]] = ...) -> None: ...

class CreateConversationMetricRequest(_message.Message):
    __slots__ = ("assistantId", "assistantConversationId", "metrics")
    ASSISTANTID_FIELD_NUMBER: _ClassVar[int]
    ASSISTANTCONVERSATIONID_FIELD_NUMBER: _ClassVar[int]
    METRICS_FIELD_NUMBER: _ClassVar[int]
    assistantId: int
    assistantConversationId: int
    metrics: _containers.RepeatedCompositeFieldContainer[_common_pb2.Metric]
    def __init__(self, assistantId: _Optional[int] = ..., assistantConversationId: _Optional[int] = ..., metrics: _Optional[_Iterable[_Union[_common_pb2.Metric, _Mapping]]] = ...) -> None: ...

class CreateConversationMetricResponse(_message.Message):
    __slots__ = ("code", "success", "data", "error")
    CODE_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    code: int
    success: bool
    data: _containers.RepeatedCompositeFieldContainer[_common_pb2.Metric]
    error: _common_pb2.Error
    def __init__(self, code: _Optional[int] = ..., success: bool = ..., data: _Optional[_Iterable[_Union[_common_pb2.Metric, _Mapping]]] = ..., error: _Optional[_Union[_common_pb2.Error, _Mapping]] = ...) -> None: ...

class CreatePhoneCallRequest(_message.Message):
    __slots__ = ("assistant", "metadata", "args", "options", "fromNumber", "toNumber")
    class MetadataEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: _any_pb2.Any
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...
    class ArgsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: _any_pb2.Any
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...
    class OptionsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: _any_pb2.Any
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...
    ASSISTANT_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    ARGS_FIELD_NUMBER: _ClassVar[int]
    OPTIONS_FIELD_NUMBER: _ClassVar[int]
    FROMNUMBER_FIELD_NUMBER: _ClassVar[int]
    TONUMBER_FIELD_NUMBER: _ClassVar[int]
    assistant: _common_pb2.AssistantDefinition
    metadata: _containers.MessageMap[str, _any_pb2.Any]
    args: _containers.MessageMap[str, _any_pb2.Any]
    options: _containers.MessageMap[str, _any_pb2.Any]
    fromNumber: str
    toNumber: str
    def __init__(self, assistant: _Optional[_Union[_common_pb2.AssistantDefinition, _Mapping]] = ..., metadata: _Optional[_Mapping[str, _any_pb2.Any]] = ..., args: _Optional[_Mapping[str, _any_pb2.Any]] = ..., options: _Optional[_Mapping[str, _any_pb2.Any]] = ..., fromNumber: _Optional[str] = ..., toNumber: _Optional[str] = ...) -> None: ...

class CreatePhoneCallResponse(_message.Message):
    __slots__ = ("code", "success", "data", "error")
    CODE_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    code: int
    success: bool
    data: _common_pb2.AssistantConversation
    error: _common_pb2.Error
    def __init__(self, code: _Optional[int] = ..., success: bool = ..., data: _Optional[_Union[_common_pb2.AssistantConversation, _Mapping]] = ..., error: _Optional[_Union[_common_pb2.Error, _Mapping]] = ...) -> None: ...

class CreateBulkPhoneCallRequest(_message.Message):
    __slots__ = ("phoneCalls",)
    PHONECALLS_FIELD_NUMBER: _ClassVar[int]
    phoneCalls: _containers.RepeatedCompositeFieldContainer[CreatePhoneCallRequest]
    def __init__(self, phoneCalls: _Optional[_Iterable[_Union[CreatePhoneCallRequest, _Mapping]]] = ...) -> None: ...

class CreateBulkPhoneCallResponse(_message.Message):
    __slots__ = ("code", "success", "data", "error")
    CODE_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    code: int
    success: bool
    data: _containers.RepeatedCompositeFieldContainer[_common_pb2.AssistantConversation]
    error: _common_pb2.Error
    def __init__(self, code: _Optional[int] = ..., success: bool = ..., data: _Optional[_Iterable[_Union[_common_pb2.AssistantConversation, _Mapping]]] = ..., error: _Optional[_Union[_common_pb2.Error, _Mapping]] = ...) -> None: ...

class TalkInput(_message.Message):
    __slots__ = ("configuration", "message", "interruption")
    CONFIGURATION_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    INTERRUPTION_FIELD_NUMBER: _ClassVar[int]
    configuration: ConversationConfiguration
    message: ConversationUserMessage
    interruption: ConversationInterruption
    def __init__(self, configuration: _Optional[_Union[ConversationConfiguration, _Mapping]] = ..., message: _Optional[_Union[ConversationUserMessage, _Mapping]] = ..., interruption: _Optional[_Union[ConversationInterruption, _Mapping]] = ...) -> None: ...

class TalkOutput(_message.Message):
    __slots__ = ("code", "success", "interruption", "assistant", "tool", "toolResult", "directive", "error")
    CODE_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    INTERRUPTION_FIELD_NUMBER: _ClassVar[int]
    ASSISTANT_FIELD_NUMBER: _ClassVar[int]
    TOOL_FIELD_NUMBER: _ClassVar[int]
    TOOLRESULT_FIELD_NUMBER: _ClassVar[int]
    DIRECTIVE_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    code: int
    success: bool
    interruption: ConversationInterruption
    assistant: ConversationAssistantMessage
    tool: ConversationToolCall
    toolResult: ConversationToolResult
    directive: ConversationDirective
    error: _common_pb2.Error
    def __init__(self, code: _Optional[int] = ..., success: bool = ..., interruption: _Optional[_Union[ConversationInterruption, _Mapping]] = ..., assistant: _Optional[_Union[ConversationAssistantMessage, _Mapping]] = ..., tool: _Optional[_Union[ConversationToolCall, _Mapping]] = ..., toolResult: _Optional[_Union[ConversationToolResult, _Mapping]] = ..., directive: _Optional[_Union[ConversationDirective, _Mapping]] = ..., error: _Optional[_Union[_common_pb2.Error, _Mapping]] = ...) -> None: ...
