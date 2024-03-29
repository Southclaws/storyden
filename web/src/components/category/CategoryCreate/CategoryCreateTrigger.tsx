import { PlusIcon } from "@heroicons/react/24/solid";
import { PropsWithChildren } from "react";

import { Button } from "src/theme/components/Button";
import { useDisclosure } from "src/utils/useDisclosure";

import { ButtonVariantProps } from "@/styled-system/recipes";

import { CategoryCreateModal } from "./CategoryCreateModal";

export function CategoryCreateTrigger(
  props: PropsWithChildren<ButtonVariantProps>,
) {
  const { onOpen, isOpen, onClose } = useDisclosure();
  return (
    <>
      {props.children ?? (
        <Button
          w="full"
          justifyContent="start"
          size="xs"
          kind="neutral"
          onClick={onOpen}
          {...props}
        >
          <PlusIcon /> Create category
        </Button>
      )}
      <CategoryCreateModal isOpen={isOpen} onClose={onClose} {...props} />
    </>
  );
}
