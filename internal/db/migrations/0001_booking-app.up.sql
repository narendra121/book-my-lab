-- migrate: disable-transaction
-- ==========================================================
-- ENUMS, FUNCTIONS, AND TABLE CREATION
-- ==========================================================

-- ==========================================================
-- FUNCTION: Auto-update "updated_at" timestamps
-- ==========================================================
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- ==========================================================
-- USERS TABLE
-- ==========================================================
CREATE TABLE IF NOT EXISTS users (
    username            VARCHAR(50) PRIMARY KEY,
    first_name          VARCHAR(100) NOT NULL,
    last_name           VARCHAR(100) NOT NULL,
    email               VARCHAR(100) UNIQUE NOT NULL,
    phone               VARCHAR(20) UNIQUE,
    password_hash       VARCHAR(255) NOT NULL,
    salt                VARCHAR(100) NOT NULL,
    profile_pic_url     TEXT,
    address             VARCHAR(255),
    role                VARCHAR(100) NOT NULL,  -- 'admin', 'partner', 'buyer'
    is_email_verified   BOOLEAN DEFAULT false,
    is_phone_verified   BOOLEAN DEFAULT false,
    refresh_token       VARCHAR(255),
    rating              NUMERIC(3,2) DEFAULT 0,  -- partner average rating
    deleted             BOOLEAN DEFAULT false,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER users_update_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();


-- ==========================================================
-- FUNCTION: Generate unique, random-suffixed username
-- ==========================================================
CREATE OR REPLACE FUNCTION generate_unique_username()
RETURNS TRIGGER AS $$
DECLARE
    base_first      VARCHAR(50);
    base_last       VARCHAR(50);
    random_suffix   VARCHAR(3);
    candidate       VARCHAR(50);
BEGIN
    IF NEW.username IS NULL OR NEW.username = '' THEN
        base_first := lower(substring(NEW.first_name FROM 1 FOR 4));
        base_last  := lower(substring(NEW.last_name FROM 1 FOR 4));

        LOOP
            random_suffix := substring(md5(random()::text || clock_timestamp()::text) FROM 1 FOR 3);
            candidate := base_first || '_' || base_last || '_' || random_suffix;
            EXIT WHEN NOT EXISTS (SELECT 1 FROM users WHERE username = candidate);
        END LOOP;

        NEW.username := candidate;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- ==========================================================
-- TRIGGER: Assign username before insert
-- ==========================================================
DROP TRIGGER IF EXISTS users_generate_unique_username ON users;

CREATE TRIGGER users_generate_unique_username
BEFORE INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION generate_unique_username();


-- ==========================================================
-- PROPERTIES TABLE
-- ==========================================================
CREATE TABLE IF NOT EXISTS properties (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    partner_username   VARCHAR(50) NOT NULL,
    title              VARCHAR(200) NOT NULL,
    description        TEXT,
    property_type      VARCHAR(50) NOT NULL,
    bedrooms           INT,
    bathrooms          INT,
    area_sqft          NUMERIC(10,2),
    price              NUMERIC(12,2),
    city               VARCHAR(100),
    state              VARCHAR(100),
    address            TEXT,
    status             VARCHAR(50) NOT NULL,
    deleted            BOOLEAN DEFAULT false,
    created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_property_owner FOREIGN KEY (partner_username)
        REFERENCES users(username) ON DELETE CASCADE
);

CREATE TRIGGER properties_update_timestamp
BEFORE UPDATE ON properties
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();


-- ==========================================================
-- PROPERTY_PHOTOS TABLE
-- ==========================================================
CREATE TABLE IF NOT EXISTS property_photos (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    property_id   BIGINT NOT NULL,
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
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_username   VARCHAR(50) NOT NULL,
    property_id     BIGINT NOT NULL,
    deleted         BOOLEAN DEFAULT false,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_favorites UNIQUE (user_username, property_id),
    CONSTRAINT fk_favorite_user FOREIGN KEY (user_username)
        REFERENCES users(username) ON DELETE CASCADE,
    CONSTRAINT fk_favorite_property FOREIGN KEY (property_id)
        REFERENCES properties(id) ON DELETE CASCADE
);


-- ==========================================================
-- VISITS TABLE
-- ==========================================================
CREATE TABLE IF NOT EXISTS visits (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    property_id       BIGINT NOT NULL,
    buyer_username    VARCHAR(50) NOT NULL,
    scheduled_time    TIMESTAMP NOT NULL,
    status            VARCHAR(100),
    reschedule_time   TIMESTAMP NULL,
    partner_note      TEXT,
    buyer_note        TEXT,
    deleted           BOOLEAN DEFAULT false,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_visit_property FOREIGN KEY (property_id)
        REFERENCES properties(id) ON DELETE CASCADE,
    CONSTRAINT fk_visit_buyer FOREIGN KEY (buyer_username)
        REFERENCES users(username) ON DELETE CASCADE
);

CREATE TRIGGER visits_update_timestamp
BEFORE UPDATE ON visits
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();


-- ==========================================================
-- RATINGS TABLE (Buyer → Partner for a Property)
-- ==========================================================
CREATE TABLE IF NOT EXISTS ratings (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    property_id      BIGINT NOT NULL,
    buyer_username   VARCHAR(50) NOT NULL,
    partner_username VARCHAR(50) NOT NULL,
    rating           INT CHECK (rating BETWEEN 1 AND 5),
    review_text      TEXT,
    deleted          BOOLEAN DEFAULT false,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ratings_property FOREIGN KEY (property_id)
        REFERENCES properties(id) ON DELETE CASCADE,
    CONSTRAINT fk_ratings_buyer FOREIGN KEY (buyer_username)
        REFERENCES users(username) ON DELETE CASCADE,
    CONSTRAINT fk_ratings_partner FOREIGN KEY (partner_username)
        REFERENCES users(username) ON DELETE CASCADE
);


-- ==========================================================
-- FUNCTION: Update partner's average rating
-- ==========================================================
CREATE OR REPLACE FUNCTION update_partner_avg_rating()
RETURNS TRIGGER AS $$
DECLARE
    avg_rating NUMERIC(3,2);
    target_partner VARCHAR(50);
BEGIN
    -- Identify the partner for update
    IF (TG_OP = 'DELETE') THEN
        target_partner := OLD.partner_username;
    ELSE
        target_partner := NEW.partner_username;
    END IF;

    -- Compute average from non-deleted ratings
    SELECT ROUND(AVG(rating)::numeric, 2)
    INTO avg_rating
    FROM ratings
    WHERE partner_username = target_partner AND deleted = false;

    -- Default to 0 if no ratings remain
    IF avg_rating IS NULL THEN
        avg_rating := 0;
    END IF;

    -- Update user's rating
    UPDATE users
    SET rating = avg_rating
    WHERE username = target_partner;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;


-- ==========================================================
-- TRIGGERS: Auto-update partner rating when ratings table changes
-- ==========================================================
DROP TRIGGER IF EXISTS ratings_insert_partner_avg_rating ON ratings;
DROP TRIGGER IF EXISTS ratings_update_partner_avg_rating ON ratings;
DROP TRIGGER IF EXISTS ratings_delete_partner_avg_rating ON ratings;

CREATE TRIGGER ratings_insert_partner_avg_rating
AFTER INSERT ON ratings
FOR EACH ROW
EXECUTE FUNCTION update_partner_avg_rating();

CREATE TRIGGER ratings_update_partner_avg_rating
AFTER UPDATE ON ratings
FOR EACH ROW
EXECUTE FUNCTION update_partner_avg_rating();

CREATE TRIGGER ratings_delete_partner_avg_rating
AFTER DELETE ON ratings
FOR EACH ROW
EXECUTE FUNCTION update_partner_avg_rating();
