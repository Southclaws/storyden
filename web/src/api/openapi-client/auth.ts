/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: v1.25.3-canary
 */
import useSwr from "swr";
import type { Arguments, Key, SWRConfiguration } from "swr";
import useSWRMutation from "swr/mutation";
import type { SWRMutationConfiguration } from "swr/mutation";

import { fetcher } from "../client";
import type {
  AccessKeyCreateBody,
  AccessKeyCreateOKResponse,
  AccessKeyListOKResponse,
  AuthEmailBody,
  AuthEmailPasswordBody,
  AuthEmailPasswordResetBody,
  AuthEmailPasswordSignupParams,
  AuthEmailSignupParams,
  AuthEmailVerifyBody,
  AuthPasswordBody,
  AuthPasswordCreateBody,
  AuthPasswordResetBody,
  AuthPasswordSignupParams,
  AuthPasswordUpdateBody,
  AuthProviderListOKResponse,
  AuthSuccessOKResponse,
  BadRequestResponse,
  ForbiddenResponse,
  InternalServerErrorResponse,
  NoContentResponse,
  NotFoundResponse,
  OAuthProviderCallbackBody,
  PhoneRequestCodeBody,
  PhoneRequestCodeParams,
  PhoneSubmitCodeBody,
  UnauthorisedResponse,
  WebAuthnGetAssertionOKResponse,
  WebAuthnMakeAssertionBody,
  WebAuthnMakeCredentialBody,
  WebAuthnMakeCredentialParams,
  WebAuthnRequestCredentialOKResponse,
} from "../openapi-schema";

/**
 * Retrieve a list of authentication providers. Storyden supports a few
ways to authenticate, from simple passwords to OAuth and WebAuthn. This
endpoint tells a client which auth capabilities are enabled.

 */
export const authProviderList = () => {
  return fetcher<AuthProviderListOKResponse>({ url: `/auth`, method: "GET" });
};

export const getAuthProviderListKey = () => [`/auth`] as const;

export type AuthProviderListQueryResult = NonNullable<
  Awaited<ReturnType<typeof authProviderList>>
>;
export type AuthProviderListQueryError =
  | BadRequestResponse
  | InternalServerErrorResponse;

export const useAuthProviderList = <
  TError = BadRequestResponse | InternalServerErrorResponse,
>(options?: {
  swr?: SWRConfiguration<
    Awaited<ReturnType<typeof authProviderList>>,
    TError
  > & { swrKey?: Key; enabled?: boolean };
}) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false;
  const swrKey =
    swrOptions?.swrKey ?? (() => (isEnabled ? getAuthProviderListKey() : null));
  const swrFn = () => authProviderList();

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(
    swrKey,
    swrFn,
    swrOptions,
  );

  return {
    swrKey,
    ...query,
  };
};
/**
 * Register a new account with a username and password.
 */
export const authPasswordSignup = (
  authPasswordBody: AuthPasswordBody,
  params?: AuthPasswordSignupParams,
) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/password/signup`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: authPasswordBody,
    params,
  });
};

export const getAuthPasswordSignupMutationFetcher = (
  params?: AuthPasswordSignupParams,
) => {
  return (
    _: Key,
    { arg }: { arg: AuthPasswordBody },
  ): Promise<AuthSuccessOKResponse> => {
    return authPasswordSignup(arg, params);
  };
};
export const getAuthPasswordSignupMutationKey = (
  params?: AuthPasswordSignupParams,
) => [`/auth/password/signup`, ...(params ? [params] : [])] as const;

export type AuthPasswordSignupMutationResult = NonNullable<
  Awaited<ReturnType<typeof authPasswordSignup>>
>;
export type AuthPasswordSignupMutationError =
  | BadRequestResponse
  | InternalServerErrorResponse;

export const useAuthPasswordSignup = <
  TError = BadRequestResponse | InternalServerErrorResponse,
>(
  params?: AuthPasswordSignupParams,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof authPasswordSignup>>,
      TError,
      Key,
      AuthPasswordBody,
      Awaited<ReturnType<typeof authPasswordSignup>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getAuthPasswordSignupMutationKey(params);
  const swrFn = getAuthPasswordSignupMutationFetcher(params);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Sign in to an existing account with a username and password.
 */
export const authPasswordSignin = (authPasswordBody: AuthPasswordBody) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/password/signin`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: authPasswordBody,
  });
};

