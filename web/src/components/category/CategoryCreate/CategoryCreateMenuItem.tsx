import { Category } from "src/api/openapi-schema";
import { useDisclosure } from "src/utils/useDisclosure";

import { Item } from "@/components/ui/menu";

import { CategoryCreateModal } from "./CategoryCreateModal";

type Props = {
  parentCategory?: Category;
};

export function CategoryCreateMenuItem({ parentCategory }: Props) {
  const useDisclosureProps = useDisclosure();

  return (
    <>
      <Item value="create-subcategory" onClick={useDisclosureProps.onOpen}>
        Create subcategory
      </Item>
      <CategoryCreateModal
        defaultParent={parentCategory?.id}
        {...useDisclosureProps}
      />
    </>
  );
}
