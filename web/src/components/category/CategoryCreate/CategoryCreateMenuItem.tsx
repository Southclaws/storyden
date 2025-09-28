import { Category } from "src/api/openapi-schema";
import { useDisclosure } from "src/utils/useDisclosure";

import { Item } from "@/components/ui/menu";

import { CategoryCreateModal } from "./CategoryCreateModal";

type Props = {
  parentCategory?: Category;
};

export function CategoryCreateMenuItem({ parentCategory }: Props) {
  const { onOpen, isOpen, onClose } = useDisclosure();

  return (
    <Item value="create-subcategory" onClick={onOpen}>
      Create subcategory
      <CategoryCreateModal
        onClose={onClose}
        isOpen={isOpen}
        defaultParent={parentCategory?.id}
      />
    </Item>
  );
}