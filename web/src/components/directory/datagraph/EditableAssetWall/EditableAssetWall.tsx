import styles from "./assetwall.module.css";

import { Box, HstackProps, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

import { Props, useEditableAssetWall } from "./useEditableAssetWall";

export function EditableAssetWall({
  assets,
  editing,
  onUpload,
  ...props
}: Props & HstackProps) {
  const { handlers } = useEditableAssetWall({
    assets,
    editing,
    onUpload,
  });

  return (
    <Box className={styles["root"]} {...props} gap="2" w="full">
      {assets?.map((a) => (
        <styled.img key={a.id} className={styles["asset"]} src={a.url} />
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

      {/* TODO: Lightbox and ability to view more */}
      {/* {hasMore > 0 && <Button type="button">{hasMore} more</Button>} */}
    </Box>
  );
}
