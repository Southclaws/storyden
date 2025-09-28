import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";

import { CategoryCreateProps, CategoryCreateScreen } from "./CategoryCreateScreen";

export function CategoryCreateModal(props: CategoryCreateProps) {
  return (
    <>
      <ModalDrawer
        isOpen={props.isOpen}
        onClose={props.onClose}
        title="Create category"
      >
        <CategoryCreateScreen {...props} />
      </ModalDrawer>
    </>
  );
}
