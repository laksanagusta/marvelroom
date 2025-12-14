--
-- PostgreSQL database initial schema
--

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

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: trigger_set_timestamp(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.trigger_set_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$;


--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: assignee_transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.assignee_transactions (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT transactions_id_not_null NOT NULL,
    assignee_id uuid CONSTRAINT transactions_assignee_id_not_null NOT NULL,
    name character varying(255) CONSTRAINT transactions_name_not_null NOT NULL,
    type character varying(50) CONSTRAINT transactions_type_not_null NOT NULL,
    subtype character varying(50),
    amount numeric(15,2) DEFAULT 0.00 CONSTRAINT transactions_amount_not_null NOT NULL,
    total_night integer,
    subtotal numeric(15,2) DEFAULT 0.00 CONSTRAINT transactions_subtotal_not_null NOT NULL,
    description text,
    transport_detail text,
    created_at timestamp without time zone DEFAULT now() CONSTRAINT transactions_created_at_not_null NOT NULL,
    updated_at timestamp without time zone DEFAULT now() CONSTRAINT transactions_updated_at_not_null NOT NULL,
    deleted_at timestamp without time zone,
    is_valid boolean DEFAULT true,
    CONSTRAINT chk_transaction_type CHECK (((type)::text = ANY ((ARRAY['accommodation'::character varying, 'transport'::character varying, 'other'::character varying, 'allowance'::character varying])::text[])))
);


COMMENT ON TABLE public.assignee_transactions IS 'Table for storing financial transactions related to business trips';
COMMENT ON COLUMN public.assignee_transactions.type IS 'Type of transaction (accommodation, transport, other, allowance)';
COMMENT ON COLUMN public.assignee_transactions.subtype IS 'Subtype of transaction (hotel, flight, train, etc.)';
COMMENT ON COLUMN public.assignee_transactions.amount IS 'Amount per unit';
COMMENT ON COLUMN public.assignee_transactions.total_night IS 'Total number of nights (for accommodation)';
COMMENT ON COLUMN public.assignee_transactions.subtotal IS 'Total amount for this transaction';
COMMENT ON COLUMN public.assignee_transactions.transport_detail IS 'Details for transportation (flight number, etc.)';
COMMENT ON COLUMN public.assignee_transactions.is_valid IS 'Indicates whether the transaction document has been validated as authentic (true) or potentially fraudulent (false)';


--
-- Name: assignees; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.assignees (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    business_trip_id uuid NOT NULL,
    name character varying(255) NOT NULL,
    spd_number character varying(100) NOT NULL,
    employee_id character varying(100) NOT NULL,
    "position" character varying(255) NOT NULL,
    rank character varying(100) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    employee_name character varying(255),
    employee_number character varying(100)
);


COMMENT ON TABLE public.assignees IS 'Table for storing employees assigned to business trips';
COMMENT ON COLUMN public.assignees.spd_number IS 'SPD number for the assignee';
COMMENT ON COLUMN public.assignees.employee_id IS 'User ID from external API system (response.id)';
COMMENT ON COLUMN public.assignees."position" IS 'Position of the employee';
COMMENT ON COLUMN public.assignees.rank IS 'Rank or grade of the employee';
COMMENT ON COLUMN public.assignees.employee_number IS 'Employee NIP/number from external API system (response.employee_id)';


--
-- Name: business_trip_histories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.business_trip_histories (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    business_trip_id uuid NOT NULL,
    change_type character varying(50) NOT NULL,
    field_name character varying(100),
    old_value text,
    new_value text,
    notes text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    changed_by character varying(255),
    CONSTRAINT chk_history_change_type CHECK (((change_type)::text = ANY ((ARRAY['status_change'::character varying, 'verification_approved'::character varying, 'verification_rejected'::character varying, 'verification_pending'::character varying])::text[])))
);


COMMENT ON TABLE public.business_trip_histories IS 'Table for tracking all important changes on business trips (status changes and verification process)';
COMMENT ON COLUMN public.business_trip_histories.business_trip_id IS 'Reference to the business trip';
COMMENT ON COLUMN public.business_trip_histories.change_type IS 'Type of change (status_change, verification_approved, verification_rejected, verification_pending)';
COMMENT ON COLUMN public.business_trip_histories.field_name IS 'Name of the field that was changed (e.g., "status", "verificator")';
COMMENT ON COLUMN public.business_trip_histories.old_value IS 'Previous value before the change';
COMMENT ON COLUMN public.business_trip_histories.new_value IS 'New value after the change';
COMMENT ON COLUMN public.business_trip_histories.notes IS 'Additional notes or comments about the change';
COMMENT ON COLUMN public.business_trip_histories.created_at IS 'Timestamp when the change was recorded';


--
-- Name: business_trip_verificators; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.business_trip_verificators (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    business_trip_id uuid NOT NULL,
    user_id character varying(100) NOT NULL,
    user_name character varying(255) NOT NULL,
    employee_number character varying(50) NOT NULL,
    "position" character varying(255) NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    verified_at timestamp without time zone,
    verification_notes text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    CONSTRAINT chk_verificator_status CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'approved'::character varying, 'rejected'::character varying])::text[])))
);


