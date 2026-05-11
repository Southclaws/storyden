import { PluginList } from "@/api/openapi-schema";
import { EmptyState } from "@/components/site/EmptyState";
import { useI18n } from "@/i18n/provider";
import { css } from "@/styled-system/css";
import { Center } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { PluginItem } from "./PluginItem";

type Props = {
  plugins: PluginList;
};

export function PluginItemList({ plugins }: Props) {
  if (plugins.length === 0) {
    return <PluginsEmptyState />;
  }

  return (
    <ul className={lstack({ gap: "3" })}>
      {plugins.map((plugin) => (
        <PluginItem key={plugin.id} plugin={plugin} />
      ))}
    </ul>
  );
}

function PluginsEmptyState() {
  const { t } = useI18n();
  return (
    <Center w="full" h="64">
      <EmptyState hideContributionLabel>
        {t("No plugins have been installed yet.")}
        <span>
          <a
            className={css({
              color: "fg.emphasized",
              _hover: { textDecoration: "underline" },
            })}
            href="https://www.storyden.org/docs/introduction/plugins"
          >
            {t("Learn more")}
          </a>{" "}
          {t("about Storyden plugins.")}
        </span>
      </EmptyState>
    </Center>
  );
}
