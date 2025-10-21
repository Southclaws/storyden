import { Plugin } from "@/api/openapi-schema";

import { PluginLogViewer } from "../../PluginLogViewer";

type Props = {
  plugin: Plugin;
};

export function LogsTab({ plugin }: Props) {
  return <PluginLogViewer plugin={plugin} />;
}