COMMENT ON TABLE public.business_trip_verificators IS 'Table for storing users assigned to verify business trips';
COMMENT ON COLUMN public.business_trip_verificators.business_trip_id IS 'Reference to the business trip being verified';
COMMENT ON COLUMN public.business_trip_verificators.user_id IS 'User ID in the system';
COMMENT ON COLUMN public.business_trip_verificators.user_name IS 'Full name of the verificator';
COMMENT ON COLUMN public.business_trip_verificators.employee_number IS 'Employee number/NIP of the verificator';
COMMENT ON COLUMN public.business_trip_verificators."position" IS 'Position of the verificator';
COMMENT ON COLUMN public.business_trip_verificators.status IS 'Verification status (pending, approved, rejected)';
COMMENT ON COLUMN public.business_trip_verificators.verified_at IS 'Timestamp when verification was completed';
COMMENT ON COLUMN public.business_trip_verificators.verification_notes IS 'Notes or comments from verificator';


--
-- Name: business_trips; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.business_trips (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    start_date timestamp without time zone NOT NULL,
    end_date timestamp without time zone NOT NULL,
    activity_purpose character varying(500) NOT NULL,
    destination_city character varying(255) NOT NULL,
    spd_date timestamp without time zone NOT NULL,
    departure_date timestamp without time zone NOT NULL,
    return_date timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    status character varying(20) DEFAULT 'draft'::character varying NOT NULL,
    document_link text,
    business_trip_number character varying(10),
    organization_id uuid,
    CONSTRAINT chk_business_trip_status CHECK (((status)::text = ANY ((ARRAY['draft'::character varying, 'ready_to_verify'::character varying, 'ongoing'::character varying, 'completed'::character varying, 'canceled'::character varying])::text[])))
);


COMMENT ON TABLE public.business_trips IS 'Main table for storing business trip information';
COMMENT ON COLUMN public.business_trips.start_date IS 'Start date of the business trip';
COMMENT ON COLUMN public.business_trips.end_date IS 'End date of the business trip';
COMMENT ON COLUMN public.business_trips.activity_purpose IS 'Purpose of the business trip';
COMMENT ON COLUMN public.business_trips.destination_city IS 'Destination city for the business trip';
COMMENT ON COLUMN public.business_trips.spd_date IS 'SPD (Surat Perintah Dinas) date';
COMMENT ON COLUMN public.business_trips.departure_date IS 'Actual departure date';
COMMENT ON COLUMN public.business_trips.return_date IS 'Actual return date';
COMMENT ON COLUMN public.business_trips.status IS 'Status of the business trip: draft, ready_to_verify, ongoing, completed, or canceled';
COMMENT ON COLUMN public.business_trips.document_link IS 'Link to Google Drive or other document storage for this business trip';
COMMENT ON COLUMN public.business_trips.business_trip_number IS 'Auto-generated business trip number in format BT-XXXX';


