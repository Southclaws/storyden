import { Button, useDisclosure } from "@chakra-ui/react";
import { PlusIcon } from "@heroicons/react/24/solid";

import { CategoryCreateModal } from "./CategoryCreateModal";

export function CategoryCreate() {
  const { onOpen, isOpen, onClose } = useDisclosure();
  return (
    <>
      <Button
        w="full"
        size="xs"
        variant="outline"
        leftIcon={<PlusIcon width="1.125em" />}
        onClick={onOpen}
      >
        New category
      </Button>

      <CategoryCreateModal onClose={onClose} isOpen={isOpen} />
    </>
  );
}
