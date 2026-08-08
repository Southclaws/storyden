import type {
  React as Reaction,
  ReactList as Reactions,
} from "@/api/openapi-schema";
import type { ReactListItem } from "@/components/ui/react-list";

export type ReactionGroup = ReactListItem & {
  reactions: Reaction[];
};

export function groupReactions(
  accountID: string | undefined,
  reactions: Reactions,
): ReactionGroup[] {
  const groups = new Map<string, Reaction[]>();

  for (const reaction of reactions) {
    const group = groups.get(reaction.emoji);
    if (group) {
      group.push(reaction);
    } else {
      groups.set(reaction.emoji, [reaction]);
    }
  }

  return [...groups.entries()].map(([emoji, groupedReactions]) => ({
    emoji,
    count: groupedReactions.length,
    hasReacted:
      accountID !== undefined &&
      groupedReactions.some((reaction) => reaction.author?.id === accountID),
    reactedBy: groupedReactions
      .map((reaction) => reaction.author?.handle)
      .filter((handle): handle is string => Boolean(handle)),
    reactions: groupedReactions,
  }));
}
