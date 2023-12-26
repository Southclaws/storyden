import { Category } from "src/api/openapi/schemas";
import { EditAction } from "src/components/site/Action/Edit";
import { useDisclosure } from "src/utils/useDisclosure";

import { CategoryEditModal } from "./CategoryEditModal";

export function CategoryEdit(props: Category) {
  const { onOpen, isOpen, onClose } = useDisclosure();
  return (
    <>
      <EditAction onClick={onOpen} />

      <CategoryEditModal onClose={onClose} isOpen={isOpen} category={props} />
    </>
  );
}
