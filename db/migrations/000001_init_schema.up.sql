CREATE TABLE
    "people" (
                 "id" UUID NOT NULL DEFAULT (uuid_generate_v4()),
                 "first_name" VARCHAR NOT NULL,
                 "surname" VARCHAR NOT NULL,
                 "email" VARCHAR NOT NULL,
                 "nickname" VARCHAR NOT NULL,
                 "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
                 "updated_at" TIMESTAMP(3) NOT NULL,
                 CONSTRAINT "people_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "people_nickname_key" ON "people"("nickname");