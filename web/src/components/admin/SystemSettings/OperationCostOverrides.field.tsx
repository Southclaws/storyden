import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import { OperationCostOverridesList } from "./OperationCostOverridesList";
import {
  DEFAULT_RATE_LIMIT,
  DEFAULT_RATE_LIMIT_PERIOD,
} from "./useSystemSettings";

type CostOverrideFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  Record<string, number> | undefined
>;

type OperationCostOverridesFieldProps<
  TFieldValues extends FieldValues,
  TName extends CostOverrideFieldName<TFieldValues> =
    CostOverrideFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> & {
  rateLimit?: number;
  rateLimitPeriod?: number;
};

export function OperationCostOverridesField<
  TFieldValues extends FieldValues,
  TName extends CostOverrideFieldName<TFieldValues> =
    CostOverrideFieldName<TFieldValues>,
>({
  rateLimit = DEFAULT_RATE_LIMIT,
  rateLimitPeriod = DEFAULT_RATE_LIMIT_PERIOD,
  ...controllerProps
}: OperationCostOverridesFieldProps<TFieldValues, TName>) {
  return (
    <Controller<TFieldValues, TName>
      {...controllerProps}
      render={({ field }) => (
        <OperationCostOverridesList
          value={
            typeof field.value === "object" && field.value !== null
              ? field.value
              : {}
          }
          rateLimit={rateLimit}
          rateLimitPeriod={rateLimitPeriod}
          onChange={field.onChange}
        />
      )}
    />
  );
}
