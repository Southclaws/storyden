import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";

import { CreateOAuthClientScreen } from "./CreateOAuthClientScreen";
import { Props, WithDisclosure } from "./useCreateOAuthClientScreen";

export function CreateOAuthClientModal(props: WithDisclosure<Props>) {
  return (
    <ModalDrawer
      isOpen={props.isOpen}
      onClose={props.onClose}
      title="Create OAuth Client"
    >
      <CreateOAuthClientScreen onClose={props.onClose || (() => {})} />
    </ModalDrawer>
  );
}
