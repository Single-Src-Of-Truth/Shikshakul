package com.modulixinfotech.iams.shikshakul.mapper.request;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.UUID;

public class AuthRequest {

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class RegisterRequest {

        @NotBlank(message = "Invite token is required")
        private UUID inviteToken;

        @NotBlank(message = "First name is required")
        @Size(min = 2, max = 100, message = "First name must be between 2 and 100 characters")
        private String firstName;

        @NotBlank(message = "Last name is required")
        @Size(min = 2, max = 100, message = "Last name must be between 2 and 100 characters")
        private String lastName;

        @NotBlank(message = "Email is required")
        @Email(message = "Invalid email format")
        private String email;

        @Pattern(regexp = "^[0-9]{10}$", message = "Mobile must be 10 digits")
        private String mobile;

        @NotBlank(message = "Password is required")
        @Size(min = 8, max = 100, message = "Password must be between 8 and 100 characters")
        @Pattern(regexp = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]{8,}$",
                message = "Password must contain at least one uppercase, one lowercase, one digit, and one special character")
        private String password;

        @NotBlank(message = "Confirm password is required")
        private String confirmPassword;
    }

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class LoginRequest {

        @NotBlank(message = "Email is required")
        @Email(message = "Invalid email format")
        private String email;

        @NotBlank(message = "Password is required")
        private String password;

        private String tenantDomain;
    }

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class ChangePasswordRequest {

        @NotBlank(message = "Current password is required")
        private String currentPassword;

        @NotBlank(message = "New password is required")
        @Size(min = 8, max = 100, message = "Password must be between 8 and 100 characters")
        @Pattern(regexp = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]{8,}$",
                message = "Password must contain at least one uppercase, one lowercase, one digit, and one special character")
        private String newPassword;

        @NotBlank(message = "Confirm password is required")
        private String confirmPassword;
    }

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class CreateInviteRequest {

        @NotBlank(message = "Email is required")
        @Email(message = "Invalid email format")
        private String email;

        @Pattern(regexp = "^[0-9]{10}$", message = "Mobile must be 10 digits")
        private String mobile;

        @NotBlank(message = "Role is required")
        private String role;
    }

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class CreateTenantRequest {

        @NotBlank(message = "Tenant name is required")
        @Size(min = 3, max = 200, message = "Tenant name must be between 3 and 200 characters")
        private String name;

        @NotBlank(message = "Domain is required")
        @Pattern(regexp = "^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$",
                message = "Invalid domain format")
        private String domain;

        @Size(max = 500, message = "Description must not exceed 500 characters")
        private String description;

        @Email(message = "Invalid email format")
        private String contactEmail;

        @Pattern(regexp = "^[0-9]{10}$", message = "Contact phone must be 10 digits")
        private String contactPhone;

        @NotBlank(message = "Admin email is required")
        @Email(message = "Invalid admin email format")
        private String adminEmail;

        @NotBlank(message = "Admin password is required")
        @Size(min = 8, max = 100, message = "Admin password must be between 8 and 100 characters")
        private String adminPassword;

        @NotBlank(message = "Admin first name is required")
        private String adminFirstName;

        @NotBlank(message = "Admin last name is required")
        private String adminLastName;
    }
}