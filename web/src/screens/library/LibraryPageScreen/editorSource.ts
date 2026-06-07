import { NodeVersion, NodeWithChildren } from "@/api/openapi-schema";

export function liveEditorSourceKey(node: NodeWithChildren) {
  return `live:${node.id}:${node.current_version_id ?? node.updatedAt}`;
}

export function directEditorSourceKey(node: NodeWithChildren) {
  return `direct:${node.id}`;
}

export function versionEditorSourceKey(version: NodeVersion) {
  return `version:${version.id}`;
}
