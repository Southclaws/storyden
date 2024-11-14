import { AuthMode } from "@/api/openapi-schema";
import { authProviderList } from "@/api/openapi-server/auth";
import { VStack } from "@/styled-system/jsx";

import { LoginEmailForm } from "./LoginEmail/LoginEmailForm";
import { LoginHandleForm } from "./LoginHandle/LoginHandleForm";
import { LoginPhoneForm } from "./LoginPhone/LoginPhoneForm";

export async function LoginScreen() {
  const { data } = await authProviderList();

  switch (data.mode) {
    case AuthMode.handle:
      return <LoginHandleForm />;

    case AuthMode.email:
      return <LoginEmailForm />;

    case AuthMode.phone:
      return <LoginPhoneForm />;

    default:
      console.error("no authentication modes available");

      return (
        <VStack>
          <p>This instance is closed.</p>
        </VStack>
      );
  }
}