export const getAuthPasswordSigninMutationFetcher = () => {
  return (
    _: Key,
    { arg }: { arg: AuthPasswordBody },
  ): Promise<AuthSuccessOKResponse> => {
    return authPasswordSignin(arg);
  };
};
export const getAuthPasswordSigninMutationKey = () =>
  [`/auth/password/signin`] as const;

export type AuthPasswordSigninMutationResult = NonNullable<
  Awaited<ReturnType<typeof authPasswordSignin>>
>;
export type AuthPasswordSigninMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useAuthPasswordSignin = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(options?: {
  swr?: SWRMutationConfiguration<
    Awaited<ReturnType<typeof authPasswordSignin>>,
    TError,
    Key,
    AuthPasswordBody,
    Awaited<ReturnType<typeof authPasswordSignin>>
  > & { swrKey?: string };
}) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getAuthPasswordSigninMutationKey();
  const swrFn = getAuthPasswordSigninMutationFetcher();

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Given the requesting account does not have a password authentication,
add a password authentication method to it with the given password.

 */
export const authPasswordCreate = (
  authPasswordCreateBody: AuthPasswordCreateBody,
) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/password`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: authPasswordCreateBody,
  });
};

export const getAuthPasswordCreateMutationFetcher = () => {
  return (
    _: Key,
    { arg }: { arg: AuthPasswordCreateBody },
  ): Promise<AuthSuccessOKResponse> => {
    return authPasswordCreate(arg);
  };
};
export const getAuthPasswordCreateMutationKey = () =>
  [`/auth/password`] as const;

export type AuthPasswordCreateMutationResult = NonNullable<
  Awaited<ReturnType<typeof authPasswordCreate>>
>;
export type AuthPasswordCreateMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useAuthPasswordCreate = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(options?: {
  swr?: SWRMutationConfiguration<
    Awaited<ReturnType<typeof authPasswordCreate>>,
    TError,
    Key,
    AuthPasswordCreateBody,
    Awaited<ReturnType<typeof authPasswordCreate>>
  > & { swrKey?: string };
}) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getAuthPasswordCreateMutationKey();
  const swrFn = getAuthPasswordCreateMutationFetcher();

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Given the requesting account has a password authentication, update the
password on file.

 */
export const authPasswordUpdate = (
  authPasswordUpdateBody: AuthPasswordUpdateBody,
) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/password`,
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    data: authPasswordUpdateBody,
  });
};

export const getAuthPasswordUpdateMutationFetcher = () => {
  return (
    _: Key,
    { arg }: { arg: AuthPasswordUpdateBody },
  ): Promise<AuthSuccessOKResponse> => {
    return authPasswordUpdate(arg);
  };
};
export const getAuthPasswordUpdateMutationKey = () =>
  [`/auth/password`] as const;

export type AuthPasswordUpdateMutationResult = NonNullable<
  Awaited<ReturnType<typeof authPasswordUpdate>>
>;
export type AuthPasswordUpdateMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useAuthPasswordUpdate = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(options?: {
  swr?: SWRMutationConfiguration<
    Awaited<ReturnType<typeof authPasswordUpdate>>,
    TError,
    Key,
    AuthPasswordUpdateBody,
    Awaited<ReturnType<typeof authPasswordUpdate>>
  > & { swrKey?: string };
}) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getAuthPasswordUpdateMutationKey();
  const swrFn = getAuthPasswordUpdateMutationFetcher();

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Complete a password-reset flow using a token that was provided to the
member via a reset request operation such as `AuthEmailPasswordReset`.

 */
export const authPasswordReset = (
  authPasswordResetBody: AuthPasswordResetBody,
) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/password/reset`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: authPasswordResetBody,
  });
};

export const getAuthPasswordResetMutationFetcher = () => {
  return (
    _: Key,
    { arg }: { arg: AuthPasswordResetBody },
  ): Promise<AuthSuccessOKResponse> => {
    return authPasswordReset(arg);
  };
};
export const getAuthPasswordResetMutationKey = () =>
  [`/auth/password/reset`] as const;

export type AuthPasswordResetMutationResult = NonNullable<
  Awaited<ReturnType<typeof authPasswordReset>>
>;
export type AuthPasswordResetMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useAuthPasswordReset = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(options?: {
  swr?: SWRMutationConfiguration<
    Awaited<ReturnType<typeof authPasswordReset>>,
    TError,
    Key,
    AuthPasswordResetBody,
    Awaited<ReturnType<typeof authPasswordReset>>
  > & { swrKey?: string };
}) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getAuthPasswordResetMutationKey();
  const swrFn = getAuthPasswordResetMutationFetcher();

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Register a new account with a email and password.
 */
export const authEmailPasswordSignup = (
  authEmailPasswordBody: AuthEmailPasswordBody,
  params?: AuthEmailPasswordSignupParams,
) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/email-password/signup`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: authEmailPasswordBody,
    params,
  });
};

