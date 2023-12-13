-- CreateEnum
CREATE TYPE "TransactionTypeEnum" AS ENUM ('tithe', 'pledge', 'offertory', 'freewill', 'busing', 'covenant_partner', 'other');

-- CreateEnum
CREATE TYPE "ReferenceTypeEnum" AS ENUM ('paystack', 'stripe', 'theteller');

-- CreateTable
CREATE TABLE "transactions" (
    "id" VARCHAR(40) NOT NULL,
    "user_id" VARCHAR(40) NOT NULL,
    "reference_id" VARCHAR(40),
    "reference_type" "ReferenceTypeEnum" DEFAULT 'paystack',
    "idempotency_key" VARCHAR(40) NOT NULL,
    "amount" DECIMAL(10,2) NOT NULL,
    "type" "TransactionTypeEnum" NOT NULL DEFAULT 'freewill',
    "description" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "transactions_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "transactions_user_id_type_reference_type_idx" ON "transactions"("user_id", "type", "reference_type");

-- AddForeignKey
ALTER TABLE "transactions" ADD CONSTRAINT "transactions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
