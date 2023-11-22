CREATE TABLE IF NOT EXISTS `assets`
(
    `version`  VARCHAR(255) NOT NULL,
    `platform` VARCHAR(255) NOT NULL,
    `major`    INT(11)      NOT NULL,
    `minor`    INT(11)      NOT NULL,
    `patch`    INT(11)      NOT NULL,
    `hash`     VARCHAR(255) NOT NULL,
    PRIMARY KEY (`version`)
);

CREATE INDEX `idx_assets_platform_major_minor_patch` ON `assets` (`platform`, `major`, `minor`, `patch`);

CREATE TABLE IF NOT EXISTS `assets_urls`
(
    `id`  INT          NOT NULL AUTO_INCREMENT,
    `url` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
);


CREATE TABLE IF NOT EXISTS `definitions`
(
    `version`  VARCHAR(255) NOT NULL,
    `platform` VARCHAR(255) NOT NULL,
    `major`    INT(11)      NOT NULL,
    `minor`    INT(11)      NOT NULL,
    `patch`    INT(11)      NOT NULL,
    `hash`     VARCHAR(255) NOT NULL,
    PRIMARY KEY (`version`)
);

CREATE INDEX `idx_definitions_platform_major_minor_patch` ON `definitions` (`platform`, `major`, `minor`, `patch`);

CREATE TABLE IF NOT EXISTS `definitions_urls`
(
    `id`  INT          NOT NULL AUTO_INCREMENT,
    `url` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
);
