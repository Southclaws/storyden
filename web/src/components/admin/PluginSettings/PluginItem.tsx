import { formatDate } from "date-fns";
import { useState } from "react";

import { handle } from "@/api/client";
import {
  usePluginDelete,
  usePluginSetActiveState,
} from "@/api/openapi-client/plugins";
import { Plugin, PluginActiveState } from "@/api/openapi-schema";
import { useConfirmation } from "@/components/site/useConfirmation";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { Switch } from "@/components/ui/switch";
import { Text } from "@/components/ui/text";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { CardBox as cardBox } from "@/styled-system/patterns";

import { getPluginActiveState, isPluginStatusError } from "./utils";

type Props = {
  plugin: Plugin;
};

export function PluginItem({ plugin }: Props) {
  const [isToggling, setIsToggling] = useState(false);

  const { trigger: setActiveState } = usePluginSetActiveState(plugin.id);
  const { trigger: deletePlugin } = usePluginDelete(plugin.id);

  const { isConfirming, handleConfirmAction, handleCancelAction } =
    useConfirmation(async () => {
      await deletePlugin({});
    });

  const activeState = getPluginActiveState(plugin);
  const isActive = activeState === PluginActiveState.active;
  const isError = activeState === PluginActiveState.error;

  const handleToggleActive = async () => {
    setIsToggling(true);

    await handle(
      async () => {
        await setActiveState({
          active: isActive
            ? PluginActiveState.inactive
            : PluginActiveState.active,
        });
      },
      {
        cleanup: async () => {
          setIsToggling(false);
        },
      },
    );
  };

  return (
    <li className={cardBox()}>
      <LStack>
        <WStack alignItems="center" justifyContent="space-between">
          <HStack alignItems="center">
            <Heading size="sm">{plugin.manifest.name}</Heading>
            <PluginVersionBadge plugin={plugin} />
          </HStack>

          {!isError && (
            <HStack>
              <Switch
                size="sm"
                checked={isActive}
                onClick={handleToggleActive}
                disabled={isToggling}
              />
              <PluginStatusBadge plugin={plugin} />
            </HStack>
          )}
        </WStack>

        <WStack alignItems="end">
          <styled.p fontSize="xs" color="fg.muted">
            Installed: <time>{formatDate(plugin.added_at, "PPpp")}</time>
          </styled.p>

          {isConfirming ? (
            <>
              <Button
                size="xs"
                variant="subtle"
                bgColor="bg.destructive"
                onClick={handleConfirmAction}
              >
                Confirm Delete
              </Button>
              <Button size="xs" variant="outline" onClick={handleCancelAction}>
                Cancel
              </Button>
            </>
          ) : (
            <Button
              size="xs"
              variant="outline"
              bgColor="bg.destructive"
              onClick={handleConfirmAction}
            >
              Delete
            </Button>
          )}
        </WStack>

        {isError && isPluginStatusError(plugin.status) && (
          <styled.p fontSize="xs" color="fg.error">
            Error: {plugin.status.message}
          </styled.p>
        )}
      </LStack>
    </li>
  );
}

function PluginVersionBadge({ plugin }: { plugin: Plugin }) {
  return <Badge size="sm">v{plugin.manifest.version}</Badge>;
}

function PluginStatusBadge({ plugin }: { plugin: Plugin }) {
  const activeState = getPluginActiveState(plugin);
  switch (activeState) {
    case PluginActiveState.active:
      return (
        <Badge size="sm" colorPalette="green">
          Active
        </Badge>
      );

    case PluginActiveState.inactive:
      return (
        <Badge size="sm" colorPalette="gray">
          Inactive
        </Badge>
      );

    case PluginActiveState.error:
      return (
        <Badge size="sm" colorPalette="red">
          Error
        </Badge>
      );

    default:
      return (
        <Badge size="sm" colorPalette="gray">
          Unknown
        </Badge>
      );
  }
}
