"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { oAuthProviderCallback } from "src/api/openapi-client/auth";
import { Unready } from "src/components/site/Unready";

import { handle } from "@/api/client";
import { OAuthCallback } from "@/api/openapi-schema";
import { deriveError } from "@/utils/error";

export type Props = {
  params: Promise<{
    provider: string;
  }>;
  searchParams: Promise<OAuthCallback>;
};

export default function Page(props: Props) {
  const router = useRouter();
  const [error, setError] = useState<unknown | null>(null);

  useEffect(() => {
    handle(
      async () => {
        const params = await props.params;
        const searchParams = await props.searchParams;

        if (error != null) {
          return;
        }

        if (params.provider === undefined || searchParams === undefined) {
          return;
        }

        const { id } = await oAuthProviderCallback(
          params.provider,
          searchParams,
        );

        router.push(`/?id=${id}`);
      },
      {
        errorToast: false,
        onError: async (e) => {
          console.error("OAuth callback error:", e);
          setError(deriveError(e));
        },
      },
    );
  }, [router, error, props]);

  return <Unready error={error} />;
}
