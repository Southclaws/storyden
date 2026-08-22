import { Plugin, PluginActiveState } from "@/api/openapi-schema";
import { StatusBadge } from "@/components/ui/status-badge";

export function PluginStatusBadge({ plugin }: { plugin: Plugin }) {
  const activeState = plugin.status.active_state;

  switch (activeState) {
    case PluginActiveState.active:
      return (
        <StatusBadge size="sm" tone="success">
          Active
        </StatusBadge>
      );

    case PluginActiveState.inactive:
      return (
        <StatusBadge size="sm" tone="neutral">
          Inactive
        </StatusBadge>
      );

    case PluginActiveState.starting:
      return (
        <StatusBadge size="sm" tone="info">
          Starting
        </StatusBadge>
      );

    case PluginActiveState.connecting:
      return (
        <StatusBadge size="sm" tone="info">
          Connecting
        </StatusBadge>
      );

    case PluginActiveState.restarting:
      return (
        <StatusBadge size="sm" tone="warning">
          Restarting
        </StatusBadge>
      );

    case PluginActiveState.error:
      return (
        <StatusBadge size="sm" tone="danger">
          Error
        </StatusBadge>
      );
  }
}
