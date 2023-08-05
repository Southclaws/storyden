import { Unready } from "src/components/Unready";

import { Collection } from "./components/Collection";
import { Props, useCollectionScreen } from "./useCollectionScreen";

export function CollectionScreen(props: Props) {
  const { data, error } = useCollectionScreen(props);

  if (!data) return <Unready {...error} />;

  return <Collection {...data} />;
}
