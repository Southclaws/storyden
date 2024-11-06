import { useDisclosure } from "src/utils/useDisclosure";

import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { CreateIcon } from "@/components/ui/icons/Create";
import { Item } from "@/components/ui/menu";

import { CategoryCreateModal } from "./CategoryCreateModal";

type Props = ButtonProps & {
  parentSlug?: string;
  hideLabel?: boolean;
};

export const CreateCategoryID = "create-category";
export const CreateCategoryLabel = "Create";
export const CreateCategoryIcon = <CreateIcon />;

export function CategoryCreateTrigger({
  parentSlug,
  hideLabel,
  ...props
}: Props) {
  const { onOpen, isOpen, onClose } = useDisclosure();

  return (
    <>
      <IconButton
        type="button"
        size="xs"
        variant="ghost"
        px={hideLabel ? "0" : "1"}
        onClick={onOpen}
        {...props}
      >
        {CreateCategoryIcon}
        {!hideLabel && (
          <>
            <span>{CreateCategoryLabel}</span>
          </>
        )}
      </IconButton>

      <CategoryCreateModal isOpen={isOpen} onClose={onClose} {...props} />
    </>
  );
}

export function CreateCategoryMenuItem({ hideLabel }: Props) {
  return (
    <Item value={CreateCategoryID}>
      {CreateCategoryIcon}
      {!hideLabel && (
        <>
          &nbsp;<span>{CreateCategoryLabel}</span>
        </>
      )}
    </Item>
  );
}
