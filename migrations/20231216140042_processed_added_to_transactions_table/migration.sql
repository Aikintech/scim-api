-- DropIndex
DROP INDEX "transactions_user_id_type_provider_idx";

-- AlterTable
ALTER TABLE "transactions" ADD COLUMN     "processed" BOOLEAN NOT NULL DEFAULT false;

-- CreateIndex
CREATE INDEX "transactions_user_id_type_provider_processed_idx" ON "transactions"("user_id", "type", "provider", "processed");
