--
-- PostgreSQL database dump
--

-- Dumped from database version 17.5
-- Dumped by pg_dump version 17.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: update_task_day_aggregates(); Type: FUNCTION; Schema: public; Owner: root
--

CREATE FUNCTION public.update_task_day_aggregates() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Handle INSERT operations (a new task is created)
    IF (TG_OP = 'INSERT') THEN
        UPDATE task_days
        SET
            count = count + 1,
            total_duration = total_duration + NEW.duration,
            -- Add to completed_duration only if the new task is already marked complete
            completed_duration = completed_duration + (CASE WHEN NEW.completed THEN NEW.duration ELSE INTERVAL '0' END)
        WHERE id = NEW.task_day_id;
        RETURN NEW;
    END IF;

    -- Handle UPDATE operations (a task is modified)
    IF (TG_OP = 'UPDATE') THEN
        UPDATE task_days
        SET
            -- Adjust total_duration by the difference in the task's duration
            total_duration = total_duration - OLD.duration + NEW.duration,
            -- Adjust completed_duration based on the old and new completion status and duration
            completed_duration = completed_duration
                                 - (CASE WHEN OLD.completed THEN OLD.duration ELSE INTERVAL '0' END) -- Subtract old completed value
                                 + (CASE WHEN NEW.completed THEN NEW.duration ELSE INTERVAL '0' END) -- Add new completed value
        WHERE id = NEW.task_day_id;
        RETURN NEW;
    END IF;

    -- Handle DELETE operations (a task is removed)
    IF (TG_OP = 'DELETE') THEN
        UPDATE task_days
        SET
            count = count - 1,
            total_duration = total_duration - OLD.duration,
            -- Subtract from completed_duration only if the deleted task was complete
            completed_duration = completed_duration - (CASE WHEN OLD.completed THEN OLD.duration ELSE INTERVAL '0' END)
        WHERE id = OLD.task_day_id;
        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$;


ALTER FUNCTION public.update_task_day_aggregates() OWNER TO root;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: challenge_entries; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.challenge_entries (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    challenge_id uuid,
    date date DEFAULT (now())::date,
    completed boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.challenge_entries OWNER TO root;

--
-- Name: challenges; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.challenges (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title character varying NOT NULL,
    user_id uuid,
    description character varying,
    start_date date DEFAULT (now())::date NOT NULL,
    end_date date,
    active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.challenges OWNER TO root;

--
-- Name: current_challenges_view; Type: VIEW; Schema: public; Owner: root
--

CREATE VIEW public.current_challenges_view AS
 SELECT id,
    title,
    user_id,
    description,
    start_date,
    end_date,
    active,
    created_at,
        CASE
            WHEN (end_date IS NOT NULL) THEN ((end_date - start_date) + 1)
            ELSE ((CURRENT_DATE - start_date) + 1)
        END AS duration
   FROM public.challenges c;


ALTER VIEW public.current_challenges_view OWNER TO root;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO root;

--
-- Name: task_days; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.task_days (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    date date DEFAULT (now())::date,
    count integer DEFAULT 0,
    total_duration interval DEFAULT '00:00:00'::interval,
    completed_duration interval DEFAULT '00:00:00'::interval
);


ALTER TABLE public.task_days OWNER TO root;

--
-- Name: tasks; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.tasks (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    task_day_id uuid,
    title character varying NOT NULL,
    description character varying,
    duration interval NOT NULL,
    completed boolean DEFAULT false
);


ALTER TABLE public.tasks OWNER TO root;

--
-- Name: users; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email character varying NOT NULL,
    name character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    timezone character varying DEFAULT 'UTC'::character varying NOT NULL,
    password_changed_at timestamp with time zone DEFAULT '0001-01-01 00:00:00+00'::timestamp with time zone NOT NULL,
    hashed_password character varying NOT NULL
);


ALTER TABLE public.users OWNER TO root;

--
-- Name: challenge_entries challenge_entries_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.challenge_entries
    ADD CONSTRAINT challenge_entries_pkey PRIMARY KEY (id);


--
-- Name: challenges challenges_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.challenges
    ADD CONSTRAINT challenges_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: task_days task_days_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.task_days
    ADD CONSTRAINT task_days_pkey PRIMARY KEY (id);


--
-- Name: tasks tasks_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_pkey PRIMARY KEY (id);


--
-- Name: task_days user_id_date_unique; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.task_days
    ADD CONSTRAINT user_id_date_unique UNIQUE (user_id, date);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: challenge_entries_challenge_id_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX challenge_entries_challenge_id_idx ON public.challenge_entries USING btree (challenge_id);


--
-- Name: challenge_entries_date_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX challenge_entries_date_idx ON public.challenge_entries USING btree (date);


--
-- Name: challenges_start_date_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX challenges_start_date_idx ON public.challenges USING btree (start_date);


--
-- Name: challenges_user_id_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX challenges_user_id_idx ON public.challenges USING btree (user_id);


--
-- Name: task_days_date_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX task_days_date_idx ON public.task_days USING btree (date);


--
-- Name: task_days_user_id_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX task_days_user_id_idx ON public.task_days USING btree (user_id);


--
-- Name: tasks_task_day_id_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX tasks_task_day_id_idx ON public.tasks USING btree (task_day_id);


--
-- Name: tasks tasks_after_change_trigger; Type: TRIGGER; Schema: public; Owner: root
--

CREATE TRIGGER tasks_after_change_trigger AFTER INSERT OR DELETE OR UPDATE ON public.tasks FOR EACH ROW EXECUTE FUNCTION public.update_task_day_aggregates();


--
-- Name: challenge_entries challenge_entries_challenge_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.challenge_entries
    ADD CONSTRAINT challenge_entries_challenge_id_fkey FOREIGN KEY (challenge_id) REFERENCES public.challenges(id) ON DELETE CASCADE;


--
-- Name: challenges challenges_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.challenges
    ADD CONSTRAINT challenges_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: task_days task_days_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.task_days
    ADD CONSTRAINT task_days_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: tasks tasks_task_day_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_task_day_id_fkey FOREIGN KEY (task_day_id) REFERENCES public.task_days(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

