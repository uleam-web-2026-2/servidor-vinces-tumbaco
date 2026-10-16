--
-- PostgreSQL database dump
--

\restrict EPUUVGemNj38I83yF1M42ExNMblMtIXrW0q8v8L3hTfooNj4GwD7xWXgdhKvcyW

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
-- Name: autors; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.autors (
    id bigint NOT NULL,
    nombre character varying(50),
    apellido character varying(50),
    nacionalidad character varying(50),
    fecha_nacimiento timestamp with time zone
);


ALTER TABLE public.autors OWNER TO postgres;

--
-- Name: autors_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.autors_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.autors_id_seq OWNER TO postgres;

--
-- Name: autors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.autors_id_seq OWNED BY public.autors.id;


--
-- Name: categoria; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categoria (
    id bigint NOT NULL,
    nombre character varying(50),
    descripcion character varying(200)
);


ALTER TABLE public.categoria OWNER TO postgres;

--
-- Name: categoria_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.categoria_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categoria_id_seq OWNER TO postgres;

--
-- Name: categoria_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categoria_id_seq OWNED BY public.categoria.id;


--
-- Name: detalle_intercambios; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.detalle_intercambios (
    id bigint NOT NULL,
    intercambio_id bigint,
    ejemplar_entregado_id bigint,
    ejemplar_recibido_id bigint,
    cantidad bigint
);


ALTER TABLE public.detalle_intercambios OWNER TO postgres;

--
-- Name: detalle_intercambios_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.detalle_intercambios_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.detalle_intercambios_id_seq OWNER TO postgres;

--
-- Name: detalle_intercambios_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.detalle_intercambios_id_seq OWNED BY public.detalle_intercambios.id;


--
-- Name: ejemplar_libros; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ejemplar_libros (
    id bigint NOT NULL,
    libro_id bigint NOT NULL,
    usuario_id bigint NOT NULL,
    estado_fisico character varying(30),
    disponibilidad character varying(30),
    fecha_registro timestamp with time zone,
    observaciones text
);


ALTER TABLE public.ejemplar_libros OWNER TO postgres;

--
-- Name: ejemplar_libros_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ejemplar_libros_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ejemplar_libros_id_seq OWNER TO postgres;

--
-- Name: ejemplar_libros_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ejemplar_libros_id_seq OWNED BY public.ejemplar_libros.id;


--
-- Name: entregas; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.entregas (
    id bigint NOT NULL,
    intercambio_id bigint,
    lugar_entrega character varying(150),
    fecha_programada timestamp with time zone,
    fecha_real timestamp with time zone,
    estado character varying(30)
);


ALTER TABLE public.entregas OWNER TO postgres;

--
-- Name: entregas_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.entregas_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.entregas_id_seq OWNER TO postgres;

--
-- Name: entregas_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.entregas_id_seq OWNED BY public.entregas.id;


