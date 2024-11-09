import { EmptyState } from "../site/EmptyState";
import { EmptyThreadsIcon } from "../ui/icons/Empty";

export function FeedEmptyState() {
  return (
    <EmptyState w="full" icon={<EmptyThreadsIcon />}>
      <p>*tumbleweed*&nbsp;there&nbsp;are&nbsp;no&nbsp;posts...</p>
    </EmptyState>
  );
}
