package com.modulixinfotech.iams.shikshakul.mapper.internal;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.Set;
import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class JwtClaimsDTO {
    private UUID userId;
    private UUID tenantId;
    private String email;
    private Set<String> roles;
    private Set<String> scopes;
    private String audience;
    private String issuer;
    private Long issuedAt;
    private Long expiresAt;
    private String jti;
}