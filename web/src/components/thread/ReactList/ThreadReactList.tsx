import { ReactList } from "@/components/ui/react-list";

import { Props, useReactionList } from "./useReactList";

export function ThreadReactList(props: Props) {
  const { data, handlers } = useReactionList(props);

  return (
    <ReactList
      isLoggedIn={data.isLoggedIn}
      pendingEmojis={data.pendingEmojis}
      reactions={data.reactions}
      onReactExisting={handlers.handleReactExisting}
      onReactPicker={handlers.handleReactPicker}
    />
  );
}
