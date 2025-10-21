import { Plugin, PluginActiveState } from "@/api/openapi-schema";
import { Badge } from "@/components/ui/badge";

export function PluginStatusBadge({ plugin }: { plugin: Plugin }) {
  const activeState = plugin.status.active_state;

  switch (activeState) {
    case PluginActiveState.active:
      return (
        <Badge
          size="sm"
          borderColor="border.success"
          backgroundColor="bg.success"
          color="fg.success"
        >
          Active
        </Badge>
      );

    case PluginActiveState.inactive:
      return (
        <Badge
          size="sm"
          borderColor="border.muted"
          backgroundColor="bg.muted"
          color="fg.muted"
        >
          Inactive
        </Badge>
      );

    case PluginActiveState.starting:
      return (
        <Badge
          size="sm"
          borderColor="border.info"
          backgroundColor="bg.info"
          color="fg.info"
        >
          Starting
        </Badge>
      );

    case PluginActiveState.connecting:
      return (
        <Badge
          size="sm"
          borderColor="border.info"
          backgroundColor="bg.info"
          color="fg.info"
        >
          Connecting
        </Badge>
      );

    case PluginActiveState.restarting:
      return (
        <Badge
          size="sm"
          borderColor="border.warning"
          backgroundColor="bg.warning"
          color="fg.warning"
        >
          Restarting
        </Badge>
      );

    case PluginActiveState.error:
      return (
        <Badge
          size="sm"
          borderColor="border.error"
          backgroundColor="bg.error"
          color="fg.error"
        >
          Error
        </Badge>
      );
  }
}
