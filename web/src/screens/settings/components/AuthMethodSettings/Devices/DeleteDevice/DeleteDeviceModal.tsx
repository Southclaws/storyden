import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";

import { RoleCreateScreen } from "./DeleteDeviceScreen";
import { Props, WithDisclosure } from "./useDeleteDeviceScreen";

export function DeleteDeviceModal(props: WithDisclosure<Props>) {
  return (
    <>
      <ModalDrawer
        isOpen={props.isOpen}
        onClose={props.onClose}
        title="Delete device"
      >
        <RoleCreateScreen onClose={props.onClose} id={props.id} />
      </ModalDrawer>
    </>
  );
}
