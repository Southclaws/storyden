import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";

import { CategoryDeleteProps, CategoryDeleteScreen } from "./CategoryDeleteScreen";

export function CategoryDeleteModal(props: CategoryDeleteProps) {
  return (
    <ModalDrawer
      isOpen={props.isOpen}
      onClose={props.onClose}
      title="Delete category"
    >
      <CategoryDeleteScreen {...props} />
    </ModalDrawer>
  );
}