CREATE TABLE IF NOT EXISTS public.tasks
(
    id              serial PRIMARY KEY,
    user_id int NOT NULL REFERENCES public.users(id),
    title           text NOT NULL,
    description     text NOT NULL,
    status text NOT NULL,
    date      timestamptz NOT NULL,
    created_date    timestamptz NOT NULL,
    updated_date    timestamptz NOT NULL,
    deleted_date    timestamptz
);
