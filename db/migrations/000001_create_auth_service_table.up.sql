create type Role as enum (
    'admin',
    'user'
);

create table users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_name varchar(100),
    role Role DEFAULT 'admin',
    email varchar(100),
    password_hash varchar(250),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);