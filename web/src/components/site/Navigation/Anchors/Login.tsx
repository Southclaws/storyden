import { LoginIcon } from "src/components/graphics/LoginIcon";
import { Link } from "src/theme/components/Link";

export function LoginAction() {
  return (
    <Link href="/register" kind="ghost" size="sm" p="0">
      <LoginIcon width="1.5em" />
    </Link>
  );
}
