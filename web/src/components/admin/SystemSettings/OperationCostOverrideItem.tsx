import { CardBox } from "@/components/ui/card-box";
import { IconButton } from "@/components/ui/icon-button";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { NumberInput } from "@/components/ui/number-input";
import { Text } from "@/components/ui/text";
import { HStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { formatSeconds } from "@/utils/date";

type OperationCostOverrideItemProps = {
  operationId: string;
  cost: number;
  rateLimit: number;
  rateLimitPeriod: number;
  onCostChange: (operationId: string, cost: number) => void;
  onRemove: (operationId: string) => void;
};

export function OperationCostOverrideItem({
  operationId,
  cost,
  rateLimit,
  rateLimitPeriod,
  onCostChange,
  onRemove,
}: OperationCostOverrideItemProps) {
  const effectiveLimit = Math.floor(rateLimit / cost);

  return (
    <CardBox
      className={lstack()}
      w="full"
      borderRadius="sm"
      bg="background.inset"
    >
      <WStack w="full" gap="2" alignItems="center">
        <styled.code fontSize="sm" fontWeight="semibold" flex="1">
          {operationId}
        </styled.code>

        <HStack gap="2">
          <NumberInput
            width="16"
            size="sm"
            min={1}
            max={100}
            value={cost.toString()}
            onValueChange={(details) => {
              const newCost = parseInt(details.value, 10);
              if (!isNaN(newCost) && newCost > 0) {
                onCostChange(operationId, newCost);
              }
            }}
          />
          <IconButton
            size="sm"
            variant="ghost"
            onClick={() => onRemove(operationId)}
          >
            <DeleteIcon />
          </IconButton>
        </HStack>
      </WStack>
      <Text variant="metadata">
        Can be performed{" "}
        <styled.strong color="status.info.content">
          {effectiveLimit}
        </styled.strong>{" "}
        times every{" "}
        <styled.strong color="status.info.content">
          {formatSeconds(rateLimitPeriod)}
        </styled.strong>
      </Text>
    </CardBox>
  );
}
