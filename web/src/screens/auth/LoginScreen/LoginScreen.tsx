import Image from "next/image";

import { LinkButton } from "@/components/ui/link-button";
import { css } from "@/styled-system/css";
import { VStack, styled } from "@/styled-system/jsx";
import { getIconURL } from "@/utils/icon";
import { getInfo } from "@/utils/info";

import { AuthSelection } from "../components/AuthSelection/AuthSelection";
import { getProviders } from "../providers";

import { LoginForm } from "./LoginForm";

export async function LoginScreen() {
  const info = await getInfo();
  const { password, webauthn } = await getProviders();

  // TODO: Phone login form.
  const isOnlyOAuth = password === false && webauthn === false; // && phone === false
  // ...
  // Then show this copy if there's neither a password or phone form.
  // {isOnlyOAuth && <styled.p textAlign="center">choose a login provider</styled.p>}

  return (
    <VStack minH="dvh" p="12">
      <VStack>
        <Image
          className={css({ width: "28" })}
          src={getIconURL("512x512")}
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
              Log in to your account with your username and password.
            </styled.p>

            <LoginForm webauthn={webauthn} />

            <styled.p>or choose a login provider</styled.p>
          </>
        )}

        {isOnlyOAuth && <styled.p>choose a login provider</styled.p>}

        <AuthSelection />
      </VStack>

      <p>
        <LinkButton size="xs" href="/register">
          Register
        </LinkButton>
      </p>
    </VStack>
  );
}
