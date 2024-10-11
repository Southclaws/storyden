import { PlusIcon } from "@heroicons/react/24/outline";

import { handle } from "@/api/client";
import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { Item } from "@/components/ui/menu";
import { useLibraryMutation } from "@/lib/library/library";

type Props = ButtonProps & {
  parentSlug?: string;
  hideLabel?: boolean;
};

export const CreatePageID = "create-page";
export const CreatePageLabel = "Create";
export const CreatePageIcon = <PlusIcon />;

export function CreatePageAction({ parentSlug, hideLabel, ...props }: Props) {
  const { createNode, revalidate } = useLibraryMutation();

  async function handleCreate() {
    await handle(
      async () => {
        await createNode({ parentSlug });
      },
      {
        promiseToast: {
          loading: "Creating page...",
          success: "Page created!",
        },
        cleanup: async () => revalidate(),
      },
    );
  }

  return (
    <IconButton
      type="button"
      size="xs"
      variant="subtle"
      px={hideLabel ? "0" : "1"}
      onClick={handleCreate}
      {...props}
    >
      {CreatePageIcon}
      {!hideLabel && (
        <>
          <span>{CreatePageLabel}</span>
        </>
      )}
    </IconButton>
  );
}

export function CreatePageMenuItem({ hideLabel }: Props) {
  return (
    <Item value={CreatePageID}>
      {CreatePageIcon}
      {!hideLabel && (
        <>
          &nbsp;<span>{CreatePageLabel}</span>
        </>
      )}
    </Item>
  );
}
