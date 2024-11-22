import { flatten, groupBy, keyBy, values } from "lodash/fp";

import {
  AccountAuthMethod,
  AccountAuthMethodList,
  AuthProvider,
  AuthProviderList,
} from "src/api/openapi-schema";

import { filterWithLink } from "./oauth";

const groupProviders = keyBy<AuthProvider>("provider");

export function groupAuthProviders(providers: AuthProviderList) {
  // pull out password and phone, if present, the rest are OAuth2 providers.
  const { password, phone, webauthn, ...rest } = groupProviders(providers);

  const oauth = filterWithLink(values(rest));

  return {
    password: Boolean(password),
    phone: Boolean(phone),
    webauthn: Boolean(webauthn),
    oauth,
  };
}

const groupMethods = groupBy((v: AccountAuthMethod) => v.provider.provider);

export function groupAuthMethods(methods: AccountAuthMethodList) {
  const { password, phone, webauthn, ...rest } = groupMethods(methods);

  return {
    password: password,
    phone: phone,
    webauthn: webauthn,
    methods: flatten(values(rest)),
  };
}
