import { RobotActivityIcon } from "@/components/ui/icons/RobotActivity";
import { HStack, styled } from "@/styled-system/jsx";

export function RobotChatLoadingStatus({
  active,
  robotName,
}: {
  active: boolean;
  robotName: string;
}) {
  if (!active) {
    return null;
  }

  return (
    <HStack
      role="status"
      aria-live="polite"
      color="fg.muted"
      fontSize="xs"
      px="1"
      gap="2"
    >
      <RobotActivityIcon size={18} aria-hidden="true" />
      <styled.span>{robotName} is responding...</styled.span>
    </HStack>
  );
}
