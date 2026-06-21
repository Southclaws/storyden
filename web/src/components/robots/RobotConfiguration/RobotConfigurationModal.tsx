import { useRobotGet } from "@/api/openapi-client/robots";
import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { Unready } from "@/components/site/Unready";
import { UseDisclosureProps } from "@/utils/useDisclosure";

import { RobotConfigurationForm } from "./RobotConfigurationForm";

type Props = {
  robotId: string;
} & UseDisclosureProps;

export function RobotConfigurationModal({
  robotId,
  onClose,
  onOpen,
  isOpen,
}: Props) {
  const { data, error } = useRobotGet(robotId, {
    swr: { enabled: isOpen },
  });

  return (
    <ModalDrawer
      onOpen={onOpen}
      isOpen={isOpen}
      onClose={onClose}
      title={data ? `Configure: ${data.name}` : "Loading..."}
    >
      {!data ? (
        <Unready error={error} />
      ) : (
        <RobotConfigurationForm robot={data} onSave={onClose} />
      )}
    </ModalDrawer>
  );
}
