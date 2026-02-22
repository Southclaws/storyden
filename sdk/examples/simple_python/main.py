from __future__ import annotations

import asyncio
import sys

from dotenv import load_dotenv
from storyden import Event, Plugin, StorydenError

load_dotenv()


async def main() -> None:
    print("Starting plugin...")
    plugin = Plugin.from_env()

    print(f"Plugin started: {plugin}")

    @plugin.on(Event.EVENTTHREADPUBLISHED)
    async def on_thread_published(event) -> None:
        print(f"thread published: {event.id}")

    try:
        await plugin.run()
    finally:
        await plugin.shutdown()


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("Stopped plugin")
    except StorydenError as exc:
        print(str(exc), file=sys.stderr)
        raise SystemExit(1) from None
