import { Fragment, PropsWithChildren } from "react";

import { Asset } from "src/api/openapi/schemas";
import { FileDrop } from "src/components/content/FileDrop/FileDrop";
import { useFileUpload } from "src/components/content/FileDrop/useFileDrop";

import styles from "./assetwall.module.css";

import { Box, HstackProps, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

type Props = {
  assets: Asset[];
  editing?: boolean;
  onUpload: (asset: Asset) => void;
} & HstackProps;

export function EditableAssetWall({ assets, editing, ...props }: Props) {
  const { upload } = useFileUpload();

  const short = assets.slice(0, 6);
  // const hasMore = assets.length - short.length;

  async function handleFile(event: React.ChangeEvent<HTMLInputElement>) {
    if (!event.target.files) return;

    for (const file of event.target.files) {
      handleAssetUpload(await upload(file));
    }
  }

  const handleAssetUpload = async (asset: Asset) => {
    props.onUpload(asset);
  };

  const Container = editing
    ? ({ children }: PropsWithChildren) => (
        <FileDrop onComplete={handleAssetUpload}>{children}</FileDrop>
      )
    : Box;

  return (
    <Container id="file-drop">
      <Box className={styles["root"]} {...props} gap="2">
        {short.map((a) => (
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
              onChange={handleFile}
            />
          </Box>
        )}

        {/* TODO: Lightbox and ability to view more */}
        {/* {hasMore > 0 && <Button type="button">{hasMore} more</Button>} */}
      </Box>
    </Container>
  );
}
