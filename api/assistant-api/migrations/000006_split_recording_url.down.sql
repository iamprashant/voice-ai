ALTER TABLE public.assistant_conversation_recordings
    ADD COLUMN recording_url character varying(200);

UPDATE public.assistant_conversation_recordings
    SET recording_url = assistant_recording_url;

ALTER TABLE public.assistant_conversation_recordings
    ALTER COLUMN recording_url SET NOT NULL;

ALTER TABLE public.assistant_conversation_recordings
    DROP COLUMN assistant_recording_url,
    DROP COLUMN user_recording_url;
