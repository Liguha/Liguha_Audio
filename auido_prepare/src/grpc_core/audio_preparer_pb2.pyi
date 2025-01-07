from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Audio(_message.Message):
    __slots__ = ("sample_rate", "data")
    SAMPLE_RATE_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    sample_rate: int
    data: _containers.RepeatedScalarFieldContainer[float]
    def __init__(self, sample_rate: _Optional[int] = ..., data: _Optional[_Iterable[float]] = ...) -> None: ...

class AudioResponse(_message.Message):
    __slots__ = ("ok", "err", "id")
    OK_FIELD_NUMBER: _ClassVar[int]
    ERR_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    ok: bool
    err: str
    id: str
    def __init__(self, ok: bool = ..., err: _Optional[str] = ..., id: _Optional[str] = ...) -> None: ...
