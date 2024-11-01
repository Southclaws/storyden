import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";

import { CollectionCreateScreen } from "./CollectionCreateScreen";
import { Props } from "./useCollectionCreate";

export function CollectionCreateModal({ session, ...props }: Props) {
  return (
    <>
      <ModalDrawer
        isOpen={props.isOpen}
        onClose={props.onClose}
        title="Create collection"
      >
        <CollectionCreateScreen
          id={props.id}
          session={session}
          onClose={props.onClose}
        />
      </ModalDrawer>
    </>
  );
}
