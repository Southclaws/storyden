import { Heading } from "@/components/ui/heading";
import { Text } from "@/components/ui/text";
import { CardBox } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { RobotConfigurationScreen } from "./RobotConfigurationScreen";
import RobotList from "./RobotList";
import { useSelectedRobot } from "./useSelectedRobot";

export function RobotsSettingsScreen() {
  const [robotID] = useSelectedRobot();

  return (
    <CardBox className={lstack()}>
      {robotID ? (
        <RobotConfigurationScreen robotId={robotID} />
      ) : (
        <>
          <Heading size="md">Robots</Heading>
          <Text color="fg.muted">
            Robots are language-model driven automations for organising your
            community.
          </Text>

          <RobotList />
        </>
      )}
    </CardBox>
  );
}
