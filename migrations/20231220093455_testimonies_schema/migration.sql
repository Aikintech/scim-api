-- CreateTable
CREATE TABLE "testimonies" (
    "id" VARCHAR(40) NOT NULL,
    "yt_reference_id" VARCHAR(40),
    "tk_reference_id" VARCHAR(40),
    "yt_url" TEXT,
    "tk_url" TEXT,
    "title" TEXT NOT NULL,
    "body" TEXT NOT NULL,
    "published" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "testimonies_pkey" PRIMARY KEY ("id")
);
