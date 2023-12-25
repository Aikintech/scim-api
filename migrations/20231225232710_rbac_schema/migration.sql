/*
  Warnings:

  - The primary key for the `podcast_playlist` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to drop the column `id` on the `podcast_playlist` table. All the data in the column will be lost.

*/
-- DropIndex
DROP INDEX "podcast_playlist_playlist_id_podcast_id_idx";

-- AlterTable
ALTER TABLE "podcast_playlist" DROP CONSTRAINT "podcast_playlist_pkey",
DROP COLUMN "id",
ADD CONSTRAINT "podcast_playlist_pkey" PRIMARY KEY ("playlist_id", "podcast_id");

-- CreateTable
CREATE TABLE "roles" (
    "id" VARCHAR(40) NOT NULL,
    "name" VARCHAR(40) NOT NULL,
    "display_name" VARCHAR(40),
    "discription" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "roles_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "permissions" (
    "id" VARCHAR(40) NOT NULL,
    "name" VARCHAR(40) NOT NULL,
    "display_name" VARCHAR(40),
    "module" VARCHAR(40),
    "description" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "permissions_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "permission_role" (
    "permission_id" VARCHAR(40) NOT NULL,
    "role_id" VARCHAR(40) NOT NULL,

    CONSTRAINT "permission_role_pkey" PRIMARY KEY ("permission_id","role_id")
);

-- CreateTable
CREATE TABLE "role_user" (
    "role_id" VARCHAR(40) NOT NULL,
    "user_id" VARCHAR(40) NOT NULL,

    CONSTRAINT "role_user_pkey" PRIMARY KEY ("role_id","user_id")
);

-- CreateTable
CREATE TABLE "permission_user" (
    "permission_id" VARCHAR(40) NOT NULL,
    "user_id" VARCHAR(40) NOT NULL,

    CONSTRAINT "permission_user_pkey" PRIMARY KEY ("permission_id","user_id")
);

-- CreateIndex
CREATE UNIQUE INDEX "roles_name_key" ON "roles"("name");

-- CreateIndex
CREATE UNIQUE INDEX "permissions_name_key" ON "permissions"("name");

-- AddForeignKey
ALTER TABLE "permission_role" ADD CONSTRAINT "permission_role_permission_id_fkey" FOREIGN KEY ("permission_id") REFERENCES "permissions"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "permission_role" ADD CONSTRAINT "permission_role_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "role_user" ADD CONSTRAINT "role_user_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "role_user" ADD CONSTRAINT "role_user_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "permission_user" ADD CONSTRAINT "permission_user_permission_id_fkey" FOREIGN KEY ("permission_id") REFERENCES "permissions"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "permission_user" ADD CONSTRAINT "permission_user_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
