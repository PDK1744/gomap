CREATE TABLE branches (
    tru_id TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO branches (tru_id, name)
VALUES
    ('sl_branch', 'Stoneledge Branch'),
    ('rr_branch', 'Richmond Road Branch'),
    ('atl_branch', 'Atlanta Branch'),
    ('nb_branch', 'New Boston Branch');

CREATE TABLE object_metadata (
    tru_id TEXT PRIMARY KEY,
    branch_tru_id TEXT NOT NULL REFERENCES branches(tru_id),
    name TEXT NOT NULL,
    object_type TEXT NOT NULL,
    assignable BOOLEAN NOT NULL DEFAULT TRUE,
    is_generator BOOLEAN NOT NULL DEFAULT FALSE,
    parent TEXT
);


CREATE TABLE assets (
    asset_id INT PRIMARY KEY,
    name TEXT NOT NULL,
    asset_type TEXT,
    branch TEXT,
    branch_id TEXT REFERENCES branches(tru_id)
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_seen_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE branch_aliases (
    id SERIAL PRIMARY KEY,
    branch_id TEXT NOT NULL
        REFERENCES branches(tru_id)
        ON DELETE CASCADE,
    alias TEXT NOT NULL
);

CREATE UNIQUE INDEX ux_branch_aliases_alias
ON branch_aliases (LOWER(alias));

CREATE TABLE object_assets (
    object_tru_id TEXT NOT NULL
        REFERENCES object_metadata(tru_id),

    asset_id TEXT NOT NULL
        REFERENCES assets(asset_id),

    assigned_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (object_tru_id, asset_id)
);