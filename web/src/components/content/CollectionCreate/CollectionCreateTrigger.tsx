import { PropsWithChildren } from "react";

import { useDisclosure } from "src/utils/useDisclosure";

import { Button } from "@/components/ui/button";
import { CreateFolderIcon } from "@/components/ui/icons/CreateFolder";
import { ButtonVariantProps } from "@/styled-system/recipes";

import { CollectionCreateModal } from "./CollectionCreateModal";
import { Props } from "./useCollectionCreate";

export function CollectionCreateTrigger(
  props: PropsWithChildren<
    ButtonVariantProps &
      Props & {
        label?: string;
      }
  >,
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
          <CreateFolderIcon /> {props.label ?? "Collection"}
        </Button>
      )}
      <CollectionCreateModal isOpen={isOpen} onClose={onClose} {...props} />
    </>
  );
}
