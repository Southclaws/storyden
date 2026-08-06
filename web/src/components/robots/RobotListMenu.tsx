import { MenuSelectionDetails } from "@ark-ui/react";
import { ComponentProps } from "react";

import {
  BUILT_IN_ROBOTS,
  DEFAULT_ROBOT_NAME,
  useRobotChat,
} from "@/components/site/CommandPalette/RobotChat/RobotChatContext";
import { Button } from "@/components/ui/button";
import * as Menu from "@/components/ui/menu";
import { Menu as ClosedMenu } from "@/components/ui/menu";

type Props = ComponentProps<typeof Button>;

export function RobotListMenu({
  size = "sm",
  variant = "ghost",
  ...buttonProps
}: Props) {
  const { setSelectedRobot, selectedRobot, robots } = useRobotChat();

  function handleSelectRobot({ value }: MenuSelectionDetails) {
    const selected =
      BUILT_IN_ROBOTS.find((r) => r.id === value) ??
      robots.find((r) => r.id === value);
    if (!selected) {
      return;
    }

    setSelectedRobot(selected);
  }

  const selectedRobotLabel = selectedRobot
    ? selectedRobot?.name
    : DEFAULT_ROBOT_NAME;

  return (
    <ClosedMenu
      contentProps={{ minW: "48", userSelect: "none" }}
      onSelect={handleSelectRobot}
      trigger={
        <Button size={size} variant={variant} {...buttonProps}>
          {selectedRobotLabel}
        </Button>
      }
    >
      <Menu.ItemGroup>
        <Menu.ItemGroupLabel>Select Robot</Menu.ItemGroupLabel>
        {BUILT_IN_ROBOTS.map((r) => {
          return (
            <Menu.Item key={r.id} value={r.id}>
              {r.name}
            </Menu.Item>
          );
        })}

        {robots.map((r) => {
          return (
            <Menu.Item key={r.id} value={r.id}>
              {r.name}
            </Menu.Item>
          );
        })}
      </Menu.ItemGroup>
    </ClosedMenu>
  );
}
