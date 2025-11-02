--
-- PostgreSQL database dump
--

-- Dumped from database version 17.6 (Debian 17.6-1.pgdg12+1)
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
-- Name: public; Type: SCHEMA; Schema: -; Owner: paydeya_user
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO paydeya_user;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: material_blocks; Type: TABLE; Schema: public; Owner: paydeya_user
--

CREATE TABLE public.material_blocks (
    id integer NOT NULL,
    material_id integer NOT NULL,
    block_id character varying(50) NOT NULL,
    type character varying(20) NOT NULL,
    content jsonb NOT NULL,
    styles jsonb,
    animation jsonb,
    "position" integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT material_blocks_type_check CHECK (((type)::text = ANY ((ARRAY['text'::character varying, 'image'::character varying, 'video'::character varying, 'formula'::character varying, 'quiz'::character varying])::text[])))
);


ALTER TABLE public.material_blocks OWNER TO paydeya_user;

--
-- Name: material_blocks_id_seq; Type: SEQUENCE; Schema: public; Owner: paydeya_user
--

CREATE SEQUENCE public.material_blocks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.material_blocks_id_seq OWNER TO paydeya_user;

--
-- Name: material_blocks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: paydeya_user
--

ALTER SEQUENCE public.material_blocks_id_seq OWNED BY public.material_blocks.id;


--
-- Name: material_ratings; Type: TABLE; Schema: public; Owner: paydeya_user
--

CREATE TABLE public.material_ratings (
    id integer NOT NULL,
    material_id integer NOT NULL,
    user_id integer NOT NULL,
    rating integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT material_ratings_rating_check CHECK (((rating >= 1) AND (rating <= 5)))
);


ALTER TABLE public.material_ratings OWNER TO paydeya_user;

--
-- Name: material_ratings_id_seq; Type: SEQUENCE; Schema: public; Owner: paydeya_user
--

CREATE SEQUENCE public.material_ratings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.material_ratings_id_seq OWNER TO paydeya_user;

--
-- Name: material_ratings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: paydeya_user
--

ALTER SEQUENCE public.material_ratings_id_seq OWNED BY public.material_ratings.id;


--
-- Name: materials; Type: TABLE; Schema: public; Owner: paydeya_user
--

CREATE TABLE public.materials (
    id integer NOT NULL,
    title character varying(500) NOT NULL,
    subject character varying(100) NOT NULL,
    author_id integer NOT NULL,
    status character varying(20) DEFAULT 'draft'::character varying NOT NULL,
    access character varying(20) DEFAULT 'open'::character varying NOT NULL,
    share_url character varying(500),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT materials_access_check CHECK (((access)::text = ANY ((ARRAY['open'::character varying, 'link'::character varying])::text[]))),
    CONSTRAINT materials_status_check CHECK (((status)::text = ANY ((ARRAY['draft'::character varying, 'published'::character varying, 'archived'::character varying])::text[])))
);


ALTER TABLE public.materials OWNER TO paydeya_user;

--
-- Name: materials_id_seq; Type: SEQUENCE; Schema: public; Owner: paydeya_user
--

CREATE SEQUENCE public.materials_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.materials_id_seq OWNER TO paydeya_user;

--
-- Name: materials_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: paydeya_user
--

ALTER SEQUENCE public.materials_id_seq OWNED BY public.materials.id;


--
-- Name: specializations; Type: TABLE; Schema: public; Owner: paydeya_user
--

