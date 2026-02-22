import { formatDate } from "date-fns";
import { useSWRConfig } from "swr";

import { mutateTransaction } from "@/api/mutate";
import {
  getPluginListKey,
  usePluginDelete,
} from "@/api/openapi-client/plugins";
import { Plugin, PluginListOKResponse } from "@/api/openapi-schema";
import { useConfirmation } from "@/components/site/useConfirmation";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { CardBox as cardBox } from "@/styled-system/patterns";

import { PluginStatusBadge } from "./PluginStatusBadge";
import { useSelectedPlugin } from "./useSelectedPlugin";
import { isPluginStatusError } from "./utils";

type Props = {
  plugin: Plugin;
};

export function PluginItem({ plugin }: Props) {
  const [_, setSelectedPlugin] = useSelectedPlugin();
  const { mutate } = useSWRConfig();

  const { trigger: deletePlugin } = usePluginDelete(plugin.id);

  const { isConfirming, handleConfirmAction, handleCancelAction } =
    useConfirmation(async () => {
      await mutateTransaction(
        mutate,
        [
          {
            key: getPluginListKey(),
            optimistic: (current: PluginListOKResponse | undefined) => {
              if (!current) return current;
              return {
                ...current,
                plugins: current.plugins.filter((p) => p.id !== plugin.id),
              };
            },
          },
        ],
        () => deletePlugin({}),
        { revalidate: true },
      );
    });

  const isError = plugin.status.active_state === "error";

  const handleSelectPlugin = () => {
    setSelectedPlugin(plugin.id);
  };

  return (
    <li className={cardBox()}>
      <LStack>
        <WStack alignItems="center" justifyContent="space-between">
          <HStack alignItems="center">
            <a href="#" onClick={handleSelectPlugin}>
              <Heading lineClamp="1" size="sm">
                {plugin.name}
              </Heading>
            </a>
            <PluginVersionBadge plugin={plugin} />
          </HStack>

          <PluginStatusBadge plugin={plugin} />
        </WStack>

        <WStack alignItems="end">
          <styled.p fontSize="xs" color="fg.muted">
            Installed: <time>{formatDate(plugin.added_at, "PPpp")}</time>
          </styled.p>

          <HStack>
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
                <Button size="xs" variant="subtle" onClick={handleCancelAction}>
                  Cancel
                </Button>
              </>
            ) : (
              <Button size="xs" variant="subtle" onClick={handleConfirmAction}>
                Delete
              </Button>
            )}
          </HStack>
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
  return <Badge size="sm">v{plugin.version}</Badge>;
}
