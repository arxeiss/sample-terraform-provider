-- Adminer 4.8.1 SQLite 3 3.35.5 dump

DROP TABLE IF EXISTS "networks";
CREATE TABLE "networks" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" text(50) NOT NULL,
  "display_name" text(255) NULL,
  "ip_range" text(18) NOT NULL,
  "use_dhcp" integer(1) NOT NULL
);

CREATE UNIQUE INDEX "network_name" ON "networks" ("name");

INSERT INTO "networks" ("id", "name", "display_name", "ip_range", "use_dhcp") VALUES (1,	'super-duper-network',	'The best network ever',	'192.168.0.0/16',	1);

DROP TABLE IF EXISTS "storages";
CREATE TABLE "storages" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" text(50) NOT NULL,
  "display_name" text(255) NULL,
  "size" integer NOT NULL,
  "network_id" integer NULL,
  "network_ip" text(15) NULL,
  "virtual_machine_id" integer NULL,
  "mount_path" text(100) NULL,
  FOREIGN KEY ("network_id") REFERENCES "networks" ("id"),
  FOREIGN KEY ("virtual_machine_id") REFERENCES "virtual_machines" ("id")
);

CREATE UNIQUE INDEX "storages_name" ON "storages" ("name");

INSERT INTO "storages" ("id", "name", "display_name", "size", "network_id", "network_ip", "virtual_machine_id", "mount_path") VALUES (1,	'super-duper-storage',	'The best storage in universe',	1048576,	1,	'192.168.200.200',	NULL,	NULL);

DROP TABLE IF EXISTS "virtual_machines";
CREATE TABLE "virtual_machines" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" text(50) NOT NULL,
  "display_name" text(255) NULL,
  "ram_size" integer NOT NULL,
  "network_id" integer NULL,
  "network_ip" text(18) NULL,
  "public_ip" text(18) NULL,
  FOREIGN KEY ("network_id") REFERENCES "networks" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE UNIQUE INDEX "virtual_machines_name" ON "virtual_machines" ("name");

INSERT INTO "virtual_machines" ("id", "name", "display_name", "ram_size", "network_id", "network_ip", "public_ip") VALUES (1,	'super-duper-vm',	'The best VM',	4096,	1,	NULL,	'123.123.123.123');

--
