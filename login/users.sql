CREATE TABLE `users` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL,
  `email` text NOT NULL,
  `salt` text NOT NULL,
  `salted` text NOT NULL,
  `created` text NOT NULL DEFAULT CURRENT_DATE,
  `updated` text DEFAULT NULL
);