export const getAuthEmailPasswordSignupMutationFetcher = (
  params?: AuthEmailPasswordSignupParams,
) => {
  return (
    _: Key,
    { arg }: { arg: AuthEmailPasswordBody },
  ): Promise<AuthSuccessOKResponse> => {
    return authEmailPasswordSignup(arg, params);
  };
};
export const getAuthEmailPasswordSignupMutationKey = (
  params?: AuthEmailPasswordSignupParams,
) => [`/auth/email-password/signup`, ...(params ? [params] : [])] as const;

export type AuthEmailPasswordSignupMutationResult = NonNullable<
  Awaited<ReturnType<typeof authEmailPasswordSignup>>
>;
export type AuthEmailPasswordSignupMutationError =
  | BadRequestResponse
  | InternalServerErrorResponse;

export const useAuthEmailPasswordSignup = <
  TError = BadRequestResponse | InternalServerErrorResponse,
>(
  params?: AuthEmailPasswordSignupParams,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof authEmailPasswordSignup>>,
      TError,
      Key,
      AuthEmailPasswordBody,
      Awaited<ReturnType<typeof authEmailPasswordSignup>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getAuthEmailPasswordSignupMutationKey(params);
  const swrFn = getAuthEmailPasswordSignupMutationFetcher(params);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Sign in to an existing account with a email and password.
 */
export const authEmailPasswordSignin = (
  authEmailPasswordBody: AuthEmailPasswordBody,
) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/email-password/signin`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: authEmailPasswordBody,
  });
};

export const getAuthEmailPasswordSigninMutationFetcher = () => {
  return (
    _: Key,
    { arg }: { arg: AuthEmailPasswordBody },
  ): Promise<AuthSuccessOKResponse> => {
    return authEmailPasswordSignin(arg);
  };
};
export const getAuthEmailPasswordSigninMutationKey = () =>
  [`/auth/email-password/signin`] as const;

export type AuthEmailPasswordSigninMutationResult = NonNullable<
  Awaited<ReturnType<typeof authEmailPasswordSignin>>
>;
export type AuthEmailPasswordSigninMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useAuthEmailPasswordSignin = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(options?: {
  swr?: SWRMutationConfiguration<
    Awaited<ReturnType<typeof authEmailPasswordSignin>>,
    TError,
    Key,
    AuthEmailPasswordBody,
    Awaited<ReturnType<typeof authEmailPasswordSignin>>
  > & { swrKey?: string };
}) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getAuthEmailPasswordSigninMutationKey();
  const swrFn = getAuthEmailPasswordSigninMutationFetcher();

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Request password reset email to be sent to the specified email address.

 */
export const authPasswordResetRequestEmail = (
  authEmailPasswordResetBody: AuthEmailPasswordResetBody,
) => {
  return fetcher<void>({
    url: `/auth/email-password/reset`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: authEmailPasswordResetBody,
  });
};

export const getAuthPasswordResetRequestEmailMutationFetcher = () => {
  return (
    _: Key,
    { arg }: { arg: AuthEmailPasswordResetBody },
  ): Promise<void> => {
    return authPasswordResetRequestEmail(arg);
  };
};
export const getAuthPasswordResetRequestEmailMutationKey = () =>
  [`/auth/email-password/reset`] as const;

export type AuthPasswordResetRequestEmailMutationResult = NonNullable<
  Awaited<ReturnType<typeof authPasswordResetRequestEmail>>
>;
export type AuthPasswordResetRequestEmailMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useAuthPasswordResetRequestEmail = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(options?: {
  swr?: SWRMutationConfiguration<
    Awaited<ReturnType<typeof authPasswordResetRequestEmail>>,
    TError,
    Key,
    AuthEmailPasswordResetBody,
    Awaited<ReturnType<typeof authPasswordResetRequestEmail>>
  > & { swrKey?: string };
}) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getAuthPasswordResetRequestEmailMutationKey();
  const swrFn = getAuthPasswordResetRequestEmailMutationFetcher();

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Register a new account with an email and optional password. The password
requirement is dependent on how the instance is configured for account
authentication with email addresses (password vs magic link.)

When the email address has not been registered, this endpoint will send
a verification email however it will also return a session cookie to
facilitate pre-verification usage of the platform. If the email address
already exists, no session cookie will be returned in order to prevent
arbitrary account control by a malicious actor. In this case, the email
will be sent again with the same OTP for the case where the user has
cleared their cookies or switched device but hasn't yet verified due to
missing the email or a delivery failure. In this sense, the endpoint can
act as a "resend verification email" operation as well as registration.

In the first case, a 200 response is provided with the session cookie,
in the second case, a 422 response is provided without a session cookie.

Given that this is an unauthenticated endpoint that triggers an email to
be sent to any public address, it MUST be heavily rate limited.

 */
export const authEmailSignup = (
  authEmailBody: AuthEmailBody,
  params?: AuthEmailSignupParams,
) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/email/signup`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: authEmailBody,
    params,
  });
};

