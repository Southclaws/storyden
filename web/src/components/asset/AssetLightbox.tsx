import { Portal, Presence, UsePresenceProps } from "@ark-ui/react";
import { useClickAway } from "@uidotdev/usehooks";
import Image from "next/image";

import { Asset } from "@/api/openapi-schema";
import { css, cx } from "@/styled-system/css";
import { Box, Center } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

type Props = UsePresenceProps & {
  asset: Asset;
  onClose: () => void;
};

const overlayStyles = css({
  position: "fixed",
  top: "0",
  left: "0",
  width: "screen",
  height: "dvh",
});

const containerStyles = css({
  zIndex: "overlay",
});

const backdropStyles = css({
  backdropBlur: "sm",
  backdropGrayscale: "0.8",
  backdropBrightness: "0.1",
  backdropFilter: "auto",
});

const backdropImageStyles = css({
  blur: "3xl",
  opacity: "2",
  brightness: "0.5",
  filter: "auto",
});

const lightboxStyles = css({
  margin: "auto",
});

const imageStyles = css({
  height: "auto",
  width: "auto",
  maxHeight: "lvh",
  borderRadius: "sm",
  boxShadow: "var(--box-shadow)",
  "--box-shadow":
    "0px 0px 40px var(--colors-black-a2), 0px 0px 1px var(--colors-gray-a7)",
});

export function AssetLightbox({ asset, onClose, ...presenceProps }: Props) {
  const url = getAssetURL(asset.path)!;

  const ref = useClickAway<HTMLImageElement>(handleClose);

  function handleClose() {
    onClose();
  }

  // NOTE: Presence doesn't work for some reason. Possibly bug in Ark UI.
  if (!presenceProps.present) {
    return null;
  }

  return (
    <Presence lazyMount present={presenceProps.present}>
      <Portal>
        <Box className={cx(overlayStyles, containerStyles)}>
          <Box className={cx(overlayStyles, backdropStyles)}>
            <Image className={backdropImageStyles} src={url} fill alt="" />
          </Box>

          <Center className={cx(overlayStyles, lightboxStyles)}>
            <Image
              ref={ref}
              className={imageStyles}
              src={url}
              width={2000}
              height={2000}
              alt=""
            />
          </Center>

          {/* <Image src={url} alt={asset.filename} width="1920" height="1080" /> */}
        </Box>
      </Portal>
    </Presence>
  );
}
