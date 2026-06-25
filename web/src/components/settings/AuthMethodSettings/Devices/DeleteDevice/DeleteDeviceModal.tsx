import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";

import { DeleteDeviceScreen } from "./DeleteDeviceScreen";
import { Props, WithDisclosure } from "./useDeleteDeviceScreen";

export function DeleteDeviceModal(props: WithDisclosure<Props>) {
  return (
    <>
      <ModalDrawer
        isOpen={props.isOpen}
        onClose={props.onClose}
        title="Delete device"
      >
        <DeleteDeviceScreen onClose={props.onClose} id={props.id} />
      </ModalDrawer>
    </>
  );
}
