import { useRobotGet } from "@/api/openapi-client/robots";
import { RobotConfigurationForm } from "@/components/robots/RobotConfiguration/RobotConfigurationForm";
import { BackAction } from "@/components/site/Action/Back";
import { UnreadyBanner } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { HStack, LStack } from "@/styled-system/jsx";

import { useSelectedRobot } from "./useSelectedRobot";

type Props = {
  robotId: string;
};

export function RobotConfigurationScreen({ robotId }: Props) {
  const { data, error } = useRobotGet(robotId);
  const [_, setSelectedRobot] = useSelectedRobot();

  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  return (
    <LStack>
      <HStack>
        <BackAction onClick={() => void setSelectedRobot(null)} />
        <Heading size="md" lineClamp="1">
          {data.name}
        </Heading>
      </HStack>

      <RobotConfigurationForm robot={data} />
    </LStack>
  );
}
