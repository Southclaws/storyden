import { useDisclosure } from "@chakra-ui/react";

import { Category } from "src/api/openapi/schemas";
import { Edit } from "src/components/site/Action/Action";

import { CategoryEditModal } from "./CategoryEditModal";

export function CategoryEdit(props: Category) {
  const { onOpen, isOpen, onClose } = useDisclosure();
  return (
    <>
      <Edit onClick={onOpen} />

      <CategoryEditModal onClose={onClose} isOpen={isOpen} category={props} />
    </>
  );
}
