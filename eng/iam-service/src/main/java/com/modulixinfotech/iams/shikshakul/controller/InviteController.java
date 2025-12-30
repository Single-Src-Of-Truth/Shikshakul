package com.modulixinfotech.iams.shikshakul.controller;

import com.modulixinfotech.iams.shikshakul.domain.model.Invite;
import com.modulixinfotech.iams.shikshakul.domain.model.Role;
import com.modulixinfotech.iams.shikshakul.domain.model.Tenant;
import com.modulixinfotech.iams.shikshakul.mapper.InviteDTO;
import com.modulixinfotech.iams.shikshakul.mapper.common.ApiResponse;
import com.modulixinfotech.iams.shikshakul.security.userdetails.CustomUserDetails;
import com.modulixinfotech.iams.shikshakul.service.InviteService;
import com.modulixinfotech.iams.shikshakul.service.RoleService;
import com.modulixinfotech.iams.shikshakul.service.TenantService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

import static com.modulixinfotech.iams.shikshakul.mapper.request.AuthRequest.CreateInviteRequest;
import static com.modulixinfotech.iams.shikshakul.mapper.response.AuthResponse.InviteResponse;

@RestController
@RequestMapping("/api/invites")
@RequiredArgsConstructor
@Slf4j
public class InviteController {

    private final InviteService inviteService;
    private final RoleService roleService;
    private final TenantService tenantService;

    @Value("${app.oauth2.issuer-url}")
    private String issuerUrl;

    @PostMapping
    public ResponseEntity<ApiResponse<InviteResponse>> createInvite(
            @Valid @RequestBody CreateInviteRequest request,
            @AuthenticationPrincipal CustomUserDetails userDetails) {

        log.info("Creating invite for email: {}", request.getEmail());

        Role role = roleService.getRoleByName(Role.RoleName.valueOf(request.getRole().toUpperCase()));
        Tenant tenant = userDetails.getUser().getTenant();

        Invite invite = inviteService.createInvite(
                tenant,
                request.getEmail(),
                request.getMobile(),
                role,
                userDetails.getUser()
        );

        String inviteUrl = issuerUrl + "/register?token=" + invite.getToken();

        InviteResponse response = InviteResponse.builder()
                .inviteId(invite.getId())
                .token(invite.getToken())
                .email(invite.getEmail())
                .mobile(invite.getMobile())
                .role(invite.getRole().getName().name())
                .expiresAt(invite.getExpiresAt())
                .inviteUrl(inviteUrl)
                .build();

        return ResponseEntity
                .status(HttpStatus.CREATED)
                .body(ApiResponse.success("Invite created successfully", response));
    }

    @GetMapping
    public ResponseEntity<ApiResponse<List<InviteDTO>>> getInvites(
            @AuthenticationPrincipal CustomUserDetails userDetails) {

        log.info("Fetching invites for tenant: {}", userDetails.getUser().getTenant().getName());

        List<InviteDTO> invites = inviteService.getInvitesByTenant(userDetails.getUser().getTenant())
                .stream()
                .map(InviteDTO::fromEntity)
                .collect(Collectors.toList());

        return ResponseEntity.ok(ApiResponse.success(invites));
    }

    @GetMapping("/pending")
    public ResponseEntity<ApiResponse<List<InviteDTO>>> getPendingInvites(
            @AuthenticationPrincipal CustomUserDetails userDetails) {

        log.info("Fetching pending invites for tenant: {}", userDetails.getUser().getTenant().getName());

        List<InviteDTO> invites = inviteService.getPendingInvitesByTenant(userDetails.getUser().getTenant())
                .stream()
                .map(InviteDTO::fromEntity)
                .collect(Collectors.toList());

        return ResponseEntity.ok(ApiResponse.success(invites));
    }

    @DeleteMapping("/{inviteId}")
    public ResponseEntity<ApiResponse<Void>> revokeInvite(
            @PathVariable UUID inviteId,
            @AuthenticationPrincipal CustomUserDetails userDetails) {

        log.info("Revoking invite: {}", inviteId);

        inviteService.revokeInvite(inviteId, userDetails.getUser());

        return ResponseEntity.ok(ApiResponse.success("Invite revoked successfully", null));
    }
}