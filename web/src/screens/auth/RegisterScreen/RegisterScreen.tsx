import { AuthMode } from "@/api/openapi-schema";
import { authProviderList } from "@/api/openapi-server/auth";
import { VStack } from "@/styled-system/jsx";

import { RegisterEmailForm } from "./RegisterEmail/RegisterEmailForm";
import { RegisterHandleForm } from "./RegisterHandle/RegisterHandleForm";
import { RegisterPhoneForm } from "./RegisterPhone/RegisterPhoneForm";

export async function RegisterScreen() {
  const { data } = await authProviderList();

  switch (data.mode) {
    case AuthMode.handle:
      return <RegisterHandleForm webauthn={false} />;

    case AuthMode.email:
      return <RegisterEmailForm />;

    case AuthMode.phone:
      return <RegisterPhoneForm />;

    default:
      console.error("no authentication modes available");

      return (
        <VStack>
          <p>This instance is private.</p>
        </VStack>
      );
  }
}
