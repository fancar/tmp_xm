create table "user" (
	id bigserial primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null,
	username character varying (100) not null,
	password_hash character varying (200) not null,
	session_ttl bigint not null,
	is_admin boolean not null
);

create unique index idx_user_username on "user"(username);
create index idx_user_username_prefix on "user"(username varchar_pattern_ops);

-- global admin (password: admin)
insert into "user" (
	created_at,
	updated_at,
	username,
	password_hash,
	session_ttl,
	is_admin
) values (
	now(),
	now(),
	'admin',
	'PBKDF2$sha512$1$l8zGKtxRESq3PA2kFhHRWA==$H3lGMxOt55wjwoc+myeOoABofJY9oDpldJa7fhqdjbh700V6FLPML75UmBOt9J5VFNjAL1AvqCozA1HJM0QVGA==',
	0,
	true
);


create table company (
	id uuid primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null,
	name character varying (15) UNIQUE not null,
	description character varying (3000) not null,
	employees_cnt bigserial not null,
	registered boolean not null,
	type bigserial not null
);

create index idx_company_name on company(name);

-- for LIKE searches we can use gin indexes
-- create index idx_company_name_trgm on company using gin (name gin_trgm_ops);	
