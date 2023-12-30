/*
  Warnings:

  - You are about to drop the column `method` on the `transactions` table. All the data in the column will be lost.

*/
-- CreateEnum
CREATE TYPE "TransactionChannel" AS ENUM ('card', 'bank', 'mobile_money');

-- AlterEnum
ALTER TYPE "TransactionStatusEnum" ADD VALUE 'canceled';

-- AlterTable
ALTER TABLE "transactions" DROP COLUMN "method",
ADD COLUMN     "channel" "TransactionChannel" NOT NULL DEFAULT 'mobile_money';

-- DropEnum
DROP TYPE "TransactionMethod";
