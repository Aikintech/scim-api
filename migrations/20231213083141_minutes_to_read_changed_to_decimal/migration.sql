/*
  Warnings:

  - You are about to alter the column `minutes_to_read` on the `posts` table. The data in that column could be lost. The data in that column will be cast from `Integer` to `Decimal(3,2)`.

*/
-- AlterTable
ALTER TABLE "posts" ALTER COLUMN "minutes_to_read" SET DEFAULT 0,
ALTER COLUMN "minutes_to_read" SET DATA TYPE DECIMAL(3,2);
