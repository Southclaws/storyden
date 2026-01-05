import { MenuSelectionDetails } from "@ark-ui/react";

import { Button } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import * as Menu from "@/components/ui/menu";
import { HStack, styled } from "@/styled-system/jsx";

import { useRobotChat } from "./RobotChatContext";

const STORYDEN_BUILDER_ROBOT_ID = "storyden-builder-robot";

export function RobotCommandPaletteStatusBar() {
  const { sessionId, setSelectedRobot, selectedRobot, robots, handleReset } =
    useRobotChat();

  function handleSelectRobot({ value }: MenuSelectionDetails) {
    if (value === STORYDEN_BUILDER_ROBOT_ID) {
      setSelectedRobot(undefined);
      return;
    }

    const selected = robots.find((r) => r.id === value);
    if (!selected) {
      return;
    }

    setSelectedRobot(selected);
  }

  const selectedRobotLabel = selectedRobot
    ? selectedRobot?.name
    : "Storyden Builder Robot";

  return (
    <>
      <styled.p>
        Robot Chat <code>session: {sessionId}</code>
      </styled.p>

      <HStack gap="0">
        <Menu.Root onSelect={handleSelectRobot}>
          <Menu.Trigger asChild>
            <Button size="xs" variant="ghost" borderRightRadius="none">
              {selectedRobotLabel}
            </Button>
          </Menu.Trigger>

          <Menu.Positioner>
            <Menu.Content minW="48" userSelect="none">
              <Menu.ItemGroup>
                <Menu.ItemGroupLabel>Select Robot</Menu.ItemGroupLabel>

                <Menu.Item value={STORYDEN_BUILDER_ROBOT_ID}>
                  Storyden Robot Builder
                </Menu.Item>

                {robots.map((r) => {
                  return (
                    <Menu.Item key={r.id} value={r.id}>
                      {r.name}
                    </Menu.Item>
                  );
                })}
              </Menu.ItemGroup>
            </Menu.Content>
          </Menu.Positioner>
        </Menu.Root>
        <IconButton
          size="xs"
          variant="ghost"
          borderLeftRadius="none"
          onClick={handleReset}
        >
          <CancelIcon />
        </IconButton>
      </HStack>
    </>
  );
}
