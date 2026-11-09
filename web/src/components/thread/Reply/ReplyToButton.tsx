import { Account, Reply, Thread } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { IconButton } from "@/components/ui/icon-button";
import { ReplyIcon } from "@/components/ui/icons/Reply";
import { canReplyToThread } from "@/lib/thread/permissions";

import { useReplyContext } from "../ReplyContext";

type Props = {
  initialSession?: Account;
  thread: Thread;
  reply: Reply;
};

export function ReplyToButton(props: Props) {
  const session = useSession(props.initialSession);
  const { setReplyTo } = useReplyContext();

  function handleClick() {
    setReplyTo(props.thread, props.reply);
  }

  if (props.thread.locked_at && !canReplyToThread(props.thread, session)) {
    return null;
  }

  return (
    <IconButton
      type="button"
      variant="ghost"
      aria-label="Reply to this"
      onClick={handleClick}
    >
      <ReplyIcon />
    </IconButton>
  );
}
