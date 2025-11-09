import { Suspense } from "react";

import { Unready } from "@/components/site/Unready";
import { PasswordResetScreen } from "@/screens/auth/PasswordResetScreen/PasswordResetScreen";

export default function Page() {
  return (
    <Suspense fallback={<Unready />}>
      <PasswordResetScreen />
    </Suspense>
  );
}
