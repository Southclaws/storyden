import { useEffect, useState } from "react";
import { useWatch } from "react-hook-form";

import { getGetInfoKey } from "@/api/openapi-client/misc";
import { InfoTip } from "@/components/site/InfoTip";
import { Admonition } from "@/components/ui/admonition";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { SliderField } from "@/components/ui/form/SliderField";
import { Heading } from "@/components/ui/heading";
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
  Props,
  formatSeconds,
  useSystemSettings,
} from "./useSystemSettings";

export function SystemSettingsForm(props: Props) {
  const { control, formState, onSubmit } = useSystemSettings(props);
  const [showError, setShowError] = useState(true);

  // Watch form values for live preview
  const rateLimit = useWatch({ control, name: "rate_limit" });
  const rateLimitPeriod = useWatch({ control, name: "rate_limit_period" });
  const rateLimitBucket = useWatch({ control, name: "rate_limit_bucket" });
  const rateLimitGuestCost = useWatch({
    control,
    name: "rate_limit_guest_cost",
  });

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
      </CardBox>
    </styled.form>
  );
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
