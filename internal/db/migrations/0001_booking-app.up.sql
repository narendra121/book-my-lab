-- migrate: disable-transaction
-- ==========================================================
-- ENUMS, FUNCTIONS, AND TABLE CREATION
-- ==========================================================

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'property_type_enum') THEN
        CREATE TYPE property_type_enum AS ENUM ('HOUSE', 'APARTMENT', 'CONDO');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'property_status_enum') THEN
        CREATE TYPE property_status_enum AS ENUM ('LISTED', 'UNLISTED', 'BOOKED');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'visit_status_enum') THEN
        CREATE TYPE visit_status_enum AS ENUM (
            'PENDING', 'ACCEPTED', 'REJECTED', 'RESCHEDULED', 'COMPLETED', 'CANCELLED'
        );
    END IF;
END$$;


-- ==========================================================
-- TRIGGER FUNCTION FOR AUTO-UPDATING "updated_at"
-- ==========================================================
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- ==========================================================
-- USERS TABLE (email as PK)
-- ==========================================================
CREATE TABLE IF NOT EXISTS users (
    email               VARCHAR(100) PRIMARY KEY,
    first_name          VARCHAR(100) NOT NULL,
    last_name           VARCHAR(100) NOT NULL,
    phone               VARCHAR(20) UNIQUE,
    password_hash       VARCHAR(255) NOT NULL,
    salt                VARCHAR(100) NOT NULL,
    profile_pic_url     TEXT,
    address             VARCHAR(255),
    role                VARCHAR(100) NOT NULL, -- 'admin', 'partner', 'buyer'
    is_email_verified   BOOLEAN DEFAULT false,
    is_phone_verified   BOOLEAN DEFAULT false,
    refresh_token       VARCHAR(255),
    deleted             BOOLEAN DEFAULT false,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER users_update_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();


-- ==========================================================
-- PROPERTIES TABLE
-- ==========================================================
CREATE TABLE IF NOT EXISTS properties (
    id SERIAL PRIMARY KEY,
    partner_email      VARCHAR(100) NOT NULL,
    title              VARCHAR(200) NOT NULL,
    description        TEXT,
    property_type      property_type_enum NOT NULL,
    bedrooms           INT,
    bathrooms          INT,
    area_sqft          NUMERIC(10,2),
    price              NUMERIC(12,2),
    city               VARCHAR(100),
    state              VARCHAR(100),
    address            TEXT,
    status             property_status_enum DEFAULT 'LISTED',
    created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_property_owner FOREIGN KEY (partner_email)
        REFERENCES users(email) ON DELETE CASCADE
);

CREATE TRIGGER properties_update_timestamp
BEFORE UPDATE ON properties
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();


-- ==========================================================
-- PROPERTY_PHOTOS TABLE
-- ==========================================================
CREATE TABLE IF NOT EXISTS property_photos (
    id SERIAL PRIMARY KEY,
    property_id   INT NOT NULL,
    image_url     TEXT NOT NULL,
    is_primary    BOOLEAN DEFAULT FALSE,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_photo_property FOREIGN KEY (property_id)
        REFERENCES properties(id) ON DELETE CASCADE
);


-- ==========================================================
-- FAVORITES TABLE
-- ==========================================================
CREATE TABLE IF NOT EXISTS favorites (
    id SERIAL PRIMARY KEY,
    user_email     VARCHAR(100) NOT NULL,
    property_id    INT NOT NULL,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_favorites UNIQUE (user_email, property_id),
    CONSTRAINT fk_favorite_user FOREIGN KEY (user_email)
        REFERENCES users(email) ON DELETE CASCADE,
    CONSTRAINT fk_favorite_property FOREIGN KEY (property_id)
        REFERENCES properties(id) ON DELETE CASCADE
);


-- ==========================================================
-- VISITS TABLE
-- ==========================================================
CREATE TABLE IF NOT EXISTS visits (
    id SERIAL PRIMARY KEY,
    property_id       INT NOT NULL,
    buyer_email       VARCHAR(100) NOT NULL,
    scheduled_time    TIMESTAMP NOT NULL,
    status            visit_status_enum DEFAULT 'PENDING',
    reschedule_time   TIMESTAMP NULL,
    partner_note      TEXT,
    buyer_note        TEXT,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_visit_property FOREIGN KEY (property_id)
        REFERENCES properties(id) ON DELETE CASCADE,
    CONSTRAINT fk_visit_buyer FOREIGN KEY (buyer_email)
        REFERENCES users(email) ON DELETE CASCADE
);

CREATE TRIGGER visits_update_timestamp
BEFORE UPDATE ON visits
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();


-- ==========================================================
-- RATINGS TABLE
-- ==========================================================
CREATE TABLE IF NOT EXISTS ratings (
    id SERIAL PRIMARY KEY,
    property_id     INT NOT NULL,
    buyer_email     VARCHAR(100) NOT NULL,
    rating          INT CHECK (rating BETWEEN 1 AND 5),
    review_text     TEXT,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ratings_property FOREIGN KEY (property_id)
        REFERENCES properties(id) ON DELETE CASCADE,
    CONSTRAINT fk_ratings_buyer FOREIGN KEY (buyer_email)
        REFERENCES users(email) ON DELETE CASCADE
);
