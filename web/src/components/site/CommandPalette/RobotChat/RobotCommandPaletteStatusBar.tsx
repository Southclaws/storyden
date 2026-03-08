import { RobotListMenu } from "@/components/robots/RobotListMenu";
import { Button } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { HStack } from "@/styled-system/jsx";

import { useCommandPalette } from "../Context";

import { RobotSessionMenu } from "./RobotSessionMenu";

export function RobotCommandPaletteStatusBar() {
  const { resetChatSession } = useCommandPalette();

  function handleReset() {
    resetChatSession();
  }

  return (
    <>
      <RobotSessionMenu />

      <HStack gap="0">
        <RobotListMenu size="xs" variant="ghost" borderRightRadius="none" />
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
