import { MenuSelectionDetails } from "@ark-ui/react";
import { ComponentProps } from "react";

import { useRobotChat } from "@/components/site/CommandPalette/RobotChat/RobotChatContext";
import { Button } from "@/components/ui/button";
import * as Menu from "@/components/ui/menu";

const STORYDEN_BUILDER_ROBOT_ID = "storyden-builder-robot";

type Props = ComponentProps<typeof Button>;

export function RobotListMenu({
  size = "xs",
  variant = "ghost",
  ...buttonProps
}: Props) {
  const { setSelectedRobot, selectedRobot, robots } = useRobotChat();

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
    <Menu.Root onSelect={handleSelectRobot}>
      <Menu.Trigger asChild>
        <Button size={size} variant={variant} {...buttonProps}>
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
  );
}