--
-- Name: countries; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.countries (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    country_code character varying(100) NOT NULL,
    country_name_id character varying(255) NOT NULL,
    country_name_en character varying(255) NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: country_vaccine_requirements; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.country_vaccine_requirements (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    country_id uuid NOT NULL,
    vaccine_id uuid NOT NULL,
    requirement_type character varying(20) NOT NULL,
    cdc_data jsonb,
    cached_at timestamp without time zone DEFAULT now() NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: master_vaccines; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.master_vaccines (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    vaccine_code character varying(20) NOT NULL,
    vaccine_name_id character varying(255) NOT NULL,
    vaccine_name_en character varying(255) NOT NULL,
    description_id text,
    description_en text,
    vaccine_type character varying(50) NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: organizations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.organizations (
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    code character varying(100) NOT NULL,
    address text,
    type character varying(100) NOT NULL,
    parent_uuid uuid,
    level integer DEFAULT 0,
    path text NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(255) NOT NULL,
    updated_at timestamp with time zone DEFAULT now(),
    updated_by character varying(255) NOT NULL,
    deleted_at timestamp with time zone,
    deleted_by character varying(255)
);


--
-- Name: permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.permissions (
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    action character varying(100) NOT NULL,
    resource character varying(100) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(255) NOT NULL,
    updated_at timestamp with time zone DEFAULT now(),
    updated_by character varying(255) NOT NULL
);


--
-- Name: role_permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.role_permissions (
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    role_uuid uuid NOT NULL,
    permission_uuid uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(255) NOT NULL,
    updated_at timestamp with time zone DEFAULT now(),
    updated_by character varying(255) NOT NULL
);


--
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    is_system boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(255) NOT NULL,
    updated_at timestamp with time zone DEFAULT now(),
    updated_by character varying(255) NOT NULL,
    deleted_at timestamp with time zone,
    deleted_by character varying(255)
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- Name: units; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.units (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    code character varying(100),
    description text,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


COMMENT ON TABLE public.units IS 'Data unit kerja/organisasi yang akan melakukan LAKIP';


--
-- Name: user_roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_roles (
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_uuid uuid NOT NULL,
    role_uuid uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(255) NOT NULL,
    updated_at timestamp with time zone DEFAULT now(),
    updated_by character varying(255) NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    username character varying(100) NOT NULL,
    email character varying(255),
    first_name character varying(255) NOT NULL,
    last_name character varying(255),
    phone_number character varying(20),
    password_hash character varying(255) NOT NULL,
    organization_uuid uuid NOT NULL,
    is_active boolean DEFAULT true,
    last_login_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(255) NOT NULL,
    updated_at timestamp with time zone DEFAULT now(),
    updated_by character varying(255) NOT NULL,
    deleted_at timestamp with time zone,
    deleted_by character varying(255),
    employee_id character varying(50)
);


COMMENT ON COLUMN public.users.employee_id IS 'Employee identifier (NIP) - unique identifier for employee';


--
-- Name: work_paper_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.work_paper_items (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT master_lakip_items_id_not_null NOT NULL,
    number character varying(50) CONSTRAINT master_lakip_items_nomor_not_null NOT NULL,
    statement text CONSTRAINT master_lakip_items_pernyataan_not_null NOT NULL,
    explanation text,
    filling_guide text,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    type character varying(10) DEFAULT 'A'::character varying CONSTRAINT master_lakip_items_type_not_null NOT NULL,
    parent_id uuid,
    level integer DEFAULT 1 CONSTRAINT master_lakip_items_level_not_null NOT NULL,
    sort_order integer DEFAULT 0 CONSTRAINT master_lakip_items_sort_order_not_null NOT NULL
);


COMMENT ON TABLE public.work_paper_items IS 'Master data untuk item-item kertas kerja LAKIP yang berisi struktur pemeriksaan';
COMMENT ON COLUMN public.work_paper_items.number IS 'Nomor urut item pemeriksaan (Kolom B)';
COMMENT ON COLUMN public.work_paper_items.statement IS 'Pernyataan/eksistensi yang harus dicek (Kolom D)';
COMMENT ON COLUMN public.work_paper_items.explanation IS 'Penjelasan detail dari pernyataan (Kolom F)';
COMMENT ON COLUMN public.work_paper_items.filling_guide IS 'Petunjuk pengisian untuk item tersebut (Kolom G)';


--
-- Name: work_paper_notes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.work_paper_notes (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT kertas_kerja_items_id_not_null NOT NULL,
    master_item_id uuid CONSTRAINT kertas_kerja_items_master_item_id_not_null NOT NULL,
    gdrive_link text,
    is_valid boolean,
    notes text,
    last_llm_response jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    work_paper_id uuid DEFAULT gen_random_uuid() NOT NULL
);


COMMENT ON TABLE public.work_paper_notes IS 'Detail items dalam sebuah kertas kerja';
COMMENT ON COLUMN public.work_paper_notes.gdrive_link IS 'Link Google Drive folder yang berisi dokumen pendukung';
COMMENT ON COLUMN public.work_paper_notes.is_valid IS 'Hasil pemeriksaan otomatis dari LLM (true/false)';
COMMENT ON COLUMN public.work_paper_notes.notes IS 'Catatan dari hasil pemeriksaan otomatis';
COMMENT ON COLUMN public.work_paper_notes.last_llm_response IS 'Response lengkap dari LLM dalam format JSON';


--
-- Name: work_paper_signatures; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.work_paper_signatures (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    work_paper_id uuid NOT NULL,
    user_id character varying(255) NOT NULL,
    user_name character varying(255) NOT NULL,
    user_email character varying(255),
    user_role character varying(100),
    signature_data jsonb,
    signature_type character varying(50) DEFAULT 'digital'::character varying,
    status character varying(50) DEFAULT 'pending'::character varying,
    notes text,
    signed_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    CONSTRAINT work_paper_signatures_signature_type_check CHECK (((signature_type)::text = ANY ((ARRAY['digital'::character varying, 'manual'::character varying, 'approval'::character varying])::text[]))),
    CONSTRAINT work_paper_signatures_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'signed'::character varying, 'rejected'::character varying])::text[])))
);


