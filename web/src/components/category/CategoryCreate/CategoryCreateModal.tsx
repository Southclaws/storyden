import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";
import { UseDisclosureProps } from "src/utils/useDisclosure";

import { CategoryCreateScreen } from "./CategoryCreateScreen";

export function CategoryCreateModal(props: UseDisclosureProps) {
  return (
    <>
      <ModalDrawer
        isOpen={props.isOpen}
        onClose={props.onClose}
        title="Create category"
      >
        <CategoryCreateScreen onClose={props.onClose} id={props.id} />
      </ModalDrawer>
    </>
  );
}
