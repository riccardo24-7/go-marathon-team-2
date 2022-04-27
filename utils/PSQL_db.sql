CREATE TABLE "metrics" (
        id SERIAL PRIMARY KEY,
        metrics_info REAL
);

CREATE TABLE "folders" (
        id SERIAL PRIMARY KEY,
        folders_uid VARCHAR(100),
        folders_name CHAR(30)
);

CREATE TABLE "metrics_to_folders" 
(
        mtf_id SERIAL PRIMARY KEY,
            metrics_id INTEGER NOT NULL,
            folders_id INTEGER NOT NULL,
        CONSTRAINT "FK_metrics_id" FOREIGN KEY ("metrics_id")
          REFERENCES "metrics" ("id"),
        CONSTRAINT "FK_folders_id" FOREIGN KEY ("folders_id")
          REFERENCES "folders" ("id")
);

CREATE UNIQUE INDEX "UI_metrics_to_folders_metrics_id_folders_id"
  ON "metrics_to_folders"
  USING btree
  ("metrics_id", "folders_id");
  
CREATE TABLE "METRIC_GROUPS" (
        id SERIAL PRIMARY KEY,
        folderName CHAR(30)
);

CREATE TABLE "METRIC_GROUPS_to_folders" 
(
        mgtf_id SERIAL PRIMARY KEY,
            mgroups_id INTEGER NOT NULL,
            folders_id INTEGER NOT NULL,
        CONSTRAINT "FK_mgroups_id" FOREIGN KEY ("mgroups_id")
          REFERENCES "METRIC_GROUPS" ("id"),
        CONSTRAINT "FK_folders_id" FOREIGN KEY ("folders_id")
          REFERENCES "folders" ("id")
);

CREATE UNIQUE INDEX "UI_METRIC_GROUPS_to_folders_mgroups_id_folders_id"
  ON "METRIC_GROUPS_to_folders"
  USING btree
  ("mgroups_id", "folders_id");