import { TrashIcon } from "@heroicons/react/24/outline";

import { Button } from "src/theme/components/Button";

import styles from "./assetwall.module.css";

import { Box, BoxProps, HStack, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

import { Props, useEditableAssetWall } from "./useEditableAssetWall";

export function EditableAssetWall({
  initialAssets,
  editing,
  onUpload,
  onRemove,
  ...props
}: Props & BoxProps) {
  const { assets, handlers } = useEditableAssetWall({
    initialAssets,
    editing,
    onUpload,
    onRemove,
  });

  return (
    <Box className={styles["root"]} {...props}>
      <Box className={styles["grid"]}>
        {assets?.map((a) => (
          <Box key={a.id} className={styles["asset"]}>
            <styled.img className={styles["asset__image"]} src={a.url} />

            {editing && (
              <HStack className={styles["asset__actions"]} justify="end" p="2">
                <Button
                  type="button"
                  size="sm"
                  kind="destructive"
                  onClick={() => handlers.handleAssetRemove(a)}
                >
                  <TrashIcon />
                </Button>
              </HStack>
            )}
          </Box>
        ))}

        {editing && (
          <Box className={styles["asset"]}>
            <styled.label className={button()} htmlFor="assetwall__file-input">
              Upload
            </styled.label>
            <styled.input
              display="none"
              id="assetwall__file-input"
              type="file"
              onChange={handlers.handleFile}
            />
          </Box>
        )}
      </Box>

      {/* TODO: Lightbox and ability to view more */}
      {/* {hasMore > 0 && <Button type="button">{hasMore} more</Button>} */}
    </Box>
  );
}
