import { FolderPlusIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { useDisclosure } from "src/utils/useDisclosure";

import { Button } from "@/components/ui/button";
import { ButtonVariantProps } from "@/styled-system/recipes";

import { CollectionCreateModal } from "./CollectionCreateModal";

type Props = {
  label?: string;
};

export function CollectionCreateTrigger(
  props: PropsWithChildren<ButtonVariantProps & Props>,
) {
  const { onOpen, isOpen, onClose } = useDisclosure();
  return (
    <>
      {props.children ?? (
        <Button
          flexShrink="0"
          minW="0"
          variant="subtle"
          justifyContent="start"
          size="sm"
          {...props}
          onClick={onOpen}
        >
          <FolderPlusIcon width="1.4rem" /> {props.label ?? "Collection"}
        </Button>
      )}
      <CollectionCreateModal isOpen={isOpen} onClose={onClose} {...props} />
    </>
  );
}
