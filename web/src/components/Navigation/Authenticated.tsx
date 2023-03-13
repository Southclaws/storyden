import { IconButton } from "@chakra-ui/react";
import { BellIcon, PlusIcon } from "@heroicons/react/24/outline";
import { Anchor } from "../site/Anchor";

type Props = { onExpand: () => void; category?: string };

export function Authenticated({ onExpand, category }: Props) {
  return (
    <>
      <Anchor href="/notifications">
        <IconButton
          aria-label=""
          borderRadius="50%"
          icon={<BellIcon width="1em" />}
        />
      </Anchor>

      <Anchor href="/categories" onClick={onExpand}>
        {category ?? "Categories"}
      </Anchor>

      <Anchor href="/new">
        <IconButton
          aria-label=""
          borderRadius="50%"
          icon={<PlusIcon width="1em" />}
        />
      </Anchor>
    </>
  );
}
