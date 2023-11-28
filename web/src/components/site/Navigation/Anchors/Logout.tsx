import { LogoutIcon } from "src/components/graphics/LogoutIcon";

import { link } from "@/styled-system/recipes";

export function LogoutAction() {
  return (
    <a
      className={link({ kind: "ghost", size: "sm" })}
      href="/logout"
      title="Log out of your session"
    >
      <LogoutIcon width="1.5em" />
    </a>
  );
}
