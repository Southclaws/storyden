-- reverse: Modify "tag_posts" table
ALTER TABLE "tag_posts" DROP CONSTRAINT "tag_posts_tag_id", DROP CONSTRAINT "tag_posts_post_id";
-- reverse: Modify "users" table
ALTER TABLE "users" DROP CONSTRAINT "users_subscriptions_user";
-- reverse: Modify "subscriptions" table
ALTER TABLE "subscriptions" DROP CONSTRAINT "subscriptions_notifications_subscription", DROP CONSTRAINT "subscriptions_users_subscriptions";
-- reverse: Modify "reacts" table
ALTER TABLE "reacts" DROP CONSTRAINT "reacts_posts_reacts", DROP CONSTRAINT "reacts_users_user", DROP CONSTRAINT "reacts_posts_Post", DROP CONSTRAINT "reacts_users_reacts";
-- reverse: Modify "posts" table
ALTER TABLE "posts" DROP CONSTRAINT "posts_categories_posts", DROP CONSTRAINT "posts_users_posts";
-- reverse: Modify "notifications" table
ALTER TABLE "notifications" DROP CONSTRAINT "notifications_subscriptions_notifications";
-- reverse: create "tag_posts" table
DROP TABLE "tag_posts";
-- reverse: Create index "users_email_key" to table: "users"
DROP INDEX "users_email_key";
-- reverse: create "users" table
DROP TABLE "users";
-- reverse: Create index "tags_name_key" to table: "tags"
DROP INDEX "tags_name_key";
-- reverse: create "tags" table
DROP TABLE "tags";
-- reverse: create "subscriptions" table
DROP TABLE "subscriptions";
-- reverse: create "reacts" table
DROP TABLE "reacts";
-- reverse: create "posts" table
DROP TABLE "posts";
-- reverse: create "notifications" table
DROP TABLE "notifications";
-- reverse: create "categories" table
DROP TABLE "categories";
