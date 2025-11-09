import { Suspense } from "react";

import { Unready } from "@/components/site/Unready";

import { EmailVerifyPage } from "./EmailVerifyPage";

type Props = {
  searchParams: Promise<{
    returnURL?: string;
  }>;
};

export default function Page(props: Props) {
  return (
    <Suspense fallback={<Unready />}>
      <EmailVerifyPage {...props} />
    </Suspense>
  );
}
