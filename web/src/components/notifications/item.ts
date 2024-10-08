import { DatagraphItem, ProfileReference } from "@/api/openapi-schema";

export type NotificationItem = {
  id: string;
  createdAt: Date;
  title: string;
  description: string;
  url: string;
  isRead: boolean;
  source?: ProfileReference;
  item?: DatagraphItem;
};
