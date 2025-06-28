import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";

import { CreateAccessKeyScreen } from "./CreateAccessKeyScreen";
import { Props, WithDisclosure } from "./useCreateAccessKeyScreen";

export function CreateAccessKeyModal(props: WithDisclosure<Props>) {
  return (
    <ModalDrawer
      isOpen={props.isOpen}
      onClose={props.onClose}
      title="Create Access Key"
    >
      <CreateAccessKeyScreen onClose={props.onClose || (() => {})} />
    </ModalDrawer>
  );
}