CREATE TABLE public.specializations (
    id integer NOT NULL,
    user_id integer,
    subject character varying(100) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.specializations OWNER TO paydeya_user;

--
-- Name: specializations_id_seq; Type: SEQUENCE; Schema: public; Owner: paydeya_user
--

CREATE SEQUENCE public.specializations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.specializations_id_seq OWNER TO paydeya_user;

--
-- Name: specializations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: paydeya_user
--

ALTER SEQUENCE public.specializations_id_seq OWNED BY public.specializations.id;


--
-- Name: subjects; Type: TABLE; Schema: public; Owner: paydeya_user
--

CREATE TABLE public.subjects (
    id character varying(50) NOT NULL,
    name character varying(100) NOT NULL,
    icon character varying(200)
);


ALTER TABLE public.subjects OWNER TO paydeya_user;

--
-- Name: teacher_specializations; Type: TABLE; Schema: public; Owner: paydeya_user
--

CREATE TABLE public.teacher_specializations (
    id integer NOT NULL,
    user_id integer NOT NULL,
    subject character varying(100) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.teacher_specializations OWNER TO paydeya_user;

--
-- Name: teacher_specializations_id_seq; Type: SEQUENCE; Schema: public; Owner: paydeya_user
--

CREATE SEQUENCE public.teacher_specializations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.teacher_specializations_id_seq OWNER TO paydeya_user;

--
-- Name: teacher_specializations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: paydeya_user
--

ALTER SEQUENCE public.teacher_specializations_id_seq OWNED BY public.teacher_specializations.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: paydeya_user
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    full_name character varying(255) NOT NULL,
    role character varying(20) NOT NULL,
    avatar_url character varying(500),
    is_verified boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_role_check CHECK (((role)::text = ANY ((ARRAY['student'::character varying, 'teacher'::character varying, 'admin'::character varying])::text[])))
);


ALTER TABLE public.users OWNER TO paydeya_user;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: paydeya_user
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO paydeya_user;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: paydeya_user
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: material_blocks id; Type: DEFAULT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.material_blocks ALTER COLUMN id SET DEFAULT nextval('public.material_blocks_id_seq'::regclass);


--
-- Name: material_ratings id; Type: DEFAULT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.material_ratings ALTER COLUMN id SET DEFAULT nextval('public.material_ratings_id_seq'::regclass);


--
-- Name: materials id; Type: DEFAULT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.materials ALTER COLUMN id SET DEFAULT nextval('public.materials_id_seq'::regclass);


--
-- Name: specializations id; Type: DEFAULT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.specializations ALTER COLUMN id SET DEFAULT nextval('public.specializations_id_seq'::regclass);


--
-- Name: teacher_specializations id; Type: DEFAULT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.teacher_specializations ALTER COLUMN id SET DEFAULT nextval('public.teacher_specializations_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: material_blocks; Type: TABLE DATA; Schema: public; Owner: paydeya_user
--

COPY public.material_blocks (id, material_id, block_id, type, content, styles, animation, "position", created_at) FROM stdin;
\.


--
-- Data for Name: material_ratings; Type: TABLE DATA; Schema: public; Owner: paydeya_user
--

COPY public.material_ratings (id, material_id, user_id, rating, created_at) FROM stdin;
\.


--
-- Data for Name: materials; Type: TABLE DATA; Schema: public; Owner: paydeya_user
--

COPY public.materials (id, title, subject, author_id, status, access, share_url, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: specializations; Type: TABLE DATA; Schema: public; Owner: paydeya_user
--

COPY public.specializations (id, user_id, subject, created_at) FROM stdin;
\.


--
-- Data for Name: subjects; Type: TABLE DATA; Schema: public; Owner: paydeya_user
--

COPY public.subjects (id, name, icon) FROM stdin;
informatics	Информатика	/icons/informatics.svg
mathematics	Математика	/icons/mathematics.svg
physics	Физика	/icons/physics.svg
programming	Программирование	/icons/programming.svg
\.


--
-- Data for Name: teacher_specializations; Type: TABLE DATA; Schema: public; Owner: paydeya_user
--

COPY public.teacher_specializations (id, user_id, subject, created_at) FROM stdin;
1	2	Информатика	2025-10-24 12:42:22.47979+00
2	2	Математика	2025-10-24 12:42:22.47979+00
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: paydeya_user
--

COPY public.users (id, email, password_hash, full_name, role, avatar_url, is_verified, created_at, updated_at) FROM stdin;
1	student@example.com	hashed_password_1	Иван Петров	student	\N	t	2025-10-24 06:59:47.314045+00	2025-10-24 06:59:47.314045+00
2	teacher@example.com	hashed_password_2	Мария Сидорова	teacher	\N	t	2025-10-24 06:59:47.314045+00	2025-10-24 06:59:47.314045+00
3	admin@example.com	hashed_password_3	Администратор Системы	admin	\N	t	2025-10-24 06:59:47.314045+00	2025-10-24 06:59:47.314045+00
19	test2@example.com	$2a$10$bJaxIgR3TBAulo5g4DW0IeUUzIZyrzFNK/FvY16siFxBmpI70PNpG	Тестовый Пользователь 2	student		t	2025-10-28 07:24:56.546852+00	2025-10-28 07:24:56.546852+00
\.


--
-- Name: material_blocks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: paydeya_user
--

SELECT pg_catalog.setval('public.material_blocks_id_seq', 1, false);


--
-- Name: material_ratings_id_seq; Type: SEQUENCE SET; Schema: public; Owner: paydeya_user
--

SELECT pg_catalog.setval('public.material_ratings_id_seq', 1, false);


--
-- Name: materials_id_seq; Type: SEQUENCE SET; Schema: public; Owner: paydeya_user
--

SELECT pg_catalog.setval('public.materials_id_seq', 1, false);


--
-- Name: specializations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: paydeya_user
--

SELECT pg_catalog.setval('public.specializations_id_seq', 1, false);


--
-- Name: teacher_specializations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: paydeya_user
--

SELECT pg_catalog.setval('public.teacher_specializations_id_seq', 30, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: paydeya_user
--

SELECT pg_catalog.setval('public.users_id_seq', 49, true);


--
-- Name: material_blocks material_blocks_material_id_block_id_key; Type: CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.material_blocks
    ADD CONSTRAINT material_blocks_material_id_block_id_key UNIQUE (material_id, block_id);


--
-- Name: material_blocks material_blocks_pkey; Type: CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.material_blocks
    ADD CONSTRAINT material_blocks_pkey PRIMARY KEY (id);


--
-- Name: material_ratings material_ratings_material_id_user_id_key; Type: CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.material_ratings
    ADD CONSTRAINT material_ratings_material_id_user_id_key UNIQUE (material_id, user_id);


--
-- Name: material_ratings material_ratings_pkey; Type: CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.material_ratings
    ADD CONSTRAINT material_ratings_pkey PRIMARY KEY (id);


--
-- Name: materials materials_pkey; Type: CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.materials
    ADD CONSTRAINT materials_pkey PRIMARY KEY (id);


--
-- Name: specializations specializations_pkey; Type: CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.specializations
    ADD CONSTRAINT specializations_pkey PRIMARY KEY (id);


--
-- Name: subjects subjects_pkey; Type: CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.subjects
    ADD CONSTRAINT subjects_pkey PRIMARY KEY (id);


--
-- Name: teacher_specializations teacher_specializations_pkey; Type: CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.teacher_specializations
    ADD CONSTRAINT teacher_specializations_pkey PRIMARY KEY (id);


--
-- Name: teacher_specializations teacher_specializations_user_id_subject_key; Type: CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.teacher_specializations
    ADD CONSTRAINT teacher_specializations_user_id_subject_key UNIQUE (user_id, subject);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_material_blocks_material_id; Type: INDEX; Schema: public; Owner: paydeya_user
--

CREATE INDEX idx_material_blocks_material_id ON public.material_blocks USING btree (material_id);


--
-- Name: idx_material_blocks_position; Type: INDEX; Schema: public; Owner: paydeya_user
--

CREATE INDEX idx_material_blocks_position ON public.material_blocks USING btree (material_id, "position");


--
-- Name: idx_material_ratings_material_id; Type: INDEX; Schema: public; Owner: paydeya_user
--

CREATE INDEX idx_material_ratings_material_id ON public.material_ratings USING btree (material_id);


--
-- Name: idx_material_ratings_user_id; Type: INDEX; Schema: public; Owner: paydeya_user
--

CREATE INDEX idx_material_ratings_user_id ON public.material_ratings USING btree (user_id);


--
-- Name: idx_materials_author_id; Type: INDEX; Schema: public; Owner: paydeya_user
--

CREATE INDEX idx_materials_author_id ON public.materials USING btree (author_id);


--
-- Name: idx_materials_status; Type: INDEX; Schema: public; Owner: paydeya_user
--

CREATE INDEX idx_materials_status ON public.materials USING btree (status);


--
-- Name: idx_materials_subject; Type: INDEX; Schema: public; Owner: paydeya_user
--

CREATE INDEX idx_materials_subject ON public.materials USING btree (subject);


--
-- Name: idx_teacher_specializations_user_id; Type: INDEX; Schema: public; Owner: paydeya_user
--

CREATE INDEX idx_teacher_specializations_user_id ON public.teacher_specializations USING btree (user_id);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: paydeya_user
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_role; Type: INDEX; Schema: public; Owner: paydeya_user
--

CREATE INDEX idx_users_role ON public.users USING btree (role);


--
-- Name: material_blocks material_blocks_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.material_blocks
    ADD CONSTRAINT material_blocks_material_id_fkey FOREIGN KEY (material_id) REFERENCES public.materials(id) ON DELETE CASCADE;


--
-- Name: material_ratings material_ratings_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.material_ratings
    ADD CONSTRAINT material_ratings_material_id_fkey FOREIGN KEY (material_id) REFERENCES public.materials(id) ON DELETE CASCADE;


--
-- Name: material_ratings material_ratings_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.material_ratings
    ADD CONSTRAINT material_ratings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: materials materials_author_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.materials
    ADD CONSTRAINT materials_author_id_fkey FOREIGN KEY (author_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: specializations specializations_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.specializations
    ADD CONSTRAINT specializations_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: teacher_specializations teacher_specializations_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: paydeya_user
--

ALTER TABLE ONLY public.teacher_specializations
    ADD CONSTRAINT teacher_specializations_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: -; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON SEQUENCES TO paydeya_user;


--
-- Name: DEFAULT PRIVILEGES FOR TYPES; Type: DEFAULT ACL; Schema: -; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON TYPES TO paydeya_user;


--
-- Name: DEFAULT PRIVILEGES FOR FUNCTIONS; Type: DEFAULT ACL; Schema: -; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON FUNCTIONS TO paydeya_user;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: -; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON TABLES TO paydeya_user;


--
-- PostgreSQL database dump complete
--

