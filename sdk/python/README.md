# Storyden Python SDK

用于构建 Storyden 插件和集成的 Python SDK。

Python SDK for building Storyden plugins and integrations.

## Highlights

## 亮点

- 通过 WebSocket 连接 Storyden RPC。
- 处理 host -> plugin RPC。
- 发送 plugin -> host RPC。
- 支持外部插件模式。
- 包含 API access helpers（`get_access`、`build_api_client`）。

- Connects to Storyden RPC over WebSocket.
- Handles host -> plugin RPCs.
- Sends plugin -> host RPCs.
- External plugin mode supported.
- Includes API access helpers (`get_access`, `build_api_client`).

## Quick Start

## 快速开始

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

### 外部插件

外部插件不由 Storyden 进程托管，它们从进程外连接到 Storyden。

External plugins are not supervised by the Storyden process, they connect from outside.

构建外部插件时，请在运行前设置环境变量 `STORYDEN_RPC_URL`，`from_env()` 会读取它。

When building an External plugin, set `STORYDEN_RPC_URL` in your environment before running, `from_env()` will pick this up.
