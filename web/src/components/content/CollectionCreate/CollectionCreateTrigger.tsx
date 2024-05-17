import { FolderPlusIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { Button } from "src/theme/components/Button";
import { useDisclosure } from "src/utils/useDisclosure";

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
          variant="subtle"
          justifyContent="start"
          size="sm"
          onClick={onOpen}
          {...props}
        >
          <FolderPlusIcon /> Create collection
        </Button>
      )}
      <CollectionCreateModal isOpen={isOpen} onClose={onClose} {...props} />
    </>
  );
}
