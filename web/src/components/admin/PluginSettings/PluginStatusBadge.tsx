import { Plugin, PluginActiveState } from "@/api/openapi-schema";
import { Badge } from "@/components/ui/badge";

export function PluginStatusBadge({ plugin }: { plugin: Plugin }) {
  const activeState = plugin.status.active_state;

  switch (activeState) {
    case PluginActiveState.active:
      return (
        <Badge
          size="sm"
          borderColor="status.success.border"
          backgroundColor="status.success.surface"
          color="status.success.content"
        >
          Active
        </Badge>
      );

    case PluginActiveState.inactive:
      return (
        <Badge
          size="sm"
          borderColor="border.strong"
          backgroundColor="control.disabledBackground"
          color="text.subtle"
        >
          Inactive
        </Badge>
      );

    case PluginActiveState.starting:
      return (
        <Badge
          size="sm"
          borderColor="status.info.border"
          backgroundColor="status.info.surface"
          color="status.info.content"
        >
          Starting
        </Badge>
      );

    case PluginActiveState.connecting:
      return (
        <Badge
          size="sm"
          borderColor="status.info.border"
          backgroundColor="status.info.surface"
          color="status.info.content"
        >
          Connecting
        </Badge>
      );

    case PluginActiveState.restarting:
      return (
        <Badge
          size="sm"
          borderColor="status.warning.border"
          backgroundColor="status.warning.surface"
          color="status.warning.content"
        >
          Restarting
        </Badge>
      );

    case PluginActiveState.error:
      return (
        <Badge
          size="sm"
          borderColor="status.danger.border"
          backgroundColor="status.danger.surface"
          color="status.danger.content"
        >
          Error
        </Badge>
      );
  }
}
