import type { ProfileReference } from "@/api/openapi-schema";

export interface MemberIdentity extends Pick<
  ProfileReference,
  "id" | "handle" | "name"
> {
  roles?: ProfileReference["roles"];
}
