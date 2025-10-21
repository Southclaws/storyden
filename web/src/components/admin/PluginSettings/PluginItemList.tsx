import { PluginList } from "@/api/openapi-schema";
import { EmptyState } from "@/components/site/EmptyState";
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
  return (
    <Center w="full" h="64">
      <EmptyState hideContributionLabel>
        No plugins have been installed yet.
        <span>
          <a
            className={css({
              color: "fg.emphasized",
              _hover: { textDecoration: "underline" },
            })}
            href="https://www.storyden.org/docs/introduction/plugins"
          >
            Learn more
          </a>{" "}
          about Storyden plugins.
        </span>
      </EmptyState>
    </Center>
  );
}
