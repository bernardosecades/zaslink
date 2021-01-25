DROP TABLE IF EXISTS sharesecret.secret;

CREATE TABLE sharesecret.secret (
    id varchar(36) NOT NULL PRIMARY KEY,
    content text NOT NULL,
    custom_pwd bool NOT NULL default 0,
    created_at timestamp NOT NULL
);

INSERT INTO `secret` (`id`, `content`, `created_at`)
VALUES
	('22e04f8a-c18d-4f80-8a34-ebd26122274b','cb98267468c271c1a09bd6d03a919a2af89e9bde934b409258e9e462e2a7b312a9e6cb4d92582155f7a7c48922', '2020-10-19 15:20:44'),
	('fa7617c3-7247-4cc9-9047-c8111440728a','cb98267468c271c1a09bd6d03a919a2af89e9bde934b409258e9e462e2a7b312a9e6cb4d92582155f7a7c48922', '2020-10-19 15:20:44'),
	('7bd3c403-fd16-47fa-ba77-87412dcef1b0','cb98267468c271c1a09bd6d03a919a2af89e9bde934b409258e9e462e2a7b312a9e6cb4d92582155f7a7c48922', '2020-10-19 15:20:44');