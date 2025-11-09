import { Suspense } from "react";

import { Unready } from "@/components/site/Unready";

import { PasswordResetVerifyPage } from "./PasswordResetVerifyPage";

type Props = {
  searchParams: Promise<{
    token: string;
  }>;
};

export default function Page(props: Props) {
  return (
    <Suspense fallback={<Unready />}>
      <PasswordResetVerifyPage {...props} />
    </Suspense>
  );
}
