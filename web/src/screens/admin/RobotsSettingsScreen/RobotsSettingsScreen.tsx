import { useState } from "react";

import { useRobotsList } from "@/api/openapi-client/robots";
import { Robot } from "@/api/openapi-schema";
import { Unready } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { Text } from "@/components/ui/text";
import { CardBox, LStack } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { useDisclosure } from "@/utils/useDisclosure";

import { RobotCard } from "./RobotCard";
import { RobotConfigurationModal } from "./RobotConfigurationModal";

export function RobotsSettingsScreen() {
  const { data, error } = useRobotsList();
  const [selectedRobot, setSelectedRobot] = useState<Robot | null>(null);
  const disclosure = useDisclosure();

  if (!data) {
    return <Unready error={error} />;
  }

  function handleRobotClick(robot: Robot) {
    setSelectedRobot(robot);
    disclosure.onOpen();
  }

  function handleClose() {
    disclosure.onClose();
    setSelectedRobot(null);
  }

  return (
    <CardBox className={lstack()}>
      <Heading size="md">Robots</Heading>
      <Text color="fg.muted">
        Robots are language-model driven automations for organising your
        community.
      </Text>

      <LStack>
        {data.robots.map((robot) => (
          <RobotCard
            key={robot.id}
            robot={robot}
            onClick={() => handleRobotClick(robot)}
          />
        ))}
      </LStack>

      {selectedRobot && (
        <RobotConfigurationModal
          robotId={selectedRobot.id}
          isOpen={disclosure.isOpen}
          onClose={handleClose}
        />
      )}
    </CardBox>
  );
}
