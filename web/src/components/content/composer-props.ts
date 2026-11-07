import { Asset } from "@/api/openapi-schema";

export type ContentComposerProps = {
  className?: string;
  disabled?: boolean;
  focusOnReady?: boolean;
  // Starting content for a new editor. Change the React key to start over.
  initialValue?: string;
  initialValueFormat?: "html" | "markdown";

  // Current content for the read-only renderer.
  value?: string;
  placeholder?: string;
  onChange?: (value: string, isEmpty: boolean) => void;
  onAssetUpload?: (asset: Asset) => void;
};
