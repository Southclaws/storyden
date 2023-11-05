import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";
import { UseDisclosureProps } from "src/theme/components";

import { CollectionCreateScreen } from "./CollectionCreateScreen";

export function CollectionCreateModal(props: UseDisclosureProps) {
  return (
    <>
      <ModalDrawer
        isOpen={props.isOpen}
        onClose={props.onClose}
        title="Create collection"
      >
        <CollectionCreateScreen onClose={props.onClose} id={props.id} />
      </ModalDrawer>
    </>
  );
}
