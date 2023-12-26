-- DropForeignKey
ALTER TABLE "comments" DROP CONSTRAINT "comments_user_id_fkey";

-- DropForeignKey
ALTER TABLE "likes" DROP CONSTRAINT "likes_user_id_fkey";

-- DropForeignKey
ALTER TABLE "permission_role" DROP CONSTRAINT "permission_role_permission_id_fkey";

-- DropForeignKey
ALTER TABLE "permission_role" DROP CONSTRAINT "permission_role_role_id_fkey";

-- DropForeignKey
ALTER TABLE "permission_user" DROP CONSTRAINT "permission_user_permission_id_fkey";

-- DropForeignKey
ALTER TABLE "permission_user" DROP CONSTRAINT "permission_user_user_id_fkey";

-- DropForeignKey
ALTER TABLE "playlists" DROP CONSTRAINT "playlists_user_id_fkey";

-- DropForeignKey
ALTER TABLE "podcast_playlist" DROP CONSTRAINT "podcast_playlist_playlist_id_fkey";

-- DropForeignKey
ALTER TABLE "podcast_playlist" DROP CONSTRAINT "podcast_playlist_podcast_id_fkey";

-- DropForeignKey
ALTER TABLE "posts" DROP CONSTRAINT "posts_user_id_fkey";

-- DropForeignKey
ALTER TABLE "prayer_requests" DROP CONSTRAINT "prayer_requests_user_id_fkey";

-- DropForeignKey
ALTER TABLE "role_user" DROP CONSTRAINT "role_user_role_id_fkey";

-- DropForeignKey
ALTER TABLE "role_user" DROP CONSTRAINT "role_user_user_id_fkey";

-- DropForeignKey
ALTER TABLE "tokens" DROP CONSTRAINT "tokens_user_id_fkey";

-- DropForeignKey
ALTER TABLE "user_devices" DROP CONSTRAINT "user_devices_user_id_fkey";

-- DropForeignKey
ALTER TABLE "verification_codes" DROP CONSTRAINT "verification_codes_user_id_fkey";

-- AddForeignKey
ALTER TABLE "playlists" ADD CONSTRAINT "playlists_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "podcast_playlist" ADD CONSTRAINT "podcast_playlist_playlist_id_fkey" FOREIGN KEY ("playlist_id") REFERENCES "playlists"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "podcast_playlist" ADD CONSTRAINT "podcast_playlist_podcast_id_fkey" FOREIGN KEY ("podcast_id") REFERENCES "podcasts"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "prayer_requests" ADD CONSTRAINT "prayer_requests_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "user_devices" ADD CONSTRAINT "user_devices_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "verification_codes" ADD CONSTRAINT "verification_codes_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "comments" ADD CONSTRAINT "comments_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "likes" ADD CONSTRAINT "likes_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "posts" ADD CONSTRAINT "posts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "tokens" ADD CONSTRAINT "tokens_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "permission_role" ADD CONSTRAINT "permission_role_permission_id_fkey" FOREIGN KEY ("permission_id") REFERENCES "permissions"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "permission_role" ADD CONSTRAINT "permission_role_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "role_user" ADD CONSTRAINT "role_user_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "role_user" ADD CONSTRAINT "role_user_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "permission_user" ADD CONSTRAINT "permission_user_permission_id_fkey" FOREIGN KEY ("permission_id") REFERENCES "permissions"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "permission_user" ADD CONSTRAINT "permission_user_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;
