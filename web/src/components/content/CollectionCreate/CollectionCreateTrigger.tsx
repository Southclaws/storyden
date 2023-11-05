import { PlusIcon } from "@heroicons/react/24/solid";
import { PropsWithChildren } from "react";

import { useDisclosure } from "src/theme/components";
import { Button } from "src/theme/components/Button";

import { ButtonVariantProps } from "@/styled-system/recipes";

import { CollectionCreateModal } from "./CollectionCreateModal";

export function CollectionCreateTrigger(
  props: PropsWithChildren<ButtonVariantProps>,
) {
  const { onOpen, isOpen, onClose } = useDisclosure();
  return (
    <>
      {props.children ?? (
        <Button
          w="full"
          justifyContent="start"
          size="sm"
          kind="primary"
          onClick={onOpen}
          {...props}
        >
          <PlusIcon /> Create collection
        </Button>
      )}
      <CollectionCreateModal isOpen={isOpen} onClose={onClose} {...props} />
    </>
  );
}
