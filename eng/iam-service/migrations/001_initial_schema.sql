-- Enable UUID extension (NeonDB usually has this, but good to be explicit)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- We define these types first so we can use them in tables.
CREATE TYPE tenant_status_type AS ENUM ('ACTIVE', 'SUSPENDED', 'PENDING_SETUP');
CREATE TYPE user_status_type AS ENUM ('INVITED', 'ACTIVE', 'SUSPENDED', 'LOCKED_OUT');
CREATE TYPE key_status_type AS ENUM ('ACTIVE', 'PASSIVE', 'RETIRED');

-- ==========================================
-- 1. TENANTS (The Schools)
-- ==========================================
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    domain TEXT UNIQUE, -- e.g., "greenwood.schoolapp.com"
    status tenant_status_type NOT NULL DEFAULT 'PENDING_SETUP',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==========================================
-- 2. USERS (The Identity)
-- ==========================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL, -- Argon2id hash

    -- "Amazon Style" Status: INVITED -> ACTIVE -> SUSPENDED
    status user_status_type NOT NULL DEFAULT 'INVITED',

    -- Security features
    is_mfa_enabled BOOLEAN DEFAULT FALSE,
    mfa_secret_encrypted TEXT, -- Store encrypted TOTP secret

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ -- Soft Delete column
);

-- CRITICAL: Unique constraint that ignores soft-deleted users.
-- This solves the "Zombie Email" problem we discussed.
CREATE UNIQUE INDEX idx_users_email_tenant
ON users(tenant_id, email)
WHERE deleted_at IS NULL;

-- ==========================================
-- 3. RBAC (Roles & Permissions)
-- ==========================================
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE, -- e.g., "USER_CREATE", "MARKS_VIEW"
    description TEXT
);

CREATE TABLE user_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    name TEXT NOT NULL, -- e.g., "Math Teacher", "Principal"
    description TEXT,
    is_system_role BOOLEAN DEFAULT FALSE, -- If TRUE, Admins cannot delete this role
    created_at TIMESTAMPTZ DEFAULT NOW(),

    -- Prevent duplicate group names within the same school
    UNIQUE(tenant_id, name)
);

-- Mapping: Which Role has which Permission?
CREATE TABLE group_permissions (
    group_id UUID NOT NULL REFERENCES user_groups(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (group_id, permission_id)
);

-- Assignment: Which User has which Role?
CREATE TABLE group_members (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES user_groups(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, group_id)
);

-- ==========================================
-- 4. SESSIONS (Refresh Tokens)
-- ==========================================
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Security: We store the HASH of the token, not the token itself.
    refresh_token_hash TEXT NOT NULL UNIQUE,

    client_ip INET,
    user_agent TEXT,
    is_revoked BOOLEAN DEFAULT FALSE, -- Emergency Kill Switch

    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index for fast token lookups during refresh flow
CREATE INDEX idx_sessions_token_hash ON sessions(refresh_token_hash);

-- ==========================================
-- 5. KEY MANAGEMENT (The "Netflix" Rotation)
-- ==========================================
CREATE TABLE jwk_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kid TEXT NOT NULL UNIQUE, -- Key ID (e.g., "2026-01-15-v1")

    -- The Public Key is safe to show.
    public_key TEXT NOT NULL,

    -- The Private Key is ENCRYPTED with the Master Key (from Env Vars).
    encrypted_private_key TEXT NOT NULL,

    algorithm VARCHAR(10) DEFAULT 'EdDSA',
    status key_status_type NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, PASSIVE (Drain), RETIRED

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ -- When this key should stop being used entirely
);

-- ==========================================
-- 6. AUDIT LOGS (Immutable History)
-- ==========================================
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trace_id TEXT NOT NULL, -- The Request ID from Middleware

    tenant_id UUID NOT NULL, -- No Foreign Key (keep log even if tenant deleted)
    actor_id UUID,           -- Who did it? (Nullable for system actions)

    action VARCHAR(100) NOT NULL, -- "USER_LOGIN", "ROLE_CREATED"
    resource_id UUID,             -- Which user/role was touched?
    metadata JSONB,               -- Extra details {"ip": "1.2.3.4"}

    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Performance: Indexes for the Admin Dashboard "Activity Log"
CREATE INDEX idx_audit_tenant_time ON audit_logs(tenant_id, created_at DESC);
