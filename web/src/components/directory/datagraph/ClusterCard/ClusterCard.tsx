import { Cluster } from "src/api/openapi/schemas";
import { Empty } from "src/components/site/Empty";
import {
  DirectoryPath,
  joinDirectoryPath,
} from "src/screens/directory/datagraph/useDirectoryPath";
import { Heading3 } from "src/theme/components/Heading/Index";

import styles from "./ClusterCard.module.css";

import { cx } from "@/styled-system/css";
import {
  Box,
  Center,
  LinkBox,
  LinkOverlay,
  VStack,
  styled,
} from "@/styled-system/jsx";
import { Card } from "@/styled-system/patterns";

export type Props = {
  cluster: Cluster;
  directoryPath: DirectoryPath;
};

export function ClusterCard({ cluster, directoryPath }: Props) {
  const slug = joinDirectoryPath(directoryPath, cluster.slug);
  const asset = cluster.assets?.[0];

  return (
    <styled.article containerType="inline-size" w="full">
      <LinkBox
        className={cx(
          Card({ kind: "edge", display: "grid" }),
          styles["container"],
        )}
        w="full"
        overflow="hidden"
      >
        {asset && (
          <Box className={styles["background-blur"]} gridRow="1" height="full">
            <styled.img
              gridRow="1"
              src={asset.url}
              width="full"
              height="full"
              objectPosition="center"
              objectFit="cover"
              blur="xl"
              opacity="3"
              filter="auto"
            />
          </Box>
        )}

        {asset ? (
          <styled.img
            className={styles["image"]}
            src={asset.url}
            width="full"
            height="full"
            objectPosition="center"
            objectFit="cover"
            zIndex="tooltip"
          />
        ) : (
          <Center display={{ base: "none", md: "flex" }}>
            <Empty>no image</Empty>
          </Center>
        )}

        <VStack
          className={styles["title"]}
          alignItems="center"
          justifyContent="start"
          background="cardBackgroundGradient"
        >
          <Box w="full" height="min" p="2" wordBreak="break-all">
            <Heading3 className="fluid-font-size" lineClamp={2}>
              <LinkOverlay href={`/directory/${slug}`}>
                {cluster.name}
              </LinkOverlay>
            </Heading3>
            {cluster.description ? (
              <styled.p lineClamp={2}>{cluster.description}</styled.p>
            ) : (
              <styled.p color="fg.subtle" fontStyle="italic">
                (no description)
              </styled.p>
            )}
          </Box>
        </VStack>
      </LinkBox>
    </styled.article>
  );
}
