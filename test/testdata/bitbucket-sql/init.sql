--
-- PostgreSQL database dump
--

-- Dumped from database version 12.4 (Debian 12.4-1.pgdg100+1)
-- Dumped by pg_dump version 12.4 (Debian 12.4-1.pgdg100+1)

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
-- Name: AO_02A6C0_REJECTED_REF; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_02A6C0_REJECTED_REF" (
    "ID" integer NOT NULL,
    "REF_DISPLAY_ID" character varying(450) NOT NULL,
    "REF_ID" character varying(450) NOT NULL,
    "REF_STATUS" integer DEFAULT 0 NOT NULL,
    "REPOSITORY_ID" integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."AO_02A6C0_REJECTED_REF" OWNER TO bitbucketuser;

--
-- Name: AO_02A6C0_REJECTED_REF_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_02A6C0_REJECTED_REF_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_02A6C0_REJECTED_REF_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_02A6C0_REJECTED_REF_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_02A6C0_REJECTED_REF_ID_seq" OWNED BY public."AO_02A6C0_REJECTED_REF"."ID";


--
-- Name: AO_02A6C0_SYNC_CONFIG; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_02A6C0_SYNC_CONFIG" (
    "IS_ENABLED" boolean NOT NULL,
    "LAST_SYNC" timestamp without time zone NOT NULL,
    "REPOSITORY_ID" integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."AO_02A6C0_SYNC_CONFIG" OWNER TO bitbucketuser;

--
-- Name: AO_0E97B5_REPOSITORY_SHORTCUT; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_0E97B5_REPOSITORY_SHORTCUT" (
    "APPLICATION_LINK_ID" character varying(255),
    "CREATED_DATE" timestamp without time zone NOT NULL,
    "ID" integer NOT NULL,
    "LABEL" character varying(255) NOT NULL,
    "PRODUCT_TYPE" character varying(255),
    "REPOSITORY_ID" integer DEFAULT 0 NOT NULL,
    "URL" character varying(450) NOT NULL
);


ALTER TABLE public."AO_0E97B5_REPOSITORY_SHORTCUT" OWNER TO bitbucketuser;

--
-- Name: AO_0E97B5_REPOSITORY_SHORTCUT_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_0E97B5_REPOSITORY_SHORTCUT_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_0E97B5_REPOSITORY_SHORTCUT_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_0E97B5_REPOSITORY_SHORTCUT_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_0E97B5_REPOSITORY_SHORTCUT_ID_seq" OWNED BY public."AO_0E97B5_REPOSITORY_SHORTCUT"."ID";


--
-- Name: AO_2AD648_INSIGHT_ANNOTATION; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_2AD648_INSIGHT_ANNOTATION" (
    "EXTERNAL_ID" character varying(450),
    "FK_REPORT_ID" bigint NOT NULL,
    "ID" bigint NOT NULL,
    "LINE" integer DEFAULT 0 NOT NULL,
    "LINK" text,
    "MESSAGE" text NOT NULL,
    "PATH" text,
    "PATH_MD5" character varying(32),
    "SEVERITY_ID" integer DEFAULT 0 NOT NULL,
    "TYPE_ID" integer
);


ALTER TABLE public."AO_2AD648_INSIGHT_ANNOTATION" OWNER TO bitbucketuser;

--
-- Name: AO_2AD648_INSIGHT_ANNOTATION_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_2AD648_INSIGHT_ANNOTATION_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_2AD648_INSIGHT_ANNOTATION_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_2AD648_INSIGHT_ANNOTATION_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_2AD648_INSIGHT_ANNOTATION_ID_seq" OWNED BY public."AO_2AD648_INSIGHT_ANNOTATION"."ID";


--
-- Name: AO_2AD648_INSIGHT_REPORT; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_2AD648_INSIGHT_REPORT" (
    "AUTHOR_ID" integer,
    "COMMIT_ID" character varying(40) NOT NULL,
    "COVERAGE_PROVIDER_KEY" character varying(450),
    "CREATED_DATE" timestamp without time zone NOT NULL,
    "DATA" text,
    "DETAILS" text,
    "ID" bigint NOT NULL,
    "LINK" text,
    "LOGO" text,
    "REPORTER" character varying(450),
    "REPORT_KEY" character varying(450) NOT NULL,
    "REPOSITORY_ID" integer NOT NULL,
    "RESULT_ID" integer,
    "TITLE" character varying(450) NOT NULL
);


ALTER TABLE public."AO_2AD648_INSIGHT_REPORT" OWNER TO bitbucketuser;

--
-- Name: AO_2AD648_INSIGHT_REPORT_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_2AD648_INSIGHT_REPORT_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_2AD648_INSIGHT_REPORT_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_2AD648_INSIGHT_REPORT_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_2AD648_INSIGHT_REPORT_ID_seq" OWNED BY public."AO_2AD648_INSIGHT_REPORT"."ID";


--
-- Name: AO_2AD648_MERGE_CHECK; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_2AD648_MERGE_CHECK" (
    "ANNOTATION_SEVERITY" character varying(255),
    "ID" bigint NOT NULL,
    "MUST_PASS" boolean,
    "REPORT_KEY" character varying(450) NOT NULL,
    "RESOURCE_ID" integer NOT NULL,
    "SCOPE_TYPE" character varying(255) NOT NULL
);


ALTER TABLE public."AO_2AD648_MERGE_CHECK" OWNER TO bitbucketuser;

--
-- Name: AO_2AD648_MERGE_CHECK_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_2AD648_MERGE_CHECK_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_2AD648_MERGE_CHECK_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_2AD648_MERGE_CHECK_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_2AD648_MERGE_CHECK_ID_seq" OWNED BY public."AO_2AD648_MERGE_CHECK"."ID";


--
-- Name: AO_33D892_COMMENT_JIRA_ISSUE; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_33D892_COMMENT_JIRA_ISSUE" (
    "COMMENT_ID" bigint DEFAULT 0 NOT NULL,
    "ID" integer NOT NULL,
    "ISSUE_KEY" character varying(450) NOT NULL
);


ALTER TABLE public."AO_33D892_COMMENT_JIRA_ISSUE" OWNER TO bitbucketuser;

--
-- Name: AO_33D892_COMMENT_JIRA_ISSUE_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_33D892_COMMENT_JIRA_ISSUE_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_33D892_COMMENT_JIRA_ISSUE_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_33D892_COMMENT_JIRA_ISSUE_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_33D892_COMMENT_JIRA_ISSUE_ID_seq" OWNED BY public."AO_33D892_COMMENT_JIRA_ISSUE"."ID";


--
-- Name: AO_38321B_CUSTOM_CONTENT_LINK; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_38321B_CUSTOM_CONTENT_LINK" (
    "CONTENT_KEY" character varying(255),
    "ID" integer NOT NULL,
    "LINK_LABEL" character varying(255),
    "LINK_URL" character varying(255),
    "SEQUENCE" integer DEFAULT 0
);


ALTER TABLE public."AO_38321B_CUSTOM_CONTENT_LINK" OWNER TO bitbucketuser;

--
-- Name: AO_38321B_CUSTOM_CONTENT_LINK_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_38321B_CUSTOM_CONTENT_LINK_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_38321B_CUSTOM_CONTENT_LINK_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_38321B_CUSTOM_CONTENT_LINK_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_38321B_CUSTOM_CONTENT_LINK_ID_seq" OWNED BY public."AO_38321B_CUSTOM_CONTENT_LINK"."ID";


--
-- Name: AO_38F373_COMMENT_LIKE; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_38F373_COMMENT_LIKE" (
    "COMMENT_ID" bigint DEFAULT 0 NOT NULL,
    "ID" bigint NOT NULL,
    "USER_ID" integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."AO_38F373_COMMENT_LIKE" OWNER TO bitbucketuser;

--
-- Name: AO_38F373_COMMENT_LIKE_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_38F373_COMMENT_LIKE_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_38F373_COMMENT_LIKE_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_38F373_COMMENT_LIKE_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_38F373_COMMENT_LIKE_ID_seq" OWNED BY public."AO_38F373_COMMENT_LIKE"."ID";


--
-- Name: AO_4789DD_HEALTH_CHECK_STATUS; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_4789DD_HEALTH_CHECK_STATUS" (
    "APPLICATION_NAME" character varying(255),
    "COMPLETE_KEY" character varying(255),
    "DESCRIPTION" text,
    "FAILED_DATE" timestamp without time zone,
    "FAILURE_REASON" text,
    "ID" integer NOT NULL,
    "IS_HEALTHY" boolean,
    "IS_RESOLVED" boolean,
    "RESOLVED_DATE" timestamp without time zone,
    "SEVERITY" character varying(255),
    "STATUS_NAME" character varying(255) NOT NULL
);


ALTER TABLE public."AO_4789DD_HEALTH_CHECK_STATUS" OWNER TO bitbucketuser;

--
-- Name: AO_4789DD_HEALTH_CHECK_STATUS_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_4789DD_HEALTH_CHECK_STATUS_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_4789DD_HEALTH_CHECK_STATUS_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_4789DD_HEALTH_CHECK_STATUS_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_4789DD_HEALTH_CHECK_STATUS_ID_seq" OWNED BY public."AO_4789DD_HEALTH_CHECK_STATUS"."ID";


--
-- Name: AO_4789DD_PROPERTIES; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_4789DD_PROPERTIES" (
    "ID" integer NOT NULL,
    "PROPERTY_NAME" character varying(255) NOT NULL,
    "PROPERTY_VALUE" character varying(255) NOT NULL
);


ALTER TABLE public."AO_4789DD_PROPERTIES" OWNER TO bitbucketuser;

--
-- Name: AO_4789DD_PROPERTIES_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_4789DD_PROPERTIES_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_4789DD_PROPERTIES_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_4789DD_PROPERTIES_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_4789DD_PROPERTIES_ID_seq" OWNED BY public."AO_4789DD_PROPERTIES"."ID";


--
-- Name: AO_4789DD_READ_NOTIFICATIONS; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_4789DD_READ_NOTIFICATIONS" (
    "ID" integer NOT NULL,
    "IS_SNOOZED" boolean,
    "NOTIFICATION_ID" integer NOT NULL,
    "SNOOZE_COUNT" integer,
    "SNOOZE_DATE" timestamp without time zone,
    "USER_KEY" character varying(255) NOT NULL
);


ALTER TABLE public."AO_4789DD_READ_NOTIFICATIONS" OWNER TO bitbucketuser;

--
-- Name: AO_4789DD_READ_NOTIFICATIONS_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_4789DD_READ_NOTIFICATIONS_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_4789DD_READ_NOTIFICATIONS_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_4789DD_READ_NOTIFICATIONS_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_4789DD_READ_NOTIFICATIONS_ID_seq" OWNED BY public."AO_4789DD_READ_NOTIFICATIONS"."ID";


--
-- Name: AO_4789DD_TASK_MONITOR; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_4789DD_TASK_MONITOR" (
    "CLUSTERED_TASK_ID" character varying(255),
    "CREATED_TIMESTAMP" bigint DEFAULT 0,
    "ID" integer NOT NULL,
    "NODE_ID" character varying(255),
    "PROGRESS_MESSAGE" text,
    "PROGRESS_PERCENTAGE" integer,
    "SERIALIZED_ERRORS" text,
    "SERIALIZED_WARNINGS" text,
    "TASK_ID" character varying(255) NOT NULL,
    "TASK_MONITOR_KIND" character varying(255),
    "TASK_STATUS" text
);


ALTER TABLE public."AO_4789DD_TASK_MONITOR" OWNER TO bitbucketuser;

--
-- Name: AO_4789DD_TASK_MONITOR_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_4789DD_TASK_MONITOR_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_4789DD_TASK_MONITOR_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_4789DD_TASK_MONITOR_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_4789DD_TASK_MONITOR_ID_seq" OWNED BY public."AO_4789DD_TASK_MONITOR"."ID";


--
-- Name: AO_616D7B_BRANCH_MODEL; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_616D7B_BRANCH_MODEL" (
    "DEV_ID" character varying(450),
    "DEV_USE_DEFAULT" boolean,
    "IS_ENABLED" boolean,
    "PROD_ID" character varying(450),
    "PROD_USE_DEFAULT" boolean,
    "REPOSITORY_ID" integer NOT NULL
);


ALTER TABLE public."AO_616D7B_BRANCH_MODEL" OWNER TO bitbucketuser;

--
-- Name: AO_616D7B_BRANCH_MODEL_CONFIG; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_616D7B_BRANCH_MODEL_CONFIG" (
    "DEV_ID" character varying(450),
    "DEV_USE_DEFAULT" boolean NOT NULL,
    "ID" integer NOT NULL,
    "PROD_ID" character varying(450),
    "PROD_USE_DEFAULT" boolean NOT NULL,
    "RESOURCE_ID" integer DEFAULT 0 NOT NULL,
    "SCOPE_TYPE" character varying(255) NOT NULL
);


ALTER TABLE public."AO_616D7B_BRANCH_MODEL_CONFIG" OWNER TO bitbucketuser;

--
-- Name: AO_616D7B_BRANCH_MODEL_CONFIG_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_616D7B_BRANCH_MODEL_CONFIG_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_616D7B_BRANCH_MODEL_CONFIG_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_616D7B_BRANCH_MODEL_CONFIG_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_616D7B_BRANCH_MODEL_CONFIG_ID_seq" OWNED BY public."AO_616D7B_BRANCH_MODEL_CONFIG"."ID";


--
-- Name: AO_616D7B_BRANCH_TYPE; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_616D7B_BRANCH_TYPE" (
    "FK_BM_ID" integer,
    "ID" integer NOT NULL,
    "IS_ENABLED" boolean,
    "PREFIX" character varying(450),
    "TYPE_ID" character varying(450)
);


ALTER TABLE public."AO_616D7B_BRANCH_TYPE" OWNER TO bitbucketuser;

--
-- Name: AO_616D7B_BRANCH_TYPE_CONFIG; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_616D7B_BRANCH_TYPE_CONFIG" (
    "BM_ID" integer NOT NULL,
    "ID" integer NOT NULL,
    "IS_ENABLED" boolean NOT NULL,
    "PREFIX" character varying(450),
    "TYPE_ID" character varying(450) NOT NULL
);


ALTER TABLE public."AO_616D7B_BRANCH_TYPE_CONFIG" OWNER TO bitbucketuser;

--
-- Name: AO_616D7B_BRANCH_TYPE_CONFIG_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_616D7B_BRANCH_TYPE_CONFIG_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_616D7B_BRANCH_TYPE_CONFIG_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_616D7B_BRANCH_TYPE_CONFIG_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_616D7B_BRANCH_TYPE_CONFIG_ID_seq" OWNED BY public."AO_616D7B_BRANCH_TYPE_CONFIG"."ID";


--
-- Name: AO_616D7B_BRANCH_TYPE_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_616D7B_BRANCH_TYPE_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_616D7B_BRANCH_TYPE_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_616D7B_BRANCH_TYPE_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_616D7B_BRANCH_TYPE_ID_seq" OWNED BY public."AO_616D7B_BRANCH_TYPE"."ID";


--
-- Name: AO_616D7B_SCOPE_AUTO_MERGE; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_616D7B_SCOPE_AUTO_MERGE" (
    "ENABLED" boolean NOT NULL,
    "ID" integer NOT NULL,
    "MERGE_CHECK_ENABLED" boolean DEFAULT false NOT NULL,
    "RESOURCE_ID" integer NOT NULL,
    "SCOPE_TYPE" character varying(255) NOT NULL
);


ALTER TABLE public."AO_616D7B_SCOPE_AUTO_MERGE" OWNER TO bitbucketuser;

--
-- Name: AO_616D7B_SCOPE_AUTO_MERGE_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_616D7B_SCOPE_AUTO_MERGE_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_616D7B_SCOPE_AUTO_MERGE_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_616D7B_SCOPE_AUTO_MERGE_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_616D7B_SCOPE_AUTO_MERGE_ID_seq" OWNED BY public."AO_616D7B_SCOPE_AUTO_MERGE"."ID";


--
-- Name: AO_6978BB_PERMITTED_ENTITY; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_6978BB_PERMITTED_ENTITY" (
    "ACCESS_KEY_ID" integer,
    "ENTITY_ID" integer NOT NULL,
    "FK_RESTRICTED_ID" integer NOT NULL,
    "GROUP_ID" character varying(255),
    "USER_ID" integer
);


ALTER TABLE public."AO_6978BB_PERMITTED_ENTITY" OWNER TO bitbucketuser;

--
-- Name: AO_6978BB_PERMITTED_ENTITY_ENTITY_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_6978BB_PERMITTED_ENTITY_ENTITY_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_6978BB_PERMITTED_ENTITY_ENTITY_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_6978BB_PERMITTED_ENTITY_ENTITY_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_6978BB_PERMITTED_ENTITY_ENTITY_ID_seq" OWNED BY public."AO_6978BB_PERMITTED_ENTITY"."ENTITY_ID";


--
-- Name: AO_6978BB_RESTRICTED_REF; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_6978BB_RESTRICTED_REF" (
    "REF_ID" integer NOT NULL,
    "REF_TYPE" character varying(255) NOT NULL,
    "REF_VALUE" character varying(255) NOT NULL,
    "RESOURCE_ID" integer NOT NULL,
    "RESTRICTION_TYPE" character varying(255) NOT NULL,
    "SCOPE_TYPE" character varying(255) NOT NULL
);


ALTER TABLE public."AO_6978BB_RESTRICTED_REF" OWNER TO bitbucketuser;

--
-- Name: AO_6978BB_RESTRICTED_REF_REF_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_6978BB_RESTRICTED_REF_REF_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_6978BB_RESTRICTED_REF_REF_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_6978BB_RESTRICTED_REF_REF_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_6978BB_RESTRICTED_REF_REF_ID_seq" OWNED BY public."AO_6978BB_RESTRICTED_REF"."REF_ID";


--
-- Name: AO_777666_JIRA_INDEX; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_777666_JIRA_INDEX" (
    "BRANCH" character varying(255) NOT NULL,
    "ID" bigint NOT NULL,
    "ISSUE" character varying(255) NOT NULL,
    "LAST_UPDATED" timestamp without time zone NOT NULL,
    "PR_ID" bigint,
    "PR_STATE" character varying(255),
    "REPOSITORY" integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."AO_777666_JIRA_INDEX" OWNER TO bitbucketuser;

--
-- Name: AO_777666_JIRA_INDEX_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_777666_JIRA_INDEX_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_777666_JIRA_INDEX_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_777666_JIRA_INDEX_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_777666_JIRA_INDEX_ID_seq" OWNED BY public."AO_777666_JIRA_INDEX"."ID";


--
-- Name: AO_777666_UPDATED_ISSUES; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_777666_UPDATED_ISSUES" (
    "ISSUE" character varying(255) NOT NULL,
    "UPDATE_TIME" bigint DEFAULT 0 NOT NULL
);


ALTER TABLE public."AO_777666_UPDATED_ISSUES" OWNER TO bitbucketuser;

--
-- Name: AO_811463_GIT_LFS_LOCK; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_811463_GIT_LFS_LOCK" (
    "DIRECTORY_HASH" character varying(64) NOT NULL,
    "ID" integer NOT NULL,
    "LOCKED_AT" timestamp without time zone NOT NULL,
    "OWNER_ID" integer NOT NULL,
    "PATH" text NOT NULL,
    "REPOSITORY_ID" integer NOT NULL,
    "REPO_PATH_HASH" character varying(75) NOT NULL
);


ALTER TABLE public."AO_811463_GIT_LFS_LOCK" OWNER TO bitbucketuser;

--
-- Name: AO_811463_GIT_LFS_LOCK_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_811463_GIT_LFS_LOCK_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_811463_GIT_LFS_LOCK_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_811463_GIT_LFS_LOCK_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_811463_GIT_LFS_LOCK_ID_seq" OWNED BY public."AO_811463_GIT_LFS_LOCK"."ID";


--
-- Name: AO_811463_GIT_LFS_REPO_CONFIG; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_811463_GIT_LFS_REPO_CONFIG" (
    "IS_ENABLED" boolean NOT NULL,
    "REPOSITORY_ID" integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."AO_811463_GIT_LFS_REPO_CONFIG" OWNER TO bitbucketuser;

--
-- Name: AO_8E6075_MIRRORING_REQUEST; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_8E6075_MIRRORING_REQUEST" (
    "ADDON_DESCRIPTOR_URI" character varying(450),
    "BASE_URL" character varying(450) NOT NULL,
    "CREATED_DATE" timestamp without time zone NOT NULL,
    "DESCRIPTOR_URL" character varying(450) NOT NULL,
    "ID" integer NOT NULL,
    "MIRROR_ID" character varying(64) NOT NULL,
    "MIRROR_NAME" character varying(64) NOT NULL,
    "MIRROR_TYPE" character varying(255) NOT NULL,
    "PRODUCT_TYPE" character varying(64) NOT NULL,
    "PRODUCT_VERSION" character varying(64) NOT NULL,
    "RESOLVED_DATE" timestamp without time zone,
    "RESOLVER_USER_ID" integer,
    "STATE" character varying(255) NOT NULL
);


ALTER TABLE public."AO_8E6075_MIRRORING_REQUEST" OWNER TO bitbucketuser;

--
-- Name: AO_8E6075_MIRRORING_REQUEST_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_8E6075_MIRRORING_REQUEST_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_8E6075_MIRRORING_REQUEST_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_8E6075_MIRRORING_REQUEST_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_8E6075_MIRRORING_REQUEST_ID_seq" OWNED BY public."AO_8E6075_MIRRORING_REQUEST"."ID";


--
-- Name: AO_8E6075_MIRROR_SERVER; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_8E6075_MIRROR_SERVER" (
    "ADD_ON_KEY" character varying(64) NOT NULL,
    "BASE_URL" character varying(255) NOT NULL,
    "ID" character varying(64) NOT NULL,
    "LAST_SEEN" timestamp without time zone NOT NULL,
    "MIRROR_TYPE" character varying(255) NOT NULL,
    "NAME" character varying(64) NOT NULL,
    "PRODUCT_TYPE" character varying(64) NOT NULL,
    "PRODUCT_VERSION" character varying(64) NOT NULL,
    "STATE" character varying(255) NOT NULL
);


ALTER TABLE public."AO_8E6075_MIRROR_SERVER" OWNER TO bitbucketuser;

--
-- Name: AO_92D5D5_REPO_NOTIFICATION; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_92D5D5_REPO_NOTIFICATION" (
    "ID" integer NOT NULL,
    "PR_NOTIFICATION_SCOPE" integer DEFAULT 0 NOT NULL,
    "PUSH_NOTIFICATION_SCOPE" integer DEFAULT 0 NOT NULL,
    "REPO_ID" integer DEFAULT 0 NOT NULL,
    "USER_ID" integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."AO_92D5D5_REPO_NOTIFICATION" OWNER TO bitbucketuser;

--
-- Name: AO_92D5D5_REPO_NOTIFICATION_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_92D5D5_REPO_NOTIFICATION_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_92D5D5_REPO_NOTIFICATION_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_92D5D5_REPO_NOTIFICATION_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_92D5D5_REPO_NOTIFICATION_ID_seq" OWNED BY public."AO_92D5D5_REPO_NOTIFICATION"."ID";


--
-- Name: AO_92D5D5_USER_NOTIFICATION; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_92D5D5_USER_NOTIFICATION" (
    "BATCH_ID" character varying(255),
    "BATCH_SENDER_ID" character varying(255) NOT NULL,
    "DATA" text NOT NULL,
    "DATE" timestamp without time zone NOT NULL,
    "ID" bigint NOT NULL,
    "USER_ID" integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."AO_92D5D5_USER_NOTIFICATION" OWNER TO bitbucketuser;

--
-- Name: AO_92D5D5_USER_NOTIFICATION_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_92D5D5_USER_NOTIFICATION_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_92D5D5_USER_NOTIFICATION_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_92D5D5_USER_NOTIFICATION_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_92D5D5_USER_NOTIFICATION_ID_seq" OWNED BY public."AO_92D5D5_USER_NOTIFICATION"."ID";


--
-- Name: AO_9DEC2A_DEFAULT_REVIEWER; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_9DEC2A_DEFAULT_REVIEWER" (
    "ENTITY_ID" integer NOT NULL,
    "FK_RESTRICTED_ID" integer NOT NULL,
    "USER_ID" integer
);


ALTER TABLE public."AO_9DEC2A_DEFAULT_REVIEWER" OWNER TO bitbucketuser;

--
-- Name: AO_9DEC2A_DEFAULT_REVIEWER_ENTITY_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_9DEC2A_DEFAULT_REVIEWER_ENTITY_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_9DEC2A_DEFAULT_REVIEWER_ENTITY_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_9DEC2A_DEFAULT_REVIEWER_ENTITY_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_9DEC2A_DEFAULT_REVIEWER_ENTITY_ID_seq" OWNED BY public."AO_9DEC2A_DEFAULT_REVIEWER"."ENTITY_ID";


--
-- Name: AO_9DEC2A_PR_CONDITION; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_9DEC2A_PR_CONDITION" (
    "PR_CONDITION_ID" integer NOT NULL,
    "REQUIRED_APPROVALS" integer DEFAULT 0,
    "RESOURCE_ID" integer NOT NULL,
    "SCOPE_TYPE" character varying(255) NOT NULL,
    "SOURCE_REF_TYPE" character varying(255) NOT NULL,
    "SOURCE_REF_VALUE" character varying(255) NOT NULL,
    "TARGET_REF_TYPE" character varying(255) NOT NULL,
    "TARGET_REF_VALUE" character varying(255) NOT NULL
);


ALTER TABLE public."AO_9DEC2A_PR_CONDITION" OWNER TO bitbucketuser;

--
-- Name: AO_9DEC2A_PR_CONDITION_PR_CONDITION_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_9DEC2A_PR_CONDITION_PR_CONDITION_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_9DEC2A_PR_CONDITION_PR_CONDITION_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_9DEC2A_PR_CONDITION_PR_CONDITION_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_9DEC2A_PR_CONDITION_PR_CONDITION_ID_seq" OWNED BY public."AO_9DEC2A_PR_CONDITION"."PR_CONDITION_ID";


--
-- Name: AO_A0B856_DAILY_COUNTS; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_A0B856_DAILY_COUNTS" (
    "DAY_SINCE_EPOCH" bigint DEFAULT 0 NOT NULL,
    "ERRORS" integer DEFAULT 0 NOT NULL,
    "EVENT_ID" character varying(64) NOT NULL,
    "FAILURES" integer DEFAULT 0 NOT NULL,
    "ID" character varying(88) NOT NULL,
    "SUCCESSES" integer DEFAULT 0 NOT NULL,
    "WEBHOOK_ID" integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."AO_A0B856_DAILY_COUNTS" OWNER TO bitbucketuser;

--
-- Name: AO_A0B856_HIST_INVOCATION; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_A0B856_HIST_INVOCATION" (
    "ERROR_CONTENT" text,
    "EVENT_ID" character varying(64) NOT NULL,
    "FINISH" bigint DEFAULT 0 NOT NULL,
    "ID" character varying(77) NOT NULL,
    "OUTCOME" character varying(255) NOT NULL,
    "REQUEST_BODY" text,
    "REQUEST_HEADERS" text,
    "REQUEST_ID" character varying(64) NOT NULL,
    "REQUEST_METHOD" character varying(16) NOT NULL,
    "REQUEST_URL" character varying(255) NOT NULL,
    "RESPONSE_BODY" text,
    "RESPONSE_HEADERS" text,
    "RESULT_DESCRIPTION" character varying(255) NOT NULL,
    "START" bigint DEFAULT 0 NOT NULL,
    "STATUS_CODE" integer,
    "WEBHOOK_ID" integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."AO_A0B856_HIST_INVOCATION" OWNER TO bitbucketuser;

--
-- Name: AO_A0B856_WEBHOOK; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_A0B856_WEBHOOK" (
    "ACTIVE" boolean,
    "CREATED" timestamp without time zone NOT NULL,
    "ID" integer NOT NULL,
    "NAME" character varying(255) NOT NULL,
    "SCOPE_ID" character varying(255),
    "SCOPE_TYPE" character varying(255) NOT NULL,
    "UPDATED" timestamp without time zone NOT NULL,
    "URL" text NOT NULL
);


ALTER TABLE public."AO_A0B856_WEBHOOK" OWNER TO bitbucketuser;

--
-- Name: AO_A0B856_WEBHOOK_CONFIG; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_A0B856_WEBHOOK_CONFIG" (
    "ID" integer NOT NULL,
    "KEY" character varying(255) NOT NULL,
    "VALUE" character varying(255) NOT NULL,
    "WEBHOOKID" integer NOT NULL
);


ALTER TABLE public."AO_A0B856_WEBHOOK_CONFIG" OWNER TO bitbucketuser;

--
-- Name: AO_A0B856_WEBHOOK_CONFIG_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_A0B856_WEBHOOK_CONFIG_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_A0B856_WEBHOOK_CONFIG_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_A0B856_WEBHOOK_CONFIG_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_A0B856_WEBHOOK_CONFIG_ID_seq" OWNED BY public."AO_A0B856_WEBHOOK_CONFIG"."ID";


--
-- Name: AO_A0B856_WEBHOOK_EVENT; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_A0B856_WEBHOOK_EVENT" (
    "EVENT_ID" character varying(255) NOT NULL,
    "ID" integer NOT NULL,
    "WEBHOOKID" integer NOT NULL
);


ALTER TABLE public."AO_A0B856_WEBHOOK_EVENT" OWNER TO bitbucketuser;

--
-- Name: AO_A0B856_WEBHOOK_EVENT_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_A0B856_WEBHOOK_EVENT_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_A0B856_WEBHOOK_EVENT_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_A0B856_WEBHOOK_EVENT_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_A0B856_WEBHOOK_EVENT_ID_seq" OWNED BY public."AO_A0B856_WEBHOOK_EVENT"."ID";


--
-- Name: AO_A0B856_WEBHOOK_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_A0B856_WEBHOOK_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_A0B856_WEBHOOK_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_A0B856_WEBHOOK_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_A0B856_WEBHOOK_ID_seq" OWNED BY public."AO_A0B856_WEBHOOK"."ID";


--
-- Name: AO_A0B856_WEB_HOOK_LISTENER_AO; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_A0B856_WEB_HOOK_LISTENER_AO" (
    "DESCRIPTION" text,
    "ENABLED" boolean,
    "EVENTS" text,
    "EXCLUDE_BODY" boolean,
    "FILTERS" text,
    "ID" integer NOT NULL,
    "LAST_UPDATED" timestamp without time zone,
    "LAST_UPDATED_USER" character varying(255),
    "NAME" text,
    "PARAMETERS" text,
    "REGISTRATION_METHOD" character varying(255),
    "URL" text
);


ALTER TABLE public."AO_A0B856_WEB_HOOK_LISTENER_AO" OWNER TO bitbucketuser;

--
-- Name: AO_A0B856_WEB_HOOK_LISTENER_AO_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_A0B856_WEB_HOOK_LISTENER_AO_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_A0B856_WEB_HOOK_LISTENER_AO_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_A0B856_WEB_HOOK_LISTENER_AO_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_A0B856_WEB_HOOK_LISTENER_AO_ID_seq" OWNED BY public."AO_A0B856_WEB_HOOK_LISTENER_AO"."ID";


--
-- Name: AO_B586BC_GPG_KEY; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_B586BC_GPG_KEY" (
    "EMAIL" character varying(255),
    "EXPIRY_DATE" timestamp without time zone,
    "FINGERPRINT" character varying(255) NOT NULL,
    "KEY_ID" bigint DEFAULT 0 NOT NULL,
    "KEY_TEXT" text,
    "USER_ID" integer
);


ALTER TABLE public."AO_B586BC_GPG_KEY" OWNER TO bitbucketuser;

--
-- Name: AO_B586BC_GPG_SUB_KEY; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_B586BC_GPG_SUB_KEY" (
    "EXPIRY_DATE" timestamp without time zone,
    "FINGERPRINT" character varying(255) NOT NULL,
    "FK_GPG_KEY_ID" character varying(255) NOT NULL,
    "KEY_ID" bigint DEFAULT 0 NOT NULL
);


ALTER TABLE public."AO_B586BC_GPG_SUB_KEY" OWNER TO bitbucketuser;

--
-- Name: AO_BD73C3_PROJECT_AUDIT; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_BD73C3_PROJECT_AUDIT" (
    "ACTION" character varying(255) NOT NULL,
    "AUDIT_ITEM_ID" integer NOT NULL,
    "DATE" timestamp without time zone NOT NULL,
    "DETAILS" text,
    "PROJECT_ID" integer NOT NULL,
    "USER" integer
);


ALTER TABLE public."AO_BD73C3_PROJECT_AUDIT" OWNER TO bitbucketuser;

--
-- Name: AO_BD73C3_PROJECT_AUDIT_AUDIT_ITEM_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_BD73C3_PROJECT_AUDIT_AUDIT_ITEM_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_BD73C3_PROJECT_AUDIT_AUDIT_ITEM_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_BD73C3_PROJECT_AUDIT_AUDIT_ITEM_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_BD73C3_PROJECT_AUDIT_AUDIT_ITEM_ID_seq" OWNED BY public."AO_BD73C3_PROJECT_AUDIT"."AUDIT_ITEM_ID";


--
-- Name: AO_BD73C3_REPOSITORY_AUDIT; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_BD73C3_REPOSITORY_AUDIT" (
    "ACTION" character varying(255) NOT NULL,
    "AUDIT_ITEM_ID" integer NOT NULL,
    "DATE" timestamp without time zone NOT NULL,
    "DETAILS" text,
    "REPOSITORY_ID" integer NOT NULL,
    "USER" integer
);


ALTER TABLE public."AO_BD73C3_REPOSITORY_AUDIT" OWNER TO bitbucketuser;

--
-- Name: AO_BD73C3_REPOSITORY_AUDIT_AUDIT_ITEM_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_BD73C3_REPOSITORY_AUDIT_AUDIT_ITEM_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_BD73C3_REPOSITORY_AUDIT_AUDIT_ITEM_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_BD73C3_REPOSITORY_AUDIT_AUDIT_ITEM_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_BD73C3_REPOSITORY_AUDIT_AUDIT_ITEM_ID_seq" OWNED BY public."AO_BD73C3_REPOSITORY_AUDIT"."AUDIT_ITEM_ID";


--
-- Name: AO_C77861_AUDIT_ENTITY; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_C77861_AUDIT_ENTITY" (
    "ACTION" character varying(255) NOT NULL,
    "AREA" character varying(255) NOT NULL,
    "ATTRIBUTES" text,
    "CATEGORY" character varying(255),
    "CHANGE_VALUES" text,
    "ENTITY_TIMESTAMP" bigint NOT NULL,
    "ID" bigint NOT NULL,
    "LEVEL" character varying(255) NOT NULL,
    "METHOD" character varying(255),
    "NODE" character varying(255),
    "PRIMARY_RESOURCE_ID" character varying(255),
    "PRIMARY_RESOURCE_TYPE" character varying(255),
    "RESOURCES" text,
    "SEARCH_STRING" text,
    "SECONDARY_RESOURCE_ID" character varying(255),
    "SECONDARY_RESOURCE_TYPE" character varying(255),
    "SOURCE" character varying(255),
    "SYSTEM_INFO" character varying(255),
    "USER_ID" character varying(255),
    "USER_NAME" character varying(255),
    "USER_TYPE" character varying(255)
);


ALTER TABLE public."AO_C77861_AUDIT_ENTITY" OWNER TO bitbucketuser;

--
-- Name: AO_C77861_AUDIT_ENTITY_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_C77861_AUDIT_ENTITY_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_C77861_AUDIT_ENTITY_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_C77861_AUDIT_ENTITY_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_C77861_AUDIT_ENTITY_ID_seq" OWNED BY public."AO_C77861_AUDIT_ENTITY"."ID";


--
-- Name: AO_CFE8FA_BUILD_STATUS; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_CFE8FA_BUILD_STATUS" (
    "CSID" character varying(40) NOT NULL,
    "DATE_ADDED" timestamp without time zone NOT NULL,
    "DESCRIPTION" character varying(255),
    "ID" integer NOT NULL,
    "KEY" character varying(255) NOT NULL,
    "NAME" character varying(255),
    "STATE" character varying(255) NOT NULL,
    "URL" character varying(450) NOT NULL
);


ALTER TABLE public."AO_CFE8FA_BUILD_STATUS" OWNER TO bitbucketuser;

--
-- Name: AO_CFE8FA_BUILD_STATUS_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_CFE8FA_BUILD_STATUS_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_CFE8FA_BUILD_STATUS_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_CFE8FA_BUILD_STATUS_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_CFE8FA_BUILD_STATUS_ID_seq" OWNED BY public."AO_CFE8FA_BUILD_STATUS"."ID";


--
-- Name: AO_D6A508_IMPORT_JOB; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_D6A508_IMPORT_JOB" (
    "CREATED_DATE" timestamp without time zone NOT NULL,
    "ID" bigint NOT NULL,
    "USER_ID" integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."AO_D6A508_IMPORT_JOB" OWNER TO bitbucketuser;

--
-- Name: AO_D6A508_IMPORT_JOB_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_D6A508_IMPORT_JOB_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_D6A508_IMPORT_JOB_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_D6A508_IMPORT_JOB_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_D6A508_IMPORT_JOB_ID_seq" OWNED BY public."AO_D6A508_IMPORT_JOB"."ID";


--
-- Name: AO_D6A508_REPO_IMPORT_TASK; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_D6A508_REPO_IMPORT_TASK" (
    "CLONE_URL" character varying(450) NOT NULL,
    "CREATED_DATE" timestamp without time zone NOT NULL,
    "EXTERNAL_REPO_NAME" character varying(450) NOT NULL,
    "FAILURE_TYPE" integer DEFAULT 0 NOT NULL,
    "ID" bigint NOT NULL,
    "IMPORT_JOB_ID" bigint DEFAULT 0 NOT NULL,
    "LAST_UPDATED" timestamp without time zone NOT NULL,
    "REPOSITORY_ID" integer DEFAULT 0 NOT NULL,
    "STATE" integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."AO_D6A508_REPO_IMPORT_TASK" OWNER TO bitbucketuser;

--
-- Name: AO_D6A508_REPO_IMPORT_TASK_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_D6A508_REPO_IMPORT_TASK_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_D6A508_REPO_IMPORT_TASK_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_D6A508_REPO_IMPORT_TASK_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_D6A508_REPO_IMPORT_TASK_ID_seq" OWNED BY public."AO_D6A508_REPO_IMPORT_TASK"."ID";


--
-- Name: AO_E5A814_ACCESS_TOKEN; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_E5A814_ACCESS_TOKEN" (
    "CREATED_DATE" timestamp without time zone NOT NULL,
    "HASHED_TOKEN" character varying(255) NOT NULL,
    "LAST_AUTHENTICATED" timestamp without time zone,
    "NAME" character varying(255) NOT NULL,
    "TOKEN_ID" character varying(255) NOT NULL,
    "USER_ID" integer DEFAULT 0 NOT NULL
);


ALTER TABLE public."AO_E5A814_ACCESS_TOKEN" OWNER TO bitbucketuser;

--
-- Name: AO_E5A814_ACCESS_TOKEN_PERM; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_E5A814_ACCESS_TOKEN_PERM" (
    "FK_ACCESS_TOKEN_ID" character varying(255) NOT NULL,
    "ID" integer NOT NULL,
    "PERMISSION" integer DEFAULT 0
);


ALTER TABLE public."AO_E5A814_ACCESS_TOKEN_PERM" OWNER TO bitbucketuser;

--
-- Name: AO_E5A814_ACCESS_TOKEN_PERM_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_E5A814_ACCESS_TOKEN_PERM_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_E5A814_ACCESS_TOKEN_PERM_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_E5A814_ACCESS_TOKEN_PERM_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_E5A814_ACCESS_TOKEN_PERM_ID_seq" OWNED BY public."AO_E5A814_ACCESS_TOKEN_PERM"."ID";


--
-- Name: AO_ED669C_SEEN_ASSERTIONS; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_ED669C_SEEN_ASSERTIONS" (
    "ASSERTION_ID" character varying(255) NOT NULL,
    "EXPIRY_TIMESTAMP" bigint DEFAULT 0 NOT NULL,
    "ID" integer NOT NULL
);


ALTER TABLE public."AO_ED669C_SEEN_ASSERTIONS" OWNER TO bitbucketuser;

--
-- Name: AO_ED669C_SEEN_ASSERTIONS_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_ED669C_SEEN_ASSERTIONS_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_ED669C_SEEN_ASSERTIONS_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_ED669C_SEEN_ASSERTIONS_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_ED669C_SEEN_ASSERTIONS_ID_seq" OWNED BY public."AO_ED669C_SEEN_ASSERTIONS"."ID";


--
-- Name: AO_F4ED3A_ADD_ON_PROPERTY_AO; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_F4ED3A_ADD_ON_PROPERTY_AO" (
    "ID" integer NOT NULL,
    "PLUGIN_KEY" character varying(80) NOT NULL,
    "PRIMARY_KEY" character varying(208) NOT NULL,
    "PROPERTY_KEY" character varying(127) NOT NULL,
    "VALUE" text NOT NULL
);


ALTER TABLE public."AO_F4ED3A_ADD_ON_PROPERTY_AO" OWNER TO bitbucketuser;

--
-- Name: AO_F4ED3A_ADD_ON_PROPERTY_AO_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_F4ED3A_ADD_ON_PROPERTY_AO_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_F4ED3A_ADD_ON_PROPERTY_AO_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_F4ED3A_ADD_ON_PROPERTY_AO_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_F4ED3A_ADD_ON_PROPERTY_AO_ID_seq" OWNED BY public."AO_F4ED3A_ADD_ON_PROPERTY_AO"."ID";


--
-- Name: AO_FB71B4_SSH_PUBLIC_KEY; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public."AO_FB71B4_SSH_PUBLIC_KEY" (
    "ENTITY_ID" integer NOT NULL,
    "KEY_MD5" character varying(32) NOT NULL,
    "KEY_TEXT" text NOT NULL,
    "KEY_TYPE" character varying(255) NOT NULL,
    "LABEL" character varying(255),
    "LABEL_LOWER" character varying(255),
    "USER_ID" integer NOT NULL
);


ALTER TABLE public."AO_FB71B4_SSH_PUBLIC_KEY" OWNER TO bitbucketuser;

--
-- Name: AO_FB71B4_SSH_PUBLIC_KEY_ENTITY_ID_seq; Type: SEQUENCE; Schema: public; Owner: bitbucketuser
--

CREATE SEQUENCE public."AO_FB71B4_SSH_PUBLIC_KEY_ENTITY_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."AO_FB71B4_SSH_PUBLIC_KEY_ENTITY_ID_seq" OWNER TO bitbucketuser;

--
-- Name: AO_FB71B4_SSH_PUBLIC_KEY_ENTITY_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bitbucketuser
--

ALTER SEQUENCE public."AO_FB71B4_SSH_PUBLIC_KEY_ENTITY_ID_seq" OWNED BY public."AO_FB71B4_SSH_PUBLIC_KEY"."ENTITY_ID";


--
-- Name: app_property; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.app_property (
    prop_key character varying(50) NOT NULL,
    prop_value character varying(2000) NOT NULL
);


ALTER TABLE public.app_property OWNER TO bitbucketuser;

--
-- Name: bb_alert; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_alert (
    details_json text,
    id bigint NOT NULL,
    issue_component_id character varying(255) NOT NULL,
    issue_id character varying(255) NOT NULL,
    issue_severity integer NOT NULL,
    node_name character varying(255) NOT NULL,
    node_name_lower character varying(255) NOT NULL,
    "timestamp" bigint NOT NULL,
    trigger_module character varying(1024),
    trigger_plugin_key character varying(255),
    trigger_plugin_key_lower character varying(255),
    trigger_plugin_version character varying(255)
);


ALTER TABLE public.bb_alert OWNER TO bitbucketuser;

--
-- Name: COLUMN bb_alert.details_json; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_alert.details_json IS 'text';


--
-- Name: bb_announcement_banner; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_announcement_banner (
    id bigint NOT NULL,
    enabled boolean NOT NULL,
    audience integer NOT NULL,
    message character varying(4000) NOT NULL
);


ALTER TABLE public.bb_announcement_banner OWNER TO bitbucketuser;

--
-- Name: bb_attachment; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_attachment (
    id bigint NOT NULL,
    repository_id integer,
    filename character varying(255) NOT NULL
);


ALTER TABLE public.bb_attachment OWNER TO bitbucketuser;

--
-- Name: bb_attachment_metadata; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_attachment_metadata (
    attachment_id bigint NOT NULL,
    metadata text
);


ALTER TABLE public.bb_attachment_metadata OWNER TO bitbucketuser;

--
-- Name: bb_clusteredjob; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_clusteredjob (
    job_id character varying(255) NOT NULL,
    job_runner_key character varying(255) NOT NULL,
    sched_type integer NOT NULL,
    interval_millis bigint,
    first_run timestamp without time zone,
    cron_expression character varying(255),
    time_zone character varying(64),
    next_run timestamp without time zone,
    version bigint,
    parameters bytea
);


ALTER TABLE public.bb_clusteredjob OWNER TO bitbucketuser;

--
-- Name: bb_cmt_disc_comment_activity; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_cmt_disc_comment_activity (
    activity_id bigint NOT NULL,
    comment_id bigint NOT NULL,
    comment_action integer NOT NULL
);


ALTER TABLE public.bb_cmt_disc_comment_activity OWNER TO bitbucketuser;

--
-- Name: bb_comment; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_comment (
    id bigint NOT NULL,
    author_id integer NOT NULL,
    comment_text text NOT NULL,
    created_timestamp timestamp without time zone NOT NULL,
    entity_version integer NOT NULL,
    thread_id bigint,
    updated_timestamp timestamp without time zone NOT NULL,
    resolved_timestamp timestamp without time zone,
    resolver_id integer,
    severity integer NOT NULL,
    state integer NOT NULL
);


ALTER TABLE public.bb_comment OWNER TO bitbucketuser;

--
-- Name: bb_comment_parent; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_comment_parent (
    comment_id bigint NOT NULL,
    parent_id bigint NOT NULL
);


ALTER TABLE public.bb_comment_parent OWNER TO bitbucketuser;

--
-- Name: bb_comment_thread; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_comment_thread (
    id bigint NOT NULL,
    commentable_id bigint NOT NULL,
    commentable_type integer NOT NULL,
    created_timestamp timestamp without time zone NOT NULL,
    entity_version integer NOT NULL,
    updated_timestamp timestamp without time zone NOT NULL,
    diff_type integer,
    file_type integer,
    from_hash character varying(40),
    from_path character varying(1024),
    is_orphaned boolean,
    line_number integer,
    line_type integer,
    to_hash character varying(40),
    to_path character varying(1024)
);


ALTER TABLE public.bb_comment_thread OWNER TO bitbucketuser;

--
-- Name: bb_data_store; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_data_store (
    id bigint NOT NULL,
    ds_path character varying(128) NOT NULL,
    ds_uuid character varying(40) NOT NULL
);


ALTER TABLE public.bb_data_store OWNER TO bitbucketuser;

--
-- Name: bb_git_pr_cached_merge; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_git_pr_cached_merge (
    id bigint NOT NULL,
    from_hash character varying(40) NOT NULL,
    to_hash character varying(40) NOT NULL,
    merge_type integer NOT NULL
);


ALTER TABLE public.bb_git_pr_cached_merge OWNER TO bitbucketuser;

--
-- Name: bb_git_pr_common_ancestor; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_git_pr_common_ancestor (
    id bigint NOT NULL,
    from_hash character varying(40) NOT NULL,
    to_hash character varying(40) NOT NULL,
    ancestor_hash character varying(40) NOT NULL
);


ALTER TABLE public.bb_git_pr_common_ancestor OWNER TO bitbucketuser;

--
-- Name: bb_hook_script; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_hook_script (
    id bigint NOT NULL,
    hook_version integer NOT NULL,
    hook_size integer NOT NULL,
    hook_type integer NOT NULL,
    created_timestamp timestamp without time zone NOT NULL,
    updated_timestamp timestamp without time zone NOT NULL,
    hook_hash character varying(64) NOT NULL,
    hook_name character varying(255) NOT NULL,
    plugin_key character varying(255) NOT NULL,
    hook_description character varying(255)
);


ALTER TABLE public.bb_hook_script OWNER TO bitbucketuser;

--
-- Name: bb_hook_script_config; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_hook_script_config (
    id bigint NOT NULL,
    script_id bigint NOT NULL,
    scope_id integer,
    scope_type integer
);


ALTER TABLE public.bb_hook_script_config OWNER TO bitbucketuser;

--
-- Name: bb_hook_script_trigger; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_hook_script_trigger (
    config_id bigint NOT NULL,
    trigger_id character varying(255) NOT NULL
);


ALTER TABLE public.bb_hook_script_trigger OWNER TO bitbucketuser;

--
-- Name: bb_integrity_event; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_integrity_event (
    event_key character varying(255) NOT NULL,
    event_timestamp timestamp without time zone NOT NULL,
    event_node character varying(255) DEFAULT '00000000-0000-0000-0000-000000000000'::character varying NOT NULL
);


ALTER TABLE public.bb_integrity_event OWNER TO bitbucketuser;

--
-- Name: bb_job; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_job (
    end_timestamp timestamp without time zone,
    id bigint NOT NULL,
    initiator_id integer,
    node_id character varying(64) NOT NULL,
    progress_percentage integer DEFAULT 0 NOT NULL,
    progress_message text,
    start_timestamp timestamp without time zone NOT NULL,
    state integer NOT NULL,
    type character varying(255) NOT NULL,
    updated_timestamp timestamp without time zone NOT NULL,
    entity_version integer NOT NULL
);


ALTER TABLE public.bb_job OWNER TO bitbucketuser;

--
-- Name: COLUMN bb_job.end_timestamp; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job.end_timestamp IS 'endDate';


--
-- Name: COLUMN bb_job.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job.id IS 'id';


--
-- Name: COLUMN bb_job.initiator_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job.initiator_id IS 'owner';


--
-- Name: COLUMN bb_job.node_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job.node_id IS 'nodeId';


--
-- Name: COLUMN bb_job.progress_percentage; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job.progress_percentage IS 'progressPercentage';


--
-- Name: COLUMN bb_job.progress_message; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job.progress_message IS 'progressMessage';


--
-- Name: COLUMN bb_job.start_timestamp; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job.start_timestamp IS 'startDate';


--
-- Name: COLUMN bb_job.state; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job.state IS 'state';


--
-- Name: COLUMN bb_job.type; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job.type IS 'type';


--
-- Name: COLUMN bb_job.updated_timestamp; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job.updated_timestamp IS 'updatedDate';


--
-- Name: COLUMN bb_job.entity_version; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job.entity_version IS 'version';


--
-- Name: bb_job_message; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_job_message (
    created_timestamp timestamp without time zone NOT NULL,
    id bigint NOT NULL,
    job_id bigint NOT NULL,
    severity integer NOT NULL,
    subject character varying(1024),
    text text NOT NULL
);


ALTER TABLE public.bb_job_message OWNER TO bitbucketuser;

--
-- Name: COLUMN bb_job_message.created_timestamp; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job_message.created_timestamp IS 'createdDate';


--
-- Name: COLUMN bb_job_message.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job_message.id IS 'id';


--
-- Name: COLUMN bb_job_message.job_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job_message.job_id IS 'job';


--
-- Name: COLUMN bb_job_message.severity; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job_message.severity IS 'severity';


--
-- Name: COLUMN bb_job_message.subject; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job_message.subject IS 'subject';


--
-- Name: COLUMN bb_job_message.text; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_job_message.text IS 'text';


--
-- Name: bb_label; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_label (
    id bigint NOT NULL,
    name character varying(255) NOT NULL
);


ALTER TABLE public.bb_label OWNER TO bitbucketuser;

--
-- Name: bb_label_mapping; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_label_mapping (
    id bigint NOT NULL,
    label_id bigint NOT NULL,
    labelable_id integer NOT NULL,
    labelable_type integer NOT NULL
);


ALTER TABLE public.bb_label_mapping OWNER TO bitbucketuser;

--
-- Name: bb_mirror_content_hash; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_mirror_content_hash (
    repository_id integer NOT NULL,
    updated_timestamp timestamp without time zone NOT NULL,
    hash character varying(64) NOT NULL
);


ALTER TABLE public.bb_mirror_content_hash OWNER TO bitbucketuser;

--
-- Name: bb_mirror_metadata_hash; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_mirror_metadata_hash (
    repository_id integer NOT NULL,
    updated_timestamp timestamp without time zone NOT NULL,
    hash character varying(64) NOT NULL
);


ALTER TABLE public.bb_mirror_metadata_hash OWNER TO bitbucketuser;

--
-- Name: bb_pr_comment_activity; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_pr_comment_activity (
    activity_id bigint NOT NULL,
    comment_id bigint NOT NULL,
    comment_action integer NOT NULL
);


ALTER TABLE public.bb_pr_comment_activity OWNER TO bitbucketuser;

--
-- Name: bb_pr_commit; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_pr_commit (
    pr_id bigint NOT NULL,
    commit_id character varying(40) NOT NULL
);


ALTER TABLE public.bb_pr_commit OWNER TO bitbucketuser;

--
-- Name: COLUMN bb_pr_commit.pr_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_pr_commit.pr_id IS 'pullRequest';


--
-- Name: bb_pr_part_status_weight; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_pr_part_status_weight (
    status_id integer NOT NULL,
    status_weight integer NOT NULL
);


ALTER TABLE public.bb_pr_part_status_weight OWNER TO bitbucketuser;

--
-- Name: bb_pr_reviewer_added; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_pr_reviewer_added (
    activity_id bigint NOT NULL,
    user_id integer NOT NULL
);


ALTER TABLE public.bb_pr_reviewer_added OWNER TO bitbucketuser;

--
-- Name: COLUMN bb_pr_reviewer_added.activity_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_pr_reviewer_added.activity_id IS 'joinActivityKey';


--
-- Name: COLUMN bb_pr_reviewer_added.user_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_pr_reviewer_added.user_id IS 'joinUserKey';


--
-- Name: bb_pr_reviewer_removed; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_pr_reviewer_removed (
    activity_id bigint NOT NULL,
    user_id integer NOT NULL
);


ALTER TABLE public.bb_pr_reviewer_removed OWNER TO bitbucketuser;

--
-- Name: COLUMN bb_pr_reviewer_removed.activity_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_pr_reviewer_removed.activity_id IS 'joinActivityKey';


--
-- Name: COLUMN bb_pr_reviewer_removed.user_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_pr_reviewer_removed.user_id IS 'joinuserKey';


--
-- Name: bb_pr_reviewer_upd_activity; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_pr_reviewer_upd_activity (
    activity_id bigint NOT NULL
);


ALTER TABLE public.bb_pr_reviewer_upd_activity OWNER TO bitbucketuser;

--
-- Name: COLUMN bb_pr_reviewer_upd_activity.activity_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_pr_reviewer_upd_activity.activity_id IS 'joinPrimaryKey';


--
-- Name: bb_proj_merge_config; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_proj_merge_config (
    id bigint NOT NULL,
    project_id integer NOT NULL,
    scm_id character varying(255) NOT NULL,
    default_strategy_id character varying(255) NOT NULL,
    commit_summaries integer NOT NULL
);


ALTER TABLE public.bb_proj_merge_config OWNER TO bitbucketuser;

--
-- Name: COLUMN bb_proj_merge_config.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_proj_merge_config.id IS 'id';


--
-- Name: COLUMN bb_proj_merge_config.project_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_proj_merge_config.project_id IS 'project.id';


--
-- Name: COLUMN bb_proj_merge_config.scm_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_proj_merge_config.scm_id IS 'scmId';


--
-- Name: COLUMN bb_proj_merge_config.default_strategy_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_proj_merge_config.default_strategy_id IS 'defaultStrategyId';


--
-- Name: bb_proj_merge_strategy; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_proj_merge_strategy (
    config_id bigint NOT NULL,
    strategy_id character varying(255) NOT NULL
);


ALTER TABLE public.bb_proj_merge_strategy OWNER TO bitbucketuser;

--
-- Name: COLUMN bb_proj_merge_strategy.config_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_proj_merge_strategy.config_id IS 'InternalProjectMergeStrategy.id';


--
-- Name: COLUMN bb_proj_merge_strategy.strategy_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_proj_merge_strategy.strategy_id IS 'ScmMergeStrategy.id';


--
-- Name: bb_project_alias; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_project_alias (
    id bigint NOT NULL,
    project_id integer NOT NULL,
    namespace character varying(128) NOT NULL,
    project_key character varying(128) NOT NULL,
    created_timestamp timestamp without time zone NOT NULL
);


ALTER TABLE public.bb_project_alias OWNER TO bitbucketuser;

--
-- Name: bb_repo_merge_config; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_repo_merge_config (
    id bigint NOT NULL,
    repository_id integer NOT NULL,
    default_strategy_id character varying(255) NOT NULL,
    commit_summaries integer NOT NULL
);


ALTER TABLE public.bb_repo_merge_config OWNER TO bitbucketuser;

--
-- Name: COLUMN bb_repo_merge_config.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_repo_merge_config.id IS 'id';


--
-- Name: COLUMN bb_repo_merge_config.repository_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_repo_merge_config.repository_id IS 'repository.id';


--
-- Name: COLUMN bb_repo_merge_config.default_strategy_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_repo_merge_config.default_strategy_id IS 'defaultStrategyId';


--
-- Name: bb_repo_merge_strategy; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_repo_merge_strategy (
    config_id bigint NOT NULL,
    strategy_id character varying(255) NOT NULL
);


ALTER TABLE public.bb_repo_merge_strategy OWNER TO bitbucketuser;

--
-- Name: COLUMN bb_repo_merge_strategy.config_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_repo_merge_strategy.config_id IS 'InternalRepositoryMergeStrategy.id';


--
-- Name: COLUMN bb_repo_merge_strategy.strategy_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_repo_merge_strategy.strategy_id IS 'ScmMergeStrategy.id';


--
-- Name: bb_repository_alias; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_repository_alias (
    id bigint NOT NULL,
    repository_id integer NOT NULL,
    project_namespace character varying(128) NOT NULL,
    project_key character varying(128) NOT NULL,
    slug character varying(128) NOT NULL,
    created_timestamp timestamp without time zone NOT NULL
);


ALTER TABLE public.bb_repository_alias OWNER TO bitbucketuser;

--
-- Name: bb_rl_reject_counter; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_rl_reject_counter (
    id bigint NOT NULL,
    user_id integer NOT NULL,
    interval_start timestamp without time zone NOT NULL,
    reject_count bigint NOT NULL
);


ALTER TABLE public.bb_rl_reject_counter OWNER TO bitbucketuser;

--
-- Name: bb_rl_user_settings; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_rl_user_settings (
    id bigint NOT NULL,
    user_id integer NOT NULL,
    capacity integer NOT NULL,
    fill_rate integer NOT NULL,
    whitelisted boolean NOT NULL
);


ALTER TABLE public.bb_rl_user_settings OWNER TO bitbucketuser;

--
-- Name: bb_scm_merge_config; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_scm_merge_config (
    id bigint NOT NULL,
    scm_id character varying(255) NOT NULL,
    default_strategy_id character varying(255) NOT NULL,
    commit_summaries integer NOT NULL
);


ALTER TABLE public.bb_scm_merge_config OWNER TO bitbucketuser;

--
-- Name: COLUMN bb_scm_merge_config.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_scm_merge_config.id IS 'id';


--
-- Name: COLUMN bb_scm_merge_config.scm_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_scm_merge_config.scm_id IS 'scmId';


--
-- Name: COLUMN bb_scm_merge_config.default_strategy_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_scm_merge_config.default_strategy_id IS 'defaultStrategyId';


--
-- Name: bb_scm_merge_strategy; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_scm_merge_strategy (
    config_id bigint NOT NULL,
    strategy_id character varying(255) NOT NULL
);


ALTER TABLE public.bb_scm_merge_strategy OWNER TO bitbucketuser;

--
-- Name: COLUMN bb_scm_merge_strategy.config_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_scm_merge_strategy.config_id IS 'InternalScmMergeStrategy.id';


--
-- Name: COLUMN bb_scm_merge_strategy.strategy_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_scm_merge_strategy.strategy_id IS 'ScmMergeStrategy.id';


--
-- Name: bb_suggestion_group; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_suggestion_group (
    comment_id bigint NOT NULL,
    state integer NOT NULL,
    applied_index integer
);


ALTER TABLE public.bb_suggestion_group OWNER TO bitbucketuser;

--
-- Name: bb_thread_root_comment; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_thread_root_comment (
    thread_id bigint NOT NULL,
    comment_id bigint NOT NULL
);


ALTER TABLE public.bb_thread_root_comment OWNER TO bitbucketuser;

--
-- Name: bb_user_dark_feature; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.bb_user_dark_feature (
    id bigint NOT NULL,
    user_id integer NOT NULL,
    is_enabled boolean,
    feature_key character varying(255) NOT NULL
);


ALTER TABLE public.bb_user_dark_feature OWNER TO bitbucketuser;

--
-- Name: COLUMN bb_user_dark_feature.is_enabled; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.bb_user_dark_feature.is_enabled IS 'enabled';


--
-- Name: changeset; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.changeset (
    id character varying(40) NOT NULL,
    author_timestamp timestamp without time zone NOT NULL
);


ALTER TABLE public.changeset OWNER TO bitbucketuser;

--
-- Name: cs_attribute; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cs_attribute (
    cs_id character varying(40) NOT NULL,
    att_name character varying(64) NOT NULL,
    att_value character varying(1024) NOT NULL
);


ALTER TABLE public.cs_attribute OWNER TO bitbucketuser;

--
-- Name: cs_indexer_state; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cs_indexer_state (
    indexer_id character varying(128) NOT NULL,
    repository_id integer NOT NULL,
    last_run bigint
);


ALTER TABLE public.cs_indexer_state OWNER TO bitbucketuser;

--
-- Name: cs_repo_membership; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cs_repo_membership (
    cs_id character varying(40) NOT NULL,
    repository_id integer NOT NULL
);


ALTER TABLE public.cs_repo_membership OWNER TO bitbucketuser;

--
-- Name: current_app; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.current_app (
    id integer NOT NULL,
    application_id character varying(255) NOT NULL,
    public_key_base64 character varying(4000) NOT NULL,
    private_key_base64 character varying(4000) NOT NULL
);


ALTER TABLE public.current_app OWNER TO bitbucketuser;

--
-- Name: cwd_app_dir_default_groups; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_app_dir_default_groups (
    id bigint NOT NULL,
    application_mapping_id bigint NOT NULL,
    group_name character varying(255) NOT NULL
);


ALTER TABLE public.cwd_app_dir_default_groups OWNER TO bitbucketuser;

--
-- Name: cwd_app_dir_group_mapping; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_app_dir_group_mapping (
    id bigint NOT NULL,
    app_dir_mapping_id bigint NOT NULL,
    application_id bigint NOT NULL,
    directory_id bigint NOT NULL,
    group_name character varying(255) NOT NULL
);


ALTER TABLE public.cwd_app_dir_group_mapping OWNER TO bitbucketuser;

--
-- Name: cwd_app_dir_mapping; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_app_dir_mapping (
    id bigint NOT NULL,
    application_id bigint NOT NULL,
    directory_id bigint NOT NULL,
    list_index integer,
    is_allow_all character(1) NOT NULL
);


ALTER TABLE public.cwd_app_dir_mapping OWNER TO bitbucketuser;

--
-- Name: cwd_app_dir_operation; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_app_dir_operation (
    app_dir_mapping_id bigint NOT NULL,
    operation_type character varying(32) NOT NULL
);


ALTER TABLE public.cwd_app_dir_operation OWNER TO bitbucketuser;

--
-- Name: cwd_app_licensed_user; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_app_licensed_user (
    id bigint NOT NULL,
    username character varying(255) NOT NULL,
    full_name character varying(255),
    email character varying(255),
    last_active timestamp without time zone,
    directory_id bigint NOT NULL,
    lower_username character varying(255) NOT NULL,
    lower_full_name character varying(255),
    lower_email character varying(255)
);


ALTER TABLE public.cwd_app_licensed_user OWNER TO bitbucketuser;

--
-- Name: cwd_app_licensing; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_app_licensing (
    id bigint NOT NULL,
    generated_on timestamp without time zone NOT NULL,
    version bigint NOT NULL,
    application_id bigint NOT NULL,
    application_subtype character varying(32) NOT NULL,
    total_users integer NOT NULL,
    max_user_limit integer NOT NULL,
    total_crowd_users integer NOT NULL,
    active character(1) NOT NULL
);


ALTER TABLE public.cwd_app_licensing OWNER TO bitbucketuser;

--
-- Name: cwd_app_licensing_dir_info; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_app_licensing_dir_info (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    directory_id bigint,
    licensing_summary_id bigint NOT NULL
);


ALTER TABLE public.cwd_app_licensing_dir_info OWNER TO bitbucketuser;

--
-- Name: cwd_application; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_application (
    id bigint NOT NULL,
    application_name character varying(255) NOT NULL,
    lower_application_name character varying(255) NOT NULL,
    created_date timestamp without time zone NOT NULL,
    updated_date timestamp without time zone NOT NULL,
    description character varying(255),
    application_type character varying(32) NOT NULL,
    credential character varying(255) NOT NULL,
    is_active character(1) NOT NULL
);


ALTER TABLE public.cwd_application OWNER TO bitbucketuser;

--
-- Name: cwd_application_address; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_application_address (
    application_id bigint NOT NULL,
    remote_address character varying(255) NOT NULL
);


ALTER TABLE public.cwd_application_address OWNER TO bitbucketuser;

--
-- Name: cwd_application_alias; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_application_alias (
    id bigint NOT NULL,
    application_id bigint NOT NULL,
    user_name character varying(255) NOT NULL,
    lower_user_name character varying(255) NOT NULL,
    alias_name character varying(255) NOT NULL,
    lower_alias_name character varying(255) NOT NULL
);


ALTER TABLE public.cwd_application_alias OWNER TO bitbucketuser;

--
-- Name: cwd_application_attribute; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_application_attribute (
    application_id bigint NOT NULL,
    attribute_name character varying(255) NOT NULL,
    attribute_value text
);


ALTER TABLE public.cwd_application_attribute OWNER TO bitbucketuser;

--
-- Name: cwd_application_saml_config; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_application_saml_config (
    application_id bigint NOT NULL,
    assertion_consumer_service character varying(255) NOT NULL,
    audience character varying(255) NOT NULL,
    enabled character(1) NOT NULL
);


ALTER TABLE public.cwd_application_saml_config OWNER TO bitbucketuser;

--
-- Name: cwd_directory; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_directory (
    id bigint NOT NULL,
    directory_name character varying(255) NOT NULL,
    lower_directory_name character varying(255) NOT NULL,
    created_date timestamp without time zone NOT NULL,
    updated_date timestamp without time zone NOT NULL,
    description character varying(255),
    impl_class character varying(255) NOT NULL,
    lower_impl_class character varying(255) NOT NULL,
    directory_type character varying(32) NOT NULL,
    is_active character(1) NOT NULL
);


ALTER TABLE public.cwd_directory OWNER TO bitbucketuser;

--
-- Name: cwd_directory_attribute; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_directory_attribute (
    directory_id bigint NOT NULL,
    attribute_name character varying(255) NOT NULL,
    attribute_value text
);


ALTER TABLE public.cwd_directory_attribute OWNER TO bitbucketuser;

--
-- Name: cwd_directory_operation; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_directory_operation (
    directory_id bigint NOT NULL,
    operation_type character varying(32) NOT NULL
);


ALTER TABLE public.cwd_directory_operation OWNER TO bitbucketuser;

--
-- Name: cwd_granted_perm; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_granted_perm (
    id bigint NOT NULL,
    created_date timestamp without time zone NOT NULL,
    permission_id integer NOT NULL,
    group_name character varying(255) NOT NULL,
    app_dir_mapping_id bigint NOT NULL
);


ALTER TABLE public.cwd_granted_perm OWNER TO bitbucketuser;

--
-- Name: cwd_group; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_group (
    id bigint NOT NULL,
    group_name character varying(255) NOT NULL,
    lower_group_name character varying(255) NOT NULL,
    created_date timestamp without time zone NOT NULL,
    updated_date timestamp without time zone NOT NULL,
    description character varying(255),
    group_type character varying(32) NOT NULL,
    directory_id bigint NOT NULL,
    is_active character(1) NOT NULL,
    is_local character(1) NOT NULL,
    external_id character varying(255)
);


ALTER TABLE public.cwd_group OWNER TO bitbucketuser;

--
-- Name: cwd_group_admin_group; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_group_admin_group (
    id bigint NOT NULL,
    group_id bigint NOT NULL,
    target_group_id bigint NOT NULL
);


ALTER TABLE public.cwd_group_admin_group OWNER TO bitbucketuser;

--
-- Name: cwd_group_admin_user; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_group_admin_user (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    target_group_id bigint NOT NULL
);


ALTER TABLE public.cwd_group_admin_user OWNER TO bitbucketuser;

--
-- Name: cwd_group_attribute; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_group_attribute (
    id bigint NOT NULL,
    group_id bigint NOT NULL,
    directory_id bigint NOT NULL,
    attribute_name character varying(255) NOT NULL,
    attribute_value character varying(255),
    attribute_lower_value character varying(255)
);


ALTER TABLE public.cwd_group_attribute OWNER TO bitbucketuser;

--
-- Name: cwd_membership; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_membership (
    id bigint NOT NULL,
    parent_id bigint,
    child_id bigint,
    membership_type character varying(32),
    group_type character varying(32) NOT NULL,
    parent_name character varying(255) NOT NULL,
    lower_parent_name character varying(255) NOT NULL,
    child_name character varying(255) NOT NULL,
    lower_child_name character varying(255) NOT NULL,
    directory_id bigint NOT NULL,
    created_date timestamp without time zone
);


ALTER TABLE public.cwd_membership OWNER TO bitbucketuser;

--
-- Name: cwd_property; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_property (
    property_key character varying(255) NOT NULL,
    property_name character varying(255) NOT NULL,
    property_value text
);


ALTER TABLE public.cwd_property OWNER TO bitbucketuser;

--
-- Name: cwd_tombstone; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_tombstone (
    id bigint NOT NULL,
    tombstone_type character varying(255) NOT NULL,
    tombstone_timestamp bigint NOT NULL,
    application_id bigint,
    directory_id bigint,
    entity_name character varying(255),
    parent character varying(255)
);


ALTER TABLE public.cwd_tombstone OWNER TO bitbucketuser;

--
-- Name: cwd_user; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_user (
    id bigint NOT NULL,
    user_name character varying(255) NOT NULL,
    lower_user_name character varying(255) NOT NULL,
    created_date timestamp without time zone NOT NULL,
    updated_date timestamp without time zone NOT NULL,
    first_name character varying(255),
    lower_first_name character varying(255),
    last_name character varying(255),
    lower_last_name character varying(255),
    display_name character varying(255),
    lower_display_name character varying(255),
    email_address character varying(255),
    lower_email_address character varying(255),
    directory_id bigint NOT NULL,
    credential character varying(255),
    is_active character(1) NOT NULL,
    external_id character varying(255)
);


ALTER TABLE public.cwd_user OWNER TO bitbucketuser;

--
-- Name: COLUMN cwd_user.external_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.cwd_user.external_id IS 'external_id';


--
-- Name: cwd_user_attribute; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_user_attribute (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    directory_id bigint NOT NULL,
    attribute_name character varying(255) NOT NULL,
    attribute_value character varying(255),
    attribute_lower_value character varying(255),
    attribute_numeric_value bigint
);


ALTER TABLE public.cwd_user_attribute OWNER TO bitbucketuser;

--
-- Name: cwd_user_credential_record; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_user_credential_record (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    password_hash character varying(255) NOT NULL,
    list_index integer
);


ALTER TABLE public.cwd_user_credential_record OWNER TO bitbucketuser;

--
-- Name: cwd_webhook; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.cwd_webhook (
    id bigint NOT NULL,
    endpoint_url character varying(255) NOT NULL,
    application_id bigint NOT NULL,
    token character varying(255),
    oldest_failure_date timestamp without time zone,
    failures_since_last_success bigint NOT NULL
);


ALTER TABLE public.cwd_webhook OWNER TO bitbucketuser;

--
-- Name: databasechangelog; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.databasechangelog (
    id character varying(255) NOT NULL,
    author character varying(255) NOT NULL,
    filename character varying(255) NOT NULL,
    dateexecuted timestamp without time zone NOT NULL,
    orderexecuted integer NOT NULL,
    exectype character varying(10) NOT NULL,
    md5sum character varying(35),
    description character varying(255),
    comments character varying(255),
    tag character varying(255),
    liquibase character varying(20),
    contexts character varying(255),
    labels character varying(255),
    deployment_id character varying(10)
);


ALTER TABLE public.databasechangelog OWNER TO bitbucketuser;

--
-- Name: databasechangeloglock; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.databasechangeloglock (
    id integer NOT NULL,
    locked boolean NOT NULL,
    lockgranted timestamp without time zone,
    lockedby character varying(255)
);


ALTER TABLE public.databasechangeloglock OWNER TO bitbucketuser;

--
-- Name: hibernate_unique_key; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.hibernate_unique_key (
    next_hi bigint NOT NULL
);


ALTER TABLE public.hibernate_unique_key OWNER TO bitbucketuser;

--
-- Name: id_sequence; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.id_sequence (
    sequence_name character varying(255) NOT NULL,
    next_val bigint NOT NULL
);


ALTER TABLE public.id_sequence OWNER TO bitbucketuser;

--
-- Name: plugin_setting; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.plugin_setting (
    namespace character varying(255) NOT NULL,
    key_name character varying(255) NOT NULL,
    key_value text NOT NULL,
    id bigint NOT NULL
);


ALTER TABLE public.plugin_setting OWNER TO bitbucketuser;

--
-- Name: COLUMN plugin_setting.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.plugin_setting.id IS 'id';


--
-- Name: plugin_state; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.plugin_state (
    name character varying(255) NOT NULL,
    enabled boolean NOT NULL,
    updated_timestamp bigint NOT NULL
);


ALTER TABLE public.plugin_state OWNER TO bitbucketuser;

--
-- Name: project; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.project (
    id integer NOT NULL,
    name character varying(128) NOT NULL,
    project_key character varying(128) NOT NULL,
    description character varying(255),
    project_type integer NOT NULL,
    namespace character varying(128) NOT NULL
);


ALTER TABLE public.project OWNER TO bitbucketuser;

--
-- Name: COLUMN project.namespace; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.project.namespace IS 'project namespace';


--
-- Name: repository; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.repository (
    id integer NOT NULL,
    slug character varying(128) NOT NULL,
    name character varying(128) NOT NULL,
    state integer NOT NULL,
    project_id integer NOT NULL,
    scm_id character varying(255) NOT NULL,
    hierarchy_id character varying(20) NOT NULL,
    is_forkable boolean NOT NULL,
    is_public boolean NOT NULL,
    store_id bigint,
    description character varying(255)
);


ALTER TABLE public.repository OWNER TO bitbucketuser;

--
-- Name: COLUMN repository.scm_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.repository.scm_id IS 'scmId';


--
-- Name: COLUMN repository.hierarchy_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.repository.hierarchy_id IS 'hierarchyId';


--
-- Name: COLUMN repository.is_forkable; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.repository.is_forkable IS 'forkable';


--
-- Name: COLUMN repository.is_public; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.repository.is_public IS 'publiclyAccessible';


--
-- Name: repository_access; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.repository_access (
    user_id integer NOT NULL,
    repository_id integer NOT NULL,
    last_accessed bigint NOT NULL
);


ALTER TABLE public.repository_access OWNER TO bitbucketuser;

--
-- Name: sta_activity; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_activity (
    id bigint NOT NULL,
    activity_type integer NOT NULL,
    created_timestamp timestamp without time zone NOT NULL,
    user_id integer NOT NULL
);


ALTER TABLE public.sta_activity OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_activity.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_activity.id IS 'id';


--
-- Name: COLUMN sta_activity.activity_type; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_activity.activity_type IS 'discriminatorColumn';


--
-- Name: COLUMN sta_activity.created_timestamp; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_activity.created_timestamp IS 'createdDate';


--
-- Name: COLUMN sta_activity.user_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_activity.user_id IS 'user';


--
-- Name: sta_cmt_disc_activity; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_cmt_disc_activity (
    activity_id bigint NOT NULL,
    discussion_id bigint NOT NULL
);


ALTER TABLE public.sta_cmt_disc_activity OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_cmt_disc_activity.activity_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_cmt_disc_activity.activity_id IS 'joinPrimaryKey';


--
-- Name: COLUMN sta_cmt_disc_activity.discussion_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_cmt_disc_activity.discussion_id IS 'discussion';


--
-- Name: sta_cmt_disc_participant; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_cmt_disc_participant (
    id bigint NOT NULL,
    discussion_id bigint NOT NULL,
    user_id integer NOT NULL
);


ALTER TABLE public.sta_cmt_disc_participant OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_cmt_disc_participant.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_cmt_disc_participant.id IS 'id';


--
-- Name: COLUMN sta_cmt_disc_participant.discussion_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_cmt_disc_participant.discussion_id IS 'discussion';


--
-- Name: COLUMN sta_cmt_disc_participant.user_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_cmt_disc_participant.user_id IS 'user';


--
-- Name: sta_cmt_discussion; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_cmt_discussion (
    id bigint NOT NULL,
    repository_id integer NOT NULL,
    parent_count integer NOT NULL,
    commit_id character varying(40) NOT NULL,
    parent_id character varying(40)
);


ALTER TABLE public.sta_cmt_discussion OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_cmt_discussion.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_cmt_discussion.id IS 'id';


--
-- Name: COLUMN sta_cmt_discussion.repository_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_cmt_discussion.repository_id IS 'repository';


--
-- Name: COLUMN sta_cmt_discussion.parent_count; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_cmt_discussion.parent_count IS 'parents';


--
-- Name: COLUMN sta_cmt_discussion.commit_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_cmt_discussion.commit_id IS 'commitId';


--
-- Name: COLUMN sta_cmt_discussion.parent_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_cmt_discussion.parent_id IS 'parentId';


--
-- Name: sta_deleted_group; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_deleted_group (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    deleted_timestamp timestamp without time zone NOT NULL
);


ALTER TABLE public.sta_deleted_group OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_deleted_group.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_deleted_group.id IS 'id';


--
-- Name: COLUMN sta_deleted_group.name; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_deleted_group.name IS 'group name';


--
-- Name: COLUMN sta_deleted_group.deleted_timestamp; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_deleted_group.deleted_timestamp IS 'deleted date';


--
-- Name: sta_drift_request; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_drift_request (
    id bigint NOT NULL,
    pr_id bigint NOT NULL,
    old_from_hash character varying(40) NOT NULL,
    old_to_hash character varying(40) NOT NULL,
    new_from_hash character varying(40) NOT NULL,
    new_to_hash character varying(40) NOT NULL,
    attempts integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.sta_drift_request OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_drift_request.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_drift_request.id IS 'id';


--
-- Name: COLUMN sta_drift_request.pr_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_drift_request.pr_id IS 'pullRequest';


--
-- Name: COLUMN sta_drift_request.old_from_hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_drift_request.old_from_hash IS 'oldFromHash';


--
-- Name: COLUMN sta_drift_request.old_to_hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_drift_request.old_to_hash IS 'oldToHash';


--
-- Name: COLUMN sta_drift_request.new_from_hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_drift_request.new_from_hash IS 'newFromHash';


--
-- Name: COLUMN sta_drift_request.new_to_hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_drift_request.new_to_hash IS 'newToHash';


--
-- Name: sta_global_permission; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_global_permission (
    id bigint NOT NULL,
    perm_id integer NOT NULL,
    group_name character varying(255),
    user_id integer
);


ALTER TABLE public.sta_global_permission OWNER TO bitbucketuser;

--
-- Name: sta_normal_project; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_normal_project (
    project_id integer NOT NULL,
    is_public boolean NOT NULL
);


ALTER TABLE public.sta_normal_project OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_normal_project.project_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_normal_project.project_id IS 'id';


--
-- Name: COLUMN sta_normal_project.is_public; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_normal_project.is_public IS 'publiclyAccessible';


--
-- Name: sta_normal_user; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_normal_user (
    user_id integer NOT NULL,
    name character varying(255) NOT NULL,
    slug character varying(255) NOT NULL,
    locale character varying(32),
    deleted_timestamp timestamp without time zone,
    time_zone character varying(64)
);


ALTER TABLE public.sta_normal_user OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_normal_user.user_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_normal_user.user_id IS 'joinPrimaryKey';


--
-- Name: COLUMN sta_normal_user.name; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_normal_user.name IS 'normal user name';


--
-- Name: COLUMN sta_normal_user.slug; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_normal_user.slug IS 'normal user slug';


--
-- Name: COLUMN sta_normal_user.locale; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_normal_user.locale IS 'user_locale';


--
-- Name: COLUMN sta_normal_user.deleted_timestamp; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_normal_user.deleted_timestamp IS 'deletedDate';


--
-- Name: COLUMN sta_normal_user.time_zone; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_normal_user.time_zone IS 'timeZone';


--
-- Name: sta_permission_type; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_permission_type (
    perm_id integer NOT NULL,
    perm_weight integer NOT NULL
);


ALTER TABLE public.sta_permission_type OWNER TO bitbucketuser;

--
-- Name: sta_personal_project; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_personal_project (
    project_id integer NOT NULL,
    owner_id integer NOT NULL
);


ALTER TABLE public.sta_personal_project OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_personal_project.project_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_personal_project.project_id IS 'id';


--
-- Name: COLUMN sta_personal_project.owner_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_personal_project.owner_id IS 'owner';


--
-- Name: sta_pr_activity; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_pr_activity (
    activity_id bigint NOT NULL,
    pr_id bigint NOT NULL,
    pr_action integer NOT NULL
);


ALTER TABLE public.sta_pr_activity OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_pr_activity.activity_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_activity.activity_id IS 'joinPrimaryKey';


--
-- Name: COLUMN sta_pr_activity.pr_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_activity.pr_id IS 'pullRequest';


--
-- Name: COLUMN sta_pr_activity.pr_action; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_activity.pr_action IS 'action';


--
-- Name: sta_pr_merge_activity; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_pr_merge_activity (
    activity_id bigint NOT NULL,
    hash character varying(40)
);


ALTER TABLE public.sta_pr_merge_activity OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_pr_merge_activity.activity_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_merge_activity.activity_id IS 'joinPrimaryKey';


--
-- Name: COLUMN sta_pr_merge_activity.hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_merge_activity.hash IS 'hash';


--
-- Name: sta_pr_participant; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_pr_participant (
    id bigint NOT NULL,
    pr_id bigint NOT NULL,
    pr_role integer NOT NULL,
    user_id integer NOT NULL,
    participant_status integer NOT NULL,
    last_reviewed_commit character varying(40)
);


ALTER TABLE public.sta_pr_participant OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_pr_participant.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_participant.id IS 'id';


--
-- Name: COLUMN sta_pr_participant.pr_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_participant.pr_id IS 'pullRequest';


--
-- Name: COLUMN sta_pr_participant.pr_role; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_participant.pr_role IS 'role';


--
-- Name: COLUMN sta_pr_participant.user_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_participant.user_id IS 'user';


--
-- Name: COLUMN sta_pr_participant.participant_status; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_participant.participant_status IS 'approved';


--
-- Name: COLUMN sta_pr_participant.last_reviewed_commit; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_participant.last_reviewed_commit IS 'lastReviewedCommit';


--
-- Name: sta_pr_rescope_activity; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_pr_rescope_activity (
    activity_id bigint NOT NULL,
    from_hash character varying(40) NOT NULL,
    to_hash character varying(40) NOT NULL,
    prev_from_hash character varying(40) NOT NULL,
    prev_to_hash character varying(40) NOT NULL,
    commits_added integer,
    commits_removed integer
);


ALTER TABLE public.sta_pr_rescope_activity OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_pr_rescope_activity.activity_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_activity.activity_id IS 'joinPrimaryKey';


--
-- Name: COLUMN sta_pr_rescope_activity.from_hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_activity.from_hash IS 'fromHash';


--
-- Name: COLUMN sta_pr_rescope_activity.to_hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_activity.to_hash IS 'toHash';


--
-- Name: COLUMN sta_pr_rescope_activity.prev_from_hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_activity.prev_from_hash IS 'previousFromHash';


--
-- Name: COLUMN sta_pr_rescope_activity.prev_to_hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_activity.prev_to_hash IS 'previousToHash';


--
-- Name: COLUMN sta_pr_rescope_activity.commits_added; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_activity.commits_added IS 'addedCommits';


--
-- Name: COLUMN sta_pr_rescope_activity.commits_removed; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_activity.commits_removed IS 'removedCommits';


--
-- Name: sta_pr_rescope_commit; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_pr_rescope_commit (
    activity_id bigint NOT NULL,
    changeset_id character varying(40) NOT NULL,
    action integer NOT NULL
);


ALTER TABLE public.sta_pr_rescope_commit OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_pr_rescope_commit.activity_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_commit.activity_id IS 'activity';


--
-- Name: COLUMN sta_pr_rescope_commit.changeset_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_commit.changeset_id IS 'changsetId';


--
-- Name: COLUMN sta_pr_rescope_commit.action; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_commit.action IS 'action';


--
-- Name: sta_pr_rescope_request; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_pr_rescope_request (
    id bigint NOT NULL,
    repo_id integer NOT NULL,
    user_id integer NOT NULL,
    created_timestamp timestamp without time zone NOT NULL
);


ALTER TABLE public.sta_pr_rescope_request OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_pr_rescope_request.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_request.id IS 'id';


--
-- Name: COLUMN sta_pr_rescope_request.repo_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_request.repo_id IS 'repository';


--
-- Name: COLUMN sta_pr_rescope_request.user_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_request.user_id IS 'user';


--
-- Name: COLUMN sta_pr_rescope_request.created_timestamp; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_request.created_timestamp IS 'createdDate';


--
-- Name: sta_pr_rescope_request_change; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_pr_rescope_request_change (
    request_id bigint NOT NULL,
    ref_id character varying(1024) NOT NULL,
    change_type integer NOT NULL,
    from_hash character varying(40),
    to_hash character varying(40)
);


ALTER TABLE public.sta_pr_rescope_request_change OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_pr_rescope_request_change.request_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_request_change.request_id IS 'rescopeRequest';


--
-- Name: COLUMN sta_pr_rescope_request_change.ref_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_request_change.ref_id IS 'refId';


--
-- Name: COLUMN sta_pr_rescope_request_change.change_type; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_request_change.change_type IS 'type';


--
-- Name: COLUMN sta_pr_rescope_request_change.from_hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_request_change.from_hash IS 'fromHash';


--
-- Name: COLUMN sta_pr_rescope_request_change.to_hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pr_rescope_request_change.to_hash IS 'toHash';


--
-- Name: sta_project_permission; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_project_permission (
    id bigint NOT NULL,
    perm_id integer NOT NULL,
    project_id integer NOT NULL,
    group_name character varying(255),
    user_id integer
);


ALTER TABLE public.sta_project_permission OWNER TO bitbucketuser;

--
-- Name: sta_pull_request; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_pull_request (
    id bigint NOT NULL,
    entity_version integer NOT NULL,
    scoped_id bigint NOT NULL,
    pr_state integer NOT NULL,
    created_timestamp timestamp without time zone NOT NULL,
    updated_timestamp timestamp without time zone NOT NULL,
    from_repository_id integer NOT NULL,
    to_repository_id integer NOT NULL,
    from_branch_fqn character varying(1024) NOT NULL,
    to_branch_fqn character varying(1024) NOT NULL,
    from_branch_name character varying(255) NOT NULL,
    to_branch_name character varying(255) NOT NULL,
    from_hash character varying(40) NOT NULL,
    to_hash character varying(40) NOT NULL,
    title character varying(255) NOT NULL,
    description text,
    locked_timestamp timestamp without time zone,
    rescoped_timestamp timestamp without time zone NOT NULL,
    closed_timestamp timestamp without time zone
);


ALTER TABLE public.sta_pull_request OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_pull_request.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.id IS 'id, globalId';


--
-- Name: COLUMN sta_pull_request.entity_version; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.entity_version IS 'version';


--
-- Name: COLUMN sta_pull_request.scoped_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.scoped_id IS 'scopedId';


--
-- Name: COLUMN sta_pull_request.pr_state; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.pr_state IS 'state';


--
-- Name: COLUMN sta_pull_request.created_timestamp; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.created_timestamp IS 'createdDate';


--
-- Name: COLUMN sta_pull_request.updated_timestamp; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.updated_timestamp IS 'updatedDate';


--
-- Name: COLUMN sta_pull_request.from_repository_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.from_repository_id IS 'fromRef.repository';


--
-- Name: COLUMN sta_pull_request.to_repository_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.to_repository_id IS 'toRef.repository';


--
-- Name: COLUMN sta_pull_request.from_branch_fqn; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.from_branch_fqn IS 'fromRef.id';


--
-- Name: COLUMN sta_pull_request.to_branch_fqn; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.to_branch_fqn IS 'toRef.id';


--
-- Name: COLUMN sta_pull_request.from_branch_name; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.from_branch_name IS 'fromRef.displayId';


--
-- Name: COLUMN sta_pull_request.to_branch_name; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.to_branch_name IS 'toRef.displayId';


--
-- Name: COLUMN sta_pull_request.from_hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.from_hash IS 'fromRef.hash';


--
-- Name: COLUMN sta_pull_request.to_hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.to_hash IS 'toRef.hash';


--
-- Name: COLUMN sta_pull_request.title; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.title IS 'title';


--
-- Name: COLUMN sta_pull_request.description; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.description IS 'description';


--
-- Name: COLUMN sta_pull_request.locked_timestamp; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.locked_timestamp IS 'lockedDate';


--
-- Name: COLUMN sta_pull_request.rescoped_timestamp; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.rescoped_timestamp IS 'rescopeDate';


--
-- Name: COLUMN sta_pull_request.closed_timestamp; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_pull_request.closed_timestamp IS 'closedDate';


--
-- Name: sta_remember_me_token; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_remember_me_token (
    id bigint NOT NULL,
    series character varying(64) NOT NULL,
    token character varying(64) NOT NULL,
    user_id integer NOT NULL,
    expiry_timestamp timestamp without time zone NOT NULL,
    claimed boolean NOT NULL,
    claimed_address character varying(255)
);


ALTER TABLE public.sta_remember_me_token OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_remember_me_token.user_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_remember_me_token.user_id IS 'userId';


--
-- Name: sta_repo_activity; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_repo_activity (
    activity_id bigint NOT NULL,
    repository_id integer NOT NULL
);


ALTER TABLE public.sta_repo_activity OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_repo_activity.activity_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_repo_activity.activity_id IS 'joinPrimaryKey';


--
-- Name: COLUMN sta_repo_activity.repository_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_repo_activity.repository_id IS 'repository';


--
-- Name: sta_repo_hook; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_repo_hook (
    id bigint NOT NULL,
    repository_id integer,
    hook_key character varying(255) NOT NULL,
    is_enabled boolean NOT NULL,
    lob_id bigint,
    project_id integer
);


ALTER TABLE public.sta_repo_hook OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_repo_hook.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_repo_hook.id IS 'id';


--
-- Name: COLUMN sta_repo_hook.repository_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_repo_hook.repository_id IS 'repository';


--
-- Name: COLUMN sta_repo_hook.hook_key; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_repo_hook.hook_key IS 'hookKey';


--
-- Name: COLUMN sta_repo_hook.is_enabled; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_repo_hook.is_enabled IS 'enabled';


--
-- Name: COLUMN sta_repo_hook.lob_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_repo_hook.lob_id IS 'settings';


--
-- Name: sta_repo_origin; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_repo_origin (
    repository_id integer NOT NULL,
    origin_id integer NOT NULL
);


ALTER TABLE public.sta_repo_origin OWNER TO bitbucketuser;

--
-- Name: sta_repo_permission; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_repo_permission (
    id bigint NOT NULL,
    perm_id integer NOT NULL,
    repo_id integer NOT NULL,
    group_name character varying(255),
    user_id integer
);


ALTER TABLE public.sta_repo_permission OWNER TO bitbucketuser;

--
-- Name: sta_repo_push_activity; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_repo_push_activity (
    activity_id bigint NOT NULL,
    trigger_id character varying(64) NOT NULL
);


ALTER TABLE public.sta_repo_push_activity OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_repo_push_activity.activity_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_repo_push_activity.activity_id IS 'joinPrimaryKey';


--
-- Name: sta_repo_push_ref; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_repo_push_ref (
    activity_id bigint NOT NULL,
    ref_id character varying(1024) NOT NULL,
    change_type integer NOT NULL,
    from_hash character varying(40) NOT NULL,
    to_hash character varying(40) NOT NULL,
    ref_update_type integer NOT NULL
);


ALTER TABLE public.sta_repo_push_ref OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_repo_push_ref.activity_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_repo_push_ref.activity_id IS 'activity';


--
-- Name: COLUMN sta_repo_push_ref.ref_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_repo_push_ref.ref_id IS 'refId';


--
-- Name: COLUMN sta_repo_push_ref.change_type; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_repo_push_ref.change_type IS 'type';


--
-- Name: COLUMN sta_repo_push_ref.from_hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_repo_push_ref.from_hash IS 'fromHash';


--
-- Name: COLUMN sta_repo_push_ref.to_hash; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_repo_push_ref.to_hash IS 'toHash';


--
-- Name: sta_repository_scoped_id; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_repository_scoped_id (
    repository_id integer NOT NULL,
    scope_type character varying(255) NOT NULL,
    next_id bigint NOT NULL
);


ALTER TABLE public.sta_repository_scoped_id OWNER TO bitbucketuser;

--
-- Name: sta_service_user; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_service_user (
    user_id integer NOT NULL,
    display_name character varying(255) NOT NULL,
    active boolean NOT NULL,
    name character varying(255) NOT NULL,
    slug character varying(255) NOT NULL,
    email_address character varying(255),
    label character varying(128) NOT NULL
);


ALTER TABLE public.sta_service_user OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_service_user.user_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_service_user.user_id IS 'joinPrimaryKey';


--
-- Name: COLUMN sta_service_user.display_name; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_service_user.display_name IS 'service user display_name';


--
-- Name: COLUMN sta_service_user.active; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_service_user.active IS 'service user active';


--
-- Name: COLUMN sta_service_user.name; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_service_user.name IS 'service user name';


--
-- Name: COLUMN sta_service_user.slug; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_service_user.slug IS 'service user slug';


--
-- Name: COLUMN sta_service_user.email_address; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_service_user.email_address IS 'service user email';


--
-- Name: COLUMN sta_service_user.label; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_service_user.label IS 'service user label';


--
-- Name: sta_shared_lob; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_shared_lob (
    id bigint NOT NULL,
    lob_data text NOT NULL
);


ALTER TABLE public.sta_shared_lob OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_shared_lob.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_shared_lob.id IS 'id';


--
-- Name: COLUMN sta_shared_lob.lob_data; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_shared_lob.lob_data IS 'data';


--
-- Name: sta_task; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_task (
    id bigint NOT NULL,
    anchor_id bigint NOT NULL,
    anchor_type integer NOT NULL,
    author_id integer NOT NULL,
    context_id bigint NOT NULL,
    context_type integer NOT NULL,
    created_timestamp timestamp without time zone NOT NULL,
    task_state integer NOT NULL,
    task_text text NOT NULL
);


ALTER TABLE public.sta_task OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_task.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_task.id IS 'id';


--
-- Name: COLUMN sta_task.anchor_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_task.anchor_id IS 'anchor.id';


--
-- Name: COLUMN sta_task.anchor_type; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_task.anchor_type IS 'discriminatorColumn';


--
-- Name: COLUMN sta_task.author_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_task.author_id IS 'author';


--
-- Name: COLUMN sta_task.context_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_task.context_id IS 'context.id';


--
-- Name: COLUMN sta_task.context_type; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_task.context_type IS 'discriminatorColumn';


--
-- Name: COLUMN sta_task.created_timestamp; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_task.created_timestamp IS 'createdDate';


--
-- Name: COLUMN sta_task.task_state; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_task.task_state IS 'state';


--
-- Name: COLUMN sta_task.task_text; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_task.task_text IS 'text';


--
-- Name: sta_user_settings; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_user_settings (
    id integer NOT NULL,
    lob_id bigint NOT NULL
);


ALTER TABLE public.sta_user_settings OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_user_settings.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_user_settings.id IS 'id';


--
-- Name: COLUMN sta_user_settings.lob_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_user_settings.lob_id IS 'settings';


--
-- Name: sta_watcher; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.sta_watcher (
    id bigint NOT NULL,
    watchable_id bigint NOT NULL,
    watchable_type integer NOT NULL,
    user_id integer NOT NULL
);


ALTER TABLE public.sta_watcher OWNER TO bitbucketuser;

--
-- Name: COLUMN sta_watcher.id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_watcher.id IS 'id';


--
-- Name: COLUMN sta_watcher.watchable_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_watcher.watchable_id IS 'watchable.id';


--
-- Name: COLUMN sta_watcher.watchable_type; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_watcher.watchable_type IS 'discriminatorColumn';


--
-- Name: COLUMN sta_watcher.user_id; Type: COMMENT; Schema: public; Owner: bitbucketuser
--

COMMENT ON COLUMN public.sta_watcher.user_id IS 'user.id';


--
-- Name: stash_user; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.stash_user (
    id integer NOT NULL
);


ALTER TABLE public.stash_user OWNER TO bitbucketuser;

--
-- Name: trusted_app; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.trusted_app (
    id integer NOT NULL,
    application_id character varying(255) NOT NULL,
    certificate_timeout bigint NOT NULL,
    public_key_base64 character varying(4000) NOT NULL
);


ALTER TABLE public.trusted_app OWNER TO bitbucketuser;

--
-- Name: trusted_app_restriction; Type: TABLE; Schema: public; Owner: bitbucketuser
--

CREATE TABLE public.trusted_app_restriction (
    id integer NOT NULL,
    trusted_app_id integer NOT NULL,
    restriction_type smallint NOT NULL,
    restriction_value character varying(255) NOT NULL
);


ALTER TABLE public.trusted_app_restriction OWNER TO bitbucketuser;

--
-- Name: AO_02A6C0_REJECTED_REF ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_02A6C0_REJECTED_REF" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_02A6C0_REJECTED_REF_ID_seq"'::regclass);


--
-- Name: AO_0E97B5_REPOSITORY_SHORTCUT ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_0E97B5_REPOSITORY_SHORTCUT" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_0E97B5_REPOSITORY_SHORTCUT_ID_seq"'::regclass);


--
-- Name: AO_2AD648_INSIGHT_ANNOTATION ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_2AD648_INSIGHT_ANNOTATION" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_2AD648_INSIGHT_ANNOTATION_ID_seq"'::regclass);


--
-- Name: AO_2AD648_INSIGHT_REPORT ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_2AD648_INSIGHT_REPORT" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_2AD648_INSIGHT_REPORT_ID_seq"'::regclass);


--
-- Name: AO_2AD648_MERGE_CHECK ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_2AD648_MERGE_CHECK" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_2AD648_MERGE_CHECK_ID_seq"'::regclass);


--
-- Name: AO_33D892_COMMENT_JIRA_ISSUE ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_33D892_COMMENT_JIRA_ISSUE" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_33D892_COMMENT_JIRA_ISSUE_ID_seq"'::regclass);


--
-- Name: AO_38321B_CUSTOM_CONTENT_LINK ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_38321B_CUSTOM_CONTENT_LINK" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_38321B_CUSTOM_CONTENT_LINK_ID_seq"'::regclass);


--
-- Name: AO_38F373_COMMENT_LIKE ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_38F373_COMMENT_LIKE" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_38F373_COMMENT_LIKE_ID_seq"'::regclass);


--
-- Name: AO_4789DD_HEALTH_CHECK_STATUS ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_4789DD_HEALTH_CHECK_STATUS" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_4789DD_HEALTH_CHECK_STATUS_ID_seq"'::regclass);


--
-- Name: AO_4789DD_PROPERTIES ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_4789DD_PROPERTIES" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_4789DD_PROPERTIES_ID_seq"'::regclass);


--
-- Name: AO_4789DD_READ_NOTIFICATIONS ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_4789DD_READ_NOTIFICATIONS" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_4789DD_READ_NOTIFICATIONS_ID_seq"'::regclass);


--
-- Name: AO_4789DD_TASK_MONITOR ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_4789DD_TASK_MONITOR" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_4789DD_TASK_MONITOR_ID_seq"'::regclass);


--
-- Name: AO_616D7B_BRANCH_MODEL_CONFIG ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_616D7B_BRANCH_MODEL_CONFIG" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_616D7B_BRANCH_MODEL_CONFIG_ID_seq"'::regclass);


--
-- Name: AO_616D7B_BRANCH_TYPE ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_616D7B_BRANCH_TYPE" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_616D7B_BRANCH_TYPE_ID_seq"'::regclass);


--
-- Name: AO_616D7B_BRANCH_TYPE_CONFIG ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_616D7B_BRANCH_TYPE_CONFIG" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_616D7B_BRANCH_TYPE_CONFIG_ID_seq"'::regclass);


--
-- Name: AO_616D7B_SCOPE_AUTO_MERGE ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_616D7B_SCOPE_AUTO_MERGE" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_616D7B_SCOPE_AUTO_MERGE_ID_seq"'::regclass);


--
-- Name: AO_6978BB_PERMITTED_ENTITY ENTITY_ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_6978BB_PERMITTED_ENTITY" ALTER COLUMN "ENTITY_ID" SET DEFAULT nextval('public."AO_6978BB_PERMITTED_ENTITY_ENTITY_ID_seq"'::regclass);


--
-- Name: AO_6978BB_RESTRICTED_REF REF_ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_6978BB_RESTRICTED_REF" ALTER COLUMN "REF_ID" SET DEFAULT nextval('public."AO_6978BB_RESTRICTED_REF_REF_ID_seq"'::regclass);


--
-- Name: AO_777666_JIRA_INDEX ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_777666_JIRA_INDEX" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_777666_JIRA_INDEX_ID_seq"'::regclass);


--
-- Name: AO_811463_GIT_LFS_LOCK ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_811463_GIT_LFS_LOCK" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_811463_GIT_LFS_LOCK_ID_seq"'::regclass);


--
-- Name: AO_8E6075_MIRRORING_REQUEST ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_8E6075_MIRRORING_REQUEST" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_8E6075_MIRRORING_REQUEST_ID_seq"'::regclass);


--
-- Name: AO_92D5D5_REPO_NOTIFICATION ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_92D5D5_REPO_NOTIFICATION" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_92D5D5_REPO_NOTIFICATION_ID_seq"'::regclass);


--
-- Name: AO_92D5D5_USER_NOTIFICATION ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_92D5D5_USER_NOTIFICATION" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_92D5D5_USER_NOTIFICATION_ID_seq"'::regclass);


--
-- Name: AO_9DEC2A_DEFAULT_REVIEWER ENTITY_ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_9DEC2A_DEFAULT_REVIEWER" ALTER COLUMN "ENTITY_ID" SET DEFAULT nextval('public."AO_9DEC2A_DEFAULT_REVIEWER_ENTITY_ID_seq"'::regclass);


--
-- Name: AO_9DEC2A_PR_CONDITION PR_CONDITION_ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_9DEC2A_PR_CONDITION" ALTER COLUMN "PR_CONDITION_ID" SET DEFAULT nextval('public."AO_9DEC2A_PR_CONDITION_PR_CONDITION_ID_seq"'::regclass);


--
-- Name: AO_A0B856_WEBHOOK ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_A0B856_WEBHOOK" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_A0B856_WEBHOOK_ID_seq"'::regclass);


--
-- Name: AO_A0B856_WEBHOOK_CONFIG ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_A0B856_WEBHOOK_CONFIG" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_A0B856_WEBHOOK_CONFIG_ID_seq"'::regclass);


--
-- Name: AO_A0B856_WEBHOOK_EVENT ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_A0B856_WEBHOOK_EVENT" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_A0B856_WEBHOOK_EVENT_ID_seq"'::regclass);


--
-- Name: AO_A0B856_WEB_HOOK_LISTENER_AO ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_A0B856_WEB_HOOK_LISTENER_AO" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_A0B856_WEB_HOOK_LISTENER_AO_ID_seq"'::regclass);


--
-- Name: AO_BD73C3_PROJECT_AUDIT AUDIT_ITEM_ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_BD73C3_PROJECT_AUDIT" ALTER COLUMN "AUDIT_ITEM_ID" SET DEFAULT nextval('public."AO_BD73C3_PROJECT_AUDIT_AUDIT_ITEM_ID_seq"'::regclass);


--
-- Name: AO_BD73C3_REPOSITORY_AUDIT AUDIT_ITEM_ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_BD73C3_REPOSITORY_AUDIT" ALTER COLUMN "AUDIT_ITEM_ID" SET DEFAULT nextval('public."AO_BD73C3_REPOSITORY_AUDIT_AUDIT_ITEM_ID_seq"'::regclass);


--
-- Name: AO_C77861_AUDIT_ENTITY ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_C77861_AUDIT_ENTITY" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_C77861_AUDIT_ENTITY_ID_seq"'::regclass);


--
-- Name: AO_CFE8FA_BUILD_STATUS ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_CFE8FA_BUILD_STATUS" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_CFE8FA_BUILD_STATUS_ID_seq"'::regclass);


--
-- Name: AO_D6A508_IMPORT_JOB ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_D6A508_IMPORT_JOB" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_D6A508_IMPORT_JOB_ID_seq"'::regclass);


--
-- Name: AO_D6A508_REPO_IMPORT_TASK ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_D6A508_REPO_IMPORT_TASK" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_D6A508_REPO_IMPORT_TASK_ID_seq"'::regclass);


--
-- Name: AO_E5A814_ACCESS_TOKEN_PERM ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_E5A814_ACCESS_TOKEN_PERM" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_E5A814_ACCESS_TOKEN_PERM_ID_seq"'::regclass);


--
-- Name: AO_ED669C_SEEN_ASSERTIONS ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_ED669C_SEEN_ASSERTIONS" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_ED669C_SEEN_ASSERTIONS_ID_seq"'::regclass);


--
-- Name: AO_F4ED3A_ADD_ON_PROPERTY_AO ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_F4ED3A_ADD_ON_PROPERTY_AO" ALTER COLUMN "ID" SET DEFAULT nextval('public."AO_F4ED3A_ADD_ON_PROPERTY_AO_ID_seq"'::regclass);


--
-- Name: AO_FB71B4_SSH_PUBLIC_KEY ENTITY_ID; Type: DEFAULT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_FB71B4_SSH_PUBLIC_KEY" ALTER COLUMN "ENTITY_ID" SET DEFAULT nextval('public."AO_FB71B4_SSH_PUBLIC_KEY_ENTITY_ID_seq"'::regclass);


--
-- Data for Name: AO_02A6C0_REJECTED_REF; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_02A6C0_REJECTED_REF" ("ID", "REF_DISPLAY_ID", "REF_ID", "REF_STATUS", "REPOSITORY_ID") FROM stdin;
\.


--
-- Data for Name: AO_02A6C0_SYNC_CONFIG; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_02A6C0_SYNC_CONFIG" ("IS_ENABLED", "LAST_SYNC", "REPOSITORY_ID") FROM stdin;
\.


--
-- Data for Name: AO_0E97B5_REPOSITORY_SHORTCUT; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_0E97B5_REPOSITORY_SHORTCUT" ("APPLICATION_LINK_ID", "CREATED_DATE", "ID", "LABEL", "PRODUCT_TYPE", "REPOSITORY_ID", "URL") FROM stdin;
\.


--
-- Data for Name: AO_2AD648_INSIGHT_ANNOTATION; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_2AD648_INSIGHT_ANNOTATION" ("EXTERNAL_ID", "FK_REPORT_ID", "ID", "LINE", "LINK", "MESSAGE", "PATH", "PATH_MD5", "SEVERITY_ID", "TYPE_ID") FROM stdin;
\.


--
-- Data for Name: AO_2AD648_INSIGHT_REPORT; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_2AD648_INSIGHT_REPORT" ("AUTHOR_ID", "COMMIT_ID", "COVERAGE_PROVIDER_KEY", "CREATED_DATE", "DATA", "DETAILS", "ID", "LINK", "LOGO", "REPORTER", "REPORT_KEY", "REPOSITORY_ID", "RESULT_ID", "TITLE") FROM stdin;
\.


--
-- Data for Name: AO_2AD648_MERGE_CHECK; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_2AD648_MERGE_CHECK" ("ANNOTATION_SEVERITY", "ID", "MUST_PASS", "REPORT_KEY", "RESOURCE_ID", "SCOPE_TYPE") FROM stdin;
\.


--
-- Data for Name: AO_33D892_COMMENT_JIRA_ISSUE; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_33D892_COMMENT_JIRA_ISSUE" ("COMMENT_ID", "ID", "ISSUE_KEY") FROM stdin;
\.


--
-- Data for Name: AO_38321B_CUSTOM_CONTENT_LINK; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_38321B_CUSTOM_CONTENT_LINK" ("CONTENT_KEY", "ID", "LINK_LABEL", "LINK_URL", "SEQUENCE") FROM stdin;
\.


--
-- Data for Name: AO_38F373_COMMENT_LIKE; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_38F373_COMMENT_LIKE" ("COMMENT_ID", "ID", "USER_ID") FROM stdin;
\.


--
-- Data for Name: AO_4789DD_HEALTH_CHECK_STATUS; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_4789DD_HEALTH_CHECK_STATUS" ("APPLICATION_NAME", "COMPLETE_KEY", "DESCRIPTION", "FAILED_DATE", "FAILURE_REASON", "ID", "IS_HEALTHY", "IS_RESOLVED", "RESOLVED_DATE", "SEVERITY", "STATUS_NAME") FROM stdin;
\.


--
-- Data for Name: AO_4789DD_PROPERTIES; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_4789DD_PROPERTIES" ("ID", "PROPERTY_NAME", "PROPERTY_VALUE") FROM stdin;
\.


--
-- Data for Name: AO_4789DD_READ_NOTIFICATIONS; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_4789DD_READ_NOTIFICATIONS" ("ID", "IS_SNOOZED", "NOTIFICATION_ID", "SNOOZE_COUNT", "SNOOZE_DATE", "USER_KEY") FROM stdin;
\.


--
-- Data for Name: AO_4789DD_TASK_MONITOR; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_4789DD_TASK_MONITOR" ("CLUSTERED_TASK_ID", "CREATED_TIMESTAMP", "ID", "NODE_ID", "PROGRESS_MESSAGE", "PROGRESS_PERCENTAGE", "SERIALIZED_ERRORS", "SERIALIZED_WARNINGS", "TASK_ID", "TASK_MONITOR_KIND", "TASK_STATUS") FROM stdin;
\.


--
-- Data for Name: AO_616D7B_BRANCH_MODEL; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_616D7B_BRANCH_MODEL" ("DEV_ID", "DEV_USE_DEFAULT", "IS_ENABLED", "PROD_ID", "PROD_USE_DEFAULT", "REPOSITORY_ID") FROM stdin;
\.


--
-- Data for Name: AO_616D7B_BRANCH_MODEL_CONFIG; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_616D7B_BRANCH_MODEL_CONFIG" ("DEV_ID", "DEV_USE_DEFAULT", "ID", "PROD_ID", "PROD_USE_DEFAULT", "RESOURCE_ID", "SCOPE_TYPE") FROM stdin;
\N	t	1	\N	f	1	PROJECT
\.


--
-- Data for Name: AO_616D7B_BRANCH_TYPE; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_616D7B_BRANCH_TYPE" ("FK_BM_ID", "ID", "IS_ENABLED", "PREFIX", "TYPE_ID") FROM stdin;
\.


--
-- Data for Name: AO_616D7B_BRANCH_TYPE_CONFIG; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_616D7B_BRANCH_TYPE_CONFIG" ("BM_ID", "ID", "IS_ENABLED", "PREFIX", "TYPE_ID") FROM stdin;
1	1	t	bugfix/	BUGFIX
1	2	t	feature/	FEATURE
1	3	t	hotfix/	HOTFIX
1	4	t	release/	RELEASE
\.


--
-- Data for Name: AO_616D7B_SCOPE_AUTO_MERGE; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_616D7B_SCOPE_AUTO_MERGE" ("ENABLED", "ID", "MERGE_CHECK_ENABLED", "RESOURCE_ID", "SCOPE_TYPE") FROM stdin;
\.


--
-- Data for Name: AO_6978BB_PERMITTED_ENTITY; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_6978BB_PERMITTED_ENTITY" ("ACCESS_KEY_ID", "ENTITY_ID", "FK_RESTRICTED_ID", "GROUP_ID", "USER_ID") FROM stdin;
\.


--
-- Data for Name: AO_6978BB_RESTRICTED_REF; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_6978BB_RESTRICTED_REF" ("REF_ID", "REF_TYPE", "REF_VALUE", "RESOURCE_ID", "RESTRICTION_TYPE", "SCOPE_TYPE") FROM stdin;
\.


--
-- Data for Name: AO_777666_JIRA_INDEX; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_777666_JIRA_INDEX" ("BRANCH", "ID", "ISSUE", "LAST_UPDATED", "PR_ID", "PR_STATE", "REPOSITORY") FROM stdin;
\.


--
-- Data for Name: AO_777666_UPDATED_ISSUES; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_777666_UPDATED_ISSUES" ("ISSUE", "UPDATE_TIME") FROM stdin;
\.


--
-- Data for Name: AO_811463_GIT_LFS_LOCK; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_811463_GIT_LFS_LOCK" ("DIRECTORY_HASH", "ID", "LOCKED_AT", "OWNER_ID", "PATH", "REPOSITORY_ID", "REPO_PATH_HASH") FROM stdin;
\.


--
-- Data for Name: AO_811463_GIT_LFS_REPO_CONFIG; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_811463_GIT_LFS_REPO_CONFIG" ("IS_ENABLED", "REPOSITORY_ID") FROM stdin;
\.


--
-- Data for Name: AO_8E6075_MIRRORING_REQUEST; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_8E6075_MIRRORING_REQUEST" ("ADDON_DESCRIPTOR_URI", "BASE_URL", "CREATED_DATE", "DESCRIPTOR_URL", "ID", "MIRROR_ID", "MIRROR_NAME", "MIRROR_TYPE", "PRODUCT_TYPE", "PRODUCT_VERSION", "RESOLVED_DATE", "RESOLVER_USER_ID", "STATE") FROM stdin;
\.


--
-- Data for Name: AO_8E6075_MIRROR_SERVER; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_8E6075_MIRROR_SERVER" ("ADD_ON_KEY", "BASE_URL", "ID", "LAST_SEEN", "MIRROR_TYPE", "NAME", "PRODUCT_TYPE", "PRODUCT_VERSION", "STATE") FROM stdin;
\.


--
-- Data for Name: AO_92D5D5_REPO_NOTIFICATION; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_92D5D5_REPO_NOTIFICATION" ("ID", "PR_NOTIFICATION_SCOPE", "PUSH_NOTIFICATION_SCOPE", "REPO_ID", "USER_ID") FROM stdin;
\.


--
-- Data for Name: AO_92D5D5_USER_NOTIFICATION; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_92D5D5_USER_NOTIFICATION" ("BATCH_ID", "BATCH_SENDER_ID", "DATA", "DATE", "ID", "USER_ID") FROM stdin;
\.


--
-- Data for Name: AO_9DEC2A_DEFAULT_REVIEWER; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_9DEC2A_DEFAULT_REVIEWER" ("ENTITY_ID", "FK_RESTRICTED_ID", "USER_ID") FROM stdin;
\.


--
-- Data for Name: AO_9DEC2A_PR_CONDITION; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_9DEC2A_PR_CONDITION" ("PR_CONDITION_ID", "REQUIRED_APPROVALS", "RESOURCE_ID", "SCOPE_TYPE", "SOURCE_REF_TYPE", "SOURCE_REF_VALUE", "TARGET_REF_TYPE", "TARGET_REF_VALUE") FROM stdin;
\.


--
-- Data for Name: AO_A0B856_DAILY_COUNTS; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_A0B856_DAILY_COUNTS" ("DAY_SINCE_EPOCH", "ERRORS", "EVENT_ID", "FAILURES", "ID", "SUCCESSES", "WEBHOOK_ID") FROM stdin;
\.


--
-- Data for Name: AO_A0B856_HIST_INVOCATION; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_A0B856_HIST_INVOCATION" ("ERROR_CONTENT", "EVENT_ID", "FINISH", "ID", "OUTCOME", "REQUEST_BODY", "REQUEST_HEADERS", "REQUEST_ID", "REQUEST_METHOD", "REQUEST_URL", "RESPONSE_BODY", "RESPONSE_HEADERS", "RESULT_DESCRIPTION", "START", "STATUS_CODE", "WEBHOOK_ID") FROM stdin;
\.


--
-- Data for Name: AO_A0B856_WEBHOOK; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_A0B856_WEBHOOK" ("ACTIVE", "CREATED", "ID", "NAME", "SCOPE_ID", "SCOPE_TYPE", "UPDATED", "URL") FROM stdin;
\.


--
-- Data for Name: AO_A0B856_WEBHOOK_CONFIG; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_A0B856_WEBHOOK_CONFIG" ("ID", "KEY", "VALUE", "WEBHOOKID") FROM stdin;
\.


--
-- Data for Name: AO_A0B856_WEBHOOK_EVENT; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_A0B856_WEBHOOK_EVENT" ("EVENT_ID", "ID", "WEBHOOKID") FROM stdin;
\.


--
-- Data for Name: AO_A0B856_WEB_HOOK_LISTENER_AO; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_A0B856_WEB_HOOK_LISTENER_AO" ("DESCRIPTION", "ENABLED", "EVENTS", "EXCLUDE_BODY", "FILTERS", "ID", "LAST_UPDATED", "LAST_UPDATED_USER", "NAME", "PARAMETERS", "REGISTRATION_METHOD", "URL") FROM stdin;
\.


--
-- Data for Name: AO_B586BC_GPG_KEY; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_B586BC_GPG_KEY" ("EMAIL", "EXPIRY_DATE", "FINGERPRINT", "KEY_ID", "KEY_TEXT", "USER_ID") FROM stdin;
\.


--
-- Data for Name: AO_B586BC_GPG_SUB_KEY; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_B586BC_GPG_SUB_KEY" ("EXPIRY_DATE", "FINGERPRINT", "FK_GPG_KEY_ID", "KEY_ID") FROM stdin;
\.


--
-- Data for Name: AO_BD73C3_PROJECT_AUDIT; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_BD73C3_PROJECT_AUDIT" ("ACTION", "AUDIT_ITEM_ID", "DATE", "DETAILS", "PROJECT_ID", "USER") FROM stdin;
\.


--
-- Data for Name: AO_BD73C3_REPOSITORY_AUDIT; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_BD73C3_REPOSITORY_AUDIT" ("ACTION", "AUDIT_ITEM_ID", "DATE", "DETAILS", "REPOSITORY_ID", "USER") FROM stdin;
\.


--
-- Data for Name: AO_C77861_AUDIT_ENTITY; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_C77861_AUDIT_ENTITY" ("ACTION", "AREA", "ATTRIBUTES", "CATEGORY", "CHANGE_VALUES", "ENTITY_TIMESTAMP", "ID", "LEVEL", "METHOD", "NODE", "PRIMARY_RESOURCE_ID", "PRIMARY_RESOURCE_TYPE", "RESOURCES", "SEARCH_STRING", "SECONDARY_RESOURCE_ID", "SECONDARY_RESOURCE_TYPE", "SOURCE", "SYSTEM_INFO", "USER_ID", "USER_NAME", "USER_TYPE") FROM stdin;
DirectoryCreatedEvent	USER_MANAGEMENT	[{"name":"target","value":"Bitbucket Internal Directory"}]	Users and groups	[]	1622534202419	1	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"Bitbucket Internal Directory","type":"MISC","uri":null,"id":null}]	bitbucket internal directory directorycreatedevent users and groups system	\N	\N	\N	\N	-1	System	system
GroupCreatedEvent	USER_MANAGEMENT	[{"name":"target","value":"stash-users"}]	Users and groups	[]	1622534203259	2	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"stash-users","type":"MISC","uri":null,"id":null}]	stash-users groupcreatedevent users and groups system	\N	\N	\N	\N	-1	System	system
DisplayNameChangedEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"details","value":"{\\"new\\":\\"Bitbucket\\",\\"old\\":null}"},{"name":"target","value":"DISPLAY_NAME"}]	Global administration	[]	1622534203725	3	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"DISPLAY_NAME","type":"MISC","uri":null,"id":null}]	display_name displaynamechangedevent global administration system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.templaterenderer.api"}]	Apps	[]	1622534264486	4	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.templaterenderer.api","type":"MISC","uri":null,"id":null}]	com.atlassian.templaterenderer.api pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.activeobjects.activeobjects-plugin"}]	Apps	[]	1622534264535	5	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.activeobjects.activeobjects-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.activeobjects.activeobjects-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.templaterenderer.atlassian-template-renderer-velocity1.6-plugin"}]	Apps	[]	1622534264634	6	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.templaterenderer.atlassian-template-renderer-velocity1.6-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.templaterenderer.atlassian-template-renderer-velocity1.6-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.rest.atlassian-rest-module"}]	Apps	[]	1622534264673	7	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.rest.atlassian-rest-module","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.rest.atlassian-rest-module pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.analytics.analytics-client"}]	Apps	[]	1622534264918	8	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.analytics.analytics-client","type":"MISC","uri":null,"id":null}]	com.atlassian.analytics.analytics-client pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.analytics.analytics-whitelist"}]	Apps	[]	1622534264930	9	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.analytics.analytics-whitelist","type":"MISC","uri":null,"id":null}]	com.atlassian.analytics.analytics-whitelist pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bundles.json-20070829.0.0.1"}]	Apps	[]	1622534264934	10	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bundles.json-20070829.0.0.1","type":"MISC","uri":null,"id":null}]	com.atlassian.bundles.json-20070829.0.0.1 pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.oauth.atlassian-oauth-consumer-spi"}]	Apps	[]	1622534264936	11	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.oauth.atlassian-oauth-consumer-spi","type":"MISC","uri":null,"id":null}]	com.atlassian.oauth.atlassian-oauth-consumer-spi pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.oauth.atlassian-oauth-service-provider-spi"}]	Apps	[]	1622534264940	12	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.oauth.atlassian-oauth-service-provider-spi","type":"MISC","uri":null,"id":null}]	com.atlassian.oauth.atlassian-oauth-service-provider-spi pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.springsource.org.jdom-1.1.0"}]	Apps	[]	1622534264942	13	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.springsource.org.jdom-1.1.0","type":"MISC","uri":null,"id":null}]	com.springsource.org.jdom-1.1.0 pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.applinks.applinks-plugin"}]	Apps	[]	1622534265483	14	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.applinks.applinks-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.applinks.applinks-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.applinks.applinks-basicauth-plugin"}]	Apps	[]	1622534265562	15	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.applinks.applinks-basicauth-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.applinks.applinks-basicauth-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.applinks.applinks-cors-plugin"}]	Apps	[]	1622534265589	16	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.applinks.applinks-cors-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.applinks.applinks-cors-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.oauth.consumer.sal"}]	Apps	[]	1622534265596	17	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.oauth.consumer.sal","type":"MISC","uri":null,"id":null}]	com.atlassian.oauth.consumer.sal pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.applinks.applinks-oauth-plugin"}]	Apps	[]	1622534265688	18	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.applinks.applinks-oauth-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.applinks.applinks-oauth-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.applinks.applinks-trustedapps-plugin"}]	Apps	[]	1622534265754	19	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.applinks.applinks-trustedapps-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.applinks.applinks-trustedapps-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.atlassian-failure-cache-plugin"}]	Apps	[]	1622534265799	20	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.atlassian-failure-cache-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.atlassian-failure-cache-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.audit.atlassian-audit-plugin"}]	Apps	[]	1622534265972	21	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.audit.atlassian-audit-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.audit.atlassian-audit-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.auiplugin"}]	Apps	[]	1622534266804	22	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.auiplugin","type":"MISC","uri":null,"id":null}]	com.atlassian.auiplugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-ao-common"}]	Apps	[]	1622534266812	23	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-ao-common","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-ao-common pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-webhooks"}]	Apps	[]	1622534266874	24	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-webhooks","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-webhooks pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-rest"}]	Apps	[]	1622534266905	25	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-rest","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-rest pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-access-tokens"}]	Apps	[]	1622534266954	26	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-access-tokens","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-access-tokens pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-analytics-whitelist"}]	Apps	[]	1622534266970	27	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-analytics-whitelist","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-analytics-whitelist pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-announcement-banner"}]	Apps	[]	1622534267026	28	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-announcement-banner","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-announcement-banner pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-audit"}]	Apps	[]	1622534267112	29	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-audit","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-audit pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-authentication"}]	Apps	[]	1622534267159	30	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-authentication","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-authentication pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.stash.ssh-plugin"}]	Apps	[]	1622534267269	31	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.stash.ssh-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.stash.ssh-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-ref-restriction"}]	Apps	[]	1622534267416	32	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-ref-restriction","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-ref-restriction pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.integration.jira.jira-integration-plugin"}]	Apps	[]	1622534267450	33	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.integration.jira.jira-integration-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.integration.jira.jira-integration-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-branch"}]	Apps	[]	1622534267571	34	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-branch","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-branch pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-build"}]	Apps	[]	1622534267650	35	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-build","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-build pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-bundled-hooks"}]	Apps	[]	1622534267734	36	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-bundled-hooks","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-bundled-hooks pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-client-web-fragments"}]	Apps	[]	1622534267992	37	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-client-web-fragments","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-client-web-fragments pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-cluster-info"}]	Apps	[]	1622534268030	38	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-cluster-info","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-cluster-info pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-code-insights"}]	Apps	[]	1622534268302	39	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-code-insights","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-code-insights pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-comment-likes"}]	Apps	[]	1622534268372	40	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-comment-likes","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-comment-likes pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-comment-properties"}]	Apps	[]	1622534268376	41	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-comment-properties","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-comment-properties pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-compare"}]	Apps	[]	1622534268392	42	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-compare","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-compare pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.jwt.jwt-plugin"}]	Apps	[]	1622534268412	43	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.jwt.jwt-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.jwt.jwt-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-connect-support"}]	Apps	[]	1622534268421	44	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-connect-support","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-connect-support pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-contributing-guidelines"}]	Apps	[]	1622534268473	45	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-contributing-guidelines","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-contributing-guidelines pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-crowd-spi"}]	Apps	[]	1622534268497	46	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-crowd-spi","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-crowd-spi pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-crowd-sso"}]	Apps	[]	1622534268503	47	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-crowd-sso","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-crowd-sso pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-dashboard"}]	Apps	[]	1622534268521	48	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-dashboard","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-dashboard pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-default-reviewers"}]	Apps	[]	1622534268627	49	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-default-reviewers","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-default-reviewers pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-emoticons"}]	Apps	[]	1622534269016	50	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-emoticons","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-emoticons pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-frontend"}]	Apps	[]	1622534270540	51	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-frontend","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-frontend pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-git"}]	Apps	[]	1622534270638	52	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-git","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-git pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.httpclient.atlassian-httpclient-plugin"}]	Apps	[]	1622534270698	53	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.httpclient.atlassian-httpclient-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.httpclient.atlassian-httpclient-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-git-lfs"}]	Apps	[]	1622534270746	54	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-git-lfs","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-git-lfs pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-git-rest"}]	Apps	[]	1622534270884	55	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-git-rest","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-git-rest pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-gpg"}]	Apps	[]	1622534271095	56	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-gpg","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-gpg pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-highlight"}]	Apps	[]	1622534271492	57	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-highlight","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-highlight pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-http-scm-protocol"}]	Apps	[]	1622534271577	58	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-http-scm-protocol","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-http-scm-protocol pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-i18n"}]	Apps	[]	1622534271598	59	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-i18n","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-i18n pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-importer"}]	Apps	[]	1622534271735	60	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-importer","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-importer pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-instance-migration"}]	Apps	[]	1622534271750	61	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-instance-migration","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-instance-migration pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-jira"}]	Apps	[]	1622534272260	62	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-jira","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-jira pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-repository-ref-sync"}]	Apps	[]	1622534272736	63	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-repository-ref-sync","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-repository-ref-sync pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugin.atlassian-spring-scanner-annotation"}]	Apps	[]	1622534272774	64	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugin.atlassian-spring-scanner-annotation","type":"MISC","uri":null,"id":null}]	com.atlassian.plugin.atlassian-spring-scanner-annotation pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.atlassian-nav-links-plugin"}]	Apps	[]	1622534273178	65	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.atlassian-nav-links-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.atlassian-nav-links-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.atlassian-remote-event-common-plugin"}]	Apps	[]	1622534273246	66	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.atlassian-remote-event-common-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.atlassian-remote-event-common-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.remote-link-aggregator-plugin"}]	Apps	[]	1622534273328	67	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.remote-link-aggregator-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.remote-link-aggregator-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-jira-development-integration"}]	Apps	[]	1622534273425	68	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-jira-development-integration","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-jira-development-integration pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-keyboard-shortcuts"}]	Apps	[]	1622534273535	69	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-keyboard-shortcuts","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-keyboard-shortcuts pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-labels"}]	Apps	[]	1622534273574	70	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-labels","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-labels pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-markup-renderers"}]	Apps	[]	1622534273618	71	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-markup-renderers","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-markup-renderers pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bundles.json-schema-validator-atlassian-bundle"}]	Apps	[]	1622534273621	72	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bundles.json-schema-validator-atlassian-bundle","type":"MISC","uri":null,"id":null}]	com.atlassian.bundles.json-schema-validator-atlassian-bundle pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"rome.rome-1.0"}]	Apps	[]	1622534273629	73	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"rome.rome-1.0","type":"MISC","uri":null,"id":null}]	rome.rome-1.0 pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.upm.atlassian-universal-plugin-manager-plugin"}]	Apps	[]	1622534273971	74	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.upm.atlassian-universal-plugin-manager-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.upm.atlassian-universal-plugin-manager-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.atlassian-connect-plugin"}]	Apps	[]	1622534274153	75	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.atlassian-connect-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.atlassian-connect-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-mirroring-upstream"}]	Apps	[]	1622534274192	76	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-mirroring-upstream","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-mirroring-upstream pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-nav-links"}]	Apps	[]	1622534274202	77	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-nav-links","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-nav-links pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-notification"}]	Apps	[]	1622534274408	78	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-notification","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-notification pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-page-data"}]	Apps	[]	1622534274467	79	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-page-data","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-page-data pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-plugin-information-provider"}]	Apps	[]	1622534274468	80	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-plugin-information-provider","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-plugin-information-provider pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-pull-request-cleanup"}]	Apps	[]	1622534274504	81	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-pull-request-cleanup","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-pull-request-cleanup pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-pull-request-properties"}]	Apps	[]	1622534274516	82	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-pull-request-properties","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-pull-request-properties pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-rate-limit"}]	Apps	[]	1622534274528	83	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-rate-limit","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-rate-limit pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-ref-metadata"}]	Apps	[]	1622534274578	84	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-ref-metadata","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-ref-metadata pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-repository-hooks"}]	Apps	[]	1622534274631	85	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-repository-hooks","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-repository-hooks pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-repository-shortcuts"}]	Apps	[]	1622534274670	86	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-repository-shortcuts","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-repository-shortcuts pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-rest-ui"}]	Apps	[]	1622534274722	87	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-rest-ui","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-rest-ui pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-sal"}]	Apps	[]	1622534274725	88	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-sal","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-sal pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-search"}]	Apps	[]	1622534275042	89	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-search","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-search pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.oauth.admin"}]	Apps	[]	1622534279627	108	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.oauth.admin","type":"MISC","uri":null,"id":null}]	com.atlassian.oauth.admin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-server-web-fragments"}]	Apps	[]	1622534275665	90	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-server-web-fragments","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-server-web-fragments pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-sourcetree"}]	Apps	[]	1622534275750	91	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-sourcetree","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-sourcetree pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-soy-functions"}]	Apps	[]	1622534276642	92	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-soy-functions","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-soy-functions pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-suggestions"}]	Apps	[]	1622534276670	93	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-suggestions","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-suggestions pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-tag"}]	Apps	[]	1622534276725	94	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-tag","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-tag pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-user-erasure"}]	Apps	[]	1622534276737	95	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-user-erasure","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-user-erasure pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-velocity-helper"}]	Apps	[]	1622534276738	96	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-velocity-helper","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-velocity-helper pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-web"}]	Apps	[]	1622534277622	97	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-web","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-web pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-web-api"}]	Apps	[]	1622534277646	98	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-web-api","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-web-api pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-web-resource-transformers"}]	Apps	[]	1622534277652	99	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-web-resource-transformers","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-web-resource-transformers pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-webpack-INTERNAL"}]	Apps	[]	1622534279373	100	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-webpack-INTERNAL","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-webpack-internal pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.bitbucket-xcode"}]	Apps	[]	1622534279389	101	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.bitbucket-xcode","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.bitbucket-xcode pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.config-wrm-data"}]	Apps	[]	1622534279475	102	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.config-wrm-data","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.config-wrm-data pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.feature-wrm-data"}]	Apps	[]	1622534279488	103	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.feature-wrm-data","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.feature-wrm-data pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.bitbucket.server.support-info-providers"}]	Apps	[]	1622534279505	104	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.bitbucket.server.support-info-providers","type":"MISC","uri":null,"id":null}]	com.atlassian.bitbucket.server.support-info-providers pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.crowd.embedded.admin"}]	Apps	[]	1622534279555	105	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.crowd.embedded.admin","type":"MISC","uri":null,"id":null}]	com.atlassian.crowd.embedded.admin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.diagnostics.atlassian-diagnostics-plugin"}]	Apps	[]	1622534279596	106	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.diagnostics.atlassian-diagnostics-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.diagnostics.atlassian-diagnostics-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.healthcheck.atlassian-healthcheck"}]	Apps	[]	1622534279625	107	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.healthcheck.atlassian-healthcheck","type":"MISC","uri":null,"id":null}]	com.atlassian.healthcheck.atlassian-healthcheck pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.oauth.consumer"}]	Apps	[]	1622534279632	109	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.oauth.consumer","type":"MISC","uri":null,"id":null}]	com.atlassian.oauth.consumer pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.oauth.serviceprovider"}]	Apps	[]	1622534279726	110	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.oauth.serviceprovider","type":"MISC","uri":null,"id":null}]	com.atlassian.oauth.serviceprovider pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.oauth.serviceprovider.sal"}]	Apps	[]	1622534279744	111	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.oauth.serviceprovider.sal","type":"MISC","uri":null,"id":null}]	com.atlassian.oauth.serviceprovider.sal pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugin.atlassian-spring-scanner-runtime"}]	Apps	[]	1622534279746	112	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugin.atlassian-spring-scanner-runtime","type":"MISC","uri":null,"id":null}]	com.atlassian.plugin.atlassian-spring-scanner-runtime pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.atlassian-chaperone"}]	Apps	[]	1622534279765	113	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.atlassian-chaperone","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.atlassian-chaperone pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.atlassian-client-resource"}]	Apps	[]	1622534279768	114	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.atlassian-client-resource","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.atlassian-client-resource pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.atlassian-clientside-extensions-runtime"}]	Apps	[]	1622534279789	115	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.atlassian-clientside-extensions-runtime","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.atlassian-clientside-extensions-runtime pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.atlassian-plugins-webresource-plugin"}]	Apps	[]	1622534279806	116	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.atlassian-plugins-webresource-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.atlassian-plugins-webresource-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.atlassian-plugins-webresource-rest"}]	Apps	[]	1622534279819	117	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.atlassian-plugins-webresource-rest","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.atlassian-plugins-webresource-rest pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.atlassian-remote-event-consumer-plugin"}]	Apps	[]	1622534279832	118	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.atlassian-remote-event-consumer-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.atlassian-remote-event-consumer-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.atlassian-remote-event-producer-plugin"}]	Apps	[]	1622534279852	119	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.atlassian-remote-event-producer-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.atlassian-remote-event-producer-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.authentication.atlassian-authentication-plugin"}]	Apps	[]	1622534279980	120	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.authentication.atlassian-authentication-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.authentication.atlassian-authentication-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.browser.metrics.browser-metrics-plugin"}]	Apps	[]	1622534280090	121	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.browser.metrics.browser-metrics-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.browser.metrics.browser-metrics-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.issue-status-plugin"}]	Apps	[]	1622534280093	122	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.issue-status-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.issue-status-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.jquery"}]	Apps	[]	1622534280095	123	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.jquery","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.jquery pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.less-transformer-plugin"}]	Apps	[]	1622534280111	124	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.less-transformer-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.less-transformer-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.shortcuts.atlassian-shortcuts-plugin"}]	Apps	[]	1622534280142	125	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.shortcuts.atlassian-shortcuts-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.shortcuts.atlassian-shortcuts-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.plugins.static-assets-url"}]	Apps	[]	1622534280186	126	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.plugins.static-assets-url","type":"MISC","uri":null,"id":null}]	com.atlassian.plugins.static-assets-url pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.prettyurls.atlassian-pretty-urls-plugin"}]	Apps	[]	1622534280199	127	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.prettyurls.atlassian-pretty-urls-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.prettyurls.atlassian-pretty-urls-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.soy.soy-template-plugin"}]	Apps	[]	1622534280214	128	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.soy.soy-template-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.soy.soy-template-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.stash.plugins.stash-remote-event-bitbucket-server-spi"}]	Apps	[]	1622534280230	129	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.stash.plugins.stash-remote-event-bitbucket-server-spi","type":"MISC","uri":null,"id":null}]	com.atlassian.stash.plugins.stash-remote-event-bitbucket-server-spi pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.troubleshooting.plugin-bitbucket"}]	Apps	[]	1622534280321	130	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.troubleshooting.plugin-bitbucket","type":"MISC","uri":null,"id":null}]	com.atlassian.troubleshooting.plugin-bitbucket pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.atlassian.webhooks.atlassian-webhooks-plugin"}]	Apps	[]	1622534280351	131	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.atlassian.webhooks.atlassian-webhooks-plugin","type":"MISC","uri":null,"id":null}]	com.atlassian.webhooks.atlassian-webhooks-plugin pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"com.springsource.net.jcip.annotations-1.0.0"}]	Apps	[]	1622534280359	132	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"com.springsource.net.jcip.annotations-1.0.0","type":"MISC","uri":null,"id":null}]	com.springsource.net.jcip.annotations-1.0.0 pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"tac.bitbucket.languages.de_DE"}]	Apps	[]	1622534280370	133	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"tac.bitbucket.languages.de_DE","type":"MISC","uri":null,"id":null}]	tac.bitbucket.languages.de_de pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"tac.bitbucket.languages.fr_FR"}]	Apps	[]	1622534280403	134	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"tac.bitbucket.languages.fr_FR","type":"MISC","uri":null,"id":null}]	tac.bitbucket.languages.fr_fr pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
PluginEnabledEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"tac.bitbucket.languages.ja_JP"}]	Apps	[]	1622534280453	135	BASE	System	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"tac.bitbucket.languages.ja_JP","type":"MISC","uri":null,"id":null}]	tac.bitbucket.languages.ja_jp pluginenabledevent apps system	\N	\N	\N	\N	-1	System	system
BaseUrlChangedEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"BASE_URL"},{"name":"details","value":"{\\"new\\":\\"http://localhost:7990\\",\\"old\\":null}"}]	Global administration	[]	1622534406218	136	BASE	Browser	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"BASE_URL","type":"MISC","uri":null,"id":null}]	base_url baseurlchangedevent global administration anonymous 172.18.0.1	\N	\N	172.18.0.1	http://localhost:7990	-2	Anonymous	user
LicenseChangedEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"System"}]	Global administration	[]	1622534406290	137	BASE	Browser	07c63f67-7a40-4518-a341-368ba29082ba	\N	\N	[]	licensechangedevent global administration anonymous 172.18.0.1	\N	\N	172.18.0.1	http://localhost:7990	-2	Anonymous	user
UserCreatedEvent	USER_MANAGEMENT	[{"name":"target","value":"admin"}]	Users and groups	[]	1622534446835	138	BASE	Browser	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"admin","type":"MISC","uri":null,"id":null}]	admin usercreatedevent users and groups anonymous 172.18.0.1	\N	\N	172.18.0.1	http://localhost:7990	-2	Anonymous	user
GroupMembershipsCreatedEvent	USER_MANAGEMENT	[{"name":"target","value":"stash-users"},{"name":"details","value":"{\\"entities\\":[\\"admin\\"],\\"membership\\":\\"GROUP_USER\\"}"}]	Users and groups	[]	1622534447149	139	BASE	Browser	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"stash-users","type":"MISC","uri":null,"id":null}]	stash-users groupmembershipscreatedevent users and groups anonymous 172.18.0.1	\N	\N	172.18.0.1	http://localhost:7990	-2	Anonymous	user
GlobalPermissionGrantRequestedEvent	PERMISSIONS	[{"name":"details","value":"{\\"permission\\":\\"SYS_ADMIN\\",\\"user\\":\\"admin\\"}"},{"name":"target","value":"Global"}]	Permissions	[]	1622534447403	140	BASE	Browser	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"Global","type":"MISC","uri":null,"id":null}]	global globalpermissiongrantrequestedevent permissions anonymous 172.18.0.1	\N	\N	172.18.0.1	http://localhost:7990	-2	Anonymous	user
GlobalPermissionGrantedEvent	PERMISSIONS	[{"name":"details","value":"{\\"permission\\":\\"SYS_ADMIN\\",\\"user\\":\\"admin\\"}"},{"name":"target","value":"Global"}]	Permissions	[]	1622534447446	141	BASE	Browser	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"Global","type":"MISC","uri":null,"id":null}]	global globalpermissiongrantedevent permissions anonymous 172.18.0.1	\N	\N	172.18.0.1	http://localhost:7990	-2	Anonymous	user
ApplicationSetupEvent	GLOBAL_CONFIG_AND_ADMINISTRATION	[{"name":"details","value":"{\\"new\\":true,\\"old\\":false}"},{"name":"target","value":"SERVER_IS_SETUP"}]	Global administration	[]	1622534447519	142	BASE	Browser	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"SERVER_IS_SETUP","type":"MISC","uri":null,"id":null}]	server_is_setup applicationsetupevent global administration anonymous 172.18.0.1	\N	\N	172.18.0.1	http://localhost:7990	-2	Anonymous	user
ProjectCreationRequestedEvent	LOCAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"ODSPIPELINETEST"}]	Projects	[]	1622534518013	143	BASE	Browser	07c63f67-7a40-4518-a341-368ba29082ba	0	PROJECT	[{"name":"ODSPIPELINETEST","type":"PROJECT","uri":null,"id":"0"}]	odspipelinetest projectcreationrequestedevent projects admin 172.18.0.1	\N	\N	172.18.0.1	http://localhost:7990	1	admin	NORMAL
ProjectPermissionGrantedEvent	PERMISSIONS	[{"name":"target","value":"ODSPIPELINETEST"},{"name":"details","value":"{\\"permission\\":\\"PROJECT_ADMIN\\",\\"user\\":\\"admin\\"}"}]	Permissions	[]	1622534518114	144	BASE	Browser	07c63f67-7a40-4518-a341-368ba29082ba	1	PROJECT	[{"name":"ODSPIPELINETEST","type":"PROJECT","uri":null,"id":"1"}]	odspipelinetest projectpermissiongrantedevent permissions admin 172.18.0.1	\N	\N	172.18.0.1	http://localhost:7990	1	admin	NORMAL
ProjectCreatedEvent	LOCAL_CONFIG_AND_ADMINISTRATION	[{"name":"target","value":"ODSPIPELINETEST"}]	Projects	[]	1622534518127	145	BASE	Browser	07c63f67-7a40-4518-a341-368ba29082ba	1	PROJECT	[{"name":"ODSPIPELINETEST","type":"PROJECT","uri":null,"id":"1"}]	odspipelinetest projectcreatedevent projects admin 172.18.0.1	\N	\N	172.18.0.1	http://localhost:7990	1	admin	NORMAL
BranchModelConfigurationCreatedEvent	LOCAL_CONFIG_AND_ADMINISTRATION	[{"name":"details","value":"{\\"development\\":{\\"refId\\":null,\\"useDefault\\":true},\\"types\\":[{\\"id\\":\\"BUGFIX\\",\\"displayName\\":\\"BUGFIX\\",\\"enabled\\":true,\\"prefix\\":\\"bugfix/\\"},{\\"id\\":\\"FEATURE\\",\\"displayName\\":\\"FEATURE\\",\\"enabled\\":true,\\"prefix\\":\\"feature/\\"},{\\"id\\":\\"HOTFIX\\",\\"displayName\\":\\"HOTFIX\\",\\"enabled\\":true,\\"prefix\\":\\"hotfix/\\"},{\\"id\\":\\"RELEASE\\",\\"displayName\\":\\"RELEASE\\",\\"enabled\\":true,\\"prefix\\":\\"release/\\"}],\\"scope\\":{\\"type\\":\\"PROJECT\\",\\"resourceId\\":1}}"},{"name":"target","value":"ODSPIPELINETEST"}]	Projects	[]	1622534518357	146	BASE	Browser	07c63f67-7a40-4518-a341-368ba29082ba	1	PROJECT	[{"name":"ODSPIPELINETEST","type":"PROJECT","uri":null,"id":"1"}]	odspipelinetest branchmodelconfigurationcreatedevent projects admin 172.18.0.1	\N	\N	172.18.0.1	http://localhost:7990	1	admin	NORMAL
AccessTokenCreatedEvent	USER_MANAGEMENT	[{"name":"details","value":"{\\"id\\":\\"754995254613\\",\\"tokenOwner\\":{\\"id\\":1,\\"name\\":\\"admin\\",\\"slug\\":\\"admin\\"},\\"name\\":\\"admin\\",\\"permissions\\":[\\"PROJECT_ADMIN\\",\\"REPO_ADMIN\\"]}"},{"name":"target","value":"GLOBAL"}]	Users and groups	[]	1622534544252	147	BASE	Browser	07c63f67-7a40-4518-a341-368ba29082ba	\N	MISC	[{"name":"GLOBAL","type":"MISC","uri":null,"id":null}]	global accesstokencreatedevent users and groups admin 172.18.0.1	\N	\N	172.18.0.1	http://localhost:7990	1	admin	NORMAL
\.


--
-- Data for Name: AO_CFE8FA_BUILD_STATUS; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_CFE8FA_BUILD_STATUS" ("CSID", "DATE_ADDED", "DESCRIPTION", "ID", "KEY", "NAME", "STATE", "URL") FROM stdin;
\.


--
-- Data for Name: AO_D6A508_IMPORT_JOB; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_D6A508_IMPORT_JOB" ("CREATED_DATE", "ID", "USER_ID") FROM stdin;
\.


--
-- Data for Name: AO_D6A508_REPO_IMPORT_TASK; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_D6A508_REPO_IMPORT_TASK" ("CLONE_URL", "CREATED_DATE", "EXTERNAL_REPO_NAME", "FAILURE_TYPE", "ID", "IMPORT_JOB_ID", "LAST_UPDATED", "REPOSITORY_ID", "STATE") FROM stdin;
\.


--
-- Data for Name: AO_E5A814_ACCESS_TOKEN; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_E5A814_ACCESS_TOKEN" ("CREATED_DATE", "HASHED_TOKEN", "LAST_AUTHENTICATED", "NAME", "TOKEN_ID", "USER_ID") FROM stdin;
2021-06-01 08:02:24.219	{PKCS5S2}bhfsoAsIB7QuBAZaf1iv8D3e6/YKNexqaYtWGni+6N/P9AG3t4y3E7dcuaLub62G	\N	admin	754995254613	1
\.


--
-- Data for Name: AO_E5A814_ACCESS_TOKEN_PERM; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_E5A814_ACCESS_TOKEN_PERM" ("FK_ACCESS_TOKEN_ID", "ID", "PERMISSION") FROM stdin;
754995254613	1	4
754995254613	2	8
\.


--
-- Data for Name: AO_ED669C_SEEN_ASSERTIONS; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_ED669C_SEEN_ASSERTIONS" ("ASSERTION_ID", "EXPIRY_TIMESTAMP", "ID") FROM stdin;
\.


--
-- Data for Name: AO_F4ED3A_ADD_ON_PROPERTY_AO; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_F4ED3A_ADD_ON_PROPERTY_AO" ("ID", "PLUGIN_KEY", "PRIMARY_KEY", "PROPERTY_KEY", "VALUE") FROM stdin;
\.


--
-- Data for Name: AO_FB71B4_SSH_PUBLIC_KEY; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public."AO_FB71B4_SSH_PUBLIC_KEY" ("ENTITY_ID", "KEY_MD5", "KEY_TEXT", "KEY_TYPE", "LABEL", "LABEL_LOWER", "USER_ID") FROM stdin;
\.


--
-- Data for Name: app_property; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.app_property (prop_key, prop_value) FROM stdin;
instance.application.mode	default
server.id	B3WW-UDNX-WQWG-ZXPI
instance.home	/var/atlassian/application-data/bitbucket/shared
last.os	LINUX
locale	en_US
instance.name	Bitbucket
instance.url	http://localhost:7990
license	AAACLg0ODAoPeNqNVEtv4jAQvudXRNpbpUSEx6FIOQBxW3ZZiCB0V1WllXEG8DbYke3A8u/XdUgVQ\r\nyg9ZvLN+HuM/e1BUHdGlNvuuEHQ73X73Y4bR4nbbgU9ZwFiD2IchcPH+8T7vXzuej9eXp68YSv45\r\nUwoASYhOeYwxTsIE7RIxtNHhwh+SP3a33D0XnntuxHsIeM5CIdwtvYxUXQPoRIF6KaC0FUGVlEB3\r\nv0hOAOWYiH9abFbgZith3i34nwOO65gsAGmZBhUbNC/nIpjhBWEcefJWelzqIDPWz/OtjmXRYv2X\r\nyqwnwueFkT57x8e4cLmbCD1QnX0UoKQoRc4EUgiaK4oZ2ECUrlZeay75sLNs2JDmZtWR8oPCfWZG\r\nwHAtjzXgIo0SqmZiKYJmsfz8QI5aI+zApuq6fqJKVPAMCPnNpk4LPW6kBWgkZb+kQAzzzS2g6Dnt\r\ne69Tqvsr4SOskIqEFOeggz1v4zrHbr0yLJR8rU64FpQpVtBy1mZxM4CnHC9Faf8tKMnTF1AiXORF\r\nixyQaWto3RZ+ncWLXtMg6EnKZZRpmQNb2R8tnJXFulCfXmXLry7TrHBWn2HNVyH8WYxj9AzmsxiN\r\nL/R88Xg6rA1lVs4QpO5titxhplJcCY2mFFZLutAZVhKipm15/VhJx36YVqyN8YP7IaGC1+lwnJ7Q\r\n5pJpNmxk5hP3qovutY8Pi4E2WIJ59esnr1p+T6eD67teBVCHf+ga+ho4/4D9YItZDAsAhQ5qQ6pA\r\nSJ+SA7YG9zthbLxRoBBEwIURQr5Zy1B8PonepyLz3UhL7kMVEs=X02q6
setup.completed	true
last.licensed.user.count	1
\.


--
-- Data for Name: bb_alert; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_alert (details_json, id, issue_component_id, issue_id, issue_severity, node_name, node_name_lower, "timestamp", trigger_module, trigger_plugin_key, trigger_plugin_key_lower, trigger_plugin_version) FROM stdin;
\.


--
-- Data for Name: bb_announcement_banner; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_announcement_banner (id, enabled, audience, message) FROM stdin;
\.


--
-- Data for Name: bb_attachment; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_attachment (id, repository_id, filename) FROM stdin;
\.


--
-- Data for Name: bb_attachment_metadata; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_attachment_metadata (attachment_id, metadata) FROM stdin;
\.


--
-- Data for Name: bb_clusteredjob; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_clusteredjob (job_id, job_runner_key, sched_type, interval_millis, first_run, cron_expression, time_zone, next_run, version, parameters) FROM stdin;
TruncateAlertsJobRunner	com.atlassian.diagnostics.internal.DefaultMonitoringService$TruncateAlertsJobRunner	1	86400000	2021-06-02 07:56:10.883	\N	\N	2021-06-02 07:56:10.883	1	\N
analytics-collection	com.atlassian.plugins.authentication.impl.analytics.StatisticsCollectionService	0	\N	\N	0 0 23 * * ?	\N	2021-06-01 23:00:00	1	\N
assertionId-cleanup	com.atlassian.plugins.authentication.impl.web.saml.SamlAssertionValidationService	1	3600000	2021-06-01 08:57:30.896	\N	\N	2021-06-01 08:57:30.896	1	\N
OidcDiscoveryRefresh	com.atlassian.plugins.authentication.impl.web.oidc.OidcDiscoveryRefreshJob-refresh	0	\N	\N	0 0 1 * * ?	\N	2021-06-02 01:00:00	1	\N
com.atlassian.audit.retention.RetentionJobRunner	com.atlassian.audit.retention.RetentionJobRunner	0	\N	\N	0 0 0 1/1 * ? *	\N	2021-06-02 00:00:00	1	\N
com.atlassian.audit.schedule.db.limit.DbLimiterJobRunner	com.atlassian.audit.schedule.db.limit.DbLimiterJobRunner	1	3600000	2021-06-01 08:57:45.985	\N	\N	2021-06-01 08:57:45.985	1	\N
com.atlassian.diagnostics.internal.analytics.DailyAlertAnalyticsJob	DailyAlertAnalyticsJob	0	\N	\N	0 19 * * * ?	\N	2021-06-01 08:19:00	1	\N
GroupCleanUpJob	com.atlassian.stash.internal.user.DefaultUserAdminService$GroupCleanUpJob	1	21600000	2021-06-01 13:58:00.994	\N	\N	2021-06-01 13:58:00.994	1	\N
UserCleanupJob	com.atlassian.stash.internal.user.DefaultUserAdminService$UserCleanupJob	1	21600000	2021-06-01 13:58:00.998	\N	\N	2021-06-01 13:58:00.998	1	\N
CleanupExpiredRememberMeTokensJob	com.atlassian.stash.internal.auth.RememberMeTokenCleanupScheduler$CleanupExpiredRememberMeTokensJob	1	18000000	2021-06-01 12:58:01.007	\N	\N	2021-06-01 12:58:01.007	1	\N
HookScriptService.cleanupJob	com.atlassian.bitbucket.hook.script.HookScriptService.cleanupJob	1	1440000	2021-06-01 07:58:01.028	\N	\N	2021-06-01 08:22:01.035	2	\N
CleanupEmptyRescopesJob	com.atlassian.stash.internal.pull.rescope.DefaultRescopeProcessor$CleanupEmptyRescopesJob	1	1800000	2021-06-01 08:28:01.035	\N	\N	2021-06-01 08:28:01.035	1	\N
HistoryCleanupJob	com.atlassian.bitbucket.internal.ratelimit.history.HistoryCleanupJob	1	86400000	2021-06-01 08:13:01.086	\N	\N	2021-06-01 08:13:01.086	1	\N
applink-status-analytics-job	com.atlassian.applinks.analytics.ApplinkStatusJob	1	86400000	\N	\N	\N	2021-06-02 08:00:47.739	2	\N
WebhookAnalyticsJobRunner	com.atlassian.bitbucket.internal.webhook.WebhookAnalyticsService$WebhookAnalyticsJobRunner	1	86400000	2021-06-01 08:30:47.761	\N	\N	2021-06-01 08:30:47.761	1	\N
AccessTokenAnalyticsJob	com.atlassian.bitbucket.internal.accesstokens.analytics.AccessTokenAnalyticsJob	1	86400000	2021-06-02 08:00:47.802	\N	\N	2021-06-02 08:00:47.802	1	\N
InsightReportCleanupJob	com.atlassian.bitbucket.internal.codeinsights.report.InsightReportCleanupJob	1	86400000	2021-06-02 08:00:48.639	\N	\N	2021-06-02 08:00:48.639	1	\N
LocalPluginLicenseNotificationJob-job	LocalPluginLicenseNotificationJob-runner	1	86400000	2021-06-01 08:00:48.802	\N	\N	2021-06-02 08:00:48.813	2	\N
PluginRequestCheckJob-job	PluginRequestCheckJob-runner	1	3600000	2021-06-01 08:00:48.816	\N	\N	2021-06-01 09:00:48.83	2	\N
PluginUpdateCheckJob-job	PluginUpdateCheckJob-runner	1	86400000	2021-06-02 01:52:02.327	\N	\N	2021-06-02 01:52:02.327	1	\N
InstanceTopologyJob-job	InstanceTopologyJob-runner	1	86400000	2021-06-01 14:02:26.778	\N	\N	2021-06-01 14:02:26.778	1	\N
Service Provider Session Remover	com.atlassian.oauth.serviceprovider.internal.ExpiredSessionRemover	1	28800000	2021-06-01 16:00:50.903	\N	\N	2021-06-01 16:00:50.903	1	\N
webhooks.history.daily.cleanup.job	webhooks.history.daily.cleanup.runner	0	\N	\N	0 22 4 1/1 * ? *	\N	2021-06-02 04:22:00	1	\N
05d99db5-9dd7-46d0-9049-4d060ca4cbe4	SEARCH_HEALTH_CHECK	1	0	2021-06-01 08:06:11.435	\N	\N	2021-06-01 08:06:11.435	1	\N
com.atlassian.crowd.manager.directory.monitor.DirectoryMonitorRefresherStarter-job	com.atlassian.crowd.manager.directory.monitor.DirectoryMonitorRefresherJob-runner	1	120000	\N	\N	\N	2021-06-01 08:04:11.187	4	\N
IndexingAnalyticsJob	com.atlassian.bitbucket.internal.search.indexing.analytics.IndexingAnalyticsJob	1	604800000	2021-06-01 08:02:37.674	\N	\N	2021-06-08 08:02:37.674	2	\N
SearchAnalyticsJob	com.atlassian.bitbucket.internal.search.common.analytics.SearchAnalyticsJob	1	604800000	2021-06-01 08:02:37.74	\N	\N	2021-06-08 08:02:37.74	2	\N
Repository	com.atlassian.stash.internal.notification.repository.batch.RepositoryBatchSender.Repository	1	60000	2021-06-01 08:01:49.655	\N	\N	2021-06-01 08:04:49.661	4	\N
PullRequest	com.atlassian.stash.internal.notification.pull.activity.PullRequestBatchSender.PullRequest	1	60000	2021-06-01 08:01:49.672	\N	\N	2021-06-01 08:04:49.672	4	\N
a1f5755c-37d1-44a1-9756-38c16ab71483	com.atlassian.bitbucket.internal.search.indexing.jobs.StartupChecksJob	1	0	2021-06-01 08:08:06.115	\N	\N	2021-06-01 08:08:06.115	1	\\xaced000573720037636f6d2e676f6f676c652e636f6d6d6f6e2e636f6c6c6563742e496d6d757461626c6542694d61702453657269616c697a6564466f726d000000000000000002000078720035636f6d2e676f6f676c652e636f6d6d6f6e2e636f6c6c6563742e496d6d757461626c654d61702453657269616c697a6564466f726d00000000000000000200025b00046b6579737400135b4c6a6176612f6c616e672f4f626a6563743b5b000676616c75657371007e00027870757200135b4c6a6176612e6c616e672e4f626a6563743b90ce589f1073296c020000787000000001740007726574726965737571007e000400000001737200116a6176612e6c616e672e496e746567657212e2a0a4f781873802000149000576616c7565787200106a6176612e6c616e672e4e756d62657286ac951d0b94e08b020000787000000003
\.


--
-- Data for Name: bb_cmt_disc_comment_activity; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_cmt_disc_comment_activity (activity_id, comment_id, comment_action) FROM stdin;
\.


--
-- Data for Name: bb_comment; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_comment (id, author_id, comment_text, created_timestamp, entity_version, thread_id, updated_timestamp, resolved_timestamp, resolver_id, severity, state) FROM stdin;
\.


--
-- Data for Name: bb_comment_parent; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_comment_parent (comment_id, parent_id) FROM stdin;
\.


--
-- Data for Name: bb_comment_thread; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_comment_thread (id, commentable_id, commentable_type, created_timestamp, entity_version, updated_timestamp, diff_type, file_type, from_hash, from_path, is_orphaned, line_number, line_type, to_hash, to_path) FROM stdin;
\.


--
-- Data for Name: bb_data_store; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_data_store (id, ds_path, ds_uuid) FROM stdin;
\.


--
-- Data for Name: bb_git_pr_cached_merge; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_git_pr_cached_merge (id, from_hash, to_hash, merge_type) FROM stdin;
\.


--
-- Data for Name: bb_git_pr_common_ancestor; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_git_pr_common_ancestor (id, from_hash, to_hash, ancestor_hash) FROM stdin;
\.


--
-- Data for Name: bb_hook_script; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_hook_script (id, hook_version, hook_size, hook_type, created_timestamp, updated_timestamp, hook_hash, hook_name, plugin_key, hook_description) FROM stdin;
\.


--
-- Data for Name: bb_hook_script_config; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_hook_script_config (id, script_id, scope_id, scope_type) FROM stdin;
\.


--
-- Data for Name: bb_hook_script_trigger; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_hook_script_trigger (config_id, trigger_id) FROM stdin;
\.


--
-- Data for Name: bb_integrity_event; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_integrity_event (event_key, event_timestamp, event_node) FROM stdin;
\.


--
-- Data for Name: bb_job; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_job (end_timestamp, id, initiator_id, node_id, progress_percentage, progress_message, start_timestamp, state, type, updated_timestamp, entity_version) FROM stdin;
\.


--
-- Data for Name: bb_job_message; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_job_message (created_timestamp, id, job_id, severity, subject, text) FROM stdin;
\.


--
-- Data for Name: bb_label; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_label (id, name) FROM stdin;
\.


--
-- Data for Name: bb_label_mapping; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_label_mapping (id, label_id, labelable_id, labelable_type) FROM stdin;
\.


--
-- Data for Name: bb_mirror_content_hash; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_mirror_content_hash (repository_id, updated_timestamp, hash) FROM stdin;
\.


--
-- Data for Name: bb_mirror_metadata_hash; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_mirror_metadata_hash (repository_id, updated_timestamp, hash) FROM stdin;
\.


--
-- Data for Name: bb_pr_comment_activity; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_pr_comment_activity (activity_id, comment_id, comment_action) FROM stdin;
\.


--
-- Data for Name: bb_pr_commit; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_pr_commit (pr_id, commit_id) FROM stdin;
\.


--
-- Data for Name: bb_pr_part_status_weight; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_pr_part_status_weight (status_id, status_weight) FROM stdin;
0	100
1	300
2	200
\.


--
-- Data for Name: bb_pr_reviewer_added; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_pr_reviewer_added (activity_id, user_id) FROM stdin;
\.


--
-- Data for Name: bb_pr_reviewer_removed; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_pr_reviewer_removed (activity_id, user_id) FROM stdin;
\.


--
-- Data for Name: bb_pr_reviewer_upd_activity; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_pr_reviewer_upd_activity (activity_id) FROM stdin;
\.


--
-- Data for Name: bb_proj_merge_config; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_proj_merge_config (id, project_id, scm_id, default_strategy_id, commit_summaries) FROM stdin;
\.


--
-- Data for Name: bb_proj_merge_strategy; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_proj_merge_strategy (config_id, strategy_id) FROM stdin;
\.


--
-- Data for Name: bb_project_alias; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_project_alias (id, project_id, namespace, project_key, created_timestamp) FROM stdin;
\.


--
-- Data for Name: bb_repo_merge_config; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_repo_merge_config (id, repository_id, default_strategy_id, commit_summaries) FROM stdin;
\.


--
-- Data for Name: bb_repo_merge_strategy; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_repo_merge_strategy (config_id, strategy_id) FROM stdin;
\.


--
-- Data for Name: bb_repository_alias; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_repository_alias (id, repository_id, project_namespace, project_key, slug, created_timestamp) FROM stdin;
\.


--
-- Data for Name: bb_rl_reject_counter; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_rl_reject_counter (id, user_id, interval_start, reject_count) FROM stdin;
\.


--
-- Data for Name: bb_rl_user_settings; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_rl_user_settings (id, user_id, capacity, fill_rate, whitelisted) FROM stdin;
\.


--
-- Data for Name: bb_scm_merge_config; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_scm_merge_config (id, scm_id, default_strategy_id, commit_summaries) FROM stdin;
\.


--
-- Data for Name: bb_scm_merge_strategy; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_scm_merge_strategy (config_id, strategy_id) FROM stdin;
\.


--
-- Data for Name: bb_suggestion_group; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_suggestion_group (comment_id, state, applied_index) FROM stdin;
\.


--
-- Data for Name: bb_thread_root_comment; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_thread_root_comment (thread_id, comment_id) FROM stdin;
\.


--
-- Data for Name: bb_user_dark_feature; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.bb_user_dark_feature (id, user_id, is_enabled, feature_key) FROM stdin;
\.


--
-- Data for Name: changeset; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.changeset (id, author_timestamp) FROM stdin;
\.


--
-- Data for Name: cs_attribute; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cs_attribute (cs_id, att_name, att_value) FROM stdin;
\.


--
-- Data for Name: cs_indexer_state; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cs_indexer_state (indexer_id, repository_id, last_run) FROM stdin;
\.


--
-- Data for Name: cs_repo_membership; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cs_repo_membership (cs_id, repository_id) FROM stdin;
\.


--
-- Data for Name: current_app; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.current_app (id, application_id, public_key_base64, private_key_base64) FROM stdin;
1	AC1200050179C69483B5593411B6A29F	MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjNrXgLEBiztZbqoudyKCIDCxAfNJcvzmleRlC5wsame8jPn2Ufnek9ovxCT4s9tQ5HU9XLo8oG8kKsmz+U/DXq3dimxzi6lqV1HJGxzOb6spV4jt82rEXHZqBXY9pp/odvOuZ73WGC4HRRHfbcVOu9Nl9brOkVTSAHABS4vjYR1fGQqsQ3CmVc30KuvYaA9d6zJIGdX2YwfCwXkTGZtLyUSCRJrLp/dT5qpuA61akfIkFlCN/hAuFiEQDW8PYDTcL0+0F0FRLoS2i7GWb7nv5f+xPdgYRNzhfk0oqvuzN/5NaFmQipVLaG24zSoXk9/G3NP2b744s1jpkHtGqQUDqQIDAQAB	MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCM2teAsQGLO1luqi53IoIgMLEB80ly/OaV5GULnCxqZ7yM+fZR+d6T2i/EJPiz21DkdT1cujygbyQqybP5T8Nerd2KbHOLqWpXUckbHM5vqylXiO3zasRcdmoFdj2mn+h2865nvdYYLgdFEd9txU6702X1us6RVNIAcAFLi+NhHV8ZCqxDcKZVzfQq69hoD13rMkgZ1fZjB8LBeRMZm0vJRIJEmsun91Pmqm4DrVqR8iQWUI3+EC4WIRANbw9gNNwvT7QXQVEuhLaLsZZvue/l/7E92BhE3OF+TSiq+7M3/k1oWZCKlUtobbjNKheT38bc0/ZvvjizWOmQe0apBQOpAgMBAAECggEADHKEa894mvS+NPzeCVIn3K9g3RLCUWKO//0Efu+oryiGrZCjV2A07qzv3q4DumUts1q29vxQQj9AG3XirSwC0FmeA88MsPFiP+Au3PIBPjYKe71ShdnQC3m9ackmrpRgBy8GoQ1SH6Xlp5FqRLmUeR2322zMN2SfAEHlo2bzy4+a2GluEFRwId6O1r+vf6i4GTFCSGOZPGwJcWKzyUpISs0dPlfEa9fif7akVS2WCSt78YSimvIk7yCsGHTzpEg75tiLTn6+6Yey3YlJDIH6k5rWnFP53IuZ0wurYTDW0TLtWT7oFJ9AYzwmjKjwKg61DDeh5K87IDZax1NiTTIVAQKBgQDA2tKNQYH1bemhZPzjMGq0eXnytXkvnnMryDCXFGDG3RDTvxwX1MISQn9mcNBY8gmDuc5fnSM9XIYdq+s11ufKdFtxOnVSFEQf7Xs3EyMKYW4yDFpeZGpcyXYP4Qg6GKaCqUu3G++ZTaKvGgqprmEDEUe9CY1wnltIx0/M1trb6QKBgQC6+VyAQ54WzTUHPUHwAiBOvoGn6XZGn5IZJ3fzWJHcMn0DFUku8mEoHCDFw4i5c9DL0AG4/qbiZC3HQWZRZ1veeUV0byCrvwNT+NfoVTp8r/qq7K+OZwYu7LDMaar691nNFNccqcxDHLyBNclaLqgF/fP2jPkSi0GGGfy/exzRwQKBgQCQzPwgMVpZo9AybTfvgS/tF/R3RsiZ93d0HRhWp2dOiTeUNT7rqcSZnzI3AWd+ESURsZYBdmO6M9lDOA0f3J8nBJyP9JuYKD1KV64XGRhLOAJcM6g6jVzLFDzACW96531GR8Tg1GnfCkqm/H+bDaIrgnMBvcVkFJJnn7cMDxo+2QKBgG0fZf4x+I0kPOO5u8cA5qwugWtnVTFIpjLqFxa+RXq3OMDY5npw2YVYTUQ+p4hc8KpS/v5iGTId953ILJgr87E3I/MdfHrgI2gZ3qDpRRZKesjRFHDO7gvq9hCHR1PrksyfciB3dRBiMB5VLuvkOQouOflWM3PANIC4oAt9JcfBAoGBALkSbIxHbIZ00K15B7/uZeNdD4Epr5cbk+A4Za/o+umBBfu6caDNAFvXLh6hUGTPDJDqFhnwjOVgOOsujyZzjC+w0DdQQxaqJ0UbZMNj/YqfblKxofTitUW5WchtHXVzb/By4e3116wxhpMbQtO0kcnrq0LlistovgUfABUFNX7s
\.


--
-- Data for Name: cwd_app_dir_default_groups; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_app_dir_default_groups (id, application_mapping_id, group_name) FROM stdin;
\.


--
-- Data for Name: cwd_app_dir_group_mapping; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_app_dir_group_mapping (id, app_dir_mapping_id, application_id, directory_id, group_name) FROM stdin;
\.


--
-- Data for Name: cwd_app_dir_mapping; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_app_dir_mapping (id, application_id, directory_id, list_index, is_allow_all) FROM stdin;
65537	1	32769	0	T
\.


--
-- Data for Name: cwd_app_dir_operation; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_app_dir_operation (app_dir_mapping_id, operation_type) FROM stdin;
65537	CREATE_GROUP
65537	UPDATE_USER
65537	DELETE_GROUP
65537	UPDATE_GROUP_ATTRIBUTE
65537	UPDATE_ROLE
65537	UPDATE_USER_ATTRIBUTE
65537	CREATE_ROLE
65537	UPDATE_GROUP
65537	DELETE_USER
65537	DELETE_ROLE
65537	CREATE_USER
65537	UPDATE_ROLE_ATTRIBUTE
\.


--
-- Data for Name: cwd_app_licensed_user; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_app_licensed_user (id, username, full_name, email, last_active, directory_id, lower_username, lower_full_name, lower_email) FROM stdin;
\.


--
-- Data for Name: cwd_app_licensing; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_app_licensing (id, generated_on, version, application_id, application_subtype, total_users, max_user_limit, total_crowd_users, active) FROM stdin;
\.


--
-- Data for Name: cwd_app_licensing_dir_info; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_app_licensing_dir_info (id, name, directory_id, licensing_summary_id) FROM stdin;
\.


--
-- Data for Name: cwd_application; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_application (id, application_name, lower_application_name, created_date, updated_date, description, application_type, credential, is_active) FROM stdin;
1	crowd-embedded	crowd-embedded	2021-06-01 07:56:40.823	2021-06-01 07:56:42.639	\N	CROWD	{PKCS5S2}dlIub/E/eM0Y4/ujf0fVzNA7Z31EfqIzPLDlNRYuHBAtlkobVZkd52RQPvJz9EAR	T
\.


--
-- Data for Name: cwd_application_address; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_application_address (application_id, remote_address) FROM stdin;
\.


--
-- Data for Name: cwd_application_alias; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_application_alias (id, application_id, user_name, lower_user_name, alias_name, lower_alias_name) FROM stdin;
\.


--
-- Data for Name: cwd_application_attribute; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_application_attribute (application_id, attribute_name, attribute_value) FROM stdin;
1	aggregateMemberships	true
1	atlassian_sha1_applied	true
\.


--
-- Data for Name: cwd_application_saml_config; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_application_saml_config (application_id, assertion_consumer_service, audience, enabled) FROM stdin;
\.


--
-- Data for Name: cwd_directory; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_directory (id, directory_name, lower_directory_name, created_date, updated_date, description, impl_class, lower_impl_class, directory_type, is_active) FROM stdin;
32769	Bitbucket Internal Directory	bitbucket internal directory	2021-06-01 07:56:41.878	2021-06-01 07:56:41.878	Bitbucket Internal Directory	com.atlassian.crowd.directory.InternalDirectory	com.atlassian.crowd.directory.internaldirectory	INTERNAL	T
\.


--
-- Data for Name: cwd_directory_attribute; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_directory_attribute (directory_id, attribute_name, attribute_value) FROM stdin;
32769	user_encryption_method	atlassian-security
\.


--
-- Data for Name: cwd_directory_operation; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_directory_operation (directory_id, operation_type) FROM stdin;
32769	CREATE_GROUP
32769	UPDATE_USER
32769	DELETE_GROUP
32769	UPDATE_GROUP_ATTRIBUTE
32769	UPDATE_ROLE
32769	UPDATE_USER_ATTRIBUTE
32769	CREATE_ROLE
32769	UPDATE_GROUP
32769	CREATE_USER
32769	DELETE_USER
32769	DELETE_ROLE
32769	UPDATE_ROLE_ATTRIBUTE
\.


--
-- Data for Name: cwd_granted_perm; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_granted_perm (id, created_date, permission_id, group_name, app_dir_mapping_id) FROM stdin;
\.


--
-- Data for Name: cwd_group; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_group (id, group_name, lower_group_name, created_date, updated_date, description, group_type, directory_id, is_active, is_local, external_id) FROM stdin;
98305	stash-users	stash-users	2021-06-01 07:56:43.245	2021-06-01 07:56:43.245	\N	GROUP	32769	T	F	\N
\.


--
-- Data for Name: cwd_group_admin_group; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_group_admin_group (id, group_id, target_group_id) FROM stdin;
\.


--
-- Data for Name: cwd_group_admin_user; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_group_admin_user (id, user_id, target_group_id) FROM stdin;
\.


--
-- Data for Name: cwd_group_attribute; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_group_attribute (id, group_id, directory_id, attribute_name, attribute_value, attribute_lower_value) FROM stdin;
\.


--
-- Data for Name: cwd_membership; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_membership (id, parent_id, child_id, membership_type, group_type, parent_name, lower_parent_name, child_name, lower_child_name, directory_id, created_date) FROM stdin;
196609	98305	131073	GROUP_USER	GROUP	stash-users	stash-users	admin	admin	32769	2021-06-01 08:00:47.026
\.


--
-- Data for Name: cwd_property; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_property (property_key, property_name, property_value) FROM stdin;
\.


--
-- Data for Name: cwd_tombstone; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_tombstone (id, tombstone_type, tombstone_timestamp, application_id, directory_id, entity_name, parent) FROM stdin;
\.


--
-- Data for Name: cwd_user; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_user (id, user_name, lower_user_name, created_date, updated_date, first_name, lower_first_name, last_name, lower_last_name, display_name, lower_display_name, email_address, lower_email_address, directory_id, credential, is_active, external_id) FROM stdin;
131073	admin	admin	2021-06-01 08:00:46.763	2021-06-01 08:00:57.066			admin	admin	admin	admin	test@opendevstack.org	test@opendevstack.org	32769	{PKCS5S2}W5AVQFa6OS2WjGMFs8Jd7s3YAdo1LVad7EYEfVsnLwCc4DoikUn+w+eMuhfg/eGl	T	5ac5b8c5-4e86-4219-97a8-95cd7ae4637c
\.


--
-- Data for Name: cwd_user_attribute; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_user_attribute (id, user_id, directory_id, attribute_name, attribute_value, attribute_lower_value, attribute_numeric_value) FROM stdin;
163841	131073	32769	requiresPasswordChange	false	false	\N
163842	131073	32769	passwordLastChanged	1622534446778	1622534446778	1622534446778
163843	131073	32769	lastAuthenticationTimestamp	1622534457032	1622534457032	1622534457032
\.


--
-- Data for Name: cwd_user_credential_record; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_user_credential_record (id, user_id, password_hash, list_index) FROM stdin;
\.


--
-- Data for Name: cwd_webhook; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.cwd_webhook (id, endpoint_url, application_id, token, oldest_failure_date, failures_since_last_success) FROM stdin;
\.


--
-- Data for Name: databasechangelog; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.databasechangelog (id, author, filename, dateexecuted, orderexecuted, exectype, md5sum, description, comments, tag, liquibase, contexts, labels, deployment_id) FROM stdin;
STASHDEV-7910-1	jhinch	liquibase/r3_4/bootstrap-upgrade.xml	2021-06-01 07:59:18.298599	1	EXECUTED	8:21f2b187815305259479d91bc4983b1d	createTable tableName=app_property	Create the 'app_property' table, only if it hasn't already been created	\N	3.6.1	\N	\N	2534357482
BBSDEV-17340-1	bturner	liquibase/r6_0/bootstrap-upgrade.xml	2021-06-01 07:59:18.312264	2	EXECUTED	8:ea223f4d97081f0c4770d48bab0d27b1	createTable tableName=bb_data_store		\N	3.6.1	\N	\N	2534357482
initial-schema-01	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.32416	3	EXECUTED	8:22bada8ae50eb9644324344154606b2b	createTable tableName=cwd_application		\N	3.6.1	production	\N	2534357482
initial-schema-02	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.333445	4	EXECUTED	8:bc4bb8e3065127ad7c567acd5691968e	createTable tableName=cwd_directory		\N	3.6.1	production	\N	2534357482
initial-schema-03	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.357501	5	EXECUTED	8:7c504702a27f70115c827052690e5abd	createTable tableName=cwd_app_dir_mapping		\N	3.6.1	production	\N	2534357482
initial-schema-04	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.373153	6	EXECUTED	8:fe5696b3c35cfac7afa2b75d6165387a	createTable tableName=cwd_app_dir_group_mapping		\N	3.6.1	production	\N	2534357482
initial-schema-05	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.387203	7	EXECUTED	8:aee13a8764f376c3243dc241b2f08106	createTable tableName=cwd_app_dir_operation; addPrimaryKey constraintName=SYS_PK_10083, tableName=cwd_app_dir_operation		\N	3.6.1	production	\N	2534357482
initial-schema-06	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.398432	8	EXECUTED	8:6f0059aa727fbaa994381e015541ed48	createTable tableName=cwd_application_address; addPrimaryKey constraintName=SYS_PK_10100, tableName=cwd_application_address		\N	3.6.1	production	\N	2534357482
initial-schema-07	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.409889	9	EXECUTED	8:6935dc9a24194c46afa6d586d4dfe836	createTable tableName=cwd_application_alias		\N	3.6.1	production	\N	2534357482
initial-schema-08	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.422704	10	EXECUTED	8:0c5bb8ef4f50a11ea5c7f3b12a72e211	createTable tableName=cwd_application_attribute; addPrimaryKey constraintName=SYS_PK_10116, tableName=cwd_application_attribute		\N	3.6.1	production	\N	2534357482
initial-schema-09	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.436559	11	EXECUTED	8:85074470c257109825490bd48b52a78f	createTable tableName=cwd_directory_attribute; addPrimaryKey constraintName=SYS_PK_10133, tableName=cwd_directory_attribute		\N	3.6.1	production	\N	2534357482
initial-schema-10	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.453728	12	EXECUTED	8:cd764c6319147bc918627e84d8967f3a	createTable tableName=cwd_directory_operation; addPrimaryKey constraintName=SYS_PK_10137, tableName=cwd_directory_operation		\N	3.6.1	production	\N	2534357482
initial-schema-11	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.463079	13	EXECUTED	8:c88222217eacecaa3d762e0dda557006	createTable tableName=cwd_group		\N	3.6.1	production	\N	2534357482
initial-schema-12	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.47216	14	EXECUTED	8:c8b269a522e972b4cdfdf1e255767d6e	createTable tableName=cwd_group_attribute		\N	3.6.1	production	\N	2534357482
initial-schema-13	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.492948	15	EXECUTED	8:14e6db9188e558580571e2c204609eb0	createTable tableName=cwd_membership		\N	3.6.1	production	\N	2534357482
initial-schema-14	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.517919	16	EXECUTED	8:0405b76442541382a0d935118cce5461	createTable tableName=cwd_property; addPrimaryKey constraintName=SYS_PK_10173, tableName=cwd_property		\N	3.6.1	production	\N	2534357482
initial-schema-15	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.528983	17	EXECUTED	8:7920d65f2808c5ee3fe2bbb0916083f2	createTable tableName=cwd_token		\N	3.6.1	production	\N	2534357482
initial-schema-16	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.539852	18	EXECUTED	8:de92580c7c581b37e00899f63074adba	createTable tableName=cwd_user		\N	3.6.1	production	\N	2534357482
initial-schema-17	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.549989	19	EXECUTED	8:2e560e9ec11844199fad27d26dc22b04	createTable tableName=cwd_user_attribute		\N	3.6.1	production	\N	2534357482
initial-schema-18	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.55768	20	EXECUTED	8:c4a5ce33f4c4333e83d244f59b94e394	createTable tableName=cwd_user_credential_record		\N	3.6.1	production	\N	2534357482
initial-schema-19	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.573214	21	EXECUTED	8:e7c44e7ab22561298380c12f8e313c09	createIndex indexName=IDX_APP_ACTIVE, tableName=cwd_application		\N	3.6.1	production	\N	2534357482
initial-schema-20	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.580968	22	EXECUTED	8:d7ea92c11d421b5a4bcd7e4fbac81c8d	createIndex indexName=IDX_APP_TYPE, tableName=cwd_application		\N	3.6.1	production	\N	2534357482
initial-schema-21	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.587148	23	EXECUTED	8:013866a6461913f1a51d8f0ac76f8a50	createIndex indexName=SYS_IDX_SYS_CT_10094_10096, tableName=cwd_application		\N	3.6.1	production	\N	2534357482
initial-schema-22	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.592532	24	EXECUTED	8:74d0bdf59e3e41e5de481a0c7efbe1a2	createIndex indexName=SYS_IDX_SYS_CT_10109_10112, tableName=cwd_application_alias		\N	3.6.1	production	\N	2534357482
initial-schema-23	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.598533	25	EXECUTED	8:0a4010a80f534ce97982811d2bb20363	createIndex indexName=SYS_IDX_SYS_CT_10110_10113, tableName=cwd_application_alias		\N	3.6.1	production	\N	2534357482
initial-schema-24	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.609008	26	EXECUTED	8:350bae7e52dece6441470254bdd54496	createIndex indexName=IDX_APP_DIR_GROUP_GROUP_DIR, tableName=cwd_app_dir_group_mapping		\N	3.6.1	production	\N	2534357482
initial-schema-25	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.618504	27	EXECUTED	8:1401e38f06dae421430b437314711ea6	createIndex indexName=SYS_IDX_SYS_CT_10070_10072, tableName=cwd_app_dir_group_mapping		\N	3.6.1	production	\N	2534357482
initial-schema-26	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.627329	28	EXECUTED	8:dc202900e71c0ca05783cc9f1ced949a	createIndex indexName=SYS_IDX_SYS_CT_10078_10080, tableName=cwd_app_dir_mapping		\N	3.6.1	production	\N	2534357482
initial-schema-27	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.641457	29	EXECUTED	8:c6f0f900137573decf94be8e1146658b	createIndex indexName=IDX_DIR_ACTIVE, tableName=cwd_directory		\N	3.6.1	production	\N	2534357482
initial-schema-28	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.652978	30	EXECUTED	8:7e5921aac95afb6e6c522f9478377b7f	createIndex indexName=IDX_DIR_L_IMPL_CLASS, tableName=cwd_directory		\N	3.6.1	production	\N	2534357482
initial-schema-29	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.659607	31	EXECUTED	8:3e10aaac3a1833b6a063e449172ffbea	createIndex indexName=IDX_DIR_TYPE, tableName=cwd_directory		\N	3.6.1	production	\N	2534357482
initial-schema-30	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.666058	32	EXECUTED	8:6c4b3489313c33602b0b7f3a35d250a3	createIndex indexName=SYS_IDX_SYS_CT_10128_10130, tableName=cwd_directory		\N	3.6.1	production	\N	2534357482
initial-schema-31	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.673649	33	EXECUTED	8:bf24c461fbe7409030eec0809b5b2378	createIndex indexName=IDX_GROUP_ACTIVE, tableName=cwd_group		\N	3.6.1	production	\N	2534357482
initial-schema-32	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.682588	34	EXECUTED	8:6cee412b6b2899d22116dd0529310990	createIndex indexName=SYS_IDX_SYS_CT_10149_10151, tableName=cwd_group		\N	3.6.1	production	\N	2534357482
initial-schema-33	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.688593	35	EXECUTED	8:1bbd513d56dea7e837103c6f30fe1842	createIndex indexName=IDX_GROUP_ATTR_DIR_NAME_LVAL, tableName=cwd_group_attribute		\N	3.6.1	production	\N	2534357482
initial-schema-34	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.695917	36	EXECUTED	8:b9bf51ae908d06fcd55926610e3dfaca	createIndex indexName=SYS_IDX_SYS_CT_10157_10159, tableName=cwd_group_attribute		\N	3.6.1	production	\N	2534357482
initial-schema-35	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.702516	37	EXECUTED	8:45b717a4ac7fb130878e802b8687af7e	createIndex indexName=IDX_MEM_DIR_CHILD, tableName=cwd_membership		\N	3.6.1	production	\N	2534357482
initial-schema-36	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.711107	38	EXECUTED	8:0e003d3fa20f553d87589947bbfffe90	createIndex indexName=IDX_MEM_DIR_PARENT, tableName=cwd_membership		\N	3.6.1	production	\N	2534357482
initial-schema-37	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.719235	39	EXECUTED	8:fceb84979c267995fba86d581754a657	createIndex indexName=IDX_MEM_DIR_PARENT_CHILD, tableName=cwd_membership		\N	3.6.1	production	\N	2534357482
initial-schema-38	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.724726	40	EXECUTED	8:a7ed9b793f363e476a56741f726aa815	createIndex indexName=SYS_IDX_SYS_CT_10168_10170, tableName=cwd_membership		\N	3.6.1	production	\N	2534357482
initial-schema-39	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.730108	41	EXECUTED	8:4fcd1cce1ef1c981fefdf4472b01ef52	createIndex indexName=IDX_TOKEN_DIR_ID, tableName=cwd_token		\N	3.6.1	production	\N	2534357482
initial-schema-40	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.740755	42	EXECUTED	8:aafa03422d46f580a5ca8020067f8143	createIndex indexName=IDX_TOKEN_KEY, tableName=cwd_token		\N	3.6.1	production	\N	2534357482
initial-schema-41	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.749645	43	EXECUTED	8:e7e234f4c47710f20fa1e5177cef7cf4	createIndex indexName=IDX_TOKEN_LAST_ACCESS, tableName=cwd_token		\N	3.6.1	production	\N	2534357482
initial-schema-42	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.758942	44	EXECUTED	8:fc4fbff981cfa393f7654ec238c59f75	createIndex indexName=IDX_TOKEN_NAME_DIR_ID, tableName=cwd_token		\N	3.6.1	production	\N	2534357482
initial-schema-43	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.766135	45	EXECUTED	8:4b6d457569807cae767016b918863e36	createIndex indexName=SYS_IDX_SYS_CT_10184_10186, tableName=cwd_token		\N	3.6.1	production	\N	2534357482
initial-schema-44	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.779039	46	EXECUTED	8:a6361f5f0bd3ced4b14905cc2e8f1b13	createIndex indexName=IDX_USER_ACTIVE, tableName=cwd_user		\N	3.6.1	production	\N	2534357482
initial-schema-45	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.787227	47	EXECUTED	8:ecc18bdbaf1d5c55e759c37527d3aac1	createIndex indexName=IDX_USER_LOWER_DISPLAY_NAME, tableName=cwd_user		\N	3.6.1	production	\N	2534357482
initial-schema-46	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.793393	48	EXECUTED	8:31704caff15e8363097fd89813e6e4ff	createIndex indexName=IDX_USER_LOWER_EMAIL_ADDRESS, tableName=cwd_user		\N	3.6.1	production	\N	2534357482
initial-schema-47	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.804835	49	EXECUTED	8:f53f69e722f8f0c5880afd9af66f51ee	createIndex indexName=IDX_USER_LOWER_FIRST_NAME, tableName=cwd_user		\N	3.6.1	production	\N	2534357482
initial-schema-48	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.81096	50	EXECUTED	8:321ca67f248aaa28bbe9d2a4c0b06b7f	createIndex indexName=IDX_USER_LOWER_LAST_NAME, tableName=cwd_user		\N	3.6.1	production	\N	2534357482
initial-schema-49	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.817317	51	EXECUTED	8:5425b405dc859e8db6b36dfbf9b28066	createIndex indexName=SYS_IDX_SYS_CT_10195_10197, tableName=cwd_user		\N	3.6.1	production	\N	2534357482
initial-schema-50	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.823334	52	EXECUTED	8:5af2e751a27caeb25816ea10a0777b29	createIndex indexName=IDX_USER_ATTR_DIR_NAME_LVAL, tableName=cwd_user_attribute		\N	3.6.1	production	\N	2534357482
initial-schema-51	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.835188	53	EXECUTED	8:a38436887e2bb770f3fcd8acf3bc0c02	createIndex indexName=SYS_IDX_SYS_CT_10203_10205, tableName=cwd_user_attribute		\N	3.6.1	production	\N	2534357482
initial-schema-52	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.857129	54	EXECUTED	8:ac90479cbe4b203756de093fa232e899	createTable tableName=project		\N	3.6.1	production	\N	2534357482
initial-schema-53	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.872743	55	EXECUTED	8:552c8c2c30eb591f6222bab4125d8190	createTable tableName=stash_user		\N	3.6.1	production	\N	2534357482
initial-schema-54	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.892074	56	EXECUTED	8:8835eaf08a3999d10341977128a0f01c	createTable tableName=repository; addUniqueConstraint constraintName=uk_slug_project_id, tableName=repository; addForeignKeyConstraint baseTableName=repository, constraintName=fk_repository_project, referencedTableName=project		\N	3.6.1	production	\N	2534357482
initial-schema-55	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.912409	57	EXECUTED	8:c22e6b9f8323bbac4238dfcf5ba9a9d7	createTable tableName=granted_permission; addForeignKeyConstraint baseTableName=granted_permission, constraintName=fk_perm_project, referencedTableName=project; addForeignKeyConstraint baseTableName=granted_permission, constraintName=fk_perm_user,...		\N	3.6.1	production	\N	2534357482
initial-schema-56	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.925297	58	EXECUTED	8:e41fc9c33c174aaf5c692e36bb76f17b	createTable tableName=trusted_app		\N	3.6.1	production	\N	2534357482
initial-schema-57	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.943289	59	EXECUTED	8:072f6402e0a2f3ee86151e9d74eb3106	createTable tableName=trusted_app_restriction; addUniqueConstraint constraintName=uk_trusted_app_restrict, tableName=trusted_app_restriction; addForeignKeyConstraint baseTableName=trusted_app_restriction, constraintName=fk_trusted_app, referencedT...		\N	3.6.1	production	\N	2534357482
initial-schema-58	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.954167	60	EXECUTED	8:666974d732dbc2409d2de93f16e0884e	createTable tableName=current_app		\N	3.6.1	production	\N	2534357482
initial-schema-59	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.964883	61	EXECUTED	8:876b33ac4ceb0163a95c966ab4c0430b	createTable tableName=persistent_logins		\N	3.6.1	production	\N	2534357482
initial-schema-60	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.967478	62	MARK_RAN	8:e4538a6e0143083914e4284c1524c665	createTable tableName=app_property		\N	3.6.1	production	\N	2534357482
initial-schema-61	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.972842	63	EXECUTED	8:4a8c7df70c8e9c1104cca7455e73b9e8	createTable tableName=hibernate_unique_key; insert tableName=hibernate_unique_key		\N	3.6.1	production	\N	2534357482
initial-schema-62	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.979296	64	EXECUTED	8:99cc426c52a1590e083b46427c2e1b82	createTable tableName=id_sequence		\N	3.6.1	production	\N	2534357482
initial-schema-63	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.991622	65	EXECUTED	8:064ca1d812ad06bfb88cc5bde67fd71a	createTable tableName=plugin_setting; addPrimaryKey tableName=plugin_setting		\N	3.6.1	production	\N	2534357482
initial-schema-64	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/initial-schema.xml	2021-06-01 07:59:18.999794	66	EXECUTED	8:62e54854cca11295bb6a1058ede27a58	createTable tableName=plugin_state		\N	3.6.1	production	\N	2534357482
m13-01	gcrain	com/atlassian/caviar/db/changelog/r1_0/m13.xml	2021-06-01 07:59:19.008754	67	MARK_RAN	8:35b27fbca70d9672f38f711f5eaf3930	modifyDataType columnName=prop_value, tableName=app_property		\N	3.6.1	production	\N	2534357482
CAV-1123-1	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/m15.xml	2021-06-01 07:59:19.043353	68	EXECUTED	8:f9ae3553c081e51a25f5f7ea3de45b93	createTable tableName=changeset		\N	3.6.1	production	\N	2534357482
CAV-1123-2	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/m15.xml	2021-06-01 07:59:19.055952	69	EXECUTED	8:62167cad023311478d85234ddcf31019	createTable tableName=cs_attribute		\N	3.6.1	production	\N	2534357482
CAV-1123-3	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/m15.xml	2021-06-01 07:59:19.063704	70	EXECUTED	8:14b9a9c3e02d767a47180c347695e2e5	createIndex indexName=idx_cs_to_attribute, tableName=cs_attribute		\N	3.6.1	production	\N	2534357482
CAV-1123-4	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/m15.xml	2021-06-01 07:59:19.071587	71	EXECUTED	8:2be5a010ca98df59592bfa18ed710274	createIndex indexName=idx_attribute_to_cs, tableName=cs_attribute		\N	3.6.1	production	\N	2534357482
CAV-1123-5	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/m15.xml	2021-06-01 07:59:19.075819	72	EXECUTED	8:1bf1fb67a36b194571a2cb7a6fa4daf8	createTable tableName=cs_repo_membership		\N	3.6.1	production	\N	2534357482
CAV-1123-6	mheemskerk	com/atlassian/caviar/db/changelog/r1_0/m15.xml	2021-06-01 07:59:19.080643	73	EXECUTED	8:ade2ed0c01a0a600c7a23c8aa5301069	createTable tableName=cs_indexer_state		\N	3.6.1	production	\N	2534357482
STASH-1842-1	gcrain	liquibase/r1_0/m17.xml	2021-06-01 07:59:19.111405	74	EXECUTED	8:fd31e95f3efaed4b0063fc37a476fb88	dropColumn columnName=description, tableName=repository		\N	3.6.1	production	\N	2534357482
STASHDEV-363-recent-repos-01	mstudman	liquibase/r1_1/m03.xml	2021-06-01 07:59:19.118849	75	EXECUTED	8:5dd725e0999a54b238cabc06857eed71	createTable tableName=repository_access		\N	3.6.1	production	\N	2534357482
STASHDEV-363-recent-repos-02	mstudman	liquibase/r1_1/m03.xml	2021-06-01 07:59:19.125276	76	EXECUTED	8:2c80f80ca9f2cb5015a6c80c383208bf	addPrimaryKey constraintName=PK_REPOSITORY_ACCESS, tableName=repository_access		\N	3.6.1	production	\N	2534357482
STASHDEV-363-recent-repos-03	mstudman	liquibase/r1_1/m03.xml	2021-06-01 07:59:19.13677	77	EXECUTED	8:c795b39064a23a222c740e0c9c0ef5f1	createIndex indexName=IDX_REPOSITORY_ACCESS_USER_ID, tableName=repository_access		\N	3.6.1	production	\N	2534357482
STASHDEV-595-1	bturner	liquibase/r1_1/m04.xml	2021-06-01 07:59:19.141679	78	EXECUTED	8:d4bb1e27bbb70577f6c46bca78676990	dropColumn columnName=allow_anon, tableName=granted_permission	Anonymous access is not currently supported in Stash. When it is, it will likely not be implemented the way\n            this column was designed to support, so there's no point in keeping the column in the table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-595-2	bturner	liquibase/r1_1/m04.xml	2021-06-01 07:59:19.159891	79	EXECUTED	8:71b2cbbc719c13fef830690530704627	createTable tableName=weighted_permission	Metadata table for adding knowledge of each's permission's relative weight to the database. This allows more\n            efficient retrieval of a user or group's "highest" permission.\n\n            See Permission.getWeight() documentation for more ...	\N	3.6.1	\N	\N	2534357482
STASHDEV-595-3	bturner	liquibase/r1_1/m04.xml	2021-06-01 07:59:19.174802	80	EXECUTED	8:7b73c4b33cb7821bf7ab22407cc86430	insert tableName=weighted_permission; insert tableName=weighted_permission; insert tableName=weighted_permission; insert tableName=weighted_permission; insert tableName=weighted_permission; insert tableName=weighted_permission; insert tableName=we...	Initial population of relative weights for all permissions. These values must match the values specified in\n            the Permission enumeration _exactly_ or the database will return incorrect results.	\N	3.6.1	\N	\N	2534357482
STASHDEV-595-4	bturner	liquibase/r1_1/m04.xml	2021-06-01 07:59:19.190948	81	EXECUTED	8:71f3aae9d32e5a3fe810475ef70a61dc	addForeignKeyConstraint baseTableName=granted_permission, constraintName=granted_perm_weight_fk, referencedTableName=weighted_permission		\N	3.6.1	\N	\N	2534357482
STASHDEV-1042-01	jhinch	liquibase/r1_2/m01.xml	2021-06-01 07:59:19.197497	82	EXECUTED	8:0f58f50049088e74e0ed3c17e095c992	addColumn tableName=trusted_app		\N	3.6.1	\N	\N	2534357482
STASHDEV-1042-02	jhinch	liquibase/r1_2/m01.xml	2021-06-01 07:59:19.20532	83	EXECUTED	8:8c8fec006125f5780a76beefbe11c356	customChange		\N	3.6.1	\N	\N	2534357482
STASHDEV-1042-03	jhinch	liquibase/r1_2/m01.xml	2021-06-01 07:59:19.210268	84	EXECUTED	8:a0d6e1601c79285191cda47e5725db58	dropColumn columnName=public_key, tableName=trusted_app		\N	3.6.1	\N	\N	2534357482
STASHDEV-1042-04	jhinch	liquibase/r1_2/m01.xml	2021-06-01 07:59:19.215913	85	EXECUTED	8:9ee929cf85b74641017adcdbcfe4fb4a	addNotNullConstraint columnName=public_key_base64, tableName=trusted_app		\N	3.6.1	\N	\N	2534357482
STASHDEV-1042-05	jhinch	liquibase/r1_2/m01.xml	2021-06-01 07:59:19.222797	86	EXECUTED	8:b3add94fa34c7ac4cc67ece2f0f9aa0b	addColumn tableName=current_app		\N	3.6.1	\N	\N	2534357482
STASHDEV-1042-06	jhinch	liquibase/r1_2/m01.xml	2021-06-01 07:59:19.226455	87	EXECUTED	8:836f099ebd05739986fba69ea122bf51	customChange		\N	3.6.1	\N	\N	2534357482
STASHDEV-1042-07	jhinch	liquibase/r1_2/m01.xml	2021-06-01 07:59:19.229865	88	EXECUTED	8:c306d0c6b97969d7dbc098f69b87b8a0	dropColumn columnName=public_key, tableName=current_app		\N	3.6.1	\N	\N	2534357482
STASHDEV-1042-08	jhinch	liquibase/r1_2/m01.xml	2021-06-01 07:59:19.234107	89	EXECUTED	8:940ed53b750e59819dba8d8086d500c0	addNotNullConstraint columnName=public_key_base64, tableName=current_app		\N	3.6.1	\N	\N	2534357482
STASHDEV-1042-09	jhinch	liquibase/r1_2/m01.xml	2021-06-01 07:59:19.238352	90	EXECUTED	8:29b8fe4e3a43a3adaab7a5e70c22eebe	addColumn tableName=current_app		\N	3.6.1	\N	\N	2534357482
STASHDEV-1042-10	jhinch	liquibase/r1_2/m01.xml	2021-06-01 07:59:19.241299	91	EXECUTED	8:5ae968d051a27e477f0bfdcfa64db71f	customChange		\N	3.6.1	\N	\N	2534357482
STASHDEV-1042-11	jhinch	liquibase/r1_2/m01.xml	2021-06-01 07:59:19.253668	92	EXECUTED	8:346af71723fe964430944fb5df69fb11	dropColumn columnName=private_key, tableName=current_app		\N	3.6.1	\N	\N	2534357482
STASHDEV-1042-12	jhinch	liquibase/r1_2/m01.xml	2021-06-01 07:59:19.258934	93	EXECUTED	8:a307af97bcf8d4fe70574510141b1cb0	addNotNullConstraint columnName=private_key_base64, tableName=current_app		\N	3.6.1	\N	\N	2534357482
STASHDEV-616-01	dpinn	liquibase/r1_2/m03.xml	2021-06-01 07:59:20.187394	94	MARK_RAN	8:f9c86265f5225322f883b236df1f235b	dropForeignKeyConstraint baseTableName=repository, constraintName=fk_repository_repository_origin	Drop the REPOSITORY.FK_REPOSITORY_REPOSITORY_ORIGIN foreign key constraint	\N	3.6.1	\N	\N	2534357482
STASHDEV-616-02	dpinn	liquibase/r1_2/m03.xml	2021-06-01 07:59:20.197437	95	EXECUTED	8:74113ee2b8736b02c1f4d2cf2f51a417	addForeignKeyConstraint baseTableName=repository, constraintName=fk_repository_origin, referencedTableName=repository	Add the REPOSITORY.FK_REPOSITORY_ORIGIN foreign key constraint\n            to replace the REPOSITORY.FK_REPOSITORY_REPOSITORY_ORIGIN foreign key constraint	\N	3.6.1	\N	\N	2534357482
STASHDEV-616-03	dpinn	liquibase/r1_2/m03.xml	2021-06-01 07:59:21.057304	96	MARK_RAN	8:edd946f16a5dc255d151d843ceb4c698	dropForeignKeyConstraint baseTableName=cs_repo_membership, constraintName=FK_CS_REPO_MEMBERSHIP_CHANGESET	Drop the CS_REPO_MEMBERSHIP.FK_CS_REPO_MEMBERSHIP_CHANGESET foreign key constraint	\N	3.6.1	\N	\N	2534357482
STASHDEV-616-04	dpinn	liquibase/r1_2/m03.xml	2021-06-01 07:59:21.071053	97	EXECUTED	8:e476c45a464d0588d3a8f7531541cfbc	addForeignKeyConstraint baseTableName=cs_repo_membership, constraintName=fk_repo_membership_changeset, referencedTableName=changeset	Add the CS_REPO_MEMBERSHIP.FK_REPO_MEMBERSHIP_CHANGESET foreign key constraint\n            to replace the CS_REPO_MEMBERSHIP.FK_CS_REPO_MEMBERSHIP_CHANGESET foreign key constraint	\N	3.6.1	\N	\N	2534357482
STASHDEV-616-05	dpinn	liquibase/r1_2/m03.xml	2021-06-01 07:59:21.963459	98	MARK_RAN	8:f89386f9c479b72367fcb69e0ae6c27d	dropForeignKeyConstraint baseTableName=cs_repo_membership, constraintName=FK_CS_REPO_MEMBERSHIP_REPOSITORY	Drop the CS_REPO_MEMBERSHIP.FK_CS_REPO_MEMBERSHIP_REPOSITORY foreign key constraint	\N	3.6.1	\N	\N	2534357482
STASHDEV-616-06	dpinn	liquibase/r1_2/m03.xml	2021-06-01 07:59:21.973966	99	EXECUTED	8:42fbfdd07ccd6ca122a082c4ae4cae32	addForeignKeyConstraint baseTableName=cs_repo_membership, constraintName=fk_repo_membership_repo, referencedTableName=repository	Add the CS_REPO_MEMBERSHIP.FK_REPO_MEMBERSHIP_REPO foreign key constraint\n            to replace the CS_REPO_MEMBERSHIP.FK_CS_REPO_MEMBERSHIP_REPOSITORY foreign key constraint	\N	3.6.1	\N	\N	2534357482
STASHDEV-616-07	dpinn	liquibase/r1_2/m03.xml	2021-06-01 07:59:22.886128	100	MARK_RAN	8:d75b8e708cce54546fda43eec9a1a8ec	dropForeignKeyConstraint baseTableName=repository_access, constraintName=FK_REPOSITORY_ACCESS_ID_STASH_USER_ID	Drop the REPOSITORY_ACCESS.FK_REPOSITORY_ACCESS_ID_STASH_USER_ID foreign key constraint	\N	3.6.1	\N	\N	2534357482
STASHDEV-616-08	dpinn	liquibase/r1_2/m03.xml	2021-06-01 07:59:22.893598	101	EXECUTED	8:66ca0c09e06ea31cf84db1297fc40df4	addForeignKeyConstraint baseTableName=repository_access, constraintName=fk_repository_access_user, referencedTableName=stash_user	Add the REPOSITORY_ACCESS.FK_REPOSITORY_ACCESS_USER foreign key constraint\n            to replace the REPOSITORY_ACCESS.FK_REPOSITORY_ACCESS_ID_STASH_USER_ID foreign key constraint	\N	3.6.1	\N	\N	2534357482
STASHDEV-616-09	dpinn	liquibase/r1_2/m03.xml	2021-06-01 07:59:23.722749	102	MARK_RAN	8:0deeb6dd3e4bb6a23fbc8ac87f1661ef	dropForeignKeyConstraint baseTableName=repository_access, constraintName=FK_REPOSITORY_ACCESS_ID_REPOSITORY_ID	Drop the REPOSITORY_ACCESS.FK_REPOSITORY_ACCESS_ID_REPOSITORY_ID foreign key constraint	\N	3.6.1	\N	\N	2534357482
STASHDEV-616-10	dpinn	liquibase/r1_2/m03.xml	2021-06-01 07:59:23.729285	103	EXECUTED	8:f1ba726b910bc74b9a6c4f5cfd5cd0e4	addForeignKeyConstraint baseTableName=repository_access, constraintName=fk_repository_access_repo, referencedTableName=repository	Add the REPOSITORY_ACCESS.FK_REPOSITORY_ACCESS_REPO foreign key constraint\n            to replace the REPOSITORY_ACCESS.FK_REPOSITORY_ACCESS_ID_REPOSITORY_ID foreign key constraint	\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-1	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.746067	104	EXECUTED	8:a0554daf3320533492d4f0386374ef32	createTable tableName=sta_comment		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-2	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.75194	105	EXECUTED	8:24f724937ff04a35746e716722867fd4	createIndex indexName=idx_sta_comment_author, tableName=sta_comment		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-3	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.757661	106	EXECUTED	8:52516788e9057cdf88cb190944349cb1	createIndex indexName=idx_sta_comment_parent, tableName=sta_comment		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-4	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.763261	107	EXECUTED	8:5282bb545c872889344993e7cf03a034	createIndex indexName=idx_sta_comment_root, tableName=sta_comment		\N	3.6.1	\N	\N	2534357482
STASHDEV-1651-5	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.765796	108	MARK_RAN	8:ba63784e24375e426834adb57e12142d	dropForeignKeyConstraint baseTableName=sta_comment, constraintName=fk_sta_comment_author	If fk_sta_comment_author was created by a previous incarnation of this changelog, drop it so it can be\n            recreated correctly.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1651-6	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.770111	109	EXECUTED	8:f7f24e40226871efd68636730d72c113	addForeignKeyConstraint baseTableName=sta_comment, constraintName=fk_sta_comment_author, referencedTableName=stash_user	Create a foreign key between comments and their authors. Note that this foreign key does not cascade\n            deletions as it is expected that Stash users will never be deleted.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1651-1	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.772557	110	MARK_RAN	8:e50ef3151843e3d229a6e2b04e3e0227	dropForeignKeyConstraint baseTableName=sta_comment, constraintName=fk_sta_comment_parent	If fk_sta_comment_parent was created by a previous incarnation of this changelog, drop it so it can be\n            recreated correctly.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1651-2	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.778314	111	EXECUTED	8:be41db6a8b52e457126c08911b98f468	addForeignKeyConstraint baseTableName=sta_comment, constraintName=fk_sta_comment_parent, referencedTableName=sta_comment	Create a foreign key between replies and their parent comment. Note that this foreign key does not cascade\n            deletions; such cascades must be handled in code.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1651-3	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.782788	112	MARK_RAN	8:7ce5857a79d122c213422a6f15610ab2	dropForeignKeyConstraint baseTableName=sta_comment, constraintName=fk_sta_comment_root	If fk_sta_comment_root was created by a previous incarnation of this changelog, drop it so it can be\n            recreated correctly.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1651-4	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.788127	113	EXECUTED	8:dee41ebe512edcce3cb30a23b1e9efc5	addForeignKeyConstraint baseTableName=sta_comment, constraintName=fk_sta_comment_root, referencedTableName=sta_comment	Create a foreign key between replies and their root comment, where a root comment is the top-level comment\n            in a thread. Note that this foreign key does not cascade deletions; such cascades must be handled in code.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-8	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.80466	114	EXECUTED	8:0e2563311d9062f809edc02cfdce56e4	createTable tableName=sta_diff_comment_anchor		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-9	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.810755	115	EXECUTED	8:e7cf068d6a09fb754768f5e801d69aec	createIndex indexName=idx_sta_diff_comment_comment, tableName=sta_diff_comment_anchor		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-10	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.818155	116	EXECUTED	8:ed031928fbf2fda31ab03b8ef8d96fec	createIndex indexName=idx_sta_diff_comment_anchors, tableName=sta_diff_comment_anchor		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-11	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.822253	117	EXECUTED	8:0ba3646712e93d43ea84e375cb338fec	addForeignKeyConstraint baseTableName=sta_diff_comment_anchor, constraintName=fk_sta_diff_comment_comment, referencedTableName=sta_comment		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-12	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.830812	118	EXECUTED	8:386ee7b025cb26e6142ab24a4a1959e6	createTable tableName=sta_activity		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-13	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.836792	119	EXECUTED	8:d6ddbfe931e98ef5f061dbfaeb597c14	createIndex indexName=idx_sta_activity_type, tableName=sta_activity		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-14	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.844724	120	EXECUTED	8:00983b5445599296ae25b6a3519e5a32	createIndex indexName=idx_sta_activity_user, tableName=sta_activity		\N	3.6.1	\N	\N	2534357482
STASHDEV-1651-7	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.84768	121	MARK_RAN	8:4177faf7e9034419e0c4bf429cdf5fd9	dropForeignKeyConstraint baseTableName=sta_activity, constraintName=fk_sta_activity_user	If fk_sta_activity_user was created by a previous incarnation of this changelog, drop it so it can be\n            recreated correctly.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1651-8	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.854321	122	EXECUTED	8:a1a21601d2791d9bde79005ce653f5ad	addForeignKeyConstraint baseTableName=sta_activity, constraintName=fk_sta_activity_user, referencedTableName=stash_user	Create a foreign key between activities and their user. Note that this foreign key does not cascade\n            deletions as it is expected that Stash users will never be deleted.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1205-1	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.866658	123	EXECUTED	8:6f1ce1af5216a543fa4941e5b0e27390	createTable tableName=sta_repository_scoped_id	Sequence generation table for creating IDs that are scoped by a repository. Multiple scope types may exist\n            for each repository, and each scope will get its own sequence of IDs, starting from 1.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1205-2	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.876467	124	EXECUTED	8:53071ad111ea0257380d49c5e2f38fab	addForeignKeyConstraint baseTableName=sta_repository_scoped_id, constraintName=fk_sta_repo_scoped_id_repo, referencedTableName=repository		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-16	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.914732	125	EXECUTED	8:583f62281fe556f4d053e681aef7ee0a	createTable tableName=sta_pull_request		\N	3.6.1	\N	\N	2534357482
STASHDEV-1205-3	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.920725	126	EXECUTED	8:7299af62cff73c46982a93c591353182	addUniqueConstraint constraintName=uq_sta_pull_request_scoped_id, tableName=sta_pull_request		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-17	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.927625	127	EXECUTED	8:7297fc6378d223b74f9c464d48ba8d25	createIndex indexName=idx_sta_pull_request_from_repo, tableName=sta_pull_request		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-18	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.932936	128	EXECUTED	8:182353aa03607f34272e0c8ddf619e54	createIndex indexName=idx_sta_pull_request_state, tableName=sta_pull_request		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-19	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.939521	129	EXECUTED	8:60391191264af3190354b4dfe6ad53fa	createIndex indexName=idx_sta_pull_request_to_repo, tableName=sta_pull_request		\N	3.6.1	\N	\N	2534357482
STASHDEV-1205-4	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.94495	130	EXECUTED	8:2786f466b34621fcf71c64da42c0823c	addForeignKeyConstraint baseTableName=sta_pull_request, constraintName=fk_sta_pull_request_from_repo, referencedTableName=repository		\N	3.6.1	\N	\N	2534357482
STASHDEV-1205-5	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.949831	131	EXECUTED	8:f241784785ca900e020113aaf78e118d	addForeignKeyConstraint baseTableName=sta_pull_request, constraintName=fk_sta_pull_request_to_repo, referencedTableName=repository		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-20	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.961668	132	EXECUTED	8:6c49f239f4bd88dac09faf25508111fa	createTable tableName=sta_pr_activity		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-21	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.982215	133	EXECUTED	8:06d0687511ef92b5e8691e5df2d3bf0e	createIndex indexName=idx_sta_pr_activity, tableName=sta_pr_activity		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-22	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.991555	134	EXECUTED	8:29832f99b8c1dce6f173a6c6d6379d7f	addForeignKeyConstraint baseTableName=sta_pr_activity, constraintName=fk_sta_pr_activity_id, referencedTableName=sta_activity		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-23	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:23.999265	135	EXECUTED	8:ada7674aaa8c7fc4bd7a27c8103e4a79	addForeignKeyConstraint baseTableName=sta_pr_activity, constraintName=fk_sta_pr_activity_pr, referencedTableName=sta_pull_request		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-24	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.010926	136	EXECUTED	8:da320c268d183a671252f1ef2447cfb6	createTable tableName=sta_pr_comment_activity		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-25	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.017817	137	EXECUTED	8:0f07dbc82b035819607a689031552eec	createIndex indexName=idx_st_pr_com_act_anchor, tableName=sta_pr_comment_activity		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-26	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.025768	138	EXECUTED	8:87c3b275df03783b803565aa95655eb1	createIndex indexName=idx_st_pr_com_act_comment, tableName=sta_pr_comment_activity		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-27	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.030586	139	EXECUTED	8:a923bbccf329f910a94cb2ddb352ecc5	addForeignKeyConstraint baseTableName=sta_pr_comment_activity, constraintName=fk_sta_pr_com_act_id, referencedTableName=sta_pr_activity		\N	3.6.1	\N	\N	2534357482
STASHDEV-4154-1	jhinch	liquibase/r2_5/upgrade.xml	2021-06-01 07:59:25.001181	282	EXECUTED	8:6cda8bdca6b69301c5979e012c8adf0d	addColumn tableName=repository	Add a column to the "repository" table for the public flag.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1651-9	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.03487	140	MARK_RAN	8:895a6d7fc326ae2cd54d84df06db4f28	dropForeignKeyConstraint baseTableName=sta_pr_comment_activity, constraintName=fk_sta_pr_com_act_anchor	If fk_sta_pr_com_act_anchor was created by a previous incarnation of this changelog, drop it so it can be\n            recreated correctly.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1651-10	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.046401	141	EXECUTED	8:60f9269bb94f057e2093a29bdf35960a	addForeignKeyConstraint baseTableName=sta_pr_comment_activity, constraintName=fk_sta_pr_com_act_anchor, referencedTableName=sta_diff_comment_anchor	Create a foreign key between comment activities and their comment anchor, if one is set. Note that this\n            foreign key does not cascade deletions as it is expected anchors will never be deleted without deleting\n            their attached ...	\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-29	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.050529	142	EXECUTED	8:1b5c30e7e9d3ed684af6934c39c188e4	addForeignKeyConstraint baseTableName=sta_pr_comment_activity, constraintName=fk_sta_pr_com_act_comment, referencedTableName=sta_comment		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-30	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.059856	143	EXECUTED	8:aaf2aa32e0dcae403e6b52be0c524850	createTable tableName=sta_pr_participant		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-31	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.065217	144	EXECUTED	8:9c8c6f26b7520b4112cfbbb2e7eaac00	createIndex indexName=idx_sta_pr_participant_pr, tableName=sta_pr_participant		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-32	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.071135	145	EXECUTED	8:b384256646d2aa4f36f6a6c107cee049	createIndex indexName=idx_sta_pr_participant_user, tableName=sta_pr_participant		\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-33	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.076376	146	EXECUTED	8:222b657ed3fa43e7f714ef18ef650f67	addForeignKeyConstraint baseTableName=sta_pr_participant, constraintName=fk_sta_pr_participant_pr, referencedTableName=sta_pull_request		\N	3.6.1	\N	\N	2534357482
STASHDEV-1651-11	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.078987	147	MARK_RAN	8:4e5dfe4b536a50e3003f24f00037eb23	dropForeignKeyConstraint baseTableName=sta_pr_participant, constraintName=fk_sta_pr_participant_user	If fk_sta_pr_participant_user was created by a previous incarnation of this changelog, drop it so it can be\n            recreated correctly.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1651-12	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.084077	148	EXECUTED	8:2a7c84a363304a1f20e4b259fdbb01a6	addForeignKeyConstraint baseTableName=sta_pr_participant, constraintName=fk_sta_pr_participant_user, referencedTableName=stash_user	Create a foreign key between participants and their user. Note that this foreign key does not cascade\n            deletions as it is expected that Stash users will never be deleted.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1023-35	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.08957	149	EXECUTED	8:45d636425e280e5ec8910702e7ecb72c	addUniqueConstraint constraintName=uq_sta_pr_participant_pr_user, tableName=sta_pr_participant		\N	3.6.1	\N	\N	2534357482
STASHDEV-1392-1	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.094779	150	EXECUTED	8:8724a4425e4663f14a2b0eaa8760f214	dropForeignKeyConstraint baseTableName=sta_repository_scoped_id, constraintName=fk_sta_repo_scoped_id_repo		\N	3.6.1	\N	\N	2534357482
STASHDEV-1392-2	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.099549	151	EXECUTED	8:e843fc6354b81d67ef8d44dff4136156	addForeignKeyConstraint baseTableName=sta_repository_scoped_id, constraintName=fk_sta_repo_scoped_id_repo, referencedTableName=repository		\N	3.6.1	\N	\N	2534357482
STASHDEV-1392-3	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.10417	152	EXECUTED	8:6aff120a6fe13a13e4d9fc4f06e49836	dropForeignKeyConstraint baseTableName=sta_pull_request, constraintName=fk_sta_pull_request_from_repo		\N	3.6.1	\N	\N	2534357482
STASHDEV-1651-13	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.107809	153	MARK_RAN	8:6aff120a6fe13a13e4d9fc4f06e49836	dropForeignKeyConstraint baseTableName=sta_pull_request, constraintName=fk_sta_pull_request_from_repo	If fk_sta_pull_request_from_repo was created by a previous incarnation of this changelog, drop it so it can\n            be recreated correctly.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1651-14	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.114424	154	EXECUTED	8:2786f466b34621fcf71c64da42c0823c	addForeignKeyConstraint baseTableName=sta_pull_request, constraintName=fk_sta_pull_request_from_repo, referencedTableName=repository	Create a foreign key between pull requests and the repository they originate from. Note that this foreign\n            key does not cascade deletions as, currently, the from and to repositories will always be the same. When\n            this changes...	\N	3.6.1	\N	\N	2534357482
STASHDEV-1392-5	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.121046	155	EXECUTED	8:2882d4863ab42ccd45b59f8f4b55bb28	dropForeignKeyConstraint baseTableName=sta_pull_request, constraintName=fk_sta_pull_request_to_repo		\N	3.6.1	\N	\N	2534357482
STASHDEV-1392-6	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.127713	156	EXECUTED	8:60ee9ece699daed2e040326d05921e84	addForeignKeyConstraint baseTableName=sta_pull_request, constraintName=fk_sta_pull_request_to_repo, referencedTableName=repository		\N	3.6.1	\N	\N	2534357482
STASHDEV-1392-7	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.131661	157	EXECUTED	8:008a21d59ccfed017f262ce8639585d5	renameColumn newColumnName=to_path, oldColumnName=file_path, tableName=sta_diff_comment_anchor		\N	3.6.1	\N	\N	2534357482
STASHDEV-1392-8	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.140405	158	EXECUTED	8:087a4312d1cb8dae5a809ac431ea5502	addColumn tableName=sta_diff_comment_anchor		\N	3.6.1	\N	\N	2534357482
STASHDEV-1455-1	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.145019	159	EXECUTED	8:7fa1bdeb7a8844747a738ffc2608e394	addColumn tableName=repository		\N	3.6.1	\N	\N	2534357482
STASHDEV-1455-2	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.148762	160	EXECUTED	8:2fa425455caa0f67ed15f7a3737e7b39	update tableName=repository		\N	3.6.1	\N	\N	2534357482
STASHDEV-1455-3	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.153247	161	EXECUTED	8:9d71d980ff7233c141033978fddeb6e5	addNotNullConstraint columnName=scm_id, tableName=repository		\N	3.6.1	\N	\N	2534357482
STASHDEV-1455-4	bturner	liquibase/r1_3/m01.xml	2021-06-01 07:59:24.156771	162	EXECUTED	8:ba4cebc7755b9e44470b224eee72e658	dropColumn columnName=scmtype, tableName=repository		\N	3.6.1	\N	\N	2534357482
STASHDEV-1201-1	bturner	liquibase/r1_3/m02.xml	2021-06-01 07:59:24.161792	163	EXECUTED	8:20230830144ff80748f87be3ac5fbea5	addColumn tableName=sta_pull_request	Add a new CLOB column to sta_pull_request to replace the previous description column.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1201-2	bturner	liquibase/r1_3/m02.xml	2021-06-01 07:59:24.165057	164	EXECUTED	8:d2039913900f0f7c1c80751e81b5c060	update tableName=sta_pull_request	Copy the VARCHAR description into the CLOB pr_description. Note that, in Postgres, this will store the text\n            directly in the column, which is the correct usage to support Unicode encodings.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1201-3	bturner	liquibase/r1_3/m02.xml	2021-06-01 07:59:24.169302	165	EXECUTED	8:f85b20c4af2fcb19ea23976a9feaba12	dropColumn columnName=description, tableName=sta_pull_request	Drop the old VARCHAR(255) description column.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1201-4	bturner	liquibase/r1_3/m02.xml	2021-06-01 07:59:24.173605	166	EXECUTED	8:134f9e386cdf00393315e8169eed3a38	renameColumn newColumnName=description, oldColumnName=pr_description, tableName=sta_pull_request	Rename the pr_description column to description. Note that, to correctly support MySQL, the data type\n            set here is MySQL-specific. That property is ignored for all other RDBMSs.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1724-2	bturner	liquibase/r1_3/m02.xml	2021-06-01 07:59:24.181609	167	EXECUTED	8:b5096aa6f7d728eb862c76087c2f4700	customChange	Rewrite comment CLOBs which have been stored as LargeObjects (OIDs) in Postgres to store the text directly\n            in the column instead.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1690-1	dpinn	liquibase/r1_3/m02.xml	2021-06-01 07:59:24.189624	168	EXECUTED	8:45729d6f643447ce83a970e97264e534	createTable tableName=sta_watcher		\N	3.6.1	\N	\N	2534357482
STASHDEV-1690-2	tpettersen	liquibase/r1_3/m02.xml	2021-06-01 07:59:24.195794	169	EXECUTED	8:f5de46a565cdca7f1efdcdab38cd93a7	addUniqueConstraint constraintName=uq_sta_watcher, tableName=sta_watcher		\N	3.6.1	\N	\N	2534357482
STASHDEV-1690-3	dpinn	liquibase/r1_3/m02.xml	2021-06-01 07:59:24.200914	170	EXECUTED	8:d69db24fbeea06f9891d619d2a7162b4	addForeignKeyConstraint baseTableName=sta_watcher, constraintName=fk_sta_watcher_user, referencedTableName=stash_user		\N	3.6.1	\N	\N	2534357482
STASHDEV-1697-1	mstudman	liquibase/r1_3/m03.xml	2021-06-01 07:59:24.205788	171	EXECUTED	8:3f31d888b661b37970925b9370845a5e	addColumn tableName=sta_pr_participant	Add a new column to sta_pr_participant to record whether the participant has approved the pull request	\N	3.6.1	\N	\N	2534357482
STASHDEV-1697-2	mstudman	liquibase/r1_3/m03.xml	2021-06-01 07:59:24.210125	172	EXECUTED	8:a977a443816130da0d085c628e09f341	addNotNullConstraint columnName=pr_approved, tableName=sta_pr_participant	Add not null constraint on sta_pr_participant.pr_approved	\N	3.6.1	\N	\N	2534357482
STASHDEV-1739-1	pepoirot	liquibase/r1_3/m03.xml	2021-06-01 07:59:24.215207	173	EXECUTED	8:6a5dd7becda6dfd5b8cbb4a103add491	dropForeignKeyConstraint baseTableName=trusted_app_restriction, constraintName=fk_trusted_app	Removes the foreign key relationship from trusted_app_restriction to trusted_app, to be able\n            to recreate it and enable deletion cascading.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1739-2	pepoirot	liquibase/r1_3/m03.xml	2021-06-01 07:59:24.219435	174	EXECUTED	8:fea8c2d9cb7d86d3024eab2aaf42e177	addForeignKeyConstraint baseTableName=trusted_app_restriction, constraintName=fk_trusted_app, referencedTableName=trusted_app	Enables deletion cascading, so that deleting a trusted application also removes the associated\n            restrictions.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-1	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.224528	175	EXECUTED	8:932f898450ccf5a90b691ee666966b52	dropIndex indexName=idx_sta_diff_comment_anchors, tableName=sta_diff_comment_anchor	This index prioritised the to_path above the from_hash. Actual usage of the index suggests searching by\n            to_hash, from_hash and then path will be more efficient, because from_hash is always provided but to_path\n            may not be (w...	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-2	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.230859	176	EXECUTED	8:a329430f888e0e29b851f1bf6982f698	createIndex indexName=idx_sta_diff_comment_from_hash, tableName=sta_diff_comment_anchor	Create a single-column index on the from_hash.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-3	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.236869	177	EXECUTED	8:326469bd7b88b43d0bdf02b0c95ac70e	createIndex indexName=idx_sta_diff_comment_to_hash, tableName=sta_diff_comment_anchor	Create a single-column index on the to_hash.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-4	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.245943	178	EXECUTED	8:0b1c8891544a45f17c3088f498e6cfef	createIndex indexName=idx_sta_diff_comment_to_path, tableName=sta_diff_comment_anchor	Create a single-column index on the to_path.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-5	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.250926	179	EXECUTED	8:dc373770791eef8b018a34fa55f0dc48	addColumn tableName=sta_diff_comment_anchor	Add a discriminator column, allowing subclasses of an InternalDiffCommentAnchor.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-6	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.25691	180	EXECUTED	8:6660618715bcf550cd7c77f91386f082	update tableName=sta_diff_comment_anchor	Set the discriminator to 1, the value for a simple InternalDiffCommentAnchor, for every existing row.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-7	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.261642	181	EXECUTED	8:d7edb89e2c1499768df788b97eff4bec	addNotNullConstraint columnName=anchor_type, tableName=sta_diff_comment_anchor	Set the NOT NULL constraint on the anchor_type.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-8	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.27328	182	EXECUTED	8:dd85fed38bf376af92f2ba78efaeee55	createTable tableName=sta_pr_diff_comment_anchor	Create the table for the InternalPullRequestDiffCommentAnchor, a joined subtype of InternalDiffCommentAnchor\n            which adds in a reference to the pull request.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-9	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.281731	183	EXECUTED	8:fa580f2cec268ea95b1781b76091aa1d	createIndex indexName=idx_sta_pr_diff_com_anc_pr, tableName=sta_pr_diff_comment_anchor	Index InternalPullRequestDiffCommentAnchors by the pull request they belong to, and tag the orphaned flag\n            into the index as well since it will be used while calculating comment drift.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-10	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.286784	184	EXECUTED	8:7eabc1a6898043adaa56cc2c1c3e3ce0	addForeignKeyConstraint baseTableName=sta_pr_diff_comment_anchor, constraintName=fk_sta_pr_diff_com_anc_id, referencedTableName=sta_diff_comment_anchor	Add a foreign key back to the InternalDiffCommentAnchor that is the base class.\n\n            Note: It's safe to cascade this deletion because if the parent row in sta_diff_comment anchor is being\n            deleted we want to remove any children ...	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-11	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.295346	185	EXECUTED	8:4287b278ba59fab8ed0aaaba064647d6	addForeignKeyConstraint baseTableName=sta_pr_diff_comment_anchor, constraintName=fk_sta_pr_diff_com_anc_pr, referencedTableName=sta_pull_request	Add a foreign key to the pull request the InternalPullRequestDiffCommentAnchor belongs to.\n\n            Note: Because this table is a "child" of sta_diff_comment_anchor, this foreign key cannot\n            cascade deletions; it would leave orphane...	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-12b	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.300182	186	EXECUTED	8:dfdf0007c39d4f08035c1dd385bda3f2	sql	Join between the sta_pr_activity and sta_pr_comment_activity table to populate sta_pr_diff_comment_anchor\n            for all existing anchors. This version is for Postgres, which uses a boolean column to for booleans.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-13	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.304718	187	EXECUTED	8:fe4f147acdebccf6464184337dc490ca	dropForeignKeyConstraint baseTableName=sta_pr_activity, constraintName=fk_sta_pr_activity_pr	Drop the previous foreign key constraint between sta_pr_activity and sta_pull_request. It was created\n            with ON DELETE CASCADE, but that is not valid for "child" tables. If the pull request is deleted, it\n            will leave phantom r...	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-14	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.310007	188	EXECUTED	8:1664de6dcad65acf5c74ff0c3042cf22	addForeignKeyConstraint baseTableName=sta_pr_activity, constraintName=fk_sta_pr_activity_pr, referencedTableName=sta_pull_request	Recreate the fk_sta_pr_activity_pr foreign key without an ON DELETE CASCADE clause.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-15	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.314823	189	EXECUTED	8:f4871836bc6391b7466fc37c98a887a3	dropForeignKeyConstraint baseTableName=sta_pr_comment_activity, constraintName=fk_sta_pr_com_act_comment	Drop the previous foreign key constraint between sta_pr_comment_activity and sta_comment. It was created\n            with ON DELETE CASCADE, but that is not valid for "child" tables. If the comment is deleted it will leave\n            phantom rows...	\N	3.6.1	\N	\N	2534357482
STASHDEV-1325-16	bturner	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.319853	190	EXECUTED	8:828d75957dd14a4494532b2e89ed948d	addForeignKeyConstraint baseTableName=sta_pr_comment_activity, constraintName=fk_sta_pr_com_act_comment, referencedTableName=sta_comment	Recreate the fk_sta_pr_com_act_comment foreign key without an ON DELETE CASCADE clause.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1700-1	mheemskerk	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.324149	191	EXECUTED	8:e50ef3151843e3d229a6e2b04e3e0227	dropForeignKeyConstraint baseTableName=sta_comment, constraintName=fk_sta_comment_parent	Drop the self referential foreign key constraint on sta_comment for the parent_id column, so we can recreate\n            it with onDelete=cascade. The constraint causes problems when comments are bulk deleted on in MySQL. MySQL\n            applies...	\N	3.6.1	\N	\N	2534357482
STASHDEV-1700-2	mheemskerk	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.32884	192	EXECUTED	8:7ce5857a79d122c213422a6f15610ab2	dropForeignKeyConstraint baseTableName=sta_comment, constraintName=fk_sta_comment_root	Drop the self referential foreign key constraint on sta_comment for the root_id column, so we can recreate\n            it with onDelete=cascade. The constraint causes problems when comments are bulk deleted on in MySQL. MySQL\n            applies t...	\N	3.6.1	\N	\N	2534357482
STASHDEV-1700-3	mheemskerk	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.333343	193	EXECUTED	8:7c3157f7234feb82cc369647c441a582	addForeignKeyConstraint baseTableName=sta_comment, constraintName=fk_sta_comment_parent, referencedTableName=sta_comment	Recreate the self referential foreign key constraint on sta_comment for the parent_id column with\n            it with onDelete=cascade. Unfortunately, self-referential cascading deletes are not supported for MSSQL.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1700-4	mheemskerk	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.337395	194	EXECUTED	8:1689d10af48b9a031fff62aacb0d97d3	addForeignKeyConstraint baseTableName=sta_comment, constraintName=fk_sta_comment_root, referencedTableName=sta_comment	Recreate the self referential foreign key constraint on sta_comment for the root_id column with\n            it with onDelete=cascade. Unfortunately, self-referential cascading deletes are not supported for MSSQL.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1771-1	mstudman	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.34635	195	EXECUTED	8:d0342be1585441ea1b0a71b2520efcfe	createTable tableName=sta_pr_merge_activity	Create the table for the InternalPullRequestMergeActivity, a joined subtype of InternalPullRequestActivity\n            which adds the possibly null hash which merged the pull request.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1771-2	mstudman	liquibase/r1_3/m04.xml	2021-06-01 07:59:24.351315	196	EXECUTED	8:b3d3a78a3839a790de380d8adca825ee	addForeignKeyConstraint baseTableName=sta_pr_merge_activity, constraintName=fk_sta_pr_mrg_act_id, referencedTableName=sta_pr_activity		\N	3.6.1	\N	\N	2534357482
STASHDEV-1674-1	bturner	liquibase/r1_3/m05.xml	2021-06-01 07:59:24.358735	197	EXECUTED	8:cf5f9e683fca4e3e15faea6070e6323b	createTable tableName=sta_pr_rescope_activity	Create the table for the InternalPullRequestRescopeActivity, a joined subtype of InternalPullRequestActivity\n            which stores the current and previous from and to hashes for the pull request.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1674-2	bturner	liquibase/r1_3/m05.xml	2021-06-01 07:59:24.363735	198	EXECUTED	8:71886278a1fc867fd08c3a51189ae798	addForeignKeyConstraint baseTableName=sta_pr_rescope_activity, constraintName=fk_sta_pr_rescope_act_id, referencedTableName=sta_pr_activity		\N	3.6.1	\N	\N	2534357482
STASHDEV-1918-1	bturner	liquibase/r1_3/m08.xml	2021-06-01 07:59:24.367662	199	EXECUTED	8:9e8a0371fb407b515ddb0ce04e929a89	dropIndex indexName=idx_sta_pull_request_from_repo, tableName=sta_pull_request	Drop the previous from repository index so we can recreate it including the branch FQN.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1918-2	bturner	liquibase/r1_3/m08.xml	2021-06-01 07:59:24.373993	200	EXECUTED	8:7b6b7af50caa0da00b50a36eeea4182f	dropIndex indexName=idx_sta_pull_request_to_repo, tableName=sta_pull_request	Drop the previous to repository index so we can recreate it including the branch FQN.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1918-3	bturner	liquibase/r1_3/m08.xml	2021-06-01 07:59:24.382228	201	EXECUTED	8:3cf070993ac47e60d1b456f00f5db803	createIndex indexName=idx_sta_pull_request_from, tableName=sta_pull_request	Create a composite index on the from repository's ID and fully-qualified branch name.	\N	3.6.1	\N	\N	2534357482
STASHDEV-1918-4	bturner	liquibase/r1_3/m08.xml	2021-06-01 07:59:24.388063	202	EXECUTED	8:e29826e48bea6ea3f98da2e123c77cfb	createIndex indexName=idx_sta_pull_request_to, tableName=sta_pull_request	Create a composite index on the to repository's ID and fully-qualified branch name.	\N	3.6.1	\N	\N	2534357482
STASHDEV-2081-1	mheemskerk	liquibase/r1_3/m08.xml	2021-06-01 07:59:24.398605	203	EXECUTED	8:9de4a5363e26eabf1e12c4eb33989ebb	createTable tableName=sta_drift_request	Create a table to store the comment drift requests in the database so they survive\n            server restarts.	\N	3.6.1	\N	\N	2534357482
STASHDEV-2081-2	mheemskerk	liquibase/r1_3/m08.xml	2021-06-01 07:59:24.411081	204	EXECUTED	8:676abace03a1a0c2043415d1a5508618	addForeignKeyConstraint baseTableName=sta_drift_request, constraintName=fk_sta_drift_request_pr, referencedTableName=sta_pull_request	Create foreign key constraint from sta_drift_request.pr_id --> sta_pull_request.id	\N	3.6.1	\N	\N	2534357482
STASHDEV-2081-3	mheemskerk	liquibase/r1_3/m08.xml	2021-06-01 07:59:24.420648	205	EXECUTED	8:6abcd9acd8729432a31a3e617a400303	createTable tableName=sta_pr_rescope_request	Create a table to store the pull request rescope requests in the database so they survive\n            server restarts.	\N	3.6.1	\N	\N	2534357482
STASHDEV-2081-4	mheemskerk	liquibase/r1_3/m08.xml	2021-06-01 07:59:24.426676	206	EXECUTED	8:aa3162e29e19c39e66c56ef9d97b2272	addForeignKeyConstraint baseTableName=sta_pr_rescope_request, constraintName=fk_sta_pr_rescope_req_repo, referencedTableName=repository	Create foreign key constraint from sta_pr_rescope_request.repo_id --> repository.id	\N	3.6.1	\N	\N	2534357482
STASHDEV-2081-5	mheemskerk	liquibase/r1_3/m08.xml	2021-06-01 07:59:24.430789	207	EXECUTED	8:4dd684803bcf63207dc81b33992ce6d4	addForeignKeyConstraint baseTableName=sta_pr_rescope_request, constraintName=fk_sta_pr_rescope_req_user, referencedTableName=stash_user	Create foreign key constraint from sta_pr_rescope_request.user_id --> stash_user.id	\N	3.6.1	\N	\N	2534357482
STASHDEV-2600-01	mstudman	liquibase/r2_0/m06.xml	2021-06-01 07:59:24.434501	208	EXECUTED	8:70d6e74c50b7342de9921b51c03e376c	customChange		\N	3.6.1	\N	\N	2534357482
STASHDEV-2716-1	tbright	liquibase/r2_1/m01.xml	2021-06-01 07:59:24.441525	209	EXECUTED	8:83b1865d5e867204cfce764b5411fb2e	modifyDataType columnName=att_value, tableName=cs_attribute	Increase the attribute value size to the maximum allowable by the DBs we use. The intention is for use to\n            be able to store JSON values here. The limit is 4000 single byte chars in Oracle.	\N	3.6.1	\N	\N	2534357482
STASHDEV-962-1	bturner	liquibase/r2_1/m01.xml	2021-06-01 07:59:24.446306	210	EXECUTED	8:6e33c977214e314f78c8b7dd345ba243	addColumn tableName=project	Add a column for the ProjectType type property, allowing the creation of new project types. The column\n            defaults to nullable to allow for existing rows, which will be adjusted by the next commit.	\N	3.6.1	\N	\N	2534357482
STASHDEV-962-2	bturner	liquibase/r2_1/m01.xml	2021-06-01 07:59:24.449808	211	EXECUTED	8:7dbca4a0d5ca11caaa26fed7567edea8	update tableName=project	All existing projects are normal projects, so a blanket update will reflect this.	\N	3.6.1	\N	\N	2534357482
STASHDEV-962-3	bturner	liquibase/r2_1/m01.xml	2021-06-01 07:59:24.45351	212	EXECUTED	8:c33af9227cd99a35adc6293c71cc7858	addNotNullConstraint columnName=project_type, tableName=project	Having set a value on any existing projects, mark the column as not null.	\N	3.6.1	\N	\N	2534357482
STASHDEV-962-4	bturner	liquibase/r2_1/m01.xml	2021-06-01 07:59:24.460131	213	EXECUTED	8:1cf66f94c01512c805a144d588aa2c9b	createIndex indexName=idx_project_type, tableName=project	Create an index on the project type to allow filtering personal projects out of the project list.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3052-01	mstudman	liquibase/r2_1/p01.xml	2021-06-01 07:59:24.46314	214	EXECUTED	8:20078d00e38d3d74a6b48a5102c1d1a8	customChange	De-duplicate cs_indexer_state entries with identical (repository_id, indexer_id) values.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3052-02	mstudman	liquibase/r2_1/p01.xml	2021-06-01 07:59:24.470906	215	EXECUTED	8:1a97e1920836d982b9abfcc89b5fa0ff	addPrimaryKey constraintName=pk_cs_indexer_state, tableName=cs_indexer_state	Create a primary key on (repository_id, indexer_id) for cs_indexer_state.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3270-1	mheemskerk	liquibase/r2_2/m01.xml	2021-06-01 07:59:24.475861	216	EXECUTED	8:c58dba400408a74b04ae8254d5d84bb8	dropTable tableName=cwd_token	Stash never needed to have the CWD_TOKEN table. It's only used by Crowd for storing SSO tokens. Embedded\n            Crowd doesn't need it.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3270-2	mheemskerk	liquibase/r2_2/m01.xml	2021-06-01 07:59:24.479783	217	EXECUTED	8:b4b2eda82b1598ce73998e0d67f3d37e	dropPrimaryKey constraintName=SYS_PK_10100, tableName=cwd_application_address	remote_address_binary and remote_address_mask are part of the primary key for cwd_application_address,\n            so the primary key must be dropped before the columns are dropped.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3270-3	mheemskerk	liquibase/r2_2/m01.xml	2021-06-01 07:59:24.483615	218	EXECUTED	8:ae72e633d823980f01375c4b78ef9540	dropColumn columnName=remote_address_binary, tableName=cwd_application_address	Crowd has removed the 'mask' and 'encodedAddressBytes' properties from the Application entity. They were\n            unused for a long time and have finally been removed. This removes the column from the database.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3270-4	mheemskerk	liquibase/r2_2/m01.xml	2021-06-01 07:59:24.487767	219	EXECUTED	8:ab296d3e2d9e5f08ddc36914fb6fca40	dropColumn columnName=remote_address_mask, tableName=cwd_application_address	Crowd has removed the 'mask' and 'encodedAddressBytes' properties from the Application entity. They were\n            unused for a long time and have finally been removed. This removes the column from the database	\N	3.6.1	\N	\N	2534357482
STASHDEV-3270-5	mheemskerk	liquibase/r2_2/m01.xml	2021-06-01 07:59:24.493899	220	EXECUTED	8:331d03a3773f0d9fd24056c82182bc64	addPrimaryKey constraintName=SYS_PK_10100, tableName=cwd_application_address	Recreate the primary key on cwd_application_address now that the binary and mask columns have been dropped.	\N	3.6.1	\N	\N	2534357482
STASHDEV-2892-1	tbright	liquibase/r2_2/m01.xml	2021-06-01 07:59:24.500822	221	EXECUTED	8:95dd119c5ee7ad4164d286882f3273ee	createTable tableName=sta_configured_hook_status		\N	3.6.1	\N	\N	2534357482
STASHDEV-2892-2	tbright	liquibase/r2_2/m01.xml	2021-06-01 07:59:24.508161	222	EXECUTED	8:86e1a68a05c2eaf11e02668ced7497db	createIndex indexName=idx_sta_config_hook_status_pk, tableName=sta_configured_hook_status		\N	3.6.1	\N	\N	2534357482
STASHDEV-2892-3	tbright	liquibase/r2_2/m01.xml	2021-06-01 07:59:24.517967	223	EXECUTED	8:fa75456c75e3fe224f5f86d9732cab7e	addUniqueConstraint constraintName=uq_sta_config_hook_status_key, tableName=sta_configured_hook_status		\N	3.6.1	\N	\N	2534357482
STASHDEV-2916-1	cofarrell_tbright	liquibase/r2_2/m01.xml	2021-06-01 07:59:24.531439	224	EXECUTED	8:8fbc23fa2b23f8998cd8344de5ee43cd	createTable tableName=sta_repo_settings		\N	3.6.1	\N	\N	2534357482
STASHDEV-2916-2	cofarrell	liquibase/r2_2/m01.xml	2021-06-01 07:59:24.546977	225	EXECUTED	8:1767226ce461ddff105cf2ddd0920197	addUniqueConstraint constraintName=uq_sta_repo_settings_key, tableName=sta_repo_settings		\N	3.6.1	\N	\N	2534357482
STASHDEV-3474-01	tbright	liquibase/r2_2/p01.xml	2021-06-01 07:59:24.550589	226	EXECUTED	8:28a4adb5eef6e2a0d67c63c1ea8b824a	customChange	De-duplicate cs_repo_membership entries with identical (cs_id, repository_id) values.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3474-02	tbright	liquibase/r2_2/p01.xml	2021-06-01 07:59:24.561372	227	EXECUTED	8:27540be017043893c44270673dea8b61	addPrimaryKey constraintName=pk_cs_repo_membership, tableName=cs_repo_membership	Create a primary key on (cs_id, repository_id) for cs_repo_membership.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3474-03	tbright	liquibase/r2_2/p01.xml	2021-06-01 07:59:24.566634	228	EXECUTED	8:ed554521ba8f716555a0d7fa3c6283a3	createIndex indexName=idx_cs_repo_membership_repo_id, tableName=cs_repo_membership	Create an index key on repository_id for cs_repo_membership.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3283-1	bturner	liquibase/r2_3/m01.xml	2021-06-01 07:59:24.570389	229	EXECUTED	8:4525e65fdb726b0839be08352141c99e	dropColumn columnName=status, tableName=repository	The system has not been setting status messages for multiple releases. Remove the column from the\n            table entirely, since it is not really used.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3283-2	bturner	liquibase/r2_3/m01.xml	2021-06-01 07:59:24.575295	230	EXECUTED	8:6b7d34102d13d033d2eb3f78a58e3fee	addColumn tableName=repository	Add a column for tracking repository hierarchies. It must be initially nullable, and will be marked\n            non-null in a subsequent changeset after it is populated.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3283-3	bturner	liquibase/r2_3/m01.xml	2021-06-01 07:59:24.580655	231	EXECUTED	8:b102e0f8d69b3b1655b81fcc6a7545f3	customChange	Run a custom change to set hierarchy IDs for all repositories.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3283-4	bturner	liquibase/r2_3/m01.xml	2021-06-01 07:59:24.585022	232	EXECUTED	8:2bbcde72917a9d0ad6b3c2d404b24b83	addNotNullConstraint columnName=hierarchy_id, tableName=repository	Switch hierarchy_id to non-null.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3283-5	bturner	liquibase/r2_3/m01.xml	2021-06-01 07:59:24.590389	233	EXECUTED	8:b308cf3d7a6f5e73d9820901191a7e96	createIndex indexName=idx_repository_hierarchy_id, tableName=repository	Create index for hierarchy ID, since it wil be used frequently to load all repositories in a hierarchy\n            when creating pull requests.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3283-6	bturner	liquibase/r2_3/m01.xml	2021-06-01 07:59:24.59963	234	EXECUTED	8:d8b7530cc5470ec196c74aca802d245c	createIndex indexName=idx_repository_origin_id, tableName=repository	Create index for origin ID, for use retrieving forks of a given repository.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3283-7	bturner	liquibase/r2_3/m01.xml	2021-06-01 07:59:24.605065	235	EXECUTED	8:5a1e01a2489277f8154f07609290b9c5	createIndex indexName=idx_repository_project_id, tableName=repository	Create index for repository project ID, which is used in almost every query against the repository table\n            to restrict results to a single project (usually joined by key).	\N	3.6.1	\N	\N	2534357482
STASHDEV-3283-8	bturner	liquibase/r2_3/m01.xml	2021-06-01 07:59:24.61444	236	EXECUTED	8:8f57fed1b1232d8e049a4596de8fa37d	createIndex indexName=idx_repository_state, tableName=repository	Create index for repository state, which will be used to filter deleted repositories out of query results.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3283-9b	bturner	liquibase/r2_3/m01.xml	2021-06-01 07:59:24.618154	237	MARK_RAN	8:3cd5f1fb0ffdd0232a66a288d811c05b	dropIndex indexName=idx_project_key, tableName=project	This index is redundant; it's covered by a unique constraint. The changeset that created it has been\n            removed from the changelog.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3283-10b	bturner	liquibase/r2_3/m01.xml	2021-06-01 07:59:24.621778	238	MARK_RAN	8:6d42e8d26af5682a44fa9366fcf25f63	dropIndex indexName=idx_project_name, tableName=project	This index is redundant; it's covered by a unique constraint. The changeset that created it has been\n            removed from the changelog.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3283-11	bturner	liquibase/r2_3/m01.xml	2021-06-01 07:59:24.629332	239	EXECUTED	8:07ba2dde67e90a6aada3e17076bee0ea	dropForeignKeyConstraint baseTableName=cs_indexer_state, constraintName=FK_CS_INDEXER_STATE_REPOSITORY	Drop the foreign key between cs_indexer_state and repository so it can be modified to cascade deletion.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3283-12	bturner	liquibase/r2_3/m01.xml	2021-06-01 07:59:24.633718	240	EXECUTED	8:6f2392c8446a29ed7f4bea2621b4aa51	addForeignKeyConstraint baseTableName=cs_indexer_state, constraintName=fk_cs_indexer_state_repository, referencedTableName=repository	Re-add fk_cs_indexer_state_repository with ON DELETE CASCADE.	\N	3.6.1	\N	\N	2534357482
STASH-3195-1a	bturner	liquibase/r2_3/m01.xml	2021-06-01 07:59:24.641374	241	EXECUTED	8:7084550ae2f069d763b2d5e7ee3efcbe	modifyDataType columnName=name, tableName=repository	Increase the repository name limit from 64 characters to 128.	\N	3.6.1	\N	\N	2534357482
STASH-3195-2a	bturner	liquibase/r2_3/m01.xml	2021-06-01 07:59:24.651548	242	EXECUTED	8:3cec69de3567ab92b2c1bfca504e4f3e	modifyDataType columnName=slug, tableName=repository	Increase the repository slug limit from 64 characters to 128, to match the new name limit.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3619-1	bturner	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.657124	243	EXECUTED	8:9dbc421e6a081ec7dbc3d209abbe8c5e	addColumn tableName=repository	Add a column to the "repository" table for the forkable flag.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3619-2	bturner	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.666649	244	EXECUTED	8:1f78076a6ff73949b1b894de707056b0	update tableName=repository	Mark all existing repositories as forkable by default.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3619-3	bturner	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.67933	245	EXECUTED	8:6c1c393b29f53042e8b3a3fa1fe0a762	addNotNullConstraint columnName=is_forkable, tableName=repository	After setting the default value, mark "repository"."is_forkable" as NOT NULL.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-1	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.684044	246	EXECUTED	8:9abda441272726cdc55b60900855c0ce	renameTable newTableName=sta_permission_type, oldTableName=weighted_permission	Rename weighted_permission to sta_permission_type	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-2	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.690522	247	EXECUTED	8:a092a30d5ab4cf8b81076f1f326ee6ff	createTable tableName=sta_global_permission	Create the table that will receive the global permissions from the 'granted_permission' table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-3	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.703242	248	EXECUTED	8:22641d61be196a999f58802d631970a4	addForeignKeyConstraint baseTableName=sta_global_permission, constraintName=fk_global_permission_user, referencedTableName=stash_user	Add the foreign key constraint between the 'stash_user' table and the global permission table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-4	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.709733	249	EXECUTED	8:620e7e9d24437822944d02e0d6be48c2	addForeignKeyConstraint baseTableName=sta_global_permission, constraintName=fk_global_permission_type, referencedTableName=sta_permission_type	Add the foreign key constraint between the 'sta_permission_type' table and the global permission table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-5	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.719056	250	EXECUTED	8:66a7f798f4bc2ea31e44e87cb7fd4fce	createIndex indexName=idx_global_permission_user, tableName=sta_global_permission	Add an index to the 'user_id' column on in the 'sta_global_permission' table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-6	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.727193	251	EXECUTED	8:33fa19c878988300509af4ab9f75b6a4	createIndex indexName=idx_global_permission_group, tableName=sta_global_permission	Add an index to the 'group_name' column on in the 'sta_global_permission' table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-7	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.733986	252	EXECUTED	8:3596473cc50f89bb41beffddec35eb26	createTable tableName=sta_project_permission	Create the table that will receive the project permissions from the 'granted_permission' table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-8	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.741468	253	EXECUTED	8:4ef88a542f32ce00c17a41a4d9fd31a4	addForeignKeyConstraint baseTableName=sta_project_permission, constraintName=fk_project_permission_user, referencedTableName=stash_user	Add the foreign key constraint between the 'stash_user' table and the project permission table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-9	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.748473	254	EXECUTED	8:3f999c1e97324ecdd14935a7c05f773f	addForeignKeyConstraint baseTableName=sta_project_permission, constraintName=fk_project_permission_project, referencedTableName=project	Add the foreign key constraint between the 'project' table and the repository permission table.\n            Deleting a repository will be cascaded to the associated project permissions.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-10	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.75647	255	EXECUTED	8:370a011e9e47ef74f15ab910ee0b53bf	addForeignKeyConstraint baseTableName=sta_project_permission, constraintName=fk_project_permission_weight, referencedTableName=sta_permission_type	Add the foreign key constraint between the 'sta_permission_type' table and the project permission table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-11	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.763835	256	EXECUTED	8:bedc12d224c735a6ade37d6710810628	createIndex indexName=idx_project_permission_user, tableName=sta_project_permission	Add an index to the 'user_id' column on in the 'sta_project_permission' table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-12	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.771349	257	EXECUTED	8:ef3f1c50e82822dc11c826715c664d0b	createIndex indexName=idx_project_permission_group, tableName=sta_project_permission	Add an index to the 'group_name' column on in the 'sta_project_permission' table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-13	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.780243	258	EXECUTED	8:4cd351a17f10e995704f831a3a365edf	createTable tableName=sta_repo_permission	Create the table that will receive the new repository permissions.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-14	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.786508	259	EXECUTED	8:4914355927d06a577784f99b79e30d79	addForeignKeyConstraint baseTableName=sta_repo_permission, constraintName=fk_repo_permission_user, referencedTableName=stash_user	Add the foreign key constraint between the 'stash_user' table and the repository permission table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-15	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.793504	260	EXECUTED	8:e8979c8f52f7f5def1629ba53dcffe81	addForeignKeyConstraint baseTableName=sta_repo_permission, constraintName=fk_repo_permission_repo, referencedTableName=repository	Add the foreign key constraint between the 'repository' table and the repository permission table.\n            Deleting a repository will be cascaded to the associated repository permissions.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-16	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.812616	261	EXECUTED	8:96e1bef37fb6c1dfa7485db7bc2f9d39	addForeignKeyConstraint baseTableName=sta_repo_permission, constraintName=fk_repo_permission_weight, referencedTableName=sta_permission_type	Add the foreign key constraint between the 'sta_permission_type' table and the repository permission table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-17	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.821193	262	EXECUTED	8:f904382404b804aced899a1762e9fdb5	createIndex indexName=idx_repo_permission_user, tableName=sta_repo_permission	Add an index to the 'user_id' column on in the 'sta_repo_permission' table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-18	jhinch	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.828206	263	EXECUTED	8:1076541e260ac271681b70f6e1df829f	createIndex indexName=idx_repo_permission_group, tableName=sta_repo_permission	Add an index to the 'group_name' column on in the 'sta_repo_permission' table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-19	pepoirot	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.833399	264	EXECUTED	8:3305f8908210b2102403c171b9796823	sql	Migrate the global permissions in 'granted_permission' to the new 'sta_global_permission' table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-20	pepoirot	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.83851	265	EXECUTED	8:2bad0660793c97e2d8abfacf14aae56b	sql	Migrate the project permissions in 'granted_permission' to the new 'sta_project_permission' table.\n            Due to an existing bug which allowed you do grant REPO_* level permissions on a project we explicitly select\n            only project pe...	\N	3.6.1	\N	\N	2534357482
STASHDEV-3458-21	pepoirot	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.844303	266	EXECUTED	8:5945d39aea6f3ef2287a41478645ce78	dropTable tableName=granted_permission	Drop the 'granted_permission' table.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3567-1	bturner	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.854608	267	EXECUTED	8:ce287e733a5c08987b1bd9a88708fe96	createTable tableName=sta_normal_project	Create the new sta_normal_project table	\N	3.6.1	\N	\N	2534357482
STASHDEV-3567-2	bturner	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.866425	268	EXECUTED	8:2216af15ed92eca3c17919427a5053e2	addForeignKeyConstraint baseTableName=sta_normal_project, constraintName=fk_sta_normal_project_id, referencedTableName=project	Create a cascading foreign key from sta_normal_project back to project	\N	3.6.1	\N	\N	2534357482
STASHDEV-3567-3	bturner	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.872358	269	EXECUTED	8:57576b4e4ac28d33d3b876154d8cda92	sql	Populate the sta_normal_project table with all of the rows from project that identify normal projects.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3567-4	bturner	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.884008	270	EXECUTED	8:4498e86140fa05a08c564d9463fc490f	createTable tableName=sta_personal_project	Create the new sta_personal_project table	\N	3.6.1	\N	\N	2534357482
STASHDEV-3567-5	bturner	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.891064	271	EXECUTED	8:d43b4b063c031eb3066fda8b25d63927	addForeignKeyConstraint baseTableName=sta_personal_project, constraintName=fk_sta_personal_project_id, referencedTableName=project	Create a cascading foreign key from sta_personal_project back to project	\N	3.6.1	\N	\N	2534357482
STASHDEV-3567-6	bturner	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.901171	272	EXECUTED	8:b1790a2e1598ab5d6e0e7b45c46a7932	addForeignKeyConstraint baseTableName=sta_personal_project, constraintName=fk_sta_personal_project_owner, referencedTableName=stash_user	Create a foreign key from sta_personal_project to its owner in stash_user	\N	3.6.1	\N	\N	2534357482
STASHDEV-3567-7	bturner	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.905803	273	EXECUTED	8:48f2f709f0aac795f713fe3dfbc4f26b	customChange	Use a custom change to populate sta_personal_project from projects and set owner links.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3734-1	mstudman	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.920078	274	EXECUTED	8:2706c03d3a16c98442075827c7477d67	addColumn tableName=stash_user	Adds the column stash_user.slug, initially nullable and non-unique, for tracking a user's slug.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3734-2	mstudman	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.931611	275	EXECUTED	8:a9f5dd6b90621777397871038ce814cd	update tableName=stash_user	Default all user slugs to the username, which is already unique.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3734-3	mstudman	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.936782	276	EXECUTED	8:2a4997a181665711a0be1b63b2366ac7	addNotNullConstraint columnName=slug, tableName=stash_user	Adds not-null constraint on stash_user.slug now that all rows have values	\N	3.6.1	\N	\N	2534357482
STASHDEV-3734-4	mstudman	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.955881	277	EXECUTED	8:6dec0d80f5edb517e0cd896b29c4dc8b	addUniqueConstraint constraintName=uq_stash_user_slug, tableName=stash_user	Adds uniqueness constraint on stash_user.slug now that all rows have values	\N	3.6.1	\N	\N	2534357482
STASHDEV-3734-6	mstudman	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.964087	278	EXECUTED	8:da9ace2d75638f2fc5dda485591b9454	customChange	Update any stash_user rows where the name is not an appropriate slug, ensuring a valid\n            slug (computed from the name) is set.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3994-2	bturner	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.979486	279	EXECUTED	8:e070ad6f651ec8638502b5764612ae0f	modifyDataType columnName=slug, tableName=stash_user	Shrink the user slug column from 255 characters to 127. SetStashUserSlug should have already ensured there\n            are no rows left in the database with long values (H2, Oracle and PostgreSQL only)	\N	3.6.1	\N	\N	2534357482
STASHDEV-3994-5	bturner	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.985669	280	EXECUTED	8:8c21daed536b10943af7678631c0c2f9	modifyDataType columnName=project_key, tableName=project	Increase the project key from 64 characters to 128.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3994-7	bturner	liquibase/r2_4/upgrade.xml	2021-06-01 07:59:24.995559	281	EXECUTED	8:b7eab1a142e3086dd50b1978aefb805b	modifyDataType columnName=name, tableName=project	Increase the project name from 64 characters to 128.	\N	3.6.1	\N	\N	2534357482
STASHDEV-4154-2	jhinch	liquibase/r2_5/upgrade.xml	2021-06-01 07:59:25.005244	283	EXECUTED	8:5a88ae75cee4ca501a42bfc99f852331	update tableName=repository	Mark all existing repositories as private by default.	\N	3.6.1	\N	\N	2534357482
STASHDEV-4154-3	jhinch	liquibase/r2_5/upgrade.xml	2021-06-01 07:59:25.010927	284	EXECUTED	8:90607e609dcab4470516d031e9b3bd2f	addNotNullConstraint columnName=is_public, tableName=repository	After setting the default value, mark "repository"."is_public" as NOT NULL.	\N	3.6.1	\N	\N	2534357482
STASHDEV-4154-4	jhinch	liquibase/r2_5/upgrade.xml	2021-06-01 07:59:25.017708	285	EXECUTED	8:88c2c53d9483b76d3511ab1de878e4b5	addColumn tableName=sta_normal_project	Add a column to the "sta_normal_project" table for the public flag.	\N	3.6.1	\N	\N	2534357482
STASHDEV-4154-5	jhinch	liquibase/r2_5/upgrade.xml	2021-06-01 07:59:25.021287	286	EXECUTED	8:01c07376c8913be0717754009c954336	update tableName=sta_normal_project	Mark all existing projects as private by default.	\N	3.6.1	\N	\N	2534357482
STASHDEV-4154-6	jhinch	liquibase/r2_5/upgrade.xml	2021-06-01 07:59:25.026227	287	EXECUTED	8:e8a87d25c861ddcdd20b490574150553	addNotNullConstraint columnName=is_public, tableName=sta_normal_project	After setting the default value, mark "sta_normal_project"."is_public" as NOT NULL.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3319-1	bturner	liquibase/r2_5/upgrade.xml	2021-06-01 07:59:25.032225	288	EXECUTED	8:fb927ab11f76131cb411a997a3888237	addColumn tableName=sta_pr_rescope_activity	Add commits_added and commits_removed columns for tracking the total number of commits added and removed\n            by a rescope activity. These may be null, in which case they're calculated on retrieval.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3319-2	bturner	liquibase/r2_5/upgrade.xml	2021-06-01 07:59:25.040806	289	EXECUTED	8:16c05e5c38a06af65588af690035b19d	createTable tableName=sta_pr_rescope_commit	Create sta_pr_rescope_commit table for recording the IDs of commits added and removed by a rescope activity.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3319-3	bturner	liquibase/r2_5/upgrade.xml	2021-06-01 07:59:25.051133	290	EXECUTED	8:28c262e5cea9a0ba248f52ea3a432165	createIndex indexName=idx_sta_pr_rescope_cmmt_act, tableName=sta_pr_rescope_commit	Add an index on the rescope activity ID, which will be used by the foreign key.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3319-4	bturner	liquibase/r2_5/upgrade.xml	2021-06-01 07:59:25.057727	291	EXECUTED	8:eefb75e57496d356d1308fa0d9da03fe	addForeignKeyConstraint baseTableName=sta_pr_rescope_commit, constraintName=fk_sta_pr_rescope_cmmt_act, referencedTableName=sta_pr_rescope_activity	Create a foreign key from rescope commits to their rescope activity, with cascading deletion.	\N	3.6.1	\N	\N	2534357482
STASH-3884-1	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.064302	292	EXECUTED	8:0cabc7b6ab3778563fe63001e1141b49	addColumn tableName=cwd_application	Add cwd_application.is_active column to replace the existing cwd_application.active column.	\N	3.6.1	\N	\N	2534357482
STASH-3884-2	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.069081	293	EXECUTED	8:0abf83088f4072b2e11967e298d4186d	update tableName=cwd_application	Populate cwd_application.is_active from active by trimming the trailing spaces.	\N	3.6.1	\N	\N	2534357482
STASH-3884-3	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.080875	294	EXECUTED	8:203e01cb3284de691d6bae2c0a1a473d	addNotNullConstraint columnName=is_active, tableName=cwd_application	Add NOT NULL constraint on cwd_application.is_active.	\N	3.6.1	\N	\N	2534357482
STASH-3804-5	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.085936	295	EXECUTED	8:55dc5da920bd1891e184ca6450bd0815	dropColumn columnName=active, tableName=cwd_application	Drop the cwd_application.active column. The Hibernate mapping now expects is_active instead.	\N	3.6.1	\N	\N	2534357482
STASH-3884-6	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.092632	296	EXECUTED	8:60495179f953e89d0bfad96c263409c1	createIndex indexName=idx_app_active, tableName=cwd_application	Add index on cwd_application.is_active, and fix its case while we're at it.	\N	3.6.1	\N	\N	2534357482
STASH-3884-7	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.098815	297	EXECUTED	8:9f4e43beb60f86f4f0dd431fcb2d617f	addColumn tableName=cwd_directory	Add cwd_directory.is_active column to replace the existing cwd_directory.active column.	\N	3.6.1	\N	\N	2534357482
STASH-3884-8	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.105341	298	EXECUTED	8:e595ba4b37bd83be594ea0546c0e4ed9	update tableName=cwd_directory	Populate cwd_directory.is_active from active by trimming the trailing spaces.	\N	3.6.1	\N	\N	2534357482
STASH-3884-9	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.120525	299	EXECUTED	8:b97e8903edeec96c292815c57263a44e	addNotNullConstraint columnName=is_active, tableName=cwd_directory	Add NOT NULL constraint on cwd_directory.is_active.	\N	3.6.1	\N	\N	2534357482
STASH-3804-11	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.125848	300	EXECUTED	8:beb73765b2d4327a48c89c2c03d00d42	dropColumn columnName=active, tableName=cwd_directory	Drop the cwd_directory.active column. The Hibernate mapping now expects is_active instead.	\N	3.6.1	\N	\N	2534357482
STASH-3884-12	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.132807	301	EXECUTED	8:ecafc5298d3c3a3e35b39002dea41c70	createIndex indexName=idx_dir_active, tableName=cwd_directory	Add index on cwd_directory.is_active, and fix its case while we're at it.	\N	3.6.1	\N	\N	2534357482
STASH-3884-13	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.137179	302	EXECUTED	8:c866fd6303f2cbe549ff543345b3c278	addColumn tableName=cwd_app_dir_mapping	Add a cwd_app_dir_mapping.is_allow_all column to replace the existing cwd_app_dir_mapping.allow_all column.	\N	3.6.1	\N	\N	2534357482
STASH-3884-14	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.142285	303	EXECUTED	8:261a3603ce9ea2b0b3bf375143b125ac	update tableName=cwd_app_dir_mapping	Populate cwd_app_dir_mapping.is_allow_all from allow_all by trimming the trailing spaces.	\N	3.6.1	\N	\N	2534357482
STASH-3884-15	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.146901	304	EXECUTED	8:3c5177ab39cb5df27a79bcd73b269f42	addNotNullConstraint columnName=is_allow_all, tableName=cwd_app_dir_mapping	Add NOT NULL constraint on cwd_app_dir_mapping.is_allow_all.	\N	3.6.1	\N	\N	2534357482
STASH-3804-16	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.15235	305	EXECUTED	8:f8141d4429faea3d05eb1953e21256b3	dropColumn columnName=allow_all, tableName=cwd_app_dir_mapping	Drop the cwd_app_dir_mapping.allow_all column. This Hibernate mapping now expects is_allow_all instead.	\N	3.6.1	\N	\N	2534357482
STASH-3884-17	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.166375	306	EXECUTED	8:ba3556cf54241c55b90b7107173cf2f7	addColumn tableName=cwd_group	Add cwd_group.is_active column to replace the existing cwd_group.active column.	\N	3.6.1	\N	\N	2534357482
STASH-3884-18	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.169803	307	EXECUTED	8:3302a56618a01e895814ece697bd337d	update tableName=cwd_group	Populate cwd_group.is_active from active by trimming the trailing spaces.	\N	3.6.1	\N	\N	2534357482
STASH-3884-19	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.178726	308	EXECUTED	8:257f166501edf3114b364480a45a767a	addNotNullConstraint columnName=is_active, tableName=cwd_group	Add NOT NULL constraint on cwd_group.is_active.	\N	3.6.1	\N	\N	2534357482
STASH-3804-20	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.196749	309	EXECUTED	8:7f5c72ae70a3e15dcae7e1bd981775f8	dropIndex indexName=IDX_GROUP_ACTIVE, tableName=cwd_group	Drop the index on cwd_group.active so it can be recreated with new columns.	\N	3.6.1	\N	\N	2534357482
STASH-3804-21	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.212103	310	EXECUTED	8:42e0979e2f6be7c78e485017ad8a9d20	dropColumn columnName=active, tableName=cwd_group	Drop the cwd_group.active column. The Hibernate mapping now expects is_active instead.	\N	3.6.1	\N	\N	2534357482
STASH-3884-22	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.217779	311	EXECUTED	8:134b17bd0b4de44ac06a3d0e26c35b87	createIndex indexName=idx_group_active, tableName=cwd_group	Add index on cwd_group.is_active, and fix its case while we're at it.	\N	3.6.1	\N	\N	2534357482
STASH-3884-23	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.225843	312	EXECUTED	8:a955f3e29e320f631281ab813e56bf62	addColumn tableName=cwd_group	Add placeholder cwd_group.tmp_local column to replace is_local.	\N	3.6.1	\N	\N	2534357482
STASH-3884-24	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.237385	313	EXECUTED	8:739a5b14053b57ce30f686885ff968ff	update tableName=cwd_group	Populate cwd_group.tmp_local from is_local by trimming the trailing spaces.	\N	3.6.1	\N	\N	2534357482
STASH-3804-25	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.24364	314	EXECUTED	8:f657d667e7863221bcf36b3ab78da2ba	dropColumn columnName=is_local, tableName=cwd_group	Drop the existing cwd_group.is_local column to make room for renaming the tmp_local column.	\N	3.6.1	\N	\N	2534357482
STASH-3884-26	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.247861	315	EXECUTED	8:8b89243b9624cebbba5a52a2bcd97286	renameColumn newColumnName=is_local, oldColumnName=tmp_local, tableName=cwd_group	Rename cwd_group.tmp_local to is_local, effectively truncating the column from CHAR(255) to CHAR(1).	\N	3.6.1	\N	\N	2534357482
STASH-3884-27	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.253624	316	EXECUTED	8:7a64e592a72fe0a5975611f55870e909	addNotNullConstraint columnName=is_local, tableName=cwd_group	Re-add the NOT NULL constraint on cwd_group.is_local.	\N	3.6.1	\N	\N	2534357482
STASH-3884-28	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.258476	317	EXECUTED	8:1271103a6f7283b060f279150a842a64	addColumn tableName=cwd_user	Add cwd_user.is_active column to replace the existing cwd_user.active column.	\N	3.6.1	\N	\N	2534357482
STASH-3884-29	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.263413	318	EXECUTED	8:3f03f3e6eafab9e0d52ca1ae7da86f58	update tableName=cwd_user	Populate cwd_user.is_active from cwd_user.active by trimming the trailing spaces.	\N	3.6.1	\N	\N	2534357482
STASH-3884-30	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.268864	319	EXECUTED	8:53b40b339ea6b6e3841dca2edef1f85f	addNotNullConstraint columnName=is_active, tableName=cwd_user	Add NOT NULL constraint on cwd_user.is_active.	\N	3.6.1	\N	\N	2534357482
STASH-3804-31	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.276066	320	EXECUTED	8:ae894fca8106555405bf02d824dac299	dropIndex indexName=IDX_USER_ACTIVE, tableName=cwd_user	Drop the index on cwd_user.active so it can be recreated with new columns.	\N	3.6.1	\N	\N	2534357482
STASH-3804-32	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.283039	321	EXECUTED	8:4995fd22d0d72314f59603b5788e27d1	dropColumn columnName=active, tableName=cwd_user	Drop the cwd_user.active column. The Hibernate mapping now expects is_active instead.	\N	3.6.1	\N	\N	2534357482
STASH-3884-33	bturner	liquibase/r2_7/p01.xml	2021-06-01 07:59:25.289388	322	EXECUTED	8:dbae5e068525197382c63e913d270fcc	createIndex indexName=idx_user_active, tableName=cwd_user	Add index on cwd_user.is_active, and fix its case while we're at it.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5102-1	cofarrell	liquibase/r2_9/upgrade.xml	2021-06-01 07:59:25.295676	323	EXECUTED	8:2f0a3e03f00771d037eb5d96019d904f	sql	Uppercase all of the indexed Jira values to allow for queries that don't have to be case insensitive	\N	3.6.1	\N	\N	2534357482
STASHDEV-5250-1	mheemskerk	liquibase/r2_9/upgrade.xml	2021-06-01 07:59:25.303058	324	EXECUTED	8:e177c116431ab97a973dd4343062f63d	addColumn tableName=cwd_user	Add a column to the "cwd_user" table for an external identifier.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5398-1	dkordonski	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.312019	325	EXECUTED	8:11e4f485b94d7485c5fb1273de7d34c5	addColumn tableName=sta_diff_comment_anchor	Add file_type column to sta_diff_comment_anchor to support anchoring comments in diffs in source or\n            destination file	\N	3.6.1	\N	\N	2534357482
STASHDEV-5398-2	dkordonski	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.319819	326	EXECUTED	8:c49e2e2677797b860d47e46fa1564a22	update tableName=sta_diff_comment_anchor; update tableName=sta_diff_comment_anchor	Update file_type to proper values according to line_type.\n            line_type = 1 (ADDED) -> file_type = 1 (TO)\n            line_type not null -> file_type = 0 (FROM)\n            for line_type null we leave file_type null as this means non-line ...	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-1	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.332619	327	EXECUTED	8:c7293b853289ca11c4111eafda840c83	createTable tableName=sta_cmt_discussion	Create the sta_cmt_discussion table for CommitDiscussion.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-2	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.340888	328	EXECUTED	8:b4cd21dc127c090e2859d66f672767fd	createIndex indexName=idx_sta_cmt_disc_repo, tableName=sta_cmt_discussion	Create an index on discussion repositories.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-3	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.35397	329	EXECUTED	8:7f1c326b4ffea2f19eb56de1338bfb73	createIndex indexName=idx_sta_cmt_disc_cmt, tableName=sta_cmt_discussion	Create an index on discussion commit IDs.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-4	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.360879	330	EXECUTED	8:407ccd4f5b38c9bbf803acc8ea9df5fa	addUniqueConstraint constraintName=uq_sta_cmt_disc_repo_cmt, tableName=sta_cmt_discussion	Create a unique constraint ensuring multiple discussions are not started on the same commit ID within\n            a single repository.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-5	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.366351	331	EXECUTED	8:4516a1d22494abbde48e443257f0b2dc	addForeignKeyConstraint baseTableName=sta_cmt_discussion, constraintName=fk_sta_cmt_disc_repo, referencedTableName=repository	Create a foreign key between discussions and their repositories, cascading deletion to remove discussions\n            when their containing repository is deleted.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-6	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.375597	332	EXECUTED	8:7ef39975f87cdbbbe38cbfd86878674d	createTable tableName=sta_cmt_disc_participant	Create the sta_cmt_disc_participant table for tracking which users have participated in a commit discussion.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-7	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.381836	333	EXECUTED	8:cc872a3fd9012f7b9cb7ca04b78b2fd5	createIndex indexName=idx_sta_cmt_disc_part_disc, tableName=sta_cmt_disc_participant	Create an index on participant discussions to speed up processing the foreign key.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-8	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.387829	334	EXECUTED	8:557e4ba345b8f3227c9f97ede80d9bf4	createIndex indexName=idx_sta_cmt_disc_part_user, tableName=sta_cmt_disc_participant	Create an index on participant users to speed up processing the foreign key.	\N	3.6.1	\N	\N	2534357482
STASHDEV-6754-5	jhinch	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.763665	377	EXECUTED	8:4284d75b7925fc70846b247e1a946114	addPrimaryKey constraintName=pk_plugin_setting, tableName=plugin_setting	Add new primary key on 'id' column on the 'plugin_setting' table	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-9	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.393248	335	EXECUTED	8:7aaed0f3e8300d29ff32b635a40a9ce2	addForeignKeyConstraint baseTableName=sta_cmt_disc_participant, constraintName=fk_sta_cmt_disc_part_disc, referencedTableName=sta_cmt_discussion	Create a foreign key between discussion participants and the discussion, cascading deletion to remove\n            participants when discussions are deleted.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-10	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.400783	336	EXECUTED	8:6aedcc0184e786a36a2cce8b04ce1d46	addForeignKeyConstraint baseTableName=sta_cmt_disc_participant, constraintName=fk_sta_cmt_disc_part_user, referencedTableName=stash_user	Create a foreign key between discussion participants and their user. Note that this foreign key does\n            not cascade deletions as it is expected that Stash users will never be deleted.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-11	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.408228	337	EXECUTED	8:5689275c6908f44dac56e8cd215e47e8	addUniqueConstraint constraintName=uq_sta_cmt_disc_part_disc_user, tableName=sta_cmt_disc_participant	Create a unique constraint ensuring a given user is not a participant in any discussion more than once.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-12	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.417118	338	EXECUTED	8:bd6c651a28bf25682df589f79c834259	createTable tableName=sta_repo_activity	Create the sta_repo_activity table for tracking repository activity streams.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-13	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.42403	339	EXECUTED	8:02978e6e31361ad351ae9d3d5ce85bcf	createIndex indexName=idx_sta_repo_activity_repo, tableName=sta_repo_activity		\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-14	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.430645	340	EXECUTED	8:a7f7680c4cd75cd1319a613a295ee296	addForeignKeyConstraint baseTableName=sta_repo_activity, constraintName=fk_sta_repo_activity_id, referencedTableName=sta_activity	Create a foreign key between repository activities and their base activities, cascading deletion to\n            simplify deleting activities.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-15	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.436162	341	EXECUTED	8:6cf99c16359b36be3d821e3617cbb9c7	addForeignKeyConstraint baseTableName=sta_repo_activity, constraintName=fk_sta_repo_activity_repo, referencedTableName=repository	Create a foreign key between activities and their repositories. This foreign key does not cascade\n            because doing so would leave orphaned partial activities	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-16	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.444083	342	EXECUTED	8:25407a1b2d9a0391243db5308b226835	createTable tableName=sta_cmt_disc_activity	Create the sta_cmt_disc_activity table for tracking commit discussion activity streams.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-17	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.453296	343	EXECUTED	8:e39586fea23d96b55d4378466af8aa25	createIndex indexName=idx_sta_cmt_disc_act_disc, tableName=sta_cmt_disc_activity	Create an index on discussion IDs to facilitate applying the foreign key to sta_cmt_discussion.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-18	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.460123	344	EXECUTED	8:215c34e2fb21eadb215d2f7e59105103	addForeignKeyConstraint baseTableName=sta_cmt_disc_activity, constraintName=fk_sta_cmt_disc_act_id, referencedTableName=sta_repo_activity	Create a foreign key between discussion activities and their base repository activities, cascading\n            deletion to simplify deleting activities.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-19	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.467298	345	EXECUTED	8:8154154c3f24fade40540f687ecda50c	addForeignKeyConstraint baseTableName=sta_cmt_disc_activity, constraintName=fk_sta_cmt_disc_act_disc, referencedTableName=sta_cmt_discussion	Create a foreign key between discussion activities and their discussions. Note that this foreign key\n            does not cascade deletions because doing so would leave orphaned rows in other activity tables.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-20	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.485498	346	EXECUTED	8:4e1303fafa80cd72aed56b63457eefa5	createTable tableName=sta_cmt_disc_comment_activity	Create the sta_cmt_disc_comment_activity for tracking commit discussion comments in the activity stream.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-21	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.496143	347	EXECUTED	8:cc649fbda8ef95905a0c7386f89c8d1e	createIndex indexName=idx_sta_cmt_disc_com_act_anc, tableName=sta_cmt_disc_comment_activity	Create an index on anchor IDs to facilitate applying the foreign key to sta_diff_comment_anchor.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-22	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.50233	348	EXECUTED	8:d11de5b456b04866723c10fa346de894	createIndex indexName=idx_sta_cmt_disc_com_act_com, tableName=sta_cmt_disc_comment_activity	Create an index on comment IDs to facilitate applying the foreign key to sta_comment.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-23	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.507545	349	EXECUTED	8:a05e7520a03a69d61af0c9c246077365	addForeignKeyConstraint baseTableName=sta_cmt_disc_comment_activity, constraintName=fk_sta_cmt_disc_com_act_id, referencedTableName=sta_cmt_disc_activity	Create a foreign key between comment activities and their base discussion activities, cascading deletion\n            to simplify deleting activities.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-24	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.514742	350	EXECUTED	8:886f660e4a5de0998a15efaad676f69a	addForeignKeyConstraint baseTableName=sta_cmt_disc_comment_activity, constraintName=fk_sta_cmt_disc_com_act_anc, referencedTableName=sta_diff_comment_anchor	Create a foreign key between comment activities and their comment anchor, if one is set. Note that this\n            foreign key does not cascade deletions as it is expected anchors will never be deleted without deleting\n            their attached ...	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-25	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.520171	351	EXECUTED	8:ca9e75f9b23d58dc30312ba3a1a56e13	addForeignKeyConstraint baseTableName=sta_cmt_disc_comment_activity, constraintName=fk_sta_cmt_disc_com_act_com, referencedTableName=sta_comment	Create a foreign key between comment activities and their comments. Note that this foreign key does not\n            cascade deletions because doing so would leave orphaned rows in other activity tables.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-26	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.529494	352	EXECUTED	8:77f518154dcc5a459f7e793d32a02e79	createTable tableName=sta_cmt_disc_comment_anchor	Create the table for the InternalChangesetDiffCommentAnchor, a joined subtype of InternalDiffCommentAnchor\n            which adds in a reference to an AnnotatedChangeset.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-27	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.536079	353	EXECUTED	8:66461c623e219119d20ed5b747a53f7e	createIndex indexName=idx_sta_cmt_com_anc_disc, tableName=sta_cmt_disc_comment_anchor	Create an index on discussion IDs to facilitate applying the foreign key to sta_cmt_discussion.	\N	3.6.1	\N	\N	2534357482
STASHDEV-6754-6	jhinch	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.769893	378	EXECUTED	8:c2d9f8e359c9eb1b762acb689d3c62cc	addUniqueConstraint constraintName=uq_plug_setting_ns_key, tableName=plugin_setting	Add a unique constraint the 'key_name' and 'namespace' columns from 'plugin_setting'	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-28	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.541148	354	EXECUTED	8:c34d0bd4d1446265cf3184a3253275d6	addForeignKeyConstraint baseTableName=sta_cmt_disc_comment_anchor, constraintName=fk_sta_cmt_disc_com_anc_id, referencedTableName=sta_diff_comment_anchor	Create a foreign key between discussion comment anchors and their base anchor, cascading deletion\n            to simplify deleting anchors.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5791-29	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.546692	355	EXECUTED	8:fca10f0112460f591f7f558b1fd5f329	addForeignKeyConstraint baseTableName=sta_cmt_disc_comment_anchor, constraintName=fk_sta_cmt_disc_com_anc_disc, referencedTableName=sta_cmt_discussion	Create a foreign key between discussion comment anchors and their discussions. Note that this foreign\n            key does not cascade deletions between doing so would leave orphaned rows in sta_diff_comment_anchor.	\N	3.6.1	\N	\N	2534357482
STASH-2642-1	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.554657	356	EXECUTED	8:72c4234f9d8608a0ed6e34b74ce55a60	createTable tableName=sta_repo_push_activity		\N	3.6.1	\N	\N	2534357482
STASH-2642-2	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.566551	357	EXECUTED	8:296b13f82fb106eed7c6071b6c8a1a90	addForeignKeyConstraint baseTableName=sta_repo_push_activity, constraintName=fk_sta_repo_push_activity_id, referencedTableName=sta_repo_activity	Create a foreign key between push activities and their base repository activities, cascading deletion\n            to simplify deleting activities.	\N	3.6.1	\N	\N	2534357482
STASH-2642-3	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.595199	358	EXECUTED	8:63ebd45614ad24a6af705c7e3b963111	createTable tableName=sta_repo_push_ref		\N	3.6.1	\N	\N	2534357482
STASH-2642-4	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.603809	359	EXECUTED	8:caf6ea6d5bb7ff9ba65ca7c75f1c7ff3	addPrimaryKey constraintName=pk_sta_repo_push_ref, tableName=sta_repo_push_ref	On all sensible databases, create a primary key between the activity ID and ref ID. No single push\n            should ever be able to update the same ref more than once.	\N	3.6.1	\N	\N	2534357482
STASH-2642-5	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.609387	360	EXECUTED	8:fdc07ccca35fc4871831a1b315d4641f	createIndex indexName=idx_sta_repo_push_ref_activity, tableName=sta_repo_push_ref	Create an index on activity IDs to facilitate applying the foreign key to sta_repo_push_activity.	\N	3.6.1	\N	\N	2534357482
STASH-2642-7	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.619623	361	EXECUTED	8:adf900f7ac0c1226387c2073a2145109	addForeignKeyConstraint baseTableName=sta_repo_push_ref, constraintName=fk_sta_repo_push_ref_act_id, referencedTableName=sta_repo_push_activity	Create a foreign key between push activities and their base repository activities, cascading deletion\n            to simplify deleting activities.	\N	3.6.1	\N	\N	2534357482
STASH-2642-8	bturner	liquibase/r2_11/upgrade.xml	2021-06-01 07:59:25.629287	362	EXECUTED	8:1c64dbb838155a0c94bea8553620ed5f	update tableName=id_sequence; update tableName=id_sequence	Update id_sequence to make room for the change in allocation sizes for activities and rescope requests.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5511-1	mstudman	liquibase/r2_12/upgrade.xml	2021-06-01 07:59:25.636246	363	EXECUTED	8:598844df38d1e3a91c481e02fcb832d7	createTable tableName=sta_service_user	Creating the sta_service_user table for the InternalServiceUser entity.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5511-2	mstudman	liquibase/r2_12/upgrade.xml	2021-06-01 07:59:25.641313	364	EXECUTED	8:f3c8c779f54116ac165f3e02631847d7	addForeignKeyConstraint baseTableName=sta_service_user, constraintName=fk_sta_service_user_id, referencedTableName=stash_user	Create a foreign key constraint between service users and their base user.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5511-3	mstudman	liquibase/r2_12/upgrade.xml	2021-06-01 07:59:25.663516	365	EXECUTED	8:af0386e23043b44289deb05bfbedfa6b	createTable tableName=sta_normal_user	Creating the sta_normal_user table for the InternalNormalUser entity.\n            The length of the 'locale' column is just largest enough to allow for the slightly longer 'ja_JP_JP' locales.\n            This column default to nullable because the...	\N	3.6.1	\N	\N	2534357482
STASHDEV-5511-4	mstudman	liquibase/r2_12/upgrade.xml	2021-06-01 07:59:25.67182	366	EXECUTED	8:3191d6b2790ec4867bcb89a92b40e7ce	addForeignKeyConstraint baseTableName=sta_normal_user, constraintName=fk_sta_normal_user_id, referencedTableName=stash_user	Create a foreign key constraint between normal users and their base user.	\N	3.6.1	\N	\N	2534357482
STASHDEV-5511-5	mstudman	liquibase/r2_12/upgrade.xml	2021-06-01 07:59:25.682751	367	EXECUTED	8:9a4256c70acafb69d1db73e53fb9c16d	sql	Insert sta_normal_user rows for all the stash_user rows	\N	3.6.1	\N	\N	2534357482
STASHDEV-5511-7	mstudman	liquibase/r2_12/upgrade.xml	2021-06-01 07:59:25.689911	368	EXECUTED	8:da47dc19ff285b42e88c748b57c27c72	dropColumn columnName=name, tableName=stash_user	Dropping column name from table stash_user as this now lives on sta_normal_user	\N	3.6.1	\N	\N	2534357482
STASHDEV-5511-9	mstudman	liquibase/r2_12/upgrade.xml	2021-06-01 07:59:25.695941	369	EXECUTED	8:bec8118f2dfa4abaa8640bd33afb8160	dropColumn columnName=slug, tableName=stash_user	Dropping column slug from table stash_user as this now lives on sta_normal_user	\N	3.6.1	\N	\N	2534357482
STASHDEV-6470-01	cszmajda	liquibase/r3_0/upgrade.xml	2021-06-01 07:59:25.699727	370	EXECUTED	8:8741916b8e4b8b5345af87bec67047af	dropForeignKeyConstraint baseTableName=cs_repo_membership, constraintName=fk_repo_membership_repo	Drop the CS_REPO_MEMBERSHIP.FK_REPO_MEMBERSHIP_REPO foreign key constraint	\N	3.6.1	\N	\N	2534357482
STASHDEV-6470-02	cszmajda	liquibase/r3_0/upgrade.xml	2021-06-01 07:59:25.70953	371	EXECUTED	8:364f85f53e9eaaaa8c8c2f9e4c7602ce	addForeignKeyConstraint baseTableName=cs_repo_membership, constraintName=fk_repo_membership_repo, referencedTableName=repository	Add back the CS_REPO_MEMBERSHIP.FK_REPO_MEMBERSHIP_REPO foreign key constraint with ON DELETE CASCADE	\N	3.6.1	\N	\N	2534357482
STASHDEV-6116-1	mheemskerk	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.718843	372	EXECUTED	8:8a58f31661861af4a8e1b063e0b7d07d	sql	Drop sta_pr_rescope_request from id_sequence table, because sta_pr_rescope_request IDs are no\n            longer generated by Hibernate	\N	3.6.1	\N	\N	2534357482
STASHDEV-6754-1	jhinch	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.727058	373	EXECUTED	8:eb196da57af5866a178d6b6f88189f33	addColumn tableName=plugin_setting	Create a nullable 'id' column on the 'plugin_setting' table	\N	3.6.1	\N	\N	2534357482
STASHDEV-6754-2	jhinch	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.73572	374	EXECUTED	8:7f3d81fc32b42ba147f86f3df58f3958	customChange	Populate the 'id' column on the 'plugin_setting' table and seed the 'id_sequence' table	\N	3.6.1	\N	\N	2534357482
STASHDEV-6754-3	jhinch	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.740322	375	EXECUTED	8:5991e7eaf8ffe07b39472eff986e6138	addNotNullConstraint columnName=id, tableName=plugin_setting	Make 'id' column on the 'plugin_setting' table not null	\N	3.6.1	\N	\N	2534357482
STASHDEV-6754-4b	jhinch	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.754395	376	EXECUTED	8:741e17f2c4efca3a45d9aa1e375b722f	dropPrimaryKey constraintName=plugin_setting_pkey, tableName=plugin_setting	Drop the primary key on the 'key_name' and 'namespace' columns from 'plugin_setting' table.\n            The constraint was not given an explicit name and liquibase fails on Postgres 8.x when attempting\n            to retrieve the correct constrain...	\N	3.6.1	\N	\N	2534357482
STASHDEV-7129-1	mheemskerk	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.774235	379	EXECUTED	8:f0512019c3277df774930eee61d6bbc2	renameTable newTableName=id_sequence_dupes, oldTableName=id_sequence	Rename id_sequence so we can recreate it and dedupe the data	\N	3.6.1	\N	\N	2534357482
STASHDEV-7129-2	mheemskerk	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.781086	380	EXECUTED	8:4355d8669cf5ee41fa6e9bf1539fc3ee	createTable tableName=id_sequence	Recreate id_sequence table with the proper constraints	\N	3.6.1	\N	\N	2534357482
STASHDEV-7129-3	mheemskerk	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.785803	381	EXECUTED	8:28eddefeb9adcff1cb10c8a40dbf8f29	sql	Copy and dedupe the old id_sequence contents to the new id_sequence table	\N	3.6.1	\N	\N	2534357482
STASHDEV-7129-4	mheemskerk	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.790897	382	EXECUTED	8:ddb388d9b1b79374ca129aadf36ad510	dropTable tableName=id_sequence_dupes	Drop id_sequence_dupes	\N	3.6.1	\N	\N	2534357482
STASHDEV-3320-1	bturner	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.80118	383	EXECUTED	8:00c75c7cbe8a23481673e8ab4a14dbdf	createTable tableName=sta_shared_lob		\N	3.6.1	\N	\N	2534357482
STASHDEV-3320-2	bturner	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.80855	384	EXECUTED	8:f481abf3a46babf445f91742c3333622	sql		\N	3.6.1	\N	\N	2534357482
STASHDEV-3320-3	bturner	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.84866	385	EXECUTED	8:ad65ed0daf6cdbfa4329377e1ec99ef7	createTable tableName=sta_repo_hook		\N	3.6.1	\N	\N	2534357482
STASHDEV-3320-4	bturner	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.857183	386	EXECUTED	8:12e033806ff0eb8414213a000abe5ff9	createIndex indexName=idx_sta_repo_hook_hook_key, tableName=sta_repo_hook		\N	3.6.1	\N	\N	2534357482
STASHDEV-3320-5	bturner	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.865398	387	EXECUTED	8:54037fd829aceb0f6e8841202a030782	createIndex indexName=idx_sta_repo_hook_lob_id, tableName=sta_repo_hook	Create an index on the LOB ID used to store settings, for use by the foreign key.	\N	3.6.1	\N	\N	2534357482
STASHDEV-3320-6	bturner	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.872092	388	EXECUTED	8:f5b236b3f25787c6fcb904289dcd1649	createIndex indexName=idx_sta_repo_hook_repo_id, tableName=sta_repo_hook		\N	3.6.1	\N	\N	2534357482
STASHDEV-3320-7	bturner	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.886516	389	EXECUTED	8:0d7302cd8928670569e17000a70a4e06	addUniqueConstraint constraintName=uq_sta_repo_hook_repo_hook_key, tableName=sta_repo_hook		\N	3.6.1	\N	\N	2534357482
STASHDEV-3320-8	bturner	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.893835	390	EXECUTED	8:8b54d1cfae14920ae03d9564ed74494b	addForeignKeyConstraint baseTableName=sta_repo_hook, constraintName=fk_sta_repo_hook_lob, referencedTableName=sta_shared_lob		\N	3.6.1	\N	\N	2534357482
STASHDEV-3320-9	bturner	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.902481	391	EXECUTED	8:735023955aa678dd886f9dabe35d1f7e	addForeignKeyConstraint baseTableName=sta_repo_hook, constraintName=fk_sta_repo_hook_repo, referencedTableName=repository		\N	3.6.1	\N	\N	2534357482
STASHDEV-3320-10	bturner	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.906549	392	EXECUTED	8:807fe956e92e2f73d7beebbb493a1996	customChange		\N	3.6.1	\N	\N	2534357482
STASHDEV-3320-11	bturner	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.912906	393	EXECUTED	8:7c187ca694c9030b6cf3159228614af5	update tableName=id_sequence; update tableName=id_sequence		\N	3.6.1	\N	\N	2534357482
STASHDEV-3320-12	bturner	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.920309	394	EXECUTED	8:e88447faccea2aff01272bf50ef15e88	dropTable tableName=sta_configured_hook_status		\N	3.6.1	\N	\N	2534357482
STASHDEV-3320-13	bturner	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.925823	395	EXECUTED	8:9824ce318cd2bf6ab8d118adca84e834	dropTable tableName=sta_repo_settings		\N	3.6.1	\N	\N	2534357482
STASHDEV-7021-1	jpalacios	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.934196	396	EXECUTED	8:25294ec9075acbaf0ec702f0b906fde9	createTable tableName=sta_user_settings		\N	3.6.1	\N	\N	2534357482
STASHDEV-7021-4	jpalacios	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.939097	397	EXECUTED	8:4aa5476db432f533752d6c7a6a0c454d	addForeignKeyConstraint baseTableName=sta_user_settings, constraintName=fk_sta_user_settings_lob, referencedTableName=sta_shared_lob		\N	3.6.1	\N	\N	2534357482
STASHDEV-7021-5	jpalacios	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.944415	398	EXECUTED	8:1a3ef4cb2788afb330278f16591054c4	addForeignKeyConstraint baseTableName=sta_user_settings, constraintName=fk_sta_user_settings_user, referencedTableName=stash_user		\N	3.6.1	\N	\N	2534357482
STASH-4631	rstocker	liquibase/r3_2/upgrade.xml	2021-06-01 07:59:25.949615	399	EXECUTED	8:f3c1da1286f22dea224153d0b05b311f	addColumn tableName=sta_normal_user	Add column 'deleted_timestamp' to 'sta_normal_user'	\N	3.6.1	\N	\N	2534357482
STASHDEV-7551-1	jthomas,pepoirot	liquibase/r3_3/upgrade.xml	2021-06-01 07:59:25.96208	400	EXECUTED	8:26414dd3dddf2bcb38b40757a4e37e7b	createTable tableName=sta_task		\N	3.6.1	\N	\N	2534357482
STASHDEV-7780-1	jthomas	liquibase/r3_3/upgrade.xml	2021-06-01 07:59:25.96869	401	EXECUTED	8:30a2b9e9d1cf8ff2e641700dc3a5c1c1	createIndex indexName=idx_sta_task_anchor, tableName=sta_task		\N	3.6.1	\N	\N	2534357482
STASHDEV-7780-2	jthomas	liquibase/r3_3/upgrade.xml	2021-06-01 07:59:25.974932	402	EXECUTED	8:60f9ced54fce1d6771a3de77c61464bd	createIndex indexName=idx_sta_task_context, tableName=sta_task		\N	3.6.1	\N	\N	2534357482
STASHDEV-7846-1	jhinch	liquibase/r3_4/upgrade.xml	2021-06-01 07:59:25.984924	403	EXECUTED	8:13b2e9c63d37453fb0ab13c374b959b5	customChange	Enable membership aggregation if it has no chance of altering effective permissions	\N	3.6.1	\N	\N	2534357482
STASHDEV-8207-1	mheemskerk	liquibase/r3_5/upgrade.xml	2021-06-01 07:59:25.991467	404	EXECUTED	8:59e7f7872592255c353f3b673307a054	addColumn tableName=sta_pull_request	Add locked_timestamp to sta_pull_request	\N	3.6.1	\N	\N	2534357482
STASHDEV-8207-2	mheemskerk	liquibase/r3_5/upgrade.xml	2021-06-01 07:59:26.007173	405	EXECUTED	8:d58a14aa41e93d458665ad09e4e09b4c	customChange	InternalRescopeRequest is managed by Hibernate again. Re-initialize id_sequence for sta_pr_rescope_request.	\N	3.6.1	\N	\N	2534357482
STASH-4413-1	jthomas	liquibase/r3_7/upgrade.xml	2021-06-01 07:59:26.017981	406	EXECUTED	8:c898cdaf05b63026de5c688b9a366a6d	customChange	Remove orphaned memberships where membership_type is equal to 'GROUP_USER'. This _should_ remove all the duplicates.\n            If they have any duplicates left over we'll have to just delete them all. Getting in a state where\n            this wi...	\N	3.6.1	\N	\N	2534357482
STASH-4413-2	jthomas	liquibase/r3_7/upgrade.xml	2021-06-01 07:59:26.024609	407	EXECUTED	8:6c9be3b577b6141963d613d742a7db0e	addUniqueConstraint constraintName=uk_mem_dir_parent_child, tableName=cwd_membership	Add a unique constraint to "cwd_membership" for the existing index "idx_mem_dir_parent_child".	\N	3.6.1	\N	\N	2534357482
STASHDEV-8755-1	mheemskerk	liquibase/r3_8/upgrade.xml	2021-06-01 07:59:26.03447	408	EXECUTED	8:67ce65efffb39567aa6a99e3e92d1527	createTable tableName=sta_remember_me_token	Creates sta_remember_me_token table	\N	3.6.1	\N	\N	2534357482
STASHDEV-8755-2	mheemskerk	liquibase/r3_8/upgrade.xml	2021-06-01 07:59:26.040572	409	EXECUTED	8:90f4b7c5a7cbbeb4adb8a287920ecdb2	addForeignKeyConstraint baseTableName=sta_remember_me_token, constraintName=fk_remember_me_user_id, referencedTableName=stash_user	Add a foreign key constraint to stash_user	\N	3.6.1	\N	\N	2534357482
STASHDEV-8755-3	mheemskerk	liquibase/r3_8/upgrade.xml	2021-06-01 07:59:26.054985	410	EXECUTED	8:9cf77e0a4cb1788dc9f6995bd4ef65c0	addUniqueConstraint constraintName=uq_remember_me_series_token, tableName=sta_remember_me_token	Add a uniqueness constraint on sta_remember_me_token (series, token)	\N	3.6.1	\N	\N	2534357482
STASHDEV-8755-4	mheemskerk	liquibase/r3_8/upgrade.xml	2021-06-01 07:59:26.061916	411	EXECUTED	8:903b2aee3666d824fe97e75e016b4b6a	dropTable tableName=persistent_logins	Drop the old persistent_logins table	\N	3.6.1	\N	\N	2534357482
STASHDEV-8452-1	mheemskerk	liquibase/r3_8/upgrade.xml	2021-06-01 07:59:26.071981	412	EXECUTED	8:3488bd1d0878efe8d5b0ace92f391c40	createTable tableName=sta_repo_origin	Create "sta_repo_origin" table to manage the repository -> origin relationship. This is step 1 in getting rid\n            of the self-referential foreign key on the "repository" table.\n            A repository can only have a single origin. As a r...	\N	3.6.1	\N	\N	2534357482
STASHDEV-8452-2	mheemskerk	liquibase/r3_8/upgrade.xml	2021-06-01 07:59:26.08081	413	EXECUTED	8:4166024415699106a1940a727ca8959d	createIndex indexName=idx_sta_repo_origin_origin_id, tableName=sta_repo_origin	Create index on "sta_repo_origin.origin_id	\N	3.6.1	\N	\N	2534357482
STASHDEV-8452-3	mheemskerk	liquibase/r3_8/upgrade.xml	2021-06-01 07:59:26.084881	414	EXECUTED	8:d8f5d8468b6babd3edd1fd57e8fbcb83	sql	Migrate the is-origin-of relationship from the "repository" table to the "sta_repo_origin" table	\N	3.6.1	\N	\N	2534357482
STASHDEV-8452-4	mheemskerk	liquibase/r3_8/upgrade.xml	2021-06-01 07:59:26.089306	415	EXECUTED	8:5265579cb56967eef47bbdfd36a27d76	dropForeignKeyConstraint baseTableName=repository, constraintName=fk_repository_origin	Drop the foreign key constraint on repository.origin_id prior to dropping the column	\N	3.6.1	\N	\N	2534357482
STASHDEV-8452-5	mheemskerk	liquibase/r3_8/upgrade.xml	2021-06-01 07:59:26.093885	416	EXECUTED	8:b475a484414f02b59600f2780c3ee064	dropIndex indexName=idx_repository_origin_id, tableName=repository	Drop the index on repository.origin_id prior to dropping the column	\N	3.6.1	\N	\N	2534357482
STASHDEV-8452-6	mheemskerk	liquibase/r3_8/upgrade.xml	2021-06-01 07:59:26.09884	417	EXECUTED	8:0a2ce6dab98083f1fe23cae1d953e51f	dropColumn columnName=origin_id, tableName=repository	Drop the "origin_id" column from the "repository" table	\N	3.6.1	\N	\N	2534357482
STASH-7261	jpalacios	liquibase/r3_9/upgrade.xml	2021-06-01 07:59:26.104438	418	EXECUTED	8:780e1e3e91abded81559bccfd39f933d	delete tableName=sta_shared_lob		\N	3.6.1	\N	\N	2534357482
STASH-7119-1	jpalacios	liquibase/r3_9/upgrade.xml	2021-06-01 07:59:26.11751	419	EXECUTED	8:1cc696168cae2fd7b6ffa32a032dadab	createTable tableName=sta_deleted_group		\N	3.6.1	\N	\N	2534357482
STASH-7119-2	jpalacios	liquibase/r3_9/upgrade.xml	2021-06-01 07:59:26.124458	420	EXECUTED	8:f54399275bbba852cb1ef3204b5fdc1d	createIndex indexName=idx_sta_deleted_group_ts, tableName=sta_deleted_group	Create an index on the deleted timestamp to filter by date during the cleanup task.	\N	3.6.1	\N	\N	2534357482
STASHDEV-9310-1	mheemskerk	liquibase/r3_9/upgrade.xml	2021-06-01 07:59:26.146766	421	EXECUTED	8:258530fb4c11392e311391aae7a88149	addColumn tableName=sta_service_user	Add columns to the sta_service_user table	\N	3.6.1	\N	\N	2534357482
STASHDEV-9310-2	mheemskerk	liquibase/r3_9/upgrade.xml	2021-06-01 07:59:26.15274	422	EXECUTED	8:b96cb5e7c38c89177515ada553117fe6	customChange	Populates the new columns in the sta_service_user_new table	\N	3.6.1	\N	\N	2534357482
STASHDEV-9310-3	mheemskerk	liquibase/r3_9/upgrade.xml	2021-06-01 07:59:26.158368	423	EXECUTED	8:c4cd8bb105e98956a8a9fc554bf03bfe	addNotNullConstraint columnName=active, tableName=sta_service_user	Add not-null constraint to the sta_service_user.active column	\N	3.6.1	\N	\N	2534357482
STASHDEV-9310-4	mheemskerk	liquibase/r3_9/upgrade.xml	2021-06-01 07:59:26.164914	424	EXECUTED	8:9e697520b61d04ac0bbb762328a44ee3	addNotNullConstraint columnName=name, tableName=sta_service_user	Add not-null constraint to the sta_service_user.name column	\N	3.6.1	\N	\N	2534357482
STASHDEV-9310-5	mheemskerk	liquibase/r3_9/upgrade.xml	2021-06-01 07:59:26.174132	425	EXECUTED	8:7f3b9e97a7ae5fcf63969897569f2f7a	addUniqueConstraint constraintName=uq_sta_service_user_name, tableName=sta_service_user	Add unique constraint to the sta_service_user.name column	\N	3.6.1	\N	\N	2534357482
STASHDEV-9310-6	mheemskerk	liquibase/r3_9/upgrade.xml	2021-06-01 07:59:26.179154	426	EXECUTED	8:ef6f326f9dc4c2d748ec2ebe233115a2	addNotNullConstraint columnName=slug, tableName=sta_service_user	Add not-null constraint to the sta_service_user.slug column	\N	3.6.1	\N	\N	2534357482
STASHDEV-9310-7	mheemskerk	liquibase/r3_9/upgrade.xml	2021-06-01 07:59:26.186661	427	EXECUTED	8:c25e854f16c1e3ee368c8f2265963d7e	addUniqueConstraint constraintName=uq_sta_service_slug, tableName=sta_service_user	Add unique constraint to the sta_service_user.slug column	\N	3.6.1	\N	\N	2534357482
STASHDEV-9310-8	mheemskerk	liquibase/r3_9/upgrade.xml	2021-06-01 07:59:26.190585	428	EXECUTED	8:318eaec27f75421743cbcce66d86ef40	addNotNullConstraint columnName=label, tableName=sta_service_user	Add not-null constraint to the sta_service_user.label column	\N	3.6.1	\N	\N	2534357482
STASH-5244-1	bturner	liquibase/r3_10/upgrade.xml	2021-06-01 07:59:26.195498	429	EXECUTED	8:0a159bf4b281184f3a3e62beaf87de56	update tableName=cwd_user		\N	3.6.1	\N	\N	2534357482
STASH-5244-3	bturner	liquibase/r3_10/upgrade.xml	2021-06-01 07:59:26.204428	430	EXECUTED	8:ff5ccbcc138e9d56c200b84e0fba8b7f	addUniqueConstraint constraintName=uq_cwd_user_dir_ext_id, tableName=cwd_user		\N	3.6.1	\N	\N	2534357482
STASHDEV-9602-1	jpalacios	liquibase/r3_12/upgrade.xml	2021-06-01 07:59:26.213184	431	EXECUTED	8:c817f0e5d79082a146ada77ec8343f81	createIndex indexName=idx_sta_pr_update_ts, tableName=sta_pull_request	Add index on sta_pull_request.updated_timestamp	\N	3.6.1	\N	\N	2534357482
STASHDEV-9602-2	jpalacios	liquibase/r3_12/upgrade.xml	2021-06-01 07:59:26.21977	432	EXECUTED	8:80d4c9f759145c46c3a99e4fd0fdc290	createIndex indexName=idx_sta_pr_to_repo_update, tableName=sta_pull_request	Add index on sta_pull_request.to_repository_id, sta_pull_request.updated_timestamp	\N	3.6.1	\N	\N	2534357482
STASHDEV-9602-3	jpalacios	liquibase/r3_12/upgrade.xml	2021-06-01 07:59:26.227657	433	EXECUTED	8:3123fe970251c9d56bad0cdff6629f12	createIndex indexName=idx_sta_pr_from_repo_update, tableName=sta_pull_request	Add index on sta_pull_request.from_repository_id, sta_pull_request.updated_timestamp	\N	3.6.1	\N	\N	2534357482
STASH-7580	rfriend	liquibase/r3_12/upgrade.xml	2021-06-01 07:59:26.232053	434	EXECUTED	8:4f82224e7614e187f878b98db235eada	update tableName=cwd_user		\N	3.6.1	\N	\N	2534357482
STASHDEV-9922-1	bturner	liquibase/r4_0/upgrade.xml	2021-06-01 07:59:26.241734	435	EXECUTED	8:771f511faf9763ff35fa87b39810ba1b	update tableName=plugin_setting; update tableName=plugin_setting; update tableName=plugin_setting; update tableName=plugin_setting; update tableName=plugin_setting; update tableName=plugin_setting; update tableName=plugin_setting; update tableName...		\N	3.6.1	\N	\N	2534357482
STASHDEV-9922-2	bturner	liquibase/r4_0/upgrade.xml	2021-06-01 07:59:26.246377	436	EXECUTED	8:134b9199ffddb144ba187eacc6379fad	update tableName=plugin_setting		\N	3.6.1	\N	\N	2534357482
STASHDEV-10475-1	sgoodhew	liquibase/r4_0/upgrade.xml	2021-06-01 07:59:26.251302	437	EXECUTED	8:3174d5e733a00ba3b5c5a44bdc4c77b1	delete tableName=sta_remember_me_token		\N	3.6.1	\N	\N	2534357482
BSERV-8242-1a	bturner	liquibase/r4_0/p07.xml	2021-06-01 07:59:26.257451	438	EXECUTED	8:7f8f4c5416e1aaa734a2c32ff45f328f	update tableName=sta_repo_hook	Update rows with the stash-bundled-hooks key to use the new bitbucket-bundled-hooks key, unless\n            a row with that key already exists for the repository.	\N	3.6.1	\N	\N	2534357482
BSERV-8242-2	bturner	liquibase/r4_0/p07.xml	2021-06-01 07:59:26.264303	439	EXECUTED	8:36ac63f4fea7873336572bd6841b40a0	delete tableName=sta_repo_hook	Delete any remaining rows with the stash-bundled-hooks key. A row with the bitbucket-bundled-hooks\n            key must already exist for each of these rows or they would have been updated.	\N	3.6.1	\N	\N	2534357482
BSERVDEV-10768-2	jpalacios	liquibase/r4_2/upgrade.xml	2021-06-01 07:59:26.287225	440	EXECUTED	8:e9a2a6aac66dd8799837c10a68a00dc6	sql	Change the type of 'approved' flag from boolean to int for postgres	\N	3.6.1	\N	\N	2534357482
BSERVDEV-10768-4	jpalacios	liquibase/r4_2/upgrade.xml	2021-06-01 07:59:26.29399	441	EXECUTED	8:4f47d5db000075c4b83cb27cee688645	renameColumn newColumnName=participant_status, oldColumnName=pr_approved, tableName=sta_pr_participant	Rename 'approved' column to 'participant_status'	\N	3.6.1	\N	\N	2534357482
BSERVDEV-10769	jpalacios	liquibase/r4_2/upgrade.xml	2021-06-01 07:59:26.298494	442	EXECUTED	8:9a7296102022a4a07a4578e69f3802e8	update tableName=sta_pr_participant	Reset participant status for authors	\N	3.6.1	\N	\N	2534357482
BSERVDEV-10791-1	mheemskerk	liquibase/r4_2/upgrade.xml	2021-06-01 07:59:26.305379	443	EXECUTED	8:7fded9f8100558db77f5b34d9192663d	addColumn tableName=project	Add namespace column to project table for support of 3-level clone URLs for mirrors	\N	3.6.1	\N	\N	2534357482
BSERVDEV-10791-2	mheemskerk	liquibase/r4_2/upgrade.xml	2021-06-01 07:59:26.31516	444	EXECUTED	8:291452200f1982f4c2faaa55eee5ee19	addNotNullConstraint columnName=namespace, tableName=project	Add not-null constraint to namespace, initializing all null values to #	\N	3.6.1	\N	\N	2534357482
BSERVDEV-10791-3	mheemskerk	liquibase/r4_2/upgrade.xml	2021-06-01 07:59:26.329443	445	EXECUTED	8:e58dbc23d5f0d77501ff547b046dca3a	dropUniqueConstraint constraintName=uk_project_name, tableName=project	Drop uk_project_name	\N	3.6.1	\N	\N	2534357482
BSERVDEV-10791-4	mheemskerk	liquibase/r4_2/upgrade.xml	2021-06-01 07:59:26.335415	446	EXECUTED	8:51bf75698a9b77584f2e89cdb8670ec9	addUniqueConstraint constraintName=uk_project_name, tableName=project	Recreate uk_project_name on namespace,name	\N	3.6.1	\N	\N	2534357482
BSERVDEV-10791-5	mheemskerk	liquibase/r4_2/upgrade.xml	2021-06-01 07:59:26.340991	447	EXECUTED	8:7b86fb5093a37b5e43e663cfe9449af2	dropUniqueConstraint constraintName=uk_project_key, tableName=project	Drop uk_project_key	\N	3.6.1	\N	\N	2534357482
BSERVDEV-10791-6	mheemskerk	liquibase/r4_2/upgrade.xml	2021-06-01 07:59:26.348809	448	EXECUTED	8:78004962386bcb9b31a54dd28086bf23	addUniqueConstraint constraintName=uk_project_key, tableName=project	Recreate uk_project_key on namespace,project_key	\N	3.6.1	\N	\N	2534357482
BSERVDEV-10186	bturner	liquibase/r4_2/upgrade.xml	2021-06-01 07:59:26.355104	449	EXECUTED	8:600fb2db3d170c7cf318d7cfe15e7a91	createIndex indexName=idx_sta_pr_rescope_req_repo, tableName=sta_pr_rescope_request	Add an index for looking up rescope requests by repository ID.	\N	3.6.1	\N	\N	2534357482
BSERVDEV-10774-1	jpalacios	liquibase/r4_4/upgrade.xml	2021-06-01 07:59:26.364316	450	EXECUTED	8:aca994a7a695a676dc31e1b8fe6184c3	createTable tableName=bb_pr_part_status_weight	Create a weighting table for sorting by participant status	\N	3.6.1	\N	\N	2534357482
BSERVDEV-10774-2	jpalacios	liquibase/r4_4/upgrade.xml	2021-06-01 07:59:26.369389	451	EXECUTED	8:e508c4b249ec6c9ecbae4742554ed921	insert tableName=bb_pr_part_status_weight; insert tableName=bb_pr_part_status_weight; insert tableName=bb_pr_part_status_weight	Populate weighting table for sorting by participant status.\n            Mapping will be: UNAPPROVED(0) -> 0, NEEDS_WORK(2) -> 1, APPROVED(1) -> 2	\N	3.6.1	\N	\N	2534357482
BSERVDEV-10820	jpalacios	liquibase/r4_4/upgrade.xml	2021-06-01 07:59:26.373635	452	EXECUTED	8:75c92543b38526217aaede7e469ff0af	update tableName=sta_pr_participant	Promote participants who have approved a PR to reviewers	\N	3.6.1	\N	\N	2534357482
BSERVDEV-11532-1	jthomas	liquibase/r4_4/upgrade.xml	2021-06-01 07:59:26.381313	453	EXECUTED	8:dc34a1ffd7fbba6d1bc8b2bb2889437b	createTable tableName=bb_pr_reviewer_upd_activity		\N	3.6.1	\N	\N	2534357482
BSERVDEV-11532-2	jthomas	liquibase/r4_4/upgrade.xml	2021-06-01 07:59:26.388681	454	EXECUTED	8:b3c118e5c24139a7ec85e2067c16b3ac	addForeignKeyConstraint baseTableName=bb_pr_reviewer_upd_activity, constraintName=fk_bb_pr_reviewer_act_id, referencedTableName=sta_pr_activity		\N	3.6.1	\N	\N	2534357482
BSERVDEV-11532-3	jthomas	liquibase/r4_4/upgrade.xml	2021-06-01 07:59:26.399556	455	EXECUTED	8:22ce3678ab1ef3b182d68b5b39084a3f	createTable tableName=bb_pr_reviewer_added		\N	3.6.1	\N	\N	2534357482
BSERVDEV-11532-5	jthomas	liquibase/r4_4/upgrade.xml	2021-06-01 07:59:26.40572	456	EXECUTED	8:edce4dac4ef1a58ac405a35b6a718b17	addForeignKeyConstraint baseTableName=bb_pr_reviewer_added, constraintName=fk_bb_pr_reviewer_added_act, referencedTableName=bb_pr_reviewer_upd_activity		\N	3.6.1	\N	\N	2534357482
BSERVDEV-11532-6	jthomas	liquibase/r4_4/upgrade.xml	2021-06-01 07:59:26.412903	457	EXECUTED	8:26774f6fef7d980b5e0dae0637bf381a	addForeignKeyConstraint baseTableName=bb_pr_reviewer_added, constraintName=fk_bb_pr_reviewer_added_user, referencedTableName=stash_user		\N	3.6.1	\N	\N	2534357482
BSERVDEV-11532-7	jthomas	liquibase/r4_4/upgrade.xml	2021-06-01 07:59:26.425346	458	EXECUTED	8:5c9781ecd75cd1fbe79beb387c3fd27f	createTable tableName=bb_pr_reviewer_removed		\N	3.6.1	\N	\N	2534357482
BSERVDEV-11532-9	jthomas	liquibase/r4_4/upgrade.xml	2021-06-01 07:59:26.429952	459	EXECUTED	8:0dc5e5d6ecefd0949f8252016e6a423b	addForeignKeyConstraint baseTableName=bb_pr_reviewer_removed, constraintName=fk_bb_pr_reviewer_removed_act, referencedTableName=bb_pr_reviewer_upd_activity		\N	3.6.1	\N	\N	2534357482
BSERVDEV-11532-10	jthomas	liquibase/r4_4/upgrade.xml	2021-06-01 07:59:26.434024	460	EXECUTED	8:b89f88c2675049d8703b67b5616d0c7f	addForeignKeyConstraint baseTableName=bb_pr_reviewer_removed, constraintName=fk_bb_pr_reviewer_removed_user, referencedTableName=stash_user		\N	3.6.1	\N	\N	2534357482
BSERV-7216-1	mheemskerk	liquibase/r4_5/upgrade.xml	2021-06-01 07:59:26.441552	461	EXECUTED	8:5af85be29ce3f740e19f1db948c87683	createTable tableName=sta_pr_rescope_request_change	Create the sta_pr_rescope_request_change table to persist the ref-changes with the rescope requests	\N	3.6.1	\N	\N	2534357482
BSERV-7216-2	mheemskerk	liquibase/r4_5/upgrade.xml	2021-06-01 07:59:26.447768	462	EXECUTED	8:03c7b6f0d2b3788fb47b9b7ae9aefe8c	addPrimaryKey constraintName=pk_sta_pr_rescope_req_change, tableName=sta_pr_rescope_request_change	On all sensible databases, create a primary key between the rescope request ID and ref ID. No single\n            rescope trigger should ever be able to change the same ref more than once.	\N	3.6.1	\N	\N	2534357482
BSERV-7216-4	mheemskerk	liquibase/r4_5/upgrade.xml	2021-06-01 07:59:26.452218	463	EXECUTED	8:aa302a2013216f6715aafd8b87067830	addForeignKeyConstraint baseTableName=sta_pr_rescope_request_change, constraintName=fk_sta_pr_rescope_ch_req_id, referencedTableName=sta_pr_rescope_request	Create a foreign key between rescope request ref changes and the rescope request, cascading deletion\n            to simplify deleting rescope requests.	\N	3.6.1	\N	\N	2534357482
BSERV-7216-5	mheemskerk	liquibase/r4_5/upgrade.xml	2021-06-01 07:59:26.456768	464	EXECUTED	8:e438705a0db905d07797b8ebf6d0063b	dropColumn columnName=branch_fqn, tableName=sta_pr_rescope_request	Drop the branch_fqn column from the rescope request table. It has been replaced by individual ref changes.\n            After the upgrade, any persisted rescope requests will be executed as a full-repository rescope.	\N	3.6.1	\N	\N	2534357482
BSERV-7216-6	mheemskerk	liquibase/r4_5/upgrade.xml	2021-06-01 07:59:26.462316	465	EXECUTED	8:637d037c42f35ea5e0f23b51e4a354a9	addColumn tableName=sta_pr_rescope_request	Add the created_date column to track when the rescope request was created. This can be used to order the\n            requests by date to ensure they are replayed in the correct order	\N	3.6.1	\N	\N	2534357482
BSERV-7216-7	mheemskerk	liquibase/r4_5/upgrade.xml	2021-06-01 07:59:26.467268	466	EXECUTED	8:5087cfe0bed425e1240faee35c8e938f	addColumn tableName=sta_pull_request	Add the rescoped_date column to track when the scope (from/to ref) of a pull request was last updated.	\N	3.6.1	\N	\N	2534357482
BSERV-7216-8	mheemskerk	liquibase/r4_5/upgrade.xml	2021-06-01 07:59:26.471083	467	EXECUTED	8:cf2d560163738bbf664d99fda715e4ce	addNotNullConstraint columnName=created_timestamp, tableName=sta_pr_rescope_request	Add non-null constraint to sta_pr_rescope_request.created_timestamp	\N	3.6.1	\N	\N	2534357482
BSERV-7216-9	mheemskerk	liquibase/r4_5/upgrade.xml	2021-06-01 07:59:26.474875	468	EXECUTED	8:23f3e39d9594fc722ab187b66be898c8	addNotNullConstraint columnName=rescoped_timestamp, tableName=sta_pull_request	Add non-null constraint to sta_pull_request.rescoped_timestamp	\N	3.6.1	\N	\N	2534357482
BSERVDEV-11909-1	cszmajda	liquibase/r4_6/upgrade.xml	2021-06-01 07:59:26.484646	469	EXECUTED	8:9a4c7ff348cad86e19a13d49917151d6	createTable tableName=bb_clusteredjob	Create clustered job table backing atlassian-scheduler-caesium	\N	3.6.1	\N	\N	2534357482
BSERVDEV-11909-2	cszmajda	liquibase/r4_6/upgrade.xml	2021-06-01 07:59:26.492482	470	EXECUTED	8:a65b71e599d5b9463d666a4655aeb061	createIndex indexName=idx_bb_clusteredjob_jrk, tableName=bb_clusteredjob	Add an index for bb_clusteredjob.job_runner_key	\N	3.6.1	\N	\N	2534357482
BSERVDEV-11909-3	cszmajda	liquibase/r4_6/upgrade.xml	2021-06-01 07:59:26.498296	471	EXECUTED	8:66b2c681a4650d75b6709ce80c1812af	createIndex indexName=idx_bb_clusteredjob_next_run, tableName=bb_clusteredjob	Add an index for bb_clusteredjob.nextRun	\N	3.6.1	\N	\N	2534357482
BSERVDEV-12154	bbain	liquibase/r4_6/upgrade.xml	2021-06-01 07:59:26.5024	472	EXECUTED	8:d1f0f6a2d7747e1111a08261c5c5554c	addColumn tableName=sta_normal_user	Add column 'time_zone' to 'sta_normal_user'	\N	3.6.1	\N	\N	2534357482
BSERVDEV-12323-1	jpalacios	liquibase/r4_7/upgrade.xml	2021-06-01 07:59:26.506743	473	EXECUTED	8:dba1089526310b095acb226d058f3ae4	addColumn tableName=sta_pr_diff_comment_anchor	Add a diff_type column to sta_pr_diff_comment_anchor and initialize it to 0 (EFFECTIVE)	\N	3.6.1	\N	\N	2534357482
BSERVDEV-12323-2	jpalacios	liquibase/r4_7/upgrade.xml	2021-06-01 07:59:26.510886	474	EXECUTED	8:79368cad96353ebb88e49db8689b06db	addNotNullConstraint columnName=diff_type, tableName=sta_pr_diff_comment_anchor		\N	3.6.1	\N	\N	2534357482
BSERVDEV-12323-3	jpalacios	liquibase/r4_7/upgrade.xml	2021-06-01 07:59:26.559246	475	EXECUTED	8:27b87f64ecb971dd36d707ba0fedd031	addColumn tableName=sta_cmt_disc_comment_anchor	Add a pr_id column to sta_cmt_disc_comment_anchor and initialize it to null	\N	3.6.1	\N	\N	2534357482
BSERVDEV-12323-4	jpalacios	liquibase/r4_7/upgrade.xml	2021-06-01 07:59:26.566859	476	EXECUTED	8:32929efeb32804760a0f039f0948b933	createIndex indexName=idx_sta_cmt_com_anc_pr, tableName=sta_cmt_disc_comment_anchor	Add a pr_id column index to sta_cmt_disc_comment_anchor	\N	3.6.1	\N	\N	2534357482
BSERVDEV-12323-5	jpalacios	liquibase/r4_7/upgrade.xml	2021-06-01 07:59:26.57207	477	EXECUTED	8:8686bbff91172d153c6af5775a811d9a	addForeignKeyConstraint baseTableName=sta_cmt_disc_comment_anchor, constraintName=fk_sta_pr_com_anc_disc, referencedTableName=sta_pull_request	Add foreign key constraint to the relationship between commit discussion anchor and the pull request	\N	3.6.1	\N	\N	2534357482
BSERVDEV-12993-1	cszmajda	liquibase/r4_8/upgrade.xml	2021-06-01 07:59:26.578003	478	EXECUTED	8:4223011359f4f780201061b7df3c942e	addColumn tableName=sta_drift_request	Add an attempts column to sta_drift_request and initialize it to 0	\N	3.6.1	\N	\N	2534357482
BSERVDEV-13151-1	fhaehnel	liquibase/r4_8/upgrade.xml	2021-06-01 07:59:26.588806	479	EXECUTED	8:041142ba26d00b94d75cf960e6964147	createTable tableName=bb_integrity_event	Create bb_integrity_event table to track events of interest to integrity checking	\N	3.6.1	\N	\N	2534357482
BSERVDEV-13238-1	cszmajda	liquibase/r4_8/upgrade.xml	2021-06-01 07:59:26.598406	480	EXECUTED	8:1d573cbd706c45a9ac123978bf81c93b	addColumn tableName=bb_integrity_event	Add bb_integrity_event.event_node column for cluster safety	\N	3.6.1	\N	\N	2534357482
BSERVDEV-12610	crolf	liquibase/r4_9/upgrade.xml	2021-06-01 07:59:26.607365	481	EXECUTED	8:b0fd31ffaea233ba45566add81d9ac80	createTable tableName=cwd_webhook		\N	3.6.1	\N	\N	2534357482
BSERVDEV-13399-1	spetrucev	liquibase/r4_9/upgrade.xml	2021-06-01 07:59:26.611482	482	EXECUTED	8:f4e5797f225e087fa4bb73b67fa47a90	addColumn tableName=sta_pr_participant	Add last_reviewed_commit column to sta_pr_participant	\N	3.6.1	\N	\N	2534357482
BSERVDEV-13438-1a	bturner	liquibase/r4_9/upgrade.xml	2021-06-01 07:59:26.622262	483	EXECUTED	8:3dd944bceeee845881d26a354fe164f1	createTable tableName=bb_proj_merge_config		\N	3.6.1	\N	\N	2534357482
BSERVDEV-13438-1b	bturner	liquibase/r4_9/upgrade.xml	2021-06-01 07:59:26.628563	484	EXECUTED	8:e2eb5482f087b4431fa852e681085684	addForeignKeyConstraint baseTableName=bb_proj_merge_config, constraintName=fk_bb_proj_merge_config, referencedTableName=project		\N	3.6.1	\N	\N	2534357482
BSERVDEV-13438-1c	bturner	liquibase/r4_9/upgrade.xml	2021-06-01 07:59:26.634223	485	EXECUTED	8:c2d54e7464077e1d80053734f097968e	addUniqueConstraint constraintName=uq_bb_proj_merge_config, tableName=bb_proj_merge_config		\N	3.6.1	\N	\N	2534357482
bturner	BSERVDEV-13438-2a	liquibase/r4_9/upgrade.xml	2021-06-01 07:59:26.641226	486	EXECUTED	8:ceb94899f5d69ad623cdd32465944c1d	createTable tableName=bb_proj_merge_strategy		\N	3.6.1	\N	\N	2534357482
BSERVDEV-13438-2b	bturner	liquibase/r4_9/upgrade.xml	2021-06-01 07:59:26.648572	487	EXECUTED	8:f4439c8646adbdb8dd3a93ad7a6b7fec	addForeignKeyConstraint baseTableName=bb_proj_merge_strategy, constraintName=fk_bb_proj_merge_strategy, referencedTableName=bb_proj_merge_config		\N	3.6.1	\N	\N	2534357482
BSERVDEV-13438-3a	bturner	liquibase/r4_9/upgrade.xml	2021-06-01 07:59:26.6677	488	EXECUTED	8:06ae8095cf6ed43bc64cdc389c826b82	createTable tableName=bb_repo_merge_config		\N	3.6.1	\N	\N	2534357482
BSERVDEV-13438-3b	bturner	liquibase/r4_9/upgrade.xml	2021-06-01 07:59:26.676764	489	EXECUTED	8:615662d6ea32f86e80e1c906f2a2d310	addForeignKeyConstraint baseTableName=bb_repo_merge_config, constraintName=fk_bb_repo_merge_config, referencedTableName=repository		\N	3.6.1	\N	\N	2534357482
bturner	BSERVDEV-13438-4a	liquibase/r4_9/upgrade.xml	2021-06-01 07:59:26.69703	490	EXECUTED	8:55e9a8b832816b93074ea93c44ad4b27	createTable tableName=bb_repo_merge_strategy		\N	3.6.1	\N	\N	2534357482
BSERVDEV-13438-4b	bturner	liquibase/r4_9/upgrade.xml	2021-06-01 07:59:26.703273	491	EXECUTED	8:bec07c93c8b6ad810522411e64d53706	addForeignKeyConstraint baseTableName=bb_repo_merge_strategy, constraintName=fk_bb_repo_merge_strategy, referencedTableName=bb_repo_merge_config		\N	3.6.1	\N	\N	2534357482
BSERVDEV-13438-5	bturner	liquibase/r4_9/upgrade.xml	2021-06-01 07:59:26.737187	492	EXECUTED	8:effb9993d75be7a342906b0391d518ad	createTable tableName=bb_scm_merge_config		\N	3.6.1	\N	\N	2534357482
bturner	BSERVDEV-13438-6a	liquibase/r4_9/upgrade.xml	2021-06-01 07:59:26.75061	493	EXECUTED	8:d0f1a71906af93b731e1b0662e85b4f1	createTable tableName=bb_scm_merge_strategy		\N	3.6.1	\N	\N	2534357482
BSERVDEV-13438-6b	bturner	liquibase/r4_9/upgrade.xml	2021-06-01 07:59:26.759698	494	EXECUTED	8:f7ecf93eb27bcd05825a5d25a5f25788	addForeignKeyConstraint baseTableName=bb_scm_merge_strategy, constraintName=fk_bb_scm_merge_strategy, referencedTableName=bb_scm_merge_config		\N	3.6.1	\N	\N	2534357482
BSERVDEV-13597-1	behumphreys	liquibase/r4_10/upgrade.xml	2021-06-01 07:59:26.768167	495	EXECUTED	8:262274c9dbeb55b168f92f653b000228	addColumn tableName=sta_pull_request	Add closed_timestamp column to sta_pull_request	\N	3.6.1	\N	\N	2534357482
BBSERVDEV-13597-2	behumphreys	liquibase/r4_10/upgrade.xml	2021-06-01 07:59:26.776764	496	EXECUTED	8:3a6ad2463ee0b590e25f4084f3778f5b	createIndex indexName=idx_sta_pr_closed_ts, tableName=sta_pull_request	Add a closed_timestamp column index to sta_pull_request	\N	3.6.1	\N	\N	2534357482
BBSERVDEV-13597-3	behumphreys	liquibase/r4_10/upgrade.xml	2021-06-01 07:59:26.782442	497	EXECUTED	8:06bf49cef8edef47aa9f5c009139da1e	update tableName=sta_pull_request	Populate the closed_timestamp column with updated_timestamp when pull request is in the closed state	\N	3.6.1	\N	\N	2534357482
BSERVDEV-13594-1	behumphreys	liquibase/r4_10/upgrade.xml	2021-06-01 07:59:26.790782	498	EXECUTED	8:293524d483c9ba0d16af049f7db8f041	createIndex indexName=idx_sta_activity_created_time, tableName=sta_activity	Add a created_timestamp column index to sta_activity	\N	3.6.1	\N	\N	2534357482
BSERVDEV-13236-1	behumphreys	liquibase/r4_11/upgrade.xml	2021-06-01 07:59:26.822179	499	EXECUTED	8:22e9b246bcd551cf35230e5dd6549422	sql	Add a text_pattern_ops index to the id column on changeset. This supports queries that attempt\n            to match a commit hash prefix. On PostgreSQL, where a locale other than 'C' is used, such queries\n            require a full table scan.	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-1	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.833666	500	EXECUTED	8:007301ca6d207a740855d5d92c8f609e	createTable tableName=bb_comment_thread	A table for InternalCommentThread instances.\n\n            This table embeds the columns for InternalCommentThreadDiffAnchor and provides no nullability constraints\n            to support pull request general comments.	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-2	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.841384	501	EXECUTED	8:0bfe2b4047cbaa4c217f62bcb3db72de	sql	Create comment threads for all the root comments created on a commit	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-3	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.846959	502	EXECUTED	8:72b568023e0f480f5c1a067be5ddc2a2	sql	Create comment threads for all the root comments created on a pull request diff	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-4	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.852101	503	EXECUTED	8:17ed86371598ffe173a24971474fb576	sql	Create comment threads for pull request general comments	\N	3.6.1	\N	\N	2534357482
BSERV-9918-1	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.856992	504	EXECUTED	8:3515eff276d584460ed2fced2f39bccd	sql	Create comment threads for commit general comments	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-5	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.863295	505	EXECUTED	8:187eb9d33ebce4c8191400caf05859e2	createIndex indexName=idx_bb_com_thr_commentable, tableName=bb_comment_thread		\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-6	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.869483	506	EXECUTED	8:3beddacfdd831bf740f507d374db79a5	createIndex indexName=idx_bb_com_thr_from_hash, tableName=bb_comment_thread	Index from_hash in bb_comment_thread	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-7	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.878877	507	EXECUTED	8:18c7660aa62e2efa47eea8a0706e93af	createIndex indexName=idx_bb_com_thr_to_hash, tableName=bb_comment_thread	Index to_hash in bb_comment_thread	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-8-1	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.887447	508	EXECUTED	8:ad364405106472ed0a1cc363fb9c9960	createIndex indexName=idx_bb_com_thr_to_path, tableName=bb_comment_thread	Index to_path in bb_comment_thread	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-9	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.896463	509	EXECUTED	8:7467cf2264c5a50cc04a98ef7e087f68	createIndex indexName=idx_bb_com_thr_diff_type, tableName=bb_comment_thread	Index the diff_type in bb_comment_thread	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-10	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.904513	510	EXECUTED	8:632b7ed1ffd73a77fd707423c39cbf1e	sql	Update the bb_comment_thread sequence id generator to avoid collisions	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-11	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.91689	511	EXECUTED	8:21e13ccd067af4eaccb2ed33473183e9	createTable tableName=bb_comment	Create bb_comment table to replace sta_comment	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-12	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.931143	512	EXECUTED	8:550dbd42a9a399fe39437b07e74c463e	sql	Copy data from sta_comment to bb_comment for root comments	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-13	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.9435	513	EXECUTED	8:de1240ff63dec5ab231f138d1f690626	createTable tableName=bb_thread_root_comment	Create "bb_thread_root_comment" table to manage the commentThread -> rootComment -> commentThread relationship.\n\n            Note that no DELETE cascade is possible on the comment_id FK. SQL Server detects a potential fk loop and\n            stops...	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-14	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.95429	514	EXECUTED	8:f86d27cb00b3f7ffb66189b281606490	sql	Populate bb_thread_root_comment	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-15	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.959291	515	EXECUTED	8:93628dab4b10dbfe09584e15de3e2033	sql	Copy data from sta_comment to bb_comment for replies	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-16	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.968019	516	EXECUTED	8:f8c9c068cfed53a7f83501fe5f7f9841	createIndex indexName=idx_bb_comment_author, tableName=bb_comment	Create index on bb_comment by author_id	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-17	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.977718	517	EXECUTED	8:f90b1b0f402c46a80b3ab4524eb16ec8	createIndex indexName=idx_bb_comment_thread, tableName=bb_comment	Create index on bb_comment by thread_id	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-18	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.987545	518	EXECUTED	8:79f9bc291e6fcf77c306e0273f915fde	sql	Update updated_timestamp in bb_comment_thread to latest comment in thread	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-19	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.99242	519	EXECUTED	8:89285d07a298d5196f6cf0551753e33a	addNotNullConstraint columnName=updated_timestamp, tableName=bb_comment_thread	Add not null constraint to bb_comment_thread.updated_timestamp	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-20	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:26.998996	520	EXECUTED	8:d273bbe2c049fc8bf3d448826f523042	sql	Update the bb_comment sequence id generator to avoid collisions	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-21	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.007126	521	EXECUTED	8:f97755ca20292f56e417f0ec3fd7cef2	createTable tableName=bb_pr_comment_activity	Create bb_pr_comment_activity table to replace sta_pr_comment_activity.\n            Essentially the same table but without the anchor_id since the anchor can be reached from the comment	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-22	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.012626	522	EXECUTED	8:647e34c09d9971a24c24c6674c0bbf78	addForeignKeyConstraint baseTableName=bb_pr_comment_activity, constraintName=fk_bb_pr_com_act_id, referencedTableName=sta_pr_activity		\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-23	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.016908	523	EXECUTED	8:66116f76919ad82fbe42f705549d0c2f	addForeignKeyConstraint baseTableName=bb_pr_comment_activity, constraintName=fk_bb_pr_com_act_comment, referencedTableName=bb_comment	Create the fk_bb_pr_com_act_comment foreign key	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-24	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.021111	524	EXECUTED	8:e63dd7d548b1912bb44d9d357333ffa5	sql	Populate bb_pr_comment_activity	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-25	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.02686	525	EXECUTED	8:3b36cb81e056bc6f015619ed6d145f17	createIndex indexName=idx_bb_pr_com_act_comment, tableName=bb_pr_comment_activity		\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-26	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.035837	526	EXECUTED	8:29decfb8fe8fa514a80f6ac0bbd25151	createTable tableName=bb_cmt_disc_comment_activity	Create the bb_cmt_disc_comment_activity to replace sta_cmt_disc_comment_activity\n            Essentially the same table but without the anchor_id since the anchor can be reached from the comment	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-27	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.040932	527	EXECUTED	8:c0e7aa72cbbc395ca4d3a2d1a83b323f	addForeignKeyConstraint baseTableName=bb_cmt_disc_comment_activity, constraintName=fk_bb_cmt_disc_com_act_id, referencedTableName=sta_cmt_disc_activity	Create a foreign key between comment activities and their base discussion activities, cascading deletion\n            to simplify deleting activities.	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-28	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.046561	528	EXECUTED	8:208331d500724f227c56d3ef3abe6eb8	addForeignKeyConstraint baseTableName=bb_cmt_disc_comment_activity, constraintName=fk_bb_cmt_disc_com_act_com, referencedTableName=bb_comment	Create a foreign key between comment activities and their comments. Note that this foreign key does not\n            cascade deletions because doing so would leave orphaned rows in other activity tables.	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-29	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.05154	529	EXECUTED	8:3163953eeb90ae5bcb405e5700cec031	sql	Populate bb_cmt_disc_comment_activity	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-30	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.059344	530	EXECUTED	8:76868468fa274114c1d3b39db6b4f78c	createIndex indexName=idx_bb_cmt_disc_com_act_com, tableName=bb_cmt_disc_comment_activity	Create an index on comment IDs to facilitate applying the foreign key to bb_comment.	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-31	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.073236	531	EXECUTED	8:de4d539421d5d1ec10cf346c7d749098	createTable tableName=bb_comment_parent	Create "bb_comment_parent" table to manage the comment -> parent relationship.\n\n            Note that no DELETE cascade is possible on the parent_id FK. SQL Server detects a potential fk loop and\n            stops us from doing it. Deletes need to...	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-32	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.076848	532	EXECUTED	8:f2d5c9d14514b2213c8514d6f0b4d3cc	sql	Populate bb_comment_parent	\N	3.6.1	\N	\N	2534357482
BSERVDEV-8489-33	jpalacios	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.087392	533	EXECUTED	8:7703dda781f2a888a6f4b92a1b5d1037	createIndex indexName=idx_bb_com_par_parent, tableName=bb_comment_parent	Create index on bb_comment_parent.parent_id	\N	3.6.1	\N	\N	2534357482
BSERV-3751-1	mheemskerk	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.095115	534	EXECUTED	8:b7feed1f9d8098a4c7219e67ba7b549d	createTable tableName=bb_repository_alias	Add bb_repository_alias table to keep track of renamed/moved repositories	\N	3.6.1	\N	\N	2534357482
BSERV-3751-2	mheemskerk	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.101799	535	EXECUTED	8:e6fe03fb6021d20d325c2b646378bec2	addForeignKeyConstraint baseTableName=bb_repository_alias, constraintName=fk_repository_alias_repo, referencedTableName=repository		\N	3.6.1	\N	\N	2534357482
BSERV-3751-3	mheemskerk	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.109819	536	EXECUTED	8:f13a6e2b84ffb6b1966ee6894bff5206	addUniqueConstraint constraintName=uq_bb_repo_alias_key_slug, tableName=bb_repository_alias	Create unique constraint on bb_repository_alias.[project_namespace,project_key,slug]	\N	3.6.1	\N	\N	2534357482
BSERV-3751-4	mheemskerk	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.11668	537	EXECUTED	8:c4ad65b38f24af68875926e44993ca00	createTable tableName=bb_project_alias	Add bb_project_alias table to keep track of renamed projects	\N	3.6.1	\N	\N	2534357482
BSERV-3751-5	mheemskerk	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.122971	538	EXECUTED	8:9dd82cb327ba4d281e427b4935e1c508	addForeignKeyConstraint baseTableName=bb_project_alias, constraintName=fk_project_alias_project, referencedTableName=project		\N	3.6.1	\N	\N	2534357482
BSERV-3751-6	mheemskerk	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.137084	539	EXECUTED	8:5a1c748c574346a3e8214ef9fdd35a05	addUniqueConstraint constraintName=uq_bb_project_alias_ns_key, tableName=bb_project_alias	Create unique constraint on bb_project_alias.[project_namespace,project_key]	\N	3.6.1	\N	\N	2534357482
BSERVDEV-14515-1	bturner	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.143367	540	EXECUTED	8:958ab29f96865624238ea1f82c28a4df	createTable tableName=cwd_granted_perm	Create cwd_granted_perm table for Crowd's new UserPermission type.	\N	3.6.1	\N	\N	2534357482
BSERVDEV-14515-2	bturner	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.148933	541	EXECUTED	8:5e62646b33f3aac394d8b75b3294da1a	addColumn tableName=cwd_group	Add external_id column to cwd_group	\N	3.6.1	\N	\N	2534357482
BSERVDEV-14515-3	bturner	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.15523	542	EXECUTED	8:e73cdb4a10afdbb346ea5300197c208f	createIndex indexName=idx_cwd_group_external_id, tableName=cwd_group	Index cwd_group.external_id for queries	\N	3.6.1	\N	\N	2534357482
BSERVDEV-14515-4	bturner	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.159232	543	EXECUTED	8:c082cbb8fca645c5ffec757156bfec6d	addColumn tableName=cwd_membership		\N	3.6.1	\N	\N	2534357482
BSERVDEV-14515-5	bturner	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.172646	544	EXECUTED	8:eea005faf181ae2eacaacae0cc9130ed	createTable tableName=cwd_tombstone		\N	3.6.1	\N	\N	2534357482
BSERVDEV-14515-6	bturner	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.181138	545	EXECUTED	8:ffd3e48d65927976941e8e46934bc4d0	createIndex indexName=idx_tombstone_type_timestamp, tableName=cwd_tombstone		\N	3.6.1	\N	\N	2534357482
BSERVDEV-13541-1	bturner	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.188625	546	EXECUTED	8:958d2f4c580b5c8fc27f2d41a517d916	modifyDataType columnName=next_hi, tableName=hibernate_unique_key	Convert hibernate_unique_key.next_hi from int to bigint, as required by Hibernate 5.2's ID generators	\N	3.6.1	\N	\N	2534357482
BSERVDEV-13541-2	bturner	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.194904	547	EXECUTED	8:b553fe198974a1a3625f22e3ca6badae	addNotNullConstraint columnName=next_hi, tableName=hibernate_unique_key	Add a NOT NULL constraint to hibernate_unique_key.next_hi and, on MySQL and SQL Server, change its\n            data type from int to bigint as required by Hibernate 5.2's ID generators	\N	3.6.1	\N	\N	2534357482
BSERVDEV-13541-3	bturner	liquibase/r5_0/upgrade.xml	2021-06-01 07:59:27.212413	548	EXECUTED	8:8e3c806c068a131950e97573efcf0b9f	modifyDataType columnName=next_val, tableName=id_sequence	Convert id_sequence.next_val from int to bigint, as required by Hibernate 5.2's ID generators	\N	3.6.1	\N	\N	2534357482
BSERVDEV-15338-1	fhaehnel	liquibase/r5_0/p02.xml	2021-06-01 07:59:27.228569	549	EXECUTED	8:76df70a7df374ccc71e2878d7e6a302a	dropForeignKeyConstraint baseTableName=sta_pr_diff_comment_anchor, constraintName=fk_sta_pr_diff_com_anc_pr	Dropping FK constraints on deprecated comment tables to allow deleting pull requests	\N	3.6.1	\N	\N	2534357482
BSERVDEV-15338-2	fhaehnel	liquibase/r5_0/p02.xml	2021-06-01 07:59:27.233759	550	EXECUTED	8:5f09e1aa50680b1f7d7959f9a3aa1c12	dropForeignKeyConstraint baseTableName=sta_cmt_disc_comment_anchor, constraintName=fk_sta_pr_com_anc_disc		\N	3.6.1	\N	\N	2534357482
BSERVDEV-15569-1	fhaehnel	liquibase/r5_0/p05.xml	2021-06-01 07:59:27.238536	551	EXECUTED	8:33a722171a978211f9802c3b37fa839e	dropForeignKeyConstraint baseTableName=sta_cmt_disc_comment_anchor, constraintName=fk_sta_cmt_disc_com_anc_disc		\N	3.6.1	\N	\N	2534357482
BSERVDEV-15569-2	fhaehnel	liquibase/r5_0/p05.xml	2021-06-01 07:59:27.244501	552	EXECUTED	8:589b66bb33310dd724a2ea6c71543ff9	sql	De-duplicate pull request commit-level review comment activities.\n            During the 5.0 migration, we failed to differentiate between Commit comment activities (with no pull request)\n            and Pull request commit comment activities; we ...	\N	3.6.1	\N	\N	2534357482
BSERVDEV-15994-1	spetrucev	liquibase/r5_0/p08.xml	2021-06-01 07:59:27.249113	553	EXECUTED	8:e50ef3151843e3d229a6e2b04e3e0227	dropForeignKeyConstraint baseTableName=sta_comment, constraintName=fk_sta_comment_parent	Dropping FK constraints on deprecated comment tables to prevent non-monotonic comment rows from\n                 interfering with database restores and migrations (see BSERVDEV-8452 also)	\N	3.6.1	\N	\N	2534357482
BSERVDEV-15994-2	spetrucev	liquibase/r5_0/p08.xml	2021-06-01 07:59:27.258375	554	EXECUTED	8:7ce5857a79d122c213422a6f15610ab2	dropForeignKeyConstraint baseTableName=sta_comment, constraintName=fk_sta_comment_root	Dropping FK constraints on deprecated comment tables to prevent non-monotonic comment rows from\n                 interfering with database restores and migrations (see BSERVDEV-8452 also)	\N	3.6.1	\N	\N	2534357482
BSERVDEV-14855-1	jthomas	liquibase/r5_2/upgrade.xml	2021-06-01 07:59:27.269609	555	EXECUTED	8:db2fbcde030f15f1150e031982790cc2	addColumn tableName=sta_repo_hook; addForeignKeyConstraint baseTableName=sta_repo_hook, constraintName=fk_sta_repo_hook_proj, referencedTableName=project	Add project column to repository hooks	\N	3.6.1	\N	\N	2534357482
BSERVDEV-14855-2	jthomas	liquibase/r5_2/upgrade.xml	2021-06-01 07:59:27.275256	556	EXECUTED	8:51421c6cc78f92fb381d988800b79850	dropNotNullConstraint columnName=repository_id, tableName=sta_repo_hook	Drop not null constraint on repository_id	\N	3.6.1	\N	\N	2534357482
BSERVDEV-14855-3	jthomas	liquibase/r5_2/upgrade.xml	2021-06-01 07:59:27.281993	557	EXECUTED	8:48b943a434c3c99520e389383148f0b6	createIndex indexName=idx_sta_repo_hook_proj_id, tableName=sta_repo_hook	Create a index on project_id	\N	3.6.1	\N	\N	2534357482
BSERVDEV-14855-4	jthomas	liquibase/r5_2/upgrade.xml	2021-06-01 07:59:27.286726	558	EXECUTED	8:c888339274c6a84ddbc5b01fdc3f7158	dropUniqueConstraint constraintName=uq_sta_repo_hook_repo_hook_key, tableName=sta_repo_hook	Drop existing unique constraint on repository_id	\N	3.6.1	\N	\N	2534357482
BSERVDEV-14855-5	jthomas	liquibase/r5_2/upgrade.xml	2021-06-01 07:59:27.295007	559	EXECUTED	8:632b08111ffe884e17202d76b0aa9eec	addUniqueConstraint constraintName=uq_sta_repo_hook_scope_hook, tableName=sta_repo_hook	Add unique constraint on project_id, repository_id and hook_key	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16024-1	spetrucev	liquibase/r5_4/upgrade.xml	2021-06-01 07:59:27.302856	560	EXECUTED	8:dedb457e79f43d1e32a195018d95fb2e	dropTable tableName=sta_cmt_disc_comment_activity	Drop the deprecated sta_cmt_disc_comment_activity table, which was replaced by bb_cmt_disc_comment_activity\n            in 5.0 (BSERVDEV-8489)	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16024-2	spetrucev	liquibase/r5_4/upgrade.xml	2021-06-01 07:59:27.309309	561	EXECUTED	8:db63438cc2255cedeec6cf9d1ce6d3dc	dropTable tableName=sta_pr_comment_activity	Drop the deprecated sta_pr_comment_activity table, which was replaced by bb_pr_comment_activity in 5.0\n            (BSERVDEV-8489)	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16024-3	spetrucev	liquibase/r5_4/upgrade.xml	2021-06-01 07:59:27.317469	562	EXECUTED	8:e794adea92533dde1b2de20e8799beaa	dropTable tableName=sta_cmt_disc_comment_anchor	Drop the deprecated sta_cmt_disc_comment_anchor table, which was replaced by bb_comment_thread in 5.0\n            (BSERVDEV-8489)	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16024-4	spetrucev	liquibase/r5_4/upgrade.xml	2021-06-01 07:59:27.326458	563	EXECUTED	8:645bde11fcbc897e0b114cd7679b1f10	dropTable tableName=sta_pr_diff_comment_anchor	Drop the deprecated sta_pr_diff_comment_anchor table, which was replaced by bb_comment_thread in 5.0\n            (BSERVDEV-8489)	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16024-5	spetrucev	liquibase/r5_4/upgrade.xml	2021-06-01 07:59:27.333621	564	EXECUTED	8:bfe3fa4d1e0342f98511d0e0263cae1e	dropTable tableName=sta_diff_comment_anchor	Drop the deprecated sta_diff_comment_anchor table, which was replaced by bb_comment_thread in 5.0\n            (BSERVDEV-8489)	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16024-6	spetrucev	liquibase/r5_4/upgrade.xml	2021-06-01 07:59:27.341698	565	EXECUTED	8:f82ad470e52f1bc95a0ed5c94fabaf5a	dropTable tableName=sta_comment	Drop the deprecated sta_comment table, which was replaced by bb_comment in 5.0 (BSERVDEV-8489)	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16024-7	spetrucev	liquibase/r5_4/upgrade.xml	2021-06-01 07:59:27.35355	566	EXECUTED	8:cf6df07155b7008a1ec9bd6936c432f8	delete tableName=id_sequence; delete tableName=id_sequence	Drop ID sequences for legacy comment tables	\N	3.6.1	\N	\N	2534357482
BSERV-10063-1	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.360745	567	EXECUTED	8:5df18fd9691087270530c4c74eaaf22f	createIndex indexName=idx_pr_reviewer_added_user_id, tableName=bb_pr_reviewer_added		\N	3.6.1	\N	\N	2534357482
BSERV-10063-2	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.367729	568	EXECUTED	8:52c4098297a2027ffbd013090648e425	createIndex indexName=idx_cwd_user_user_id, tableName=cwd_user_credential_record		\N	3.6.1	\N	\N	2534357482
BSERV-10063-3	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.375672	569	EXECUTED	8:1fd0dfc5b78ea909792adcedae010b58	createIndex indexName=idx_cwd_webhook_application_id, tableName=cwd_webhook		\N	3.6.1	\N	\N	2534357482
BSERV-10063-4	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.381882	570	EXECUTED	8:25d0c0014c825fabda9d5d567217cd01	createIndex indexName=idx_sta_user_settings_lob_id, tableName=sta_user_settings		\N	3.6.1	\N	\N	2534357482
BSERV-10063-5	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.387948	571	EXECUTED	8:292ae725aa9f18bf1dd6812ac4135218	createIndex indexName=idx_granted_perm_group_mapping, tableName=cwd_granted_perm		\N	3.6.1	\N	\N	2534357482
BSERV-10063-6	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.412993	572	EXECUTED	8:1c85f0cd452953f7b0f3d1512d9884ae	createIndex indexName=idx_cwd_group_directory_id, tableName=cwd_group		\N	3.6.1	\N	\N	2534357482
BSERV-10063-7	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.419796	573	EXECUTED	8:f8fe664394ca76405a641bc0e5a93ba6	createIndex indexName=idx_pr_review_removed_user_id, tableName=bb_pr_reviewer_removed		\N	3.6.1	\N	\N	2534357482
BSERV-10063-8	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.429719	574	EXECUTED	8:76463daeb41d488ee9f49e45ece50e86	createIndex indexName=idx_app_address_app_id, tableName=cwd_application_address		\N	3.6.1	\N	\N	2534357482
BSERV-10063-9	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.435406	575	EXECUTED	8:de5d709810b57a8c8379413a7c368e30	createIndex indexName=idx_project_permission_perm_id, tableName=sta_project_permission		\N	3.6.1	\N	\N	2534357482
BSERV-10063-10	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.440994	576	EXECUTED	8:df08474697894ceb51cf1ec64df7d9e4	createIndex indexName=idx_rep_alias_repo_id, tableName=bb_repository_alias		\N	3.6.1	\N	\N	2534357482
BSERV-10063-11	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.447071	577	EXECUTED	8:67f658667f6f26516af4072597c6571e	createIndex indexName=idx_app_dir_grp_mapping_app_id, tableName=cwd_app_dir_group_mapping		\N	3.6.1	\N	\N	2534357482
BSERV-10063-12	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.453504	578	EXECUTED	8:d63b9e01c1dcffec25758246a7c99c83	createIndex indexName=idx_sta_global_perm_perm_id, tableName=sta_global_permission		\N	3.6.1	\N	\N	2534357482
BSERV-10063-13	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.460791	579	EXECUTED	8:fed7a44a84f9967ef760e1f385b558f6	createIndex indexName=idx_sta_repo_perm_repo_id, tableName=sta_repo_permission		\N	3.6.1	\N	\N	2534357482
BSERV-10063-14	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.471463	580	EXECUTED	8:5d213e2f52914c06c6892f8406d2ea53	createIndex indexName=idx_sta_drift_request_pr_id, tableName=sta_drift_request		\N	3.6.1	\N	\N	2534357482
BSERV-10063-15	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.477924	581	EXECUTED	8:d92085ce06f3532cea4ff6c6856da56a	createIndex indexName=idx_pr_rescope_request_pr_id, tableName=sta_pr_rescope_request		\N	3.6.1	\N	\N	2534357482
BSERV-10063-16	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.491042	582	EXECUTED	8:c98e306fe5cf0f6110a69a6bb50475ea	createIndex indexName=idx_cwd_app_dir_mapping_dir_id, tableName=cwd_app_dir_mapping		\N	3.6.1	\N	\N	2534357482
BSERV-10063-17	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.497664	583	EXECUTED	8:22b78ece06b4aa95322c31f7fc930e7a	createIndex indexName=idx_cwd_membership_dir_id, tableName=cwd_membership		\N	3.6.1	\N	\N	2534357482
BSERV-10063-18	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.502998	584	EXECUTED	8:cf6b45b3b59dc8969bae10565bb95357	createIndex indexName=idx_repo_access_repo_id, tableName=repository_access		\N	3.6.1	\N	\N	2534357482
BSERV-10063-19	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.508971	585	EXECUTED	8:d8f36d3f1f55c4cdac189a7a4c2c2e2b	createIndex indexName=idx_sta_watcher_user_id, tableName=sta_watcher		\N	3.6.1	\N	\N	2534357482
BSERV-10063-20	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.516819	586	EXECUTED	8:c64bf84485a3d5d97a13dd1031a5d5f0	createIndex indexName=idx_bb_proj_alias_proj_id, tableName=bb_project_alias		\N	3.6.1	\N	\N	2534357482
BSERV-10063-21	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.52651	587	EXECUTED	8:488db6466cee2ada1c97139a7a17672e	createIndex indexName=idx_sta_proj_perm_pro_id, tableName=sta_project_permission		\N	3.6.1	\N	\N	2534357482
BSERV-10063-22	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.534499	588	EXECUTED	8:db59c7d122bd3efb711ce807fd4aea6e	createIndex indexName=idx_sta_repo_perm_perm_id, tableName=sta_repo_permission		\N	3.6.1	\N	\N	2534357482
BSERV-10063-23	dkjellin	liquibase/r5_5/upgrade.xml	2021-06-01 07:59:27.540728	589	EXECUTED	8:289e2be6367ce225521227755c630559	createIndex indexName=idx_remember_me_token_user_id, tableName=sta_remember_me_token		\N	3.6.1	\N	\N	2534357482
BSERVDEV-16158-1	bturner	liquibase/r5_6/upgrade.xml	2021-06-01 07:59:27.546162	590	EXECUTED	8:2805c79f2946d16cc01365e89850c0df	addColumn tableName=plugin_state		\N	3.6.1	\N	\N	2534357482
BSERVDEV-16158-2	bturner	liquibase/r5_6/upgrade.xml	2021-06-01 07:59:27.550233	591	EXECUTED	8:815ebf4635f5690780b36a3918640426	update tableName=plugin_state		\N	3.6.1	\N	\N	2534357482
BSERVDEV-16158-3	bturner	liquibase/r5_6/upgrade.xml	2021-06-01 07:59:27.554052	592	EXECUTED	8:5d236a31e2e49fc5839f5095a7dad1f8	addNotNullConstraint columnName=updated_timestamp, tableName=plugin_state		\N	3.6.1	\N	\N	2534357482
BSERVDEV-16266-1	mstudman	liquibase/r5_9/upgrade.xml	2021-06-01 07:59:27.563731	593	EXECUTED	8:42b4955ac561651491725484c628740e	createTable tableName=bb_alert	A table for InternalAlert instances.	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16266-2	mstudman	liquibase/r5_9/upgrade.xml	2021-06-01 07:59:27.571798	594	EXECUTED	8:44517501762ff03e129b8bc76c5b65fe	createIndex indexName=bb_alert_timestamp, tableName=bb_alert		\N	3.6.1	\N	\N	2534357482
BSERVDEV-16266-3	mstudman	liquibase/r5_9/upgrade.xml	2021-06-01 07:59:27.584684	595	EXECUTED	8:c0d927cbebaf86981748533f0118eb2f	createIndex indexName=bb_alert_issue, tableName=bb_alert		\N	3.6.1	\N	\N	2534357482
BSERVDEV-16266-4	mstudman	liquibase/r5_9/upgrade.xml	2021-06-01 07:59:27.594414	596	EXECUTED	8:12dbcd2718fc9bd0e68343c02d750ef9	createIndex indexName=bb_alert_issue_component, tableName=bb_alert		\N	3.6.1	\N	\N	2534357482
BSERVDEV-16266-5	mstudman	liquibase/r5_9/upgrade.xml	2021-06-01 07:59:27.602609	597	EXECUTED	8:14d79c5ce7918c6d536d415ffe05bdc9	createIndex indexName=bb_alert_node_lower, tableName=bb_alert		\N	3.6.1	\N	\N	2534357482
BSERVDEV-16266-6	mstudman	liquibase/r5_9/upgrade.xml	2021-06-01 07:59:27.610765	598	EXECUTED	8:0a7f0430de9574348a0d2d9261da87b8	createIndex indexName=bb_alert_plugin_lower, tableName=bb_alert		\N	3.6.1	\N	\N	2534357482
BSERVDEV-16266-7	mstudman	liquibase/r5_9/upgrade.xml	2021-06-01 07:59:27.617543	599	EXECUTED	8:8e971209270bf8063576e7564e38323d	createIndex indexName=bb_alert_severity, tableName=bb_alert		\N	3.6.1	\N	\N	2534357482
BSERV-2642-1	bturner	liquibase/r5_10/upgrade.xml	2021-06-01 07:59:27.625766	600	EXECUTED	8:b6de3c1469eb4c325797624055f13b55	addColumn tableName=sta_repo_push_activity		\N	3.6.1	\N	\N	2534357482
BSERV-2642-2	bturner	liquibase/r5_10/upgrade.xml	2021-06-01 07:59:27.629817	601	EXECUTED	8:a861741d0d6c93b917d248bc684f01a5	update tableName=sta_repo_push_activity		\N	3.6.1	\N	\N	2534357482
BSERV-2642-3	bturner	liquibase/r5_10/upgrade.xml	2021-06-01 07:59:27.633597	602	EXECUTED	8:0a7246ad4bd3f03ad2a17c3de11d9841	addNotNullConstraint columnName=trigger_id, tableName=sta_repo_push_activity		\N	3.6.1	\N	\N	2534357482
BBSDEV-16221-1	bplump	liquibase/r5_11/upgrade.xml	2021-06-01 07:59:27.640932	603	EXECUTED	8:1ce5bc986583f9d249a3dd64d92a12e6	createTable tableName=bb_pr_commit		\N	3.6.1	\N	\N	2534357482
BBSDEV-16221-2	bplump	liquibase/r5_11/upgrade.xml	2021-06-01 07:59:27.647847	604	EXECUTED	8:2595f5516e46c310ac6452a843de57a3	createIndex indexName=idx_bb_pr_commit_commit_id, tableName=bb_pr_commit		\N	3.6.1	\N	\N	2534357482
BBSDEV-17719-1	tkenis	liquibase/r5_12/upgrade.xml	2021-06-01 07:59:27.654565	605	EXECUTED	8:d3b8689b39ca273e761181ee73df375c	createTable tableName=bb_label		\N	3.6.1	\N	\N	2534357482
BBSDEV-17719-2	tkenis	liquibase/r5_12/upgrade.xml	2021-06-01 07:59:27.66077	606	EXECUTED	8:0616e3bf5ef7502ee3b94a255e41e875	createTable tableName=bb_label_mapping		\N	3.6.1	\N	\N	2534357482
BBSDEV-17719-3	tkenis	liquibase/r5_12/upgrade.xml	2021-06-01 07:59:27.667973	607	EXECUTED	8:1b4199b4adc00a7069961a77461d0353	createIndex indexName=idx_bb_label_mapping_label_id, tableName=bb_label_mapping		\N	3.6.1	\N	\N	2534357482
BBSDEV-17719-4	tkenis	liquibase/r5_12/upgrade.xml	2021-06-01 07:59:27.673746	608	EXECUTED	8:903c228dbfc2b1312adcc1d1a36beb9c	createIndex indexName=idx_bb_label_map_labelable_id, tableName=bb_label_mapping		\N	3.6.1	\N	\N	2534357482
BBSDEV-17719-5	tkenis	liquibase/r5_12/upgrade.xml	2021-06-01 07:59:27.679669	609	EXECUTED	8:c00f827d5b419bb09ece5da47cffc8ce	sql		\N	3.6.1	\N	\N	2534357482
BSERVDEV-16285-1	mstudman	liquibase/r5_13/upgrade.xml	2021-06-01 07:59:27.691768	610	EXECUTED	8:b8f904b398cb361d0df30cb5fd70043a	createTable tableName=bb_job	Create the bb_job table for tracking state and progress of long running jobs in different features	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16285-2	mstudman	liquibase/r5_13/upgrade.xml	2021-06-01 07:59:27.698015	611	EXECUTED	8:792273fe349689ed484a9900093114e6	createIndex indexName=idx_bb_job_type, tableName=bb_job	Create the idx_bb_job_type index on the bb_job table	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16285-3	mstudman	liquibase/r5_13/upgrade.xml	2021-06-01 07:59:27.70479	612	EXECUTED	8:7fcbc5647f3c6156abbf030e78428ebc	createIndex indexName=idx_bb_job_state_type, tableName=bb_job	Create the idx_bb_job_state_type index on the bb_job table	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16285-4	mstudman	liquibase/r5_13/upgrade.xml	2021-06-01 07:59:27.70972	613	EXECUTED	8:21f5132d82593d220f6e0912b284d620	addForeignKeyConstraint baseTableName=bb_job, constraintName=fk_bb_job_initiator, referencedTableName=stash_user	Create the foreign key constraint fk_bb_job_initiator for bb_job	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16285-5	mstudman	liquibase/r5_13/upgrade.xml	2021-06-01 07:59:27.725373	614	EXECUTED	8:20d368eae840884a755c6cba223f731f	createTable tableName=bb_job_message	Create the bb_job_message table for tracking messages logged by a Job while executing, associated with a\n            scope - an entity such as a repository or project or the global scope	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16285-6	mstudman	liquibase/r5_13/upgrade.xml	2021-06-01 07:59:27.732979	615	EXECUTED	8:94f428843f9715ab6019ca39243af94d	addForeignKeyConstraint baseTableName=bb_job_message, constraintName=fk_bb_job_msg_job, referencedTableName=bb_job	Create the foreign key constraints for bb_job_message	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16285-7	mstudman	liquibase/r5_13/upgrade.xml	2021-06-01 07:59:27.743366	616	EXECUTED	8:28109080d1c7c22a072b728713a7497f	createIndex indexName=idx_bb_job_msg_job_severity, tableName=bb_job_message	Create the indexes on the bb_job_message table	\N	3.6.1	\N	\N	2534357482
BSERVDEV-16285-8	mstudman	liquibase/r5_13/upgrade.xml	2021-06-01 07:59:27.751248	617	EXECUTED	8:0984e7c09020779441861c164538839c	createIndex indexName=idx_bb_job_stash_user, tableName=bb_job	Create the indexes on the bb_job_message table	\N	3.6.1	\N	\N	2534357482
BBSDEV-17340-2	bturner	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.761459	618	EXECUTED	8:2155441f8f85cb619613ec61b606b57e	addColumn tableName=repository		\N	3.6.1	\N	\N	2534357482
BBSDEV-17340-3	bturner	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.767387	619	EXECUTED	8:d9349db917818e71d83b7414b7df42b6	createIndex indexName=idx_repository_store_id, tableName=repository		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-1	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.77429	620	EXECUTED	8:bcae45981822009fa54c18e48213a18f	createTable tableName=cwd_app_dir_default_groups		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-2	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.780223	621	EXECUTED	8:6eda013ecfb77d075bbf41a518584ef3	addUniqueConstraint constraintName=uk_appmapping_groupname, tableName=cwd_app_dir_default_groups		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-3	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.785452	622	EXECUTED	8:50d57ea2ffea3852ac7a145bb2008792	addForeignKeyConstraint baseTableName=cwd_app_dir_default_groups, constraintName=fk_app_mapping, referencedTableName=cwd_app_dir_mapping		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-4	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.79005	623	EXECUTED	8:0bf508fd82b03d785b7ec2fd37697bf7	dropIndex indexName=idx_granted_perm_group_mapping, tableName=cwd_granted_perm		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-5	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.795713	624	EXECUTED	8:1abc73b75423360c039281a91d9bed7b	addColumn tableName=cwd_granted_perm		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-6	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.800587	625	EXECUTED	8:1506a9722294a710734c990d3066db52	delete tableName=cwd_granted_perm		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-7	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.804783	626	EXECUTED	8:e67ce1a4250503bbf868950bfadd5a98	addNotNullConstraint columnName=group_name, tableName=cwd_granted_perm		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-8	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.808795	627	EXECUTED	8:176016c5376797af87e571a50f685cc7	addNotNullConstraint columnName=app_dir_mapping_id, tableName=cwd_granted_perm		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-9	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.815149	628	EXECUTED	8:21c474dd4aeb0ba07763c5a30946b73a	addForeignKeyConstraint baseTableName=cwd_granted_perm, constraintName=fk_granted_perm_dir_mapping, referencedTableName=cwd_app_dir_mapping		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-10	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.820364	629	EXECUTED	8:29486aed4bfeb969adc34d3e8443b163	dropForeignKeyConstraint baseTableName=cwd_granted_perm, constraintName=fk_cwd_granted_perm_grp_map		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-11	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.824801	630	EXECUTED	8:c98711e5a3ad54fa09af9e0c2552fdbc	dropColumn columnName=group_mapping, tableName=cwd_granted_perm		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-12	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.829117	631	EXECUTED	8:ae343d35a2302ad495fc1aae0967d4ff	addColumn tableName=cwd_user_attribute		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-13	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.835827	632	EXECUTED	8:d6142e3ac8e9f84f2bdd0a48bd589159	createIndex indexName=idx_user_attr_nval, tableName=cwd_user_attribute		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-14	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.839508	633	EXECUTED	8:59956d2f087ecc4ba93b07d57f15212d	delete tableName=cwd_directory_attribute		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-15	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.845623	634	EXECUTED	8:227347e940dbddce12a0f0d21eb35ea5	addColumn tableName=cwd_directory_attribute		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-16	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.849197	635	EXECUTED	8:694062d7248b19e8097f1efb82d426c0	update tableName=cwd_directory_attribute		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-17	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.85408	636	EXECUTED	8:86528d11cc2ce5dfa35d747c64854d52	dropColumn columnName=attribute_value, tableName=cwd_directory_attribute		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-18	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.859711	637	EXECUTED	8:82b5e09a1f827d36c27294e7b5fb24a3	renameColumn newColumnName=attribute_value, oldColumnName=attribute_value_clob, tableName=cwd_directory_attribute		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-19	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.869052	638	EXECUTED	8:e24c260d1e97a9152ce9de162bea9c20	createTable tableName=cwd_group_admin_group		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-20	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.875872	639	EXECUTED	8:6d6b6879bf173a11e7b57696cc09cc0c	addForeignKeyConstraint baseTableName=cwd_group_admin_group, constraintName=fk_admin_group, referencedTableName=cwd_group		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-21	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.881379	640	EXECUTED	8:552448a32c88c2ab047fe32ff86fe2a0	addForeignKeyConstraint baseTableName=cwd_group_admin_group, constraintName=fk_group_target_group, referencedTableName=cwd_group		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-22	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.893786	641	EXECUTED	8:4143a3d5aae82e4b71f7a581df3dd152	createIndex indexName=idx_admin_group, tableName=cwd_group_admin_group		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-23	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.902384	642	EXECUTED	8:ce9686830ce40f5c7db59cc06dd9a7a5	createIndex indexName=idx_group_target_group, tableName=cwd_group_admin_group		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-24	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.908124	643	EXECUTED	8:9782e8522ac6f83947a27a45b2aeeb89	addUniqueConstraint constraintName=uk_group_and_target_group, tableName=cwd_group_admin_group		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-25	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.914911	644	EXECUTED	8:8c9dd52cc981c60836d8393d711c0dbf	createTable tableName=cwd_group_admin_user		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-26	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.920193	645	EXECUTED	8:2052ebb980fafac34c2a13c2f91897cb	addForeignKeyConstraint baseTableName=cwd_group_admin_user, constraintName=fk_admin_user, referencedTableName=cwd_user		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-27	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.924806	646	EXECUTED	8:eaca499737a2d5a02dddcec52031946b	addForeignKeyConstraint baseTableName=cwd_group_admin_user, constraintName=fk_user_target_group, referencedTableName=cwd_group		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-28	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.932127	647	EXECUTED	8:b4c34d6b2c9d68b8bd351b49bd04e863	createIndex indexName=idx_admin_user, tableName=cwd_group_admin_user		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-29	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.937542	648	EXECUTED	8:317aa4a8ffc3665b077ad75d14d1c40d	createIndex indexName=idx_user_target_group, tableName=cwd_group_admin_user		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-30	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.943449	649	EXECUTED	8:e36053ebfea20f245e0eba897ead4aca	addUniqueConstraint constraintName=uk_user_and_target_group, tableName=cwd_group_admin_user		\N	3.6.1	\N	\N	2534357482
BBSDEV-16950-31	fdoherty	liquibase/r6_0/upgrade.xml	2021-06-01 07:59:27.949423	650	EXECUTED	8:2783ad634f95f591998f4f8c45178820	createIndex indexName=idx_granted_perm_dir_map_id, tableName=cwd_granted_perm		\N	3.6.1	\N	\N	2534357482
BBSDEV-19515-1	tkenis	liquibase/r6_2/upgrade.xml	2021-06-01 07:59:27.957114	651	EXECUTED	8:ec9d8e4675361e2fa98f1e69abbde307	createTable tableName=bb_user_dark_feature		\N	3.6.1	\N	\N	2534357482
BBSDEV-19515-2	tkenis	liquibase/r6_2/upgrade.xml	2021-06-01 07:59:27.964588	652	EXECUTED	8:7a9952683f7e803c0fc3d39760870611	addUniqueConstraint constraintName=uq_bb_user_dark_feat_user_feat, tableName=bb_user_dark_feature		\N	3.6.1	\N	\N	2534357482
BBSDEV-18348	tkenis	liquibase/r6_2/upgrade.xml	2021-06-01 07:59:27.968465	653	EXECUTED	8:3e6ce15a37c32178b015afb0306090cb	addColumn tableName=repository		\N	3.6.1	\N	\N	2534357482
BBSDEV-19379-1	bturner	liquibase/r6_2/upgrade.xml	2021-06-01 07:59:27.976723	654	EXECUTED	8:7561dce32e0e8ed7e1539f50db3692a3	createTable tableName=bb_hook_script		\N	3.6.1	\N	\N	2534357482
BBSDEV-19379-2	bturner	liquibase/r6_2/upgrade.xml	2021-06-01 07:59:27.982549	655	EXECUTED	8:0dcafd4430f4fdec6c93f4e427815215	createIndex indexName=idx_bb_hook_script_plugin, tableName=bb_hook_script		\N	3.6.1	\N	\N	2534357482
BBSDEV-19379-3	bturner	liquibase/r6_2/upgrade.xml	2021-06-01 07:59:27.988421	656	EXECUTED	8:f5b69791c7c33ab4332f555f0540e778	createIndex indexName=idx_bb_hook_script_type, tableName=bb_hook_script		\N	3.6.1	\N	\N	2534357482
BBSDEV-19380-1	bturner	liquibase/r6_2/upgrade.xml	2021-06-01 07:59:28.001109	657	EXECUTED	8:c73be4966b975d23de51c862da114d1a	createTable tableName=bb_hook_script_config		\N	3.6.1	\N	\N	2534357482
BBSDEV-19380-2	bturner	liquibase/r6_2/upgrade.xml	2021-06-01 07:59:28.009044	658	EXECUTED	8:7a7841ed1a2d57ba20deddcc662752f5	createIndex indexName=bb_hook_script_cfg_scope, tableName=bb_hook_script_config		\N	3.6.1	\N	\N	2534357482
BBSDEV-19380-3	bturner	liquibase/r6_2/upgrade.xml	2021-06-01 07:59:28.016116	659	EXECUTED	8:856ae3bbd9dfa18e1737a3c74c68f869	createIndex indexName=bb_hook_script_cfg_script, tableName=bb_hook_script_config		\N	3.6.1	\N	\N	2534357482
BBSDEV-19380-4	bturner	liquibase/r6_2/upgrade.xml	2021-06-01 07:59:28.026519	660	EXECUTED	8:3ac8fe2f6fa1be023603eff4654f7175	addUniqueConstraint constraintName=uq_bb_hook_script_config, tableName=bb_hook_script_config		\N	3.6.1	\N	\N	2534357482
BBSDEV-19380-5	bturner	liquibase/r6_2/upgrade.xml	2021-06-01 07:59:28.042694	661	EXECUTED	8:30c0a9088ff9b002b4eb80cb02fbff01	createTable tableName=bb_hook_script_trigger		\N	3.6.1	\N	\N	2534357482
BBSDEV-19380-6	bturner	liquibase/r6_2/upgrade.xml	2021-06-01 07:59:28.050168	662	EXECUTED	8:363a2f5531a00c5edf64191f6c283572	createIndex indexName=idx_bb_hook_script_trgr_config, tableName=bb_hook_script_trigger		\N	3.6.1	\N	\N	2534357482
BBSDEV-19263-1	bturner	liquibase/r6_4/upgrade.xml	2021-06-01 07:59:28.057553	663	EXECUTED	8:147d352425d6b69f3d60d04a29f832e1	createTable tableName=bb_git_pr_cached_merge		\N	3.6.1	\N	\N	2534357482
BBSDEV-19263-2	bturner	liquibase/r6_4/upgrade.xml	2021-06-01 07:59:28.068639	664	EXECUTED	8:d83cfa7d94c4de63349b3d47b081a6e7	createTable tableName=bb_git_pr_common_ancestor		\N	3.6.1	\N	\N	2534357482
BBSDEV-18333	mhenschke	liquibase/r6_4/upgrade.xml	2021-06-01 07:59:28.078265	665	EXECUTED	8:1490c661e0c9f75f08f16b6001e1bdf8	createTable tableName=bb_announcement_banner		\N	3.6.1	\N	\N	2534357482
BBSDEV-19543-1	aermolenko	liquibase/r6_5/upgrade.xml	2021-06-01 07:59:28.092849	666	EXECUTED	8:448e716a62628efc49b470282cf55740	createTable tableName=bb_rl_reject_counter		\N	3.6.1	\N	\N	2534357482
BBSDEV-19543-2	aermolenko	liquibase/r6_5/upgrade.xml	2021-06-01 07:59:28.09925	667	EXECUTED	8:2509d7b29946713ed8ba40923cdf3ffa	createIndex indexName=bb_rl_rej_cntr_intvl_start, tableName=bb_rl_reject_counter		\N	3.6.1	\N	\N	2534357482
BBSDEV-19543-3	aermolenko	liquibase/r6_5/upgrade.xml	2021-06-01 07:59:28.116902	668	EXECUTED	8:94e8fac31fd23e3fc9d69d8d9cfc80b6	createIndex indexName=bb_rl_rej_cntr_usr_id, tableName=bb_rl_reject_counter		\N	3.6.1	\N	\N	2534357482
BBSDEV-19620-1	akord	liquibase/r6_5/upgrade.xml	2021-06-01 07:59:28.138292	669	EXECUTED	8:28f3e07ca02ddd9e93eb3eeb65ce23d7	createTable tableName=bb_rl_user_settings		\N	3.6.1	\N	\N	2534357482
BSERV-11820-1	behumphreys	liquibase/r6_5/upgrade.xml	2021-06-01 07:59:28.150685	670	EXECUTED	8:27e2bd941d5d8ad63131b1daf85d44c1	createIndex indexName=idx_project_lower_name, tableName=project		\N	3.6.1	\N	\N	2534357482
BBSDEV-18815-1	bplump	liquibase/r6_6/upgrade.xml	2021-06-01 07:59:28.164445	671	EXECUTED	8:410481098b066410fab5f5b6ad376914	createTable tableName=bb_suggestion_group		\N	3.6.1	\N	\N	2534357482
BBSDEV-18815-2	fdoherty	liquibase/r6_6/upgrade.xml	2021-06-01 07:59:28.172558	672	EXECUTED	8:6b4b83a219307a1f7528a7dece7d18ff	addColumn tableName=bb_suggestion_group		\N	3.6.1	\N	\N	2534357482
BBSDEV-20617-1	mheemskerk	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.203014	673	MARK_RAN	8:8a7fe8fc51fa3c8cecb2aa6688d9bb38	renameTable newTableName=AO_A0B856_DAILY_COUNTS, oldTableName=AO_371AEF_DAILY_COUNTS		\N	3.6.1	\N	\N	2534357482
BBSDEV-20617-2	mheemskerk	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.246528	674	MARK_RAN	8:b07fc54c5cf8d1e3009c9ca25dae1f6c	renameTable newTableName=AO_A0B856_HIST_INVOCATION, oldTableName=AO_371AEF_HIST_INVOCATION		\N	3.6.1	\N	\N	2534357482
BSERV-10559-1	bturner	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.262974	675	EXECUTED	8:f2eaa8aafe7d40248be2aacf8daa8d25	addColumn tableName=bb_proj_merge_config		\N	3.6.1	\N	\N	2534357482
BSERV-10559-2	bturner	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.270496	676	EXECUTED	8:7b66f23f767ed667ab80a594cc115232	update tableName=bb_proj_merge_config		\N	3.6.1	\N	\N	2534357482
BSERV-10559-3	bturner	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.276366	677	EXECUTED	8:3d1726845c78bd1c5756810335a3dec1	addNotNullConstraint columnName=commit_summaries, tableName=bb_proj_merge_config		\N	3.6.1	\N	\N	2534357482
BSERV-10559-4	bturner	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.28596	678	EXECUTED	8:ca9855efd81bf0fe2f82b53dbe3d555c	addColumn tableName=bb_repo_merge_config		\N	3.6.1	\N	\N	2534357482
BSERV-10559-5	bturner	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.290718	679	EXECUTED	8:c44b0c19d9d955341dfec62f9d1eb0c2	update tableName=bb_repo_merge_config		\N	3.6.1	\N	\N	2534357482
BSERV-10559-6	bturner	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.307307	680	EXECUTED	8:c7b0e19f0d58347ca4b5e45028318863	addNotNullConstraint columnName=commit_summaries, tableName=bb_repo_merge_config		\N	3.6.1	\N	\N	2534357482
BSERV-10559-7	bturner	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.320587	681	EXECUTED	8:018365f3370abecf41b0796f176ce2c9	addColumn tableName=bb_scm_merge_config		\N	3.6.1	\N	\N	2534357482
BSERV-10559-8	bturner	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.332742	682	EXECUTED	8:45e035e56178aef8d9fb2a27d99af2b6	update tableName=bb_scm_merge_config		\N	3.6.1	\N	\N	2534357482
BSERV-10559-9	bturner	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.336717	683	EXECUTED	8:8338793f0ac52e600b77668157c9a38b	addNotNullConstraint columnName=commit_summaries, tableName=bb_scm_merge_config		\N	3.6.1	\N	\N	2534357482
BBSDEV-18969-1	fdoherty	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.354776	684	EXECUTED	8:6222181554a06e3bcaee45ff51f5271c	addColumn tableName=bb_comment		\N	3.6.1	\N	\N	2534357482
BBSDEV-18969-2	fdoherty	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.360251	685	EXECUTED	8:0b205c682defb1f936d9755042c3e4b8	update tableName=bb_comment		\N	3.6.1	\N	\N	2534357482
BBSDEV-18969-3	fdoherty	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.374034	686	EXECUTED	8:13b8e5711c645073db80c2c3ba224158	addNotNullConstraint columnName=severity, tableName=bb_comment		\N	3.6.1	\N	\N	2534357482
BBSDEV-18969-4	fdoherty	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.390011	687	EXECUTED	8:ec98d329679be8d4e259113880088004	addNotNullConstraint columnName=state, tableName=bb_comment		\N	3.6.1	\N	\N	2534357482
BBSDEV-18969-5	fdoherty	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.41197	688	EXECUTED	8:a2fa8bc9b7f12e29a558e8585d8a9aa7	createIndex indexName=idx_bb_comment_severity, tableName=bb_comment		\N	3.6.1	\N	\N	2534357482
BBSDEV-18969-6	fdoherty	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.429739	689	EXECUTED	8:2ae9dcbc206bd87f8f75bf781d02e8b4	createIndex indexName=idx_bb_comment_state, tableName=bb_comment		\N	3.6.1	\N	\N	2534357482
BBSDEV-18969-7	fdoherty	liquibase/r6_7/upgrade.xml	2021-06-01 07:59:28.440013	690	EXECUTED	8:d6360b22aba1075ee0c5ff0d1b1d55a0	createIndex indexName=idx_bb_comment_resolver, tableName=bb_comment		\N	3.6.1	\N	\N	2534357482
BBSDEV-19983-1	bturner	liquibase/r6_6/upgrade.xml	2021-06-01 07:59:28.456354	691	EXECUTED	8:220985c855693e066de4cdaa6d850fa1	createTable tableName=bb_mirror_content_hash		\N	3.6.1	\N	\N	2534357482
BBSDEV-19983-2	bturner	liquibase/r6_6/upgrade.xml	2021-06-01 07:59:28.46877	692	EXECUTED	8:b20e881e29da409af07adcb57c869dce	createTable tableName=bb_mirror_metadata_hash		\N	3.6.1	\N	\N	2534357482
BBSDEV-19009-1	fdoherty	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.486885	693	EXECUTED	8:3673bc72e121daf7bba1bb9b7727fdf4	sql		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-1	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.489432	694	EXECUTED	8:d41d8cd98f00b204e9800998ecf8427e	empty		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-2	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.497906	695	EXECUTED	8:98d2737610294608b9ac7a85328ae359	createTable tableName=cwd_application_saml_config		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-3	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.503225	696	EXECUTED	8:b0f41cd6dc1fad730d7154c0ac07778b	addForeignKeyConstraint baseTableName=cwd_application_saml_config, constraintName=fk_app_sso_config, referencedTableName=cwd_application		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-4	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.508238	697	EXECUTED	8:a760b2556bddcf06400a9e160f38c49f	sql		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-5	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.51216	698	EXECUTED	8:89432466c43ab0662d0f40e0abf8197f	sql		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-6	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.518388	699	EXECUTED	8:a6567da73e39b321cc310dd6308f90cb	createTable tableName=cwd_app_licensing		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-7	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.524639	700	EXECUTED	8:c5916dd0acf38d4fa6d926bf4548ac6b	createIndex indexName=idx_app_id, tableName=cwd_app_licensing		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-8	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.533977	701	EXECUTED	8:ed28bf5a47c408f27c68660153392467	createIndex indexName=idx_app_id_subtype_version, tableName=cwd_app_licensing		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-9	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.548132	702	EXECUTED	8:775b745d2d765a06a2fa283a8885c0db	addForeignKeyConstraint baseTableName=cwd_app_licensing, constraintName=fk_app_id, referencedTableName=cwd_application		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-10	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.556029	703	EXECUTED	8:c88d3dab0b20c7a64bc37f61c397ef78	createTable tableName=cwd_app_licensing_dir_info		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-11	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.567256	704	EXECUTED	8:dca800f147349071fdc490d1da6ff393	createIndex indexName=idx_dir_id, tableName=cwd_app_licensing_dir_info		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-12	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.583445	705	EXECUTED	8:afb41e7b19a30ade50796e00260c04a2	createIndex indexName=idx_summary_id, tableName=cwd_app_licensing_dir_info		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-13	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.591069	706	EXECUTED	8:85d7333b8cc5475398d4e0fdea3d7a33	addForeignKeyConstraint baseTableName=cwd_app_licensing_dir_info, constraintName=fk_licensing_dir_dir_id, referencedTableName=cwd_directory		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-14	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.595568	707	EXECUTED	8:cac7e9618417d8e436ec56bbd3368e0f	addForeignKeyConstraint baseTableName=cwd_app_licensing_dir_info, constraintName=fk_licensing_dir_summary_id, referencedTableName=cwd_app_licensing		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-15	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.607565	708	EXECUTED	8:697f0ce1548cff7e5c0624bd428b6960	createTable tableName=cwd_app_licensed_user		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-16	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.614707	709	EXECUTED	8:7eb2bc42384f809ace0f3920da0a7e14	createIndex indexName=idx_directory_id, tableName=cwd_app_licensed_user		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-17	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.619252	710	EXECUTED	8:e2f2f8dce42688bcd7e89cefec5ce9ba	addForeignKeyConstraint baseTableName=cwd_app_licensed_user, constraintName=fk_licensed_user_dir_id, referencedTableName=cwd_app_licensing_dir_info		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-18	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.62369	711	EXECUTED	8:0283b27be3ece36edd11addfd4148880	addColumn tableName=cwd_property		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-19	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.627814	712	EXECUTED	8:68b352415f1556d7602fa24947a5addc	update tableName=cwd_property		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-20	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.633568	713	EXECUTED	8:e739d5483775cf77a6a1a585474960ed	dropColumn columnName=property_value, tableName=cwd_property		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-21	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.6377	714	EXECUTED	8:f111b9b91fb61a44920bf373ef9a78ee	renameColumn newColumnName=property_value, oldColumnName=property_value_clob, tableName=cwd_property		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-22	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.641308	715	EXECUTED	8:a67e66497432d093e86753ba39fe2cc9	addColumn tableName=cwd_application_attribute		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-23	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.644696	716	EXECUTED	8:ad10da6df0b9099c8e79e9b16f693195	update tableName=cwd_application_attribute		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-24	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.649394	717	EXECUTED	8:31b426aeaa42118fc5052f291e52183e	dropColumn columnName=attribute_value, tableName=cwd_application_attribute		\N	3.6.1	\N	\N	2534357482
BBSDEV-20781-25	wkritzinger	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.653654	718	EXECUTED	8:c417844615935a48deab4dec666b9641	renameColumn newColumnName=attribute_value, oldColumnName=attribute_value_clob, tableName=cwd_application_attribute		\N	3.6.1	\N	\N	2534357482
BBSDEV-21611	bturner	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.703264	719	EXECUTED	8:016633fd38646e075945d1770ffc66e9	dropColumn columnName=merge_hash, tableName=bb_git_pr_cached_merge		\N	3.6.1	\N	\N	2534357482
BBSDEV-19380-1	bplump	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.707567	720	EXECUTED	8:c1fad03ac7ce110ecba787ad22b55faa	dropIndex indexName=bb_hook_script_cfg_scope, tableName=bb_hook_script_config		\N	3.6.1	\N	\N	2534357482
BBSDEV-19380-2	bplump	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.714797	721	EXECUTED	8:5011e05fae1440ba1eaef9e785a36b2a	createIndex indexName=idx_bb_hook_script_cfg_scope, tableName=bb_hook_script_config		\N	3.6.1	\N	\N	2534357482
BBSDEV-19380-3	bplump	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.718826	722	EXECUTED	8:512ae9f7c60a49f42f942c2b53793a04	dropIndex indexName=bb_hook_script_cfg_script, tableName=bb_hook_script_config		\N	3.6.1	\N	\N	2534357482
BBSDEV-19380-4	bplump	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.726507	723	EXECUTED	8:f70120e184d38bd0ac1c79c62d6dc988	createIndex indexName=idx_bb_hook_script_cfg_script, tableName=bb_hook_script_config		\N	3.6.1	\N	\N	2534357482
BBSDEV-21161-1	ckochovski	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.734194	724	EXECUTED	8:766fea306a4ce30a14eee95e4c74bbd9	createTable tableName=bb_attachment		\N	3.6.1	\N	\N	2534357482
BBSDEV-21161-2	ckochovski	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.740049	725	EXECUTED	8:9664e2d51d17c94c84aa7adf467a35bb	createIndex indexName=idx_attachment_repo_id, tableName=bb_attachment		\N	3.6.1	\N	\N	2534357482
BBSDEV-21166-1	ckochovski	liquibase/r7_0/upgrade.xml	2021-06-01 07:59:28.748569	726	EXECUTED	8:368c2fa07c04a56ec22ef3cc9d9ccb3e	createTable tableName=bb_attachment_metadata		\N	3.6.1	\N	\N	2534357482
BBSDEV-21559-1	acarlton	liquibase/r7_1/upgrade.xml	2021-06-01 07:59:28.752531	727	EXECUTED	8:d396d9714e9fe97a245fd8e23ea29847	addColumn tableName=sta_repo_push_ref		\N	3.6.1	\N	\N	2534357482
BBSDEV-21559-2	acarlton	liquibase/r7_1/upgrade.xml	2021-06-01 07:59:28.755841	728	EXECUTED	8:8b40c39c5b93b27aed144c88db26e75a	update tableName=sta_repo_push_ref		\N	3.6.1	\N	\N	2534357482
BBSDEV-21559-3	acarlton	liquibase/r7_1/upgrade.xml	2021-06-01 07:59:28.759573	729	EXECUTED	8:660d8287bb4c9fdc989c6ffd26be1842	addNotNullConstraint columnName=ref_update_type, tableName=sta_repo_push_ref		\N	3.6.1	\N	\N	2534357482
\.


--
-- Data for Name: databasechangeloglock; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.databasechangeloglock (id, locked, lockgranted, lockedby) FROM stdin;
1	f	\N	\N
\.


--
-- Data for Name: hibernate_unique_key; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.hibernate_unique_key (next_hi) FROM stdin;
7
\.


--
-- Data for Name: id_sequence; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.id_sequence (sequence_name, next_val) FROM stdin;
current_app	20
stash_user	100
sta_shared_lob	40
granted_permission	30
plugin_setting	80
sta_remember_me_token	40
project	40
\.


--
-- Data for Name: plugin_setting; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.plugin_setting (namespace, key_name, key_value, id) FROM stdin;
bitbucket.global.settings	com.atlassian.bitbucket.server.bitbucket-git:build	8	1
bitbucket.global.settings	com.atlassian.bitbucket.server.bitbucket-git-lfs.tokensecret	11ehsc80vsvqe5hjgofq9uasieph8ocankbcturi0saq6gg93ghb	2
bitbucket.global.settings	AO_FB71B4_#	6	3
bitbucket.global.settings	AO_6978BB_#	2	4
bitbucket.global.settings	AO_616D7B_#	2	5
bitbucket.global.settings	AO_2AD648_#	1	6
bitbucket.global.settings	AO_9DEC2A_#	1	7
bitbucket.global.settings	AO_811463_#	1	8
bitbucket.global.settings	AO_02A6C0_#	1	9
bitbucket.global.settings	AO_8E6075_#	1	10
bitbucket.global.settings	AO_A0B856_#	1	11
bitbucket.global.settings	com.atlassian.analytics.client.configuration.uuid	a7730256-e856-4a0f-ba15-c4059cec52d0	22
bitbucket.global.settings	com.atlassian.analytics.client.configuration.serverid	BG9T-XUV4-KZZH-B01W	23
bitbucket.global.settings	com.atlassian.analytics.client.configuration..policy_acknowledged	true	24
bitbucket.global.settings	com.atlassian.analytics.client.configuration..analytics_enabled	true	25
bitbucket.global.settings	com.atlassian.upm:notifications:notification-plugin.request	#java.util.List\n	26
bitbucket.global.settings	com.atlassian.upm:notifications:notification-edition.mismatch	#java.util.List\n	27
bitbucket.global.settings	com.atlassian.upm.log.PluginSettingsAuditLogService:log:upm_audit_log_v3	#java.util.List\n{"userKey":"Bitbucket","date":1622534449274,"i18nKey":"upm.auditLog.upm.startup","entryType":"UPM_STARTUP","params":[]}	28
bitbucket.global.settings	com.atlassian.upm:notifications:notification-evaluation.expired	#java.util.List\n	29
bitbucket.global.settings	plugins.lastVersion.server	7002000	30
bitbucket.global.settings	com.atlassian.upm:notifications:notification-evaluation.nearlyexpired	#java.util.List\njira.product.jira-core\njira.product.jira-servicedesk\njira.product.jira-software	31
bitbucket.global.settings	plugins.lastVersion.plugins	2.0.0-m2	32
bitbucket.global.settings	com.atlassian.upm:notifications:notification-maintenance.expired	#java.util.List\n	33
bitbucket.global.settings	com.atlassian.upm:notifications:notification-maintenance.nearlyexpired	#java.util.List\n	34
bitbucket.global.settings	com.atlassian.upm:notifications:notification-license.expired	#java.util.List\n	35
bitbucket.global.settings	com.atlassian.upm:notifications:notification-license.nearlyexpired	#java.util.List\n	36
bitbucket.global.settings	com.atlassian.bitbucket.audit:migration.triggered	false	38
bitbucket.global.settings	com.atlassian.audit.atlassian-audit-plugin:build	2	37
bitbucket.global.settings	com.atlassian.crowd.embedded.admin:build	3	39
bitbucket.global.settings	com.atlassian.bitbucket.server.bitbucket-build:build	2	40
bitbucket.global.settings	com.atlassian.bitbucket.server.bitbucket-bundled-hooks:build	2	41
bitbucket.global.settings	com.atlassian.upm.atlassian-universal-plugin-manager-plugin:build	5	42
bitbucket.global.settings	com.atlassian.plugins.authentication.sso.config.sso-type	NONE	44
bitbucket.global.settings	com.atlassian.plugins.authentication.atlassian-authentication-plugin:build	4	43
bitbucket.global.settings	com.atlassian.bitbucket.server.bitbucket-jira-development-integration:build	5	45
bitbucket.global.settings	com.atlassian.plugins.custom_apps.hasCustomOrder	false	46
bitbucket.global.settings	com.atlassian.plugins.atlassian-nav-links-plugin:build	1	47
bitbucket.global.settings	com.atlassian.upm:notifications:notification-update	#java.util.List\ncom.atlassian.plugins.authentication.atlassian-authentication-plugin	48
bitbucket.global.settings	com.atlassian.analytics.client.configuration..logged_base_analytics_data	true	49
bitbucket-search-checks	search-index-check-ran-bitbucket-search-BSERVDEV-14367	1622534471375	50
\.


--
-- Data for Name: plugin_state; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.plugin_state (name, enabled, updated_timestamp) FROM stdin;
\.


--
-- Data for Name: project; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.project (id, name, project_key, description, project_type, namespace) FROM stdin;
1	ODS Pipeline Test	ODSPIPELINETEST	\N	0	#
\.


--
-- Data for Name: repository; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.repository (id, slug, name, state, project_id, scm_id, hierarchy_id, is_forkable, is_public, store_id, description) FROM stdin;
\.


--
-- Data for Name: repository_access; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.repository_access (user_id, repository_id, last_accessed) FROM stdin;
\.


--
-- Data for Name: sta_activity; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_activity (id, activity_type, created_timestamp, user_id) FROM stdin;
\.


--
-- Data for Name: sta_cmt_disc_activity; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_cmt_disc_activity (activity_id, discussion_id) FROM stdin;
\.


--
-- Data for Name: sta_cmt_disc_participant; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_cmt_disc_participant (id, discussion_id, user_id) FROM stdin;
\.


--
-- Data for Name: sta_cmt_discussion; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_cmt_discussion (id, repository_id, parent_count, commit_id, parent_id) FROM stdin;
\.


--
-- Data for Name: sta_deleted_group; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_deleted_group (id, name, deleted_timestamp) FROM stdin;
\.


--
-- Data for Name: sta_drift_request; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_drift_request (id, pr_id, old_from_hash, old_to_hash, new_from_hash, new_to_hash, attempts) FROM stdin;
\.


--
-- Data for Name: sta_global_permission; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_global_permission (id, perm_id, group_name, user_id) FROM stdin;
1	9	stash-users	\N
12	7	\N	1
\.


--
-- Data for Name: sta_normal_project; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_normal_project (project_id, is_public) FROM stdin;
1	f
\.


--
-- Data for Name: sta_normal_user; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_normal_user (user_id, name, slug, locale, deleted_timestamp, time_zone) FROM stdin;
1	admin	admin	\N	\N	\N
\.


--
-- Data for Name: sta_permission_type; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_permission_type (perm_id, perm_weight) FROM stdin;
0	1000
1	3000
2	2000
3	4000
4	6000
5	7000
6	9000
7	10000
8	5000
9	0
\.


--
-- Data for Name: sta_personal_project; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_personal_project (project_id, owner_id) FROM stdin;
\.


--
-- Data for Name: sta_pr_activity; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_pr_activity (activity_id, pr_id, pr_action) FROM stdin;
\.


--
-- Data for Name: sta_pr_merge_activity; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_pr_merge_activity (activity_id, hash) FROM stdin;
\.


--
-- Data for Name: sta_pr_participant; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_pr_participant (id, pr_id, pr_role, user_id, participant_status, last_reviewed_commit) FROM stdin;
\.


--
-- Data for Name: sta_pr_rescope_activity; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_pr_rescope_activity (activity_id, from_hash, to_hash, prev_from_hash, prev_to_hash, commits_added, commits_removed) FROM stdin;
\.


--
-- Data for Name: sta_pr_rescope_commit; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_pr_rescope_commit (activity_id, changeset_id, action) FROM stdin;
\.


--
-- Data for Name: sta_pr_rescope_request; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_pr_rescope_request (id, repo_id, user_id, created_timestamp) FROM stdin;
\.


--
-- Data for Name: sta_pr_rescope_request_change; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_pr_rescope_request_change (request_id, ref_id, change_type, from_hash, to_hash) FROM stdin;
\.


--
-- Data for Name: sta_project_permission; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_project_permission (id, perm_id, project_id, group_name, user_id) FROM stdin;
13	4	1	\N	1
\.


--
-- Data for Name: sta_pull_request; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_pull_request (id, entity_version, scoped_id, pr_state, created_timestamp, updated_timestamp, from_repository_id, to_repository_id, from_branch_fqn, to_branch_fqn, from_branch_name, to_branch_name, from_hash, to_hash, title, description, locked_timestamp, rescoped_timestamp, closed_timestamp) FROM stdin;
\.


--
-- Data for Name: sta_remember_me_token; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_remember_me_token (id, series, token, user_id, expiry_timestamp, claimed, claimed_address) FROM stdin;
1	f4aec8836e83ffa80860162958d7c5a076ea9b28	8b198da416143a697727afe2b00436907c860bd5	1	2021-07-01 08:00:57.094	f	\N
\.


--
-- Data for Name: sta_repo_activity; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_repo_activity (activity_id, repository_id) FROM stdin;
\.


--
-- Data for Name: sta_repo_hook; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_repo_hook (id, repository_id, hook_key, is_enabled, lob_id, project_id) FROM stdin;
\.


--
-- Data for Name: sta_repo_origin; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_repo_origin (repository_id, origin_id) FROM stdin;
\.


--
-- Data for Name: sta_repo_permission; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_repo_permission (id, perm_id, repo_id, group_name, user_id) FROM stdin;
\.


--
-- Data for Name: sta_repo_push_activity; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_repo_push_activity (activity_id, trigger_id) FROM stdin;
\.


--
-- Data for Name: sta_repo_push_ref; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_repo_push_ref (activity_id, ref_id, change_type, from_hash, to_hash, ref_update_type) FROM stdin;
\.


--
-- Data for Name: sta_repository_scoped_id; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_repository_scoped_id (repository_id, scope_type, next_id) FROM stdin;
\.


--
-- Data for Name: sta_service_user; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_service_user (user_id, display_name, active, name, slug, email_address, label) FROM stdin;
\.


--
-- Data for Name: sta_shared_lob; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_shared_lob (id, lob_data) FROM stdin;
1	{"user.created.version":"7.2.0"}
\.


--
-- Data for Name: sta_task; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_task (id, anchor_id, anchor_type, author_id, context_id, context_type, created_timestamp, task_state, task_text) FROM stdin;
\.


--
-- Data for Name: sta_user_settings; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_user_settings (id, lob_id) FROM stdin;
1	1
\.


--
-- Data for Name: sta_watcher; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.sta_watcher (id, watchable_id, watchable_type, user_id) FROM stdin;
\.


--
-- Data for Name: stash_user; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.stash_user (id) FROM stdin;
1
\.


--
-- Data for Name: trusted_app; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.trusted_app (id, application_id, certificate_timeout, public_key_base64) FROM stdin;
\.


--
-- Data for Name: trusted_app_restriction; Type: TABLE DATA; Schema: public; Owner: bitbucketuser
--

COPY public.trusted_app_restriction (id, trusted_app_id, restriction_type, restriction_value) FROM stdin;
\.


--
-- Name: AO_02A6C0_REJECTED_REF_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_02A6C0_REJECTED_REF_ID_seq"', 1, false);


--
-- Name: AO_0E97B5_REPOSITORY_SHORTCUT_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_0E97B5_REPOSITORY_SHORTCUT_ID_seq"', 1, false);


--
-- Name: AO_2AD648_INSIGHT_ANNOTATION_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_2AD648_INSIGHT_ANNOTATION_ID_seq"', 1, false);


--
-- Name: AO_2AD648_INSIGHT_REPORT_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_2AD648_INSIGHT_REPORT_ID_seq"', 1, false);


--
-- Name: AO_2AD648_MERGE_CHECK_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_2AD648_MERGE_CHECK_ID_seq"', 1, false);


--
-- Name: AO_33D892_COMMENT_JIRA_ISSUE_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_33D892_COMMENT_JIRA_ISSUE_ID_seq"', 1, false);


--
-- Name: AO_38321B_CUSTOM_CONTENT_LINK_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_38321B_CUSTOM_CONTENT_LINK_ID_seq"', 1, false);


--
-- Name: AO_38F373_COMMENT_LIKE_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_38F373_COMMENT_LIKE_ID_seq"', 1, false);


--
-- Name: AO_4789DD_HEALTH_CHECK_STATUS_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_4789DD_HEALTH_CHECK_STATUS_ID_seq"', 1, false);


--
-- Name: AO_4789DD_PROPERTIES_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_4789DD_PROPERTIES_ID_seq"', 1, false);


--
-- Name: AO_4789DD_READ_NOTIFICATIONS_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_4789DD_READ_NOTIFICATIONS_ID_seq"', 1, false);


--
-- Name: AO_4789DD_TASK_MONITOR_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_4789DD_TASK_MONITOR_ID_seq"', 1, false);


--
-- Name: AO_616D7B_BRANCH_MODEL_CONFIG_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_616D7B_BRANCH_MODEL_CONFIG_ID_seq"', 1, true);


--
-- Name: AO_616D7B_BRANCH_TYPE_CONFIG_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_616D7B_BRANCH_TYPE_CONFIG_ID_seq"', 4, true);


--
-- Name: AO_616D7B_BRANCH_TYPE_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_616D7B_BRANCH_TYPE_ID_seq"', 1, false);


--
-- Name: AO_616D7B_SCOPE_AUTO_MERGE_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_616D7B_SCOPE_AUTO_MERGE_ID_seq"', 1, false);


--
-- Name: AO_6978BB_PERMITTED_ENTITY_ENTITY_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_6978BB_PERMITTED_ENTITY_ENTITY_ID_seq"', 1, false);


--
-- Name: AO_6978BB_RESTRICTED_REF_REF_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_6978BB_RESTRICTED_REF_REF_ID_seq"', 1, false);


--
-- Name: AO_777666_JIRA_INDEX_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_777666_JIRA_INDEX_ID_seq"', 1, false);


--
-- Name: AO_811463_GIT_LFS_LOCK_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_811463_GIT_LFS_LOCK_ID_seq"', 1, false);


--
-- Name: AO_8E6075_MIRRORING_REQUEST_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_8E6075_MIRRORING_REQUEST_ID_seq"', 1, false);


--
-- Name: AO_92D5D5_REPO_NOTIFICATION_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_92D5D5_REPO_NOTIFICATION_ID_seq"', 1, false);


--
-- Name: AO_92D5D5_USER_NOTIFICATION_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_92D5D5_USER_NOTIFICATION_ID_seq"', 1, false);


--
-- Name: AO_9DEC2A_DEFAULT_REVIEWER_ENTITY_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_9DEC2A_DEFAULT_REVIEWER_ENTITY_ID_seq"', 1, false);


--
-- Name: AO_9DEC2A_PR_CONDITION_PR_CONDITION_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_9DEC2A_PR_CONDITION_PR_CONDITION_ID_seq"', 1, false);


--
-- Name: AO_A0B856_WEBHOOK_CONFIG_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_A0B856_WEBHOOK_CONFIG_ID_seq"', 1, false);


--
-- Name: AO_A0B856_WEBHOOK_EVENT_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_A0B856_WEBHOOK_EVENT_ID_seq"', 1, false);


--
-- Name: AO_A0B856_WEBHOOK_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_A0B856_WEBHOOK_ID_seq"', 1, false);


--
-- Name: AO_A0B856_WEB_HOOK_LISTENER_AO_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_A0B856_WEB_HOOK_LISTENER_AO_ID_seq"', 1, false);


--
-- Name: AO_BD73C3_PROJECT_AUDIT_AUDIT_ITEM_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_BD73C3_PROJECT_AUDIT_AUDIT_ITEM_ID_seq"', 1, false);


--
-- Name: AO_BD73C3_REPOSITORY_AUDIT_AUDIT_ITEM_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_BD73C3_REPOSITORY_AUDIT_AUDIT_ITEM_ID_seq"', 1, false);


--
-- Name: AO_C77861_AUDIT_ENTITY_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_C77861_AUDIT_ENTITY_ID_seq"', 147, true);


--
-- Name: AO_CFE8FA_BUILD_STATUS_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_CFE8FA_BUILD_STATUS_ID_seq"', 1, false);


--
-- Name: AO_D6A508_IMPORT_JOB_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_D6A508_IMPORT_JOB_ID_seq"', 1, false);


--
-- Name: AO_D6A508_REPO_IMPORT_TASK_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_D6A508_REPO_IMPORT_TASK_ID_seq"', 1, false);


--
-- Name: AO_E5A814_ACCESS_TOKEN_PERM_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_E5A814_ACCESS_TOKEN_PERM_ID_seq"', 2, true);


--
-- Name: AO_ED669C_SEEN_ASSERTIONS_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_ED669C_SEEN_ASSERTIONS_ID_seq"', 1, false);


--
-- Name: AO_F4ED3A_ADD_ON_PROPERTY_AO_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_F4ED3A_ADD_ON_PROPERTY_AO_ID_seq"', 1, false);


--
-- Name: AO_FB71B4_SSH_PUBLIC_KEY_ENTITY_ID_seq; Type: SEQUENCE SET; Schema: public; Owner: bitbucketuser
--

SELECT pg_catalog.setval('public."AO_FB71B4_SSH_PUBLIC_KEY_ENTITY_ID_seq"', 1, false);


--
-- Name: AO_02A6C0_REJECTED_REF AO_02A6C0_REJECTED_REF_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_02A6C0_REJECTED_REF"
    ADD CONSTRAINT "AO_02A6C0_REJECTED_REF_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_02A6C0_SYNC_CONFIG AO_02A6C0_SYNC_CONFIG_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_02A6C0_SYNC_CONFIG"
    ADD CONSTRAINT "AO_02A6C0_SYNC_CONFIG_pkey" PRIMARY KEY ("REPOSITORY_ID");


--
-- Name: AO_0E97B5_REPOSITORY_SHORTCUT AO_0E97B5_REPOSITORY_SHORTCUT_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_0E97B5_REPOSITORY_SHORTCUT"
    ADD CONSTRAINT "AO_0E97B5_REPOSITORY_SHORTCUT_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_2AD648_INSIGHT_ANNOTATION AO_2AD648_INSIGHT_ANNOTATION_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_2AD648_INSIGHT_ANNOTATION"
    ADD CONSTRAINT "AO_2AD648_INSIGHT_ANNOTATION_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_2AD648_INSIGHT_REPORT AO_2AD648_INSIGHT_REPORT_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_2AD648_INSIGHT_REPORT"
    ADD CONSTRAINT "AO_2AD648_INSIGHT_REPORT_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_2AD648_MERGE_CHECK AO_2AD648_MERGE_CHECK_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_2AD648_MERGE_CHECK"
    ADD CONSTRAINT "AO_2AD648_MERGE_CHECK_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_33D892_COMMENT_JIRA_ISSUE AO_33D892_COMMENT_JIRA_ISSUE_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_33D892_COMMENT_JIRA_ISSUE"
    ADD CONSTRAINT "AO_33D892_COMMENT_JIRA_ISSUE_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_38321B_CUSTOM_CONTENT_LINK AO_38321B_CUSTOM_CONTENT_LINK_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_38321B_CUSTOM_CONTENT_LINK"
    ADD CONSTRAINT "AO_38321B_CUSTOM_CONTENT_LINK_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_38F373_COMMENT_LIKE AO_38F373_COMMENT_LIKE_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_38F373_COMMENT_LIKE"
    ADD CONSTRAINT "AO_38F373_COMMENT_LIKE_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_4789DD_HEALTH_CHECK_STATUS AO_4789DD_HEALTH_CHECK_STATUS_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_4789DD_HEALTH_CHECK_STATUS"
    ADD CONSTRAINT "AO_4789DD_HEALTH_CHECK_STATUS_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_4789DD_PROPERTIES AO_4789DD_PROPERTIES_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_4789DD_PROPERTIES"
    ADD CONSTRAINT "AO_4789DD_PROPERTIES_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_4789DD_READ_NOTIFICATIONS AO_4789DD_READ_NOTIFICATIONS_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_4789DD_READ_NOTIFICATIONS"
    ADD CONSTRAINT "AO_4789DD_READ_NOTIFICATIONS_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_4789DD_TASK_MONITOR AO_4789DD_TASK_MONITOR_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_4789DD_TASK_MONITOR"
    ADD CONSTRAINT "AO_4789DD_TASK_MONITOR_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_616D7B_BRANCH_MODEL_CONFIG AO_616D7B_BRANCH_MODEL_CONFIG_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_616D7B_BRANCH_MODEL_CONFIG"
    ADD CONSTRAINT "AO_616D7B_BRANCH_MODEL_CONFIG_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_616D7B_BRANCH_MODEL AO_616D7B_BRANCH_MODEL_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_616D7B_BRANCH_MODEL"
    ADD CONSTRAINT "AO_616D7B_BRANCH_MODEL_pkey" PRIMARY KEY ("REPOSITORY_ID");


--
-- Name: AO_616D7B_BRANCH_TYPE_CONFIG AO_616D7B_BRANCH_TYPE_CONFIG_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_616D7B_BRANCH_TYPE_CONFIG"
    ADD CONSTRAINT "AO_616D7B_BRANCH_TYPE_CONFIG_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_616D7B_BRANCH_TYPE AO_616D7B_BRANCH_TYPE_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_616D7B_BRANCH_TYPE"
    ADD CONSTRAINT "AO_616D7B_BRANCH_TYPE_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_616D7B_SCOPE_AUTO_MERGE AO_616D7B_SCOPE_AUTO_MERGE_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_616D7B_SCOPE_AUTO_MERGE"
    ADD CONSTRAINT "AO_616D7B_SCOPE_AUTO_MERGE_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_6978BB_PERMITTED_ENTITY AO_6978BB_PERMITTED_ENTITY_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_6978BB_PERMITTED_ENTITY"
    ADD CONSTRAINT "AO_6978BB_PERMITTED_ENTITY_pkey" PRIMARY KEY ("ENTITY_ID");


--
-- Name: AO_6978BB_RESTRICTED_REF AO_6978BB_RESTRICTED_REF_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_6978BB_RESTRICTED_REF"
    ADD CONSTRAINT "AO_6978BB_RESTRICTED_REF_pkey" PRIMARY KEY ("REF_ID");


--
-- Name: AO_777666_JIRA_INDEX AO_777666_JIRA_INDEX_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_777666_JIRA_INDEX"
    ADD CONSTRAINT "AO_777666_JIRA_INDEX_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_777666_UPDATED_ISSUES AO_777666_UPDATED_ISSUES_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_777666_UPDATED_ISSUES"
    ADD CONSTRAINT "AO_777666_UPDATED_ISSUES_pkey" PRIMARY KEY ("ISSUE");


--
-- Name: AO_811463_GIT_LFS_LOCK AO_811463_GIT_LFS_LOCK_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_811463_GIT_LFS_LOCK"
    ADD CONSTRAINT "AO_811463_GIT_LFS_LOCK_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_811463_GIT_LFS_REPO_CONFIG AO_811463_GIT_LFS_REPO_CONFIG_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_811463_GIT_LFS_REPO_CONFIG"
    ADD CONSTRAINT "AO_811463_GIT_LFS_REPO_CONFIG_pkey" PRIMARY KEY ("REPOSITORY_ID");


--
-- Name: AO_8E6075_MIRRORING_REQUEST AO_8E6075_MIRRORING_REQUEST_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_8E6075_MIRRORING_REQUEST"
    ADD CONSTRAINT "AO_8E6075_MIRRORING_REQUEST_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_8E6075_MIRROR_SERVER AO_8E6075_MIRROR_SERVER_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_8E6075_MIRROR_SERVER"
    ADD CONSTRAINT "AO_8E6075_MIRROR_SERVER_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_92D5D5_REPO_NOTIFICATION AO_92D5D5_REPO_NOTIFICATION_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_92D5D5_REPO_NOTIFICATION"
    ADD CONSTRAINT "AO_92D5D5_REPO_NOTIFICATION_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_92D5D5_USER_NOTIFICATION AO_92D5D5_USER_NOTIFICATION_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_92D5D5_USER_NOTIFICATION"
    ADD CONSTRAINT "AO_92D5D5_USER_NOTIFICATION_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_9DEC2A_DEFAULT_REVIEWER AO_9DEC2A_DEFAULT_REVIEWER_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_9DEC2A_DEFAULT_REVIEWER"
    ADD CONSTRAINT "AO_9DEC2A_DEFAULT_REVIEWER_pkey" PRIMARY KEY ("ENTITY_ID");


--
-- Name: AO_9DEC2A_PR_CONDITION AO_9DEC2A_PR_CONDITION_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_9DEC2A_PR_CONDITION"
    ADD CONSTRAINT "AO_9DEC2A_PR_CONDITION_pkey" PRIMARY KEY ("PR_CONDITION_ID");


--
-- Name: AO_A0B856_DAILY_COUNTS AO_A0B856_DAILY_COUNTS_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_A0B856_DAILY_COUNTS"
    ADD CONSTRAINT "AO_A0B856_DAILY_COUNTS_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_A0B856_HIST_INVOCATION AO_A0B856_HIST_INVOCATION_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_A0B856_HIST_INVOCATION"
    ADD CONSTRAINT "AO_A0B856_HIST_INVOCATION_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_A0B856_WEBHOOK_CONFIG AO_A0B856_WEBHOOK_CONFIG_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_A0B856_WEBHOOK_CONFIG"
    ADD CONSTRAINT "AO_A0B856_WEBHOOK_CONFIG_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_A0B856_WEBHOOK_EVENT AO_A0B856_WEBHOOK_EVENT_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_A0B856_WEBHOOK_EVENT"
    ADD CONSTRAINT "AO_A0B856_WEBHOOK_EVENT_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_A0B856_WEBHOOK AO_A0B856_WEBHOOK_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_A0B856_WEBHOOK"
    ADD CONSTRAINT "AO_A0B856_WEBHOOK_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_A0B856_WEB_HOOK_LISTENER_AO AO_A0B856_WEB_HOOK_LISTENER_AO_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_A0B856_WEB_HOOK_LISTENER_AO"
    ADD CONSTRAINT "AO_A0B856_WEB_HOOK_LISTENER_AO_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_B586BC_GPG_KEY AO_B586BC_GPG_KEY_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_B586BC_GPG_KEY"
    ADD CONSTRAINT "AO_B586BC_GPG_KEY_pkey" PRIMARY KEY ("FINGERPRINT");


--
-- Name: AO_B586BC_GPG_SUB_KEY AO_B586BC_GPG_SUB_KEY_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_B586BC_GPG_SUB_KEY"
    ADD CONSTRAINT "AO_B586BC_GPG_SUB_KEY_pkey" PRIMARY KEY ("FINGERPRINT");


--
-- Name: AO_BD73C3_PROJECT_AUDIT AO_BD73C3_PROJECT_AUDIT_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_BD73C3_PROJECT_AUDIT"
    ADD CONSTRAINT "AO_BD73C3_PROJECT_AUDIT_pkey" PRIMARY KEY ("AUDIT_ITEM_ID");


--
-- Name: AO_BD73C3_REPOSITORY_AUDIT AO_BD73C3_REPOSITORY_AUDIT_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_BD73C3_REPOSITORY_AUDIT"
    ADD CONSTRAINT "AO_BD73C3_REPOSITORY_AUDIT_pkey" PRIMARY KEY ("AUDIT_ITEM_ID");


--
-- Name: AO_C77861_AUDIT_ENTITY AO_C77861_AUDIT_ENTITY_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_C77861_AUDIT_ENTITY"
    ADD CONSTRAINT "AO_C77861_AUDIT_ENTITY_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_CFE8FA_BUILD_STATUS AO_CFE8FA_BUILD_STATUS_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_CFE8FA_BUILD_STATUS"
    ADD CONSTRAINT "AO_CFE8FA_BUILD_STATUS_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_D6A508_IMPORT_JOB AO_D6A508_IMPORT_JOB_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_D6A508_IMPORT_JOB"
    ADD CONSTRAINT "AO_D6A508_IMPORT_JOB_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_D6A508_REPO_IMPORT_TASK AO_D6A508_REPO_IMPORT_TASK_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_D6A508_REPO_IMPORT_TASK"
    ADD CONSTRAINT "AO_D6A508_REPO_IMPORT_TASK_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_E5A814_ACCESS_TOKEN_PERM AO_E5A814_ACCESS_TOKEN_PERM_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_E5A814_ACCESS_TOKEN_PERM"
    ADD CONSTRAINT "AO_E5A814_ACCESS_TOKEN_PERM_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_E5A814_ACCESS_TOKEN AO_E5A814_ACCESS_TOKEN_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_E5A814_ACCESS_TOKEN"
    ADD CONSTRAINT "AO_E5A814_ACCESS_TOKEN_pkey" PRIMARY KEY ("TOKEN_ID");


--
-- Name: AO_ED669C_SEEN_ASSERTIONS AO_ED669C_SEEN_ASSERTIONS_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_ED669C_SEEN_ASSERTIONS"
    ADD CONSTRAINT "AO_ED669C_SEEN_ASSERTIONS_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_F4ED3A_ADD_ON_PROPERTY_AO AO_F4ED3A_ADD_ON_PROPERTY_AO_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_F4ED3A_ADD_ON_PROPERTY_AO"
    ADD CONSTRAINT "AO_F4ED3A_ADD_ON_PROPERTY_AO_pkey" PRIMARY KEY ("ID");


--
-- Name: AO_FB71B4_SSH_PUBLIC_KEY AO_FB71B4_SSH_PUBLIC_KEY_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_FB71B4_SSH_PUBLIC_KEY"
    ADD CONSTRAINT "AO_FB71B4_SSH_PUBLIC_KEY_pkey" PRIMARY KEY ("ENTITY_ID");


--
-- Name: app_property app_property_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.app_property
    ADD CONSTRAINT app_property_pkey PRIMARY KEY (prop_key);


--
-- Name: bb_alert bb_alert_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_alert
    ADD CONSTRAINT bb_alert_pkey PRIMARY KEY (id);


--
-- Name: bb_clusteredjob bb_clusteredjob_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_clusteredjob
    ADD CONSTRAINT bb_clusteredjob_pkey PRIMARY KEY (job_id);


--
-- Name: bb_comment_thread bb_comment_thread_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_comment_thread
    ADD CONSTRAINT bb_comment_thread_pkey PRIMARY KEY (id);


--
-- Name: bb_integrity_event bb_integrity_event_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_integrity_event
    ADD CONSTRAINT bb_integrity_event_pkey PRIMARY KEY (event_key);


--
-- Name: bb_label bb_label_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_label
    ADD CONSTRAINT bb_label_pkey PRIMARY KEY (id);


--
-- Name: bb_pr_part_status_weight bb_pr_part_status_weight_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_part_status_weight
    ADD CONSTRAINT bb_pr_part_status_weight_pkey PRIMARY KEY (status_id);


--
-- Name: bb_pr_part_status_weight bb_pr_part_status_weight_status_weight_key; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_part_status_weight
    ADD CONSTRAINT bb_pr_part_status_weight_status_weight_key UNIQUE (status_weight);


--
-- Name: bb_rl_user_settings bb_rl_user_settings_user_id_key; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_rl_user_settings
    ADD CONSTRAINT bb_rl_user_settings_user_id_key UNIQUE (user_id);


--
-- Name: changeset changeset_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.changeset
    ADD CONSTRAINT changeset_pkey PRIMARY KEY (id);


--
-- Name: current_app current_app_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.current_app
    ADD CONSTRAINT current_app_pkey PRIMARY KEY (id);


--
-- Name: cwd_app_dir_default_groups cwd_app_dir_default_groups_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_dir_default_groups
    ADD CONSTRAINT cwd_app_dir_default_groups_pkey PRIMARY KEY (id);


--
-- Name: cwd_app_licensed_user cwd_app_licensed_user_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_licensed_user
    ADD CONSTRAINT cwd_app_licensed_user_pkey PRIMARY KEY (id);


--
-- Name: cwd_app_licensing_dir_info cwd_app_licensing_dir_info_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_licensing_dir_info
    ADD CONSTRAINT cwd_app_licensing_dir_info_pkey PRIMARY KEY (id);


--
-- Name: cwd_app_licensing cwd_app_licensing_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_licensing
    ADD CONSTRAINT cwd_app_licensing_pkey PRIMARY KEY (id);


--
-- Name: cwd_application_saml_config cwd_application_saml_config_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_application_saml_config
    ADD CONSTRAINT cwd_application_saml_config_pkey PRIMARY KEY (application_id);


--
-- Name: cwd_group_admin_group cwd_group_admin_group_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_group_admin_group
    ADD CONSTRAINT cwd_group_admin_group_pkey PRIMARY KEY (id);


--
-- Name: cwd_group_admin_user cwd_group_admin_user_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_group_admin_user
    ADD CONSTRAINT cwd_group_admin_user_pkey PRIMARY KEY (id);


--
-- Name: cwd_webhook cwd_webhook_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_webhook
    ADD CONSTRAINT cwd_webhook_pkey PRIMARY KEY (id);


--
-- Name: databasechangeloglock databasechangeloglock_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.databasechangeloglock
    ADD CONSTRAINT databasechangeloglock_pkey PRIMARY KEY (id);


--
-- Name: bb_attachment pk_attachment_id; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_attachment
    ADD CONSTRAINT pk_attachment_id PRIMARY KEY (id);


--
-- Name: bb_attachment_metadata pk_attachment_metadata; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_attachment_metadata
    ADD CONSTRAINT pk_attachment_metadata PRIMARY KEY (attachment_id);


--
-- Name: bb_announcement_banner pk_bb_announcement_banner; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_announcement_banner
    ADD CONSTRAINT pk_bb_announcement_banner PRIMARY KEY (id);


--
-- Name: bb_cmt_disc_comment_activity pk_bb_cmt_disc_com_activity; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_cmt_disc_comment_activity
    ADD CONSTRAINT pk_bb_cmt_disc_com_activity PRIMARY KEY (activity_id);


--
-- Name: bb_comment_parent pk_bb_com_par_comment; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_comment_parent
    ADD CONSTRAINT pk_bb_com_par_comment PRIMARY KEY (comment_id);


--
-- Name: bb_comment pk_bb_comment; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_comment
    ADD CONSTRAINT pk_bb_comment PRIMARY KEY (id);


--
-- Name: bb_data_store pk_bb_data_store; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_data_store
    ADD CONSTRAINT pk_bb_data_store PRIMARY KEY (id);


--
-- Name: bb_git_pr_cached_merge pk_bb_git_pr_cached_merge; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_git_pr_cached_merge
    ADD CONSTRAINT pk_bb_git_pr_cached_merge PRIMARY KEY (id);


--
-- Name: bb_git_pr_common_ancestor pk_bb_git_pr_common_ancestor; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_git_pr_common_ancestor
    ADD CONSTRAINT pk_bb_git_pr_common_ancestor PRIMARY KEY (id);


--
-- Name: bb_hook_script pk_bb_hook_script; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_hook_script
    ADD CONSTRAINT pk_bb_hook_script PRIMARY KEY (id);


--
-- Name: bb_hook_script_config pk_bb_hook_script_config; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_hook_script_config
    ADD CONSTRAINT pk_bb_hook_script_config PRIMARY KEY (id);


--
-- Name: bb_hook_script_trigger pk_bb_hook_script_trigger; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_hook_script_trigger
    ADD CONSTRAINT pk_bb_hook_script_trigger PRIMARY KEY (config_id, trigger_id);


--
-- Name: bb_job pk_bb_job; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_job
    ADD CONSTRAINT pk_bb_job PRIMARY KEY (id);


--
-- Name: bb_job_message pk_bb_job_message; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_job_message
    ADD CONSTRAINT pk_bb_job_message PRIMARY KEY (id);


--
-- Name: bb_mirror_content_hash pk_bb_mirror_content_hash; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_mirror_content_hash
    ADD CONSTRAINT pk_bb_mirror_content_hash PRIMARY KEY (repository_id);


--
-- Name: bb_mirror_metadata_hash pk_bb_mirror_metadata_hash; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_mirror_metadata_hash
    ADD CONSTRAINT pk_bb_mirror_metadata_hash PRIMARY KEY (repository_id);


--
-- Name: bb_pr_comment_activity pk_bb_pr_comment_activity; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_comment_activity
    ADD CONSTRAINT pk_bb_pr_comment_activity PRIMARY KEY (activity_id);


--
-- Name: bb_pr_commit pk_bb_pr_commit; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_commit
    ADD CONSTRAINT pk_bb_pr_commit PRIMARY KEY (pr_id, commit_id);


--
-- Name: bb_pr_reviewer_added pk_bb_pr_reviewer_added; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_reviewer_added
    ADD CONSTRAINT pk_bb_pr_reviewer_added PRIMARY KEY (activity_id, user_id);


--
-- Name: bb_pr_reviewer_removed pk_bb_pr_reviewer_removed; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_reviewer_removed
    ADD CONSTRAINT pk_bb_pr_reviewer_removed PRIMARY KEY (activity_id, user_id);


--
-- Name: bb_pr_reviewer_upd_activity pk_bb_pr_reviewer_upd_act; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_reviewer_upd_activity
    ADD CONSTRAINT pk_bb_pr_reviewer_upd_act PRIMARY KEY (activity_id);


--
-- Name: bb_proj_merge_config pk_bb_proj_merge_config; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_proj_merge_config
    ADD CONSTRAINT pk_bb_proj_merge_config PRIMARY KEY (id);


--
-- Name: bb_proj_merge_strategy pk_bb_proj_merge_strategy; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_proj_merge_strategy
    ADD CONSTRAINT pk_bb_proj_merge_strategy PRIMARY KEY (config_id, strategy_id);


--
-- Name: bb_project_alias pk_bb_project_alias; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_project_alias
    ADD CONSTRAINT pk_bb_project_alias PRIMARY KEY (id);


--
-- Name: bb_repo_merge_config pk_bb_repo_merge_config; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_repo_merge_config
    ADD CONSTRAINT pk_bb_repo_merge_config PRIMARY KEY (id);


--
-- Name: bb_repo_merge_strategy pk_bb_repo_merge_strategy; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_repo_merge_strategy
    ADD CONSTRAINT pk_bb_repo_merge_strategy PRIMARY KEY (config_id, strategy_id);


--
-- Name: bb_repository_alias pk_bb_repository_alias; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_repository_alias
    ADD CONSTRAINT pk_bb_repository_alias PRIMARY KEY (id);


--
-- Name: bb_rl_reject_counter pk_bb_rl_reject_counter; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_rl_reject_counter
    ADD CONSTRAINT pk_bb_rl_reject_counter PRIMARY KEY (id);


--
-- Name: bb_rl_user_settings pk_bb_rl_user_settings; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_rl_user_settings
    ADD CONSTRAINT pk_bb_rl_user_settings PRIMARY KEY (id);


--
-- Name: bb_scm_merge_config pk_bb_scm_merge_config; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_scm_merge_config
    ADD CONSTRAINT pk_bb_scm_merge_config PRIMARY KEY (id);


--
-- Name: bb_scm_merge_strategy pk_bb_scm_merge_strategy; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_scm_merge_strategy
    ADD CONSTRAINT pk_bb_scm_merge_strategy PRIMARY KEY (config_id, strategy_id);


--
-- Name: bb_suggestion_group pk_bb_suggestion_group; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_suggestion_group
    ADD CONSTRAINT pk_bb_suggestion_group PRIMARY KEY (comment_id);


--
-- Name: bb_thread_root_comment pk_bb_thr_root_com_comment; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_thread_root_comment
    ADD CONSTRAINT pk_bb_thr_root_com_comment PRIMARY KEY (thread_id);


--
-- Name: bb_user_dark_feature pk_bb_user_dark_feature; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_user_dark_feature
    ADD CONSTRAINT pk_bb_user_dark_feature PRIMARY KEY (id);


--
-- Name: cs_indexer_state pk_cs_indexer_state; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cs_indexer_state
    ADD CONSTRAINT pk_cs_indexer_state PRIMARY KEY (repository_id, indexer_id);


--
-- Name: cs_repo_membership pk_cs_repo_membership; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cs_repo_membership
    ADD CONSTRAINT pk_cs_repo_membership PRIMARY KEY (cs_id, repository_id);


--
-- Name: cwd_granted_perm pk_cwd_granted_perm; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_granted_perm
    ADD CONSTRAINT pk_cwd_granted_perm PRIMARY KEY (id);


--
-- Name: cwd_tombstone pk_cwd_tombstone; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_tombstone
    ADD CONSTRAINT pk_cwd_tombstone PRIMARY KEY (id);


--
-- Name: sta_global_permission pk_global_permission; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_global_permission
    ADD CONSTRAINT pk_global_permission PRIMARY KEY (id);


--
-- Name: id_sequence pk_id_sequence; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.id_sequence
    ADD CONSTRAINT pk_id_sequence PRIMARY KEY (sequence_name);


--
-- Name: plugin_setting pk_plugin_setting; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.plugin_setting
    ADD CONSTRAINT pk_plugin_setting PRIMARY KEY (id);


--
-- Name: sta_project_permission pk_project_permission; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_project_permission
    ADD CONSTRAINT pk_project_permission PRIMARY KEY (id);


--
-- Name: sta_remember_me_token pk_remember_me_id; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_remember_me_token
    ADD CONSTRAINT pk_remember_me_id PRIMARY KEY (id);


--
-- Name: sta_repo_permission pk_repo_permission; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_permission
    ADD CONSTRAINT pk_repo_permission PRIMARY KEY (id);


--
-- Name: repository_access pk_repository_access; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.repository_access
    ADD CONSTRAINT pk_repository_access PRIMARY KEY (user_id, repository_id);


--
-- Name: sta_activity pk_sta_activity; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_activity
    ADD CONSTRAINT pk_sta_activity PRIMARY KEY (id);


--
-- Name: sta_cmt_disc_activity pk_sta_cmt_disc_activity; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_cmt_disc_activity
    ADD CONSTRAINT pk_sta_cmt_disc_activity PRIMARY KEY (activity_id);


--
-- Name: sta_cmt_disc_participant pk_sta_cmt_disc_participant; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_cmt_disc_participant
    ADD CONSTRAINT pk_sta_cmt_disc_participant PRIMARY KEY (id);


--
-- Name: sta_cmt_discussion pk_sta_cmt_discussion; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_cmt_discussion
    ADD CONSTRAINT pk_sta_cmt_discussion PRIMARY KEY (id);


--
-- Name: sta_deleted_group pk_sta_deleted_group; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_deleted_group
    ADD CONSTRAINT pk_sta_deleted_group PRIMARY KEY (id);


--
-- Name: sta_drift_request pk_sta_drift_request; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_drift_request
    ADD CONSTRAINT pk_sta_drift_request PRIMARY KEY (id);


--
-- Name: sta_normal_project pk_sta_normal_project; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_normal_project
    ADD CONSTRAINT pk_sta_normal_project PRIMARY KEY (project_id);


--
-- Name: sta_normal_user pk_sta_normal_user; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_normal_user
    ADD CONSTRAINT pk_sta_normal_user PRIMARY KEY (user_id);


--
-- Name: sta_personal_project pk_sta_personal_project; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_personal_project
    ADD CONSTRAINT pk_sta_personal_project PRIMARY KEY (project_id);


--
-- Name: sta_pr_activity pk_sta_pr_activity; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_activity
    ADD CONSTRAINT pk_sta_pr_activity PRIMARY KEY (activity_id);


--
-- Name: sta_pr_merge_activity pk_sta_pr_merge_activity; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_merge_activity
    ADD CONSTRAINT pk_sta_pr_merge_activity PRIMARY KEY (activity_id);


--
-- Name: sta_pr_participant pk_sta_pr_participant; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_participant
    ADD CONSTRAINT pk_sta_pr_participant PRIMARY KEY (id);


--
-- Name: sta_pr_rescope_activity pk_sta_pr_rescope_activity; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_rescope_activity
    ADD CONSTRAINT pk_sta_pr_rescope_activity PRIMARY KEY (activity_id);


--
-- Name: sta_pr_rescope_commit pk_sta_pr_rescope_commit; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_rescope_commit
    ADD CONSTRAINT pk_sta_pr_rescope_commit PRIMARY KEY (activity_id, changeset_id);


--
-- Name: sta_pr_rescope_request_change pk_sta_pr_rescope_req_change; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_rescope_request_change
    ADD CONSTRAINT pk_sta_pr_rescope_req_change PRIMARY KEY (request_id, ref_id);


--
-- Name: sta_pr_rescope_request pk_sta_pr_rescope_request; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_rescope_request
    ADD CONSTRAINT pk_sta_pr_rescope_request PRIMARY KEY (id);


--
-- Name: sta_pull_request pk_sta_pull_request; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pull_request
    ADD CONSTRAINT pk_sta_pull_request PRIMARY KEY (id);


--
-- Name: sta_repo_activity pk_sta_repo_activity; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_activity
    ADD CONSTRAINT pk_sta_repo_activity PRIMARY KEY (activity_id);


--
-- Name: sta_repo_hook pk_sta_repo_hook; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_hook
    ADD CONSTRAINT pk_sta_repo_hook PRIMARY KEY (id);


--
-- Name: sta_repo_origin pk_sta_repo_origin; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_origin
    ADD CONSTRAINT pk_sta_repo_origin PRIMARY KEY (repository_id);


--
-- Name: sta_repo_push_activity pk_sta_repo_push_activity; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_push_activity
    ADD CONSTRAINT pk_sta_repo_push_activity PRIMARY KEY (activity_id);


--
-- Name: sta_repo_push_ref pk_sta_repo_push_ref; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_push_ref
    ADD CONSTRAINT pk_sta_repo_push_ref PRIMARY KEY (activity_id, ref_id);


--
-- Name: sta_repository_scoped_id pk_sta_repository_scoped_id; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repository_scoped_id
    ADD CONSTRAINT pk_sta_repository_scoped_id PRIMARY KEY (repository_id, scope_type);


--
-- Name: sta_service_user pk_sta_service_user; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_service_user
    ADD CONSTRAINT pk_sta_service_user PRIMARY KEY (user_id);


--
-- Name: sta_shared_lob pk_sta_shared_lob; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_shared_lob
    ADD CONSTRAINT pk_sta_shared_lob PRIMARY KEY (id);


--
-- Name: sta_task pk_sta_task; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_task
    ADD CONSTRAINT pk_sta_task PRIMARY KEY (id);


--
-- Name: sta_user_settings pk_sta_user_settings; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_user_settings
    ADD CONSTRAINT pk_sta_user_settings PRIMARY KEY (id);


--
-- Name: sta_watcher pk_sta_watcher; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_watcher
    ADD CONSTRAINT pk_sta_watcher PRIMARY KEY (id);


--
-- Name: plugin_state plugin_state_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.plugin_state
    ADD CONSTRAINT plugin_state_pkey PRIMARY KEY (name);


--
-- Name: project project_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.project
    ADD CONSTRAINT project_pkey PRIMARY KEY (id);


--
-- Name: repository repository_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.repository
    ADD CONSTRAINT repository_pkey PRIMARY KEY (id);


--
-- Name: stash_user stash_user_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.stash_user
    ADD CONSTRAINT stash_user_pkey PRIMARY KEY (id);


--
-- Name: cwd_app_dir_group_mapping sys_pk_10069; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_dir_group_mapping
    ADD CONSTRAINT sys_pk_10069 PRIMARY KEY (id);


--
-- Name: cwd_app_dir_mapping sys_pk_10077; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_dir_mapping
    ADD CONSTRAINT sys_pk_10077 PRIMARY KEY (id);


--
-- Name: cwd_app_dir_operation sys_pk_10083; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_dir_operation
    ADD CONSTRAINT sys_pk_10083 PRIMARY KEY (app_dir_mapping_id, operation_type);


--
-- Name: cwd_application sys_pk_10093; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_application
    ADD CONSTRAINT sys_pk_10093 PRIMARY KEY (id);


--
-- Name: cwd_application_address sys_pk_10100; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_application_address
    ADD CONSTRAINT sys_pk_10100 PRIMARY KEY (remote_address);


--
-- Name: cwd_application_alias sys_pk_10108; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_application_alias
    ADD CONSTRAINT sys_pk_10108 PRIMARY KEY (id);


--
-- Name: cwd_application_attribute sys_pk_10116; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_application_attribute
    ADD CONSTRAINT sys_pk_10116 PRIMARY KEY (application_id, attribute_name);


--
-- Name: cwd_directory sys_pk_10127; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_directory
    ADD CONSTRAINT sys_pk_10127 PRIMARY KEY (id);


--
-- Name: cwd_directory_attribute sys_pk_10133; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_directory_attribute
    ADD CONSTRAINT sys_pk_10133 PRIMARY KEY (directory_id, attribute_name);


--
-- Name: cwd_directory_operation sys_pk_10137; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_directory_operation
    ADD CONSTRAINT sys_pk_10137 PRIMARY KEY (directory_id, operation_type);


--
-- Name: cwd_group sys_pk_10148; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_group
    ADD CONSTRAINT sys_pk_10148 PRIMARY KEY (id);


--
-- Name: cwd_group_attribute sys_pk_10156; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_group_attribute
    ADD CONSTRAINT sys_pk_10156 PRIMARY KEY (id);


--
-- Name: cwd_membership sys_pk_10167; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_membership
    ADD CONSTRAINT sys_pk_10167 PRIMARY KEY (id);


--
-- Name: cwd_property sys_pk_10173; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_property
    ADD CONSTRAINT sys_pk_10173 PRIMARY KEY (property_key, property_name);


--
-- Name: cwd_user sys_pk_10194; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_user
    ADD CONSTRAINT sys_pk_10194 PRIMARY KEY (id);


--
-- Name: cwd_user_attribute sys_pk_10202; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_user_attribute
    ADD CONSTRAINT sys_pk_10202 PRIMARY KEY (id);


--
-- Name: cwd_user_credential_record sys_pk_10209; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_user_credential_record
    ADD CONSTRAINT sys_pk_10209 PRIMARY KEY (id);


--
-- Name: trusted_app trusted_app_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.trusted_app
    ADD CONSTRAINT trusted_app_pkey PRIMARY KEY (id);


--
-- Name: trusted_app_restriction trusted_app_restriction_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.trusted_app_restriction
    ADD CONSTRAINT trusted_app_restriction_pkey PRIMARY KEY (id);


--
-- Name: AO_4789DD_TASK_MONITOR u_ao_4789dd_task_mo1827547914; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_4789DD_TASK_MONITOR"
    ADD CONSTRAINT u_ao_4789dd_task_mo1827547914 UNIQUE ("TASK_ID");


--
-- Name: AO_811463_GIT_LFS_LOCK u_ao_811463_git_lfs121924061; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_811463_GIT_LFS_LOCK"
    ADD CONSTRAINT u_ao_811463_git_lfs121924061 UNIQUE ("REPO_PATH_HASH");


--
-- Name: AO_8E6075_MIRROR_SERVER u_ao_8e6075_mirror_881127116; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_8E6075_MIRROR_SERVER"
    ADD CONSTRAINT u_ao_8e6075_mirror_881127116 UNIQUE ("ADD_ON_KEY");


--
-- Name: AO_ED669C_SEEN_ASSERTIONS u_ao_ed669c_seen_as1055534769; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_ED669C_SEEN_ASSERTIONS"
    ADD CONSTRAINT u_ao_ed669c_seen_as1055534769 UNIQUE ("ASSERTION_ID");


--
-- Name: AO_F4ED3A_ADD_ON_PROPERTY_AO u_ao_f4ed3a_add_on_1238639798; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_F4ED3A_ADD_ON_PROPERTY_AO"
    ADD CONSTRAINT u_ao_f4ed3a_add_on_1238639798 UNIQUE ("PRIMARY_KEY");


--
-- Name: cwd_app_dir_default_groups uk_appmapping_groupname; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_dir_default_groups
    ADD CONSTRAINT uk_appmapping_groupname UNIQUE (application_mapping_id, group_name);


--
-- Name: cwd_group_admin_group uk_group_and_target_group; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_group_admin_group
    ADD CONSTRAINT uk_group_and_target_group UNIQUE (group_id, target_group_id);


--
-- Name: cwd_membership uk_mem_dir_parent_child; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_membership
    ADD CONSTRAINT uk_mem_dir_parent_child UNIQUE (directory_id, lower_child_name, lower_parent_name, membership_type);


--
-- Name: project uk_project_key; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.project
    ADD CONSTRAINT uk_project_key UNIQUE (namespace, project_key);


--
-- Name: project uk_project_name; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.project
    ADD CONSTRAINT uk_project_name UNIQUE (namespace, name);


--
-- Name: repository uk_slug_project_id; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.repository
    ADD CONSTRAINT uk_slug_project_id UNIQUE (slug, project_id);


--
-- Name: trusted_app_restriction uk_trusted_app_restrict; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.trusted_app_restriction
    ADD CONSTRAINT uk_trusted_app_restrict UNIQUE (trusted_app_id, restriction_type, restriction_value);


--
-- Name: cwd_group_admin_user uk_user_and_target_group; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_group_admin_user
    ADD CONSTRAINT uk_user_and_target_group UNIQUE (user_id, target_group_id);


--
-- Name: bb_data_store uq_bb_data_store_path; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_data_store
    ADD CONSTRAINT uq_bb_data_store_path UNIQUE (ds_path);


--
-- Name: bb_data_store uq_bb_data_store_uuid; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_data_store
    ADD CONSTRAINT uq_bb_data_store_uuid UNIQUE (ds_uuid);


--
-- Name: bb_hook_script_config uq_bb_hook_script_config; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_hook_script_config
    ADD CONSTRAINT uq_bb_hook_script_config UNIQUE (script_id, scope_id, scope_type);


--
-- Name: bb_proj_merge_config uq_bb_proj_merge_config; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_proj_merge_config
    ADD CONSTRAINT uq_bb_proj_merge_config UNIQUE (project_id, scm_id);


--
-- Name: bb_project_alias uq_bb_project_alias_ns_key; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_project_alias
    ADD CONSTRAINT uq_bb_project_alias_ns_key UNIQUE (namespace, project_key);


--
-- Name: bb_repository_alias uq_bb_repo_alias_key_slug; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_repository_alias
    ADD CONSTRAINT uq_bb_repo_alias_key_slug UNIQUE (project_namespace, project_key, slug);


--
-- Name: bb_repo_merge_config uq_bb_repo_merge_config; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_repo_merge_config
    ADD CONSTRAINT uq_bb_repo_merge_config UNIQUE (repository_id);


--
-- Name: bb_scm_merge_config uq_bb_scm_merge_config; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_scm_merge_config
    ADD CONSTRAINT uq_bb_scm_merge_config UNIQUE (scm_id);


--
-- Name: bb_thread_root_comment uq_bb_thr_root_com_comment_id; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_thread_root_comment
    ADD CONSTRAINT uq_bb_thr_root_com_comment_id UNIQUE (comment_id);


--
-- Name: bb_user_dark_feature uq_bb_user_dark_feat_user_feat; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_user_dark_feature
    ADD CONSTRAINT uq_bb_user_dark_feat_user_feat UNIQUE (user_id, feature_key);


--
-- Name: cwd_user uq_cwd_user_dir_ext_id; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_user
    ADD CONSTRAINT uq_cwd_user_dir_ext_id UNIQUE (directory_id, external_id);


--
-- Name: sta_normal_user uq_normal_user_name; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_normal_user
    ADD CONSTRAINT uq_normal_user_name UNIQUE (name);


--
-- Name: sta_normal_user uq_normal_user_slug; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_normal_user
    ADD CONSTRAINT uq_normal_user_slug UNIQUE (slug);


--
-- Name: plugin_setting uq_plug_setting_ns_key; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.plugin_setting
    ADD CONSTRAINT uq_plug_setting_ns_key UNIQUE (key_name, namespace);


--
-- Name: sta_remember_me_token uq_remember_me_series_token; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_remember_me_token
    ADD CONSTRAINT uq_remember_me_series_token UNIQUE (series, token);


--
-- Name: sta_cmt_disc_participant uq_sta_cmt_disc_part_disc_user; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_cmt_disc_participant
    ADD CONSTRAINT uq_sta_cmt_disc_part_disc_user UNIQUE (discussion_id, user_id);


--
-- Name: sta_cmt_discussion uq_sta_cmt_disc_repo_cmt; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_cmt_discussion
    ADD CONSTRAINT uq_sta_cmt_disc_repo_cmt UNIQUE (repository_id, commit_id);


--
-- Name: sta_deleted_group uq_sta_deleted_group_name; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_deleted_group
    ADD CONSTRAINT uq_sta_deleted_group_name UNIQUE (name);


--
-- Name: sta_personal_project uq_sta_personal_project_owner; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_personal_project
    ADD CONSTRAINT uq_sta_personal_project_owner UNIQUE (owner_id);


--
-- Name: sta_pr_participant uq_sta_pr_participant_pr_user; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_participant
    ADD CONSTRAINT uq_sta_pr_participant_pr_user UNIQUE (pr_id, user_id);


--
-- Name: sta_pull_request uq_sta_pull_request_scoped_id; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pull_request
    ADD CONSTRAINT uq_sta_pull_request_scoped_id UNIQUE (to_repository_id, scoped_id);


--
-- Name: sta_repo_hook uq_sta_repo_hook_scope_hook; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_hook
    ADD CONSTRAINT uq_sta_repo_hook_scope_hook UNIQUE (project_id, repository_id, hook_key);


--
-- Name: sta_service_user uq_sta_service_slug; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_service_user
    ADD CONSTRAINT uq_sta_service_slug UNIQUE (slug);


--
-- Name: sta_service_user uq_sta_service_user_name; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_service_user
    ADD CONSTRAINT uq_sta_service_user_name UNIQUE (name);


--
-- Name: sta_watcher uq_sta_watcher; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_watcher
    ADD CONSTRAINT uq_sta_watcher UNIQUE (watchable_id, watchable_type, user_id);


--
-- Name: sta_permission_type weighted_permission_perm_weight_key; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_permission_type
    ADD CONSTRAINT weighted_permission_perm_weight_key UNIQUE (perm_weight);


--
-- Name: sta_permission_type weighted_permission_pkey; Type: CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_permission_type
    ADD CONSTRAINT weighted_permission_pkey PRIMARY KEY (perm_id);


--
-- Name: bb_alert_issue; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX bb_alert_issue ON public.bb_alert USING btree (issue_id);


--
-- Name: bb_alert_issue_component; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX bb_alert_issue_component ON public.bb_alert USING btree (issue_component_id);


--
-- Name: bb_alert_node_lower; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX bb_alert_node_lower ON public.bb_alert USING btree (node_name_lower);


--
-- Name: bb_alert_plugin_lower; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX bb_alert_plugin_lower ON public.bb_alert USING btree (trigger_plugin_key_lower);


--
-- Name: bb_alert_severity; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX bb_alert_severity ON public.bb_alert USING btree (issue_severity);


--
-- Name: bb_alert_timestamp; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX bb_alert_timestamp ON public.bb_alert USING btree ("timestamp");


--
-- Name: bb_rl_rej_cntr_intvl_start; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX bb_rl_rej_cntr_intvl_start ON public.bb_rl_reject_counter USING btree (interval_start);


--
-- Name: bb_rl_rej_cntr_usr_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX bb_rl_rej_cntr_usr_id ON public.bb_rl_reject_counter USING btree (user_id);


--
-- Name: idx_admin_group; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_admin_group ON public.cwd_group_admin_group USING btree (group_id);


--
-- Name: idx_admin_user; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_admin_user ON public.cwd_group_admin_user USING btree (user_id);


--
-- Name: idx_app_active; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_app_active ON public.cwd_application USING btree (is_active);


--
-- Name: idx_app_address_app_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_app_address_app_id ON public.cwd_application_address USING btree (application_id);


--
-- Name: idx_app_dir_group_group_dir; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_app_dir_group_group_dir ON public.cwd_app_dir_group_mapping USING btree (directory_id, group_name);


--
-- Name: idx_app_dir_grp_mapping_app_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_app_dir_grp_mapping_app_id ON public.cwd_app_dir_group_mapping USING btree (application_id);


--
-- Name: idx_app_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_app_id ON public.cwd_app_licensing USING btree (application_id);


--
-- Name: idx_app_id_subtype_version; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE UNIQUE INDEX idx_app_id_subtype_version ON public.cwd_app_licensing USING btree (application_id, application_subtype, version);


--
-- Name: idx_app_type; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_app_type ON public.cwd_application USING btree (application_type);


--
-- Name: idx_attachment_repo_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_attachment_repo_id ON public.bb_attachment USING btree (repository_id);


--
-- Name: idx_attribute_to_cs; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_attribute_to_cs ON public.cs_attribute USING btree (att_name, att_value);


--
-- Name: idx_bb_clusteredjob_jrk; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_clusteredjob_jrk ON public.bb_clusteredjob USING btree (job_runner_key);


--
-- Name: idx_bb_clusteredjob_next_run; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_clusteredjob_next_run ON public.bb_clusteredjob USING btree (next_run);


--
-- Name: idx_bb_cmt_disc_com_act_com; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_cmt_disc_com_act_com ON public.bb_cmt_disc_comment_activity USING btree (comment_id);


--
-- Name: idx_bb_com_par_parent; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_com_par_parent ON public.bb_comment_parent USING btree (parent_id);


--
-- Name: idx_bb_com_thr_commentable; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_com_thr_commentable ON public.bb_comment_thread USING btree (commentable_type, commentable_id);


--
-- Name: idx_bb_com_thr_diff_type; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_com_thr_diff_type ON public.bb_comment_thread USING btree (diff_type);


--
-- Name: idx_bb_com_thr_from_hash; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_com_thr_from_hash ON public.bb_comment_thread USING btree (from_hash);


--
-- Name: idx_bb_com_thr_to_hash; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_com_thr_to_hash ON public.bb_comment_thread USING btree (to_hash);


--
-- Name: idx_bb_com_thr_to_path; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_com_thr_to_path ON public.bb_comment_thread USING btree (to_path);


--
-- Name: idx_bb_comment_author; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_comment_author ON public.bb_comment USING btree (author_id);


--
-- Name: idx_bb_comment_resolver; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_comment_resolver ON public.bb_comment USING btree (resolver_id);


--
-- Name: idx_bb_comment_severity; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_comment_severity ON public.bb_comment USING btree (severity);


--
-- Name: idx_bb_comment_state; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_comment_state ON public.bb_comment USING btree (state);


--
-- Name: idx_bb_comment_thread; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_comment_thread ON public.bb_comment USING btree (thread_id);


--
-- Name: idx_bb_hook_script_cfg_scope; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_hook_script_cfg_scope ON public.bb_hook_script_config USING btree (scope_id, scope_type);


--
-- Name: idx_bb_hook_script_cfg_script; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_hook_script_cfg_script ON public.bb_hook_script_config USING btree (script_id);


--
-- Name: idx_bb_hook_script_plugin; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_hook_script_plugin ON public.bb_hook_script USING btree (plugin_key);


--
-- Name: idx_bb_hook_script_trgr_config; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_hook_script_trgr_config ON public.bb_hook_script_trigger USING btree (config_id);


--
-- Name: idx_bb_hook_script_type; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_hook_script_type ON public.bb_hook_script USING btree (hook_type);


--
-- Name: idx_bb_job_msg_job_severity; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_job_msg_job_severity ON public.bb_job_message USING btree (job_id, severity);


--
-- Name: idx_bb_job_stash_user; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_job_stash_user ON public.bb_job USING btree (initiator_id);


--
-- Name: idx_bb_job_state_type; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_job_state_type ON public.bb_job USING btree (state, type);


--
-- Name: idx_bb_job_type; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_job_type ON public.bb_job USING btree (type);


--
-- Name: idx_bb_label_map_labelable_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_label_map_labelable_id ON public.bb_label_mapping USING btree (labelable_id);


--
-- Name: idx_bb_label_mapping_label_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_label_mapping_label_id ON public.bb_label_mapping USING btree (label_id);


--
-- Name: idx_bb_pr_com_act_comment; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_pr_com_act_comment ON public.bb_pr_comment_activity USING btree (comment_id);


--
-- Name: idx_bb_pr_commit_commit_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_pr_commit_commit_id ON public.bb_pr_commit USING btree (commit_id);


--
-- Name: idx_bb_proj_alias_proj_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_bb_proj_alias_proj_id ON public.bb_project_alias USING btree (project_id);


--
-- Name: idx_changeset_id_text; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_changeset_id_text ON public.changeset USING btree (id text_pattern_ops);


--
-- Name: idx_cs_repo_membership_repo_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_cs_repo_membership_repo_id ON public.cs_repo_membership USING btree (repository_id);


--
-- Name: idx_cs_to_attribute; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_cs_to_attribute ON public.cs_attribute USING btree (cs_id, att_name);


--
-- Name: idx_cwd_app_dir_mapping_dir_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_cwd_app_dir_mapping_dir_id ON public.cwd_app_dir_mapping USING btree (directory_id);


--
-- Name: idx_cwd_group_directory_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_cwd_group_directory_id ON public.cwd_group USING btree (directory_id);


--
-- Name: idx_cwd_group_external_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_cwd_group_external_id ON public.cwd_group USING btree (external_id);


--
-- Name: idx_cwd_membership_dir_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_cwd_membership_dir_id ON public.cwd_membership USING btree (directory_id);


--
-- Name: idx_cwd_user_user_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_cwd_user_user_id ON public.cwd_user_credential_record USING btree (user_id);


--
-- Name: idx_cwd_webhook_application_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_cwd_webhook_application_id ON public.cwd_webhook USING btree (application_id);


--
-- Name: idx_dir_active; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_dir_active ON public.cwd_directory USING btree (is_active);


--
-- Name: idx_dir_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_dir_id ON public.cwd_app_licensing_dir_info USING btree (directory_id);


--
-- Name: idx_dir_l_impl_class; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_dir_l_impl_class ON public.cwd_directory USING btree (lower_impl_class);


--
-- Name: idx_dir_type; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_dir_type ON public.cwd_directory USING btree (directory_type);


--
-- Name: idx_directory_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_directory_id ON public.cwd_app_licensed_user USING btree (directory_id);


--
-- Name: idx_global_permission_group; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_global_permission_group ON public.sta_global_permission USING btree (group_name);


--
-- Name: idx_global_permission_user; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_global_permission_user ON public.sta_global_permission USING btree (user_id);


--
-- Name: idx_granted_perm_dir_map_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_granted_perm_dir_map_id ON public.cwd_granted_perm USING btree (app_dir_mapping_id);


--
-- Name: idx_group_active; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_group_active ON public.cwd_group USING btree (is_active, directory_id);


--
-- Name: idx_group_attr_dir_name_lval; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_group_attr_dir_name_lval ON public.cwd_group_attribute USING btree (directory_id, attribute_name, attribute_lower_value);


--
-- Name: idx_group_target_group; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_group_target_group ON public.cwd_group_admin_group USING btree (target_group_id);


--
-- Name: idx_label_name; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_label_name ON public.bb_label USING btree (name text_pattern_ops);


--
-- Name: idx_mem_dir_child; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_mem_dir_child ON public.cwd_membership USING btree (membership_type, lower_child_name, directory_id);


--
-- Name: idx_mem_dir_parent; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_mem_dir_parent ON public.cwd_membership USING btree (membership_type, lower_parent_name, directory_id);


--
-- Name: idx_mem_dir_parent_child; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_mem_dir_parent_child ON public.cwd_membership USING btree (membership_type, lower_parent_name, lower_child_name, directory_id);


--
-- Name: idx_pr_rescope_request_pr_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_pr_rescope_request_pr_id ON public.sta_pr_rescope_request USING btree (user_id);


--
-- Name: idx_pr_review_removed_user_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_pr_review_removed_user_id ON public.bb_pr_reviewer_removed USING btree (user_id);


--
-- Name: idx_pr_reviewer_added_user_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_pr_reviewer_added_user_id ON public.bb_pr_reviewer_added USING btree (user_id);


--
-- Name: idx_project_lower_name; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_project_lower_name ON public.project USING btree (lower((name)::text));


--
-- Name: idx_project_permission_group; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_project_permission_group ON public.sta_project_permission USING btree (group_name);


--
-- Name: idx_project_permission_perm_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_project_permission_perm_id ON public.sta_project_permission USING btree (perm_id);


--
-- Name: idx_project_permission_user; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_project_permission_user ON public.sta_project_permission USING btree (user_id);


--
-- Name: idx_project_type; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_project_type ON public.project USING btree (project_type);


--
-- Name: idx_remember_me_token_user_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_remember_me_token_user_id ON public.sta_remember_me_token USING btree (user_id);


--
-- Name: idx_rep_alias_repo_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_rep_alias_repo_id ON public.bb_repository_alias USING btree (repository_id);


--
-- Name: idx_repo_access_repo_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_repo_access_repo_id ON public.repository_access USING btree (repository_id);


--
-- Name: idx_repo_permission_group; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_repo_permission_group ON public.sta_repo_permission USING btree (group_name);


--
-- Name: idx_repo_permission_user; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_repo_permission_user ON public.sta_repo_permission USING btree (user_id);


--
-- Name: idx_repository_access_user_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_repository_access_user_id ON public.repository_access USING btree (user_id);


--
-- Name: idx_repository_hierarchy_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_repository_hierarchy_id ON public.repository USING btree (hierarchy_id);


--
-- Name: idx_repository_project_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_repository_project_id ON public.repository USING btree (project_id);


--
-- Name: idx_repository_state; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_repository_state ON public.repository USING btree (state);


--
-- Name: idx_repository_store_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_repository_store_id ON public.repository USING btree (store_id);


--
-- Name: idx_sta_activity_created_time; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_activity_created_time ON public.sta_activity USING btree (created_timestamp);


--
-- Name: idx_sta_activity_type; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_activity_type ON public.sta_activity USING btree (activity_type);


--
-- Name: idx_sta_activity_user; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_activity_user ON public.sta_activity USING btree (user_id);


--
-- Name: idx_sta_cmt_disc_act_disc; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_cmt_disc_act_disc ON public.sta_cmt_disc_activity USING btree (discussion_id);


--
-- Name: idx_sta_cmt_disc_cmt; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_cmt_disc_cmt ON public.sta_cmt_discussion USING btree (commit_id);


--
-- Name: idx_sta_cmt_disc_part_disc; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_cmt_disc_part_disc ON public.sta_cmt_disc_participant USING btree (discussion_id);


--
-- Name: idx_sta_cmt_disc_part_user; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_cmt_disc_part_user ON public.sta_cmt_disc_participant USING btree (user_id);


--
-- Name: idx_sta_cmt_disc_repo; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_cmt_disc_repo ON public.sta_cmt_discussion USING btree (repository_id);


--
-- Name: idx_sta_deleted_group_ts; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_deleted_group_ts ON public.sta_deleted_group USING btree (deleted_timestamp);


--
-- Name: idx_sta_drift_request_pr_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_drift_request_pr_id ON public.sta_drift_request USING btree (pr_id);


--
-- Name: idx_sta_global_perm_perm_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_global_perm_perm_id ON public.sta_global_permission USING btree (perm_id);


--
-- Name: idx_sta_pr_activity; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_pr_activity ON public.sta_pr_activity USING btree (pr_id, pr_action);


--
-- Name: idx_sta_pr_closed_ts; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_pr_closed_ts ON public.sta_pull_request USING btree (closed_timestamp);


--
-- Name: idx_sta_pr_from_repo_update; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_pr_from_repo_update ON public.sta_pull_request USING btree (from_repository_id, updated_timestamp);


--
-- Name: idx_sta_pr_participant_pr; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_pr_participant_pr ON public.sta_pr_participant USING btree (pr_id);


--
-- Name: idx_sta_pr_participant_user; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_pr_participant_user ON public.sta_pr_participant USING btree (user_id);


--
-- Name: idx_sta_pr_rescope_cmmt_act; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_pr_rescope_cmmt_act ON public.sta_pr_rescope_commit USING btree (activity_id);


--
-- Name: idx_sta_pr_rescope_req_repo; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_pr_rescope_req_repo ON public.sta_pr_rescope_request USING btree (repo_id);


--
-- Name: idx_sta_pr_to_repo_update; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_pr_to_repo_update ON public.sta_pull_request USING btree (to_repository_id, updated_timestamp);


--
-- Name: idx_sta_pr_update_ts; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_pr_update_ts ON public.sta_pull_request USING btree (updated_timestamp);


--
-- Name: idx_sta_proj_perm_pro_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_proj_perm_pro_id ON public.sta_project_permission USING btree (project_id);


--
-- Name: idx_sta_pull_request_from; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_pull_request_from ON public.sta_pull_request USING btree (from_repository_id, from_branch_fqn);


--
-- Name: idx_sta_pull_request_state; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_pull_request_state ON public.sta_pull_request USING btree (pr_state);


--
-- Name: idx_sta_pull_request_to; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_pull_request_to ON public.sta_pull_request USING btree (to_repository_id, to_branch_fqn);


--
-- Name: idx_sta_repo_activity_repo; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_repo_activity_repo ON public.sta_repo_activity USING btree (repository_id);


--
-- Name: idx_sta_repo_hook_hook_key; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_repo_hook_hook_key ON public.sta_repo_hook USING btree (hook_key);


--
-- Name: idx_sta_repo_hook_lob_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_repo_hook_lob_id ON public.sta_repo_hook USING btree (lob_id);


--
-- Name: idx_sta_repo_hook_proj_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_repo_hook_proj_id ON public.sta_repo_hook USING btree (project_id);


--
-- Name: idx_sta_repo_hook_repo_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_repo_hook_repo_id ON public.sta_repo_hook USING btree (repository_id);


--
-- Name: idx_sta_repo_origin_origin_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_repo_origin_origin_id ON public.sta_repo_origin USING btree (origin_id);


--
-- Name: idx_sta_repo_perm_perm_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_repo_perm_perm_id ON public.sta_repo_permission USING btree (perm_id);


--
-- Name: idx_sta_repo_perm_repo_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_repo_perm_repo_id ON public.sta_repo_permission USING btree (repo_id);


--
-- Name: idx_sta_repo_push_ref_activity; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_repo_push_ref_activity ON public.sta_repo_push_ref USING btree (activity_id);


--
-- Name: idx_sta_task_anchor; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_task_anchor ON public.sta_task USING btree (anchor_type, anchor_id);


--
-- Name: idx_sta_task_context; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_task_context ON public.sta_task USING btree (context_type, context_id);


--
-- Name: idx_sta_user_settings_lob_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_user_settings_lob_id ON public.sta_user_settings USING btree (lob_id);


--
-- Name: idx_sta_watcher_user_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_sta_watcher_user_id ON public.sta_watcher USING btree (user_id);


--
-- Name: idx_summary_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_summary_id ON public.cwd_app_licensing_dir_info USING btree (licensing_summary_id);


--
-- Name: idx_tombstone_type_timestamp; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_tombstone_type_timestamp ON public.cwd_tombstone USING btree (tombstone_type, tombstone_timestamp);


--
-- Name: idx_user_active; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_user_active ON public.cwd_user USING btree (is_active, directory_id);


--
-- Name: idx_user_attr_dir_name_lval; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_user_attr_dir_name_lval ON public.cwd_user_attribute USING btree (directory_id, attribute_name, attribute_lower_value);


--
-- Name: idx_user_attr_nval; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_user_attr_nval ON public.cwd_user_attribute USING btree (attribute_name, attribute_numeric_value);


--
-- Name: idx_user_lower_display_name; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_user_lower_display_name ON public.cwd_user USING btree (lower_display_name, directory_id);


--
-- Name: idx_user_lower_email_address; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_user_lower_email_address ON public.cwd_user USING btree (lower_email_address, directory_id);


--
-- Name: idx_user_lower_first_name; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_user_lower_first_name ON public.cwd_user USING btree (lower_first_name, directory_id);


--
-- Name: idx_user_lower_last_name; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_user_lower_last_name ON public.cwd_user USING btree (lower_last_name, directory_id);


--
-- Name: idx_user_target_group; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX idx_user_target_group ON public.cwd_group_admin_user USING btree (target_group_id);


--
-- Name: index_ao_02a6c0_rej1887153917; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_02a6c0_rej1887153917 ON public."AO_02A6C0_REJECTED_REF" USING btree ("REF_ID");


--
-- Name: index_ao_02a6c0_rej2030165690; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_02a6c0_rej2030165690 ON public."AO_02A6C0_REJECTED_REF" USING btree ("REPOSITORY_ID");


--
-- Name: index_ao_0e97b5_rep1393549559; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_0e97b5_rep1393549559 ON public."AO_0E97B5_REPOSITORY_SHORTCUT" USING btree ("REPOSITORY_ID");


--
-- Name: index_ao_0e97b5_rep250640664; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_0e97b5_rep250640664 ON public."AO_0E97B5_REPOSITORY_SHORTCUT" USING btree ("APPLICATION_LINK_ID");


--
-- Name: index_ao_0e97b5_rep309643510; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_0e97b5_rep309643510 ON public."AO_0E97B5_REPOSITORY_SHORTCUT" USING btree ("URL");


--
-- Name: index_ao_2ad648_ins1731502975; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_2ad648_ins1731502975 ON public."AO_2AD648_INSIGHT_ANNOTATION" USING btree ("FK_REPORT_ID");


--
-- Name: index_ao_2ad648_ins1796380851; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_2ad648_ins1796380851 ON public."AO_2AD648_INSIGHT_REPORT" USING btree ("CREATED_DATE");


--
-- Name: index_ao_2ad648_ins282325602; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_2ad648_ins282325602 ON public."AO_2AD648_INSIGHT_REPORT" USING btree ("REPORT_KEY");


--
-- Name: index_ao_2ad648_ins395294165; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_2ad648_ins395294165 ON public."AO_2AD648_INSIGHT_REPORT" USING btree ("COMMIT_ID");


--
-- Name: index_ao_2ad648_ins940577476; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_2ad648_ins940577476 ON public."AO_2AD648_INSIGHT_ANNOTATION" USING btree ("EXTERNAL_ID");


--
-- Name: index_ao_2ad648_ins998130206; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_2ad648_ins998130206 ON public."AO_2AD648_INSIGHT_REPORT" USING btree ("REPOSITORY_ID");


--
-- Name: index_ao_2ad648_mer169660680; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_2ad648_mer169660680 ON public."AO_2AD648_MERGE_CHECK" USING btree ("RESOURCE_ID");


--
-- Name: index_ao_2ad648_mer693118112; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_2ad648_mer693118112 ON public."AO_2AD648_MERGE_CHECK" USING btree ("RESOURCE_ID", "SCOPE_TYPE");


--
-- Name: index_ao_33d892_com451567676; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_33d892_com451567676 ON public."AO_33D892_COMMENT_JIRA_ISSUE" USING btree ("COMMENT_ID");


--
-- Name: index_ao_38321b_cus1828044926; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_38321b_cus1828044926 ON public."AO_38321B_CUSTOM_CONTENT_LINK" USING btree ("CONTENT_KEY");


--
-- Name: index_ao_38f373_com1776971882; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_38f373_com1776971882 ON public."AO_38F373_COMMENT_LIKE" USING btree ("COMMENT_ID");


--
-- Name: index_ao_4789dd_tas42846517; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_4789dd_tas42846517 ON public."AO_4789DD_TASK_MONITOR" USING btree ("TASK_MONITOR_KIND");


--
-- Name: index_ao_616d7b_bra1650093697; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_616d7b_bra1650093697 ON public."AO_616D7B_BRANCH_TYPE_CONFIG" USING btree ("BM_ID");


--
-- Name: index_ao_6978bb_per1013425949; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_6978bb_per1013425949 ON public."AO_6978BB_PERMITTED_ENTITY" USING btree ("FK_RESTRICTED_ID");


--
-- Name: index_ao_6978bb_res1734713733; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_6978bb_res1734713733 ON public."AO_6978BB_RESTRICTED_REF" USING btree ("RESOURCE_ID", "SCOPE_TYPE");


--
-- Name: index_ao_6978bb_res2050912801; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_6978bb_res2050912801 ON public."AO_6978BB_RESTRICTED_REF" USING btree ("REF_TYPE");


--
-- Name: index_ao_6978bb_res847341420; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_6978bb_res847341420 ON public."AO_6978BB_RESTRICTED_REF" USING btree ("REF_VALUE");


--
-- Name: index_ao_777666_jir1125550666; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_777666_jir1125550666 ON public."AO_777666_JIRA_INDEX" USING btree ("PR_ID");


--
-- Name: index_ao_777666_jir1131965929; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_777666_jir1131965929 ON public."AO_777666_JIRA_INDEX" USING btree ("ISSUE");


--
-- Name: index_ao_777666_jir1850935500; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_777666_jir1850935500 ON public."AO_777666_JIRA_INDEX" USING btree ("REPOSITORY");


--
-- Name: index_ao_777666_upd291711106; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_777666_upd291711106 ON public."AO_777666_UPDATED_ISSUES" USING btree ("UPDATE_TIME");


--
-- Name: index_ao_811463_git849077583; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_811463_git849077583 ON public."AO_811463_GIT_LFS_LOCK" USING btree ("REPOSITORY_ID");


--
-- Name: index_ao_811463_git865917084; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_811463_git865917084 ON public."AO_811463_GIT_LFS_LOCK" USING btree ("DIRECTORY_HASH");


--
-- Name: index_ao_8e6075_mir347013871; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_8e6075_mir347013871 ON public."AO_8E6075_MIRRORING_REQUEST" USING btree ("STATE");


--
-- Name: index_ao_8e6075_mir555689735; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_8e6075_mir555689735 ON public."AO_8E6075_MIRRORING_REQUEST" USING btree ("MIRROR_ID");


--
-- Name: index_ao_92d5d5_rep1921630497; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_92d5d5_rep1921630497 ON public."AO_92D5D5_REPO_NOTIFICATION" USING btree ("REPO_ID");


--
-- Name: index_ao_92d5d5_rep679913000; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_92d5d5_rep679913000 ON public."AO_92D5D5_REPO_NOTIFICATION" USING btree ("USER_ID");


--
-- Name: index_ao_92d5d5_rep7900273; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_92d5d5_rep7900273 ON public."AO_92D5D5_REPO_NOTIFICATION" USING btree ("REPO_ID", "USER_ID");


--
-- Name: index_ao_92d5d5_use1759417856; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_92d5d5_use1759417856 ON public."AO_92D5D5_USER_NOTIFICATION" USING btree ("BATCH_SENDER_ID");


--
-- Name: index_ao_9dec2a_def2001216546; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_9dec2a_def2001216546 ON public."AO_9DEC2A_DEFAULT_REVIEWER" USING btree ("FK_RESTRICTED_ID");


--
-- Name: index_ao_9dec2a_pr_1505644799; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_9dec2a_pr_1505644799 ON public."AO_9DEC2A_PR_CONDITION" USING btree ("SOURCE_REF_VALUE");


--
-- Name: index_ao_9dec2a_pr_1891129876; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_9dec2a_pr_1891129876 ON public."AO_9DEC2A_PR_CONDITION" USING btree ("SOURCE_REF_TYPE");


--
-- Name: index_ao_9dec2a_pr_1950938186; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_9dec2a_pr_1950938186 ON public."AO_9DEC2A_PR_CONDITION" USING btree ("RESOURCE_ID", "SCOPE_TYPE");


--
-- Name: index_ao_9dec2a_pr_387063498; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_9dec2a_pr_387063498 ON public."AO_9DEC2A_PR_CONDITION" USING btree ("TARGET_REF_TYPE");


--
-- Name: index_ao_9dec2a_pr_887062261; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_9dec2a_pr_887062261 ON public."AO_9DEC2A_PR_CONDITION" USING btree ("TARGET_REF_VALUE");


--
-- Name: index_ao_a0b856_web1050270930; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_a0b856_web1050270930 ON public."AO_A0B856_WEBHOOK_CONFIG" USING btree ("WEBHOOKID");


--
-- Name: index_ao_a0b856_web110787824; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_a0b856_web110787824 ON public."AO_A0B856_WEBHOOK_EVENT" USING btree ("WEBHOOKID");


--
-- Name: index_ao_b586bc_gpg1041851355; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_b586bc_gpg1041851355 ON public."AO_B586BC_GPG_KEY" USING btree ("USER_ID");


--
-- Name: index_ao_b586bc_gpg1471083652; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_b586bc_gpg1471083652 ON public."AO_B586BC_GPG_SUB_KEY" USING btree ("KEY_ID");


--
-- Name: index_ao_b586bc_gpg594452429; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_b586bc_gpg594452429 ON public."AO_B586BC_GPG_SUB_KEY" USING btree ("FK_GPG_KEY_ID");


--
-- Name: index_ao_b586bc_gpg_key_key_id; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_b586bc_gpg_key_key_id ON public."AO_B586BC_GPG_KEY" USING btree ("KEY_ID");


--
-- Name: index_ao_bd73c3_pro578890136; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_bd73c3_pro578890136 ON public."AO_BD73C3_PROJECT_AUDIT" USING btree ("PROJECT_ID");


--
-- Name: index_ao_bd73c3_rep1781875940; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_bd73c3_rep1781875940 ON public."AO_BD73C3_REPOSITORY_AUDIT" USING btree ("REPOSITORY_ID");


--
-- Name: index_ao_c77861_aud1490488814; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_c77861_aud1490488814 ON public."AO_C77861_AUDIT_ENTITY" USING btree ("USER_ID", "ENTITY_TIMESTAMP");


--
-- Name: index_ao_c77861_aud2071685161; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_c77861_aud2071685161 ON public."AO_C77861_AUDIT_ENTITY" USING btree ("ENTITY_TIMESTAMP", "ID");


--
-- Name: index_ao_c77861_aud237541374; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_c77861_aud237541374 ON public."AO_C77861_AUDIT_ENTITY" USING btree ("PRIMARY_RESOURCE_ID", "PRIMARY_RESOURCE_TYPE", "ENTITY_TIMESTAMP");


--
-- Name: index_ao_c77861_aud470300084; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_c77861_aud470300084 ON public."AO_C77861_AUDIT_ENTITY" USING btree ("SECONDARY_RESOURCE_ID", "SECONDARY_RESOURCE_TYPE", "ENTITY_TIMESTAMP");


--
-- Name: index_ao_cfe8fa_bui803074819; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_cfe8fa_bui803074819 ON public."AO_CFE8FA_BUILD_STATUS" USING btree ("CSID");


--
-- Name: index_ao_d6a508_rep1236870352; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_d6a508_rep1236870352 ON public."AO_D6A508_REPO_IMPORT_TASK" USING btree ("REPOSITORY_ID");


--
-- Name: index_ao_d6a508_rep1615191599; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_d6a508_rep1615191599 ON public."AO_D6A508_REPO_IMPORT_TASK" USING btree ("FAILURE_TYPE");


--
-- Name: index_ao_e5a814_acc680162561; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_e5a814_acc680162561 ON public."AO_E5A814_ACCESS_TOKEN" USING btree ("USER_ID");


--
-- Name: index_ao_e5a814_acc834148545; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_e5a814_acc834148545 ON public."AO_E5A814_ACCESS_TOKEN_PERM" USING btree ("FK_ACCESS_TOKEN_ID");


--
-- Name: index_ao_ed669c_see20117222; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_ed669c_see20117222 ON public."AO_ED669C_SEEN_ASSERTIONS" USING btree ("EXPIRY_TIMESTAMP");


--
-- Name: index_ao_f4ed3a_add50909668; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_f4ed3a_add50909668 ON public."AO_F4ED3A_ADD_ON_PROPERTY_AO" USING btree ("PLUGIN_KEY");


--
-- Name: index_ao_fb71b4_ssh120529590; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_fb71b4_ssh120529590 ON public."AO_FB71B4_SSH_PUBLIC_KEY" USING btree ("LABEL_LOWER");


--
-- Name: index_ao_fb71b4_ssh1382560526; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_fb71b4_ssh1382560526 ON public."AO_FB71B4_SSH_PUBLIC_KEY" USING btree ("KEY_MD5");


--
-- Name: index_ao_fb71b4_ssh714567837; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE INDEX index_ao_fb71b4_ssh714567837 ON public."AO_FB71B4_SSH_PUBLIC_KEY" USING btree ("USER_ID");


--
-- Name: sys_idx_sys_ct_10070_10072; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE UNIQUE INDEX sys_idx_sys_ct_10070_10072 ON public.cwd_app_dir_group_mapping USING btree (app_dir_mapping_id, group_name);


--
-- Name: sys_idx_sys_ct_10078_10080; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE UNIQUE INDEX sys_idx_sys_ct_10078_10080 ON public.cwd_app_dir_mapping USING btree (application_id, directory_id);


--
-- Name: sys_idx_sys_ct_10094_10096; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE UNIQUE INDEX sys_idx_sys_ct_10094_10096 ON public.cwd_application USING btree (lower_application_name);


--
-- Name: sys_idx_sys_ct_10109_10112; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE UNIQUE INDEX sys_idx_sys_ct_10109_10112 ON public.cwd_application_alias USING btree (application_id, lower_user_name);


--
-- Name: sys_idx_sys_ct_10110_10113; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE UNIQUE INDEX sys_idx_sys_ct_10110_10113 ON public.cwd_application_alias USING btree (application_id, lower_alias_name);


--
-- Name: sys_idx_sys_ct_10128_10130; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE UNIQUE INDEX sys_idx_sys_ct_10128_10130 ON public.cwd_directory USING btree (lower_directory_name);


--
-- Name: sys_idx_sys_ct_10149_10151; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE UNIQUE INDEX sys_idx_sys_ct_10149_10151 ON public.cwd_group USING btree (lower_group_name, directory_id);


--
-- Name: sys_idx_sys_ct_10157_10159; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE UNIQUE INDEX sys_idx_sys_ct_10157_10159 ON public.cwd_group_attribute USING btree (group_id, attribute_name, attribute_lower_value);


--
-- Name: sys_idx_sys_ct_10168_10170; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE UNIQUE INDEX sys_idx_sys_ct_10168_10170 ON public.cwd_membership USING btree (parent_id, child_id, membership_type);


--
-- Name: sys_idx_sys_ct_10195_10197; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE UNIQUE INDEX sys_idx_sys_ct_10195_10197 ON public.cwd_user USING btree (lower_user_name, directory_id);


--
-- Name: sys_idx_sys_ct_10203_10205; Type: INDEX; Schema: public; Owner: bitbucketuser
--

CREATE UNIQUE INDEX sys_idx_sys_ct_10203_10205 ON public.cwd_user_attribute USING btree (user_id, attribute_name, attribute_lower_value);


--
-- Name: cwd_group_admin_group fk_admin_group; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_group_admin_group
    ADD CONSTRAINT fk_admin_group FOREIGN KEY (group_id) REFERENCES public.cwd_group(id);


--
-- Name: cwd_group_admin_user fk_admin_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_group_admin_user
    ADD CONSTRAINT fk_admin_user FOREIGN KEY (user_id) REFERENCES public.cwd_user(id) ON DELETE CASCADE;


--
-- Name: cwd_application_alias fk_alias_app_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_application_alias
    ADD CONSTRAINT fk_alias_app_id FOREIGN KEY (application_id) REFERENCES public.cwd_application(id);


--
-- Name: AO_2AD648_INSIGHT_ANNOTATION fk_ao_2ad648_insight_annotation_fk_report_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_2AD648_INSIGHT_ANNOTATION"
    ADD CONSTRAINT fk_ao_2ad648_insight_annotation_fk_report_id FOREIGN KEY ("FK_REPORT_ID") REFERENCES public."AO_2AD648_INSIGHT_REPORT"("ID");


--
-- Name: AO_616D7B_BRANCH_TYPE_CONFIG fk_ao_616d7b_branch_type_config_bm_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_616D7B_BRANCH_TYPE_CONFIG"
    ADD CONSTRAINT fk_ao_616d7b_branch_type_config_bm_id FOREIGN KEY ("BM_ID") REFERENCES public."AO_616D7B_BRANCH_MODEL_CONFIG"("ID");


--
-- Name: AO_9DEC2A_DEFAULT_REVIEWER fk_ao_9dec2a_default_reviewer_fk_restricted_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_9DEC2A_DEFAULT_REVIEWER"
    ADD CONSTRAINT fk_ao_9dec2a_default_reviewer_fk_restricted_id FOREIGN KEY ("FK_RESTRICTED_ID") REFERENCES public."AO_9DEC2A_PR_CONDITION"("PR_CONDITION_ID");


--
-- Name: AO_A0B856_WEBHOOK_CONFIG fk_ao_a0b856_webhook_config_webhookid; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_A0B856_WEBHOOK_CONFIG"
    ADD CONSTRAINT fk_ao_a0b856_webhook_config_webhookid FOREIGN KEY ("WEBHOOKID") REFERENCES public."AO_A0B856_WEBHOOK"("ID");


--
-- Name: AO_A0B856_WEBHOOK_EVENT fk_ao_a0b856_webhook_event_webhookid; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_A0B856_WEBHOOK_EVENT"
    ADD CONSTRAINT fk_ao_a0b856_webhook_event_webhookid FOREIGN KEY ("WEBHOOKID") REFERENCES public."AO_A0B856_WEBHOOK"("ID");


--
-- Name: AO_B586BC_GPG_SUB_KEY fk_ao_b586bc_gpg_sub_key_fk_gpg_key_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_B586BC_GPG_SUB_KEY"
    ADD CONSTRAINT fk_ao_b586bc_gpg_sub_key_fk_gpg_key_id FOREIGN KEY ("FK_GPG_KEY_ID") REFERENCES public."AO_B586BC_GPG_KEY"("FINGERPRINT");


--
-- Name: AO_E5A814_ACCESS_TOKEN_PERM fk_ao_e5a814_access_token_perm_fk_access_token_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public."AO_E5A814_ACCESS_TOKEN_PERM"
    ADD CONSTRAINT fk_ao_e5a814_access_token_perm_fk_access_token_id FOREIGN KEY ("FK_ACCESS_TOKEN_ID") REFERENCES public."AO_E5A814_ACCESS_TOKEN"("TOKEN_ID");


--
-- Name: cwd_app_dir_mapping fk_app_dir_app; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_dir_mapping
    ADD CONSTRAINT fk_app_dir_app FOREIGN KEY (application_id) REFERENCES public.cwd_application(id);


--
-- Name: cwd_app_dir_mapping fk_app_dir_dir; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_dir_mapping
    ADD CONSTRAINT fk_app_dir_dir FOREIGN KEY (directory_id) REFERENCES public.cwd_directory(id);


--
-- Name: cwd_app_dir_group_mapping fk_app_dir_group_app; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_dir_group_mapping
    ADD CONSTRAINT fk_app_dir_group_app FOREIGN KEY (application_id) REFERENCES public.cwd_application(id);


--
-- Name: cwd_app_dir_group_mapping fk_app_dir_group_dir; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_dir_group_mapping
    ADD CONSTRAINT fk_app_dir_group_dir FOREIGN KEY (directory_id) REFERENCES public.cwd_directory(id);


--
-- Name: cwd_app_dir_group_mapping fk_app_dir_group_mapping; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_dir_group_mapping
    ADD CONSTRAINT fk_app_dir_group_mapping FOREIGN KEY (app_dir_mapping_id) REFERENCES public.cwd_app_dir_mapping(id);


--
-- Name: cwd_app_dir_operation fk_app_dir_mapping; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_dir_operation
    ADD CONSTRAINT fk_app_dir_mapping FOREIGN KEY (app_dir_mapping_id) REFERENCES public.cwd_app_dir_mapping(id);


--
-- Name: cwd_app_licensing fk_app_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_licensing
    ADD CONSTRAINT fk_app_id FOREIGN KEY (application_id) REFERENCES public.cwd_application(id);


--
-- Name: cwd_app_dir_default_groups fk_app_mapping; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_dir_default_groups
    ADD CONSTRAINT fk_app_mapping FOREIGN KEY (application_mapping_id) REFERENCES public.cwd_app_dir_mapping(id);


--
-- Name: cwd_application_saml_config fk_app_sso_config; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_application_saml_config
    ADD CONSTRAINT fk_app_sso_config FOREIGN KEY (application_id) REFERENCES public.cwd_application(id) ON DELETE CASCADE;


--
-- Name: cwd_application_address fk_application_address; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_application_address
    ADD CONSTRAINT fk_application_address FOREIGN KEY (application_id) REFERENCES public.cwd_application(id);


--
-- Name: cwd_application_attribute fk_application_attribute; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_application_attribute
    ADD CONSTRAINT fk_application_attribute FOREIGN KEY (application_id) REFERENCES public.cwd_application(id);


--
-- Name: bb_attachment_metadata fk_attachment_metadata_attach; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_attachment_metadata
    ADD CONSTRAINT fk_attachment_metadata_attach FOREIGN KEY (attachment_id) REFERENCES public.bb_attachment(id) ON DELETE CASCADE;


--
-- Name: bb_attachment fk_attachment_repo; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_attachment
    ADD CONSTRAINT fk_attachment_repo FOREIGN KEY (repository_id) REFERENCES public.repository(id) ON DELETE CASCADE;


--
-- Name: bb_cmt_disc_comment_activity fk_bb_cmt_disc_com_act_com; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_cmt_disc_comment_activity
    ADD CONSTRAINT fk_bb_cmt_disc_com_act_com FOREIGN KEY (comment_id) REFERENCES public.bb_comment(id);


--
-- Name: bb_cmt_disc_comment_activity fk_bb_cmt_disc_com_act_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_cmt_disc_comment_activity
    ADD CONSTRAINT fk_bb_cmt_disc_com_act_id FOREIGN KEY (activity_id) REFERENCES public.sta_cmt_disc_activity(activity_id) ON DELETE CASCADE;


--
-- Name: bb_comment_parent fk_bb_com_par_comment_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_comment_parent
    ADD CONSTRAINT fk_bb_com_par_comment_id FOREIGN KEY (comment_id) REFERENCES public.bb_comment(id) ON DELETE CASCADE;


--
-- Name: bb_comment_parent fk_bb_com_par_parent_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_comment_parent
    ADD CONSTRAINT fk_bb_com_par_parent_id FOREIGN KEY (parent_id) REFERENCES public.bb_comment(id);


--
-- Name: bb_comment fk_bb_comment_author; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_comment
    ADD CONSTRAINT fk_bb_comment_author FOREIGN KEY (author_id) REFERENCES public.stash_user(id);


--
-- Name: bb_comment fk_bb_comment_resolver; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_comment
    ADD CONSTRAINT fk_bb_comment_resolver FOREIGN KEY (resolver_id) REFERENCES public.stash_user(id);


--
-- Name: bb_comment fk_bb_comment_thread; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_comment
    ADD CONSTRAINT fk_bb_comment_thread FOREIGN KEY (thread_id) REFERENCES public.bb_comment_thread(id) ON DELETE CASCADE;


--
-- Name: bb_git_pr_cached_merge fk_bb_git_pr_cached_merge_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_git_pr_cached_merge
    ADD CONSTRAINT fk_bb_git_pr_cached_merge_id FOREIGN KEY (id) REFERENCES public.sta_pull_request(id) ON DELETE CASCADE;


--
-- Name: bb_git_pr_common_ancestor fk_bb_git_pr_cmn_ancstr_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_git_pr_common_ancestor
    ADD CONSTRAINT fk_bb_git_pr_cmn_ancstr_id FOREIGN KEY (id) REFERENCES public.sta_pull_request(id) ON DELETE CASCADE;


--
-- Name: bb_hook_script_config fk_bb_hook_script_cfg_script; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_hook_script_config
    ADD CONSTRAINT fk_bb_hook_script_cfg_script FOREIGN KEY (script_id) REFERENCES public.bb_hook_script(id) ON DELETE CASCADE;


--
-- Name: bb_hook_script_trigger fk_bb_hook_script_trgr_config; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_hook_script_trigger
    ADD CONSTRAINT fk_bb_hook_script_trgr_config FOREIGN KEY (config_id) REFERENCES public.bb_hook_script_config(id) ON DELETE CASCADE;


--
-- Name: bb_job fk_bb_job_initiator; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_job
    ADD CONSTRAINT fk_bb_job_initiator FOREIGN KEY (initiator_id) REFERENCES public.stash_user(id) ON DELETE SET NULL;


--
-- Name: bb_job_message fk_bb_job_msg_job; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_job_message
    ADD CONSTRAINT fk_bb_job_msg_job FOREIGN KEY (job_id) REFERENCES public.bb_job(id) ON DELETE CASCADE;


--
-- Name: bb_label_mapping fk_bb_label; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_label_mapping
    ADD CONSTRAINT fk_bb_label FOREIGN KEY (label_id) REFERENCES public.bb_label(id);


--
-- Name: bb_mirror_content_hash fk_bb_mirror_content_hash_repo; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_mirror_content_hash
    ADD CONSTRAINT fk_bb_mirror_content_hash_repo FOREIGN KEY (repository_id) REFERENCES public.repository(id) ON DELETE CASCADE;


--
-- Name: bb_mirror_metadata_hash fk_bb_mirror_mdata_hash_repo; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_mirror_metadata_hash
    ADD CONSTRAINT fk_bb_mirror_mdata_hash_repo FOREIGN KEY (repository_id) REFERENCES public.repository(id) ON DELETE CASCADE;


--
-- Name: bb_pr_comment_activity fk_bb_pr_com_act_comment; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_comment_activity
    ADD CONSTRAINT fk_bb_pr_com_act_comment FOREIGN KEY (comment_id) REFERENCES public.bb_comment(id);


--
-- Name: bb_pr_comment_activity fk_bb_pr_com_act_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_comment_activity
    ADD CONSTRAINT fk_bb_pr_com_act_id FOREIGN KEY (activity_id) REFERENCES public.sta_pr_activity(activity_id) ON DELETE CASCADE;


--
-- Name: bb_pr_commit fk_bb_pr_commit_pr; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_commit
    ADD CONSTRAINT fk_bb_pr_commit_pr FOREIGN KEY (pr_id) REFERENCES public.sta_pull_request(id) ON DELETE CASCADE;


--
-- Name: bb_pr_reviewer_upd_activity fk_bb_pr_reviewer_act_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_reviewer_upd_activity
    ADD CONSTRAINT fk_bb_pr_reviewer_act_id FOREIGN KEY (activity_id) REFERENCES public.sta_pr_activity(activity_id) ON DELETE CASCADE;


--
-- Name: bb_pr_reviewer_added fk_bb_pr_reviewer_added_act; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_reviewer_added
    ADD CONSTRAINT fk_bb_pr_reviewer_added_act FOREIGN KEY (activity_id) REFERENCES public.bb_pr_reviewer_upd_activity(activity_id) ON DELETE CASCADE;


--
-- Name: bb_pr_reviewer_added fk_bb_pr_reviewer_added_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_reviewer_added
    ADD CONSTRAINT fk_bb_pr_reviewer_added_user FOREIGN KEY (user_id) REFERENCES public.stash_user(id);


--
-- Name: bb_pr_reviewer_removed fk_bb_pr_reviewer_removed_act; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_reviewer_removed
    ADD CONSTRAINT fk_bb_pr_reviewer_removed_act FOREIGN KEY (activity_id) REFERENCES public.bb_pr_reviewer_upd_activity(activity_id) ON DELETE CASCADE;


--
-- Name: bb_pr_reviewer_removed fk_bb_pr_reviewer_removed_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_pr_reviewer_removed
    ADD CONSTRAINT fk_bb_pr_reviewer_removed_user FOREIGN KEY (user_id) REFERENCES public.stash_user(id);


--
-- Name: bb_proj_merge_config fk_bb_proj_merge_config; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_proj_merge_config
    ADD CONSTRAINT fk_bb_proj_merge_config FOREIGN KEY (project_id) REFERENCES public.project(id) ON DELETE CASCADE;


--
-- Name: bb_proj_merge_strategy fk_bb_proj_merge_strategy; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_proj_merge_strategy
    ADD CONSTRAINT fk_bb_proj_merge_strategy FOREIGN KEY (config_id) REFERENCES public.bb_proj_merge_config(id) ON DELETE CASCADE;


--
-- Name: bb_repo_merge_config fk_bb_repo_merge_config; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_repo_merge_config
    ADD CONSTRAINT fk_bb_repo_merge_config FOREIGN KEY (repository_id) REFERENCES public.repository(id) ON DELETE CASCADE;


--
-- Name: bb_repo_merge_strategy fk_bb_repo_merge_strategy; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_repo_merge_strategy
    ADD CONSTRAINT fk_bb_repo_merge_strategy FOREIGN KEY (config_id) REFERENCES public.bb_repo_merge_config(id) ON DELETE CASCADE;


--
-- Name: bb_rl_reject_counter fk_bb_rl_rejected_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_rl_reject_counter
    ADD CONSTRAINT fk_bb_rl_rejected_user FOREIGN KEY (user_id) REFERENCES public.stash_user(id) ON DELETE CASCADE;


--
-- Name: bb_rl_user_settings fk_bb_rl_user_settings_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_rl_user_settings
    ADD CONSTRAINT fk_bb_rl_user_settings_user FOREIGN KEY (user_id) REFERENCES public.stash_user(id) ON DELETE CASCADE;


--
-- Name: bb_scm_merge_strategy fk_bb_scm_merge_strategy; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_scm_merge_strategy
    ADD CONSTRAINT fk_bb_scm_merge_strategy FOREIGN KEY (config_id) REFERENCES public.bb_scm_merge_config(id) ON DELETE CASCADE;


--
-- Name: bb_suggestion_group fk_bb_sugg_grp_comment; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_suggestion_group
    ADD CONSTRAINT fk_bb_sugg_grp_comment FOREIGN KEY (comment_id) REFERENCES public.bb_comment(id) ON DELETE CASCADE;


--
-- Name: bb_thread_root_comment fk_bb_thr_root_com_comment_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_thread_root_comment
    ADD CONSTRAINT fk_bb_thr_root_com_comment_id FOREIGN KEY (comment_id) REFERENCES public.bb_comment(id);


--
-- Name: bb_thread_root_comment fk_bb_thr_root_com_thread_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_thread_root_comment
    ADD CONSTRAINT fk_bb_thr_root_com_thread_id FOREIGN KEY (thread_id) REFERENCES public.bb_comment_thread(id) ON DELETE CASCADE;


--
-- Name: bb_user_dark_feature fk_bb_user_dark_feature_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_user_dark_feature
    ADD CONSTRAINT fk_bb_user_dark_feature_user FOREIGN KEY (user_id) REFERENCES public.stash_user(id);


--
-- Name: cs_attribute fk_cs_attribute_changeset; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cs_attribute
    ADD CONSTRAINT fk_cs_attribute_changeset FOREIGN KEY (cs_id) REFERENCES public.changeset(id);


--
-- Name: cs_indexer_state fk_cs_indexer_state_repository; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cs_indexer_state
    ADD CONSTRAINT fk_cs_indexer_state_repository FOREIGN KEY (repository_id) REFERENCES public.repository(id) ON DELETE CASCADE;


--
-- Name: cwd_directory_attribute fk_directory_attribute; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_directory_attribute
    ADD CONSTRAINT fk_directory_attribute FOREIGN KEY (directory_id) REFERENCES public.cwd_directory(id);


--
-- Name: cwd_group fk_directory_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_group
    ADD CONSTRAINT fk_directory_id FOREIGN KEY (directory_id) REFERENCES public.cwd_directory(id);


--
-- Name: cwd_directory_operation fk_directory_operation; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_directory_operation
    ADD CONSTRAINT fk_directory_operation FOREIGN KEY (directory_id) REFERENCES public.cwd_directory(id);


--
-- Name: sta_global_permission fk_global_permission_type; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_global_permission
    ADD CONSTRAINT fk_global_permission_type FOREIGN KEY (perm_id) REFERENCES public.sta_permission_type(perm_id);


--
-- Name: sta_global_permission fk_global_permission_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_global_permission
    ADD CONSTRAINT fk_global_permission_user FOREIGN KEY (user_id) REFERENCES public.stash_user(id);


--
-- Name: cwd_granted_perm fk_granted_perm_dir_mapping; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_granted_perm
    ADD CONSTRAINT fk_granted_perm_dir_mapping FOREIGN KEY (app_dir_mapping_id) REFERENCES public.cwd_app_dir_mapping(id);


--
-- Name: cwd_group_attribute fk_group_attr_dir_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_group_attribute
    ADD CONSTRAINT fk_group_attr_dir_id FOREIGN KEY (directory_id) REFERENCES public.cwd_directory(id);


--
-- Name: cwd_group_attribute fk_group_attr_id_group_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_group_attribute
    ADD CONSTRAINT fk_group_attr_id_group_id FOREIGN KEY (group_id) REFERENCES public.cwd_group(id);


--
-- Name: cwd_group_admin_group fk_group_target_group; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_group_admin_group
    ADD CONSTRAINT fk_group_target_group FOREIGN KEY (target_group_id) REFERENCES public.cwd_group(id) ON DELETE CASCADE;


--
-- Name: cwd_app_licensed_user fk_licensed_user_dir_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_licensed_user
    ADD CONSTRAINT fk_licensed_user_dir_id FOREIGN KEY (directory_id) REFERENCES public.cwd_app_licensing_dir_info(id);


--
-- Name: cwd_app_licensing_dir_info fk_licensing_dir_dir_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_licensing_dir_info
    ADD CONSTRAINT fk_licensing_dir_dir_id FOREIGN KEY (directory_id) REFERENCES public.cwd_directory(id);


--
-- Name: cwd_app_licensing_dir_info fk_licensing_dir_summary_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_app_licensing_dir_info
    ADD CONSTRAINT fk_licensing_dir_summary_id FOREIGN KEY (licensing_summary_id) REFERENCES public.cwd_app_licensing(id);


--
-- Name: cwd_membership fk_membership_dir; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_membership
    ADD CONSTRAINT fk_membership_dir FOREIGN KEY (directory_id) REFERENCES public.cwd_directory(id);


--
-- Name: bb_project_alias fk_project_alias_project; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_project_alias
    ADD CONSTRAINT fk_project_alias_project FOREIGN KEY (project_id) REFERENCES public.project(id) ON DELETE CASCADE;


--
-- Name: sta_project_permission fk_project_permission_project; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_project_permission
    ADD CONSTRAINT fk_project_permission_project FOREIGN KEY (project_id) REFERENCES public.project(id) ON DELETE CASCADE;


--
-- Name: sta_project_permission fk_project_permission_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_project_permission
    ADD CONSTRAINT fk_project_permission_user FOREIGN KEY (user_id) REFERENCES public.stash_user(id);


--
-- Name: sta_project_permission fk_project_permission_weight; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_project_permission
    ADD CONSTRAINT fk_project_permission_weight FOREIGN KEY (perm_id) REFERENCES public.sta_permission_type(perm_id);


--
-- Name: sta_remember_me_token fk_remember_me_user_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_remember_me_token
    ADD CONSTRAINT fk_remember_me_user_id FOREIGN KEY (user_id) REFERENCES public.stash_user(id) ON DELETE CASCADE;


--
-- Name: cs_repo_membership fk_repo_membership_changeset; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cs_repo_membership
    ADD CONSTRAINT fk_repo_membership_changeset FOREIGN KEY (cs_id) REFERENCES public.changeset(id);


--
-- Name: cs_repo_membership fk_repo_membership_repo; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cs_repo_membership
    ADD CONSTRAINT fk_repo_membership_repo FOREIGN KEY (repository_id) REFERENCES public.repository(id) ON DELETE CASCADE;


--
-- Name: sta_repo_permission fk_repo_permission_repo; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_permission
    ADD CONSTRAINT fk_repo_permission_repo FOREIGN KEY (repo_id) REFERENCES public.repository(id) ON DELETE CASCADE;


--
-- Name: sta_repo_permission fk_repo_permission_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_permission
    ADD CONSTRAINT fk_repo_permission_user FOREIGN KEY (user_id) REFERENCES public.stash_user(id);


--
-- Name: sta_repo_permission fk_repo_permission_weight; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_permission
    ADD CONSTRAINT fk_repo_permission_weight FOREIGN KEY (perm_id) REFERENCES public.sta_permission_type(perm_id);


--
-- Name: repository_access fk_repository_access_repo; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.repository_access
    ADD CONSTRAINT fk_repository_access_repo FOREIGN KEY (repository_id) REFERENCES public.repository(id);


--
-- Name: repository_access fk_repository_access_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.repository_access
    ADD CONSTRAINT fk_repository_access_user FOREIGN KEY (user_id) REFERENCES public.stash_user(id);


--
-- Name: bb_repository_alias fk_repository_alias_repo; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.bb_repository_alias
    ADD CONSTRAINT fk_repository_alias_repo FOREIGN KEY (repository_id) REFERENCES public.repository(id) ON DELETE CASCADE;


--
-- Name: repository fk_repository_project; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.repository
    ADD CONSTRAINT fk_repository_project FOREIGN KEY (project_id) REFERENCES public.project(id);


--
-- Name: repository fk_repository_store_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.repository
    ADD CONSTRAINT fk_repository_store_id FOREIGN KEY (store_id) REFERENCES public.bb_data_store(id);


--
-- Name: sta_activity fk_sta_activity_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_activity
    ADD CONSTRAINT fk_sta_activity_user FOREIGN KEY (user_id) REFERENCES public.stash_user(id);


--
-- Name: sta_cmt_disc_activity fk_sta_cmt_disc_act_disc; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_cmt_disc_activity
    ADD CONSTRAINT fk_sta_cmt_disc_act_disc FOREIGN KEY (discussion_id) REFERENCES public.sta_cmt_discussion(id);


--
-- Name: sta_cmt_disc_activity fk_sta_cmt_disc_act_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_cmt_disc_activity
    ADD CONSTRAINT fk_sta_cmt_disc_act_id FOREIGN KEY (activity_id) REFERENCES public.sta_repo_activity(activity_id) ON DELETE CASCADE;


--
-- Name: sta_cmt_disc_participant fk_sta_cmt_disc_part_disc; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_cmt_disc_participant
    ADD CONSTRAINT fk_sta_cmt_disc_part_disc FOREIGN KEY (discussion_id) REFERENCES public.sta_cmt_discussion(id) ON DELETE CASCADE;


--
-- Name: sta_cmt_disc_participant fk_sta_cmt_disc_part_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_cmt_disc_participant
    ADD CONSTRAINT fk_sta_cmt_disc_part_user FOREIGN KEY (user_id) REFERENCES public.stash_user(id);


--
-- Name: sta_cmt_discussion fk_sta_cmt_disc_repo; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_cmt_discussion
    ADD CONSTRAINT fk_sta_cmt_disc_repo FOREIGN KEY (repository_id) REFERENCES public.repository(id) ON DELETE CASCADE;


--
-- Name: sta_drift_request fk_sta_drift_request_pr; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_drift_request
    ADD CONSTRAINT fk_sta_drift_request_pr FOREIGN KEY (pr_id) REFERENCES public.sta_pull_request(id) ON DELETE CASCADE;


--
-- Name: sta_normal_project fk_sta_normal_project_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_normal_project
    ADD CONSTRAINT fk_sta_normal_project_id FOREIGN KEY (project_id) REFERENCES public.project(id) ON DELETE CASCADE;


--
-- Name: sta_normal_user fk_sta_normal_user_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_normal_user
    ADD CONSTRAINT fk_sta_normal_user_id FOREIGN KEY (user_id) REFERENCES public.stash_user(id) ON DELETE CASCADE;


--
-- Name: sta_personal_project fk_sta_personal_project_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_personal_project
    ADD CONSTRAINT fk_sta_personal_project_id FOREIGN KEY (project_id) REFERENCES public.project(id) ON DELETE CASCADE;


--
-- Name: sta_personal_project fk_sta_personal_project_owner; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_personal_project
    ADD CONSTRAINT fk_sta_personal_project_owner FOREIGN KEY (owner_id) REFERENCES public.stash_user(id);


--
-- Name: sta_pr_activity fk_sta_pr_activity_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_activity
    ADD CONSTRAINT fk_sta_pr_activity_id FOREIGN KEY (activity_id) REFERENCES public.sta_activity(id) ON DELETE CASCADE;


--
-- Name: sta_pr_activity fk_sta_pr_activity_pr; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_activity
    ADD CONSTRAINT fk_sta_pr_activity_pr FOREIGN KEY (pr_id) REFERENCES public.sta_pull_request(id);


--
-- Name: sta_pr_merge_activity fk_sta_pr_mrg_act_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_merge_activity
    ADD CONSTRAINT fk_sta_pr_mrg_act_id FOREIGN KEY (activity_id) REFERENCES public.sta_pr_activity(activity_id) ON DELETE CASCADE;


--
-- Name: sta_pr_participant fk_sta_pr_participant_pr; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_participant
    ADD CONSTRAINT fk_sta_pr_participant_pr FOREIGN KEY (pr_id) REFERENCES public.sta_pull_request(id) ON DELETE CASCADE;


--
-- Name: sta_pr_participant fk_sta_pr_participant_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_participant
    ADD CONSTRAINT fk_sta_pr_participant_user FOREIGN KEY (user_id) REFERENCES public.stash_user(id);


--
-- Name: sta_pr_rescope_activity fk_sta_pr_rescope_act_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_rescope_activity
    ADD CONSTRAINT fk_sta_pr_rescope_act_id FOREIGN KEY (activity_id) REFERENCES public.sta_pr_activity(activity_id);


--
-- Name: sta_pr_rescope_request_change fk_sta_pr_rescope_ch_req_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_rescope_request_change
    ADD CONSTRAINT fk_sta_pr_rescope_ch_req_id FOREIGN KEY (request_id) REFERENCES public.sta_pr_rescope_request(id) ON DELETE CASCADE;


--
-- Name: sta_pr_rescope_commit fk_sta_pr_rescope_cmmt_act; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_rescope_commit
    ADD CONSTRAINT fk_sta_pr_rescope_cmmt_act FOREIGN KEY (activity_id) REFERENCES public.sta_pr_rescope_activity(activity_id) ON DELETE CASCADE;


--
-- Name: sta_pr_rescope_request fk_sta_pr_rescope_req_repo; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_rescope_request
    ADD CONSTRAINT fk_sta_pr_rescope_req_repo FOREIGN KEY (repo_id) REFERENCES public.repository(id) ON DELETE CASCADE;


--
-- Name: sta_pr_rescope_request fk_sta_pr_rescope_req_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pr_rescope_request
    ADD CONSTRAINT fk_sta_pr_rescope_req_user FOREIGN KEY (user_id) REFERENCES public.stash_user(id) ON DELETE CASCADE;


--
-- Name: sta_pull_request fk_sta_pull_request_from_repo; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pull_request
    ADD CONSTRAINT fk_sta_pull_request_from_repo FOREIGN KEY (from_repository_id) REFERENCES public.repository(id);


--
-- Name: sta_pull_request fk_sta_pull_request_to_repo; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_pull_request
    ADD CONSTRAINT fk_sta_pull_request_to_repo FOREIGN KEY (to_repository_id) REFERENCES public.repository(id) ON DELETE CASCADE;


--
-- Name: sta_repo_activity fk_sta_repo_activity_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_activity
    ADD CONSTRAINT fk_sta_repo_activity_id FOREIGN KEY (activity_id) REFERENCES public.sta_activity(id) ON DELETE CASCADE;


--
-- Name: sta_repo_activity fk_sta_repo_activity_repo; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_activity
    ADD CONSTRAINT fk_sta_repo_activity_repo FOREIGN KEY (repository_id) REFERENCES public.repository(id);


--
-- Name: sta_repo_hook fk_sta_repo_hook_lob; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_hook
    ADD CONSTRAINT fk_sta_repo_hook_lob FOREIGN KEY (lob_id) REFERENCES public.sta_shared_lob(id);


--
-- Name: sta_repo_hook fk_sta_repo_hook_proj; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_hook
    ADD CONSTRAINT fk_sta_repo_hook_proj FOREIGN KEY (project_id) REFERENCES public.project(id) ON DELETE CASCADE;


--
-- Name: sta_repo_hook fk_sta_repo_hook_repo; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_hook
    ADD CONSTRAINT fk_sta_repo_hook_repo FOREIGN KEY (repository_id) REFERENCES public.repository(id) ON DELETE CASCADE;


--
-- Name: sta_repo_origin fk_sta_repo_origin_origin_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_origin
    ADD CONSTRAINT fk_sta_repo_origin_origin_id FOREIGN KEY (origin_id) REFERENCES public.repository(id);


--
-- Name: sta_repo_origin fk_sta_repo_origin_repo_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_origin
    ADD CONSTRAINT fk_sta_repo_origin_repo_id FOREIGN KEY (repository_id) REFERENCES public.repository(id) ON DELETE CASCADE;


--
-- Name: sta_repo_push_activity fk_sta_repo_push_activity_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_push_activity
    ADD CONSTRAINT fk_sta_repo_push_activity_id FOREIGN KEY (activity_id) REFERENCES public.sta_repo_activity(activity_id) ON DELETE CASCADE;


--
-- Name: sta_repo_push_ref fk_sta_repo_push_ref_act_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repo_push_ref
    ADD CONSTRAINT fk_sta_repo_push_ref_act_id FOREIGN KEY (activity_id) REFERENCES public.sta_repo_push_activity(activity_id) ON DELETE CASCADE;


--
-- Name: sta_repository_scoped_id fk_sta_repo_scoped_id_repo; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_repository_scoped_id
    ADD CONSTRAINT fk_sta_repo_scoped_id_repo FOREIGN KEY (repository_id) REFERENCES public.repository(id) ON DELETE CASCADE;


--
-- Name: sta_service_user fk_sta_service_user_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_service_user
    ADD CONSTRAINT fk_sta_service_user_id FOREIGN KEY (user_id) REFERENCES public.stash_user(id) ON DELETE CASCADE;


--
-- Name: sta_user_settings fk_sta_user_settings_lob; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_user_settings
    ADD CONSTRAINT fk_sta_user_settings_lob FOREIGN KEY (lob_id) REFERENCES public.sta_shared_lob(id);


--
-- Name: sta_user_settings fk_sta_user_settings_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_user_settings
    ADD CONSTRAINT fk_sta_user_settings_user FOREIGN KEY (id) REFERENCES public.stash_user(id);


--
-- Name: sta_watcher fk_sta_watcher_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.sta_watcher
    ADD CONSTRAINT fk_sta_watcher_user FOREIGN KEY (user_id) REFERENCES public.stash_user(id);


--
-- Name: trusted_app_restriction fk_trusted_app; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.trusted_app_restriction
    ADD CONSTRAINT fk_trusted_app FOREIGN KEY (trusted_app_id) REFERENCES public.trusted_app(id) ON DELETE CASCADE;


--
-- Name: cwd_user_attribute fk_user_attr_dir_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_user_attribute
    ADD CONSTRAINT fk_user_attr_dir_id FOREIGN KEY (directory_id) REFERENCES public.cwd_directory(id);


--
-- Name: cwd_user_attribute fk_user_attribute_id_user_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_user_attribute
    ADD CONSTRAINT fk_user_attribute_id_user_id FOREIGN KEY (user_id) REFERENCES public.cwd_user(id);


--
-- Name: cwd_user_credential_record fk_user_cred_user; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_user_credential_record
    ADD CONSTRAINT fk_user_cred_user FOREIGN KEY (user_id) REFERENCES public.cwd_user(id);


--
-- Name: cwd_user fk_user_dir_id; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_user
    ADD CONSTRAINT fk_user_dir_id FOREIGN KEY (directory_id) REFERENCES public.cwd_directory(id);


--
-- Name: cwd_group_admin_user fk_user_target_group; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_group_admin_user
    ADD CONSTRAINT fk_user_target_group FOREIGN KEY (target_group_id) REFERENCES public.cwd_group(id) ON DELETE CASCADE;


--
-- Name: cwd_webhook fk_webhook_app; Type: FK CONSTRAINT; Schema: public; Owner: bitbucketuser
--

ALTER TABLE ONLY public.cwd_webhook
    ADD CONSTRAINT fk_webhook_app FOREIGN KEY (application_id) REFERENCES public.cwd_application(id);


--
-- PostgreSQL database dump complete
--