export const getAuthEmailSignupMutationFetcher = (
  params?: AuthEmailSignupParams,
) => {
  return (
    _: Key,
    { arg }: { arg: AuthEmailBody },
  ): Promise<AuthSuccessOKResponse> => {
    return authEmailSignup(arg, params);
  };
};
export const getAuthEmailSignupMutationKey = (params?: AuthEmailSignupParams) =>
  [`/auth/email/signup`, ...(params ? [params] : [])] as const;

export type AuthEmailSignupMutationResult = NonNullable<
  Awaited<ReturnType<typeof authEmailSignup>>
>;
export type AuthEmailSignupMutationError =
  | BadRequestResponse
  | void
  | InternalServerErrorResponse;

export const useAuthEmailSignup = <
  TError = BadRequestResponse | void | InternalServerErrorResponse,
>(
  params?: AuthEmailSignupParams,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof authEmailSignup>>,
      TError,
      Key,
      AuthEmailBody,
      Awaited<ReturnType<typeof authEmailSignup>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getAuthEmailSignupMutationKey(params);
  const swrFn = getAuthEmailSignupMutationFetcher(params);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Sign in to an existing account with an email and optional password. The
behaviour of this endpoint depends on how the instance is configured. If
email+password is the preferred method, a cookie is returned on success
but if magic links are preferred, the endpoint will start the code flow.

 */
export const authEmailSignin = (authEmailBody: AuthEmailBody) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/email/signin`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: authEmailBody,
  });
};

export const getAuthEmailSigninMutationFetcher = () => {
  return (
    _: Key,
    { arg }: { arg: AuthEmailBody },
  ): Promise<AuthSuccessOKResponse> => {
    return authEmailSignin(arg);
  };
};
export const getAuthEmailSigninMutationKey = () =>
  [`/auth/email/signin`] as const;

export type AuthEmailSigninMutationResult = NonNullable<
  Awaited<ReturnType<typeof authEmailSignin>>
>;
export type AuthEmailSigninMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useAuthEmailSignin = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(options?: {
  swr?: SWRMutationConfiguration<
    Awaited<ReturnType<typeof authEmailSignin>>,
    TError,
    Key,
    AuthEmailBody,
    Awaited<ReturnType<typeof authEmailSignin>>
  > & { swrKey?: string };
}) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getAuthEmailSigninMutationKey();
  const swrFn = getAuthEmailSigninMutationFetcher();

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Verify an email address using a token that was emailed to one of the
account's email addresses either set via sign up or added later.

 */
export const authEmailVerify = (authEmailVerifyBody: AuthEmailVerifyBody) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/email/verify`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: authEmailVerifyBody,
  });
};

export const getAuthEmailVerifyMutationFetcher = () => {
  return (
    _: Key,
    { arg }: { arg: AuthEmailVerifyBody },
  ): Promise<AuthSuccessOKResponse> => {
    return authEmailVerify(arg);
  };
};
export const getAuthEmailVerifyMutationKey = () =>
  [`/auth/email/verify`] as const;

