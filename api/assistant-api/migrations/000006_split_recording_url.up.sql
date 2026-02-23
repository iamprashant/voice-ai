ALTER TABLE public.assistant_conversation_recordings
    ADD COLUMN assistant_recording_url character varying(500),
    ADD COLUMN user_recording_url character varying(500);

UPDATE public.assistant_conversation_recordings
    SET assistant_recording_url = recording_url,
        user_recording_url = recording_url;

ALTER TABLE public.assistant_conversation_recordings
    ALTER COLUMN assistant_recording_url SET NOT NULL,
    ALTER COLUMN user_recording_url SET NOT NULL;

ALTER TABLE public.assistant_conversation_recordings
    DROP COLUMN recording_url;
