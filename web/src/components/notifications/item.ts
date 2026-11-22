import type {
  NotificationItemReference,
  NotificationSource,
} from "@/api/openapi-schema";

export type NotificationItem = {
  id: string;
  createdAt: Date;
  title: string;
  description: string;
  body: string;
  url: string;
  isRead: boolean;
  source?: NotificationSource;
  item?: NotificationItemReference;
};
