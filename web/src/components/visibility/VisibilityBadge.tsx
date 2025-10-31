import { Visibility } from "@/api/openapi-schema";
import { Badge, BadgeProps } from "@/components/ui/badge";
import { visibilityColour } from "@/lib/library/visibilityColours";

type Props = {
  visibility: Visibility;
  size?: BadgeProps["size"];
};

export function VisibilityBadge({ visibility, size = "sm" }: Props) {
  const colorPalette = visibilityColour(visibility);

  return (
    <Badge
      size={size}
      colorPalette={colorPalette}
      backgroundColor="colorPalette.bg"
      borderColor="colorPalette.border"
      color="colorPalette.fg"
      textTransform="capitalize"
    >
      {visibility}
    </Badge>
  );
}
