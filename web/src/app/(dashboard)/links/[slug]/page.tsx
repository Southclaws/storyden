import { Suspense } from "react";

import { Unready } from "@/components/site/Unready";

import { LinkPage } from "./LinkPage";

type Props = {
  params: Promise<{
    slug: string;
  }>;
};

export default function Page(props: Props) {
  return (
    <Suspense fallback={<Unready />}>
      <LinkPage {...props} />
    </Suspense>
  );
}
