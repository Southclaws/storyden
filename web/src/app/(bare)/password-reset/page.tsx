import { PasswordResetScreen } from "@/screens/auth/PasswordResetScreen/PasswordResetScreen";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

export default function Page() {
  return <PasswordResetScreen />;
}
