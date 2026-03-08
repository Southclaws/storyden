from __future__ import annotations

import asyncio
import contextlib
import inspect
import json
import logging
import os
import secrets
import time
from collections.abc import Awaitable, Callable, Mapping
from dataclasses import dataclass
from enum import StrEnum
from typing import Any, TypeAlias, cast, overload
from urllib.parse import parse_qs, urlparse

import httpx
import websockets
from pydantic import BaseModel, TypeAdapter, ValidationError
from websockets.exceptions import ConnectionClosed, InvalidStatus

from .errors import (
    ConnectionFailedError,
    NotAuthorisedError,
    PluginConnectionClosedError,
    RPCError,
    StorydenError,
)
from .rpc import models
from .rpc.events import EventName

EventHandler: TypeAlias = Callable[[models.EventPayload], Awaitable[None] | None]
ConfigureHandler: TypeAlias = Callable[[dict[str, Any]], Awaitable[bool | None] | bool | None]

_EXTERNAL_TOKEN_PREFIX = "sdprt_"
_INITIAL_RECONNECT_WAIT_SECONDS = 0.25
_MAX_RECONNECT_WAIT_SECONDS = 10.0
_DEFAULT_RPC_TIMEOUT_SECONDS = 30.0
_XID_ALPHABET = "0123456789abcdefghijklmnopqrstuv"

_HOST_TO_PLUGIN_REQUEST_ADAPTER = TypeAdapter(models.HostToPluginRequest)


class PluginMode(StrEnum):
    SUPERVISED = "supervised"
    EXTERNAL = "external"


@dataclass(slots=True, frozen=True)
class AccessCredentials:
    access_key: str
    api_base_url: str


@dataclass(slots=True)
class _OutboundWrite:
    data: str
    ack: asyncio.Future[None]


@dataclass(slots=True)
class _ConnectionState:
    websocket: Any
    outbound: asyncio.Queue[_OutboundWrite]
    done: asyncio.Event


