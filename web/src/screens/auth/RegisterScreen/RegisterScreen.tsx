import Image from "next/image";

import { Link } from "src/theme/components/Link";
import { getInfo } from "src/utils/info";

import { AuthSelection } from "../components/AuthSelection/AuthSelection";
import { getProviders } from "../providers";

import { css } from "@/styled-system/css";
import { VStack, styled } from "@/styled-system/jsx";

import { RegisterForm } from "./RegisterForm";

export async function RegisterScreen() {
  const info = await getInfo();
  const { password, webauthn } = await getProviders();

  // TODO: Phone login form.
  const isOnlyOAuth = password === false && webauthn === false; // && phone === false
  // ...
  // Then show this copy if there's neither a password or phone form.
  // {isOnlyOAuth && <styled.p textAlign="center">choose a login provider</styled.p>}

  return (
    <VStack
      minH="dvh"
      justify="center"
      p="12"
      background={
        "linear-gradient(180deg, var(--accent-colour-flat-fill-100) 0%, var(--accent-colour-flat-fill-200) 100%)" as any
      }
    >
      <VStack>
        <Image
          className={css({ width: "28" })}
          src="/api/v1/info/icon/512x512"
          width="512"
          height="512"
          alt={`The ${info.title} logo`}
        />

        <styled.h1 fontWeight="bold" fontSize="lg">
          {info.title}
        </styled.h1>
      </VStack>

      <VStack w="full" maxW="xs" textAlign="center">
        {password && (
          <>
            <styled.p>
              Join&nbsp;the&nbsp;community&nbsp;with&nbsp;a&nbsp;
              <wbr />
              username&nbsp;and&nbsp;password
            </styled.p>

            <RegisterForm webauthn={webauthn} />

            <styled.p>or choose a login provider</styled.p>
          </>
        )}

        {isOnlyOAuth && <styled.p>choose a login provider</styled.p>}

        <AuthSelection />
      </VStack>

      <p>
        <Link size="xs" href="/login">
          Sign in
        </Link>
      </p>
    </VStack>
  );
}
