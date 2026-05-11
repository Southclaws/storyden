import { AuthMode } from "@/api/openapi-schema";
import { authProviderList } from "@/api/openapi-server/auth";
import { Unready } from "@/components/site/Unready";
import { tServer } from "@/i18n/server";
import { VStack } from "@/styled-system/jsx";

import { PasswordResetEmailScreen } from "./PasswordResetEmailScreen";

export async function PasswordResetScreen() {
  const { data } = await authProviderList({
    cache: "no-store",
  });
  const unavailable = await tServer(
    "Password reset is currently disabled. Please contact the site administrator.",
  );

  switch (data.mode) {
    case AuthMode.handle:
      return <Warning message={unavailable} />;

    case AuthMode.email:
      return <PasswordResetEmailScreen />;

    case AuthMode.phone:
      return <Warning message={unavailable} />;

    default:
      return (
        <Unready error="No authentication modes are currently enabled." />
      );
  }
}

function Warning({ message }: { message: string }) {
  return (
    <VStack textWrap="balance" textAlign="center">
      <p>{message}</p>
    </VStack>
  );
}
