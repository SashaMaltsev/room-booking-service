BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NULL,
    role TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    
    CONSTRAINT users_email_not_blank_chk CHECK (BTRIM(email) <> ''),
    CONSTRAINT users_email_lowercase_chk CHECK (email = LOWER(email)),
    CONSTRAINT users_role_chk CHECK (role IN ('admin', 'user'))
);

CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT NULL,
    capacity INTEGER NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT rooms_name_not_blank_chk CHECK (BTRIM(name) <> ''),
    CONSTRAINT rooms_capacity_positive_chk CHECK (capacity IS NULL OR capacity > 0)
);

CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL,
    days_of_week SMALLINT[] NOT NULL,
    start_time_utc TIME NOT NULL,
    end_time_utc TIME NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT schedules_room_fk 
        FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    CONSTRAINT schedules_room_unique 
        UNIQUE (room_id),
    CONSTRAINT schedules_days_not_empty_chk 
        CHECK (cardinality(days_of_week) > 0),
    CONSTRAINT schedules_days_without_nulls_chk 
        CHECK (array_position(days_of_week, NULL) IS NULL),
    CONSTRAINT schedules_days_allowed_chk 
        CHECK (days_of_week <@ ARRAY[1, 2, 3, 4, 5, 6, 7]::SMALLINT[]),
    CONSTRAINT schedules_time_range_chk 
        CHECK (start_time_utc < end_time_utc),
    CONSTRAINT schedules_slot_alignment_chk  -- Длительность слота должна быть кратна 30 минутам (1800 секунд). 30 мин, 1 час, 1.5 часа, 2 часа...
        CHECK (
            MOD(EXTRACT(EPOCH FROM (end_time_utc - start_time_utc))::INTEGER, 1800) = 0
        )
);

CREATE TABLE slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL,
    schedule_id UUID NOT NULL,
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT slots_room_fk
        FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    
    CONSTRAINT slots_schedule_fk
        FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE,
    
    CONSTRAINT slots_time_range_chk 
        CHECK (start_at < end_at),
    
    CONSTRAINT slots_duration_chk 
        CHECK (end_at = start_at + INTERVAL '30 minutes'),  -- все слоты в системе имеют одинаковую длину
   
    CONSTRAINT slots_exact_interval_unique 
        UNIQUE (room_id, start_at, end_at),
    
    CONSTRAINT slots_no_overlap  -- Запрещает пересекающиеся слоты в одной комнате
        EXCLUDE USING GIST (
            room_id WITH =,
            tstzrange(start_at, end_at, '[)') WITH &&
        )
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slot_id UUID NOT NULL,
    user_id UUID NOT NULL,
    status TEXT NOT NULL,
    conference_link TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    cancelled_at TIMESTAMPTZ NULL,

    CONSTRAINT bookings_slot_fk
        FOREIGN KEY (slot_id) REFERENCES slots(id) ON DELETE CASCADE,

    CONSTRAINT bookings_user_fk
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    
    CONSTRAINT bookings_status_chk 
        CHECK (status IN ('active', 'cancelled')),
    
    CONSTRAINT bookings_cancelled_state_chk 
        CHECK (
            (status = 'active' AND cancelled_at IS NULL)
            OR
            (status = 'cancelled' AND cancelled_at IS NOT NULL)
        )
);

CREATE UNIQUE INDEX bookings_one_active_per_slot_idx
    ON bookings (slot_id)
    WHERE status = 'active';

INSERT INTO users (id, email, role)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'dummy-admin@example.com', 'admin'),
    ('00000000-0000-0000-0000-000000000002', 'dummy-user@example.com', 'user');

COMMIT;