--
-- Name: intercambios; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.intercambios (
    id bigint NOT NULL,
    libro_id bigint,
    persona text,
    estado text DEFAULT 'pendiente'::text NOT NULL,
    usuario_ofrece_id bigint,
    usuario_recibe_id bigint,
    fecha_solicitud timestamp with time zone,
    fecha_confirmacion timestamp with time zone,
    fecha_finalizacion timestamp with time zone,
    observaciones text,
    solicitante_id bigint,
    creado timestamp with time zone
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
-- Name: libro_autors; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.libro_autors (
    libro_id bigint NOT NULL,
    autor_id bigint NOT NULL,
    rol character varying(30)
);


ALTER TABLE public.libro_autors OWNER TO postgres;

--
-- Name: libro_categoria; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.libro_categoria (
    libro_id bigint NOT NULL,
    categoria_id bigint NOT NULL
);


ALTER TABLE public.libro_categoria OWNER TO postgres;

--
-- Name: libros; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.libros (
    id bigint NOT NULL,
    titulo character varying(150) NOT NULL,
    autor text,
    genero text,
    estado text,
    isbn character varying(20),
    descripcion text,
    anio_publicacion bigint,
    editorial character varying(100),
    idioma character varying(50),
    paginas bigint,
    disponible boolean DEFAULT true NOT NULL,
    usuario_id bigint
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
-- Name: resenas; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resenas (
    id bigint NOT NULL,
    intercambio_id bigint,
    usuario_autor_id bigint,
    usuario_evaluado_id bigint,
    puntuacion bigint,
    comentario text,
    fecha timestamp with time zone
);


ALTER TABLE public.resenas OWNER TO postgres;

--
-- Name: resenas_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.resenas_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.resenas_id_seq OWNER TO postgres;

--
-- Name: resenas_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.resenas_id_seq OWNED BY public.resenas.id;


--
-- Name: solicitud_intercambios; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.solicitud_intercambios (
    id bigint NOT NULL,
    usuario_solicitante_id bigint,
    ejemplar_ofrecido_id bigint,
    ejemplar_deseado_id bigint,
    fecha_solicitud timestamp with time zone,
    estado character varying(30),
    mensaje text
);


ALTER TABLE public.solicitud_intercambios OWNER TO postgres;

--
-- Name: solicitud_intercambios_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.solicitud_intercambios_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.solicitud_intercambios_id_seq OWNER TO postgres;

--
-- Name: solicitud_intercambios_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.solicitud_intercambios_id_seq OWNED BY public.solicitud_intercambios.id;


--
-- Name: usuarios; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.usuarios (
    id bigint NOT NULL,
    nombre character varying(50) NOT NULL,
    apellido character varying(50),
    correo character varying(100) NOT NULL,
    contrasena character varying(255),
    telefono character varying(20),
    direccion character varying(150),
    ciudad character varying(50),
    fecha_registro timestamp with time zone,
    estado character varying(20),
    rol text NOT NULL
);


ALTER TABLE public.usuarios OWNER TO postgres;

--
-- Name: usuarios_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.usuarios_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.usuarios_id_seq OWNER TO postgres;

--
-- Name: usuarios_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.usuarios_id_seq OWNED BY public.usuarios.id;


--
-- Name: autors id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.autors ALTER COLUMN id SET DEFAULT nextval('public.autors_id_seq'::regclass);


--
-- Name: categoria id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categoria ALTER COLUMN id SET DEFAULT nextval('public.categoria_id_seq'::regclass);


--
-- Name: detalle_intercambios id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detalle_intercambios ALTER COLUMN id SET DEFAULT nextval('public.detalle_intercambios_id_seq'::regclass);


--
-- Name: ejemplar_libros id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ejemplar_libros ALTER COLUMN id SET DEFAULT nextval('public.ejemplar_libros_id_seq'::regclass);


--
-- Name: entregas id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.entregas ALTER COLUMN id SET DEFAULT nextval('public.entregas_id_seq'::regclass);


--
-- Name: intercambios id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intercambios ALTER COLUMN id SET DEFAULT nextval('public.intercambios_id_seq'::regclass);


--
-- Name: libros id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.libros ALTER COLUMN id SET DEFAULT nextval('public.libros_id_seq'::regclass);


--
-- Name: resenas id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resenas ALTER COLUMN id SET DEFAULT nextval('public.resenas_id_seq'::regclass);


--
-- Name: solicitud_intercambios id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.solicitud_intercambios ALTER COLUMN id SET DEFAULT nextval('public.solicitud_intercambios_id_seq'::regclass);


--
-- Name: usuarios id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.usuarios ALTER COLUMN id SET DEFAULT nextval('public.usuarios_id_seq'::regclass);


--
-- Data for Name: autors; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.autors (id, nombre, apellido, nacionalidad, fecha_nacimiento) FROM stdin;
\.


--
-- Data for Name: categoria; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categoria (id, nombre, descripcion) FROM stdin;
\.


--
-- Data for Name: detalle_intercambios; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.detalle_intercambios (id, intercambio_id, ejemplar_entregado_id, ejemplar_recibido_id, cantidad) FROM stdin;
\.


--
-- Data for Name: ejemplar_libros; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.ejemplar_libros (id, libro_id, usuario_id, estado_fisico, disponibilidad, fecha_registro, observaciones) FROM stdin;
1	1	2	bueno	disponible	2026-10-14 21:40:07.653542+00	\N
2	2	1	bueno	disponible	2026-10-14 21:40:07.653542+00	\N
3	3	1	bueno	disponible	2026-10-14 21:40:07.653542+00	\N
\.


--
-- Data for Name: entregas; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.entregas (id, intercambio_id, lugar_entrega, fecha_programada, fecha_real, estado) FROM stdin;
\.


--
-- Data for Name: intercambios; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.intercambios (id, libro_id, persona, estado, usuario_ofrece_id, usuario_recibe_id, fecha_solicitud, fecha_confirmacion, fecha_finalizacion, observaciones, solicitante_id, creado) FROM stdin;
2	2	Carlos	aceptado	\N	\N	\N	\N	\N	\N	\N	\N
3	3	María	entregado	\N	\N	\N	\N	\N	\N	\N	\N
1	1	Ana	aceptado	\N	\N	\N	\N	\N	\N	\N	\N
\.


--
-- Data for Name: libro_autors; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.libro_autors (libro_id, autor_id, rol) FROM stdin;
\.


--
-- Data for Name: libro_categoria; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.libro_categoria (libro_id, categoria_id) FROM stdin;
\.


--
-- Data for Name: libros; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.libros (id, titulo, autor, genero, estado, isbn, descripcion, anio_publicacion, editorial, idioma, paginas, disponible, usuario_id) FROM stdin;
3	Don Quijote de la Mancha	Miguel de Cervantes	Novela	\N	\N	\N	\N	\N	\N	\N	t	1
1	El Principito	Antoine de Saint-Exupery	Novela	prestado	\N	\N	\N	\N	\N	\N	t	2
2	Cien anos de soledad	Gabriel Garcia Marquez	Fantasía	disponible	\N	\N	\N	\N	\N	\N	t	1
\.


--
-- Data for Name: resenas; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.resenas (id, intercambio_id, usuario_autor_id, usuario_evaluado_id, puntuacion, comentario, fecha) FROM stdin;
\.


--
-- Data for Name: solicitud_intercambios; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.solicitud_intercambios (id, usuario_solicitante_id, ejemplar_ofrecido_id, ejemplar_deseado_id, fecha_solicitud, estado, mensaje) FROM stdin;
\.


--
-- Data for Name: usuarios; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.usuarios (id, nombre, apellido, correo, contrasena, telefono, direccion, ciudad, fecha_registro, estado, rol) FROM stdin;
1	Ana		ana@email.com					0001-01-01 00:00:00+00		usuario
2	Carlos		carlos@email.com					0001-01-01 00:00:00+00		usuario
3	Administrador		admin@email.com					0001-01-01 00:00:00+00		administrador
\.


--
-- Name: autors_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.autors_id_seq', 1, false);


--
-- Name: categoria_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categoria_id_seq', 1, false);


--
-- Name: detalle_intercambios_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.detalle_intercambios_id_seq', 1, false);


--
-- Name: ejemplar_libros_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ejemplar_libros_id_seq', 3, true);


--
-- Name: entregas_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.entregas_id_seq', 1, false);


--
-- Name: intercambios_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.intercambios_id_seq', 4, true);


--
-- Name: libros_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.libros_id_seq', 8, true);


--
-- Name: resenas_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.resenas_id_seq', 1, false);


--
-- Name: solicitud_intercambios_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.solicitud_intercambios_id_seq', 1, false);


--
-- Name: usuarios_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.usuarios_id_seq', 3, true);


--
-- Name: autors autors_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.autors
    ADD CONSTRAINT autors_pkey PRIMARY KEY (id);


--
-- Name: categoria categoria_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categoria
    ADD CONSTRAINT categoria_pkey PRIMARY KEY (id);


--
-- Name: detalle_intercambios detalle_intercambios_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detalle_intercambios
    ADD CONSTRAINT detalle_intercambios_pkey PRIMARY KEY (id);


--
-- Name: ejemplar_libros ejemplar_libros_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ejemplar_libros
    ADD CONSTRAINT ejemplar_libros_pkey PRIMARY KEY (id);


--
-- Name: entregas entregas_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.entregas
    ADD CONSTRAINT entregas_pkey PRIMARY KEY (id);


--
-- Name: intercambios intercambios_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intercambios
    ADD CONSTRAINT intercambios_pkey PRIMARY KEY (id);


--
-- Name: libro_autors libro_autors_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.libro_autors
    ADD CONSTRAINT libro_autors_pkey PRIMARY KEY (libro_id, autor_id);


--
-- Name: libro_categoria libro_categoria_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.libro_categoria
    ADD CONSTRAINT libro_categoria_pkey PRIMARY KEY (libro_id, categoria_id);


--
-- Name: libros libros_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.libros
    ADD CONSTRAINT libros_pkey PRIMARY KEY (id);


--
-- Name: resenas resenas_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resenas
    ADD CONSTRAINT resenas_pkey PRIMARY KEY (id);


--
-- Name: solicitud_intercambios solicitud_intercambios_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.solicitud_intercambios
    ADD CONSTRAINT solicitud_intercambios_pkey PRIMARY KEY (id);


--
-- Name: usuarios uni_usuarios_correo; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.usuarios
    ADD CONSTRAINT uni_usuarios_correo UNIQUE (correo);


--
-- Name: usuarios usuarios_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.usuarios
    ADD CONSTRAINT usuarios_pkey PRIMARY KEY (id);


--
-- Name: idx_entregas_intercambio_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_entregas_intercambio_id ON public.entregas USING btree (intercambio_id);


--
-- Name: idx_libros_isbn; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_libros_isbn ON public.libros USING btree (isbn);


--
-- Name: intercambios fk_libros_intercambios; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intercambios
    ADD CONSTRAINT fk_libros_intercambios FOREIGN KEY (libro_id) REFERENCES public.libros(id);


--
-- PostgreSQL database dump complete
--

\unrestrict EPUUVGemNj38I83yF1M42ExNMblMtIXrW0q8v8L3hTfooNj4GwD7xWXgdhKvcyW

