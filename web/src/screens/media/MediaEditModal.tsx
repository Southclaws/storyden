import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";
import { UseDisclosureProps } from "src/utils/useDisclosure";

import { MediaEditScreen, Props } from "./MediaEditScreen";

export function MediaEditModal(props: UseDisclosureProps & Props) {
  return (
    <>
      <ModalDrawer isOpen={props.isOpen} onClose={props.onClose} title="Upload">
        <MediaEditScreen {...props} />
      </ModalDrawer>
    </>
  );
}
