-- migrate:up
CREATE TABLE IF NOT EXISTS public.messages (
  id SERIAL PRIMARY KEY,
  text TEXT,
  owner_id int REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS owner_message_idx ON public.messages(owner_id);

CREATE TABLE IF NOT EXISTS public.dialogs (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  last_message_id int REFERENCES messages(id) ON DELETE CASCADE ON UPDATE CASCADE,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE public.messages
ADD COLUMN dialog_id int REFERENCES public.dialogs (id) ON DELETE CASCADE ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS dialog_message_idx ON public.messages(dialog_id);

CREATE TABLE IF NOT EXISTS public.user_dialogs (
  id SERIAL PRIMARY KEY,
  user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  dialog_id int NOT NULL REFERENCES public.dialogs (id) ON DELETE CASCADE ON UPDATE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, dialog_id)
);

CREATE INDEX IF NOT EXISTS user_user_dialogs_idx ON public.user_dialogs(user_id);

CREATE INDEX IF NOT EXISTS dialog_user_dialogs_idx ON public.user_dialogs(dialog_id);

-- migrate:down
DROP INDEX IF EXISTS user_user_dialog_idx;

DROP INDEX IF EXISTS dialog_user_dialogs_idx;

DROP INDEX IF EXISTS owner_user_idx;

DROP INDEX IF EXISTS dialog_message_idx;

DROP TABLE IF EXISTS public.user_dialogs;

DROP TABLE IF EXISTS public.dialogs CASCADE;

DROP TABLE IF EXISTS public.messages CASCADE;