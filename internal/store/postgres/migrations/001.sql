-- Create messages table
CREATE TABLE messages
(
    id           BIGSERIAL PRIMARY KEY,
    to_recipient VARCHAR(20)  NOT NULL,
    content      VARCHAR(500) NOT NULL,
    status       smallint     NOT NULL DEFAULT 0,
    created_at   TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_messages_status_created ON messages (status, created_at);