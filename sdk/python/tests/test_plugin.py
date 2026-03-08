from __future__ import annotations

import asyncio
import json

import pytest
import websockets

from storyden import Plugin


def _server_port(server: websockets.asyncio.server.Server) -> int:
    sockets = server.sockets
    assert sockets is not None
    sock = next(iter(sockets), None)
    assert sock is not None
    return int(sock.getsockname()[1])


@pytest.mark.asyncio
async def test_get_access_and_build_client() -> None:
    stop_event = asyncio.Event()
    received_id = asyncio.Event()

    async def handler(websocket) -> None:
        while True:
            message = await websocket.recv()
            payload = json.loads(message)

            assert payload["method"] == "access_get"
            assert isinstance(payload["id"], str)
            assert len(payload["id"]) == 20

            await websocket.send(
                json.dumps(
                    {
                        "jsonrpc": "2.0",
                        "id": payload["id"],
                        "result": {
                            "method": "access_get",
                            "result": {
                                "access_key": "sdbak_test_access_key",
                                "apibase_url": "http://localhost:8100",
                            },
                        },
                    }
                )
            )
            received_id.set()

            if stop_event.is_set():
                return

    server = await websockets.serve(handler, "127.0.0.1", 0)
    port = _server_port(server)

    plugin = Plugin(f"ws://127.0.0.1:{port}/rpc?token=sdprt_external_test_token")
    run_task = asyncio.create_task(plugin.run())

    try:
        await plugin.wait_until_connected(timeout=2.0)
        access = await plugin.get_access()
        assert access.access_key == "sdbak_test_access_key"
        assert access.api_base_url == "http://localhost:8100/"

        client = await plugin.build_api_client()
        try:
            assert str(client.base_url) == "http://localhost:8100/api/"
            assert client.headers["Authorization"] == "Bearer sdbak_test_access_key"
        finally:
            await client.aclose()

        await asyncio.wait_for(received_id.wait(), timeout=2.0)
    finally:
        await plugin.shutdown()
        await run_task
        stop_event.set()
        server.close()
        await server.wait_closed()


@pytest.mark.asyncio
async def test_event_dispatch_and_ack() -> None:
    stop_event = asyncio.Event()
    event_received = asyncio.Event()
    ack_received = asyncio.Event()

    async def handler(websocket) -> None:
        await websocket.send(
            json.dumps(
                {
                    "jsonrpc": "2.0",
                    "id": "9m4e2mr0ui3e8a215n4g",
                    "method": "event",
                    "params": {
                        "event": "EventThreadPublished",
                        "id": "d41d8cd98f00b204e9800998ecf8427e",
                    },
                }
            )
        )

        response = json.loads(await websocket.recv())
        assert response["id"] == "9m4e2mr0ui3e8a215n4g"
        assert response["result"]["method"] == "event"
        assert response["result"]["ok"] is True

        ack_received.set()
        await stop_event.wait()

    server = await websockets.serve(handler, "127.0.0.1", 0)
    port = _server_port(server)

    plugin = Plugin(f"ws://127.0.0.1:{port}/rpc?token=sdprt_external_test_token")

    @plugin.on("EventThreadPublished")
    async def _on_thread_published(event) -> None:
        assert event.id == "d41d8cd98f00b204e9800998ecf8427e"
        event_received.set()

    run_task = asyncio.create_task(plugin.run())

    try:
        await plugin.wait_until_connected(timeout=2.0)
        await asyncio.wait_for(event_received.wait(), timeout=2.0)
        await asyncio.wait_for(ack_received.wait(), timeout=2.0)
    finally:
        await plugin.shutdown()
        await run_task
        stop_event.set()
        server.close()
        await server.wait_closed()


@pytest.mark.asyncio
async def test_configure_dispatch_and_ack() -> None:
    stop_event = asyncio.Event()
    configure_received = asyncio.Event()
    ack_received = asyncio.Event()

    async def handler(websocket) -> None:
        await websocket.send(
            json.dumps(
                {
                    "jsonrpc": "2.0",
                    "id": "9m4e2mr0ui3e8a215n4h",
                    "method": "configure",
                    "params": {
                        "name": "configured",
                        "enabled": True,
                    },
                }
            )
        )

        response = json.loads(await websocket.recv())
        assert response["id"] == "9m4e2mr0ui3e8a215n4h"
        assert response["result"]["method"] == "configure"
        assert response["result"]["ok"] is True

        ack_received.set()
        await stop_event.wait()

    server = await websockets.serve(handler, "127.0.0.1", 0)
    port = _server_port(server)

    plugin = Plugin(f"ws://127.0.0.1:{port}/rpc?token=sdprt_external_test_token")

    @plugin.on_configure()
    async def _on_configure(config: dict[str, object]) -> bool:
        assert config["name"] == "configured"
        assert config["enabled"] is True
        configure_received.set()
        return True

    run_task = asyncio.create_task(plugin.run())

    try:
        await plugin.wait_until_connected(timeout=2.0)
        await asyncio.wait_for(configure_received.wait(), timeout=2.0)
        await asyncio.wait_for(ack_received.wait(), timeout=2.0)
    finally:
        await plugin.shutdown()
        await run_task
        stop_event.set()
        server.close()
        await server.wait_closed()


@pytest.mark.asyncio
async def test_configure_handler_rejection_returns_not_ok() -> None:
    stop_event = asyncio.Event()
    ack_received = asyncio.Event()

    async def handler(websocket) -> None:
        await websocket.send(
            json.dumps(
                {
                    "jsonrpc": "2.0",
                    "id": "9m4e2mr0ui3e8a215n4i",
                    "method": "configure",
                    "params": {
                        "reject": True,
                    },
                }
            )
        )

        response = json.loads(await websocket.recv())
        assert response["id"] == "9m4e2mr0ui3e8a215n4i"
        assert response["result"]["method"] == "configure"
        assert response["result"]["ok"] is False

        ack_received.set()
        await stop_event.wait()

    server = await websockets.serve(handler, "127.0.0.1", 0)
    port = _server_port(server)

    plugin = Plugin(f"ws://127.0.0.1:{port}/rpc?token=sdprt_external_test_token")

    @plugin.on_configure
    async def _on_configure(_config: dict[str, object]) -> bool:
        return False

    run_task = asyncio.create_task(plugin.run())

    try:
        await plugin.wait_until_connected(timeout=2.0)
        await asyncio.wait_for(ack_received.wait(), timeout=2.0)
    finally:
        await plugin.shutdown()
        await run_task
        stop_event.set()
        server.close()
        await server.wait_closed()
