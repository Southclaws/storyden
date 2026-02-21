# Storyden Python SDK

Python SDK for building Storyden plugins and integrations.

## Highlights

- Connects to Storyden RPC over WebSocket.
- Handles host -> plugin RPCs.
- Sends plugin -> host RPCs.
- External plugin mode supported.
- Includes API access helpers (`get_access`, `build_api_client`).

## Quick Start

```python
import asyncio
from storyden import Event, Plugin


async def main() -> None:
    plugin = Plugin.from_env()

    @plugin.on(Event.EVENTTHREADPUBLISHED)
    async def on_thread_published(event):
        print("Thread published:", event.id)

    await plugin.run()


if __name__ == "__main__":
    asyncio.run(main())
```

### External Plugins

External plugins are not supervised by the Storyden process, they connect from outside.

When building an External plugin, set `STORYDEN_RPC_URL` in your environment before running, `from_env()` will pick this up.
