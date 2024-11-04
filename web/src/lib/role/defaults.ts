import { Role } from "@/api/openapi-schema";

const DefaultRoleEveryoneID = "00000000000000000010";
const DefaultRoleAdminID = "00000000000000000020";

export function isDefaultRole(role: Role) {
  switch (role.id) {
    case DefaultRoleAdminID:
    case DefaultRoleEveryoneID:
      return true;

    default:
      return false;
  }
}
