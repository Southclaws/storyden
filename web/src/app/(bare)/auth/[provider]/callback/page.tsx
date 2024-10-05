"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useRef, useState } from "react";

import { oAuthProviderCallback } from "src/api/openapi-client/auth";
import { UnreadyBanner } from "src/components/site/Unready";

import { LinkButton } from "@/components/ui/link-button";
import { HStack, VStack } from "@/styled-system/jsx";
import { deriveError } from "@/utils/error";

export type Props = {
  params: {
    provider: string;
  };
};

export default function Page(props: Props) {
  const initialized = useRef(false);
  const [error, setError] = useState<string | undefined>(undefined);
  const searchParams = useSearchParams();
  const router = useRouter();

  const params = Object.fromEntries(searchParams.entries());

  useEffect(() => {
    if (!initialized.current) {
      initialized.current = true;
    }

    (async () => {
      try {
        const { id } = await oAuthProviderCallback(
          props.params.provider,
          params as any,
        );

        return router.push(`/?id=${id}`);
      } catch (e) {
        console.log(e);

        const message = deriveError(e);

        setError(message);
      }
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <VStack w="full" height="dvh" justify="center" p="10">
      <UnreadyBanner error={error} />;
      <HStack>
        <LinkButton href="/register">Back to register</LinkButton>
        <LinkButton href="/login">Back to login</LinkButton>
      </HStack>
    </VStack>
  );
}