class Plugin:
    def __init__(
        self,
        rpc_url: str,
        *,
        logger: logging.Logger | None = None,
    ) -> None:
        self._rpc_url = rpc_url
        self._logger = logger or logging.getLogger("storyden.plugin")
        self._mode = _mode_from_rpc_url(rpc_url)
        self._start_time: float | None = None

        self._handlers: dict[EventName, EventHandler] = {}
        self._configure_handler: ConfigureHandler | None = None
        self._pending: dict[str, asyncio.Future[dict[str, Any]]] = {}
        self._state: _ConnectionState | None = None

        self._shutdown_event = asyncio.Event()
        self._connected_event = asyncio.Event()
        self._run_completed_event = asyncio.Event()
        self._run_started = False

    @classmethod
    def from_env(cls, *, logger: logging.Logger | None = None) -> Plugin:
        rpc_url = os.getenv("STORYDEN_RPC_URL")
        if not rpc_url:
            msg = "STORYDEN_RPC_URL environment variable is not set"
            raise ValueError(msg)
        return cls(rpc_url, logger=logger)

    @property
    def mode(self) -> PluginMode:
        return self._mode

    @property
    def connected(self) -> bool:
        return self._state is not None and not self._state.done.is_set()

    @overload
    def on(self, event_type: EventName) -> Callable[[EventHandler], EventHandler]: ...

    @overload
    def on(self, event_type: EventName, handler: EventHandler) -> EventHandler: ...

    @overload
    def on(self, event_type: models.Event) -> Callable[[EventHandler], EventHandler]: ...

    @overload
    def on(self, event_type: models.Event, handler: EventHandler) -> EventHandler: ...

    def on(
        self,
        event_type: models.Event | EventName,
        handler: EventHandler | None = None,
    ) -> EventHandler | Callable[[EventHandler], EventHandler]:
        event_name = cast(EventName, event_type.value) if isinstance(event_type, models.Event) else event_type
        if handler is None:
            def decorator(fn: EventHandler) -> EventHandler:
                self._handlers[event_name] = fn
                self._logger.debug("registered handler", extra={"event_type": event_name})
                return fn

            return decorator

        self._handlers[event_name] = handler
        self._logger.debug("registered handler", extra={"event_type": event_name})
        return handler

    @overload
    def on_configure(self) -> Callable[[ConfigureHandler], ConfigureHandler]: ...

    @overload
    def on_configure(self, handler: ConfigureHandler) -> ConfigureHandler: ...

    def on_configure(
        self, handler: ConfigureHandler | None = None
    ) -> ConfigureHandler | Callable[[ConfigureHandler], ConfigureHandler]:
        if handler is None:

            def decorator(fn: ConfigureHandler) -> ConfigureHandler:
                self._configure_handler = fn
                self._logger.debug("registered configure handler")
                return fn

            return decorator

        self._configure_handler = handler
        self._logger.debug("registered configure handler")
        return handler

    async def wait_until_connected(self, timeout: float | None = None) -> None:
        if timeout is None:
            await self._connected_event.wait()
            return

        await asyncio.wait_for(self._connected_event.wait(), timeout=timeout)

    async def run(self) -> None:
        if self._run_started:
            msg = "plugin is already running"
            raise RuntimeError(msg)

        self._run_started = True
        self._run_completed_event.clear()

        retry_wait = _INITIAL_RECONNECT_WAIT_SECONDS

        try:
            while not self._shutdown_event.is_set():
                try:
                    websocket = await websockets.connect(self._rpc_url)
                except InvalidStatus as exc:
                    status_code = _extract_status_code(exc)
                    if status_code == 401:
                        raise NotAuthorisedError("401 Not Authorised (check your token)") from None

                    if status_code == 403:
                        raise NotAuthorisedError("403 Forbidden (token lacks permission)") from None

                    if (
                        self._mode is PluginMode.EXTERNAL
                        and _should_retry_dial(exc)
                        and not self._shutdown_event.is_set()
                    ):
                        self._logger.warning(
                            "failed to connect, retrying",
                            extra={"retry_in": retry_wait, "error": str(exc)},
                        )
                        await asyncio.sleep(retry_wait)
                        retry_wait = min(retry_wait * 2, _MAX_RECONNECT_WAIT_SECONDS)
                        continue

                    status = f"HTTP {status_code}" if status_code is not None else "unknown status"
                    raise ConnectionFailedError(f"WebSocket rejected with {status}") from None
                except Exception as exc:
                    if (
                        self._mode is PluginMode.EXTERNAL
                        and _should_retry_dial(exc)
                        and not self._shutdown_event.is_set()
                    ):
                        self._logger.warning(
                            "failed to connect, retrying",
                            extra={"retry_in": retry_wait, "error": str(exc)},
                        )
                        await asyncio.sleep(retry_wait)
                        retry_wait = min(retry_wait * 2, _MAX_RECONNECT_WAIT_SECONDS)
                        continue

                    raise ConnectionFailedError(f"Failed to connect to {self._rpc_url}") from None

                retry_wait = _INITIAL_RECONNECT_WAIT_SECONDS
                if self._start_time is None:
                    self._start_time = time.monotonic()

                state = _ConnectionState(
                    websocket=websocket,
                    outbound=asyncio.Queue(),
                    done=asyncio.Event(),
                )
                self._set_connection_state(state)

                disconnect_error: Exception | None = None
                try:
                    disconnect_error = await self._run_connection_loops(state)
                finally:
                    await self._clear_connection_state(state)
                    await self._close_websocket(websocket)
                    self._clear_pending(PluginConnectionClosedError("connection closed"))

                if self._mode is PluginMode.SUPERVISED:
                    return

                if self._shutdown_event.is_set():
                    return

                if _should_retry_disconnect(disconnect_error):
                    self._logger.warning(
                        "connection dropped, reconnecting",
                        extra={"retry_in": retry_wait, "error": str(disconnect_error)},
                    )
                    await asyncio.sleep(retry_wait)
                    retry_wait = min(retry_wait * 2, _MAX_RECONNECT_WAIT_SECONDS)
                    continue

                if _is_graceful_disconnect(disconnect_error):
                    self._logger.warning("connection closed", extra={"disconnect": _disconnect_summary(disconnect_error)})
                    return

                if disconnect_error is not None:
                    raise _map_disconnect_error(disconnect_error)

                self._logger.warning(
                    "connection closed without an explicit error",
                )
                return
        finally:
            self._connected_event.clear()
            self._run_completed_event.set()

    async def shutdown(self) -> None:
        self._shutdown_event.set()

        state = self._state
        if state is not None:
            state.done.set()
            await self._close_websocket(state.websocket)

        self._clear_pending(PluginConnectionClosedError("plugin shutting down"))

        if not self._run_started:
            return

        try:
            await asyncio.wait_for(self._run_completed_event.wait(), timeout=5.0)
        except TimeoutError:
            self._logger.warning("timeout waiting for plugin run loop shutdown")

    async def send(
        self,
        payload: models.PluginToHostRequest | Mapping[str, Any] | BaseModel,
        *,
        timeout: float | None = _DEFAULT_RPC_TIMEOUT_SECONDS,
    ) -> dict[str, Any]:
        request_id, body = _marshal_request_with_generated_id(payload)

        response_future: asyncio.Future[dict[str, Any]] = asyncio.get_running_loop().create_future()
        self._pending[request_id] = response_future

        data = json.dumps(body)

        try:
            await self._enqueue_write(data, timeout=timeout)
            response = await _await_with_timeout(response_future, timeout)
        except Exception:
            self._pending.pop(request_id, None)
            raise

        error = response.get("error")
        if isinstance(error, Mapping):
            message = str(error.get("message") or "rpc error")
            code = error.get("code")
            raise RPCError(message=message, code=code if isinstance(code, int) else None)

        result = response.get("result")
        if not isinstance(result, Mapping):
            msg = "rpc response missing result"
            raise RuntimeError(msg)

        return dict(result)

    async def get_access(
        self,
        *,
        timeout: float | None = _DEFAULT_RPC_TIMEOUT_SECONDS,
    ) -> AccessCredentials:
        result = await self.send(
            {
                "jsonrpc": "2.0",
                "method": "access_get",
            },
            timeout=timeout,
        )

        if result.get("method") != "access_get":
            msg = f"unexpected RPC response method: {result.get('method')!r}"
            raise RuntimeError(msg)

        method_error = result.get("error")
        if isinstance(method_error, Mapping):
            message = str(method_error.get("message") or "access_get error")
            code = method_error.get("code")
            raise RPCError(message=message, code=code if isinstance(code, int) else None)

        response_payload = result.get("result")
        if not isinstance(response_payload, Mapping):
            msg = "access_get response missing result"
            raise RuntimeError(msg)

        parsed = models.RPCResponseAccessGetResult.model_validate(response_payload)
        return AccessCredentials(access_key=parsed.access_key, api_base_url=str(parsed.apibase_url))

    async def build_api_client(self) -> httpx.AsyncClient:
        access = await self.get_access()
        api_base = access.api_base_url.rstrip("/") + "/api"

        return httpx.AsyncClient(
            base_url=api_base,
            headers={"Authorization": f"Bearer {access.access_key}"},
        )

    async def _run_connection_loops(self, state: _ConnectionState) -> Exception | None:
        reader_task = asyncio.create_task(self._read_loop(state))
        writer_task = asyncio.create_task(self._write_loop(state))
        tasks = {reader_task, writer_task}

        try:
            done, pending = await asyncio.wait(
                tasks,
                return_when=asyncio.FIRST_COMPLETED,
            )
        except asyncio.CancelledError:
            state.done.set()
            for task in tasks:
                task.cancel()
            await asyncio.gather(*tasks, return_exceptions=True)
            raise

        disconnect_error: Exception | None = None

        for task in done:
            try:
                task.result()
            except asyncio.CancelledError:
                pass
            except Exception as exc:  # noqa: BLE001
                if disconnect_error is None:
                    disconnect_error = exc

        for task in pending:
            task.cancel()

        pending_results = await asyncio.gather(*pending, return_exceptions=True)
        for result in pending_results:
            if isinstance(result, Exception) and disconnect_error is None:
                disconnect_error = result

        state.done.set()
        return disconnect_error

    def _set_connection_state(self, state: _ConnectionState) -> None:
        self._state = state
        self._connected_event.set()

    async def _clear_connection_state(self, state: _ConnectionState) -> None:
        if self._state is state:
            self._state = None
        self._connected_event.clear()

    def _clear_pending(self, exc: Exception) -> None:
        pending = self._pending
        self._pending = {}
        for fut in pending.values():
            if not fut.done():
                fut.set_exception(exc)

    async def _enqueue_write(self, data: str, *, timeout: float | None) -> None:
        state = self._state
        if state is None or state.done.is_set():
            raise PluginConnectionClosedError("connection closed")

        ack: asyncio.Future[None] = asyncio.get_running_loop().create_future()
        await state.outbound.put(_OutboundWrite(data=data, ack=ack))
        await _await_with_timeout(ack, timeout)

    async def _read_loop(self, state: _ConnectionState) -> None:
        while not state.done.is_set() and not self._shutdown_event.is_set():
            message = await state.websocket.recv()
            if isinstance(message, bytes):
                message = message.decode("utf-8")

            await self._handle_message(message)

    async def _write_loop(self, state: _ConnectionState) -> None:
        while not state.done.is_set() and not self._shutdown_event.is_set():
            write = await state.outbound.get()
            await state.websocket.send(write.data)
            if not write.ack.done():
                write.ack.set_result(None)

    async def _handle_message(self, message: str) -> None:
        try:
            payload = json.loads(message)
        except json.JSONDecodeError:
            self._logger.warning("invalid message payload", extra={"payload": message})
            return

        if not isinstance(payload, Mapping):
            self._logger.warning("ignoring non-object rpc payload")
            return

        if "method" in payload:
            await self._handle_host_request(payload)
            return

        raw_id = payload.get("id")
        if not isinstance(raw_id, str):
            self._logger.warning("response missing id")
            return

        fut = self._pending.pop(raw_id, None)
        if fut is None:
            self._logger.warning("response for unknown request id", extra={"id": raw_id})
            return

        if not fut.done():
            fut.set_result(dict(payload))

    async def _handle_host_request(self, payload: Mapping[str, Any]) -> None:
        try:
            request = _HOST_TO_PLUGIN_REQUEST_ADAPTER.validate_python(payload)
        except ValidationError as exc:
            self._logger.warning("invalid host request payload", extra={"error": str(exc)})
            return

        if isinstance(request, models.RPCRequestEvent):
            await self._handle_event_request(request)
            return

        if isinstance(request, models.RPCRequestConfigure):
            await self._handle_configure_request(request)
            return

        if isinstance(request, models.RPCRequestPing):
            await self._handle_ping_request(request)
            return

        self._logger.warning("unknown host request type", extra={"type": type(request).__name__})

    async def _handle_event_request(self, request: models.RPCRequestEvent) -> None:
        event_type = getattr(request.params, "event", "")
        handler = self._handlers.get(cast(EventName, event_type))

        if handler is None:
            self._logger.warning("no handler for event", extra={"event_type": event_type})
            await self._send_result(request.id, {"method": "event", "ok": True})
            return

        try:
            result = handler(request.params)
            if inspect.isawaitable(result):
                await result
        except Exception as exc:  # noqa: BLE001
            await self._send_error(request.id, code=-32000, message=f"handler error: {exc}")
            return

        await self._send_result(request.id, {"method": "event", "ok": True})

    async def _handle_configure_request(self, request: models.RPCRequestConfigure) -> None:
        handler = self._configure_handler
        if handler is None:
            await self._send_result(request.id, {"method": "configure", "ok": True})
            return

        ok = True
        try:
            result = handler(dict(request.params))
            if inspect.isawaitable(result):
                result = await result
            if result is not None:
                ok = bool(result)
        except Exception as exc:  # noqa: BLE001
            self._logger.warning("configure handler error", extra={"error": str(exc)})
            ok = False

        await self._send_result(request.id, {"method": "configure", "ok": ok})

    async def _handle_ping_request(self, request: models.RPCRequestPing) -> None:
        start_time = self._start_time or time.monotonic()
        uptime_seconds = max(time.monotonic() - start_time, 0.0)

        await self._send_result(
            request.id,
            {
                "method": "ping",
                "pong": True,
                "status": "healthy",
                "uptime_seconds": uptime_seconds,
            },
        )

    async def _send_result(self, request_id: str, result: Mapping[str, Any]) -> None:
        await self._enqueue_write(
            json.dumps(
                {
                    "jsonrpc": "2.0",
                    "id": request_id,
                    "result": dict(result),
                }
            ),
            timeout=None,
        )

    async def _send_error(self, request_id: str, *, code: int, message: str) -> None:
        await self._enqueue_write(
            json.dumps(
                {
                    "jsonrpc": "2.0",
                    "id": request_id,
                    "error": {
                        "code": code,
                        "message": message,
                    },
                }
            ),
            timeout=None,
        )

    async def _close_websocket(self, websocket: Any) -> None:
        with contextlib.suppress(Exception):
            await websocket.close(code=1000, reason="shutting down")


