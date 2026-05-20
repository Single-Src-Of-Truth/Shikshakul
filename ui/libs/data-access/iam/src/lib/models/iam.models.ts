// ─── Generic Response Wrapper ───────────────────────────────────────────────

export interface IamApiResponse<T> {
    data: T;
    message?: string;
    success?: boolean;
}

// ─── User & Profile ──────────────────────────────────────────────────────────

export type UserStatus = 'ACTIVE' | 'SUSPENDED' | 'PENDING';
export type IdentifierType = 'EMAIL' | 'MOBILE';

export interface UserProfile {
    id: string;
    tenant_id?: string | null;
    first_name: string;
    last_name: string;
    identifier: string;
    identifier_type: IdentifierType;
    status: UserStatus;
    roles: string[];
    permissions: string[];
}

export interface Session {
    ip_address: string;
    browser: string;
    os: string;
    device_type: string;
    location: string;
    last_active_at: string;
    is_current: boolean;
}

// ─── Tenant (Super Admin only) ───────────────────────────────────────────────

export interface Tenant {
    id: string;
    name: string;
    domain: string;
    is_active: boolean;
}

// ─── Roles & Permissions ─────────────────────────────────────────────────────

export interface Role {
    id: string;
    name: string;
    description: string;
    is_system: boolean;
    permissions?: string[];
}

export interface Permission {
    id: string;
    description: string;
}

// ─── Auth Request/Response shapes ────────────────────────────────────────────

export interface LoginRequest {
    tenant_id?: string | null;
    identifier: string;
    password: string;
}

export interface LoginResponse {
    user_id: string;
    first_name: string;
    status: UserStatus;
    redirect_command: string | null;
}

export interface ForgotPasswordRequest {
    tenant_id?: string | null;
    identifier: string;
}

export interface ResetPasswordRequest {
    tenant_id?: string | null;
    identifier: string;
    /** 6-character OTP sent via email/SMS */
    otp: string;
    new_password: string;
}

export interface ChangePasswordRequest {
    old_password: string;
    new_password: string;
}

// ─── Profile Request shapes ───────────────────────────────────────────────────

export interface UpdateProfileRequest {
    first_name: string;
    last_name?: string;
}

// ─── Tenant Request shapes (Super Admin) ─────────────────────────────────────

export interface CreateTenantRequest {
    name: string;
    domain: string;
    auth_config?: string;
}

export interface UpdateTenantRequest {
    name?: string;
    domain?: string;
    auth_config?: string;
    is_active?: boolean;
}

// ─── RBAC Request shapes ──────────────────────────────────────────────────────

export interface CreateRoleRequest {
    name: string;
    description?: string;
    permission_ids: string[];
}

export interface UpdateRoleRequest {
    name?: string;
    description?: string;
    permission_ids?: string[];
}

// ─── User Management Request shapes ──────────────────────────────────────────

export interface UpdateUserStatusRequest {
    status: 'ACTIVE' | 'SUSPENDED';
}

export interface ChangeUserRoleRequest {
    role_id: string;
}

// ─── Invite Request/Response shapes ──────────────────────────────────────────

export interface GenerateInviteRequest {
    tenant_id: string;
    role_id: string;
    /** Email or mobile number to invite */
    identifier: string;
}

export interface GenerateInviteResponse {
    invite_id: string;
    invite_token: string;
    identifier: string;
}

export interface AcceptInviteRequest {
    invite_token: string;
    password: string;
    first_name: string;
    last_name?: string;
}

export interface AcceptInviteResponse {
    user_id: string;
    identifier: string;
}