export type AuthEmailVerifyMutationResult = NonNullable<
  Awaited<ReturnType<typeof authEmailVerify>>
>;
export type AuthEmailVerifyMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useAuthEmailVerify = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(options?: {
  swr?: SWRMutationConfiguration<
    Awaited<ReturnType<typeof authEmailVerify>>,
    TError,
    Key,
    AuthEmailVerifyBody,
    Awaited<ReturnType<typeof authEmailVerify>>
  > & { swrKey?: string };
}) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getAuthEmailVerifyMutationKey();
  const swrFn = getAuthEmailVerifyMutationFetcher();

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * OAuth2 callback.
 */
export const oAuthProviderCallback = (
  oauthProvider: string,
  oAuthProviderCallbackBody: OAuthProviderCallbackBody,
) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/oauth/${oauthProvider}/callback`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: oAuthProviderCallbackBody,
  });
};

export const getOAuthProviderCallbackMutationFetcher = (
  oauthProvider: string,
) => {
  return (
    _: Key,
    { arg }: { arg: OAuthProviderCallbackBody },
  ): Promise<AuthSuccessOKResponse> => {
    return oAuthProviderCallback(oauthProvider, arg);
  };
};
export const getOAuthProviderCallbackMutationKey = (oauthProvider: string) =>
  [`/auth/oauth/${oauthProvider}/callback`] as const;

export type OAuthProviderCallbackMutationResult = NonNullable<
  Awaited<ReturnType<typeof oAuthProviderCallback>>
>;
export type OAuthProviderCallbackMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useOAuthProviderCallback = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  oauthProvider: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof oAuthProviderCallback>>,
      TError,
      Key,
      OAuthProviderCallbackBody,
      Awaited<ReturnType<typeof oAuthProviderCallback>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getOAuthProviderCallbackMutationKey(oauthProvider);
  const swrFn = getOAuthProviderCallbackMutationFetcher(oauthProvider);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Start the WebAuthn registration process by requesting a credential.

 */
export const webAuthnRequestCredential = (accountHandle: string) => {
  return fetcher<WebAuthnRequestCredentialOKResponse>({
    url: `/auth/webauthn/make/${accountHandle}`,
    method: "GET",
  });
};

export const getWebAuthnRequestCredentialKey = (accountHandle: string) =>
  [`/auth/webauthn/make/${accountHandle}`] as const;

export type WebAuthnRequestCredentialQueryResult = NonNullable<
  Awaited<ReturnType<typeof webAuthnRequestCredential>>
>;
export type WebAuthnRequestCredentialQueryError =
  | BadRequestResponse
  | InternalServerErrorResponse;

export const useWebAuthnRequestCredential = <
  TError = BadRequestResponse | InternalServerErrorResponse,
>(
  accountHandle: string,
  options?: {
    swr?: SWRConfiguration<
      Awaited<ReturnType<typeof webAuthnRequestCredential>>,
      TError
    > & { swrKey?: Key; enabled?: boolean };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false && !!accountHandle;
  const swrKey =
    swrOptions?.swrKey ??
    (() => (isEnabled ? getWebAuthnRequestCredentialKey(accountHandle) : null));
  const swrFn = () => webAuthnRequestCredential(accountHandle);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(
    swrKey,
    swrFn,
    swrOptions,
  );

  return {
    swrKey,
    ...query,
  };
};
/**
 * Complete WebAuthn registration by creating a new credential.
 */
export const webAuthnMakeCredential = (
  webAuthnMakeCredentialBody: WebAuthnMakeCredentialBody,
  params?: WebAuthnMakeCredentialParams,
) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/webauthn/make`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: webAuthnMakeCredentialBody,
    params,
  });
};

export const getWebAuthnMakeCredentialMutationFetcher = (
  params?: WebAuthnMakeCredentialParams,
) => {
  return (
    _: Key,
    { arg }: { arg: WebAuthnMakeCredentialBody },
  ): Promise<AuthSuccessOKResponse> => {
    return webAuthnMakeCredential(arg, params);
  };
};
export const getWebAuthnMakeCredentialMutationKey = (
  params?: WebAuthnMakeCredentialParams,
) => [`/auth/webauthn/make`, ...(params ? [params] : [])] as const;

