CREATE TABLE `devspore_offset_table`
(
    `id`        bigint(11) NOT NULL AUTO_INCREMENT,
    `group_id`  varchar(255),
    `topic`     varchar(255),
    `partition` int(255),
    `offset`    bigint(255),
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_group_id_topic_partition` (`group_id`,`topic`,`partition`)
) ENGINE=InnoDB AUTO_INCREMENT=1;