import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { MoreIcon } from "@/components/ui/icons/More";
import { useI18n } from "@/i18n/provider";

export function MoreAction(props: ButtonProps) {
  const { t } = useI18n();

  return (
    <IconButton variant="ghost" aria-label={t("More options")} {...props}>
      <MoreIcon />
    </IconButton>
  );
}
