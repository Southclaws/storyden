import { formatDate } from "date-fns";
import { useQueryState } from "nuqs";

import { Plugin, PluginStatusError } from "@/api/openapi-schema";
import * as Alert from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { Text } from "@/components/ui/text";
import { Box, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

type Props = {
  plugin: Plugin;
};

export function OverviewTab({ plugin }: Props) {
  const [, setTab] = useQueryState("plugin-tab");
  const manifestID = manifestString(plugin.manifest, "id");
  const author = manifestString(plugin.manifest, "author");
  const command = manifestString(plugin.manifest, "command") || "-";
  const args = manifestArgs(plugin.manifest) || "-";
  const events = manifestStringList(plugin.manifest, "events_consumed");
  const statusError = getStatusError(plugin);
  const statusErrorDetails = statusError?.details ?? {};
  const hasStatusErrorDetails = Object.keys(statusErrorDetails).length > 0;
  const canViewLogs = plugin.connection.mode === "supervised";

  return (
    <LStack gap="4">
      {statusError && (
        <Alert.Root>
          <Alert.Icon asChild>
            <WarningIcon />
          </Alert.Icon>
          <Alert.Content w="full">
            <WStack>
              <Alert.Title>Plugin error</Alert.Title>
              {canViewLogs && (
                <HStack>
                  <Button variant="subtle" onClick={() => setTab("logs")}>
                    View logs
                  </Button>
                </HStack>
              )}
            </WStack>
            <Alert.Description>{statusError.message}</Alert.Description>

            {hasStatusErrorDetails && (
              <styled.details>
                <styled.summary fontSize="sm">Technical details</styled.summary>
                <styled.pre
                  mt="2"
                  p="2"
                  borderWidth="thin"
                  borderColor="border.default"
                  borderRadius="sm"
                  fontSize="xs"
                  overflowX="auto"
                >
                  {JSON.stringify(statusErrorDetails, null, 2)}
                </styled.pre>
              </styled.details>
            )}
          </Alert.Content>
        </Alert.Root>
      )}

      <Text variant="supporting" color="text.default">
        {plugin.description || "No description provided."}
      </Text>

      <LStack gap="2">
        <OverviewField label="Plugin ID" value={manifestID || "-"} monospace />
        <OverviewField label="Author" value={author || "-"} monospace />
        <OverviewField
          label="Version"
          value={plugin.version || "-"}
          monospace
        />

        <OverviewField label="Mode" value={plugin.connection.mode} monospace />

        <SubsectionTitle>Command</SubsectionTitle>
        <CommandLine command={command} args={args} />

        <SubsectionTitle>Events consumed</SubsectionTitle>
        {events.length === 0 ? (
          <Text variant="supporting">
            This plugin does not consume any events.
          </Text>
        ) : (
          <LStack gap="1" maxH="64" overflowY="auto">
            {events.map((eventName) => (
              <Text
                variant="metadata"
                key={eventName}
                fontFamily="mono"
                p="2"
                bgColor="background.inset"
                borderRadius="sm"
              >
                {eventName}
              </Text>
            ))}
          </LStack>
        )}
      </LStack>

      <WStack>
        <Text variant="metadata">
          installed:&nbsp;
          <styled.code color="text.subtle">
            {formatDate(plugin.added_at, "PPp")}
          </styled.code>
        </Text>

        <Text variant="metadata">
          id:&nbsp;<styled.code color="text.subtle">{plugin.id}</styled.code>
        </Text>
      </WStack>
    </LStack>
  );
}

function OverviewField({
  label,
  value,
  monospace,
}: {
  label: string;
  value: string;
  monospace?: boolean;
}) {
  return (
    <WStack justifyContent="space-between" alignItems="start" gap="4">
      <Text variant="metadata">{label}</Text>
      <Text
        variant="supporting"
        color="text.default"
        fontFamily={monospace ? "mono" : undefined}
        textAlign="right"
        wordBreak="break-word"
      >
        {value}
      </Text>
    </WStack>
  );
}

function SectionTitle({ children }: { children: string }) {
  return (
    <styled.h3 fontSize="sm" fontWeight="semibold" color="text.default">
      {children}
    </styled.h3>
  );
}

function SubsectionTitle({ children }: { children: string }) {
  return (
    <styled.h4 fontSize="sm" fontWeight="medium" color="text.muted">
      {children}
    </styled.h4>
  );
}

function CommandLine({ command, args }: { command: string; args: string }) {
  return (
    <Box
      borderWidth="thin"
      borderColor="border.default"
      borderRadius="md"
      bgColor="background.inset"
      overflow="hidden"
    >
      <WStack gap="0" alignItems="stretch">
        <Text
          variant="metadata"
          px="3"
          py="2"
          fontFamily="mono"
          fontWeight="semibold"
          bgColor="background.surface"
          borderRightWidth="thin"
          borderRightColor="border.default"
          whiteSpace="nowrap"
        >
          {command}
        </Text>
        <Text
          variant="metadata"
          px="3"
          py="2"
          fontFamily="mono"
          color={args === "-" ? "text.subtle" : "text.default"}
          wordBreak="break-word"
        >
          {args}
        </Text>
      </WStack>
    </Box>
  );
}

function manifestString(
  manifest: Record<string, unknown>,
  key: string,
): string | null {
  const value = manifest[key];
  return typeof value === "string" && value.trim() !== "" ? value : null;
}

function manifestArgs(manifest: Record<string, unknown>): string | null {
  const value = manifest["args"];
  if (Array.isArray(value)) {
    return value.map(String).join(" ");
  }
  return typeof value === "string" ? value : null;
}

function manifestStringList(
  manifest: Record<string, unknown>,
  key: string,
): string[] {
  const value = manifest[key];
  if (!Array.isArray(value)) {
    return [];
  }
  return value.filter((entry): entry is string => typeof entry === "string");
}

function getStatusError(plugin: Plugin): PluginStatusError | null {
  if (plugin.status.active_state !== "error") {
    return null;
  }
  return plugin.status;
}