export type WebAuthnMakeCredentialMutationResult = NonNullable<
  Awaited<ReturnType<typeof webAuthnMakeCredential>>
>;
export type WebAuthnMakeCredentialMutationError =
  | BadRequestResponse
  | InternalServerErrorResponse;

export const useWebAuthnMakeCredential = <
  TError = BadRequestResponse | InternalServerErrorResponse,
>(
  params?: WebAuthnMakeCredentialParams,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof webAuthnMakeCredential>>,
      TError,
      Key,
      WebAuthnMakeCredentialBody,
      Awaited<ReturnType<typeof webAuthnMakeCredential>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getWebAuthnMakeCredentialMutationKey(params);
  const swrFn = getWebAuthnMakeCredentialMutationFetcher(params);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Start the WebAuthn assertion for an existing account.
 */
export const webAuthnGetAssertion = (accountHandle: string) => {
  return fetcher<WebAuthnGetAssertionOKResponse>({
    url: `/auth/webauthn/assert/${accountHandle}`,
    method: "GET",
  });
};

export const getWebAuthnGetAssertionKey = (accountHandle: string) =>
  [`/auth/webauthn/assert/${accountHandle}`] as const;

export type WebAuthnGetAssertionQueryResult = NonNullable<
  Awaited<ReturnType<typeof webAuthnGetAssertion>>
>;
export type WebAuthnGetAssertionQueryError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useWebAuthnGetAssertion = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(
  accountHandle: string,
  options?: {
    swr?: SWRConfiguration<
      Awaited<ReturnType<typeof webAuthnGetAssertion>>,
      TError
    > & { swrKey?: Key; enabled?: boolean };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false && !!accountHandle;
  const swrKey =
    swrOptions?.swrKey ??
    (() => (isEnabled ? getWebAuthnGetAssertionKey(accountHandle) : null));
  const swrFn = () => webAuthnGetAssertion(accountHandle);

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(
    swrKey,
    swrFn,
    swrOptions,
  );

  return {
    swrKey,
    ...query,
  };
};
/**
 * Complete the credential assertion and sign in to an account.
 */
export const webAuthnMakeAssertion = (
  webAuthnMakeAssertionBody: WebAuthnMakeAssertionBody,
) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/webauthn/assert`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: webAuthnMakeAssertionBody,
  });
};

export const getWebAuthnMakeAssertionMutationFetcher = () => {
  return (
    _: Key,
    { arg }: { arg: WebAuthnMakeAssertionBody },
  ): Promise<AuthSuccessOKResponse> => {
    return webAuthnMakeAssertion(arg);
  };
};
export const getWebAuthnMakeAssertionMutationKey = () =>
  [`/auth/webauthn/assert`] as const;

export type WebAuthnMakeAssertionMutationResult = NonNullable<
  Awaited<ReturnType<typeof webAuthnMakeAssertion>>
>;
export type WebAuthnMakeAssertionMutationError =
  | UnauthorisedResponse
  | NotFoundResponse
  | InternalServerErrorResponse;

export const useWebAuthnMakeAssertion = <
  TError =
    | UnauthorisedResponse
    | NotFoundResponse
    | InternalServerErrorResponse,
>(options?: {
  swr?: SWRMutationConfiguration<
    Awaited<ReturnType<typeof webAuthnMakeAssertion>>,
    TError,
    Key,
    WebAuthnMakeAssertionBody,
    Awaited<ReturnType<typeof webAuthnMakeAssertion>>
  > & { swrKey?: string };
}) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getWebAuthnMakeAssertionMutationKey();
  const swrFn = getWebAuthnMakeAssertionMutationFetcher();

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Start the authentication flow with a phone number. The handler will send
a one-time code to the provided phone number which must then be sent to
the other phone endpoint to verify the number and validate the account.

 */
export const phoneRequestCode = (
  phoneRequestCodeBody: PhoneRequestCodeBody,
  params?: PhoneRequestCodeParams,
) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/phone`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: phoneRequestCodeBody,
    params,
  });
};

export const getPhoneRequestCodeMutationFetcher = (
  params?: PhoneRequestCodeParams,
) => {
  return (
    _: Key,
    { arg }: { arg: PhoneRequestCodeBody },
  ): Promise<AuthSuccessOKResponse> => {
    return phoneRequestCode(arg, params);
  };
};
export const getPhoneRequestCodeMutationKey = (
  params?: PhoneRequestCodeParams,
) => [`/auth/phone`, ...(params ? [params] : [])] as const;

