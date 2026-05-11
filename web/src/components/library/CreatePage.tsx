import { handle } from "@/api/client";
import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { CreateIcon } from "@/components/ui/icons/Create";
import { Item } from "@/components/ui/menu";
import { useI18n } from "@/i18n/provider";
import { useLibraryMutation } from "@/lib/library/library";

type Props = ButtonProps & {
  parentSlug?: string;
  hideLabel?: boolean;
  disableRedirect?: boolean;
  onComplete?: () => void;
};

export const CreatePageID = "create-page";
export const CreatePageLabel = "Create";
export const CreatePageIcon = <CreateIcon />;

export function CreatePageAction({
  parentSlug,
  hideLabel,
  disableRedirect,
  onComplete,
  ...props
}: Props) {
  const { t } = useI18n();
  const { createNode, revalidate } = useLibraryMutation();

  async function handleCreate() {
    await handle(
      async () => {
        await createNode({ parentSlug, disableRedirect });
      },
      {
        promiseToast: {
          loading: t("Creating page..."),
          success: t("Page created!"),
        },
        cleanup: async () => {
          await revalidate();
          onComplete?.();
        },
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
          <span>{t(CreatePageLabel)}</span>
        </>
      )}
    </IconButton>
  );
}

export function CreatePageMenuItem({ hideLabel }: Props) {
  const { t } = useI18n();

  return (
    <Item value={CreatePageID}>
      {CreatePageIcon}
      {!hideLabel && (
        <>
          &nbsp;<span>{t(CreatePageLabel)}</span>
        </>
      )}
    </Item>
  );
}
