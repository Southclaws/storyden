import { ProfileReference } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { Anchor } from "src/components/site/Anchor";
import { Avatar } from "src/components/site/Avatar/Avatar";

import { css } from "@/styled-system/css";
import { Box } from "@/styled-system/jsx";

import { Handle } from "./Handle";

type Props = {
  profileReference: ProfileReference;
  showHandle?: boolean;
  size?: "sm" | "lg";
};

export function ProfilePill({
  profileReference,
  showHandle = true,
  size = "sm",
}: Props) {
  const account = useSession();
  const self = account?.id === profileReference.id;
  const title = self ? `Your profile` : `${profileReference.handle}'s profile`;
  const large = size === "lg";

  return (
    <Anchor
      title={title}
      href={`/p/${profileReference.handle}`}
      className={css({
        flexShrink: 1,
        pr: showHandle ? "1" : "0",
        borderRadius: "full",
        minW: "0",
        maxW: "full",
        display: "flex",
        gap: "1",
      })}
    >
      <Avatar
        flexShrink={0}
        handle={profileReference.handle}
        width={large ? "8" : "6"}
      />
      {showHandle && (
        <Box minW="0" flexShrink={1}>
          <Handle profileReference={profileReference} size={size} />
        </Box>
      )}
    </Anchor>
  );
}
