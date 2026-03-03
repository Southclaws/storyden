import * as Alert from "@/components/ui/alert";
import { MotdMetadata, ParsedMotd } from "@/lib/settings/settings";

type Props = {
  motd: ParsedMotd | undefined;
};

export function MotdBanner({ motd }: Props) {
  if (!motd?.content || !isActive(motd.start_at, motd.end_at)) {
    return null;
  }

  return (
    <Alert.Root
      mb="4"
      colorPalette={getAlertPalette(motd.metadata)}
      backgroundColor="colorPalette.4"
      borderColor="colorPalette.6"
    >
      <Alert.Content>
        <Alert.Description
          dangerouslySetInnerHTML={{
            __html: motd.content,
          }}
        />
      </Alert.Content>
    </Alert.Root>
  );
}

function isActive(startAt?: string, endAt?: string) {
  const now = Date.now();

  if (startAt) {
    const start = new Date(startAt).getTime();
    if (!Number.isNaN(start) && now < start) {
      return false;
    }
  }

  if (endAt) {
    const end = new Date(endAt).getTime();
    if (!Number.isNaN(end) && now > end) {
      return false;
    }
  }

  return true;
}

function getAlertPalette(metadata?: MotdMetadata) {
  const type = metadata?.type ?? "information";

  switch (type) {
    case "celebration":
      return "green";
    case "alert":
      return "red";
    default:
      return "blue";
  }
}
