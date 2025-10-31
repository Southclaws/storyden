import { Portal, Presence, UsePresenceProps } from "@ark-ui/react";
import { ArrowRight } from "lucide-react";
import Image from "next/image";
import { useQueryState } from "nuqs";

import { Asset } from "@/api/openapi-schema";
import { css, cx } from "@/styled-system/css";
import { Box, Center, HStack } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";
import { useClickAway } from "@/utils/useClickAway";

import { IconButton } from "../ui/icon-button";
import { ArrowLeftIcon, ArrowRightIcon } from "../ui/icons/Arrow";

type Props = UsePresenceProps & {
  asset: Asset;
  set?: Asset[];
  setIndex?: number;
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
  padding: "4",
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

export function AssetLightbox({
  asset,
  set,
  setIndex,
  onClose,
  ...presenceProps
}: Props) {
  const url = getAssetURL(asset.path)!;

  const [view, setView] = useQueryState<string | null>("view", {
    defaultValue: null,
    clearOnDefault: true,
    parse: (value) => (value === "" ? null : value),
  });
  const ref = useClickAway<HTMLImageElement>(handleClose);

  function handleClose() {
    onClose();
  }

  function handlePrevious() {
    if (set === undefined || setIndex === undefined) return;

    const nextIndex = setIndex - 1 < 0 ? set.length - 1 : setIndex - 1;
    const next = set[nextIndex];

    if (next) {
      setView(next.id);
    }
  }

  function handleNext() {
    if (set === undefined || setIndex === undefined) return;

    const nextIndex = setIndex + 1 === set.length ? 0 : setIndex + 1;
    const next = set[nextIndex];

    if (next) {
      setView(next.id);
    }
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
            <HStack ref={ref}>
              <Box position="fixed" left="8">
                <IconButton variant="subtle" size="sm" onClick={handlePrevious}>
                  <ArrowLeftIcon />
                </IconButton>
              </Box>

              <Box position="fixed" right="8">
                <IconButton variant="subtle" size="sm" onClick={handleNext}>
                  <ArrowRightIcon />
                </IconButton>
              </Box>

              <Image
                className={imageStyles}
                src={url}
                width={2000}
                height={2000}
                alt=""
              />
            </HStack>
          </Center>

          {/* <Image src={url} alt={asset.filename} width="1920" height="1080" /> */}
        </Box>
      </Portal>
    </Presence>
  );
}
