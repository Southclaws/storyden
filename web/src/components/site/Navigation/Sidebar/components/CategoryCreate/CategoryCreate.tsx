import { PlusIcon } from "@heroicons/react/24/solid";

import { Add } from "src/components/site/Action/Action";
import { Button, useDisclosure } from "src/theme/components";

import { CategoryCreateModal } from "./CategoryCreateModal";

type Props = {
  action: "text" | "icon";
};

export function CategoryCreate(props: Props) {
  const { onOpen, isOpen, onClose } = useDisclosure();
  return (
    <>
      {props.action === "text" ? (
        <Button
          w="full"
          size="xs"
          variant="outline"
          leftIcon={<PlusIcon width="1.125em" />}
          onClick={onOpen}
        >
          New category
        </Button>
      ) : (
        <Add onClick={onOpen} />
      )}

      <CategoryCreateModal onClose={onClose} isOpen={isOpen} />
    </>
  );
}
