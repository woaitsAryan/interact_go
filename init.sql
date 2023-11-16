CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE messages (
    message_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_id UUID NOT NULL,
    user_id UUID NOT NULL,
    post_id UUID,
    project_id UUID,
    opening_id UUID,
    profile_id UUID,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp,
    read BOOLEAN DEFAULT false
);

CREATE TABLE chats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creating_user_id UUID NOT NULL,
    accepting_user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp,
    last_reset_by_creating_user TIMESTAMP DEFAULT current_timestamp,
    last_reset_by_accepting_user TIMESTAMP DEFAULT current_timestamp,
    blocked_by_creating_user BOOLEAN DEFAULT false,
    blocked_by_accepting_user BOOLEAN DEFAULT false,
    latest_message_id UUID,
    accepted BOOLEAN DEFAULT false,
    last_read_message_by_creating_user_id UUID,
    last_read_message_by_accepting_user_id UUID
);

CREATE TABLE group_chats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(50),
    description TEXT,
    admin_only BOOLEAN DEFAULT false,
    cover_pic TEXT DEFAULT 'default.jpg',
    user_id UUID NOT NULL,
    organization_id UUID,
    project_id UUID,
    created_at TIMESTAMP DEFAULT current_timestamp,
    latest_message_id UUID,
);

CREATE TABLE group_chat_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_id UUID NOT NULL,
    user_id UUID NOT NULL,
    content TEXT NOT NULL,
    post_id UUID,
    project_id UUID,
    opening_id UUID,
    profile_id UUID,
    created_at TIMESTAMP DEFAULT current_timestamp,
)