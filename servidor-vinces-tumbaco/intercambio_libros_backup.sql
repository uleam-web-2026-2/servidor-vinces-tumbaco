--
-- PostgreSQL database dump
--

\restrict MbZddFNsnWll4AyOjlRrcLtDCtSLE91HZ4scHPATGIezV9Ts2zSZRfTZRzDetWg

-- Dumped from database version 16.15 (Debian 16.15-1.pgdg13+2)
-- Dumped by pg_dump version 16.15 (Debian 16.15-1.pgdg13+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
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
-- Name: intercambios; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.intercambios (
    id bigint NOT NULL,
    libro_id bigint,
    persona text,
    estado text
);


ALTER TABLE public.intercambios OWNER TO postgres;

--
-- Name: intercambios_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.intercambios_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.intercambios_id_seq OWNER TO postgres;

--
-- Name: intercambios_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.intercambios_id_seq OWNED BY public.intercambios.id;


--
-- Name: libros; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.libros (
    id bigint NOT NULL,
    titulo text,
    autor text,
    genero text,
    estado text
);


ALTER TABLE public.libros OWNER TO postgres;

--
-- Name: libros_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.libros_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.libros_id_seq OWNER TO postgres;

--
-- Name: libros_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.libros_id_seq OWNED BY public.libros.id;


--
-- Name: prestamos; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.prestamos (
    id bigint NOT NULL,
    libro_id bigint,
    usuario text,
    estado text
);


ALTER TABLE public.prestamos OWNER TO postgres;

--
-- Name: prestamos_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.prestamos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.prestamos_id_seq OWNER TO postgres;

--
-- Name: prestamos_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.prestamos_id_seq OWNED BY public.prestamos.id;


--
-- Name: intercambios id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intercambios ALTER COLUMN id SET DEFAULT nextval('public.intercambios_id_seq'::regclass);


--
-- Name: libros id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.libros ALTER COLUMN id SET DEFAULT nextval('public.libros_id_seq'::regclass);


--
-- Name: prestamos id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.prestamos ALTER COLUMN id SET DEFAULT nextval('public.prestamos_id_seq'::regclass);


--
-- Data for Name: intercambios; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.intercambios (id, libro_id, persona, estado) FROM stdin;
2	2	Carlos	aceptado
3	3	María	entregado
1	1	Ana	aceptado
\.


--
-- Data for Name: libros; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.libros (id, titulo, autor, genero, estado) FROM stdin;
3	Don Quijote de la Mancha	Miguel de Cervantes	Novela	\N
1	El Principito	Antoine de Saint-Exupery	Novela	prestado
2	Cien anos de soledad	Gabriel Garcia Marquez	Fantasía	disponible
\.


--
-- Data for Name: prestamos; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.prestamos (id, libro_id, usuario, estado) FROM stdin;
\.


--
-- Name: intercambios_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.intercambios_id_seq', 4, true);


--
-- Name: libros_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.libros_id_seq', 6, true);


--
-- Name: prestamos_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.prestamos_id_seq', 1, false);


--
-- Name: intercambios intercambios_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intercambios
    ADD CONSTRAINT intercambios_pkey PRIMARY KEY (id);


--
-- Name: libros libros_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.libros
    ADD CONSTRAINT libros_pkey PRIMARY KEY (id);


--
-- Name: prestamos prestamos_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.prestamos
    ADD CONSTRAINT prestamos_pkey PRIMARY KEY (id);


--
-- Name: intercambios fk_libros_intercambios; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intercambios
    ADD CONSTRAINT fk_libros_intercambios FOREIGN KEY (libro_id) REFERENCES public.libros(id);


--
-- Name: prestamos fk_libros_prestamos; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.prestamos
    ADD CONSTRAINT fk_libros_prestamos FOREIGN KEY (libro_id) REFERENCES public.libros(id);


--
-- PostgreSQL database dump complete
--

\unrestrict MbZddFNsnWll4AyOjlRrcLtDCtSLE91HZ4scHPATGIezV9Ts2zSZRfTZRzDetWg

