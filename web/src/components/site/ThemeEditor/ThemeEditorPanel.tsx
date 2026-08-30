"use client";

import * as Alert from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import * as FloatingPanel from "@/components/ui/floating-panel";
import { AddIcon } from "@/components/ui/icons/Add";
import { ArrowDownIcon, ArrowUpIcon } from "@/components/ui/icons/Arrow";
import { CloseIcon } from "@/components/ui/icons/Close";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { WarningIcon } from "@/components/ui/icons/Warning";
import * as Tabs from "@/components/ui/tabs";
import { Text } from "@/components/ui/text";
import { Textarea } from "@/components/ui/textarea";
import { HStack, LStack, styled } from "@/styled-system/jsx";

import type {
  ThemeEditorDocument,
  ThemeEditorKind,
} from "./theme-editor-model";

type Props = {
  documents: ThemeEditorDocument[] | null;
  selectedKey?: string;
  dirty: boolean;
  busy: boolean;
  runtimeDisabled: boolean;
  error?: string;
  panelSize?: { width: number; height: number };
  panelPosition?: { x: number; y: number };
  onSelect: (key: string) => void;
  onAdd: (kind: ThemeEditorKind) => void;
  onChange: (key: string, source: string) => void;
  onMove: (key: string, offset: -1 | 1) => void;
  onDelete: (key: string) => void;
  onSave: () => void;
  onExit: () => void;
};

