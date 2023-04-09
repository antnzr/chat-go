-- migrate:up
CREATE TABLE IF NOT EXISTS public.messages (
  id SERIAL PRIMARY KEY,
  text TEXT,
  owner_id int REFERENCES public.users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS owner_message_idx ON public.messages(owner_id);

CREATE TABLE IF NOT EXISTS public.chats (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  last_message_id int REFERENCES public.messages(id) ON DELETE CASCADE ON UPDATE CASCADE,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE public.messages
ADD COLUMN chat_id int REFERENCES public.chats (id) ON DELETE CASCADE ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS chat_message_idx ON public.messages(chat_id);

CREATE TABLE IF NOT EXISTS public.user_chats (
  id SERIAL PRIMARY KEY,
  user_id int NOT NULL REFERENCES public.users(id) ON DELETE SET NULL ON UPDATE CASCADE,
  chat_id int NOT NULL REFERENCES public.chats (id) ON DELETE CASCADE ON UPDATE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, chat_id)
);

CREATE INDEX IF NOT EXISTS user_user_chat_idx ON public.user_chats(user_id);

CREATE INDEX IF NOT EXISTS chat_user_chat_idx ON public.user_chats(chat_id);

-- migrate:down
DROP INDEX IF EXISTS user_user_chat_idx;

DROP INDEX IF EXISTS chat_user_chat_idx;

DROP INDEX IF EXISTS owner_user_idx;

DROP INDEX IF EXISTS chat_message_idx;

DROP TABLE IF EXISTS public.user_chats;

DROP TABLE IF EXISTS public.chats CASCADE;

DROP TABLE IF EXISTS public.messages CASCADE;
