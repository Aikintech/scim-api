/*
  Warnings:

  - You are about to alter the column `minutes_to_read` on the `posts` table. The data in that column could be lost. The data in that column will be cast from `Decimal(3,2)` to `SmallInt`.

*/
-- AlterTable
ALTER TABLE "posts" ALTER COLUMN "minutes_to_read" SET DEFAULT 0,
ALTER COLUMN "minutes_to_read" SET DATA TYPE SMALLINT;
