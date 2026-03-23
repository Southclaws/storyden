import { createListCollection } from "@ark-ui/react";
import { useEffect, useState } from "react";
import { useWatch } from "react-hook-form";

import { adminSettingsGet } from "@/api/openapi-client/admin";
import { getGetInfoKey, getSession } from "@/api/openapi-client/misc";
import { ClientInfo, NetworkHeadersSample } from "@/api/openapi-schema";
import { InfoTip } from "@/components/site/InfoTip";
import { Admonition } from "@/components/ui/admonition";
import * as Alert from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { SelectField } from "@/components/ui/form/SelectField";
import { SliderField } from "@/components/ui/form/SliderField";
import { Heading } from "@/components/ui/heading";
import { Input } from "@/components/ui/input";
import { API_ADDRESS } from "@/config";
import { CardBox, HStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { deriveError } from "@/utils/error";

import { OperationCostOverrides } from "./OperationCostOverrides";
import {
  DEFAULT_RATE_LIMIT,
  DEFAULT_RATE_LIMIT_BUCKET,
  DEFAULT_RATE_LIMIT_GUEST_COST,
  DEFAULT_RATE_LIMIT_PERIOD,
  Form,
  Props,
  formatSeconds,
  useSystemSettings,
} from "./useSystemSettings";

const CLIENT_IP_MODE_COLLECTION = createListCollection({
  items: [
    {
      label: "Raw IP (default)",
      value: "remote_addr",
    },
    {
      label: "Single trusted header",
      value: "single_header",
    },
    {
      label: "X-Forwarded-For with trusted proxy CIDRs",
      value: "xff_trusted_proxies",
    },
  ],
});

export function SystemSettingsForm(props: Props) {
  const { control, register, formState, onSubmit } = useSystemSettings(props);
  const [showError, setShowError] = useState(true);

  // Watch form values for live preview
  const rateLimit = useWatch({ control, name: "rate_limit" });
  const rateLimitPeriod = useWatch({ control, name: "rate_limit_period" });
  const rateLimitBucket = useWatch({ control, name: "rate_limit_bucket" });
  const rateLimitGuestCost = useWatch({
    control,
    name: "rate_limit_guest_cost",
  });
  const clientIPMode = useWatch({ control, name: "client_ip_mode" });

  const guestRateLimit = Math.floor(
    (rateLimit ?? DEFAULT_RATE_LIMIT) /
      (rateLimitGuestCost ?? DEFAULT_RATE_LIMIT_GUEST_COST),
  );

  const memberRequestsPerMinute = Math.round(
    ((rateLimit ?? DEFAULT_RATE_LIMIT) /
      (rateLimitPeriod ?? DEFAULT_RATE_LIMIT_PERIOD)) *
      60,
  );

  const guestRequestsPerMinute = Math.round(
    (guestRateLimit / (rateLimitPeriod ?? DEFAULT_RATE_LIMIT_PERIOD)) * 60,
  );

  const hasErrors = Object.keys(formState.errors).length > 0;
  const errorMessages = Object.entries(formState.errors)
    .map(([field, error]) => `${field}: ${deriveError(error?.message)}`)
    .join(", ");

  // Reset error visibility when errors change
  useEffect(() => {
    if (hasErrors) {
      setShowError(true);
    }
  }, [hasErrors]);

  return (
    <styled.form
      width="full"
      display="flex"
      flexDirection="column"
      gap="4"
      onSubmit={onSubmit}
    >
      <CardBox className={lstack()} gap="4">
        <WStack>
          <Heading size="md">System settings</Heading>
          <Button type="submit" loading={formState.isSubmitting}>
            Save
          </Button>
        </WStack>

        {hasErrors && showError && (
          <Admonition
            value={true}
            kind="failure"
            title="Form validation error"
            onChange={() => setShowError(false)}
          >
            {errorMessages}
          </Admonition>
        )}

        <Heading>Rate limits</Heading>

        <p>
          Rate limits help protect your installation from spam, DDoS attacks and
          content scraping. This is achieved by limiting the number of{" "}
          <strong>Operations</strong>{" "}
          <InfoTip title="What is an operation?">
            An "Operation" is a request to Storyden's backend. Loading a screen
            such as Home, a Thread or a Library Page usually involves 10-30
            request operations.
          </InfoTip>{" "}
          in a time period.
        </p>

        <p>
          Members: <styled.strong color="fg.info">{rateLimit}</styled.strong>{" "}
          operations every{" "}
          <styled.strong color="fg.info">
            {formatSeconds(rateLimitPeriod ?? DEFAULT_RATE_LIMIT_PERIOD)}
          </styled.strong>{" "}
          (~
          <styled.strong color="fg.info">
            {memberRequestsPerMinute}
          </styled.strong>{" "}
          requests per minute).
        </p>

        <p>
          Guests:{" "}
          <styled.strong color="fg.info">{guestRateLimit}</styled.strong>{" "}
          operations every{" "}
          <styled.strong color="fg.info">
            {formatSeconds(rateLimitPeriod ?? DEFAULT_RATE_LIMIT_PERIOD)}
          </styled.strong>{" "}
          (~
          <styled.strong color="fg.info">
            {guestRequestsPerMinute}
          </styled.strong>{" "}
          requests per minute).
        </p>

        <RateLimitTester />

        <FormControl>
          <SliderField
            control={control}
            name="rate_limit"
            label={`Rate limit: ${rateLimit} request units`}
            min={10}
            max={20000}
            step={10}
            sliderDefaultValue={DEFAULT_RATE_LIMIT}
            marks={[
              {
                value: DEFAULT_RATE_LIMIT,
                label: "Default",
              },
            ]}
          />
          <FormHelperText>
            The amount of requests that a user can make within the
            `rate_limit_period`.
          </FormHelperText>
        </FormControl>

        <FormControl>
          <SliderField
            control={control}
            name="rate_limit_period"
            label={`Rate limit period: ${formatSeconds(rateLimitPeriod)}`}
            min={60}
            max={86400}
            step={60}
            sliderDefaultValue={DEFAULT_RATE_LIMIT_PERIOD}
            marks={[
              {
                value: DEFAULT_RATE_LIMIT_PERIOD,
                label: "Default",
              },
            ]}
          />
          <FormHelperText>
            The period of time in which the `rate_limit` is applied. This is a
            sliding window, so the `rate_limit` is applied to the last
            `rate_limit_period` of requests.
          </FormHelperText>
        </FormControl>

        <FormControl>
          <SliderField
            control={control}
            name="rate_limit_bucket"
            label={`Rate limit bucket size: ${rateLimitBucket} seconds`}
            min={0}
            max={1200}
            step={60}
            sliderDefaultValue={DEFAULT_RATE_LIMIT_BUCKET}
            marks={[
              {
                value: DEFAULT_RATE_LIMIT_BUCKET,
                label: "Default",
              },
            ]}
          />
          <FormHelperText>
            The granularity of rate limit counter buckets. Lower values use more
            memory but provide more accurate rate limiting. Higher values use
            less memory but may allow short bursts of traffic above the rate
            limit.
          </FormHelperText>
        </FormControl>

        <FormControl>
          <SliderField
            control={control}
            name="rate_limit_guest_cost"
            label="Guest rate limit cost multiplier"
            min={1}
            max={10}
            step={1}
            sliderDefaultValue={DEFAULT_RATE_LIMIT_GUEST_COST}
          />
          <FormHelperText>
            The cost multiplier applied to unauthenticated guest visitors. For
            example, a value of 5 means each operation consumes 5 units from the
            guest&apos;s rate limit instead of 1, applying stricter limits to
            non-authenticated traffic.
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>Operation cost overrides</FormLabel>
          <FormHelperText>
            Configure custom cost multipliers for specific API operations.
            Higher costs reduce the number of requests allowed within the rate
            limit period.
          </FormHelperText>
          <OperationCostOverrides
            control={control}
            name="cost_overrides"
            rateLimit={rateLimit ?? DEFAULT_RATE_LIMIT}
            rateLimitPeriod={rateLimitPeriod ?? DEFAULT_RATE_LIMIT_PERIOD}
          />
        </FormControl>

        <Heading>Client IP strategy</Heading>

        <FormControl>
          <FormLabel>Client IP mode</FormLabel>
          <SelectField<Form, (typeof CLIENT_IP_MODE_COLLECTION.items)[number]>
            control={control}
            name="client_ip_mode"
            collection={CLIENT_IP_MODE_COLLECTION}
            placeholder="Select client IP mode"
          />
          <FormHelperText>
            Choose how Storyden resolves client addresses for request context.
            The default uses only RemoteAddr and does not trust forwarded
            headers. Header-based modes should only be used when your edge
            proxy/CDN strips or overwrites client-provided forwarding headers.
          </FormHelperText>
        </FormControl>

        {clientIPMode === "single_header" && (
          <FormControl>
            <FormLabel>Client IP header</FormLabel>
            <Input {...register("client_ip_header")} />
            <FormHelperText>
              Header to trust for the canonical client IP (for example
              CF-Connecting-IP, Fly-Client-IP, X-Real-IP). Do not use this mode
              unless this header is guaranteed to be injected by trusted
              infrastructure.
            </FormHelperText>
          </FormControl>
        )}

        {clientIPMode === "xff_trusted_proxies" && (
          <FormControl>
            <FormLabel>Trusted proxy CIDRs</FormLabel>
            <Input
              {...register("trusted_proxy_cidrs")}
              placeholder="10.0.0.0/8, 172.16.0.0/12"
            />
            <FormHelperText>
              Comma-separated CIDR ranges that are allowed to append XFF hops.
              Storyden will only trust XFF when RemoteAddr is in these ranges.
              Include every proxy hop in your chain to avoid collapsing users to
              a shared proxy IP.
            </FormHelperText>
          </FormControl>
        )}

        <ClientIPTester
          canRun={!formState.isDirty && !formState.isSubmitting}
          initialHeaders={props.settings.headers}
        />
      </CardBox>
    </styled.form>
  );
}

type ClientIPTesterProps = {
  canRun: boolean;
  initialHeaders?: NetworkHeadersSample;
};

function formatNetworkHeaderSample(
  label: string,
  headers: NetworkHeadersSample | null,
) {
  const direct = headers?.headers ?? {};
  const ssr = headers?.headers_ssr ?? {};
  const rawClientAddress = headers?.raw_client_address?.trim() ?? "";

  const directEntries = Object.entries(direct).sort(([a], [b]) =>
    a.localeCompare(b),
  );
  const ssrEntries = Object.entries(ssr).sort(([a], [b]) => a.localeCompare(b));

  if (
    directEntries.length === 0 &&
    ssrEntries.length === 0 &&
    !rawClientAddress
  ) {
    return (
      <>
        {label}
        <br />
        (none)
      </>
    );
  }

  return (
    <>
      {label}
      <br />
      Raw client address:
      <br />
      {rawClientAddress || "(none)"}
      <br />
      <br />
      Browser/API headers:
      <br />
      {directEntries.length === 0 && "(none)"}
      {directEntries.map(([name, value]) => (
        <span key={`${label}-direct-${name}`}>
          {name}: {value}
          <br />
        </span>
      ))}
      <br />
      SSR-forwarded headers:
      <br />
      {ssrEntries.length === 0 && "(none)"}
      {ssrEntries.map(([name, value]) => (
        <span key={`${label}-ssr-${name}`}>
          {name}: {value}
          <br />
        </span>
      ))}
    </>
  );
}

function ClientIPTester({ canRun, initialHeaders }: ClientIPTesterProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [ssrClientInfo, setSSRClientInfo] = useState<ClientInfo | null>(null);
  const [browserClientInfo, setBrowserClientInfo] = useState<ClientInfo | null>(
    null,
  );
  const [browserHeaderSample, setBrowserHeaderSample] =
    useState<NetworkHeadersSample | null>(initialHeaders ?? null);
  const [ssrHeaderSample, setSSRHeaderSample] =
    useState<NetworkHeadersSample | null>(null);

  async function loadClientIPTest() {
    setLoading(true);
    setError(null);

    try {
      const [ssrResp, browserResp, browserAdminSettings] = await Promise.all([
        fetch("/client-ip-test", {
          credentials: "include",
          cache: "no-store",
          mode: "cors",
        }),
        getSession(),
        adminSettingsGet(),
      ]);

      if (!ssrResp.ok) {
        const data = (await ssrResp.json().catch(() => undefined)) as
          | { message?: string }
          | undefined;
        throw new Error(
          data?.message ?? `SSR test request failed with ${ssrResp.status}`,
        );
      }

      const ssrData = (await ssrResp.json()) as {
        client: ClientInfo | null;
        headers?: NetworkHeadersSample | null;
      };
      setSSRClientInfo(ssrData.client ?? null);
      setBrowserClientInfo(browserResp.client ?? null);
      setSSRHeaderSample(ssrData.headers ?? null);
      setBrowserHeaderSample(browserAdminSettings.headers ?? null);
    } catch (err) {
      setSSRClientInfo(null);
      setBrowserClientInfo(null);
      setSSRHeaderSample(null);
      setBrowserHeaderSample(null);
      setError(deriveError(err));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void (async () => {
      await loadClientIPTest();
    })();
  }, []);

  const warnings = getClientIPWarnings(ssrClientInfo, browserClientInfo);

  return (
    <CardBox bgColor="bg.subtle" fontSize="xs" display="flex" gap="2">
      <styled.p>
        This client IP test runs automatically and compares what Storyden sees
        from an SSR-origin call and a browser-origin call.
      </styled.p>

      {error && <styled.p color="fg.error">{error}</styled.p>}

      {warnings.length > 0 && (
        <Alert.Root>
          <Alert.Content>
            <Alert.Title>Potential IP config issue</Alert.Title>
            <Alert.Description>
              <styled.ul>
                {warnings.map((warning) => (
                  <li key={warning}>{warning}</li>
                ))}
              </styled.ul>
            </Alert.Description>
          </Alert.Content>
        </Alert.Root>
      )}

      <styled.pre textWrap="wrap">
        Server HTML Render client.ip_address_ssr = '
        {ssrClientInfo?.ip_address_ssr ?? ""}'
        <br />
        Browser React Render client.ip_address = '
        {browserClientInfo?.ip_address ?? ""}'
        <br />
        <br />
        These are sampled network headers seen by Storyden while running the
        client IP test.
        <br />
        Use them to configure trusted client IP settings.
        <br />
        <br />
        {formatNetworkHeaderSample(
          "Browser/API sample (admin settings request):",
          browserHeaderSample,
        )}
        <br />
        <br />
        {formatNetworkHeaderSample(
          "SSR/API sample (SSR subrequest):",
          ssrHeaderSample,
        )}
      </styled.pre>

      <HStack justify="end">
        <Button
          type="button"
          variant="subtle"
          size="xs"
          onClick={() => {
            if (!canRun) {
              setError("Save settings before refreshing the client IP test.");
              return;
            }
            void loadClientIPTest();
          }}
          loading={loading}
          disabled={!canRun}
        >
          Refresh Client IP Test
        </Button>
      </HStack>
    </CardBox>
  );
}

function getClientIPWarnings(
  ssrClientInfo: ClientInfo | null,
  browserClientInfo: ClientInfo | null,
): string[] {
  const warnings: string[] = [];

  const ssrIP = ssrClientInfo?.ip_address_ssr?.trim() ?? "";
  const browserIP = browserClientInfo?.ip_address?.trim() ?? "";

  if (ssrIP && browserIP && ssrIP !== browserIP) {
    warnings.push(
      `SSR and browser resolved IPs differ (${ssrIP} vs ${browserIP}).`,
    );
  }

  const ipValues: Array<{ label: string; value: string }> = [
    {
      label: "SSR request resolved IP",
      value: ssrIP,
    },
    { label: "Browser resolved IP", value: browserIP },
  ];

  const internalByIP = new Map<string, string[]>();
  for (const candidate of ipValues) {
    const value = candidate.value.trim();
    if (value && isLikelyInternalIP(value)) {
      const labels = internalByIP.get(value) ?? [];
      labels.push(candidate.label);
      internalByIP.set(value, labels);
    }
  }

  for (const [value, labels] of internalByIP) {
    if (labels.length === 1) {
      warnings.push(`${labels[0]} looks internal/private (${value}).`);
    } else {
      warnings.push(
        `${labels.join(" and ")} look internal/private (${value}).`,
      );
    }
  }

  return warnings;
}

function isLikelyInternalIP(ip: string): boolean {
  const trimmed = ip.trim();
  if (!trimmed) return false;

  if (trimmed === "::1") return true;
  if (trimmed.includes(":")) {
    const lower = trimmed.toLowerCase();
    return (
      lower.startsWith("fc") ||
      lower.startsWith("fd") ||
      lower.startsWith("fe80:")
    );
  }

  const parts = trimmed.split(".");
  if (parts.length !== 4) return false;

  const octets = parts.map((p) => Number.parseInt(p, 10));
  if (octets.some((n) => Number.isNaN(n) || n < 0 || n > 255)) return false;

  const a = octets[0] ?? -1;
  const b = octets[1] ?? -1;
  if (a === 10) return true;
  if (a === 127) return true;
  if (a === 192 && b === 168) return true;
  if (a === 169 && b === 254) return true;
  if (a === 172 && b >= 16 && b <= 31) return true;
  if (a === 100 && b >= 64 && b <= 127) return true;

  return false;
}

function RateLimitTester() {
  const [xRateLimit, setXRateLimit] = useState<Record<string, any>>({});

  async function run(auth: boolean) {
    const resp = await fetch(`${API_ADDRESS}/api${getGetInfoKey()[0]}`, {
      credentials: auth ? "include" : "omit",
      mode: "cors",
    });
    const rateLimitLimit = resp.headers.get("x-ratelimit-limit");
    const rateLimitRemaining = resp.headers.get("x-ratelimit-remaining");
    const rateLimitReset = resp.headers.get("x-ratelimit-reset");
    setXRateLimit({
      rateLimitLimit: rateLimitLimit,
      rateLimitRemaining: rateLimitRemaining,
      rateLimitReset: rateLimitReset,
    });
  }

  const rateLimitLimit = xRateLimit["rateLimitLimit"];
  const rateLimitRemaining = xRateLimit["rateLimitRemaining"];
  const rateLimitReset = xRateLimit["rateLimitReset"];

  let estimatedRequestsPerMinute: number | null = null;
  if (rateLimitRemaining && rateLimitReset) {
    const resetDate = new Date(rateLimitReset);
    const now = new Date();
    const secondsUntilReset = Math.max(
      0,
      (resetDate.getTime() - now.getTime()) / 1000,
    );
    const minutesUntilReset = secondsUntilReset / 60;

    if (minutesUntilReset > 0) {
      estimatedRequestsPerMinute = Math.round(
        parseInt(rateLimitRemaining) / minutesUntilReset,
      );
    }
  }

  return (
    <CardBox bgColor="bg.subtle" fontSize="xs" display="flex" gap="2">
      <styled.p>
        This is your current rate limit status. Click the &quot;Test&quot;
        button to consume one request.
      </styled.p>

      <styled.pre textWrap="wrap">
        x-rate-limit-limit = '{rateLimitLimit}'<br />
        x-rate-limit-remaining = '{rateLimitRemaining}'<br />
        x-rate-limit-reset ={" "}
        <styled.span textWrap="nowrap">'{rateLimitReset}'</styled.span>
        <br />
        {estimatedRequestsPerMinute !== null && (
          <>
            ~{estimatedRequestsPerMinute} requests per minute until period
            reset.
          </>
        )}
      </styled.pre>

      <HStack justify="end">
        <Button
          type="button"
          variant="subtle"
          size="xs"
          onClick={() => run(true)}
        >
          Test as Member
        </Button>
        <Button
          type="button"
          variant="subtle"
          size="xs"
          onClick={() => run(false)}
        >
          Test as Guest
        </Button>
      </HStack>
    </CardBox>
  );
}
