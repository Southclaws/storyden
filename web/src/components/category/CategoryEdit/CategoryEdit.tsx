import { Category } from "src/api/openapi-schema";
import { useDisclosure } from "src/utils/useDisclosure";

import { MoreAction } from "@/components/site/Action/More";
import { Item } from "@/components/ui/menu";

import { CategoryEditModal } from "./CategoryEditModal";

export function CategoryEdit(props: Category) {
  const { onOpen, isOpen, onClose } = useDisclosure();
  return (
    <>
      <MoreAction size="xs" onClick={onOpen} />

      <CategoryEditModal onClose={onClose} isOpen={isOpen} category={props} />
    </>
  );
}

export function CategoryEditMenuItem(props: Category) {
  const { onOpen, isOpen, onClose } = useDisclosure();

  return (
    <Item value="edit" onClick={onOpen}>
      Edit
      <CategoryEditModal onClose={onClose} isOpen={isOpen} category={props} />
    </Item>
  );
}
