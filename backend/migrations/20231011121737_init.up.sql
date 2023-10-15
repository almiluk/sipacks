CREATE TABLE
    "author" (
        "id" serial NOT NULL,
        "nickname" varchar(20) NOT NULL,
        CONSTRAINT "author_pk" PRIMARY KEY ("id")
    )
WITH (OIDS = FALSE);

CREATE TABLE
    "tag" (
        "id" serial NOT NULL,
        "name" varchar(10) NOT NULL,
        CONSTRAINT "tag_pk" PRIMARY KEY ("id")
    )
WITH (OIDS = FALSE);

CREATE TABLE
    "pack" (
        "id" serial NOT NULL,
        "name" varchar(20) NOT NULL,
        "author_id" integer NOT NULL,
        "creation_date" DATE NOT NULL,
        "downloads_num" integer NOT NULL,
        "file_size" integer NOT NULL,
        CONSTRAINT "pack_pk" PRIMARY KEY ("id")
    )
WITH (OIDS = FALSE);

CREATE TABLE
    "pack_tag" (
        "pack_id" integer NOT NULL,
        "tag_id" integer NOT NULL
    )
WITH (OIDS = FALSE);

ALTER TABLE "pack"
ADD
    CONSTRAINT "pack_fk0" FOREIGN KEY ("author_id") REFERENCES "author"("id") ON DELETE CASCADE;

ALTER TABLE "pack_tag"
ADD
    CONSTRAINT "pack_tag_fk0" FOREIGN KEY ("pack_id") REFERENCES "pack"("id") ON DELETE CASCADE;

ALTER TABLE "pack_tag"
ADD
    CONSTRAINT "pack_tag_fk1" FOREIGN KEY ("tag_id") REFERENCES "tag"("id") ON DELETE CASCADE;