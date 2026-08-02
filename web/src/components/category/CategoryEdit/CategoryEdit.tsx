import { Category } from "@/api/openapi-schema";
import { Item } from "@/components/ui/menu";
import { useDisclosure } from "@/utils/useDisclosure";

import { CategoryEditModal } from "./CategoryEditModal";

export function CategoryEditMenuItem(props: Category) {
  const disclosure = useDisclosure();

  return (
    <>
      <Item value="edit" onClick={disclosure.onOpen}>
        Edit
      </Item>
      <CategoryEditModal {...disclosure} category={props} />
    </>
  );
}
