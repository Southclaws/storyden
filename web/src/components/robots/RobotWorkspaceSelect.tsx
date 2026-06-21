"use client";

import { SelectValueChangeDetails, createListCollection } from "@ark-ui/react";
import type { ComponentProps } from "react";
import { useMemo } from "react";

import { useRobotChat } from "@/components/site/CommandPalette/RobotChat/RobotChatContext";
import { CheckIcon } from "@/components/ui/icons/Check";
import { SelectIcon } from "@/components/ui/icons/Select";
import * as Select from "@/components/ui/select";
import { styled } from "@/styled-system/jsx";

const NO_WORKSPACE_VALUE = "__none";

type Props = {
  size?: "xs" | "sm" | "md" | "lg";
  variant?: "outline" | "ghost";
  minW?: ComponentProps<typeof Select.Root>["minW"];
};

export function RobotWorkspaceSelect({
  size = "xs",
  variant = "ghost",
  minW = "44",
}: Props) {
  const {
    selectedWorkspaceID,
    setSelectedWorkspaceID,
    workspaces,
    workspacesReady,
  } = useRobotChat();

  const collection = useMemo(
    () =>
      createListCollection({
        items: [
          { label: "No workspace", value: NO_WORKSPACE_VALUE },
          ...workspaces.map((workspace) => ({
            label: workspace.name,
            value: workspace.id,
            description: workspace.provider,
          })),
        ],
      }),
    [workspaces],
  );

  const value = selectedWorkspaceID ?? NO_WORKSPACE_VALUE;
  const isDisabled = !workspacesReady || workspaces.length === 0;
  const placeholder = workspacesReady ? "Select workspace" : "Loading...";

  function handleWorkspaceChange({ value }: SelectValueChangeDetails) {
    const [selected] = value;
    if (!selected || selected === NO_WORKSPACE_VALUE) {
      setSelectedWorkspaceID(undefined);
      return;
    }

    setSelectedWorkspaceID(selected);
  }

  return (
    <Select.Root
      size={size}
      variant={variant}
      collection={collection}
      value={[value]}
      positioning={{ sameWidth: false }}
      onValueChange={handleWorkspaceChange}
      disabled={isDisabled}
      minW={minW}
      width="auto"
    >
      <Select.Control>
        <Select.Trigger>
          <Select.ValueText placeholder={placeholder} />
          <SelectIcon />
        </Select.Trigger>
      </Select.Control>
      <Select.Positioner>
        <Select.Content minW={minW}>
          {collection.items.map((item) => (
            <Select.Item key={item.value} item={item}>
              <Select.ItemText>
                {item.label}
                {"description" in item && item.description ? (
                  <styled.span color="fg.muted" ml="1">
                    {item.description}
                  </styled.span>
                ) : null}
              </Select.ItemText>
              <Select.ItemIndicator>
                <CheckIcon />
              </Select.ItemIndicator>
            </Select.Item>
          ))}
        </Select.Content>
      </Select.Positioner>
    </Select.Root>
  );
}
