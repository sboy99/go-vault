--
-- PostgreSQL database dump
--
--
-- Database version: PostgreSQL 15.2 (Debian 15.2-1.pgdg110+1) on x86_64-pc-linux-gnu, compiled by gcc (Debian 10.2.1-6) 10.2.1 20210110, 64-bit
--
SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = FALSE;
SET client_min_messages = warning;

--
-- Create tables
--
CREATE TABLE migrations (
    id integer,
    timestamp bigint,
    name character varying
);

CREATE TABLE rabble_users (
    id uuid,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    source text,
    rabble_aptos_wallet text,
    utm text,
    telegram_id text,
    rabble_eth_wallet text,
    rabble_sol_wallet text
);

CREATE TABLE telegram_users (
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    telegram_id text,
    first_name text,
    last_name text,
    username text
);

CREATE TABLE sessions (
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    user_id uuid,
    id uuid,
    source text,
    type text,
    hash text
);

CREATE TABLE earn (
    rabble_user_id uuid,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    user_score double precision,
    channel_admin_count numeric,
    channel_owner_count numeric,
    dm_count numeric,
    bots_count numeric,
    groups_count numeric,
    groups_member_count numeric,
    groups_admin_count numeric,
    groups_owner_count numeric,
    channel_count numeric,
    channel_member_count numeric,
    transaction_hash text,
    category text
);

CREATE TABLE referrals (
    id uuid,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    referrer_id text,
    referree_id text
);

CREATE TABLE events (
    id uuid,
    user_id uuid,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    event_type text
);

CREATE TABLE earns (
    id uuid,
    tg_user_score numeric,
    tg_dm_count numeric,
    tg_bots_count numeric,
    tg_groups_count numeric,
    tg_groups_member_count numeric,
    tg_groups_admin_count numeric,
    tg_groups_owner_count numeric,
    tg_channel_count numeric,
    tg_channel_member_count numeric,
    tg_channel_admin_count numeric,
    tg_channel_owner_count numeric,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    rabble_user_id uuid,
    event_score numeric,
    referral_score numeric
);

CREATE TABLE user_referral_info (
    id uuid,
    referral_code text
);

CREATE TABLE transactions (
    user_id uuid,
    amount numeric,
    gas numeric,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    id uuid,
    network text,
    currency text,
    tx_hash text,
    type text,
    from_address text,
    to_address text,
    status text
);

--
-- Create sequences
--
CREATE SEQUENCE public.migrations_id_seq;
ALTER TABLE public.migrations ALTER COLUMN id SET DEFAULT nextval('public.migrations_id_seq'::regclass);




