export function ThemeEditorPanel(props: Props) {
  const size = props.panelSize ?? { width: 540, height: 560 };
  const selected = props.documents?.find(
    ({ key }) => key === props.selectedKey,
  );

  return (
    <FloatingPanel.Root
      open
      allowOverflow={false}
      closeOnEscape
      defaultSize={size}
      defaultPosition={props.panelPosition ?? { x: 16, y: 16 }}
      minSize={{ width: 320, height: 300 }}
      maxSize={{ width: 960, height: 960 }}
      restoreFocus={false}
      onOpenChange={({ open }) => {
        if (!open) props.onExit();
      }}
    >
      <FloatingPanel.Positioner data-sd-theme-editor>
        <FloatingPanel.Content>
          <FloatingPanel.DragTrigger>
            <FloatingPanel.Header>
              <HStack gap="2" minWidth="0">
                <FloatingPanel.Title>Theme editor</FloatingPanel.Title>
                <Text variant="metadata" color="text.muted">
                  {props.busy ? "Saving…" : props.dirty ? "Unsaved" : "Saved"}
                </Text>
              </HStack>
              <FloatingPanel.Control>
                <FloatingPanel.CloseTrigger asChild>
                  <Button
                    aria-label="Exit theme editing"
                    size="sm"
                    variant="ghost"
                  >
                    <CloseIcon />
                  </Button>
                </FloatingPanel.CloseTrigger>
              </FloatingPanel.Control>
            </FloatingPanel.Header>
          </FloatingPanel.DragTrigger>

          <FloatingPanel.Body>
            {props.runtimeDisabled && (
              <Alert.Root tone="warning" margin="2" marginBottom="0">
                <Alert.Icon asChild>
                  <WarningIcon />
                </Alert.Icon>
                <Alert.Content>
                  <Alert.Title>Runtime themes are disabled</Alert.Title>
                  <Alert.Description>
                    CUSTOM_THEMES_DISABLE must be cleared before changes can be
                    published.
                  </Alert.Description>
                </Alert.Content>
              </Alert.Root>
            )}

            {props.documents === null ? (
              <LStack flex="1" alignItems="center" justifyContent="center">
                <Text>Loading the live theme…</Text>
              </LStack>
            ) : props.documents.length === 0 ? (
              <LStack
                flex="1"
                gap="3"
                alignItems="center"
                justifyContent="center"
                padding="6"
                textAlign="center"
              >
                <Text fontWeight="semibold">Start with CSS or JavaScript</Text>
                <Text variant="supporting" maxWidth="sm">
                  Add a tab, paste code, then save it directly to the live site.
                </Text>
                <HStack gap="2">
                  <Button onClick={() => props.onAdd("stylesheet")}>
                    <AddIcon /> Add CSS
                  </Button>
                  <Button
                    variant="outline"
                    onClick={() => props.onAdd("script")}
                  >
                    <AddIcon /> Add JavaScript
                  </Button>
                </HStack>
              </LStack>
            ) : (
              <Tabs.Root
                value={props.selectedKey}
                onValueChange={({ value }) => props.onSelect(value)}
                flex="1"
                minHeight="0"
              >
                <HStack
                  borderBottomColor="border.default"
                  borderBottomWidth="thin"
                  gap="2"
                  paddingInline="2"
                >
                  <Tabs.List flex="1" minWidth="0">
                    {props.documents.map((document) => (
                      <Tabs.Trigger key={document.key} value={document.key}>
                        {document.label}
                      </Tabs.Trigger>
                    ))}
                    <Tabs.Indicator />
                  </Tabs.List>
                  <Button
                    aria-label="Add CSS"
                    size="sm"
                    variant="ghost"
                    onClick={() => props.onAdd("stylesheet")}
                  >
                    <AddIcon /> CSS
                  </Button>
                  <Button
                    aria-label="Add JavaScript"
                    size="sm"
                    variant="ghost"
                    onClick={() => props.onAdd("script")}
                  >
                    <AddIcon /> JS
                  </Button>
                </HStack>

                {props.documents.map((document) => (
                  <Tabs.Content
                    key={document.key}
                    value={document.key}
                    display="flex"
                    flex="1"
                    flexDirection="column"
                    minHeight="0"
                    padding="2"
                    paddingTop="2"
                  >
                    <HStack gap="1" justifyContent="end" paddingBottom="2">
                      <Button
                        aria-label={`Move ${document.label} earlier`}
                        size="sm"
                        variant="ghost"
                        disabled={props.busy}
                        onClick={() => props.onMove(document.key, -1)}
                      >
                        <ArrowUpIcon />
                      </Button>
                      <Button
                        aria-label={`Move ${document.label} later`}
                        size="sm"
                        variant="ghost"
                        disabled={props.busy}
                        onClick={() => props.onMove(document.key, 1)}
                      >
                        <ArrowDownIcon />
                      </Button>
                      <Button
                        aria-label={`Delete ${document.label}`}
                        size="sm"
                        variant="ghost"
                        intent="destructive"
                        disabled={props.busy}
                        onClick={() => props.onDelete(document.key)}
                      >
                        <DeleteIcon />
                      </Button>
                    </HStack>
                    <Textarea
                      aria-label={`${document.label} source`}
                      autoCapitalize="off"
                      autoCorrect="off"
                      flex="1"
                      fontFamily="mono"
                      minHeight="0"
                      resize="none"
                      spellCheck={false}
                      value={document.source}
                      variant="inset"
                      onChange={(event) =>
                        props.onChange(document.key, event.currentTarget.value)
                      }
                    />
                    {document.kind === "script" && (
                      <Text
                        variant="metadata"
                        color="text.muted"
                        paddingTop="2"
                      >
                        Saved JavaScript runs with Storyden&apos;s browser
                        privileges and reloads this page.
                      </Text>
                    )}
                  </Tabs.Content>
                ))}
              </Tabs.Root>
            )}

            <styled.footer
              borderTopColor="border.default"
              borderTopWidth="thin"
              padding="2"
            >
              <HStack gap="3" justifyContent="space-between">
                <Text
                  variant="metadata"
                  color={props.error ? "text.danger" : "text.muted"}
                  lineClamp="2"
                >
                  {props.error ??
                    (selected
                      ? `Editing ${selected.label} on this page`
                      : "Changes apply to every Storyden page")}
                </Text>
                <Button
                  variant="solid"
                  loading={props.busy}
                  disabled={!props.dirty || props.busy || props.runtimeDisabled}
                  onClick={props.onSave}
                >
                  Save live theme
                </Button>
              </HStack>
            </styled.footer>
          </FloatingPanel.Body>

          <FloatingPanel.ResizeTrigger axis="e" />
          <FloatingPanel.ResizeTrigger axis="s" />
          <FloatingPanel.ResizeTrigger axis="se" />
        </FloatingPanel.Content>
      </FloatingPanel.Positioner>
    </FloatingPanel.Root>
  );
}
