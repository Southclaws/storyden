import { Box, IconButton } from "@chakra-ui/react";
import { UserIcon } from "@heroicons/react/24/outline";
import { Anchor } from "../site/Anchor";

type Props = { onExpand: () => void };

export function Unauthenticated({ onExpand }: Props) {
  return (
    <>
      <Box width="1em" />

      <Anchor href="/categories" onClick={onExpand}>
        General
      </Anchor>

      <Anchor href="/auth">
        <IconButton
          aria-label=""
          borderRadius="50%"
          icon={<UserIcon width="1em" />}
        />
      </Anchor>
    </>
  );
}
