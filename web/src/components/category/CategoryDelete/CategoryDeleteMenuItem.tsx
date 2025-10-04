import { Category } from "src/api/openapi-schema";
import { useDisclosure } from "src/utils/useDisclosure";

import { Item } from "@/components/ui/menu";

import { CategoryDeleteModal } from "./CategoryDeleteModal";

export function CategoryDeleteMenuItem(props: Category) {
  const { onOpen, isOpen, onClose } = useDisclosure();

  return (
    <>
      <Item value="delete" onClick={onOpen}>
        Delete
      </Item>
      <CategoryDeleteModal
        onClose={onClose}
        isOpen={isOpen}
        categorySlug={props.slug}
        categoryName={props.name}
      />
    </>
  );
}