--
-- Name: work_papers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.work_papers (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT kertas_kerja_id_not_null NOT NULL,
    organization_id uuid CONSTRAINT kertas_kerja_unit_id_not_null NOT NULL,
    year integer CONSTRAINT kertas_kerja_tahun_not_null NOT NULL,
    semester integer CONSTRAINT kertas_kerja_semester_not_null NOT NULL,
    status character varying(50) DEFAULT 'draft'::character varying,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT kertas_kerja_semester_check CHECK ((semester = ANY (ARRAY[1, 2])))
);


COMMENT ON TABLE public.work_papers IS 'Header kertas kerja per unit per semester';


--
-- Primary Keys
--

ALTER TABLE ONLY public.assignees
    ADD CONSTRAINT assignees_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.assignees
    ADD CONSTRAINT assignees_business_trip_id_spd_number_key UNIQUE (business_trip_id, spd_number);

ALTER TABLE ONLY public.business_trip_histories
    ADD CONSTRAINT business_trip_histories_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.business_trip_verificators
    ADD CONSTRAINT business_trip_verificators_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.business_trip_verificators
    ADD CONSTRAINT business_trip_verificators_business_trip_id_user_id_key UNIQUE (business_trip_id, user_id);

ALTER TABLE ONLY public.business_trips
    ADD CONSTRAINT business_trips_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.countries
    ADD CONSTRAINT countries_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.countries
    ADD CONSTRAINT countries_country_code_key UNIQUE (country_code);

ALTER TABLE ONLY public.country_vaccine_requirements
    ADD CONSTRAINT country_vaccine_requirements_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.country_vaccine_requirements
    ADD CONSTRAINT country_vaccine_requirements_country_id_vaccine_id_key UNIQUE (country_id, vaccine_id);

ALTER TABLE ONLY public.work_paper_notes
    ADD CONSTRAINT kertas_kerja_items_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.work_paper_items
    ADD CONSTRAINT master_lakip_items_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.master_vaccines
    ADD CONSTRAINT master_vaccines_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.master_vaccines
    ADD CONSTRAINT master_vaccines_vaccine_code_key UNIQUE (vaccine_code);

