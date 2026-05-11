import { UnreadyBanner } from "src/components/site/Unready";

import { LinkButton } from "@/components/ui/link-button";
import { useI18n } from "@/i18n/provider";
import { VStack } from "@/styled-system/jsx";

export function NodeNotFoundError() {
  const { t } = useI18n();

  return (
    <VStack p="4" h="dvh" justify="center">
      <VStack maxW="sm" minH="60" gap="8">
        <UnreadyBanner
          error={t("The link to this page did not lead anywhere.")}
        />
        <LinkButton variant="subtle" href="/l">
          {t("Library")}
        </LinkButton>
      </VStack>
    </VStack>
  );
}
