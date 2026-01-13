import {
  Controller,
  type ControllerProps,
  type FieldValues,
} from "react-hook-form";

import { OperationCostOverridesList } from "./OperationCostOverridesList";
import {
  DEFAULT_RATE_LIMIT,
  DEFAULT_RATE_LIMIT_PERIOD,
} from "./useSystemSettings";

type OperationCostOverridesProps<T extends FieldValues> = Omit<
  ControllerProps<T>,
  "render"
> & {
  rateLimit?: number;
  rateLimitPeriod?: number;
};

export function OperationCostOverrides<T extends FieldValues>({
  rateLimit = DEFAULT_RATE_LIMIT,
  rateLimitPeriod = DEFAULT_RATE_LIMIT_PERIOD,
  ...controllerProps
}: OperationCostOverridesProps<T>) {
  return (
    <Controller<T>
      {...controllerProps}
      render={({ field }) => (
        <OperationCostOverridesList
          value={(field.value as Record<string, number>) || {}}
          rateLimit={rateLimit}
          rateLimitPeriod={rateLimitPeriod}
          onChange={field.onChange}
        />
      )}
    />
  );
}
