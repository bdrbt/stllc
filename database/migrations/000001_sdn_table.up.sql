-- sdn records table
CREATE TABLE sdn_records
(
    id bigint NOT NULL,
    first_name varchar(255),
    last_name varchar(255),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

