import {
  MenuSelectionDetails,
  Portal,
  SliderValueChangeDetails,
} from "@ark-ui/react";

import { AssetUploadAction } from "@/components/asset/AssetUploadAction";
import { LayoutIcon } from "@/components/ui/icons/Layout";
import { MediaAddIcon } from "@/components/ui/icons/Media";
import { SizeIcon } from "@/components/ui/icons/Size";
import * as Menu from "@/components/ui/menu";
import { Slider } from "@/components/ui/slider";
import { LibraryPageBlockTypeAssetsLayout } from "@/lib/library/metadata";
import { HStack } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";
import { useBlock } from "../useBlock";

import { useLibraryPageAssetsBlock } from "./useLibraryPageAssetsBlock";

export function LibraryPageAssetsBlockMenuItems() {
  const { handleUpload } = useLibraryPageAssetsBlock();

  return (
    <>
      <LayoutMenu />

      <GridSizeControl />

      <Menu.Item value="add-media">
        <AssetUploadAction
          title="Upload media"
          operation="add"
          onFinish={handleUpload}
          hideLabel
          width="full"
        >
          <HStack w="full" gap="1">
            <MediaAddIcon />
            <span>Add media</span>
          </HStack>
        </AssetUploadAction>
      </Menu.Item>
    </>
  );
}

function LayoutMenu() {
  const { store } = useLibraryPageContext();
  const block = useBlock("assets");
  if (block === undefined) {
    throw new Error("LayoutMenu rendered in a page without an Assets block.");
  }

  const { overwriteBlock } = store.getState();

  const handleSelect = ({ value }: MenuSelectionDetails) => {
    const newBlockConfig = {
      layout: value as LibraryPageBlockTypeAssetsLayout,
    };

    overwriteBlock({
      type: "assets",
      config: newBlockConfig,
    });
  };

  return (
    <Menu.Root lazyMount onSelect={handleSelect}>
      <Menu.Trigger asChild>
        <Menu.Item value="add">
          <LayoutIcon />
          &nbsp;Layout
        </Menu.Item>
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            <Menu.Item value="strip">Strip</Menu.Item>
            <Menu.Item value="grid">Grid</Menu.Item>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}

function GridSizeControl() {
  const { handleChangeSize, config } = useLibraryPageAssetsBlock();

  const defaultValue = config?.gridSize ?? 3;

  function handleChange({ value }: SliderValueChangeDetails) {
    const size = value[0];
    if (size === undefined) {
      return;
    }

    handleChangeSize(size);
  }

  return (
    <Menu.Item value="size">
      <HStack w="full" gap="1">
        <SizeIcon flexShrink="0" w="4" h="4" />
        Size
        <Slider
          size="sm"
          minWidth="0"
          min={1}
          max={4}
          defaultValue={[defaultValue]}
          onValueChange={handleChange}
        />
      </HStack>
    </Menu.Item>
  );
}
