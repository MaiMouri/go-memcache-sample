
-- CREATE TABLE `users`
-- (
--     id   bigint auto_increment,
--     name varchar(255) NOT NULL,
--     PRIMARY KEY (`id`)
-- );

-- INSERT INTO `users` (`name`)
-- VALUES ('Solomon'),
--        ('Menelik');

CREATE TABLE `members`
(
    id int NOT NULL,
    gender int NOT NULL,
    status int NOT NULL,
    deleted_at datetime,
    created_at datetime,
    updated_at datetime
);

INSERT INTO `members` (`id`, `gender`, `status`,`created_at`)
VALUES ('1', '0', '0', '2022-12-31 00:00:00'),
       ('2', '0', '0', '2022-12-31 00:00:00'),
       ('3', '0', '0', '2022-12-31 00:00:00'),
       ('4', '1', '0', '2022-12-31 00:00:00'),
       ('5', '1', '0', '2022-12-31 00:00:00'),
       ('6', '1', '0', '2022-12-31 00:00:00'),
       ('7', '0', '0', '2023-1-16 00:00:00'),
       ('8', '0', '0', '2023-1-16 00:00:00'),
       ('9', '1', '0', '2023-1-16 00:00:00'),
       ('10', '1', '0', '2023-1-17 00:00:00'),
       ('11', '1', '0', '2023-1-17 00:00:00');
