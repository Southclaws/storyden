import { Category } from "src/api/openapi-schema";
import { useDisclosure } from "src/utils/useDisclosure";

import { MoreAction } from "@/components/site/Action/More";

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
