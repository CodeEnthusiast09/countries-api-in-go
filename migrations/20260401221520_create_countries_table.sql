-- Create "countries" table
CREATE TABLE `countries` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `capital` varchar(255) NOT NULL DEFAULT '',
  `region` varchar(100) NOT NULL DEFAULT '',
  `population` bigint NOT NULL DEFAULT 0,
  `currency_code` varchar(10) NULL,
  `exchange_rate` double NULL,
  `estimated_gdp` double NULL,
  `flag_url` varchar(500) NOT NULL DEFAULT '',
  `last_refreshed_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_countries_name` (`name`)
) COLLATE utf8mb4_uca1400_ai_ci;