ALTER TABLE ONLY public.organizations
    ADD CONSTRAINT organizations_pkey PRIMARY KEY (uuid);

ALTER TABLE ONLY public.organizations
    ADD CONSTRAINT organizations_code_key UNIQUE (code);

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (uuid);

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (uuid);

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_uuid_permission_uuid_key UNIQUE (role_uuid, permission_uuid);

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (uuid);

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);

ALTER TABLE ONLY public.assignee_transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.units
    ADD CONSTRAINT units_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.units
    ADD CONSTRAINT units_name_key UNIQUE (name);

ALTER TABLE ONLY public.units
    ADD CONSTRAINT units_code_key UNIQUE (code);

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (uuid);

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_uuid_role_uuid_key UNIQUE (user_uuid, role_uuid);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (uuid);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_employee_id_unique UNIQUE (employee_id);

ALTER TABLE ONLY public.work_paper_signatures
    ADD CONSTRAINT work_paper_signatures_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.work_papers
    ADD CONSTRAINT work_papers_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.work_papers
    ADD CONSTRAINT work_papers_organization_year_semester_key UNIQUE (organization_id, year, semester);


--
-- Indexes
--

CREATE INDEX idx_assignee_transactions_is_valid ON public.assignee_transactions USING btree (is_valid);
CREATE INDEX idx_assignees_business_trip_id ON public.assignees USING btree (business_trip_id);
CREATE INDEX idx_assignees_deleted_at ON public.assignees USING btree (deleted_at);
CREATE INDEX idx_assignees_employee_id ON public.assignees USING btree (employee_id);
CREATE INDEX idx_assignees_spd_number ON public.assignees USING btree (spd_number);
CREATE INDEX idx_business_trip_histories_business_trip_id ON public.business_trip_histories USING btree (business_trip_id);
CREATE INDEX idx_business_trip_histories_change_type ON public.business_trip_histories USING btree (change_type);
CREATE INDEX idx_business_trip_histories_changed_by ON public.business_trip_histories USING btree (changed_by);
CREATE INDEX idx_business_trip_histories_created_at ON public.business_trip_histories USING btree (created_at DESC);
CREATE INDEX idx_business_trip_verificators_business_trip_id ON public.business_trip_verificators USING btree (business_trip_id);
CREATE INDEX idx_business_trip_verificators_deleted_at ON public.business_trip_verificators USING btree (deleted_at);
CREATE INDEX idx_business_trip_verificators_employee_number ON public.business_trip_verificators USING btree (employee_number);
CREATE INDEX idx_business_trip_verificators_status ON public.business_trip_verificators USING btree (status);
CREATE INDEX idx_business_trip_verificators_user_id ON public.business_trip_verificators USING btree (user_id);
CREATE INDEX idx_business_trips_activity_purpose ON public.business_trips USING gin (to_tsvector('english'::regconfig, (activity_purpose)::text));
CREATE UNIQUE INDEX idx_business_trips_business_trip_number ON public.business_trips USING btree (business_trip_number) WHERE (business_trip_number IS NOT NULL);
CREATE INDEX idx_business_trips_created_at ON public.business_trips USING btree (created_at);
CREATE INDEX idx_business_trips_dates ON public.business_trips USING btree (start_date, end_date);
CREATE INDEX idx_business_trips_deleted_at ON public.business_trips USING btree (deleted_at);
CREATE INDEX idx_business_trips_destination_city ON public.business_trips USING btree (destination_city);
CREATE INDEX idx_business_trips_organization_id ON public.business_trips USING btree (organization_id);
CREATE INDEX idx_business_trips_status ON public.business_trips USING btree (status);
CREATE INDEX idx_countries_active ON public.countries USING btree (is_active);
CREATE INDEX idx_countries_code ON public.countries USING btree (country_code);
CREATE INDEX idx_country_vaccine_requirements_country ON public.country_vaccine_requirements USING btree (country_id);
CREATE INDEX idx_country_vaccine_requirements_expires ON public.country_vaccine_requirements USING btree (expires_at);
CREATE INDEX idx_country_vaccine_requirements_vaccine ON public.country_vaccine_requirements USING btree (vaccine_id);
CREATE INDEX idx_kertas_kerja_deleted_at ON public.work_papers USING btree (deleted_at);
CREATE INDEX idx_kertas_kerja_items_deleted_at ON public.work_paper_notes USING btree (deleted_at);
CREATE INDEX idx_kertas_kerja_status ON public.work_papers USING btree (status);
CREATE INDEX idx_master_lakip_items_active ON public.work_paper_items USING btree (is_active);
CREATE INDEX idx_master_lakip_items_deleted_at ON public.work_paper_items USING btree (deleted_at);
CREATE INDEX idx_master_lakip_items_nomor ON public.work_paper_items USING btree (number);
CREATE INDEX idx_master_vaccines_active ON public.master_vaccines USING btree (is_active);
CREATE INDEX idx_master_vaccines_code ON public.master_vaccines USING btree (vaccine_code);
CREATE INDEX idx_master_vaccines_type ON public.master_vaccines USING btree (vaccine_type);
CREATE INDEX idx_transactions_assignee_id ON public.assignee_transactions USING btree (assignee_id);
CREATE INDEX idx_transactions_deleted_at ON public.assignee_transactions USING btree (deleted_at);
CREATE INDEX idx_transactions_subtype ON public.assignee_transactions USING btree (subtype);
CREATE INDEX idx_transactions_type ON public.assignee_transactions USING btree (type);
CREATE INDEX idx_units_active ON public.units USING btree (is_active);
CREATE INDEX idx_units_code ON public.units USING btree (code);
CREATE INDEX idx_units_deleted_at ON public.units USING btree (deleted_at);
CREATE INDEX idx_units_name ON public.units USING btree (name);
CREATE INDEX idx_users_employee_id ON public.users USING btree (employee_id);
CREATE INDEX idx_work_paper_items_level ON public.work_paper_items USING btree (level) WHERE (deleted_at IS NULL);
CREATE INDEX idx_work_paper_items_number_active ON public.work_paper_items USING btree (number) WHERE ((deleted_at IS NULL) AND (number IS NOT NULL));
CREATE INDEX idx_work_paper_items_parent_id ON public.work_paper_items USING btree (parent_id) WHERE (deleted_at IS NULL);
CREATE INDEX idx_work_paper_items_type ON public.work_paper_items USING btree (type) WHERE (deleted_at IS NULL);
CREATE INDEX idx_work_paper_items_type_parent_number ON public.work_paper_items USING btree (type, parent_id, number) WHERE (deleted_at IS NULL);
CREATE INDEX idx_work_paper_notes_gdrive_link ON public.work_paper_notes USING btree (gdrive_link) WHERE ((deleted_at IS NULL) AND (gdrive_link IS NOT NULL));
CREATE INDEX idx_work_paper_notes_is_valid ON public.work_paper_notes USING btree (is_valid) WHERE (deleted_at IS NULL);
CREATE INDEX idx_work_paper_notes_master_item_id ON public.work_paper_notes USING btree (master_item_id) WHERE (deleted_at IS NULL);
CREATE INDEX idx_work_paper_notes_work_paper_id ON public.work_paper_notes USING btree (work_paper_id) WHERE (deleted_at IS NULL);
CREATE INDEX idx_work_paper_signatures_status ON public.work_paper_signatures USING btree (status) WHERE (deleted_at IS NULL);
CREATE UNIQUE INDEX idx_work_paper_signatures_unique_user_paper ON public.work_paper_signatures USING btree (work_paper_id, user_id) WHERE ((deleted_at IS NULL) AND ((status)::text = ANY ((ARRAY['pending'::character varying, 'signed'::character varying, 'rejected'::character varying])::text[])));
CREATE INDEX idx_work_paper_signatures_user_id ON public.work_paper_signatures USING btree (user_id) WHERE (deleted_at IS NULL);
CREATE INDEX idx_work_paper_signatures_work_paper_id ON public.work_paper_signatures USING btree (work_paper_id) WHERE (deleted_at IS NULL);
CREATE INDEX idx_work_papers_organization_id ON public.work_papers USING btree (organization_id) WHERE (deleted_at IS NULL);
CREATE INDEX idx_work_papers_organization_year_semester ON public.work_papers USING btree (organization_id, year, semester) WHERE (deleted_at IS NULL);
CREATE INDEX idx_work_papers_year_semester ON public.work_papers USING btree (year, semester) WHERE (deleted_at IS NULL);


