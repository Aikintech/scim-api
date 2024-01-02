-- CreateTable
CREATE TABLE "user_event" (
    "user_id" VARCHAR(40) NOT NULL,
    "event_id" VARCHAR(40) NOT NULL,

    CONSTRAINT "user_event_pkey" PRIMARY KEY ("user_id","event_id")
);

-- AddForeignKey
ALTER TABLE "user_event" ADD CONSTRAINT "user_event_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "user_event" ADD CONSTRAINT "user_event_event_id_fkey" FOREIGN KEY ("event_id") REFERENCES "events"("id") ON DELETE CASCADE ON UPDATE CASCADE;
