import { Reply, Thread } from "@/api/openapi-schema";
import { IconButton } from "@/components/ui/icon-button";
import { ReplyIcon } from "@/components/ui/icons/Reply";
import { useI18n } from "@/i18n/provider";

import { useReplyContext } from "../ReplyContext";

type Props = {
  thread: Thread;
  reply: Reply;
};

export function ReplyToButton(props: Props) {
  const { setReplyTo } = useReplyContext();
  const { t } = useI18n();

  function handleClick() {
    setReplyTo(props.thread, props.reply);
  }

  return (
    <IconButton
      type="button"
      size="xs"
      variant="ghost"
      aria-label={t("Reply to this")}
      onClick={handleClick}
    >
      <ReplyIcon />
    </IconButton>
  );
}