--
-- Triggers
--

CREATE TRIGGER set_timestamp_kertas_kerja_items BEFORE UPDATE ON public.work_paper_notes FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();
CREATE TRIGGER set_timestamp_master_lakip_items BEFORE UPDATE ON public.work_paper_items FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();
CREATE TRIGGER set_timestamp_units BEFORE UPDATE ON public.units FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();
CREATE TRIGGER update_assignees_updated_at BEFORE UPDATE ON public.assignees FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_business_trip_verificators_updated_at BEFORE UPDATE ON public.business_trip_verificators FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_business_trips_updated_at BEFORE UPDATE ON public.business_trips FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON public.assignee_transactions FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_work_paper_items_updated_at BEFORE UPDATE ON public.work_paper_items FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_work_paper_notes_updated_at BEFORE UPDATE ON public.work_paper_notes FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_work_paper_signatures_updated_at BEFORE UPDATE ON public.work_paper_signatures FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
CREATE TRIGGER update_work_papers_updated_at BEFORE UPDATE ON public.work_papers FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Foreign Keys
--

ALTER TABLE ONLY public.assignees
    ADD CONSTRAINT assignees_business_trip_id_fkey FOREIGN KEY (business_trip_id) REFERENCES public.business_trips(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.business_trip_histories
    ADD CONSTRAINT business_trip_histories_business_trip_id_fkey FOREIGN KEY (business_trip_id) REFERENCES public.business_trips(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.business_trip_verificators
    ADD CONSTRAINT business_trip_verificators_business_trip_id_fkey FOREIGN KEY (business_trip_id) REFERENCES public.business_trips(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.country_vaccine_requirements
    ADD CONSTRAINT country_vaccine_requirements_country_id_fkey FOREIGN KEY (country_id) REFERENCES public.countries(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.country_vaccine_requirements
    ADD CONSTRAINT country_vaccine_requirements_vaccine_id_fkey FOREIGN KEY (vaccine_id) REFERENCES public.master_vaccines(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.work_paper_notes
    ADD CONSTRAINT kertas_kerja_items_master_item_id_fkey FOREIGN KEY (master_item_id) REFERENCES public.work_paper_items(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.work_paper_items
    ADD CONSTRAINT master_lakip_items_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.work_paper_items(id) ON DELETE SET NULL;

ALTER TABLE ONLY public.organizations
    ADD CONSTRAINT organizations_parent_uuid_fkey FOREIGN KEY (parent_uuid) REFERENCES public.organizations(uuid);

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_permission_uuid_fkey FOREIGN KEY (permission_uuid) REFERENCES public.permissions(uuid) ON DELETE CASCADE;

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_uuid_fkey FOREIGN KEY (role_uuid) REFERENCES public.roles(uuid) ON DELETE CASCADE;

ALTER TABLE ONLY public.assignee_transactions
    ADD CONSTRAINT transactions_assignee_id_fkey FOREIGN KEY (assignee_id) REFERENCES public.assignees(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_role_uuid_fkey FOREIGN KEY (role_uuid) REFERENCES public.roles(uuid) ON DELETE CASCADE;

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_uuid_fkey FOREIGN KEY (user_uuid) REFERENCES public.users(uuid) ON DELETE CASCADE;

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_organization_uuid_fkey FOREIGN KEY (organization_uuid) REFERENCES public.organizations(uuid);

ALTER TABLE ONLY public.work_paper_notes
    ADD CONSTRAINT work_paper_notes_master_item_id_fkey FOREIGN KEY (master_item_id) REFERENCES public.work_paper_items(id) ON DELETE RESTRICT;

ALTER TABLE ONLY public.work_paper_signatures
    ADD CONSTRAINT work_paper_signatures_work_paper_id_fkey FOREIGN KEY (work_paper_id) REFERENCES public.work_papers(id) ON DELETE CASCADE;
