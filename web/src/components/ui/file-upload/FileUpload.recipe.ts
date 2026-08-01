import { fileUploadAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const fileUpload = defineSlotRecipe({
  className: "fileUpload",
  slots: fileUploadAnatomy.keys(),
  base: {
    root: {
      display: "flex",
      flexDirection: "column",
      gap: "4",
      width: "100%",
    },
    label: {
      fontWeight: "medium",
      textStyle: "sm",
    },
    dropzone: {
      alignItems: "center",
      background: "background.inset",
      borderColor: "border.default",
      borderRadius: "md",
      borderWidth: "1px",
      color: "text.default",
      display: "flex",
      flexDirection: "column",
      gap: "3",
      justifyContent: "center",
      minHeight: "xs",
      px: "6",
      py: "4",
    },
    item: {
      animation: "fadeIn 0.25s ease-out",
      background: "background.inset",
      borderColor: "border.default",
      borderRadius: "md",
      borderWidth: "1px",
      color: "text.default",
      columnGap: "3",
      display: "grid",
      gridTemplateColumns: "auto 1fr auto",
      gridTemplateAreas: `
        "preview name delete"
        "preview size delete"
        `,
      p: "4",
    },
    itemGroup: {
      display: "flex",
      flexDirection: "column",
      gap: "3",
    },
    itemName: {
      color: "text.default",
      fontWeight: "medium",
      gridArea: "name",
      textStyle: "sm",
    },
    itemSizeText: {
      color: "text.muted",
      gridArea: "size",
      textStyle: "sm",
    },
    itemDeleteTrigger: {
      alignSelf: "flex-start",
      gridArea: "delete",
    },
    itemPreview: {
      gridArea: "preview",
    },
    itemPreviewImage: {
      aspectRatio: "1",
      height: "10",
      objectFit: "scale-down",
      width: "10",
    },
  },
});
