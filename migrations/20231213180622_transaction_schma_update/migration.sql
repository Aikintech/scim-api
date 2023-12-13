/*
  Warnings:

  - You are about to drop the column `reference_type` on the `transactions` table. All the data in the column will be lost.

*/
-- CreateEnum
CREATE TYPE "TransactionProviderEnum" AS ENUM ('paystack', 'stripe', 'theteller');

-- CreateEnum
CREATE TYPE "TransactionStatusEnum" AS ENUM ('pending', 'success', 'failed');

-- CreateEnum
CREATE TYPE "TransactionMethod" AS ENUM ('card', 'bank', 'mobile_money');

-- CreateEnum
CREATE TYPE "CurrencyEnum" AS ENUM ('GHS', 'USD', 'EUR', 'GBP');

-- DropIndex
DROP INDEX "transactions_user_id_type_reference_type_idx";

-- AlterTable
ALTER TABLE "transactions" DROP COLUMN "reference_type",
ADD COLUMN     "currency" "CurrencyEnum" NOT NULL DEFAULT 'GHS',
ADD COLUMN     "method" "TransactionMethod" NOT NULL DEFAULT 'mobile_money',
ADD COLUMN     "provider" "TransactionProviderEnum" DEFAULT 'paystack',
ADD COLUMN     "status" "TransactionStatusEnum" NOT NULL DEFAULT 'pending';

-- DropEnum
DROP TYPE "ReferenceTypeEnum";

-- CreateIndex
CREATE INDEX "transactions_user_id_type_provider_idx" ON "transactions"("user_id", "type", "provider");
