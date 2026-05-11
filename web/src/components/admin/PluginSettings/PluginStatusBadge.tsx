import { Plugin, PluginActiveState } from "@/api/openapi-schema";
import { Badge } from "@/components/ui/badge";
import { useI18n } from "@/i18n/provider";

export function PluginStatusBadge({ plugin }: { plugin: Plugin }) {
  const { t } = useI18n();
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
          {t("Active")}
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
          {t("Inactive")}
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
          {t("Starting")}
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
          {t("Connecting")}
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
          {t("Restarting")}
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
          {t("Error")}
        </Badge>
      );
  }
}
