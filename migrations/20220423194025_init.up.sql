-- create "categories" table
CREATE TABLE "categories" ("id" uuid NOT NULL, "name" character varying NOT NULL, "description" character varying NOT NULL DEFAULT '(No description)', "colour" character varying NOT NULL DEFAULT '#8577ce', "sort" bigint NOT NULL DEFAULT -1, "admin" boolean NOT NULL DEFAULT false, PRIMARY KEY ("id"));
-- create "notifications" table
CREATE TABLE "notifications" ("id" character varying NOT NULL, "title" character varying NOT NULL, "description" character varying NOT NULL, "link" character varying NOT NULL, "read" boolean NOT NULL, "created_at" timestamp(0)with time zone NOT NULL, "subscription_id" character varying NOT NULL, "subscription_notifications" character varying NULL, PRIMARY KEY ("id"));
-- create "posts" table
CREATE TABLE "posts" ("id" uuid NOT NULL, "first" boolean NOT NULL, "title" character varying NULL, "slug" character varying NULL, "pinned" boolean NOT NULL DEFAULT false, "body" character varying NOT NULL, "short" character varying NOT NULL, "created_at" timestamp(0)with time zone NOT NULL, "updated_at" timestamp(0)with time zone NOT NULL, "deleted_at" timestamp(0)with time zone NULL, "category_id" uuid NULL, "root_post_id" uuid NULL, "reply_to_post_id" uuid NULL, "user_posts" uuid NOT NULL, PRIMARY KEY ("id"));
-- create "reacts" table
CREATE TABLE "reacts" ("id" uuid NOT NULL, "emoji" character varying NOT NULL, "created_at" timestamp(0)with time zone NOT NULL, "post_reacts" uuid NULL, "react_user" uuid NULL, "react_post" uuid NULL, "user_reacts" uuid NULL, PRIMARY KEY ("id"));
-- create "subscriptions" table
CREATE TABLE "subscriptions" ("id" character varying NOT NULL, "refers_type" character varying NOT NULL, "refers_to" character varying NOT NULL, "created_at" timestamp(0)with time zone NOT NULL, "updated_at" timestamp(0)with time zone NOT NULL, "deleted_at" timestamp(0)with time zone NULL, "user_id" character varying NOT NULL, "notification_subscription" character varying NULL, "user_subscriptions" uuid NULL, PRIMARY KEY ("id"));
-- create "tags" table
CREATE TABLE "tags" ("id" uuid NOT NULL, "name" character varying NOT NULL, PRIMARY KEY ("id"));
-- Create index "tags_name_key" to table: "tags"
CREATE UNIQUE INDEX "tags_name_key" ON "tags" ("name");
-- create "users" table
CREATE TABLE "users" ("id" uuid NOT NULL, "email" character varying NOT NULL, "name" character varying NOT NULL, "bio" character varying NULL, "admin" boolean NOT NULL DEFAULT false, "created_at" timestamp(0)with time zone NOT NULL, "updated_at" timestamp(0)with time zone NOT NULL, "deleted_at" timestamp(0)with time zone NULL, "subscription_user" character varying NULL, PRIMARY KEY ("id"));
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
-- create "tag_posts" table
CREATE TABLE "tag_posts" ("tag_id" uuid NOT NULL, "post_id" uuid NOT NULL, PRIMARY KEY ("tag_id", "post_id"));
-- Modify "notifications" table
ALTER TABLE "notifications" ADD CONSTRAINT "notifications_subscriptions_notifications" FOREIGN KEY ("subscription_notifications") REFERENCES "subscriptions" ("id") ON DELETE SET NULL;
-- Modify "posts" table
ALTER TABLE "posts" ADD CONSTRAINT "posts_categories_posts" FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON DELETE SET NULL, ADD CONSTRAINT "posts_users_posts" FOREIGN KEY ("user_posts") REFERENCES "users" ("id") ON DELETE NO ACTION;
-- Modify "reacts" table
ALTER TABLE "reacts" ADD CONSTRAINT "reacts_posts_reacts" FOREIGN KEY ("post_reacts") REFERENCES "posts" ("id") ON DELETE SET NULL, ADD CONSTRAINT "reacts_users_user" FOREIGN KEY ("react_user") REFERENCES "users" ("id") ON DELETE SET NULL, ADD CONSTRAINT "reacts_posts_Post" FOREIGN KEY ("react_post") REFERENCES "posts" ("id") ON DELETE SET NULL, ADD CONSTRAINT "reacts_users_reacts" FOREIGN KEY ("user_reacts") REFERENCES "users" ("id") ON DELETE SET NULL;
-- Modify "subscriptions" table
ALTER TABLE "subscriptions" ADD CONSTRAINT "subscriptions_notifications_subscription" FOREIGN KEY ("notification_subscription") REFERENCES "notifications" ("id") ON DELETE SET NULL, ADD CONSTRAINT "subscriptions_users_subscriptions" FOREIGN KEY ("user_subscriptions") REFERENCES "users" ("id") ON DELETE SET NULL;
-- Modify "users" table
ALTER TABLE "users" ADD CONSTRAINT "users_subscriptions_user" FOREIGN KEY ("subscription_user") REFERENCES "subscriptions" ("id") ON DELETE SET NULL;
-- Modify "tag_posts" table
ALTER TABLE "tag_posts" ADD CONSTRAINT "tag_posts_tag_id" FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON DELETE CASCADE, ADD CONSTRAINT "tag_posts_post_id" FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE;
