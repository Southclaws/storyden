import { useEffect, useState } from "react";

import { getSpec } from "@/api/openapi-client/misc";
import { Unready } from "@/components/site/Unready";
import { Admonition } from "@/components/ui/admonition";
import { Stack } from "@/styled-system/jsx";

import { OperationCostOverrideAddForm } from "./OperationCostOverrideAddForm";
import { OperationCostOverrideItem } from "./OperationCostOverrideItem";

type OperationCostOverridesListProps = {
  value: Record<string, number>;
  rateLimit: number;
  rateLimitPeriod: number;
  onChange: (value: Record<string, number>) => void;
};

type OperationOverride = {
  operationId: string;
  cost: number;
};

export function OperationCostOverridesList({
  value,
  rateLimit,
  rateLimitPeriod,
  onChange,
}: OperationCostOverridesListProps) {
  const [operations, setOperations] = useState<string[] | null>(null);
  const [selectedOperation, setSelectedOperation] = useState<string>("");
  const [costValue, setCostValue] = useState<number>(1);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    getSpec()
      .then((spec) => {
        const operationIds: string[] = [];

        if (spec["paths"]) {
          Object.values(spec["paths"]).forEach((pathItem: any) => {
            ["get", "post", "put", "patch", "delete"].forEach((method) => {
              if (pathItem[method]?.operationId) {
                operationIds.push(pathItem[method].operationId);
              }
            });
          });
        }

        setOperations(operationIds.sort());
      })
      .catch(() => {
        setError(
          "Failed to load API operations. Please refresh the page to try again.",
        );
        setOperations([]);
      });
  }, []);

  const overrides = value || {};
  const overridesList: OperationOverride[] = Object.entries(overrides).map(
    ([operationId, cost]) => ({
      operationId,
      cost,
    }),
  );

  const handleAdd = () => {
    if (!selectedOperation || costValue <= 0) return;

    onChange({
      ...overrides,
      [selectedOperation]: costValue,
    });

    setSelectedOperation("");
    setCostValue(1);
  };

  const handleCostChange = (operationId: string, cost: number) => {
    onChange({
      ...overrides,
      [operationId]: cost,
    });
  };

  const handleRemove = (operationId: string) => {
    const newOverrides = { ...overrides };
    delete newOverrides[operationId];
    onChange(newOverrides);
  };

  const availableOperations = operations
    ? operations.filter((op) => !overrides[op])
    : [];

  return (
    <Stack gap="2">
      {error && (
        <Admonition
          value={true}
          kind="failure"
          title="Error loading operations"
          onChange={() => setError(null)}
        >
          {error}
        </Admonition>
      )}

      {!operations && !error && <Unready />}

      {operations && (
        <OperationCostOverrideAddForm
          rateLimit={rateLimit}
          rateLimitPeriod={rateLimitPeriod}
          availableOperations={availableOperations}
          selectedOperation={selectedOperation}
          costValue={costValue}
          onOperationChange={setSelectedOperation}
          onCostChange={setCostValue}
          onAdd={handleAdd}
        />
      )}

      {overridesList.length > 0 && (
        <Stack gap="2">
          {overridesList.map((override) => (
            <OperationCostOverrideItem
              key={override.operationId}
              operationId={override.operationId}
              cost={override.cost}
              rateLimit={rateLimit}
              rateLimitPeriod={rateLimitPeriod}
              onCostChange={handleCostChange}
              onRemove={handleRemove}
            />
          ))}
        </Stack>
      )}
    </Stack>
  );
}
