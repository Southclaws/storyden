import { Category } from "src/api/openapi-schema";
import { useDisclosure } from "src/utils/useDisclosure";

import { Item } from "@/components/ui/menu";

import { CategoryEditModal } from "./CategoryEditModal";

export function CategoryEditMenuItem(props: Category) {
  const { onOpen, isOpen, onClose } = useDisclosure();

  return (
    <Item value="edit" onClick={onOpen}>
      Edit
      <CategoryEditModal onClose={onClose} isOpen={isOpen} category={props} />
    </Item>
  );
}