def _mode_from_rpc_url(rpc_url: str) -> PluginMode:
    token = parse_qs(urlparse(rpc_url).query).get("token", [""])[0]
    if token.startswith(_EXTERNAL_TOKEN_PREFIX):
        return PluginMode.EXTERNAL
    return PluginMode.SUPERVISED


def _extract_status_code(exc: Exception) -> int | None:
    status_code = getattr(exc, "status_code", None)
    if isinstance(status_code, int):
        return status_code

    response = getattr(exc, "response", None)
    if response is not None:
        response_status = getattr(response, "status_code", None)
        if isinstance(response_status, int):
            return response_status

    return None


def _should_retry_dial(exc: Exception) -> bool:
    status_code = _extract_status_code(exc)
    if status_code is not None and 400 <= status_code < 500:
        return False

    return not (isinstance(exc, ConnectionClosed) and exc.code in {1003, 1008})


def _should_retry_disconnect(exc: Exception | None) -> bool:
    if exc is None:
        return False

    if isinstance(exc, asyncio.CancelledError):
        return False

    if isinstance(exc, ConnectionClosed):
        # External plugin sessions can be intentionally restarted by the host
        # (e.g. manifest update flow), which may use normal close semantics.
        if exc.code == 1000:
            return True

        # Connection can be rejected transiently while a plugin session is
        # being unloaded/reloaded. Retry this case for better DX.
        if exc.code == 1008 and exc.reason == "connection rejected":
            return True

        return exc.code in {1001, 1006, 1011, 1012, 1013}

    return isinstance(exc, OSError)