export type PhoneRequestCodeMutationResult = NonNullable<
  Awaited<ReturnType<typeof phoneRequestCode>>
>;
export type PhoneRequestCodeMutationError =
  | BadRequestResponse
  | InternalServerErrorResponse;

export const usePhoneRequestCode = <
  TError = BadRequestResponse | InternalServerErrorResponse,
>(
  params?: PhoneRequestCodeParams,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof phoneRequestCode>>,
      TError,
      Key,
      PhoneRequestCodeBody,
      Awaited<ReturnType<typeof phoneRequestCode>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getPhoneRequestCodeMutationKey(params);
  const swrFn = getPhoneRequestCodeMutationFetcher(params);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Complete the phone number authentication flow by submitting the one-time
code that was sent to the user's phone.

 */
export const phoneSubmitCode = (
  accountHandle: string,
  phoneSubmitCodeBody: PhoneSubmitCodeBody,
) => {
  return fetcher<AuthSuccessOKResponse>({
    url: `/auth/phone/${accountHandle}`,
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    data: phoneSubmitCodeBody,
  });
};

export const getPhoneSubmitCodeMutationFetcher = (accountHandle: string) => {
  return (
    _: Key,
    { arg }: { arg: PhoneSubmitCodeBody },
  ): Promise<AuthSuccessOKResponse> => {
    return phoneSubmitCode(accountHandle, arg);
  };
};
export const getPhoneSubmitCodeMutationKey = (accountHandle: string) =>
  [`/auth/phone/${accountHandle}`] as const;

export type PhoneSubmitCodeMutationResult = NonNullable<
  Awaited<ReturnType<typeof phoneSubmitCode>>
>;
export type PhoneSubmitCodeMutationError =
  | BadRequestResponse
  | InternalServerErrorResponse;

export const usePhoneSubmitCode = <
  TError = BadRequestResponse | InternalServerErrorResponse,
>(
  accountHandle: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof phoneSubmitCode>>,
      TError,
      Key,
      PhoneSubmitCodeBody,
      Awaited<ReturnType<typeof phoneSubmitCode>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getPhoneSubmitCodeMutationKey(accountHandle);
  const swrFn = getPhoneSubmitCodeMutationFetcher(accountHandle);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * List all access keys for the authenticated account or all access keys
that have been issued for the entire instance if and only if the request
parameters specify all keys and the requesting account is an admin.

 */
export const accessKeyList = () => {
  return fetcher<AccessKeyListOKResponse>({
    url: `/auth/access-keys`,
    method: "GET",
  });
};

export const getAccessKeyListKey = () => [`/auth/access-keys`] as const;

export type AccessKeyListQueryResult = NonNullable<
  Awaited<ReturnType<typeof accessKeyList>>
>;
export type AccessKeyListQueryError =
  | BadRequestResponse
  | ForbiddenResponse
  | InternalServerErrorResponse;

export const useAccessKeyList = <
  TError = BadRequestResponse | ForbiddenResponse | InternalServerErrorResponse,
>(options?: {
  swr?: SWRConfiguration<Awaited<ReturnType<typeof accessKeyList>>, TError> & {
    swrKey?: Key;
    enabled?: boolean;
  };
}) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false;
  const swrKey =
    swrOptions?.swrKey ?? (() => (isEnabled ? getAccessKeyListKey() : null));
  const swrFn = () => accessKeyList();

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(
    swrKey,
    swrFn,
    swrOptions,
  );

  return {
    swrKey,
    ...query,
  };
};
/**
 * Create a new access key for the authenticated account. Access keys are
used to authenticate API requests on behalf of the account in a more
granular and service-friendly way than a session cookie.

Access keys share the same roles and permissions as the owning account
and only provide a way to use an `Authorization` header as an way of
interacting with the Storyden API.

Access keys also allow an expiry date to be set to limit how long a key
can be used to authenticate against the API.

 */
export const accessKeyCreate = (accessKeyCreateBody: AccessKeyCreateBody) => {
  return fetcher<AccessKeyCreateOKResponse>({
    url: `/auth/access-keys`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: accessKeyCreateBody,
  });
};

export const getAccessKeyCreateMutationFetcher = () => {
  return (
    _: Key,
    { arg }: { arg: AccessKeyCreateBody },
  ): Promise<AccessKeyCreateOKResponse> => {
    return accessKeyCreate(arg);
  };
};
export const getAccessKeyCreateMutationKey = () =>
  [`/auth/access-keys`] as const;

export type AccessKeyCreateMutationResult = NonNullable<
  Awaited<ReturnType<typeof accessKeyCreate>>
>;
export type AccessKeyCreateMutationError =
  | BadRequestResponse
  | ForbiddenResponse
  | InternalServerErrorResponse;

export const useAccessKeyCreate = <
  TError = BadRequestResponse | ForbiddenResponse | InternalServerErrorResponse,
>(options?: {
  swr?: SWRMutationConfiguration<
    Awaited<ReturnType<typeof accessKeyCreate>>,
    TError,
    Key,
    AccessKeyCreateBody,
    Awaited<ReturnType<typeof accessKeyCreate>>
  > & { swrKey?: string };
}) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey = swrOptions?.swrKey ?? getAccessKeyCreateMutationKey();
  const swrFn = getAccessKeyCreateMutationFetcher();

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Revoke an access key. This will immediately invalidate the key and it
will no longer be usable for authentication.

 */
export const accessKeyDelete = (accessKeyId: string) => {
  return fetcher<NoContentResponse>({
    url: `/auth/access-keys/${accessKeyId}`,
    method: "DELETE",
  });
};

export const getAccessKeyDeleteMutationFetcher = (accessKeyId: string) => {
  return (_: Key, __: { arg: Arguments }): Promise<NoContentResponse> => {
    return accessKeyDelete(accessKeyId);
  };
};
export const getAccessKeyDeleteMutationKey = (accessKeyId: string) =>
  [`/auth/access-keys/${accessKeyId}`] as const;

export type AccessKeyDeleteMutationResult = NonNullable<
  Awaited<ReturnType<typeof accessKeyDelete>>
>;
export type AccessKeyDeleteMutationError =
  | BadRequestResponse
  | ForbiddenResponse
  | InternalServerErrorResponse;

export const useAccessKeyDelete = <
  TError = BadRequestResponse | ForbiddenResponse | InternalServerErrorResponse,
>(
  accessKeyId: string,
  options?: {
    swr?: SWRMutationConfiguration<
      Awaited<ReturnType<typeof accessKeyDelete>>,
      TError,
      Key,
      Arguments,
      Awaited<ReturnType<typeof accessKeyDelete>>
    > & { swrKey?: string };
  },
) => {
  const { swr: swrOptions } = options ?? {};

  const swrKey =
    swrOptions?.swrKey ?? getAccessKeyDeleteMutationKey(accessKeyId);
  const swrFn = getAccessKeyDeleteMutationFetcher(accessKeyId);

  const query = useSWRMutation(swrKey, swrFn, swrOptions);

  return {
    swrKey,
    ...query,
  };
};
/**
 * Remove cookies from requesting client.
 */
export const authProviderLogout = () => {
  return fetcher<void>({ url: `/auth/logout`, method: "GET" });
};

export const getAuthProviderLogoutKey = () => [`/auth/logout`] as const;

export type AuthProviderLogoutQueryResult = NonNullable<
  Awaited<ReturnType<typeof authProviderLogout>>
>;
export type AuthProviderLogoutQueryError =
  | BadRequestResponse
  | InternalServerErrorResponse;

export const useAuthProviderLogout = <
  TError = BadRequestResponse | InternalServerErrorResponse,
>(options?: {
  swr?: SWRConfiguration<
    Awaited<ReturnType<typeof authProviderLogout>>,
    TError
  > & { swrKey?: Key; enabled?: boolean };
}) => {
  const { swr: swrOptions } = options ?? {};

  const isEnabled = swrOptions?.enabled !== false;
  const swrKey =
    swrOptions?.swrKey ??
    (() => (isEnabled ? getAuthProviderLogoutKey() : null));
  const swrFn = () => authProviderLogout();

  const query = useSwr<Awaited<ReturnType<typeof swrFn>>, TError>(
    swrKey,
    swrFn,
    swrOptions,
  );

  return {
    swrKey,
    ...query,
  };
};
