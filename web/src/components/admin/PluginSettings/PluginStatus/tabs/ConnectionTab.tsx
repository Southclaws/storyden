import { ClipboardIcon } from "lucide-react";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import { mutateTransaction } from "@/api/mutate";
import {
  getPluginGetKey,
  usePluginCycleToken,
} from "@/api/openapi-client/plugins";
import { Plugin, PluginExternalProps } from "@/api/openapi-schema";
import { InfoTip } from "@/components/site/InfoTip";
import { useConfirmation } from "@/components/site/useConfirmation";
import * as Alert from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import * as Clipboard from "@/components/ui/clipboard";
import { IconButton } from "@/components/ui/icon-button";
import { CheckIcon } from "@/components/ui/icons/Check";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { Input } from "@/components/ui/input";
import { API_ADDRESS } from "@/config";
import { useI18n } from "@/i18n/provider";
import { HStack, LStack, styled } from "@/styled-system/jsx";

type Props = {
  plugin: Plugin & { connection: PluginExternalProps };
};

export function ConnectionTab({ plugin }: Props) {
  const { t } = useI18n();
  const { mutate } = useSWRConfig();
  const { trigger: cycleToken } = usePluginCycleToken(plugin.id);
  const envURL = `STORYDEN_RPC_URL=${buildRPCURL(plugin.connection.token)}`;

  const { isConfirming, handleConfirmAction, handleCancelAction } =
    useConfirmation(async () => {
      await handle(async () => {
        await mutateTransaction(
          mutate,
          [
            {
              key: getPluginGetKey(plugin.id),
              optimistic: (current) => current,
              commit: (current, result) => {
                if (!current || current.connection.mode !== "external") {
                  return current;
                }

                return {
                  ...current,
                  connection: {
                    ...current.connection,
                    token: result.token,
                  },
                };
              },
            },
          ],
          () => cycleToken({}),
          { revalidate: true },
        );
      });
    });

  return (
    <LStack gap="4">
      <styled.p fontSize="sm" color="fg.muted">
        {t(
          "This plugin is an External plugin. This means Storyden does not manage its process lifecycle and cannot provide connection information. Use this token to connect the plugin to Storyden via RPC.",
        )}
      </styled.p>
      <styled.p fontSize="sm" color="fg.muted">
        {t(
          "External plugins are responsible for handling their own restarting and reconnection.",
        )}
      </styled.p>

      <Alert.Root>
        <Alert.Icon asChild>
          <WarningIcon />
        </Alert.Icon>
        <Alert.Content>
          <Alert.Title>{t("Security notice")}</Alert.Title>
          <Alert.Description>
            {t(
              "External plugin tokens grant full access over RPC. Keep them secret and only run plugins from trusted sources.",
            )}
          </Alert.Description>
        </Alert.Content>
      </Alert.Root>

      <LStack gap="1">
        <styled.p fontSize="xs" color="fg.muted">
          {t("Plugin token")}
        </styled.p>

        <Clipboard.Root w="full" value={plugin.connection.token}>
          <Clipboard.Control gap="0">
            <Clipboard.Input asChild>
              <Input size="sm" borderRightRadius="none" />
            </Clipboard.Input>
            <Clipboard.Trigger asChild>
              <IconButton size="sm" variant="subtle" borderLeftRadius="none">
                <Clipboard.Indicator copied={<CheckIcon />}>
                  <ClipboardIcon />
                </Clipboard.Indicator>
              </IconButton>
            </Clipboard.Trigger>
          </Clipboard.Control>
        </Clipboard.Root>
      </LStack>

      <LStack gap="1">
        <styled.p fontSize="xs" color="fg.muted">
          {t("Development environment variable")}
        </styled.p>
        <Clipboard.Root w="full" value={envURL}>
          <Clipboard.Control gap="0">
            <Clipboard.Input asChild>
              <Input size="sm" borderRightRadius="none" />
            </Clipboard.Input>
            <Clipboard.Trigger asChild>
              <IconButton size="sm" variant="subtle" borderLeftRadius="none">
                <Clipboard.Indicator copied={<CheckIcon />}>
                  <ClipboardIcon />
                </Clipboard.Indicator>
              </IconButton>
            </Clipboard.Trigger>
          </Clipboard.Control>
        </Clipboard.Root>
      </LStack>

      <HStack w="full" justify="end">
        <InfoTip title={t("Regenerating Plugin Token")}>
          <styled.p fontSize="sm" color="fg.muted">
            {t(
              "This will immediately invalidate the old token and force the plugin to disconnect if it's currently connected.",
            )}
          </styled.p>
        </InfoTip>

        {isConfirming ? (
          <HStack gap="2">
            <Button size="sm" variant="subtle" onClick={handleConfirmAction}>
              {t("Confirm regenerate")}
            </Button>
            <Button size="sm" variant="outline" onClick={handleCancelAction}>
              {t("Cancel")}
            </Button>
          </HStack>
        ) : (
          <Button
            size="sm"
            variant="subtle"
            flexShrink="0"
            onClick={handleConfirmAction}
          >
            {t("Regenerate token")}
          </Button>
        )}
      </HStack>
    </LStack>
  );
}

function buildRPCURL(token: string): string {
  try {
    const base = new URL(API_ADDRESS);
    base.protocol = base.protocol === "https:" ? "wss:" : "ws:";
    base.pathname = "/rpc";
    base.search = `token=${token}`;
    return base.toString();
  } catch {
    const host = API_ADDRESS.replace(/^https?:\/\//, "");
    return `ws://${host}/rpc?token=${token}`;
  }
}
