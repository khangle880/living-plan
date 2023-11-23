-- +goose Up
CREATE TYPE status_enum AS ENUM ('Todo', 'Process', 'Done');
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    username TEXT UNIQUE NOT NULL,
    fullname TEXT NOT NULL,
    hashed_password TEXT NOT NULL,
    bio TEXT,
    avatar TEXT,
    email TEXT UNIQUE NOT NULL,
    UNIQUE(username, email)
);
ALTER TABLE users
ADD COLUMN api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (
        encode(sha256(random()::text::bytea), 'hex')
    );
CREATE TABLE projects(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    creator UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    purpose TEXT,
    description TEXT,
    background TEXT
);
CREATE TABLE tasks(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    datetime_exc TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    purpose TEXT,
    description TEXT,
    images TEXT [],
    urls TEXT [],
    status status_enum NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES tasks(id) ON DELETE CASCADE
);
CREATE TABLE follow_projects(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    UNIQUE(user_id, project_id)
);
CREATE TABLE habits(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    creator UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    purpose TEXT,
    description TEXT,
    icon TEXT,
    background TEXT,
    images TEXT [],
    urls TEXT [],
    time_in_day TIME NOT NULL,
    loop_week INT [],
    loop_month INT [],
    recommend_duration INT NOT NULL
);
CREATE TABLE follow_habits(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    habit_id UUID NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    promise_from TIMESTAMP NOT NULL,
    promise_end TIMESTAMP NOT NULL,
    processes TIMESTAMP [],
    included_report BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE(user_id, habit_id)
);
CREATE TABLE diaries(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    datetime_exc TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    images TEXT [],
    urls TEXT [],
    icon TEXT,
    background TEXT
);
-- +goose Down
DROP TYPE status_enum;
DROP TABLE users;
DROP TABLE tasks;
DROP TABLE projects;
DROP TABLE follow_projects;
DROP TABLE habits;
DROP TABLE follow_habits;
DROP TABLE diaries;