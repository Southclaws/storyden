class StorydenError(Exception):
    """Base SDK exception."""


class ConnectionFailedError(StorydenError):
    """Raised when the plugin cannot connect to Storyden RPC."""


class NotAuthorisedError(ConnectionFailedError):
    """Raised when Storyden rejects RPC authentication."""


class PluginConnectionClosedError(StorydenError):
    """Raised when trying to use a closed plugin connection."""


class RPCError(StorydenError):
    """Raised when the host returns an RPC-level error."""

    def __init__(self, message: str, code: int | None = None) -> None:
        super().__init__(message)
        self.code = code
