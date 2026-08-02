import { RobotActivityIcon } from "@/components/ui/icons/RobotActivity";
import { HStack } from "@/styled-system/jsx";

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
      color="text.subtle"
      fontSize="xs"
      px="1"
      gap="2"
    >
      <RobotActivityIcon size={18} aria-hidden="true" />
      <span>{robotName} is responding...</span>
    </HStack>
  );
}
