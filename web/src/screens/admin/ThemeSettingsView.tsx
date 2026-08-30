"use client";

import type { ThemeAsset } from "@/api/openapi-schema";
import * as Alert from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { PageHeader } from "@/components/ui/page-header";
import { StatusBadge } from "@/components/ui/status-badge";
import * as Table from "@/components/ui/table";
import { Text } from "@/components/ui/text";
import { HStack, LStack, styled } from "@/styled-system/jsx";

type Props = {
  assets: ThemeAsset[];
  active: boolean;
  runtimeDisabled: boolean;
  editingEnabled: boolean;
  busy?: boolean;
  error?: string;
  onEnableEditing?: () => void;
  onExitEditing?: () => void;
  onDisable?: () => void;
};

export function ThemeSettingsView(props: Props) {
  const { assets, active, runtimeDisabled, editingEnabled, busy, error } =
    props;
  const scriptCount = assets.filter(
    ({ mime_type }) => mime_type === "application/javascript",
  ).length;

  return (
    <LStack gap="6" maxW="layout.contentWide">
      <PageHeader
        title="Custom theme"
        description="Edit installation-wide CSS and trusted JavaScript directly on the real Storyden interface."
        badge={
          <StatusBadge
            tone={
              runtimeDisabled
                ? "warning"
                : editingEnabled
                  ? "info"
                  : active
                    ? "success"
                    : "neutral"
            }
          >
            {runtimeDisabled
              ? "Runtime disabled"
              : editingEnabled
                ? "Editing"
                : active
                  ? "Active"
                  : "Disabled"}
          </StatusBadge>
        }
      />

      <Alert.Root tone="danger">
        <Alert.Icon asChild>
          <WarningIcon />
        </Alert.Icon>
        <Alert.Content>
          <Alert.Title>Theme code is trusted administrator code</Alert.Title>
          <Alert.Description>
            Saved scripts execute with Storyden&apos;s browser privileges for
            every visitor, including administrators. They are not sandboxed and
            may access authenticated pages or initiate network requests.
          </Alert.Description>
        </Alert.Content>
      </Alert.Root>

      {runtimeDisabled && (
        <Alert.Root tone="warning">
          <Alert.Icon asChild>
            <WarningIcon />
          </Alert.Icon>
          <Alert.Content>
            <Alert.Title>Disabled by deployment configuration</Alert.Title>
            <Alert.Description>
              The configured manifest is preserved, but CUSTOM_THEMES_DISABLE
              prevents it from loading or being edited live.
            </Alert.Description>
          </Alert.Content>
        </Alert.Root>
      )}

      {error && (
        <Alert.Root tone="danger">
          <Alert.Icon asChild>
            <WarningIcon />
          </Alert.Icon>
          <Alert.Content>
            <Alert.Title>Theme operation failed</Alert.Title>
            <Alert.Description>{error}</Alert.Description>
          </Alert.Content>
        </Alert.Root>
      )}

      <section>
        <LStack gap="3">
          <div>
            <styled.h2 fontWeight="semibold">Live editing</styled.h2>
            <Text variant="supporting" maxW="2xl">
              Enable the floating editor, then browse anywhere in Storyden while
              it follows you. Inspect the real page with browser developer
              tools, edit CSS or JavaScript tabs, and save changes directly to
              the live theme.
            </Text>
          </div>
          <CardBox>
            <HStack
              alignItems={{ base: "stretch", md: "center" }}
              flexDirection={{ base: "column", md: "row" }}
              gap="4"
              justifyContent="space-between"
            >
              <LStack gap="1">
                <Text fontWeight="semibold">
                  {editingEnabled
                    ? "The theme editor is active in this browser"
                    : "Open the editor on the current page"}
                </Text>
                <Text variant="supporting">
                  The editing flag is stored locally and never changes the SSR
                  response for other visitors.
                </Text>
              </LStack>
              {editingEnabled ? (
                <Button variant="outline" onClick={props.onExitEditing}>
                  Exit theme editing
                </Button>
              ) : (
                <Button
                  variant="solid"
                  disabled={runtimeDisabled}
                  onClick={props.onEnableEditing}
                >
                  Enable theme editing
                </Button>
              )}
            </HStack>
          </CardBox>
        </LStack>
      </section>

      <section>
        <LStack gap="3">
          <div>
            <styled.h2 fontWeight="semibold">Published assets</styled.h2>
            <Text variant="supporting">
              {assets.length === 0
                ? "No custom CSS or JavaScript is currently published."
                : `${assets.length} immutable asset${assets.length === 1 ? " is" : "s are"} active${scriptCount > 0 ? `, including ${scriptCount} trusted script${scriptCount === 1 ? "" : "s"}` : ""}.`}
            </Text>
          </div>
          <ThemeAssetOverview assets={assets} />
        </LStack>
      </section>

      <HStack gap="3" flexWrap="wrap">
        <Button
          variant="outline"
          intent="destructive"
          disabled={!active || busy}
          loading={busy}
          onClick={props.onDisable}
        >
          Disable live theme
        </Button>
        <Text variant="supporting">
          Disabling removes the live manifest without deleting its immutable
          assets.
        </Text>
      </HStack>
    </LStack>
  );
}

function ThemeAssetOverview({ assets }: { assets: ThemeAsset[] }) {
  if (assets.length === 0) {
    return (
      <CardBox>
        <Text>The default Storyden theme is active.</Text>
      </CardBox>
    );
  }

  return (
    <CardBox overflowX="auto" padding="0">
      <Table.Root size="sm" width="full">
        <Table.Head>
          <Table.Row>
            <Table.Header>Order</Table.Header>
            <Table.Header>Type</Table.Header>
            <Table.Header>Size</Table.Header>
            <Table.Header>Integrity</Table.Header>
          </Table.Row>
        </Table.Head>
        <Table.Body>
          {assets.map((asset, index) => (
            <Table.Row key={asset.id}>
              <Table.Cell>{index + 1}</Table.Cell>
              <Table.Cell>
                {asset.mime_type === "text/css" ? "CSS" : "JavaScript"}
              </Table.Cell>
              <Table.Cell>{formatBytes(asset.size)}</Table.Cell>
              <Table.Cell>
                <Text as="code" variant="metadata">
                  {asset.integrity.slice(0, 28)}…
                </Text>
              </Table.Cell>
            </Table.Row>
          ))}
        </Table.Body>
      </Table.Root>
    </CardBox>
  );
}

function formatBytes(size: number) {
  return size < 1024 ? `${size} B` : `${(size / 1024).toFixed(1)} KiB`;
}
