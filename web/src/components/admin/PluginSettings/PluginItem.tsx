import { formatDate } from "date-fns";
import Link from "next/link";
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
import { CardBox } from "@/components/ui/card-box";
import { Text } from "@/components/ui/text";
import { HStack, LStack, WStack } from "@/styled-system/jsx";
import { linkOverlay } from "@/styled-system/patterns";

import { PluginStatusBadge } from "./PluginStatusBadge";
import { isPluginStatusError } from "./utils";

type Props = {
  plugin: Plugin;
};

export function PluginItem({ plugin }: Props) {
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

  return (
    <CardBox as="li" position="relative">
      <LStack>
        <WStack alignItems="center" justifyContent="space-between">
          <HStack alignItems="center">
            <Link
              className={linkOverlay()}
              href={`/admin/plugins/${plugin.id}`}
            >
              <Text
                variant="supporting"
                color="text.default"
                fontWeight="semibold"
                lineClamp="1"
              >
                {plugin.name}
              </Text>
            </Link>
            <PluginVersionBadge plugin={plugin} />
          </HStack>

          <PluginStatusBadge plugin={plugin} />
        </WStack>

        <WStack alignItems="end">
          <Text variant="metadata">
            Installed: <time>{formatDate(plugin.added_at, "PPpp")}</time>
          </Text>

          <HStack position="relative" zIndex="docked">
            {isConfirming ? (
              <>
                <Button
                  variant="subtle"
                  bgColor="status.danger.surface"
                  onClick={handleConfirmAction}
                >
                  Confirm Delete
                </Button>
                <Button variant="subtle" onClick={handleCancelAction}>
                  Cancel
                </Button>
              </>
            ) : (
              <Button variant="subtle" onClick={handleConfirmAction}>
                Delete
              </Button>
            )}
          </HStack>
        </WStack>

        {isError && isPluginStatusError(plugin.status) && (
          <Text variant="metadata" color="status.danger.content">
            Error: {plugin.status.message}
          </Text>
        )}
      </LStack>
    </CardBox>
  );
}

function PluginVersionBadge({ plugin }: { plugin: Plugin }) {
  return <Badge size="sm">v{plugin.version}</Badge>;
}
