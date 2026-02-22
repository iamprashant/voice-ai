CREATE TABLE public.call_contexts (
    id bigint PRIMARY KEY,
    context_id character varying(36) NOT NULL,
    status character varying(20) DEFAULT 'pending' NOT NULL,
    assistant_id bigint NOT NULL,
    conversation_id bigint NOT NULL,
    project_id bigint NOT NULL DEFAULT 0,
    organization_id bigint NOT NULL DEFAULT 0,
    auth_token text NOT NULL DEFAULT '',
    auth_type character varying(50) NOT NULL DEFAULT '',
    provider character varying(50) NOT NULL DEFAULT '',
    direction character varying(20) NOT NULL DEFAULT '',
    caller_number character varying(50) NOT NULL DEFAULT '',
    callee_number character varying(50) NOT NULL DEFAULT '',
    from_number character varying(50) NOT NULL DEFAULT '',
    assistant_provider_id bigint NOT NULL DEFAULT 0,
    channel_uuid character varying(200) NOT NULL DEFAULT '',
    created_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_date timestamp without time zone
);

CREATE UNIQUE INDEX call_contexts_context_id_idx ON public.call_contexts (context_id);
CREATE INDEX call_contexts_status_idx ON public.call_contexts (status);
CREATE INDEX call_contexts_assistant_id_idx ON public.call_contexts (assistant_id);
CREATE INDEX call_contexts_conversation_id_idx ON public.call_contexts (conversation_id);
CREATE INDEX call_contexts_created_date_idx ON public.call_contexts (created_date);
