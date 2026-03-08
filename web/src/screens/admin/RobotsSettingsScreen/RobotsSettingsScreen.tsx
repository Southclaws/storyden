import { Heading } from "@/components/ui/heading";
import { Text } from "@/components/ui/text";
import { CardBox } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import RobotListScreen from "./RobotListScreen";

export function RobotsSettingsScreen() {
  return (
    <CardBox className={lstack()}>
      <Heading size="md">Robots</Heading>
      <Text color="fg.muted">
        Robots are language-model driven automations for organising your
        community.
      </Text>

      <RobotListScreen />
    </CardBox>
  );
}
