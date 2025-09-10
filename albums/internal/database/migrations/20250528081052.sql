-- Modify "users" table
ALTER TABLE `users` ADD COLUMN `revoked` bool NULL DEFAULT 0;
