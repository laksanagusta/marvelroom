-- Drop all foreign keys
ALTER TABLE IF EXISTS public.work_paper_signatures DROP CONSTRAINT IF EXISTS work_paper_signatures_work_paper_id_fkey;
ALTER TABLE IF EXISTS public.work_paper_notes DROP CONSTRAINT IF EXISTS work_paper_notes_master_item_id_fkey;
ALTER TABLE IF EXISTS public.work_paper_notes DROP CONSTRAINT IF EXISTS kertas_kerja_items_master_item_id_fkey;
ALTER TABLE IF EXISTS public.users DROP CONSTRAINT IF EXISTS users_organization_uuid_fkey;
ALTER TABLE IF EXISTS public.user_roles DROP CONSTRAINT IF EXISTS user_roles_user_uuid_fkey;
ALTER TABLE IF EXISTS public.user_roles DROP CONSTRAINT IF EXISTS user_roles_role_uuid_fkey;
ALTER TABLE IF EXISTS public.assignee_transactions DROP CONSTRAINT IF EXISTS transactions_assignee_id_fkey;
ALTER TABLE IF EXISTS public.role_permissions DROP CONSTRAINT IF EXISTS role_permissions_role_uuid_fkey;
ALTER TABLE IF EXISTS public.role_permissions DROP CONSTRAINT IF EXISTS role_permissions_permission_uuid_fkey;
ALTER TABLE IF EXISTS public.organizations DROP CONSTRAINT IF EXISTS organizations_parent_uuid_fkey;
ALTER TABLE IF EXISTS public.work_paper_items DROP CONSTRAINT IF EXISTS master_lakip_items_parent_id_fkey;
ALTER TABLE IF EXISTS public.country_vaccine_requirements DROP CONSTRAINT IF EXISTS country_vaccine_requirements_vaccine_id_fkey;
ALTER TABLE IF EXISTS public.country_vaccine_requirements DROP CONSTRAINT IF EXISTS country_vaccine_requirements_country_id_fkey;
ALTER TABLE IF EXISTS public.business_trip_verificators DROP CONSTRAINT IF EXISTS business_trip_verificators_business_trip_id_fkey;
ALTER TABLE IF EXISTS public.business_trip_histories DROP CONSTRAINT IF EXISTS business_trip_histories_business_trip_id_fkey;
ALTER TABLE IF EXISTS public.assignees DROP CONSTRAINT IF EXISTS assignees_business_trip_id_fkey;

-- Drop all triggers
DROP TRIGGER IF EXISTS update_work_papers_updated_at ON public.work_papers;
DROP TRIGGER IF EXISTS update_work_paper_signatures_updated_at ON public.work_paper_signatures;
DROP TRIGGER IF EXISTS update_work_paper_notes_updated_at ON public.work_paper_notes;
DROP TRIGGER IF EXISTS update_work_paper_items_updated_at ON public.work_paper_items;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON public.assignee_transactions;
DROP TRIGGER IF EXISTS update_business_trips_updated_at ON public.business_trips;
DROP TRIGGER IF EXISTS update_business_trip_verificators_updated_at ON public.business_trip_verificators;
DROP TRIGGER IF EXISTS update_assignees_updated_at ON public.assignees;
DROP TRIGGER IF EXISTS set_timestamp_units ON public.units;
DROP TRIGGER IF EXISTS set_timestamp_master_lakip_items ON public.work_paper_items;
DROP TRIGGER IF EXISTS set_timestamp_kertas_kerja_items ON public.work_paper_notes;

-- Drop all tables
DROP TABLE IF EXISTS public.work_papers CASCADE;
DROP TABLE IF EXISTS public.work_paper_signatures CASCADE;
DROP TABLE IF EXISTS public.work_paper_notes CASCADE;
DROP TABLE IF EXISTS public.work_paper_items CASCADE;
DROP TABLE IF EXISTS public.users CASCADE;
DROP TABLE IF EXISTS public.user_roles CASCADE;
DROP TABLE IF EXISTS public.units CASCADE;
DROP TABLE IF EXISTS public.schema_migrations CASCADE;
DROP TABLE IF EXISTS public.roles CASCADE;
DROP TABLE IF EXISTS public.role_permissions CASCADE;
DROP TABLE IF EXISTS public.permissions CASCADE;
DROP TABLE IF EXISTS public.organizations CASCADE;
DROP TABLE IF EXISTS public.master_vaccines CASCADE;
DROP TABLE IF EXISTS public.country_vaccine_requirements CASCADE;
DROP TABLE IF EXISTS public.countries CASCADE;
DROP TABLE IF EXISTS public.business_trips CASCADE;
DROP TABLE IF EXISTS public.business_trip_verificators CASCADE;
DROP TABLE IF EXISTS public.business_trip_histories CASCADE;
DROP TABLE IF EXISTS public.assignees CASCADE;
DROP TABLE IF EXISTS public.assignee_transactions CASCADE;

-- Drop functions
DROP FUNCTION IF EXISTS public.update_updated_at_column();
DROP FUNCTION IF EXISTS public.trigger_set_timestamp();

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
