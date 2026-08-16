import { RobotWorkspaceSelect } from "@/components/robots/RobotWorkspaceSelect";
import { IconButton } from "@/components/ui/icon-button";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { HStack } from "@/styled-system/jsx";

import { useCommandPalette } from "../Context";

import { useRobotChat } from "./RobotChatContext";
import { RobotSessionMenu } from "./RobotSessionMenu";

export function RobotCommandPaletteStatusBar() {
  const { resetChatSession } = useCommandPalette();
  const { cancelActiveTurn, canCancelActiveTurn, isCancelling } =
    useRobotChat();

  function handleReset() {
    resetChatSession();
  }

  return (
    <>
      <RobotSessionMenu />

      <HStack gap="0">
        <RobotWorkspaceSelect variant="ghost" minW="40" />
        {canCancelActiveTurn && (
          <IconButton
            aria-label="Cancel Robot response"
            variant="ghost"
            loading={isCancelling}
            onClick={() => void cancelActiveTurn()}
          >
            <CancelIcon />
          </IconButton>
        )}
        <IconButton
          aria-label="Start a new Robot chat"
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
