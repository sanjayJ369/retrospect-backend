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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: challenge_entries; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.challenge_entries (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    challenge_id uuid,
    user_id uuid,
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
    duration integer,
    active boolean,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.challenges OWNER TO root;

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
    total_duration time without time zone
);


ALTER TABLE public.task_days OWNER TO root;

--
-- Name: tasks; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.tasks (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    task_day_id uuid,
    user_id uuid,
    title character varying NOT NULL,
    description character varying,
    duration time without time zone NOT NULL,
    completed boolean
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
    updated_at timestamp without time zone
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
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: challenge_entries_date_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX challenge_entries_date_idx ON public.challenge_entries USING btree (date);


--
-- Name: challenges_start_date_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX challenges_start_date_idx ON public.challenges USING btree (start_date);


--
-- Name: task_days_date_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX task_days_date_idx ON public.task_days USING btree (date);


--
-- Name: users_id_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX users_id_idx ON public.users USING btree (id);


--
-- Name: challenge_entries challenge_entries_challenge_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.challenge_entries
    ADD CONSTRAINT challenge_entries_challenge_id_fkey FOREIGN KEY (challenge_id) REFERENCES public.challenges(id);


--
-- Name: challenge_entries challenge_entries_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.challenge_entries
    ADD CONSTRAINT challenge_entries_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: challenges challenges_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.challenges
    ADD CONSTRAINT challenges_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: task_days task_days_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.task_days
    ADD CONSTRAINT task_days_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: tasks tasks_task_day_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_task_day_id_fkey FOREIGN KEY (task_day_id) REFERENCES public.task_days(id);


--
-- Name: tasks tasks_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