def _is_graceful_disconnect(exc: Exception | None) -> bool:
    if exc is None:
        return True

    if isinstance(exc, asyncio.CancelledError):
        return True

    if isinstance(exc, ConnectionClosed):
        return exc.code == 1000

    return False


def _map_disconnect_error(exc: Exception) -> StorydenError:
    if isinstance(exc, ConnectionClosed):
        reason = f"WebSocket closed ({exc.code})"
        if exc.reason:
            reason = f"{reason}: {exc.reason}"

        if exc.code in {1008, 4001, 4401, 4403}:
            return NotAuthorisedError(reason)

        return ConnectionFailedError(reason)

    if isinstance(exc, OSError):
        return ConnectionFailedError(str(exc))

    return ConnectionFailedError(str(exc))


def _disconnect_summary(exc: Exception | None) -> str:
    if exc is None:
        return "closed without error"

    if isinstance(exc, ConnectionClosed):
        message = f"code={exc.code}"
        if exc.reason:
            return f"{message} reason={exc.reason}"
        return message

    return str(exc)


async def _await_with_timeout(
    fut: asyncio.Future[Any] | asyncio.Task[Any],
    timeout: float | None,
) -> Any:
    if timeout is None:
        return await fut

    return await asyncio.wait_for(fut, timeout=timeout)


def _marshal_request_with_generated_id(
    payload: models.PluginToHostRequest | Mapping[str, Any] | BaseModel,
) -> tuple[str, dict[str, Any]]:
    body = payload.model_dump(mode="json", exclude_none=True) if isinstance(payload, BaseModel) else dict(payload)

    if not body:
        msg = "request payload cannot be empty"
        raise ValueError(msg)

    request_id = _new_xid()
    body["id"] = request_id
    body.setdefault("jsonrpc", "2.0")
    return request_id, body


def _new_xid() -> str:
    value = int.from_bytes(secrets.token_bytes(12), byteorder="big", signed=False)
    chars = ["0"] * 20
    for i in range(19, -1, -1):
        chars[i] = _XID_ALPHABET[value & 0x1F]
        value >>= 5
    return "".join(chars)
