import { HomeIcon } from "@heroicons/react/24/outline";

import { LinkButton } from "@/components/ui/link-button";

export function HomeAction() {
  return (
    <LinkButton href="/" variant="ghost" size="sm" p="0">
      <HomeIcon />
    </LinkButton>
  );
}
