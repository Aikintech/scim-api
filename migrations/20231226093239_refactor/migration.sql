/*
  Warnings:

  - You are about to drop the column `discription` on the `roles` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "roles" DROP COLUMN "discription",
ADD COLUMN     "description" TEXT;
